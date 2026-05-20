package pool

import (
	"context"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"

	"pswitch/internal/config"
	"pswitch/internal/upstream"
)

type Mode string

const (
	ModeSequential    Mode = "sequential"
	ModeRoundRobin    Mode = "round_robin"
	ModeLeastFailures Mode = "least_failures"
)

type ProviderSnapshot struct {
	Name    string
	BaseURL *url.URL
	APIKey  string
}

type ProviderStatus struct {
	Healthy             bool
	ConsecutiveFailures int
	NextProbeAt         time.Time
}

type ProbeEvent struct {
	Provider string
	Kind     string
}

type Pool struct {
	mu               sync.RWMutex
	mode             Mode
	failureThreshold int
	cooldown         time.Duration
	providers        []*providerState
	indexByName      map[string]int
	cursor           int
}

type providerState struct {
	snapshot            ProviderSnapshot
	healthy             bool
	consecutiveFailures int
	totalFailures       int
	nextProbeAt         time.Time
}

func New(cfg config.Config) (*Pool, error) {
	mode := parseMode(cfg.Mode)
	states := make([]*providerState, 0, len(cfg.Providers))
	indexByName := make(map[string]int, len(cfg.Providers))

	for _, provider := range cfg.Providers {
		baseURL, err := url.Parse(provider.BaseURL)
		if err != nil {
			return nil, err
		}

		indexByName[provider.Name] = len(states)
		states = append(states, &providerState{
			snapshot: ProviderSnapshot{
				Name:    provider.Name,
				BaseURL: baseURL,
				APIKey:  provider.APIKey,
			},
			healthy: true,
		})
	}

	return &Pool{
		mode:             mode,
		failureThreshold: cfg.FailureThreshold,
		cooldown:         cfg.Cooldown,
		providers:        states,
		indexByName:      indexByName,
	}, nil
}

func (p *Pool) Mode() Mode {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.mode
}

func (p *Pool) Candidates(mode Mode, now time.Time) []ProviderSnapshot {
	p.mu.RLock()
	defer p.mu.RUnlock()

	order := p.orderedProvidersLocked(parseMode(string(mode)))
	healthy := make([]ProviderSnapshot, 0, len(order))
	for _, state := range order {
		if state.healthy {
			healthy = append(healthy, state.snapshot)
		}
	}
	if len(healthy) > 0 {
		return healthy
	}

	due := make([]ProviderSnapshot, 0, len(order))
	for _, state := range order {
		if !state.healthy && !state.nextProbeAt.After(now) {
			due = append(due, state.snapshot)
		}
	}
	return due
}

func (p *Pool) MarkSuccess(name string, _ time.Time) {
	p.mu.Lock()
	defer p.mu.Unlock()

	idx, ok := p.indexByName[name]
	if !ok {
		return
	}

	state := p.providers[idx]
	state.healthy = true
	state.consecutiveFailures = 0
	state.nextProbeAt = time.Time{}
	if p.mode == ModeRoundRobin {
		p.cursor = (idx + 1) % len(p.providers)
	}
}

func (p *Pool) MarkFailure(name string, now time.Time) {
	p.mu.Lock()
	defer p.mu.Unlock()

	idx, ok := p.indexByName[name]
	if !ok {
		return
	}

	state := p.providers[idx]
	state.consecutiveFailures++
	state.totalFailures++
	if state.consecutiveFailures >= p.failureThreshold {
		state.healthy = false
		state.nextProbeAt = now.Add(p.cooldown)
	}
}

func (p *Pool) Status(name string) (ProviderStatus, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	idx, ok := p.indexByName[name]
	if !ok {
		return ProviderStatus{}, false
	}

	state := p.providers[idx]
	return ProviderStatus{
		Healthy:             state.healthy,
		ConsecutiveFailures: state.consecutiveFailures,
		NextProbeAt:         state.nextProbeAt,
	}, true
}

func (p *Pool) ProbeDue(ctx context.Context, client *http.Client, now time.Time) []ProbeEvent {
	candidates := p.probeCandidates(now)
	events := make([]ProbeEvent, 0, len(candidates))
	for _, candidate := range candidates {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, buildTargetURL(candidate.BaseURL, "/v1/models", "").String(), nil)
		if err != nil {
			p.MarkFailure(candidate.Name, now)
			continue
		}
		req.Header.Set("Authorization", "Bearer "+candidate.APIKey)

		resp, err := client.Do(req)
		if err != nil {
			p.MarkFailure(candidate.Name, now)
			continue
		}

		_ = resp.Body.Close()
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			p.MarkSuccess(candidate.Name, now)
			events = append(events, ProbeEvent{Provider: candidate.Name, Kind: "recovered"})
			continue
		}
		p.MarkFailure(candidate.Name, now)
	}
	return events
}

func (p *Pool) probeCandidates(now time.Time) []ProviderSnapshot {
	p.mu.RLock()
	defer p.mu.RUnlock()

	out := make([]ProviderSnapshot, 0, len(p.providers))
	for _, state := range p.providers {
		if state.healthy {
			continue
		}
		if state.nextProbeAt.After(now) {
			continue
		}
		out = append(out, state.snapshot)
	}
	return out
}

func (p *Pool) orderedProvidersLocked(mode Mode) []*providerState {
	if mode == "" {
		mode = p.mode
	}
	if len(p.providers) <= 1 {
		return append([]*providerState(nil), p.providers...)
	}

	if mode == ModeLeastFailures {
		out := append([]*providerState(nil), p.providers...)
		sort.SliceStable(out, func(i, j int) bool {
			return out[i].totalFailures < out[j].totalFailures
		})
		return out
	}

	if mode != ModeRoundRobin {
		return append([]*providerState(nil), p.providers...)
	}

	out := make([]*providerState, 0, len(p.providers))
	start := p.cursor % len(p.providers)
	for i := 0; i < len(p.providers); i++ {
		out = append(out, p.providers[(start+i)%len(p.providers)])
	}
	return out
}

func parseMode(mode string) Mode {
	switch strings.ToLower(strings.TrimSpace(mode)) {
	case "sequential":
		return ModeSequential
	case "least_failures", "leastfailures", "least-failures", "lf":
		return ModeLeastFailures
	case "round_robin", "roundrobin", "rr", "":
		return ModeRoundRobin
	default:
		return ModeRoundRobin
	}
}

func buildTargetURL(base *url.URL, path, rawQuery string) *url.URL {
	target := *base
	target.Path = upstream.JoinPaths(base.Path, path)
	target.RawQuery = rawQuery
	target.Fragment = ""
	return &target
}
