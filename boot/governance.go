package boot

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/stellhub/stellar/config"
	"github.com/stellhub/stellar/governance"
	governanceauth "github.com/stellhub/stellar/governance/auth"
	governancebreaker "github.com/stellhub/stellar/governance/circuitbreaker"
	governanceratelimit "github.com/stellhub/stellar/governance/ratelimit"
	governanceorbit "github.com/stellhub/stellar/governance/stellorbit"
	stellpulsar "github.com/stellhub/stellpulsar-go-sdk"
)

func configureGovernanceStarter(ctx context.Context, app *App, cfg config.Config) error {
	if cfg.Governance == nil || !enabledByDefault(cfg.Governance.Enabled) {
		return nil
	}
	if err := validateGovernanceConfig(cfg.Governance); err != nil {
		return err
	}
	app.governanceRouteEnabled = featureEnabled(cfg.Governance.Route.Enabled)
	metrics := governance.NewMetrics(app.observability)
	app.governanceMetrics = metrics
	orbitClient, err := governanceorbit.NewClientFromConfig(cfg.Governance, app.logger)
	if err != nil {
		return fmt.Errorf("configure stellorbit governance client: %w", err)
	}
	app.registry.Set(governanceorbit.DefaultName, orbitClient)
	app.RegisterStarters(&stellorbitGovernanceStarter{
		cfg:          *cfg.Governance,
		client:       orbitClient,
		store:        app.governance,
		metrics:      metrics,
		syncInterval: governanceSyncInterval(cfg.Governance),
	})
	if featureEnabled(cfg.Governance.Route.Enabled) {
		app.RegisterStarters(&governanceRouteStarter{})
	}
	if featureEnabled(cfg.Governance.CircuitBreaker.Enabled) {
		app.RegisterStarters(&governanceCircuitBreakerStarter{metrics: metrics})
	}
	if featureEnabled(cfg.Governance.RateLimit.Enabled) {
		starter, err := newGovernanceRateLimitStarter(ctx, cfg.Governance, orbitClient, metrics)
		if err != nil {
			return err
		}
		app.RegisterStarters(starter)
	}
	if featureEnabled(cfg.Governance.Auth.Enabled) {
		app.RegisterStarters(&governanceAuthStarter{
			cfg:     cfg.Governance.Auth,
			metrics: metrics,
			logger:  app.logger,
		})
	}
	return nil
}

func validateGovernanceConfig(cfg *config.GovernanceConfig) error {
	if cfg == nil {
		return nil
	}
	if !strings.EqualFold(cfg.Adapter, "stellorbit") {
		return fmt.Errorf("stellar: unsupported governance adapter %q", cfg.Adapter)
	}
	if !strings.EqualFold(cfg.ConfigCenter.Adapter, "stellnula") {
		return fmt.Errorf("stellar: stellorbit governance requires config_center adapter stellnula")
	}
	if strings.TrimSpace(cfg.ConfigCenter.Endpoint) == "" {
		return fmt.Errorf("stellar: stellorbit governance requires config_center endpoint")
	}
	return nil
}

