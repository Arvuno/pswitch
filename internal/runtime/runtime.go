package runtime

import (
	"fmt"
	"sync"
	"time"

	"pswitch/internal/config"
	"pswitch/internal/metrics"
	"pswitch/internal/pool"
)

type Manager struct {
	mu         sync.RWMutex
	settingsPath string
	cfg        config.Config
	pool       *pool.Pool
	metrics    *metrics.Store
}

type ProviderStatus struct {
	Name                string    `json:"name"`
	BaseURL             string    `json:"base_url"`
	Healthy             bool      `json:"healthy"`
	ConsecutiveFailures int       `json:"consecutive_failures"`
	NextProbeAt         time.Time `json:"next_probe_at"`
}

type UpdateResult struct {
	RequiresRestart bool     `json:"requires_restart"`
	Messages        []string `json:"messages"`
}

func New(settingsPath, metricsPath string, cfg config.Config) (*Manager, error) {
	providerPool, err := pool.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("build provider pool: %w", err)
	}
	metricsStore, err := metrics.Open(metricsPath)
	if err != nil {
		return nil, fmt.Errorf("open metrics store: %w", err)
	}

	return &Manager{
		settingsPath: settingsPath,
		cfg:        cloneConfig(cfg),
		pool:       providerPool,
		metrics:    metricsStore,
	}, nil
}

func (m *Manager) Config() config.Config {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return cloneConfig(m.cfg)
}

func (m *Manager) Pool() *pool.Pool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.pool
}

func (m *Manager) Snapshot() (config.Config, *pool.Pool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return cloneConfig(m.cfg), m.pool
}

func (m *Manager) Metrics() *metrics.Store {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.metrics
}

func (m *Manager) MetricsSnapshot(now time.Time) metrics.Snapshot {
	m.mu.RLock()
	store := m.metrics
	m.mu.RUnlock()
	return store.Snapshot(now)
}

func (m *Manager) UpdateConfig(next config.Config) (UpdateResult, error) {
	if err := next.Validate(); err != nil {
		return UpdateResult{}, err
	}

	nextPool, err := pool.New(next)
	if err != nil {
		return UpdateResult{}, fmt.Errorf("build provider pool: %w", err)
	}

	current := m.Config()
	result := UpdateResult{}
	if next.Listen != current.Listen {
		result.RequiresRestart = true
		result.Messages = append(result.Messages, "listen changed and requires process restart to take effect")
	}

	if err := config.WriteJSON(m.settingsPath, next); err != nil {
		return UpdateResult{}, err
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	m.cfg = cloneConfig(next)
	m.pool = nextPool
	return result, nil
}

func (m *Manager) ProviderStatuses() []ProviderStatus {
	cfg, providerPool := m.Snapshot()
	statuses := make([]ProviderStatus, 0, len(cfg.Providers))
	for _, provider := range cfg.Providers {
		status, ok := providerPool.Status(provider.Name)
		if !ok {
			continue
		}
		statuses = append(statuses, ProviderStatus{
			Name:                provider.Name,
			BaseURL:             provider.BaseURL,
			Healthy:             status.Healthy,
			ConsecutiveFailures: status.ConsecutiveFailures,
			NextProbeAt:         status.NextProbeAt,
		})
	}
	return statuses
}

func cloneConfig(cfg config.Config) config.Config {
	out := cfg
	out.Routes = append([]config.Route(nil), cfg.Routes...)
	out.Providers = append([]config.Provider(nil), cfg.Providers...)
	return out
}
