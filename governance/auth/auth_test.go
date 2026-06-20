package auth

import (
	"testing"

	"github.com/stellhub/stellar/governance"
	"github.com/stellhub/stellar/interceptor"
)

func TestPolicyMatchesOnlyAuthRulesInScope(t *testing.T) {
	store := governance.NewStore(governance.Snapshot{
		Rules: []governance.Rule{
			{
				ID:      "authz-orders",
				Kind:    governance.RuleKindAuthorization,
				Enabled: true,
				Scope: governance.Scope{
					Transport: "http.server",
					Service:   "order-service",
					Path:      "/orders",
				},
			},
			{
				ID:      "authn-users",
				Kind:    governance.RuleKindAuthentication,
				Enabled: true,
				Scope: governance.Scope{
					Transport: "http.server",
					Service:   "user-service",
				},
			},
			{
				ID:      "route",
				Kind:    governance.RuleKindRoute,
				Enabled: true,
				Scope: governance.Scope{
					Service: "order-service",
				},
			},
		},
	})
	policy := &Policy{store: store}
	rules := policy.matchingRules(&interceptor.Invocation{
		Kind:     interceptor.KindHTTPServer,
		Protocol: "http",
		Service:  "order-service",
		Path:     "/orders/123",
	})

	if len(rules) != 1 || rules[0].ID != "authz-orders" {
		t.Fatalf("unexpected matching auth rules %#v", rules)
	}
}