func newGovernanceRateLimitStarter(ctx context.Context, cfg *config.GovernanceConfig, orbitClient *governanceorbit.Client, metrics *governance.Metrics) (*governanceRateLimitStarter, error) {
	configRules, err := configuredHeaderRateLimitRules(cfg)
	if err != nil {
		return nil, err
	}
	var pulsarClient stellpulsar.StellpulsarClient
	var distributed governanceLimiter
	if distributedRateLimitClientEnabled(cfg.RateLimit) {
		if strings.TrimSpace(cfg.RateLimit.Distributed.Address) == "" {
			return nil, fmt.Errorf("stellar: distributed rate limit requires governance.rate_limit.distributed.address")
		}
		timeout, err := time.ParseDuration(cfg.RateLimit.Distributed.Timeout)
		if err != nil {
			return nil, fmt.Errorf("stellar: invalid distributed rate limit timeout %q: %w", cfg.RateLimit.Distributed.Timeout, err)
		}
		retryDelay, err := time.ParseDuration(cfg.RateLimit.Distributed.RetryDelay)
		if err != nil {
			return nil, fmt.Errorf("stellar: invalid distributed rate limit retry_delay %q: %w", cfg.RateLimit.Distributed.RetryDelay, err)
		}
		client, err := stellpulsar.NewClientWithContext(ctx, stellpulsar.Options{
			ApplicationCode:    cfg.AppID,
			ClientID:           cfg.ClientID,
			Namespace:          cfg.ConfigCenter.Namespace,
			ServiceName:        "stellpulsar-service",
			StellpulsarAddress: cfg.RateLimit.Distributed.Address,
			APIToken:           cfg.RateLimit.Distributed.APIToken,
			GRPCPlaintext:      cfg.RateLimit.Distributed.GRPCPlaintext,
			GRPCDeadline:       timeout,
			MaxAcquireAttempts: cfg.RateLimit.Distributed.MaxAcquireAttempts,
			RetryDelay:         retryDelay,
			DefaultFailPolicy:  stellpulsar.ParseFailPolicy(cfg.RateLimit.Distributed.Fallback, stellpulsar.FailOpen),
			Labels:             cfg.Labels,
			StellorbitClient:   orbitClient,
		})
		if err != nil {
			return nil, fmt.Errorf("configure stellpulsar distributed rate limit client: %w", err)
		}
		pulsarClient = client
		distributed = governanceratelimit.NewDistributedLimiter(governanceratelimit.DistributedOptions{
			Client:   client,
			Timeout:  timeout,
			Fallback: cfg.RateLimit.Distributed.Fallback,
			AppID:    cfg.AppID,
		})
	}
	return &governanceRateLimitStarter{
		cfg:          cfg.RateLimit,
		metrics:      metrics,
		distributed:  distributed,
		pulsarClient: pulsarClient,
		configRules:  configRules,
	}, nil
}

type governanceLimiter = governanceratelimit.Limiter

type stellorbitGovernanceStarter struct {
	cfg          config.GovernanceConfig
	client       *governanceorbit.Client
	store        *governance.Store
	metrics      *governance.Metrics
	syncInterval time.Duration
	cancel       context.CancelFunc
}

func (s *stellorbitGovernanceStarter) Name() string {
	return governanceorbit.DefaultName
}

func (s *stellorbitGovernanceStarter) Condition(StarterContext) bool {
	return true
}

func (s *stellorbitGovernanceStarter) Init(context.Context, *App) error {
	return nil
}

func (s *stellorbitGovernanceStarter) Start(ctx context.Context) error {
	if err := s.client.Start(ctx); err != nil {
		s.recordSync(ctx, "error")
		return err
	}
	if err := s.initialSyncError(s.sync(ctx)); err != nil {
		return err
	}
	loopCtx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel
	go s.syncLoop(loopCtx)
	return nil
}

func (s *stellorbitGovernanceStarter) initialSyncError(err error) error {
	if err == nil {
		return nil
	}
	if s.cfg.FailFastOnBootstrap {
		return err
	}
	return nil
}

func (s *stellorbitGovernanceStarter) Stop(context.Context) error {
	if s.cancel != nil {
		s.cancel()
	}
	if s.client == nil {
		return nil
	}
	return s.client.Close()
}

func (s *stellorbitGovernanceStarter) syncLoop(ctx context.Context) {
	interval := s.syncInterval
	if interval <= 0 {
		interval = 5 * time.Second
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			_ = s.sync(ctx)
		case <-ctx.Done():
			return
		}
	}
}

func (s *stellorbitGovernanceStarter) sync(ctx context.Context) error {
	if _, err := governanceorbit.SyncStore(ctx, s.store, s.client); err != nil {
		s.recordSync(ctx, "error")
		return err
	}
	s.recordSync(ctx, "ok")
	return nil
}

func (s *stellorbitGovernanceStarter) recordSync(ctx context.Context, outcome string) {
	if s.metrics == nil {
		return
	}
	s.metrics.RecordRuleSync(ctx, governance.MetricAttrs{
		Adapter: "stellorbit",
		Outcome: outcome,
	})
}

type governanceRouteStarter struct{}

func (s *governanceRouteStarter) Name() string {
	return "governance-route"
}

func (s *governanceRouteStarter) Condition(StarterContext) bool {
	return true
}

func (s *governanceRouteStarter) Init(context.Context, *App) error {
	return nil
}

func (s *governanceRouteStarter) Start(context.Context) error {
	return nil
}

func (s *governanceRouteStarter) Stop(context.Context) error {
	return nil
}

type governanceCircuitBreakerStarter struct {
	metrics *governance.Metrics
}

func (s *governanceCircuitBreakerStarter) Name() string {
	return governancebreaker.Name
}

