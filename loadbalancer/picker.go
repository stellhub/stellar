package loadbalancer

import (
	"context"
	"hash/fnv"
	"math"
	"math/rand/v2"
	"sort"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/stellhub/stellar/discovery"
)

type Picker struct {
	policy string
	next   atomic.Uint64
	stats  sync.Map
}

func New(policy string) *Picker {
	return &Picker{policy: NormalizePolicy(policy)}
}

func (p *Picker) Pick(_ context.Context, request Request, endpoints []discovery.Endpoint) (Pick, error) {
	if len(endpoints) == 0 {
		return Pick{}, discovery.ErrNoAvailableEndpoint
	}
	endpoints = CloneEndpoints(endpoints)
	var endpoint discovery.Endpoint
	switch p.policy {
	case PolicyRoundRobin:
		endpoint = p.pickRoundRobin(endpoints)
	case PolicyRandom:
		endpoint = endpoints[rand.IntN(len(endpoints))]
	case PolicyWeightedRound:
		endpoint = p.pickWeightedRoundRobin(endpoints)
	case PolicyLeastRequest:
		endpoint = p.pickLeastRequest(endpoints)
	case PolicyConsistentHash:
		endpoint = p.pickConsistentHash(request, endpoints)
	default:
		endpoint = p.pickP2C(endpoints)
	}
	state := p.state(endpoint)
	state.start()
	var doneOnce sync.Once
	return Pick{
		Endpoint: endpoint,
		Done: func(result Result) {
			doneOnce.Do(func() {
				state.done(result)
			})
		},
	}, nil
}

func (p *Picker) pickRoundRobin(endpoints []discovery.Endpoint) discovery.Endpoint {
	index := p.next.Add(1)
	return endpoints[int((index-1)%uint64(len(endpoints)))]
}

func (p *Picker) pickWeightedRoundRobin(endpoints []discovery.Endpoint) discovery.Endpoint {
	total := 0
	for _, endpoint := range endpoints {
		total += endpointWeight(endpoint)
	}
	if total <= 0 {
		return endpoints[0]
	}
	cursor := int(p.next.Add(1)%uint64(total)) + 1
	for _, endpoint := range endpoints {
		cursor -= endpointWeight(endpoint)
		if cursor <= 0 {
			return endpoint
		}
	}
	return endpoints[len(endpoints)-1]
}

func (p *Picker) pickLeastRequest(endpoints []discovery.Endpoint) discovery.Endpoint {
	best := endpoints[0]
	bestScore := p.score(best)
	for _, endpoint := range endpoints[1:] {
		score := p.score(endpoint)
		if score < bestScore {
			best = endpoint
			bestScore = score
		}
	}
	return best
}

func (p *Picker) pickP2C(endpoints []discovery.Endpoint) discovery.Endpoint {
	if len(endpoints) == 1 {
		return endpoints[0]
	}
	firstIndex := rand.IntN(len(endpoints))
	secondIndex := rand.IntN(len(endpoints) - 1)
	if secondIndex >= firstIndex {
		secondIndex++
	}
	first := endpoints[firstIndex]
	second := endpoints[secondIndex]
	if p.score(second) < p.score(first) {
		return second
	}
	return first
}

func (p *Picker) pickConsistentHash(request Request, endpoints []discovery.Endpoint) discovery.Endpoint {
	key := request.HashKey
	if key == "" {
		key = HeaderValue(request.Headers, "x-stellar-lb-key")
	}
	if key == "" {
		key = request.Target
	}
	if key == "" {
		key = request.Path
	}
	if key == "" {
		return p.pickP2C(endpoints)
	}
	ring := buildHashRing(endpoints)
	if len(ring) == 0 {
		return p.pickP2C(endpoints)
	}
	hash := hashString(key)
	index := sort.Search(len(ring), func(i int) bool {
		return ring[i].hash >= hash
	})
	if index == len(ring) {
		index = 0
	}
	return ring[index].endpoint
}

func (p *Picker) score(endpoint discovery.Endpoint) float64 {
	state := p.state(endpoint)
	weight := float64(endpointWeight(endpoint))
	active := float64(state.active.Load())
	latency := float64(state.latencyMicros()) / float64(time.Millisecond/time.Microsecond)
	errors := float64(state.errors.Load())
	return (active*1000 + latency + errors*100) / weight
}

func (p *Picker) state(endpoint discovery.Endpoint) *endpointState {
	key := EndpointKey(endpoint)
	value, _ := p.stats.LoadOrStore(key, &endpointState{})
	state, ok := value.(*endpointState)
	if !ok {
		return &endpointState{}
	}
	return state
}

type endpointState struct {
	active atomic.Int64
	errors atomic.Uint64
	mu     sync.RWMutex
	ewma   float64
}

func (s *endpointState) start() {
	s.active.Add(1)
}

func (s *endpointState) done(result Result) {
	s.active.Add(-1)
	if result.Err != nil || result.StatusCode >= 500 {
		s.errors.Add(1)
	}
	if result.Duration <= 0 {
		return
	}
	micros := float64(result.Duration.Microseconds())
	s.mu.Lock()
	if s.ewma == 0 {
		s.ewma = micros
	} else {
		s.ewma = s.ewma*0.8 + micros*0.2
	}
	s.mu.Unlock()
}

func (s *endpointState) latencyMicros() int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.ewma <= 0 || math.IsNaN(s.ewma) || math.IsInf(s.ewma, 0) {
		return 0
	}
	return int64(s.ewma)
}

func endpointWeight(endpoint discovery.Endpoint) int {
	if endpoint.Weight <= 0 {
		return 100
	}
	return endpoint.Weight
}

type hashRingNode struct {
	hash     uint64
	endpoint discovery.Endpoint
}

func buildHashRing(endpoints []discovery.Endpoint) []hashRingNode {
	ring := make([]hashRingNode, 0, len(endpoints)*10)
	for _, endpoint := range endpoints {
		endpoint = discovery.NormalizeEndpoint(endpoint)
		replicas := endpointWeight(endpoint) / 10
		if replicas < 1 {
			replicas = 1
		}
		if replicas > 256 {
			replicas = 256
		}
		key := EndpointKey(endpoint)
		for i := 0; i < replicas; i++ {
			ring = append(ring, hashRingNode{
				hash:     hashString(key + "#" + strconv.Itoa(i)),
				endpoint: endpoint,
			})
		}
	}
	sort.Slice(ring, func(i, j int) bool {
		return ring[i].hash < ring[j].hash
	})
	return ring
}

func hashString(value string) uint64 {
	hash := fnv.New64a()
	_, _ = hash.Write([]byte(value))
	return hash.Sum64()
}
