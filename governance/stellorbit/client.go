package stellorbit

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/stellhub/stellar/config"
	"github.com/stellhub/stellar/governance"
	orbitsdk "github.com/stellhub/stellorbit-go-sdk"
	orbitgovernance "github.com/stellhub/stellorbit-go-sdk/governance"
)

const DefaultName = "governance-stellorbit"

var ErrSyncTargetRequired = errors.New("stellar: governance sync requires store and stellorbit client")

type Client = orbitsdk.Client

type Logger struct {
	Logger *slog.Logger
}

func (l Logger) Printf(format string, args ...any) {
	if l.Logger == nil {
		return
	}
	l.Logger.Info(fmt.Sprintf(format, args...))
}

func NewClientFromConfig(cfg *config.GovernanceConfig, logger *slog.Logger) (*Client, error) {
	if cfg == nil {
		return nil, fmt.Errorf("stellar: governance config is required")
	}
	timeout, err := time.ParseDuration(cfg.RequestTimeout)
	if err != nil {
		return nil, fmt.Errorf("stellar: invalid governance request_timeout %q: %w", cfg.RequestTimeout, err)
	}
	options := orbitsdk.Options{
		Endpoint:                 cfg.Endpoint,
		APIKey:                   cfg.APIKey,
		Timeout:                  timeout,
		StellnulaEndpoint:        cfg.ConfigCenter.Endpoint,
		StellnulaGRPCEndpoint:    cfg.ConfigCenter.GRPCEndpoint,
		StellnulaGRPCPlaintext:   cfg.ConfigCenter.GRPCPlaintext,
		StellnulaAPIToken:        cfg.ConfigCenter.APIToken,
		AppID:                    cfg.AppID,
		ClientID:                 cfg.ClientID,
		Env:                      cfg.Env,
		Region:                   cfg.Region,
		Zone:                     cfg.Zone,
		Cluster:                  cfg.Cluster,
		RuleNamespace:            cfg.ConfigCenter.Namespace,
		RuleGroup:                cfg.ConfigCenter.Group,
		WatchEnabled:             cfg.ConfigCenter.WatchEnabled,
		FailFastOnBootstrap:      cfg.ConfigCenter.FailFastOnBootstrap,
		SnapshotDirectory:        cfg.ConfigCenter.SnapshotDirectory,
		Labels:                   cfg.ConfigCenter.Labels,
		AcceptLargeFileReference: cfg.ConfigCenter.AcceptLargeFileReference,
		Logger:                   Logger{Logger: logger},
	}
	return orbitsdk.NewClient(options)
}

func SyncStore(ctx context.Context, store *governance.Store, client *Client) (governance.Snapshot, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if store == nil || client == nil {
		return governance.Snapshot{}, ErrSyncTargetRequired
	}
	if err := ctx.Err(); err != nil {
		return governance.Snapshot{}, err
	}
	snapshot := SnapshotFromRegistry(client.Rules())
	if err := ctx.Err(); err != nil {
		return governance.Snapshot{}, err
	}
	store.Replace(snapshot)
	return snapshot, nil
}

func SnapshotFromRegistry(registry orbitgovernance.Registry) governance.Snapshot {
	rules := make([]governance.Rule, 0, len(registry.Rules))
	for _, rule := range registry.Rules {
		converted, ok := ConvertRule(rule)
		if ok {
			rules = append(rules, converted)
		}
	}
	version := registry.Checksum
	if version == "" && registry.Revision > 0 {
		version = strconv.FormatInt(registry.Revision, 10)
	}
	return governance.Snapshot{
		Version:   version,
		UpdatedAt: time.Now(),
		Rules:     rules,
	}
}

