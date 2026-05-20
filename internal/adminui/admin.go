package adminui

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"

	"pswitch/internal/config"
	"pswitch/internal/metrics"
	pruntime "pswitch/internal/runtime"
)

const headerAdminToken = "X-PSwitch-Admin-Token"

//go:embed assets/*
var assets embed.FS

type configPayload struct {
	Listen              string            `json:"listen"`
	Mode                string            `json:"mode"`
	FailureThreshold    int               `json:"failure_threshold"`
	Cooldown            string            `json:"cooldown"`
	HealthCheckInterval string            `json:"health_check_interval"`
	HealthCheckTimeout  string            `json:"health_check_timeout"`
	Routes              []routePayload    `json:"routes"`
	Providers           []providerPayload `json:"providers"`
}

type routePayload struct {
	Prefix        string `json:"prefix"`
	Kind          string `json:"type"`
	Model         string `json:"model"`
	UpstreamModel string `json:"upstream_model"`
	Enabled       bool   `json:"enabled"`
}

type providerPayload struct {
	Name    string `json:"name"`
	BaseURL string `json:"base_url"`
	APIKey  string `json:"api_key"`
	Enabled bool   `json:"enabled"`
}

type stateResponse struct {
	Config     configPayload             `json:"config"`
	Providers  []pruntime.ProviderStatus `json:"providers"`
	ServerTime time.Time                 `json:"server_time"`
}

type saveResponse struct {
	RequiresRestart bool          `json:"requires_restart"`
	Messages        []string      `json:"messages"`
	State           stateResponse `json:"state"`
}

type statsResponse struct {
	Overview   statsOverviewPayload   `json:"overview"`
	Windows    statsWindowsPayload    `json:"windows"`
	Providers  []statsProviderPayload `json:"providers"`
	Models     []statsModelPayload    `json:"models"`
	Series     metrics.SeriesSnapshot `json:"series"`
	ServerTime time.Time              `json:"server_time"`
}

type statsOverviewPayload struct {
	TotalRequests     int64 `json:"total_requests"`
	TotalFailures     int64 `json:"total_failures"`
	TotalInputTokens  int64 `json:"total_input_tokens"`
	TotalOutputTokens int64 `json:"total_output_tokens"`
	TotalTokens       int64 `json:"total_tokens"`
	HealthyProviders  int   `json:"healthy_providers"`
	ProvidersCount    int   `json:"providers_count"`
}

type statsWindowsPayload struct {
	Last24h metrics.WindowSnapshot `json:"last_24h"`
	Last7d  metrics.WindowSnapshot `json:"last_7d"`
}

type statsProviderPayload struct {
	Name                string                 `json:"name"`
	BaseURL             string                 `json:"base_url"`
	Healthy             bool                   `json:"healthy"`
	ConsecutiveFailures int                    `json:"consecutive_failures"`
	NextProbeAt         time.Time              `json:"next_probe_at"`
	RequestCount        int64                  `json:"request_count"`
	FailureCount        int64                  `json:"failure_count"`
	InputTokens         int64                  `json:"input_tokens"`
	OutputTokens        int64                  `json:"output_tokens"`
	TotalTokens         int64                  `json:"total_tokens"`
	Last24h             metrics.WindowSnapshot `json:"last_24h"`
	Last7d              metrics.WindowSnapshot `json:"last_7d"`
}

type statsModelPayload struct {
	Name         string                 `json:"name"`
	RequestCount int64                  `json:"request_count"`
	FailureCount int64                  `json:"failure_count"`
	InputTokens  int64                  `json:"input_tokens"`
	OutputTokens int64                  `json:"output_tokens"`
	TotalTokens  int64                  `json:"total_tokens"`
	Last24h      metrics.WindowSnapshot `json:"last_24h"`
	Last7d       metrics.WindowSnapshot `json:"last_7d"`
}

