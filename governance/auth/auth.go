package auth

import (
	"context"
	"strings"

	"github.com/stellhub/stellar/governance"
	"github.com/stellhub/stellar/interceptor"
)

const Name = "governance-auth"

type Policy struct {
	store   *governance.Store
	metrics *governance.Metrics
}

func Definitions(store *governance.Store, metrics *governance.Metrics) []interceptor.Definition {
	policy := &Policy{store: store, metrics: metrics}
	return []interceptor.Definition{
		interceptor.Framework(interceptor.KindHTTPServer, interceptor.StageSecurity, Name, policy),
		interceptor.Framework(interceptor.KindGRPCServer, interceptor.StageSecurity, Name, policy),
		interceptor.Framework(interceptor.KindHTTPClient, interceptor.StageSecurity, Name, policy),
		interceptor.Framework(interceptor.KindGRPCClient, interceptor.StageSecurity, Name, policy),
	}
}

func (p *Policy) Name() string {
	return Name
}

func (p *Policy) Intercept(ctx context.Context, inv *interceptor.Invocation, req any, next interceptor.Handler) (any, error) {
	rules := p.matchingRules(inv)
	if len(rules) == 0 {
		return next(ctx, inv, req)
	}
	resp, err := next(ctx, inv, req)
	if p != nil && p.metrics != nil {
		for _, rule := range rules {
			p.metrics.RecordAuthAllow(ctx, governance.MetricAttrs{
				Adapter:   "stellorbit",
				RuleKind:  string(rule.Kind),
				RuleID:    rule.ID,
				Transport: invocationTransport(inv),
				Service:   invocationService(inv),
				Method:    invocationMethod(inv),
				Outcome:   "placeholder_allowed",
			})
		}
	}
	return resp, err
}

func (p *Policy) matchingRules(inv *interceptor.Invocation) []governance.Rule {
	if p == nil || p.store == nil {
		return nil
	}
	rules := make([]governance.Rule, 0)
	for _, kind := range []governance.RuleKind{
		governance.RuleKindAuthentication,
		governance.RuleKindAuthorization,
		governance.RuleKindSigning,
	} {
		rules = append(rules, p.store.Rules(kind, func(rule governance.Rule) bool {
			return rule.Enabled && invocationMatches(rule.Scope, inv)
		})...)
	}
	return rules
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