func ConvertRule(rule orbitgovernance.Rule) (governance.Rule, bool) {
	kind, ok := convertRuleKind(rule)
	if !ok {
		return governance.Rule{}, false
	}
	spec := normalizeSpec(rule.Content)
	metadata := map[string]string{
		"source":         "stellorbit",
		"config_key":     rule.ConfigKey,
		"rule_type":      string(rule.RuleType),
		"rule_name":      rule.RuleName,
		"rule_revision":  strconv.FormatInt(rule.Revision, 10),
		"rule_checksum":  rule.Checksum,
		"target_service": rule.TargetService,
	}
	version := strconv.FormatInt(rule.Revision, 10)
	if rule.Checksum != "" {
		version = rule.Checksum
	}
	return governance.Rule{
		ID:       rule.RuleID,
		Kind:     kind,
		Enabled:  rule.Active(),
		Scope:    scopeFromRule(rule, spec),
		Priority: rule.Priority,
		Version:  version,
		Metadata: metadata,
		Spec:     spec,
	}, true
}

func convertRuleKind(rule orbitgovernance.Rule) (governance.RuleKind, bool) {
	switch rule.RuleType {
	case orbitgovernance.RuleTypeRoute:
		return governance.RuleKindRoute, true
	case orbitgovernance.RuleTypeCircuitBreaker, orbitgovernance.RuleTypeDegrade:
		return governance.RuleKindCircuitBreaker, true
	case orbitgovernance.RuleTypeRateLimit:
		return governance.RuleKindRateLimit, true
	case orbitgovernance.RuleTypeAuth:
		return authRuleKind(rule.Content), true
	default:
		return "", false
	}
}

func authRuleKind(content map[string]any) governance.RuleKind {
	value := strings.ToLower(strings.TrimSpace(fmt.Sprint(firstNonNil(
		content["kind"],
		content["auth_kind"],
		content["policy"],
		content["type"],
	))))
	switch strings.ReplaceAll(value, "-", "_") {
	case "authentication", "authenticate", "authn":
		return governance.RuleKindAuthentication
	case "signing", "signature", "sign":
		return governance.RuleKindSigning
	default:
		return governance.RuleKindAuthorization
	}
}

func scopeFromRule(rule orbitgovernance.Rule, spec map[string]any) governance.Scope {
	return governance.Scope{
		Transport: stringValue(firstNonNil(spec["transport"], nested(spec, "scope", "transport"))),
		Service:   defaultString(stringValue(firstNonNil(spec["service"], nested(spec, "scope", "service"))), rule.TargetService),
		Method:    stringValue(firstNonNil(spec["method"], nested(spec, "scope", "method"))),
		Path:      stringValue(firstNonNil(spec["path"], nested(spec, "scope", "path"))),
		Target:    stringValue(firstNonNil(spec["target"], nested(spec, "scope", "target"))),
	}
}

func normalizeSpec(content map[string]any) map[string]any {
	spec := copyAnyMap(content)
	promoteMap(spec, "scope")
	promoteMap(spec, "filter")
	promoteMap(spec, "route")
	promoteMap(spec, "limit")
	promoteMap(spec, "rateLimit")
	return spec
}

func promoteMap(target map[string]any, key string) {
	nestedValue := mapValue(target[key])
	for nestedKey, nestedItem := range nestedValue {
		if _, exists := target[nestedKey]; !exists {
			target[nestedKey] = nestedItem
		}
	}
}

func copyAnyMap(values map[string]any) map[string]any {
	if len(values) == 0 {
		return map[string]any{}
	}
	copied := make(map[string]any, len(values))
	for key, value := range values {
		copied[key] = deepCopyAny(value)
	}
	return copied
}

func deepCopyAny(value any) any {
	switch typed := value.(type) {
	case map[string]any:
		return copyAnyMap(typed)
	case map[string]string:
		copied := make(map[string]any, len(typed))
		for key, item := range typed {
			copied[key] = item
		}
		return copied
	case []any:
		copied := make([]any, 0, len(typed))
		for _, item := range typed {
			copied = append(copied, deepCopyAny(item))
		}
		return copied
	default:
		return typed
	}
}

func nested(content map[string]any, first string, second string) any {
	return mapValue(content[first])[second]
}

func mapValue(value any) map[string]any {
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

func firstNonNil(values ...any) any {
	for _, value := range values {
		if value != nil {
			return value
		}
	}
	return nil
}

func stringValue(value any) string {
	if value == nil {
		return ""
	}
	return strings.TrimSpace(fmt.Sprint(value))
}

func defaultString(value string, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return strings.TrimSpace(fallback)
	}
	return strings.TrimSpace(value)
}