func New(manager *pruntime.Manager, adminToken string) http.Handler {
	router := chi.NewRouter()
	router.Get("/api/meta", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{
			"token_required":          adminToken != "",
			"restart_required_fields": []string{"listen"},
			"admin_prefix":            "/dashboard",
		})
	})

	router.Group(func(r chi.Router) {
		if adminToken != "" {
			r.Use(requireToken(adminToken))
		}
		r.Get("/api/state", func(w http.ResponseWriter, r *http.Request) {
			writeJSON(w, http.StatusOK, buildState(manager))
		})
		r.Get("/api/stats", func(w http.ResponseWriter, r *http.Request) {
			writeJSON(w, http.StatusOK, buildStats(manager))
		})
		r.Put("/api/config", func(w http.ResponseWriter, r *http.Request) {
			var payload configPayload
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				writeError(w, http.StatusBadRequest, fmt.Errorf("decode config payload: %w", err))
				return
			}

			cfg, err := payload.toConfig()
			if err != nil {
				writeError(w, http.StatusBadRequest, err)
				return
			}

			result, err := manager.UpdateConfig(cfg)
			if err != nil {
				writeError(w, http.StatusBadRequest, err)
				return
			}

			writeJSON(w, http.StatusOK, saveResponse{
				RequiresRestart: result.RequiresRestart,
				Messages:        result.Messages,
				State:           buildState(manager),
			})
		})
	})

	staticFS, err := fs.Sub(assets, "assets")
	if err != nil {
		panic(err)
	}
	router.Get("/", serveAsset(staticFS, "index.html"))
	router.Get("/app.js", serveAsset(staticFS, "app.js"))
	router.Get("/style.css", serveAsset(staticFS, "style.css"))
	router.Get("/icon.svg", serveAsset(staticFS, "icon.svg"))
	router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		name := strings.TrimPrefix(path.Clean(r.URL.Path), "/")
		if name == "" || name == "." {
			name = "index.html"
		}
		serveAsset(staticFS, name)(w, r)
	})

	return router
}

