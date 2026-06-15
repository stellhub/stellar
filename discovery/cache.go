package discovery

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/stellhub/stellar/config"
	"github.com/stellhub/stellar/observability"
)

type CacheOption func(*CachedResolver)

type CachedResolver struct {
	base            Resolver
	target          Target
	refreshInterval time.Duration
	staleTTL        time.Duration
	mu              sync.RWMutex
	endpoints       []Endpoint
	updatedAt       time.Time
	lastErr         error
	observer        *observability.Provider
	loadBalance     string
	metricsEnabled  bool
	startOnce       sync.Once
	stopOnce        sync.Once
	cancel          context.CancelFunc
}

func NewCachedResolver(base Resolver, target Target, options ...CacheOption) *CachedResolver {
	resolver := &CachedResolver{
		base:            base,
		target:          NormalizeTarget(target),
		refreshInterval: defaultRefreshInterval,
		staleTTL:        defaultStaleTTL,
		loadBalance:     DefaultLoadBalance,
		metricsEnabled:  true,
	}
	for _, option := range options {
		option(resolver)
	}
	if resolver.observer == nil {
		resolver.observer = observability.New()
	}
	return resolver
}

func NewCachedResolverFromConfig(ctx context.Context, cfg *config.DiscoveryConfig, target Target, options ...CacheOption) (*CachedResolver, error) {
	prepared := &CachedResolver{metricsEnabled: true}
	for _, option := range options {
		option(prepared)
	}
	base, err := NewRegistryResolverFromConfig(ctx, cfg, prepared.observer)
	if err != nil {
		return nil, err
	}
	refreshInterval, err := RefreshInterval(cfg)
	if err != nil {
		return nil, err
	}
	staleTTL, err := StaleTTL(cfg)
	if err != nil {
		return nil, err
	}
	loadBalance := DefaultLoadBalance
	if cfg != nil && cfg.LoadBalance != "" {
		loadBalance = cfg.LoadBalance
	}
	defaultOptions := []CacheOption{
		WithRefreshInterval(refreshInterval),
		WithStaleTTL(staleTTL),
		WithLoadBalance(loadBalance),
	}
	if cfg != nil && cfg.Observability.Metrics != nil {
		defaultOptions = append(defaultOptions, WithMetricsEnabled(*cfg.Observability.Metrics))
	}
	defaultOptions = append(defaultOptions, options...)
	return NewCachedResolver(base, target, defaultOptions...), nil
}

func WithRefreshInterval(interval time.Duration) CacheOption {
	return func(resolver *CachedResolver) {
		if interval > 0 {
			resolver.refreshInterval = interval
		}
	}
}

func WithStaleTTL(ttl time.Duration) CacheOption {
	return func(resolver *CachedResolver) {
		if ttl > 0 {
			resolver.staleTTL = ttl
		}
	}
}

func WithObservability(provider *observability.Provider) CacheOption {
	return func(resolver *CachedResolver) {
		if provider != nil {
			resolver.observer = provider
		}
	}
}

func WithLoadBalance(policy string) CacheOption {
	return func(resolver *CachedResolver) {
		if policy != "" {
			resolver.loadBalance = policy
		}
	}
}

func WithMetricsEnabled(enabled bool) CacheOption {
	return func(resolver *CachedResolver) {
		resolver.metricsEnabled = enabled
	}
}

func (r *CachedResolver) Resolve(ctx context.Context, target Target) ([]Endpoint, error) {
	if r == nil || r.base == nil {
		return nil, ErrNoAvailableEndpoint
	}
	ctx, finish := r.startDiscovery(ctx, "resolve")
	r.start()
	if endpoints, ok := r.snapshot(false); ok {
		finish(observability.DiscoveryResult{
			Result:    observability.DiscoveryResultCacheHit,
			Endpoints: len(endpoints),
			CacheAge:  r.cacheAge(),
			Stale:     false,
		})
		return endpoints, nil
	}
	if err := r.refresh(ctx); err != nil {
		if endpoints, ok := r.snapshot(true); ok {
			finish(observability.DiscoveryResult{
				Result:    observability.DiscoveryResultStale,
				Endpoints: len(endpoints),
				CacheAge:  r.cacheAge(),
				Stale:     true,
			})
			return endpoints, nil
		}
		finish(observability.DiscoveryResult{Err: err})
		return nil, err
	}
	if endpoints, ok := r.snapshot(true); ok {
		finish(observability.DiscoveryResult{
			Result:    observability.DiscoveryResultRefresh,
			Endpoints: len(endpoints),
			CacheAge:  r.cacheAge(),
			Stale:     false,
		})
		return endpoints, nil
	}
	finish(observability.DiscoveryResult{Err: ErrNoAvailableEndpoint})
	return nil, ErrNoAvailableEndpoint
}

