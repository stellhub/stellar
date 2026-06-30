package ratelimit

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	apperrors "github.com/stellhub/stellar/errors"
	"github.com/stellhub/stellar/governance"
	"github.com/stellhub/stellar/interceptor"
	stellpulsar "github.com/stellhub/stellpulsar-go-sdk"
)

const Name = "governance-rate-limit"

type LimitBehavior string

const (
	LimitBehaviorReject LimitBehavior = "reject"
	LimitBehaviorBlock  LimitBehavior = "block"
)

type Limiter interface {
	Allow(ctx context.Context, request LimitRequest) (LimitDecision, error)
	Wait(ctx context.Context, request LimitRequest) (LimitDecision, error)
}

type LimitRequest struct {
	Rule       governance.Rule
	Transport  string
	Service    string
	Method     string
	Path       string
	Target     string
	Resource   string
	QuotaKey   string
	Cost       int64
	Rate       int64
	Burst      int64
	Attributes map[string]string
}

type LimitDecision struct {
	Allowed    bool
	Limited    bool
	RuleID     string
	Reason     string
	RetryAfter time.Duration
	Fallback   bool
	RemoteCode string
}

type PolicyOptions struct {
	Store               *governance.Store
	Rules               []governance.Rule
	Metrics             *governance.Metrics
	Local               Limiter
	Distributed         Limiter
	DefaultMode         string
	DefaultBehavior     LimitBehavior
	DistributedFallback string
	DefaultRate         int64
	DefaultBurst        int64
}

type Policy struct {
	store           *governance.Store
	staticRules     []governance.Rule
	metrics         *governance.Metrics
	local           Limiter
	distributed     Limiter
	defaultMode     string
	defaultBehavior LimitBehavior
	defaultRate     int64
	defaultBurst    int64
}

func NewPolicy(options PolicyOptions) *Policy {
	defaultMode := strings.ToLower(strings.TrimSpace(options.DefaultMode))
	if defaultMode == "" {
		defaultMode = "local"
	}
	behavior := options.DefaultBehavior
	if behavior == "" {
		behavior = LimitBehaviorReject
	}
	local := options.Local
	if local == nil {
		local = NewLocalLimiter()
	}
	distributed := options.Distributed
	if distributed == nil {
		distributed = NewDistributedLimiter(DistributedOptions{Fallback: options.DistributedFallback})
	}
	defaultRate := options.DefaultRate
	if defaultRate <= 0 {
		defaultRate = 100
	}
	defaultBurst := options.DefaultBurst
	if defaultBurst <= 0 {
		defaultBurst = defaultRate
	}
	return &Policy{
		store:           options.Store,
		staticRules:     append([]governance.Rule(nil), options.Rules...),
		metrics:         options.Metrics,
		local:           local,
		distributed:     distributed,
		defaultMode:     defaultMode,
		defaultBehavior: behavior,
		defaultRate:     defaultRate,
		defaultBurst:    defaultBurst,
	}
}

func Definitions(options PolicyOptions) []interceptor.Definition {
	policy := NewPolicy(options)
	return []interceptor.Definition{
		interceptor.Framework(interceptor.KindHTTPServer, interceptor.StageAdmission, Name, policy),
		interceptor.Framework(interceptor.KindGRPCServer, interceptor.StageAdmission, Name, policy),
		interceptor.Framework(interceptor.KindHTTPClient, interceptor.StageAdmission, Name, policy),
		interceptor.Framework(interceptor.KindGRPCClient, interceptor.StageAdmission, Name, policy),
	}
}

func (p *Policy) Name() string {
	return Name
}

func (p *Policy) Intercept(ctx context.Context, inv *interceptor.Invocation, req any, next interceptor.Handler) (any, error) {
	if p == nil {
		return next(ctx, inv, req)
	}
	rules := p.matchingRules(inv)
	for _, rule := range rules {
		request := p.limitRequest(inv, rule)
		behavior := behaviorFromSpec(rule.Spec, p.defaultBehavior)
		mode := modeFromSpec(rule.Spec, p.defaultMode)
		if !isDistributedMode(mode) && !localLimitModeSupported(rule.Spec) {
			decision := unsupportedLocalLimitModeDecision(request)
			if !decision.Allowed {
				p.recordReject(ctx, request, behavior, decision)
				return nil, apperrors.New(apperrors.CodeUnavailable, "stellar: rate limit exceeded", http.StatusTooManyRequests)
			}
			p.recordAllow(ctx, request, behavior, decision)
			continue
		}
		limiter := p.limiterFor(mode)
		if limiter == nil {
			continue
		}
		start := time.Now()
		decision, err := decide(ctx, limiter, behavior, request)
		if behavior == LimitBehaviorBlock {
			p.recordWait(ctx, request, behavior, time.Since(start), decision)
		}
		if err != nil {
			return nil, err
		}
		if !decision.Allowed {
			p.recordReject(ctx, request, behavior, decision)
			return nil, apperrors.New(apperrors.CodeUnavailable, "stellar: rate limit exceeded", http.StatusTooManyRequests)
		}
		p.recordAllow(ctx, request, behavior, decision)
	}
	return next(ctx, inv, req)
}