func requireToken(expected string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := strings.TrimSpace(r.Header.Get(headerAdminToken))
			if token == "" {
				auth := strings.TrimSpace(r.Header.Get("Authorization"))
				if strings.HasPrefix(strings.ToLower(auth), "bearer ") {
					token = strings.TrimSpace(auth[7:])
				}
			}
			if token != expected {
				writeError(w, http.StatusUnauthorized, errors.New("admin token required"))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func serveAsset(staticFS fs.FS, name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := fs.ReadFile(staticFS, name)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		switch path.Ext(name) {
		case ".css":
			w.Header().Set("Content-Type", "text/css; charset=utf-8")
		case ".js":
			w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		case ".svg":
			w.Header().Set("Content-Type", "image/svg+xml")
		default:
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
		}
		_, _ = w.Write(data)
	}
}

func buildState(manager *pruntime.Manager) stateResponse {
	cfg := manager.Config()
	return stateResponse{
		Config:     configFromModel(cfg),
		Providers:  manager.ProviderStatuses(),
		ServerTime: time.Now().UTC(),
	}
}

func buildStats(manager *pruntime.Manager) statsResponse {
	now := time.Now().UTC()
	snapshot := manager.MetricsSnapshot(now)
	statuses := manager.ProviderStatuses()

	statusByName := make(map[string]pruntime.ProviderStatus, len(statuses))
	healthyProviders := 0
	for _, status := range statuses {
		statusByName[status.Name] = status
		if status.Healthy {
			healthyProviders++
		}
	}

	cfg := manager.Config()
	providers := make([]statsProviderPayload, 0, len(cfg.Providers))
	for _, provider := range cfg.Providers {
		status := statusByName[provider.Name]
		metric := snapshot.Provider(provider.Name)
		providers = append(providers, statsProviderPayload{
			Name:                provider.Name,
			BaseURL:             provider.BaseURL,
			Healthy:             status.Healthy,
			ConsecutiveFailures: status.ConsecutiveFailures,
			NextProbeAt:         status.NextProbeAt,
			RequestCount:        metric.RequestCount,
			FailureCount:        metric.FailureCount,
			InputTokens:         metric.InputTokens,
			OutputTokens:        metric.OutputTokens,
			TotalTokens:         metric.TotalTokens,
			Last24h:             metric.Last24h,
			Last7d:              metric.Last7d,
		})
	}

	models := make([]statsModelPayload, 0, len(snapshot.Models))
	for _, model := range snapshot.Models {
		models = append(models, statsModelPayload{
			Name:         model.Name,
			RequestCount: model.RequestCount,
			FailureCount: model.FailureCount,
			InputTokens:  model.InputTokens,
			OutputTokens: model.OutputTokens,
			TotalTokens:  model.TotalTokens,
			Last24h:      model.Last24h,
			Last7d:       model.Last7d,
		})
	}
	sort.Slice(models, func(i, j int) bool {
		if models[i].TotalTokens == models[j].TotalTokens {
			return models[i].Name < models[j].Name
		}
		return models[i].TotalTokens > models[j].TotalTokens
	})

	return statsResponse{
		Overview: statsOverviewPayload{
			TotalRequests:     snapshot.Overview.TotalRequests,
			TotalFailures:     snapshot.Overview.TotalFailures,
			TotalInputTokens:  snapshot.Overview.TotalInputTokens,
			TotalOutputTokens: snapshot.Overview.TotalOutputTokens,
			TotalTokens:       snapshot.Overview.TotalTokens,
			HealthyProviders:  healthyProviders,
			ProvidersCount:    len(cfg.Providers),
		},
		Windows: statsWindowsPayload{
			Last24h: snapshot.Windows.Last24h,
			Last7d:  snapshot.Windows.Last7d,
		},
		Providers:  providers,
		Models:     models,
		Series:     snapshot.Series,
		ServerTime: now,
	}
}

func configFromModel(cfg config.Config) configPayload {
	payload := configPayload{
		Listen:              cfg.Listen,
		Mode:                cfg.Mode,
		FailureThreshold:    cfg.FailureThreshold,
		Cooldown:            cfg.Cooldown.String(),
		HealthCheckInterval: cfg.HealthCheckInterval.String(),
		HealthCheckTimeout:  cfg.HealthCheckTimeout.String(),
		Routes:              make([]routePayload, 0, len(cfg.Routes)),
		Providers:           make([]providerPayload, 0, len(cfg.Providers)),
	}
	for _, route := range cfg.Routes {
		payload.Routes = append(payload.Routes, routePayload{
			Prefix:        route.Prefix,
			Kind:          route.Kind,
			Model:         route.Model,
			UpstreamModel: route.UpstreamModel,
			Enabled:       route.Enabled,
		})
	}
	for _, provider := range cfg.Providers {
		payload.Providers = append(payload.Providers, providerPayload{
			Name:    provider.Name,
			BaseURL: provider.BaseURL,
			APIKey:  provider.APIKey,
			Enabled: provider.Enabled,
		})
	}
	return payload
}

func (p configPayload) toConfig() (config.Config, error) {
	cooldown, err := time.ParseDuration(strings.TrimSpace(p.Cooldown))
	if err != nil {
		return config.Config{}, fmt.Errorf("parse cooldown: %w", err)
	}
	healthCheckInterval, err := time.ParseDuration(strings.TrimSpace(p.HealthCheckInterval))
	if err != nil {
		return config.Config{}, fmt.Errorf("parse health_check_interval: %w", err)
	}
	healthCheckTimeout, err := time.ParseDuration(strings.TrimSpace(p.HealthCheckTimeout))
	if err != nil {
		return config.Config{}, fmt.Errorf("parse health_check_timeout: %w", err)
	}

	cfg := config.Config{
		Listen:              strings.TrimSpace(p.Listen),
		Mode:                strings.TrimSpace(p.Mode),
		FailureThreshold:    p.FailureThreshold,
		Cooldown:            cooldown,
		HealthCheckInterval: healthCheckInterval,
		HealthCheckTimeout:  healthCheckTimeout,
		Routes:              make([]config.Route, 0, len(p.Routes)),
		Providers:           make([]config.Provider, 0, len(p.Providers)),
	}
	for _, route := range p.Routes {
		cfg.Routes = append(cfg.Routes, config.Route{
			Prefix:        route.Prefix,
			Kind:          route.Kind,
			Model:         route.Model,
			UpstreamModel: route.UpstreamModel,
			Enabled:       route.Enabled,
		})
	}
	for _, provider := range p.Providers {
		cfg.Providers = append(cfg.Providers, config.Provider{
			Name:    provider.Name,
			BaseURL: provider.BaseURL,
			APIKey:  provider.APIKey,
			Enabled: provider.Enabled,
		})
	}
	return cfg, nil
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}