func (s *governanceCircuitBreakerStarter) Condition(StarterContext) bool {
	return true
}

func (s *governanceCircuitBreakerStarter) Init(_ context.Context, app *App) error {
	app.interceptors.Register(governancebreaker.Definitions(app.governance, s.metrics)...)
	return nil
}

func (s *governanceCircuitBreakerStarter) Start(context.Context) error {
	return nil
}

func (s *governanceCircuitBreakerStarter) Stop(context.Context) error {
	return nil
}

type governanceRateLimitStarter struct {
	cfg          config.GovernanceRateLimitConfig
	metrics      *governance.Metrics
	distributed  governanceLimiter
	pulsarClient stellpulsar.StellpulsarClient
	configRules  []governance.Rule
}

func (s *governanceRateLimitStarter) Name() string {
	return governanceratelimit.Name
}

func (s *governanceRateLimitStarter) Condition(StarterContext) bool {
	return true
}

func (s *governanceRateLimitStarter) Init(_ context.Context, app *App) error {
	app.interceptors.Register(governanceratelimit.Definitions(governanceratelimit.PolicyOptions{
		Store:               app.governance,
		Rules:               s.configRules,
		Metrics:             s.metrics,
		Distributed:         s.distributed,
		DefaultMode:         s.cfg.Mode,
		DefaultBehavior:     governanceratelimit.LimitBehavior(s.cfg.Behavior),
		DistributedFallback: s.cfg.Distributed.Fallback,
		DefaultRate:         s.cfg.Local.DefaultRate,
		DefaultBurst:        s.cfg.Local.DefaultBurst,
	})...)
	return nil
}

func (s *governanceRateLimitStarter) Start(ctx context.Context) error {
	if s.pulsarClient == nil {
		return nil
	}
	return s.pulsarClient.Start(ctx)
}

func (s *governanceRateLimitStarter) Stop(context.Context) error {
	if s.pulsarClient == nil {
		return nil
	}
	return s.pulsarClient.Close()
}

type governanceAuthStarter struct {
	cfg     config.GovernanceAuthConfig
	metrics *governance.Metrics
	logger  *slog.Logger
}

func (s *governanceAuthStarter) Name() string {
	return governanceauth.Name
}

func (s *governanceAuthStarter) Condition(StarterContext) bool {
	return true
}

func (s *governanceAuthStarter) Init(_ context.Context, app *App) error {
	if !s.cfg.KeyProvider.Placeholder {
		return fmt.Errorf("stellar: governance auth requires stellguard-agent key provider; set placeholder: true for development placeholder mode")
	}
	if s.logger != nil {
		s.logger.Warn("stellar governance auth uses placeholder key provider", "adapter", s.cfg.KeyProvider.Adapter)
	}
	app.interceptors.Register(governanceauth.Definitions(app.governance, s.metrics)...)
	return nil
}

func (s *governanceAuthStarter) Start(context.Context) error {
	return nil
}

func (s *governanceAuthStarter) Stop(context.Context) error {
	return nil
}

func governanceSyncInterval(cfg *config.GovernanceConfig) time.Duration {
	if cfg == nil || strings.TrimSpace(cfg.SyncInterval) == "" {
		return 5 * time.Second
	}
	duration, err := time.ParseDuration(cfg.SyncInterval)
	if err != nil || duration <= 0 {
		return 5 * time.Second
	}
	return duration
}

func featureEnabled(value *bool) bool {
	return value != nil && *value
}

func enabledByDefault(value *bool) bool {
	return value == nil || *value
}

func distributedRateLimitClientEnabled(cfg config.GovernanceRateLimitConfig) bool {
	if rateLimitDefaultModeRequiresDistributed(cfg.Mode) {
		return true
	}
	for _, header := range cfg.Headers {
		if enabledByDefault(header.Enabled) && headerRateLimitRequiresDistributed(header) {
			return true
		}
	}
	if cfg.Distributed.Enabled != nil {
		return *cfg.Distributed.Enabled
	}
	if strings.TrimSpace(cfg.Distributed.Address) != "" {
		return true
	}
	return false
}