func (p *Policy) matchingRules(inv *interceptor.Invocation) []governance.Rule {
	rules := make([]governance.Rule, 0)
	if p.store != nil {
		rules = append(rules, p.store.Rules(governance.RuleKindRateLimit, func(rule governance.Rule) bool {
			return rule.Enabled && invocationMatches(rule.Scope, inv)
		})...)
	}
	for _, rule := range p.staticRules {
		if rule.Enabled && rule.Kind == governance.RuleKindRateLimit && invocationMatches(rule.Scope, inv) {
			rules = append(rules, rule)
		}
	}
	return rules
}

func decide(ctx context.Context, limiter Limiter, behavior LimitBehavior, request LimitRequest) (LimitDecision, error) {
	switch behavior {
	case LimitBehaviorBlock:
		return limiter.Wait(ctx, request)
	default:
		return limiter.Allow(ctx, request)
	}
}

func (p *Policy) limiterFor(mode string) Limiter {
	if isDistributedMode(mode) {
		return p.distributed
	}
	return p.local
}

func (p *Policy) limitRequest(inv *interceptor.Invocation, rule governance.Rule) LimitRequest {
	resource := resourceKey(inv, rule)
	attrs := invocationAttributes(inv)
	mode := modeFromSpec(rule.Spec, p.defaultMode)
	cost := costFromSpec(rule.Spec)
	rate := rateFromSpec(rule.Spec, p.defaultRate)
	burst := burstFromSpec(rule.Spec, p.defaultBurst)
	if burst <= 0 {
		burst = rate
	}
	return LimitRequest{
		Rule:       rule,
		Transport:  invocationTransport(inv),
		Service:    invocationService(inv),
		Method:     invocationMethod(inv),
		Path:       invocationPath(inv),
		Target:     invocationTarget(inv),
		Resource:   resource,
		QuotaKey:   quotaKey(attrs, inv, rule, resource, isDistributedMode(mode)),
		Cost:       cost,
		Rate:       rate,
		Burst:      burst,
		Attributes: attrs,
	}
}

func (p *Policy) recordAllow(ctx context.Context, request LimitRequest, behavior LimitBehavior, decision LimitDecision) {
	if p.metrics == nil {
		return
	}
	p.metrics.RecordRateLimitAllow(ctx, metricAttrs(request, behavior, decision, "allowed"))
}

func (p *Policy) recordReject(ctx context.Context, request LimitRequest, behavior LimitBehavior, decision LimitDecision) {
	if p.metrics == nil {
		return
	}
	p.metrics.RecordRateLimitReject(ctx, metricAttrs(request, behavior, decision, "rejected"))
}

func (p *Policy) recordWait(ctx context.Context, request LimitRequest, behavior LimitBehavior, duration time.Duration, decision LimitDecision) {
	if p.metrics == nil {
		return
	}
	p.metrics.RecordRateLimitWait(ctx, duration, metricAttrs(request, behavior, decision, "wait"))
}

func metricAttrs(request LimitRequest, behavior LimitBehavior, decision LimitDecision, outcome string) governance.MetricAttrs {
	if decision.Fallback {
		outcome = outcome + ".fallback"
	}
	return governance.MetricAttrs{
		Adapter:   "stellorbit",
		RuleKind:  string(request.Rule.Kind),
		RuleID:    request.Rule.ID,
		Transport: request.Transport,
		Service:   request.Service,
		Method:    request.Method,
		Resource:  request.Resource,
		Mode:      metricModeFromSpec(request.Rule.Spec),
		Behavior:  string(behavior),
		Outcome:   outcome,
	}
}

type LocalLimiter struct {
	buckets sync.Map
}

