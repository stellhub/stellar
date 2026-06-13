package interceptor

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Kind string

const (
	KindHTTPServer Kind = "http.server"
	KindHTTPClient Kind = "http.client"
	KindGRPCServer Kind = "grpc.server"
	KindGRPCClient Kind = "grpc.client"
)

type Stage string

const (
	StageRecovery       Stage = "recovery"
	StageRouteResolve   Stage = "route_resolve"
	StageObserve        Stage = "observe"
	StageDeadline       Stage = "deadline"
	StageAdmission      Stage = "admission"
	StageSecurity       Stage = "security"
	StageDecodeValidate Stage = "decode_validate"
	StageRetry          Stage = "retry"
	StageBusiness       Stage = "business"
)

var (
	ServerInboundStages = []Stage{
		StageRecovery,
		StageRouteResolve,
		StageObserve,
		StageDeadline,
		StageAdmission,
		StageSecurity,
		StageDecodeValidate,
		StageBusiness,
	}

	ServerInboundPreDecodeStages = []Stage{
		StageRecovery,
		StageRouteResolve,
		StageObserve,
		StageDeadline,
		StageAdmission,
		StageSecurity,
	}

	ServerInboundDecodeStages = []Stage{
		StageDecodeValidate,
	}

	ClientOutboundStages = []Stage{
		StageRecovery,
		StageObserve,
		StageDeadline,
		StageAdmission,
		StageRetry,
		StageBusiness,
		StageSecurity,
	}
)

type Invocation struct {
	Kind       Kind
	Protocol   string
	Service    string
	Operation  string
	Method     string
	Path       string
	Target     string
	RequestID  string
	Headers    Header
	Attributes map[string]any
	Raw        any
}

func (i *Invocation) SetAttribute(key string, value any) {
	if i == nil || strings.TrimSpace(key) == "" {
		return
	}
	if i.Attributes == nil {
		i.Attributes = map[string]any{}
	}
	i.Attributes[key] = value
}

func (i *Invocation) Attribute(key string) (any, bool) {
	if i == nil || i.Attributes == nil {
		return nil, false
	}
	value, ok := i.Attributes[key]
	return value, ok
}

type Header map[string][]string

func HeaderFromHTTP(header http.Header) Header {
	values := make(Header, len(header))
	for key, items := range header {
		values[key] = append([]string(nil), items...)
	}
	return values
}

func (h Header) Get(key string) string {
	if len(h) == 0 {
		return ""
	}
	for current, values := range h {
		if strings.EqualFold(current, key) && len(values) > 0 {
			return values[0]
		}
	}
	return ""
}

type Handler func(ctx context.Context, inv *Invocation, req any) (any, error)

type Interceptor interface {
	Name() string
	Intercept(ctx context.Context, inv *Invocation, req any, next Handler) (any, error)
}

type Func func(ctx context.Context, inv *Invocation, req any, next Handler) (any, error)

type namedInterceptor struct {
	name string
	fn   Func
}

func New(name string, fn Func) Interceptor {
	return namedInterceptor{name: strings.TrimSpace(name), fn: fn}
}

func (i namedInterceptor) Name() string {
	if i.name == "" {
		return "anonymous"
	}
	return i.name
}

func (i namedInterceptor) Intercept(ctx context.Context, inv *Invocation, req any, next Handler) (any, error) {
	if i.fn == nil {
		return next(ctx, inv, req)
	}
	return i.fn(ctx, inv, req, next)
}

type Definition struct {
	Name        string
	Kind        Kind
	Stage       Stage
	Order       int
	Framework   bool
	Interceptor Interceptor
}

func Framework(kind Kind, stage Stage, name string, it Interceptor) Definition {
	return Definition{
		Name:        strings.TrimSpace(name),
		Kind:        kind,
		Stage:       stage,
		Framework:   true,
		Interceptor: it,
	}
}

func Business(kind Kind, name string, order int, it Interceptor) Definition {
	return Definition{
		Name:        strings.TrimSpace(name),
		Kind:        kind,
		Stage:       StageBusiness,
		Order:       order,
		Framework:   false,
		Interceptor: it,
	}
}

type Registry struct {
	mu        sync.RWMutex
	framework map[Kind]map[Stage][]Definition
	business  map[Kind][]Definition
}

type RegistryOption func(*Registry)

func NewRegistry(options ...RegistryOption) *Registry {
	r := &Registry{
		framework: map[Kind]map[Stage][]Definition{},
		business:  map[Kind][]Definition{},
	}
	for _, option := range options {
		if option != nil {
			option(r)
		}
	}
	return r
}

func WithDefaultFramework() RegistryOption {
	return func(r *Registry) {
		r.Register(DefaultFrameworkDefinitions()...)
	}
}