func configuredHeaderRateLimitRules(cfg *config.GovernanceConfig) ([]governance.Rule, error) {
	if cfg == nil || len(cfg.RateLimit.Headers) == 0 {
		return nil, nil
	}
	rules := make([]governance.Rule, 0, len(cfg.RateLimit.Headers))
	for index, header := range cfg.RateLimit.Headers {
		if !enabledByDefault(header.Enabled) {
			continue
		}
		if strings.TrimSpace(header.Header) == "" {
			return nil, fmt.Errorf("stellar: governance.rate_limit.headers[%d].header is required", index)
		}
		transport := header.Transport
		if transport == "" {
			transport = "http"
		}
		transport = strings.ToLower(strings.TrimSpace(transport))
		if transport != "http" && transport != "grpc" {
			return nil, fmt.Errorf("stellar: governance.rate_limit.headers[%d].transport must be http or grpc", index)
		}
		source := "HEADER"
		limitType := "HEADER"
		if transport == "grpc" {
			source = "GRPC_METADATA"
			limitType = "GRPC_METADATA"
		}
		coordinationMode, err := rateLimitCoordinationMode(header.CoordinationMode)
		if err != nil {
			return nil, fmt.Errorf("stellar: governance.rate_limit.headers[%d].coordination_mode: %w", index, err)
		}
		ruleID := header.Name
		if ruleID == "" {
			ruleID = "config-header-" + transport + "-" + normalizeRuleIDToken(header.Header)
		}
		rate := header.Rate
		if rate <= 0 {
			rate = cfg.RateLimit.Local.DefaultRate
		}
		burst := header.Burst
		if burst <= 0 {
			burst = rate
		}
		spec := map[string]any{
			"limitMode":         "HEADER",
			"limit_mode":        "header",
			"limitType":         limitType,
			"trafficProtocol":   strings.ToUpper(transport),
			"executionLocation": "APPLICATION",
			"coordinationMode":  coordinationMode,
			"mode":              rateLimitModeForCoordination(coordinationMode),
			"resource":          "header:" + transport + ":" + strings.ToLower(strings.TrimSpace(header.Header)),
			"rate":              rate,
			"burst":             burst,
			"behavior":          header.Behavior,
			"quotaConfig": map[string]any{
				"limit": rate,
			},
			"burstConfig": map[string]any{
				"capacity": burst,
			},
			"keyExtractor": map[string]any{
				"keys": []any{map[string]any{
					"name":      strings.TrimSpace(header.Header),
					"source":    source,
					"key":       strings.TrimSpace(header.Header),
					"required":  headerRequired(header.Required),
					"normalize": header.Normalize,
				}},
			},
		}
		rules = append(rules, governance.Rule{
			ID:      ruleID,
			Kind:    governance.RuleKindRateLimit,
			Enabled: true,
			Scope: governance.Scope{
				Transport: transport + ".server",
				Service:   header.Service,
				Method:    header.Method,
				Path:      header.Path,
			},
			Priority: 1000 + index,
			Version:  "application.yaml",
			Metadata: map[string]string{
				"source": "application.yaml",
			},
			Spec: spec,
		})
	}
	return rules, nil
}

func rateLimitCoordinationMode(value string) (string, error) {
	switch strings.ToLower(strings.ReplaceAll(strings.TrimSpace(value), "_", "-")) {
	case "", "local", "local-only":
		return "LOCAL_ONLY", nil
	case "distributed", "global", "global-sync":
		return "GLOBAL_SYNC", nil
	case "global-quota":
		return "GLOBAL_QUOTA", nil
	default:
		return "", fmt.Errorf("unsupported coordination mode %q", value)
	}
}

func rateLimitModeForCoordination(value string) string {
	mode, err := rateLimitCoordinationMode(value)
	if err != nil {
		return ""
	}
	switch mode {
	case "GLOBAL_SYNC", "GLOBAL_QUOTA":
		return "distributed"
	default:
		return "local"
	}
}

func rateLimitDefaultModeRequiresDistributed(value string) bool {
	switch strings.ToLower(strings.ReplaceAll(strings.TrimSpace(value), "_", "-")) {
	case "distributed", "global", "global-sync", "global-quota", "edge":
		return true
	default:
		return false
	}
}

func headerRateLimitRequiresDistributed(header config.GovernanceHeaderRateLimitConfig) bool {
	return rateLimitModeForCoordination(header.CoordinationMode) == "distributed"
}

func normalizeRuleIDToken(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	replacer := strings.NewReplacer(".", "-", "_", "-", ":", "-", "/", "-", " ", "-")
	value = replacer.Replace(value)
	value = strings.Trim(value, "-")
	if value == "" {
		return "header"
	}
	return value
}

func headerRequired(value *bool) bool {
	return value != nil && *value
}