func (r *CachedResolver) Pick(ctx context.Context, picker Picker) (Endpoint, error) {
	if picker == nil {
		picker = NewPicker(DefaultLoadBalance)
	}
	ctx, finish := r.startDiscovery(ctx, "pick")
	endpoints, err := r.Resolve(ctx, r.target)
	if err != nil {
		finish(observability.DiscoveryResult{Err: err})
		return Endpoint{}, err
	}
	endpoint, err := picker.Pick(endpoints)
	finish(observability.DiscoveryResult{
		Endpoints: len(endpoints),
		CacheAge:  r.cacheAge(),
		Stale:     r.isStale(),
		Err:       err,
	})
	return endpoint, err
}

func (r *CachedResolver) Watch(ctx context.Context, target Target) (Watcher, error) {
	if r == nil || r.base == nil {
		return nil, ErrNoAvailableEndpoint
	}
	return r.base.Watch(ctx, target)
}

func (r *CachedResolver) Close(ctx context.Context) error {
	if r == nil {
		return nil
	}
	r.stopOnce.Do(func() {
		if r.cancel != nil {
			r.cancel()
		}
	})
	if r.base != nil {
		return r.base.Close(ctx)
	}
	return nil
}

func (r *CachedResolver) start() {
	r.startOnce.Do(func() {
		ctx, cancel := context.WithCancel(context.Background())
		r.cancel = cancel
		go r.refreshLoop(ctx)
		go r.watchLoop(ctx)
	})
}

func (r *CachedResolver) refreshLoop(ctx context.Context) {
	ticker := time.NewTicker(r.refreshInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			_ = r.refresh(ctx)
		case <-ctx.Done():
			return
		}
	}
}

func (r *CachedResolver) watchLoop(ctx context.Context) {
	for {
		watcher, err := r.base.Watch(ctx, r.target)
		if err != nil {
			r.setError(err)
			if !sleepOrDone(ctx, r.refreshInterval) {
				return
			}
			continue
		}
		events := watcher.Events()
		for {
			select {
			case event, ok := <-events:
				if !ok {
					_ = watcher.Close()
					goto reconnect
				}
				r.apply(event)
			case <-ctx.Done():
				_ = watcher.Close()
				return
			}
		}
	reconnect:
		if !sleepOrDone(ctx, time.Second) {
			return
		}
	}
}

func (r *CachedResolver) refresh(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}
	ctx, finish := r.startDiscovery(ctx, "refresh")
	endpoints, err := r.base.Resolve(ctx, r.target)
	if err != nil {
		r.setError(err)
		finish(observability.DiscoveryResult{Err: err})
		return err
	}
	r.mu.Lock()
	r.endpoints = cloneEndpoints(endpoints)
	r.updatedAt = time.Now()
	r.lastErr = nil
	r.mu.Unlock()
	finish(observability.DiscoveryResult{
		Result:    observability.DiscoveryResultRefresh,
		Endpoints: len(endpoints),
		CacheAge:  0,
		Stale:     false,
	})
	return nil
}

func (r *CachedResolver) snapshot(allowStale bool) ([]Endpoint, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if len(r.endpoints) == 0 {
		return nil, false
	}
	if !allowStale && r.staleTTL > 0 && !r.updatedAt.IsZero() && time.Since(r.updatedAt) > r.staleTTL {
		return nil, false
	}
	return cloneEndpoints(r.endpoints), true
}

func (r *CachedResolver) setError(err error) {
	if err == nil {
		return
	}
	r.mu.Lock()
	r.lastErr = err
	r.mu.Unlock()
}

