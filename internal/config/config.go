package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pelletier/go-toml/v2"
)

const (
	defaultListen              = "0.0.0.0:8080"
	defaultMode                = "round_robin"
	defaultFailureThreshold    = 1
	defaultCooldown            = 20 * time.Second
	defaultHealthCheckInterval = 15 * time.Second
	defaultHealthCheckTimeout  = 3 * time.Second
)

type Config struct {
	Listen              string
	Mode                string
	FailureThreshold    int
	Cooldown            time.Duration
	HealthCheckInterval time.Duration
	HealthCheckTimeout  time.Duration
	Routes              []Route
	Providers           []Provider
}

type Route struct {
	Prefix        string
	Kind          string
	Model         string
	UpstreamModel string
	Enabled       bool
}

type Provider struct {
	Name    string
	BaseURL string
	APIKey  string
	Enabled bool
}

type rawConfig struct {
	Listen              string        `toml:"listen" json:"listen"`
	Mode                string        `toml:"mode" json:"mode"`
	FailureThreshold    int           `toml:"failure_threshold" json:"failure_threshold"`
	Cooldown            string        `toml:"cooldown" json:"cooldown"`
	HealthCheckInterval string        `toml:"health_check_interval" json:"health_check_interval"`
	HealthCheckTimeout  string        `toml:"health_check_timeout" json:"health_check_timeout"`
	Routes              []rawRoute    `toml:"routes" json:"routes"`
	Providers           []rawProvider `toml:"providers" json:"providers"`
}

type rawRoute struct {
	Prefix        string `toml:"prefix" json:"prefix"`
	Kind          string `toml:"type" json:"type"`
	Model         string `toml:"model" json:"model"`
	UpstreamModel string `toml:"upstream_model" json:"upstream_model"`
	Enabled       *bool  `toml:"enabled" json:"enabled"`
}

type rawProvider struct {
	Name    string `toml:"name" json:"name"`
	BaseURL string `toml:"base_url" json:"base_url"`
	APIKey  string `toml:"api_key" json:"api_key"`
	Enabled *bool  `toml:"enabled" json:"enabled"`
}

func Load(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("read config: %w", err)
	}

	var raw rawConfig
	if err := toml.Unmarshal(data, &raw); err != nil {
		return Config{}, fmt.Errorf("decode toml: %w", err)
	}
	return configFromRaw(raw)
}

func Default() Config {
	return Config{
		Listen:              defaultListen,
		Mode:                defaultMode,
		FailureThreshold:    defaultFailureThreshold,
		Cooldown:            defaultCooldown,
		HealthCheckInterval: defaultHealthCheckInterval,
		HealthCheckTimeout:  defaultHealthCheckTimeout,
		Routes: []Route{
			{
				Prefix:  "/codex",
				Kind:    "openai",
				Enabled: true,
			},
		},
		Providers: nil,
	}
}

func LoadJSON(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("read config: %w", err)
	}

	var raw rawConfig
	if err := json.Unmarshal(data, &raw); err != nil {
		return Config{}, fmt.Errorf("decode json: %w", err)
	}
	return configFromRaw(raw)
}

