package grpcgoadapter

import (
	"context"
	stderrors "errors"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/stellhub/stellar/interceptor"
	"github.com/stellhub/stellar/observability"
	stellargrpc "github.com/stellhub/stellar/transport/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const Name = "rpc-grpc-go"

type Adapter struct {
	addr         string
	server       *grpc.Server
	options      []grpc.ServerOption
	observer     *observability.Provider
	interceptors *interceptor.Registry
	services     []stellargrpc.Service
	listener     net.Listener
	errCh        chan error
	mu           sync.Mutex
}

type Option func(*Adapter)

func New(addr string, opts ...Option) *Adapter {
	if addr == "" {
		addr = ":9090"
	}
	adapter := &Adapter{
		addr:     addr,
		observer: observability.New(),
		errCh:    make(chan error, 1),
	}
	for _, opt := range opts {
		opt(adapter)
	}
	return adapter
}

func WithServerOption(options ...grpc.ServerOption) Option {
	return func(adapter *Adapter) {
		adapter.options = append(adapter.options, options...)
	}
}

func WithServerInterceptors(registry *interceptor.Registry) Option {
	return func(adapter *Adapter) {
		adapter.interceptors = registry
	}
}

func (a *Adapter) Name() string {
	return Name
}

func (a *Adapter) Addr() string {
	return a.addr
}

func (a *Adapter) SetAddr(addr string) {
	if addr != "" {
		a.addr = addr
	}
}

func (a *Adapter) UseObservability(provider *observability.Provider) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if provider != nil {
		a.observer = provider
	}
}

func (a *Adapter) UseInterceptors(registry *interceptor.Registry) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.interceptors = registry
}

func (a *Adapter) Register(service stellargrpc.Service) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if service.Description == nil || service.Implementation == nil {
		return fmt.Errorf("grpc-go: service description and implementation are required")
	}
	desc, ok := service.Description.(*grpc.ServiceDesc)
	if !ok {
		return fmt.Errorf("grpc-go: service description must be *grpc.ServiceDesc")
	}
	if a.server != nil {
		a.server.RegisterService(desc, service.Implementation)
	}
	a.services = append(a.services, service)
	return nil
}

func (a *Adapter) Server() *grpc.Server {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.server == nil {
		a.buildServerLocked()
	}
	return a.server
}

func (a *Adapter) Start(ctx context.Context) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.server == nil {
		a.buildServerLocked()
	}

	listener, err := net.Listen("tcp", a.addr)
	if err != nil {
		return err
	}
	a.listener = listener

	go func() {
		if err := a.server.Serve(listener); err != nil && !stderrors.Is(err, grpc.ErrServerStopped) {
			a.errCh <- err
		}
	}()

	select {
	case err := <-a.errCh:
		return err
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(100 * time.Millisecond):
		return nil
	}
}

func (a *Adapter) Stop(ctx context.Context) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.server == nil {
		return nil
	}

	done := make(chan struct{})
	go func() {
		a.server.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		a.server.Stop()
		return ctx.Err()
	}
}

func (a *Adapter) buildServerLocked() {
	options := []grpc.ServerOption{
		grpc.StatsHandler(a.observer.GRPCServerStatsHandler()),
	}
	unaryInterceptors := []grpc.UnaryServerInterceptor{a.observer.UnaryServerInterceptor()}
	streamInterceptors := []grpc.StreamServerInterceptor{a.observer.StreamServerInterceptor()}
	if a.interceptors != nil {
		unaryInterceptors = append(unaryInterceptors, unaryServerInterceptor(a.interceptors))
		streamInterceptors = append(streamInterceptors, streamServerInterceptor(a.interceptors))
	}
	options = append(options,
		grpc.ChainUnaryInterceptor(unaryInterceptors...),
		grpc.ChainStreamInterceptor(streamInterceptors...),
	)
	options = append(options, a.options...)
	a.server = grpc.NewServer(options...)
	for _, service := range a.services {
		desc := service.Description.(*grpc.ServiceDesc)
		a.server.RegisterService(desc, service.Implementation)
	}
}

func unaryServerInterceptor(registry *interceptor.Registry) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		inv := grpcInvocation(interceptor.KindGRPCServer, info.FullMethod, headersFromIncomingContext(ctx), info)
		chain := registry.Chain(interceptor.KindGRPCServer, func(ctx context.Context, _ *interceptor.Invocation, payload any) (any, error) {
			return handler(ctx, payload)
		})
		return chain(ctx, inv, req)
	}
}

func streamServerInterceptor(registry *interceptor.Registry) grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := stream.Context()
		inv := grpcInvocation(interceptor.KindGRPCServer, info.FullMethod, headersFromIncomingContext(ctx), info)
		chain := registry.Chain(interceptor.KindGRPCServer, func(ctx context.Context, _ *interceptor.Invocation, payload any) (any, error) {
			serverStream, ok := payload.(grpc.ServerStream)
			if !ok {
				return nil, fmt.Errorf("stellar: unexpected grpc server stream type %T", payload)
			}
			return nil, handler(srv, &serverStreamWithContext{ServerStream: serverStream, ctx: ctx})
		})
		_, err := chain(ctx, inv, stream)
		return err
	}
}

type serverStreamWithContext struct {
	grpc.ServerStream
	ctx context.Context
}

func (s *serverStreamWithContext) Context() context.Context {
	return s.ctx
}

func grpcInvocation(kind interceptor.Kind, fullMethod string, headers interceptor.Header, raw any) *interceptor.Invocation {
	service, method := splitFullMethod(fullMethod)
	return &interceptor.Invocation{
		Kind:      kind,
		Protocol:  "grpc",
		Service:   service,
		Operation: fullMethod,
		Method:    method,
		Path:      fullMethod,
		Target:    fullMethod,
		Headers:   headers,
		Raw:       raw,
	}
}

func headersFromIncomingContext(ctx context.Context) interceptor.Header {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil
	}
	headers := make(interceptor.Header, len(md))
	for key, values := range md {
		headers[key] = append([]string(nil), values...)
	}
	return headers
}

func splitFullMethod(fullMethod string) (string, string) {
	fullMethod = strings.Trim(fullMethod, "/")
	service, method, ok := strings.Cut(fullMethod, "/")
	if !ok {
		return "", fullMethod
	}
	return service, method
}
