package governance

import (
	"sort"
	"strings"
	"sync/atomic"
	"time"
)

type RuleKind string

const (
	RuleKindRoute            RuleKind = "route"
	RuleKindRateLimit        RuleKind = "rate_limit"
	RuleKindCircuitBreaker   RuleKind = "circuit_breaker"
	RuleKindLoadShedding     RuleKind = "load_shedding"
	RuleKindConcurrencyLimit RuleKind = "concurrency_limit"
	RuleKindAuthentication   RuleKind = "authentication"
	RuleKindAuthorization    RuleKind = "authorization"
	RuleKindRetry            RuleKind = "retry"
	RuleKindSigning          RuleKind = "signing"
	RuleKindQuota            RuleKind = "quota"
)

type Scope struct {
	Transport string
	Service   string
	Method    string
	Path      string
	Target    string
}

type Rule struct {
	ID       string
	Kind     RuleKind
	Enabled  bool
	Scope    Scope
	Priority int
	Version  string
	Metadata map[string]string
	Spec     map[string]any
}

type Snapshot struct {
	Version   string
	UpdatedAt time.Time
	Rules     []Rule
}

type Store struct {
	value atomic.Value
}

func NewStore(initial ...Snapshot) *Store {
	store := &Store{}
	snapshot := Snapshot{UpdatedAt: time.Now()}
	if len(initial) > 0 {
		snapshot = initial[0]
	}
	store.Replace(snapshot)
	return store
}

func (s *Store) Snapshot() Snapshot {
	if s == nil {
		return Snapshot{}
	}
	value := s.value.Load()
	if value == nil {
		return Snapshot{}
	}
	snapshot, _ := value.(Snapshot)
	return cloneSnapshot(snapshot)
}

func (s *Store) Replace(snapshot Snapshot) {
	if s == nil {
		return
	}
	if snapshot.UpdatedAt.IsZero() {
		snapshot.UpdatedAt = time.Now()
	}
	s.value.Store(normalizeSnapshot(snapshot))
}

func (s *Store) Rules(kind RuleKind, match func(Rule) bool) []Rule {
	snapshot := s.Snapshot()
	rules := make([]Rule, 0)
	for _, rule := range snapshot.Rules {
		if kind != "" && rule.Kind != kind {
			continue
		}
		if match != nil && !match(rule) {
			continue
		}
		rules = append(rules, rule)
	}
	return rules
}

func normalizeSnapshot(snapshot Snapshot) Snapshot {
	snapshot = cloneSnapshot(snapshot)
	sort.SliceStable(snapshot.Rules, func(i, j int) bool {
		if snapshot.Rules[i].Priority == snapshot.Rules[j].Priority {
			return snapshot.Rules[i].ID < snapshot.Rules[j].ID
		}
		return snapshot.Rules[i].Priority < snapshot.Rules[j].Priority
	})
	return snapshot
}

func cloneSnapshot(snapshot Snapshot) Snapshot {
	rules := make([]Rule, 0, len(snapshot.Rules))
	for _, rule := range snapshot.Rules {
		rules = append(rules, cloneRule(rule))
	}
	snapshot.Rules = rules
	return snapshot
}

func cloneRule(rule Rule) Rule {
	rule.ID = strings.TrimSpace(rule.ID)
	if rule.Metadata != nil {
		metadata := make(map[string]string, len(rule.Metadata))
		for key, value := range rule.Metadata {
			metadata[key] = value
		}
		rule.Metadata = metadata
	}
	if rule.Spec != nil {
		spec := make(map[string]any, len(rule.Spec))
		for key, value := range rule.Spec {
			spec[key] = value
		}
		rule.Spec = spec
	}
	return rule
}