func (r *CachedResolver) apply(event Event) {
	switch event.Type {
	case EventUpsert:
		if event.Endpoint == nil {
			return
		}
		r.upsert(*event.Endpoint)
	case EventDelete:
		if event.Endpoint == nil {
			return
		}
		r.delete(*event.Endpoint)
	default:
		r.mu.Lock()
		r.endpoints = cloneEndpoints(event.Endpoints)
		r.updatedAt = time.Now()
		r.lastErr = nil
		r.mu.Unlock()
	}
	r.recordWatchEvent(event)
}

func (r *CachedResolver) upsert(endpoint Endpoint) {
	endpoint = NormalizeEndpoint(endpoint)
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, existing := range r.endpoints {
		if sameEndpoint(existing, endpoint) {
			r.endpoints[i] = endpoint
			r.updatedAt = time.Now()
			return
		}
	}
	r.endpoints = append(r.endpoints, endpoint)
	r.updatedAt = time.Now()
}

func (r *CachedResolver) delete(endpoint Endpoint) {
	endpoint = NormalizeEndpoint(endpoint)
	r.mu.Lock()
	defer r.mu.Unlock()
	filtered := r.endpoints[:0]
	for _, existing := range r.endpoints {
		if sameEndpoint(existing, endpoint) {
			continue
		}
		filtered = append(filtered, existing)
	}
	r.endpoints = append([]Endpoint(nil), filtered...)
	r.updatedAt = time.Now()
}

func (r *CachedResolver) startDiscovery(ctx context.Context, operation string) (context.Context, func(observability.DiscoveryResult)) {
	provider := (*observability.Provider)(nil)
	target := Target{}
	loadBalance := ""
	if r != nil {
		provider = r.observer
		target = r.target
		loadBalance = r.loadBalance
		if !r.metricsEnabled {
			if ctx == nil {
				ctx = context.Background()
			}
			return ctx, func(observability.DiscoveryResult) {}
		}
	}
	if provider == nil {
		provider = observability.New()
	}
	target = NormalizeTarget(target)
	return provider.StartDiscovery(ctx, observability.DiscoveryRequest{
		Resolver:    "registry",
		Operation:   operation,
		Namespace:   target.Namespace,
		Service:     target.Service,
		Protocol:    target.Protocol,
		LoadBalance: loadBalance,
	})
}

func (r *CachedResolver) recordWatchEvent(event Event) {
	if r == nil || r.observer == nil || !r.metricsEnabled {
		return
	}
	r.observer.RecordDiscoveryWatchEvent(context.Background(), observability.DiscoveryRequest{
		Resolver:    "registry",
		Operation:   "watch",
		Namespace:   r.target.Namespace,
		Service:     r.target.Service,
		Protocol:    r.target.Protocol,
		LoadBalance: r.loadBalance,
	}, event.Type, r.endpointCount())
}

func (r *CachedResolver) endpointCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.endpoints)
}

func (r *CachedResolver) cacheAge() time.Duration {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.updatedAt.IsZero() {
		return 0
	}
	return time.Since(r.updatedAt)
}

func (r *CachedResolver) isStale() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.staleTTL > 0 && !r.updatedAt.IsZero() && time.Since(r.updatedAt) > r.staleTTL
}

func sameEndpoint(left Endpoint, right Endpoint) bool {
	left = NormalizeEndpoint(left)
	right = NormalizeEndpoint(right)
	return left.InstanceID == right.InstanceID &&
		left.Name == right.Name &&
		left.Protocol == right.Protocol &&
		left.Host == right.Host &&
		left.Port == right.Port
}

func cloneEndpoints(values []Endpoint) []Endpoint {
	copied := make([]Endpoint, 0, len(values))
	for _, value := range values {
		copied = append(copied, NormalizeEndpoint(value))
	}
	return copied
}

func sleepOrDone(ctx context.Context, duration time.Duration) bool {
	if duration <= 0 {
		duration = time.Second
	}
	timer := time.NewTimer(duration)
	defer timer.Stop()
	select {
	case <-timer.C:
		return true
	case <-ctx.Done():
		return false
	}
}

func IsNoAvailableEndpoint(err error) bool {
	return errors.Is(err, ErrNoAvailableEndpoint)
}