type bucket struct {
	mu     sync.Mutex
	tokens float64
	last   time.Time
}

func NewLocalLimiter() *LocalLimiter {
	return &LocalLimiter{}
}

func (l *LocalLimiter) Allow(_ context.Context, request LimitRequest) (LimitDecision, error) {
	current := l.bucketFor(request)
	allowed, retryAfter := current.take(request.Rate, request.Burst, request.Cost)
	if allowed {
		return allowedDecision(request.Rule.ID), nil
	}
	return limitedDecision(request.Rule.ID, "local_rate_limited", retryAfter), nil
}

func (l *LocalLimiter) Wait(ctx context.Context, request LimitRequest) (LimitDecision, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	for {
		decision, err := l.Allow(ctx, request)
		if err != nil || decision.Allowed {
			return decision, err
		}
		delay := decision.RetryAfter
		if delay <= 0 {
			delay = 10 * time.Millisecond
		}
		timer := time.NewTimer(delay)
		select {
		case <-ctx.Done():
			timer.Stop()
			return limitedDecision(request.Rule.ID, ctx.Err().Error(), 0), nil
		case <-timer.C:
		}
	}
}

func (l *LocalLimiter) bucketFor(request LimitRequest) *bucket {
	key := request.Rule.ID + ":" + request.Resource + ":" + request.QuotaKey
	value, _ := l.buckets.LoadOrStore(key, &bucket{tokens: float64(request.Burst), last: time.Now()})
	return value.(*bucket)
}