func Write(path string, cfg Config) error {
	if err := cfg.Validate(); err != nil {
		return err
	}

	data, err := toml.Marshal(rawConfigFromConfig(cfg))
	if err != nil {
		return fmt.Errorf("encode toml: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}

	tempFile, err := os.CreateTemp(filepath.Dir(path), ".pswitch-config-*.toml")
	if err != nil {
		return fmt.Errorf("create temp config: %w", err)
	}
	tempPath := tempFile.Name()
	defer func() {
		_ = os.Remove(tempPath)
	}()

	if _, err := tempFile.Write(data); err != nil {
		_ = tempFile.Close()
		return fmt.Errorf("write temp config: %w", err)
	}
	if err := tempFile.Chmod(0o600); err != nil {
		_ = tempFile.Close()
		return fmt.Errorf("chmod temp config: %w", err)
	}
	if err := tempFile.Close(); err != nil {
		return fmt.Errorf("close temp config: %w", err)
	}
	if err := os.Rename(tempPath, path); err != nil {
		return fmt.Errorf("replace config: %w", err)
	}
	return nil
}

func WriteJSON(path string, cfg Config) error {
	if err := cfg.Validate(); err != nil {
		return err
	}

	data, err := json.MarshalIndent(rawConfigFromConfig(cfg), "", "  ")
	if err != nil {
		return fmt.Errorf("encode json: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}

	tempFile, err := os.CreateTemp(filepath.Dir(path), ".pswitch-settings-*.json")
	if err != nil {
		return fmt.Errorf("create temp settings: %w", err)
	}
	tempPath := tempFile.Name()
	defer func() {
		_ = os.Remove(tempPath)
	}()

	if _, err := tempFile.Write(data); err != nil {
		_ = tempFile.Close()
		return fmt.Errorf("write temp settings: %w", err)
	}
	if err := tempFile.Chmod(0o600); err != nil {
		_ = tempFile.Close()
		return fmt.Errorf("chmod temp settings: %w", err)
	}
	if err := tempFile.Close(); err != nil {
		return fmt.Errorf("close temp settings: %w", err)
	}
	if err := os.Rename(tempPath, path); err != nil {
		return fmt.Errorf("replace settings: %w", err)
	}
	return nil
}

func (c Config) Validate() error {
	var problems []error

	if strings.TrimSpace(c.Listen) == "" {
		problems = append(problems, errors.New("listen is required"))
	}

	switch normalizeMode(c.Mode) {
	case "sequential", "round_robin", "least_failures":
	default:
		problems = append(problems, fmt.Errorf("unsupported mode %q", c.Mode))
	}

	if c.FailureThreshold <= 0 {
		problems = append(problems, errors.New("failure_threshold must be greater than zero"))
	}
	if c.Cooldown < 0 {
		problems = append(problems, errors.New("cooldown must be zero or positive"))
	}
	if c.HealthCheckInterval <= 0 {
		problems = append(problems, errors.New("health_check_interval must be greater than zero"))
	}
	if c.HealthCheckTimeout <= 0 {
		problems = append(problems, errors.New("health_check_timeout must be greater than zero"))
	}

	seen := make(map[string]struct{})
	seenPrefixes := make(map[string]struct{})
	enabledRoutes := 0
	rootRouteEnabled := false
	for _, route := range c.Routes {
		if !route.Enabled {
			continue
		}
		enabledRoutes++
		if route.Prefix == "" {
			problems = append(problems, errors.New("route prefix is required"))
		}
		switch route.Kind {
		case "openai", "anthropic":
		default:
			problems = append(problems, fmt.Errorf("unsupported route type %q", route.Kind))
		}
		if route.Kind == "anthropic" && route.UpstreamModel == "" {
			problems = append(problems, fmt.Errorf("anthropic route %q upstream_model is required", route.Prefix))
		}
		key := strings.ToLower(route.Prefix)
		if _, ok := seenPrefixes[key]; ok {
			problems = append(problems, fmt.Errorf("duplicate route prefix %q", route.Prefix))
		}
		seenPrefixes[key] = struct{}{}
		if route.Prefix == "/" {
			rootRouteEnabled = true
		}
	}
	if enabledRoutes == 0 {
		problems = append(problems, errors.New("at least one enabled route is required"))
	}
	if rootRouteEnabled && enabledRoutes > 1 {
		problems = append(problems, errors.New("route prefix \"/\" cannot be combined with other routes"))
	}

	enabledCount := 0
	for _, provider := range c.Providers {
		if !provider.Enabled {
			continue
		}
		enabledCount++

		if provider.Name == "" {
			problems = append(problems, errors.New("provider name is required"))
		}
		if provider.BaseURL == "" {
			problems = append(problems, fmt.Errorf("provider %q base_url is required", provider.Name))
		}
		if provider.APIKey == "" {
			problems = append(problems, fmt.Errorf("provider %q api_key is required", provider.Name))
		}

		key := strings.ToLower(provider.Name)
		if _, ok := seen[key]; ok {
			problems = append(problems, fmt.Errorf("duplicate provider name %q", provider.Name))
		}
		seen[key] = struct{}{}
	}

	return errors.Join(problems...)
}

func normalizeMode(mode string) string {
	switch strings.ToLower(strings.TrimSpace(mode)) {
	case "":
		return defaultMode
	case "sequential":
		return "sequential"
	case "round_robin", "roundrobin", "rr":
		return "round_robin"
	case "least_failures", "leastfailures", "least-failures", "lf":
		return "least_failures"
	default:
		return strings.ToLower(strings.TrimSpace(mode))
	}
}

func defaultString(value, fallback string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback
	}
	return value
}

func normalizeRouteKind(kind string) string {
	switch strings.ToLower(strings.TrimSpace(kind)) {
	case "", "openai":
		return "openai"
	case "anthropic", "claude":
		return "anthropic"
	default:
		return strings.ToLower(strings.TrimSpace(kind))
	}
}

func normalizeRoutePrefix(prefix string) string {
	prefix = strings.TrimSpace(prefix)
	if prefix == "" {
		return "/"
	}
	if !strings.HasPrefix(prefix, "/") {
		prefix = "/" + prefix
	}
	if len(prefix) > 1 {
		prefix = strings.TrimSuffix(prefix, "/")
	}
	return prefix
}

func rawConfigFromConfig(cfg Config) rawConfig {
	raw := rawConfig{
		Listen:              strings.TrimSpace(cfg.Listen),
		Mode:                normalizeMode(cfg.Mode),
		FailureThreshold:    cfg.FailureThreshold,
		Cooldown:            cfg.Cooldown.String(),
		HealthCheckInterval: cfg.HealthCheckInterval.String(),
		HealthCheckTimeout:  cfg.HealthCheckTimeout.String(),
		Routes:              make([]rawRoute, 0, len(cfg.Routes)),
		Providers:           make([]rawProvider, 0, len(cfg.Providers)),
	}

	for _, route := range cfg.Routes {
		enabled := route.Enabled
		raw.Routes = append(raw.Routes, rawRoute{
			Prefix:        normalizeRoutePrefix(route.Prefix),
			Kind:          normalizeRouteKind(route.Kind),
			Model:         strings.TrimSpace(route.Model),
			UpstreamModel: strings.TrimSpace(route.UpstreamModel),
			Enabled:       &enabled,
		})
	}

	for _, provider := range cfg.Providers {
		enabled := provider.Enabled
		raw.Providers = append(raw.Providers, rawProvider{
			Name:    strings.TrimSpace(provider.Name),
			BaseURL: strings.TrimSpace(provider.BaseURL),
			APIKey:  strings.TrimSpace(provider.APIKey),
			Enabled: &enabled,
		})
	}

	return raw
}

func configFromRaw(raw rawConfig) (Config, error) {
	var err error
	cfg := Config{
		Listen:              defaultString(raw.Listen, defaultListen),
		Mode:                normalizeMode(defaultString(raw.Mode, defaultMode)),
		FailureThreshold:    raw.FailureThreshold,
		Cooldown:            defaultCooldown,
		HealthCheckInterval: defaultHealthCheckInterval,
		HealthCheckTimeout:  defaultHealthCheckTimeout,
	}

	if cfg.FailureThreshold <= 0 {
		cfg.FailureThreshold = defaultFailureThreshold
	}

	if raw.Cooldown != "" {
		cfg.Cooldown, err = time.ParseDuration(raw.Cooldown)
		if err != nil {
			return Config{}, fmt.Errorf("parse cooldown: %w", err)
		}
	}
	if raw.HealthCheckInterval != "" {
		cfg.HealthCheckInterval, err = time.ParseDuration(raw.HealthCheckInterval)
		if err != nil {
			return Config{}, fmt.Errorf("parse health_check_interval: %w", err)
		}
	}
	if raw.HealthCheckTimeout != "" {
		cfg.HealthCheckTimeout, err = time.ParseDuration(raw.HealthCheckTimeout)
		if err != nil {
			return Config{}, fmt.Errorf("parse health_check_timeout: %w", err)
		}
	}

	cfg.Routes = make([]Route, 0, len(raw.Routes))
	for _, item := range raw.Routes {
		enabled := true
		if item.Enabled != nil {
			enabled = *item.Enabled
		}
		if !enabled {
			continue
		}

		cfg.Routes = append(cfg.Routes, Route{
			Prefix:        normalizeRoutePrefix(item.Prefix),
			Kind:          normalizeRouteKind(item.Kind),
			Model:         strings.TrimSpace(item.Model),
			UpstreamModel: strings.TrimSpace(item.UpstreamModel),
			Enabled:       enabled,
		})
	}
	if len(cfg.Routes) == 0 {
		cfg.Routes = []Route{{
			Prefix:  "/",
			Kind:    "openai",
			Enabled: true,
		}}
	}

	cfg.Providers = make([]Provider, 0, len(raw.Providers))
	for _, item := range raw.Providers {
		enabled := true
		if item.Enabled != nil {
			enabled = *item.Enabled
		}
		if !enabled {
			continue
		}

		cfg.Providers = append(cfg.Providers, Provider{
			Name:    strings.TrimSpace(item.Name),
			BaseURL: strings.TrimSpace(item.BaseURL),
			APIKey:  strings.TrimSpace(item.APIKey),
			Enabled: enabled,
		})
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