func (r *Registry) Register(definitions ...Definition) {
	if r == nil {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, def := range definitions {
		def = normalizeDefinition(def)
		if !def.valid() {
			continue
		}
		if def.Framework {
			if r.framework[def.Kind] == nil {
				r.framework[def.Kind] = map[Stage][]Definition{}
			}
			r.framework[def.Kind][def.Stage] = append(r.framework[def.Kind][def.Stage], def)
			continue
		}
		r.business[def.Kind] = append(r.business[def.Kind], def)
	}
}

func (r *Registry) Chain(kind Kind, final Handler) Handler {
	return r.chain(kind, final, true, true)
}

func (r *Registry) FrameworkChain(kind Kind, final Handler) Handler {
	return r.chain(kind, final, true, false)
}

func (r *Registry) FrameworkChainForStages(kind Kind, stages []Stage, final Handler) Handler {
	return r.chainWithStages(kind, final, true, false, stages)
}

func (r *Registry) BusinessChain(kind Kind, final Handler) Handler {
	return r.chain(kind, final, false, true)
}

func (r *Registry) chain(kind Kind, final Handler, includeFramework bool, includeBusiness bool) Handler {
	return r.chainWithStages(kind, final, includeFramework, includeBusiness, stagesForKind(kind))
}

func (r *Registry) chainWithStages(kind Kind, final Handler, includeFramework bool, includeBusiness bool, stages []Stage) Handler {
	if final == nil {
		final = func(context.Context, *Invocation, any) (any, error) {
			return nil, nil
		}
	}
	if r == nil {
		return final
	}

	definitions := r.definitions(kind, includeFramework, includeBusiness, stages)
	for i := len(definitions) - 1; i >= 0; i-- {
		def := definitions[i]
		next := final
		final = func(ctx context.Context, inv *Invocation, req any) (any, error) {
			return def.Interceptor.Intercept(ctx, inv, req, next)
		}
	}
	return final
}

func (r *Registry) Definitions(kind Kind) []Definition {
	return r.definitions(kind, true, true, stagesForKind(kind))
}

func (r *Registry) definitions(kind Kind, includeFramework bool, includeBusiness bool, stages []Stage) []Definition {
	if r == nil {
		return nil
	}
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]Definition, 0)
	framework := r.framework[kind]
	for _, stage := range stages {
		if stage == StageBusiness {
			if !includeBusiness {
				continue
			}
			business := append([]Definition(nil), r.business[kind]...)
			sort.SliceStable(business, func(i, j int) bool {
				if business[i].Order == business[j].Order {
					return business[i].Name < business[j].Name
				}
				return business[i].Order < business[j].Order
			})
			result = append(result, business...)
			continue
		}
		if !includeFramework {
			continue
		}
		definitions := append([]Definition(nil), framework[stage]...)
		sort.SliceStable(definitions, func(i, j int) bool {
			if definitions[i].Order == definitions[j].Order {
				return definitions[i].Name < definitions[j].Name
			}
			return definitions[i].Order < definitions[j].Order
		})
		result = append(result, definitions...)
	}
	return result
}

func DefaultFrameworkDefinitions() []Definition {
	definitions := make([]Definition, 0, 8)
	for _, kind := range []Kind{KindHTTPServer, KindHTTPClient, KindGRPCServer, KindGRPCClient} {
		definitions = append(definitions,
			Framework(kind, StageRecovery, "stellar.recovery", Recovery()),
			Framework(kind, StageObserve, "stellar.request_id", RequestID("x-request-id")),
		)
	}
	return definitions
}

func Recovery() Interceptor {
	return New("stellar.recovery", func(ctx context.Context, inv *Invocation, req any, next Handler) (resp any, err error) {
		defer func() {
			if recovered := recover(); recovered != nil {
				operation := ""
				if inv != nil {
					operation = inv.Operation
				}
				err = fmt.Errorf("stellar: panic in %s interceptor chain: %v", operation, recovered)
			}
		}()
		return next(ctx, inv, req)
	})
}

type requestIDKey struct{}

var requestIDCounter atomic.Uint64

func RequestID(headerName string) Interceptor {
	if strings.TrimSpace(headerName) == "" {
		headerName = "x-request-id"
	}
	return New("stellar.request_id", func(ctx context.Context, inv *Invocation, req any, next Handler) (any, error) {
		requestID := ""
		if inv != nil {
			requestID = inv.Headers.Get(headerName)
			if requestID == "" {
				requestID = inv.RequestID
			}
			if requestID == "" {
				requestID = nextRequestID()
			}
			inv.RequestID = requestID
			inv.SetAttribute("request_id", requestID)
		}
		if requestID != "" {
			ctx = context.WithValue(ctx, requestIDKey{}, requestID)
		}
		return next(ctx, inv, req)
	})
}

func RequestIDFromContext(ctx context.Context) string {
	value, _ := ctx.Value(requestIDKey{}).(string)
	return value
}

func normalizeDefinition(def Definition) Definition {
	if def.Interceptor == nil {
		return Definition{}
	}
	if def.Name == "" {
		def.Name = def.Interceptor.Name()
	}
	if !def.Framework {
		def.Stage = StageBusiness
	}
	return def
}

func (d Definition) valid() bool {
	if d.Kind == "" || d.Stage == "" || d.Interceptor == nil {
		return false
	}
	if d.Framework {
		for _, stage := range stagesForKind(d.Kind) {
			if stage == d.Stage {
				return true
			}
		}
		return false
	}
	return d.Stage == StageBusiness
}

func stagesForKind(kind Kind) []Stage {
	switch kind {
	case KindHTTPServer, KindGRPCServer:
		return ServerInboundStages
	case KindHTTPClient, KindGRPCClient:
		return ClientOutboundStages
	default:
		return nil
	}
}

func nextRequestID() string {
	return fmt.Sprintf("stellar-%d-%d", time.Now().UnixNano(), requestIDCounter.Add(1))
}