func (b *bucket) take(rate int64, burst int64, cost int64) (bool, time.Duration) {
	if rate <= 0 {
		return true, 0
	}
	if burst <= 0 {
		burst = rate
	}
	if cost <= 0 {
		cost = 1
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	now := time.Now()
	if b.last.IsZero() {
		b.last = now
	}
	elapsed := now.Sub(b.last).Seconds()
	b.tokens += elapsed * float64(rate)
	if maxTokens := float64(burst); b.tokens > maxTokens {
		b.tokens = maxTokens
	}
	b.last = now
	if b.tokens >= float64(cost) {
		b.tokens -= float64(cost)
		return true, 0
	}
	missing := float64(cost) - b.tokens
	wait := time.Duration(missing / float64(rate) * float64(time.Second))
	if wait <= 0 {
		wait = time.Millisecond
	}
	return false, wait
}

type DistributedLimiter struct {
	client   stellpulsar.StellpulsarClient
	timeout  time.Duration
	fallback string
	appID    string
}

type DistributedOptions struct {
	Client   stellpulsar.StellpulsarClient
	Timeout  time.Duration
	Fallback string
	AppID    string
}

func NewDistributedLimiter(options DistributedOptions) *DistributedLimiter {
	timeout := options.Timeout
	if timeout <= 0 {
		timeout = 100 * time.Millisecond
	}
	fallback := strings.ToLower(strings.TrimSpace(options.Fallback))
	if fallback == "" {
		fallback = "fail_open"
	}
	return &DistributedLimiter{
		client:   options.Client,
		timeout:  timeout,
		fallback: fallback,
		appID:    options.AppID,
	}
}

func (l *DistributedLimiter) Allow(ctx context.Context, request LimitRequest) (LimitDecision, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if l == nil || l.client == nil {
		return l.fallbackDecision(request, "stellpulsar_client_missing"), nil
	}
	callCtx := ctx
	cancel := func() {}
	if l.timeout > 0 {
		callCtx, cancel = context.WithTimeout(ctx, l.timeout)
	}
	defer cancel()
	result, err := l.client.TryAcquire(callCtx, stellpulsar.RateLimitRequest{
		RequestID:       requestID(request),
		ApplicationCode: l.applicationCode(request),
		TargetService:   request.Service,
		Resource:        request.Resource,
		Method:          request.Method,
		TenantID:        request.Attributes["tenant_id"],
		UserID:          request.Attributes["user_id"],
		QuotaKey:        request.QuotaKey,
		Cost:            request.Cost,
		Attributes:      request.Attributes,
	})
	if err != nil {
		return l.fallbackDecision(request, err.Error()), nil
	}
	return decisionFromPulsar(request, result), nil
}

func (l *DistributedLimiter) Wait(ctx context.Context, request LimitRequest) (LimitDecision, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	for {
		decision, err := l.Allow(ctx, request)
		if err != nil || decision.Allowed {
			return decision, err
		}
		delay := decision.RetryAfter
		if delay <= 0 {
			delay = 50 * time.Millisecond
		}
		timer := time.NewTimer(delay)
		select {
		case <-ctx.Done():
			timer.Stop()
			return l.fallbackDecision(request, ctx.Err().Error()), nil
		case <-timer.C:
		}
	}
}

func (l *DistributedLimiter) applicationCode(request LimitRequest) string {
	if strings.TrimSpace(l.appID) != "" {
		return strings.TrimSpace(l.appID)
	}
	return request.Service
}

func (l *DistributedLimiter) fallbackDecision(request LimitRequest, reason string) LimitDecision {
	fallback := ""
	if l != nil {
		fallback = l.fallback
	}
	allowed := fallbackPermits(request.Rule.Spec, fallback)
	return LimitDecision{
		Allowed:    allowed,
		Limited:    !allowed,
		RuleID:     request.Rule.ID,
		Reason:     reason,
		Fallback:   true,
		RemoteCode: "CLIENT_ERROR",
	}
}

func decisionFromPulsar(request LimitRequest, result *stellpulsar.RateLimitResult) LimitDecision {
	if result == nil {
		return LimitDecision{Allowed: true, RuleID: request.Rule.ID, Reason: "no_result"}
	}
	return LimitDecision{
		Allowed:    result.Permitted,
		Limited:    result.Limited,
		RuleID:     firstNonBlank(result.RuleID, request.Rule.ID),
		Reason:     result.Reason,
		RetryAfter: time.Duration(result.RetryAfterMS) * time.Millisecond,
		Fallback:   result.Fallback,
		RemoteCode: string(result.Decision),
	}
}

func allowedDecision(ruleID string) LimitDecision {
	return LimitDecision{Allowed: true, RuleID: ruleID, Reason: "allowed"}
}

func limitedDecision(ruleID string, reason string, retryAfter time.Duration) LimitDecision {
	return LimitDecision{Allowed: false, Limited: true, RuleID: ruleID, Reason: reason, RetryAfter: retryAfter}
}

func invocationMatches(scope governance.Scope, inv *interceptor.Invocation) bool {
	if inv == nil {
		return true
	}
	if scope.Transport != "" && !strings.EqualFold(scope.Transport, string(inv.Kind)) && !strings.EqualFold(scope.Transport, inv.Protocol+".client") && !strings.EqualFold(scope.Transport, inv.Protocol+".server") {
		return false
	}
	if scope.Service != "" && !strings.EqualFold(scope.Service, inv.Service) {
		return false
	}
	if scope.Method != "" && !strings.EqualFold(scope.Method, inv.Method) {
		return false
	}
	if scope.Path != "" && !strings.HasPrefix(inv.Path, scope.Path) {
		return false
	}
	if scope.Target != "" && !strings.EqualFold(scope.Target, inv.Target) {
		return false
	}
	return true
}

func behaviorFromSpec(spec map[string]any, fallback LimitBehavior) LimitBehavior {
	switch strings.ToLower(strings.TrimSpace(stringSpec(spec, "behavior", "limit_behavior"))) {
	case string(LimitBehaviorBlock):
		return LimitBehaviorBlock
	case string(LimitBehaviorReject):
		return LimitBehaviorReject
	default:
		if fallback == "" {
			return LimitBehaviorReject
		}
		return fallback
	}
}

func modeFromSpec(spec map[string]any, fallback string) string {
	if value := stringSpec(spec, "coordinationMode", "coordination_mode"); value != "" {
		return modeFromCoordination(value)
	}
	if value := stringSpec(spec, "enforcementMode", "enforcement_mode"); value != "" {
		return modeFromCoordination(value)
	}
	if value := stringSpec(spec, "mode", "backend"); value != "" {
		return normalizeMode(value)
	}
	return normalizeMode(fallback)
}

func modeFromCoordination(value string) string {
	switch normalizeMode(value) {
	case "distributed", "global", "global-sync", "global-quota", "edge":
		return "distributed"
	case "local", "local-only":
		return "local"
	default:
		return normalizeMode(value)
	}
}

func isDistributedMode(mode string) bool {
	switch normalizeMode(mode) {
	case "distributed", "global", "global-sync", "global-quota", "edge":
		return true
	default:
		return false
	}
}

func normalizeMode(value string) string {
	return strings.ToLower(strings.ReplaceAll(strings.TrimSpace(value), "_", "-"))
}

func limitModeFromSpec(spec map[string]any) string {
	return normalizeMode(stringSpec(spec, "limitMode", "limit_mode"))
}

func metricModeFromSpec(spec map[string]any) string {
	if mode := limitModeFromSpec(spec); mode != "" {
		return mode
	}
	return modeFromSpec(spec, "")
}

func localLimitModeSupported(spec map[string]any) bool {
	switch limitModeFromSpec(spec) {
	case "", "qps", "header", "hot-key", "quota", "bandwidth", "model":
		return true
	case "custom", "concurrency", "connection":
		return false
	default:
		return false
	}
}

func unsupportedLocalLimitModeDecision(request LimitRequest) LimitDecision {
	reason := "unsupported_local_limit_mode"
	if mode := limitModeFromSpec(request.Rule.Spec); mode != "" {
		reason = reason + ":" + mode
	}
	allowed := fallbackPermits(request.Rule.Spec, "fail_open")
	return LimitDecision{
		Allowed:  allowed,
		Limited:  !allowed,
		RuleID:   request.Rule.ID,
		Reason:   reason,
		Fallback: true,
	}
}

func resourceKey(inv *interceptor.Invocation, rule governance.Rule) string {
	if value := stringSpec(rule.Spec, "resource", "resource_key", "resourceKey"); value != "" {
		return value
	}
	if inv == nil {
		return rule.ID
	}
	return strings.Trim(strings.Join([]string{inv.Service, inv.Method, inv.Path}, ":"), ":")
}

func quotaKey(attrs map[string]string, inv *interceptor.Invocation, rule governance.Rule, resource string, deferExtractor bool) string {
	if value := stringSpec(rule.Spec, "quota_key", "quotaKey"); value != "" {
		return value
	}
	if hasKeyExtractor(rule.Spec) {
		if deferExtractor {
			return ""
		}
		if value := extractorQuotaKey(rule.Spec, attrs); value != "" {
			return value
		}
	}
	if inv != nil {
		if value := inv.Headers.Get("x-stellar-quota-key"); value != "" {
			return value
		}
		if value := inv.Headers.Get("x-tenant-id"); value != "" {
			return value
		}
		if inv.Attributes != nil {
			for _, key := range []string{"quota_key", "quotaKey", "tenant_id", "tenantId", "user_id", "userId"} {
				if value := fmt.Sprint(inv.Attributes[key]); value != "" && value != "<nil>" {
					return strings.TrimSpace(value)
				}
			}
		}
	}
	for _, key := range []string{"quota_key", "quotaKey", "tenant_id", "tenantId", "user_id", "userId"} {
		if value := strings.TrimSpace(attrs[key]); value != "" {
			return value
		}
	}
	return resource
}

func invocationAttributes(inv *interceptor.Invocation) map[string]string {
	attrs := map[string]string{}
	if inv == nil {
		return attrs
	}
	put(attrs, "transport", string(inv.Kind))
	put(attrs, "protocol", inv.Protocol)
	put(attrs, "service", inv.Service)
	put(attrs, "serviceName", inv.Service)
	put(attrs, "operation", inv.Operation)
	put(attrs, "method", inv.Method)
	put(attrs, "path", inv.Path)
	put(attrs, "pathTemplate", inv.Path)
	put(attrs, "target", inv.Target)
	put(attrs, "request_id", inv.RequestID)
	put(attrs, "requestId", inv.RequestID)
	for key, values := range inv.Headers {
		if len(values) == 0 {
			continue
		}
		value := values[0]
		lowerKey := strings.ToLower(strings.TrimSpace(key))
		put(attrs, key, value)
		put(attrs, lowerKey, value)
		put(attrs, "header."+lowerKey, value)
		put(attrs, "metadata."+lowerKey, value)
	}
	put(attrs, "tenant_id", inv.Headers.Get("x-tenant-id"))
	put(attrs, "tenantId", inv.Headers.Get("x-tenant-id"))
	put(attrs, "user_id", inv.Headers.Get("x-user-id"))
	put(attrs, "userId", inv.Headers.Get("x-user-id"))
	put(attrs, "quota_key", inv.Headers.Get("x-stellar-quota-key"))
	put(attrs, "quotaKey", inv.Headers.Get("x-stellar-quota-key"))
	for key, value := range inv.Attributes {
		put(attrs, key, fmt.Sprint(value))
	}
	if value := attrs["http.route"]; value != "" {
		put(attrs, "pathTemplate", value)
	}
	if value := attrs["grpc.service"]; value != "" {
		put(attrs, "grpcService", value)
	}
	return attrs
}

func invocationTransport(inv *interceptor.Invocation) string {
	if inv == nil {
		return ""
	}
	return string(inv.Kind)
}

func invocationService(inv *interceptor.Invocation) string {
	if inv == nil {
		return ""
	}
	return inv.Service
}

func invocationMethod(inv *interceptor.Invocation) string {
	if inv == nil {
		return ""
	}
	return inv.Method
}

func invocationPath(inv *interceptor.Invocation) string {
	if inv == nil {
		return ""
	}
	return inv.Path
}

func invocationTarget(inv *interceptor.Invocation) string {
	if inv == nil {
		return ""
	}
	return inv.Target
}

func requestID(request LimitRequest) string {
	if value := request.Attributes["request_id"]; value != "" {
		return value
	}
	return request.Rule.ID + "-" + strconv.FormatInt(time.Now().UnixNano(), 10)
}

func costFromSpec(spec map[string]any) int64 {
	if value := firstPositiveInt64(
		int64Spec(spec, 0, "cost"),
		int64FromMap(mapAny(spec["quotaConfig"]), "cost"),
		int64FromMap(mapAny(spec["modelLimitConfig"]), "cost", "tokenCost", "requestCost"),
	); value > 0 {
		return value
	}
	return 1
}

func rateFromSpec(spec map[string]any, fallback int64) int64 {
	if value := firstPositiveInt64(
		int64Spec(spec, 0, "rate", "qps", "quota", "limit"),
		int64FromMap(mapAny(spec["quotaConfig"]), "limit", "quota", "qps", "rate"),
		int64FromMap(mapAny(spec["windowConfig"]), "rate", "qps", "limit"),
		int64FromMap(mapAny(spec["concurrencyConfig"]), "maxConcurrent", "max_concurrent", "limit"),
		int64FromMap(mapAny(spec["hotspotConfig"]), "threshold", "limit"),
		int64FromMap(mapAny(spec["modelLimitConfig"]), "requestLimit", "request_limit", "tokenLimit", "token_limit", "costLimit", "cost_limit"),
	); value > 0 {
		return value
	}
	return fallback
}

func burstFromSpec(spec map[string]any, fallback int64) int64 {
	if value := firstPositiveInt64(
		int64Spec(spec, 0, "burst"),
		int64FromMap(mapAny(spec["burstConfig"]), "capacity", "burst", "maxBurst", "max_burst"),
		int64FromMap(mapAny(spec["quotaConfig"]), "burst"),
	); value > 0 {
		return value
	}
	return fallback
}

func firstPositiveInt64(values ...int64) int64 {
	for _, value := range values {
		if value > 0 {
			return value
		}
	}
	return 0
}

func int64FromMap(values map[string]any, keys ...string) int64 {
	if len(values) == 0 {
		return 0
	}
	for _, key := range keys {
		if value, ok := values[key]; ok {
			if parsed, ok := positiveInt64(value); ok {
				return parsed
			}
		}
	}
	return 0
}

func positiveInt64(value any) (int64, bool) {
	switch typed := value.(type) {
	case int:
		if typed > 0 {
			return int64(typed), true
		}
	case int8:
		if typed > 0 {
			return int64(typed), true
		}
	case int16:
		if typed > 0 {
			return int64(typed), true
		}
	case int32:
		if typed > 0 {
			return int64(typed), true
		}
	case int64:
		if typed > 0 {
			return typed, true
		}
	case uint:
		if typed > 0 {
			return int64(typed), true
		}
	case uint8:
		if typed > 0 {
			return int64(typed), true
		}
	case uint16:
		if typed > 0 {
			return int64(typed), true
		}
	case uint32:
		if typed > 0 {
			return int64(typed), true
		}
	case uint64:
		if typed > 0 && typed <= uint64(^uint64(0)>>1) {
			return int64(typed), true
		}
	case float32:
		if typed > 0 {
			return int64(typed), true
		}
	case float64:
		if typed > 0 {
			return int64(typed), true
		}
	case string:
		parsed, err := strconv.ParseInt(strings.TrimSpace(typed), 10, 64)
		if err == nil && parsed > 0 {
			return parsed, true
		}
	}
	return 0, false
}

func hasKeyExtractor(spec map[string]any) bool {
	return len(keyExtractorMap(spec)) > 0
}

func extractorQuotaKey(spec map[string]any, attrs map[string]string) string {
	extractor := keyExtractorMap(spec)
	keys := extractorKeys(extractor)
	if len(keys) == 0 {
		return ""
	}
	failOnMissing := boolFromAny(extractor["failOnMissing"]) || boolFromAny(extractor["fail_on_missing"])
	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		value := extractorValue(key, attrs)
		if value == "" {
			if boolFromAny(key["required"]) || failOnMissing {
				return ""
			}
			continue
		}
		label := firstNonBlank(stringAny(key["name"]), stringAny(key["key"]), strings.ToLower(enumToken(stringAny(key["source"]))))
		if label == "" {
			label = "key"
		}
		parts = append(parts, label+"="+value)
	}
	if len(parts) == 0 {
		return ""
	}
	return strings.Join(parts, ":")
}

func keyExtractorMap(spec map[string]any) map[string]any {
	if spec == nil {
		return nil
	}
	if value := mapAny(spec["keyExtractor"]); len(value) > 0 {
		return value
	}
	return mapAny(spec["key_extractor"])
}

func extractorKeys(extractor map[string]any) []map[string]any {
	values := sliceAny(extractor["keys"])
	keys := make([]map[string]any, 0, len(values))
	for _, value := range values {
		if key := mapAny(value); len(key) > 0 {
			keys = append(keys, key)
		}
	}
	return keys
}

func extractorValue(key map[string]any, attrs map[string]string) string {
	source := enumToken(stringAny(key["source"]))
	lookupKey := firstNonBlank(
		stringAny(key["key"]),
		stringAny(key["attribute"]),
		stringAny(key["field"]),
		stringAny(key["path"]),
		stringAny(key["name"]),
	)
	var value string
	switch source {
	case "TENANT":
		value = firstNonBlank(lookupAttribute(attrs, lookupKey), attrs["tenant_id"], attrs["tenantId"])
	case "USER":
		value = firstNonBlank(lookupAttribute(attrs, lookupKey), attrs["user_id"], attrs["userId"])
	case "CALLER":
		value = firstNonBlank(lookupAttribute(attrs, lookupKey), attrs["caller"], attrs["callerId"])
	case "HTTP_PATH", "PATH", "ENDPOINT":
		value = firstNonBlank(lookupAttribute(attrs, lookupKey), attrs["pathTemplate"], attrs["path"], attrs["resource"])
	case "HTTP_METHOD", "GRPC_METHOD", "METHOD":
		value = firstNonBlank(lookupAttribute(attrs, lookupKey), attrs["method"])
	case "GRPC_SERVICE":
		value = firstNonBlank(lookupAttribute(attrs, lookupKey), attrs["grpc.service"], attrs["grpcService"], attrs["serviceName"])
	case "REMOTE_IP", "IP":
		value = firstNonBlank(lookupAttribute(attrs, lookupKey), attrs["remoteIp"], attrs["remoteIP"], attrs["clientIp"], attrs["clientIP"])
	case "HEADER", "GRPC_METADATA", "QUERY", "COOKIE", "JWT_CLAIM", "BODY_JSON_PATH", "API_KEY", "CUSTOM_KEY", "CUSTOM_EXPRESSION":
		value = lookupAttribute(attrs, lookupKey)
	case "TOPIC":
		value = firstNonBlank(lookupAttribute(attrs, lookupKey), attrs["topic"], attrs["resource"])
	case "MODEL", "MODEL_REQUEST", "MODEL_TOKEN", "MODEL_COST", "MODEL_CONCURRENCY":
		value = firstNonBlank(lookupAttribute(attrs, lookupKey), attrs["model"], attrs["modelName"], attrs["resource"])
	default:
		value = firstNonBlank(lookupAttribute(attrs, lookupKey), lookupAttribute(attrs, stringAny(key["name"])))
	}
	return normalizeExtractorValue(value, key["normalize"])
}

func lookupAttribute(values map[string]string, key string) string {
	key = strings.TrimSpace(key)
	if key == "" {
		return ""
	}
	if value := values[key]; value != "" {
		return value
	}
	lowerKey := strings.ToLower(key)
	if value := values[lowerKey]; value != "" {
		return value
	}
	if value := values["header."+lowerKey]; value != "" {
		return value
	}
	if value := values["metadata."+lowerKey]; value != "" {
		return value
	}
	for candidate, value := range values {
		if strings.EqualFold(candidate, key) {
			return value
		}
	}
	return ""
}

func normalizeExtractorValue(value string, normalize any) string {
	value = strings.TrimSpace(value)
	switch enumToken(normalizeName(normalize)) {
	case "", "NONE", "TRIM":
		return value
	case "LOWERCASE", "LOWER_CASE":
		return strings.ToLower(value)
	case "UPPERCASE", "UPPER_CASE":
		return strings.ToUpper(value)
	default:
		return value
	}
}

func fallbackPermits(spec map[string]any, fallback string) bool {
	return failPolicyFromSpec(spec, fallback) != "fail-closed"
}

func failPolicyFromSpec(spec map[string]any, fallback string) string {
	value := firstNonBlank(
		stringSpec(spec, "failPolicy", "fail_policy"),
		stringFromMap(mapAny(spec["fallbackPolicy"]), "failPolicy", "fail_policy", "policy", "mode"),
		stringFromMap(mapAny(spec["customPolicy"]), "failPolicy", "fail_policy"),
		fallback,
	)
	switch normalizeMode(value) {
	case "fail-closed", "closed":
		return "fail-closed"
	default:
		return "fail-open"
	}
}

func stringFromMap(values map[string]any, keys ...string) string {
	for _, key := range keys {
		if value := stringAny(values[key]); value != "" {
			return value
		}
	}
	return ""
}

func normalizeName(value any) string {
	switch typed := value.(type) {
	case map[string]any:
		return firstNonBlank(stringAny(typed["type"]), stringAny(typed["mode"]), stringAny(typed["name"]))
	case map[string]string:
		return firstNonBlank(typed["type"], typed["mode"], typed["name"])
	default:
		return stringAny(value)
	}
}

func mapAny(value any) map[string]any {
	switch typed := value.(type) {
	case map[string]any:
		return typed
	case map[string]string:
		values := make(map[string]any, len(typed))
		for key, item := range typed {
			values[key] = item
		}
		return values
	default:
		return nil
	}
}

func sliceAny(value any) []any {
	switch typed := value.(type) {
	case []any:
		return typed
	case []map[string]any:
		values := make([]any, 0, len(typed))
		for _, item := range typed {
			values = append(values, item)
		}
		return values
	case []map[string]string:
		values := make([]any, 0, len(typed))
		for _, item := range typed {
			values = append(values, item)
		}
		return values
	default:
		return nil
	}
}

func boolFromAny(value any) bool {
	switch typed := value.(type) {
	case bool:
		return typed
	case string:
		return strings.EqualFold(strings.TrimSpace(typed), "true")
	default:
		return strings.EqualFold(strings.TrimSpace(fmt.Sprint(value)), "true")
	}
}

func stringAny(value any) string {
	if value == nil {
		return ""
	}
	text := strings.TrimSpace(fmt.Sprint(value))
	if text == "<nil>" {
		return ""
	}
	return text
}

func enumToken(value string) string {
	return strings.ToUpper(strings.ReplaceAll(strings.ReplaceAll(strings.TrimSpace(value), "-", "_"), ".", "_"))
}

func int64Spec(spec map[string]any, fallback int64, keys ...string) int64 {
	for _, key := range keys {
		if spec == nil {
			return fallback
		}
		value, ok := spec[key]
		if !ok {
			continue
		}
		switch typed := value.(type) {
		case int:
			if typed > 0 {
				return int64(typed)
			}
		case int64:
			if typed > 0 {
				return typed
			}
		case float64:
			if typed > 0 {
				return int64(typed)
			}
		case string:
			parsed, err := strconv.ParseInt(strings.TrimSpace(typed), 10, 64)
			if err == nil && parsed > 0 {
				return parsed
			}
		}
	}
	return fallback
}

func stringSpec(spec map[string]any, keys ...string) string {
	for _, key := range keys {
		if spec == nil {
			return ""
		}
		value, ok := spec[key]
		if !ok {
			continue
		}
		text := strings.TrimSpace(fmt.Sprint(value))
		if text != "" && text != "<nil>" {
			return text
		}
	}
	return ""
}

func put(values map[string]string, key string, value string) {
	key = strings.TrimSpace(key)
	value = strings.TrimSpace(value)
	if key != "" && value != "" && value != "<nil>" {
		values[key] = value
	}
}

func firstNonBlank(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}
