package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadParsesProvidersAndDurations(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")

	err := os.WriteFile(path, []byte(`
listen = "127.0.0.1:0"
mode = "round_robin"
failure_threshold = 2
cooldown = "30s"
health_check_interval = "5s"
health_check_timeout = "2s"

[[providers]]
name = "first"
base_url = "http://127.0.0.1:10001"
api_key = "k1"

[[providers]]
name = "second"
base_url = "http://127.0.0.1:10002"
api_key = "k2"
`), 0o600)
	if err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := cfg.Listen, "127.0.0.1:0"; got != want {
		t.Fatalf("listen = %q, want %q", got, want)
	}
	if got, want := cfg.Mode, "round_robin"; got != want {
		t.Fatalf("mode = %q, want %q", got, want)
	}
	if got, want := cfg.FailureThreshold, 2; got != want {
		t.Fatalf("failure_threshold = %d, want %d", got, want)
	}
	if got, want := cfg.Cooldown, 30*time.Second; got != want {
		t.Fatalf("cooldown = %s, want %s", got, want)
	}
	if got, want := cfg.HealthCheckInterval, 5*time.Second; got != want {
		t.Fatalf("health_check_interval = %s, want %s", got, want)
	}
	if got, want := cfg.HealthCheckTimeout, 2*time.Second; got != want {
		t.Fatalf("health_check_timeout = %s, want %s", got, want)
	}
	if got, want := len(cfg.Providers), 2; got != want {
		t.Fatalf("providers = %d, want %d", got, want)
	}
	if got, want := cfg.Providers[0].Name, "first"; got != want {
		t.Fatalf("provider[0].name = %q, want %q", got, want)
	}
}

func TestLoadAppliesDefaults(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")

	err := os.WriteFile(path, []byte(`
[[providers]]
name = "only"
base_url = "http://127.0.0.1:10001"
api_key = "k1"
`), 0o600)
	if err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := cfg.Listen, "0.0.0.0:8080"; got != want {
		t.Fatalf("listen = %q, want %q", got, want)
	}
	if got, want := cfg.Mode, "round_robin"; got != want {
		t.Fatalf("mode = %q, want %q", got, want)
	}
	if got, want := cfg.FailureThreshold, 1; got != want {
		t.Fatalf("failure_threshold = %d, want %d", got, want)
	}
	if got, want := len(cfg.Routes), 1; got != want {
		t.Fatalf("routes = %d, want %d", got, want)
	}
	if got, want := cfg.Routes[0].Prefix, "/"; got != want {
		t.Fatalf("route[0].prefix = %q, want %q", got, want)
	}
}

func TestLoadAcceptsLeastFailuresMode(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")

	err := os.WriteFile(path, []byte(`
mode = "least_failures"

[[providers]]
name = "only"
base_url = "http://127.0.0.1:10001"
api_key = "k1"
`), 0o600)
	if err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := cfg.Mode, "least_failures"; got != want {
		t.Fatalf("mode = %q, want %q", got, want)
	}
}

func TestLoadAllowsEmptyFileWithDefaultRoutesAndNoProviders(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")

	if err := os.WriteFile(path, []byte(""), 0o600); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := cfg.Listen, "0.0.0.0:8080"; got != want {
		t.Fatalf("listen = %q, want %q", got, want)
	}
	if got, want := len(cfg.Routes), 1; got != want {
		t.Fatalf("routes = %d, want %d", got, want)
	}
	if got, want := cfg.Routes[0].Prefix, "/"; got != want {
		t.Fatalf("route[0].prefix = %q, want %q", got, want)
	}
	if got, want := len(cfg.Providers), 0; got != want {
		t.Fatalf("providers = %d, want %d", got, want)
	}
}

func TestLoadParsesRoutes(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")

	err := os.WriteFile(path, []byte(`
[[routes]]
prefix = "/codex"
type = "openai"

[[routes]]
prefix = "/claude"
type = "anthropic"
model = "claude-sonnet-4-20250514"
upstream_model = "gpt-5.4"

[[providers]]
name = "only"
base_url = "http://127.0.0.1:10001"
api_key = "k1"
`), 0o600)
	if err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := len(cfg.Routes), 2; got != want {
		t.Fatalf("routes = %d, want %d", got, want)
	}
	if got, want := cfg.Routes[0].Prefix, "/codex"; got != want {
		t.Fatalf("route[0].prefix = %q, want %q", got, want)
	}
	if got, want := cfg.Routes[1].Kind, "anthropic"; got != want {
		t.Fatalf("route[1].kind = %q, want %q", got, want)
	}
}

func TestWriteRoundTripsConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")

	cfg := Config{
		Listen:              "127.0.0.1:8080",
		Mode:                "sequential",
		FailureThreshold:    2,
		Cooldown:            30 * time.Second,
		HealthCheckInterval: 5 * time.Second,
		HealthCheckTimeout:  2 * time.Second,
		Routes: []Route{
			{Prefix: "/codex", Kind: "openai", Enabled: true},
			{Prefix: "/claude", Kind: "anthropic", Model: "claude-sonnet-4-20250514", UpstreamModel: "gpt-5.4", Enabled: true},
		},
		Providers: []Provider{
			{Name: "one", BaseURL: "http://127.0.0.1:10001", APIKey: "k1", Enabled: true},
			{Name: "two", BaseURL: "http://127.0.0.1:10002", APIKey: "k2", Enabled: true},
		},
	}

	if err := Write(path, cfg); err != nil {
		t.Fatal(err)
	}

	got, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}

	if got.Listen != cfg.Listen {
		t.Fatalf("listen = %q, want %q", got.Listen, cfg.Listen)
	}
	if got.Mode != cfg.Mode {
		t.Fatalf("mode = %q, want %q", got.Mode, cfg.Mode)
	}
	if got.FailureThreshold != cfg.FailureThreshold {
		t.Fatalf("failure_threshold = %d, want %d", got.FailureThreshold, cfg.FailureThreshold)
	}
	if got.Cooldown != cfg.Cooldown {
		t.Fatalf("cooldown = %s, want %s", got.Cooldown, cfg.Cooldown)
	}
	if got.HealthCheckInterval != cfg.HealthCheckInterval {
		t.Fatalf("health_check_interval = %s, want %s", got.HealthCheckInterval, cfg.HealthCheckInterval)
	}
	if got.HealthCheckTimeout != cfg.HealthCheckTimeout {
		t.Fatalf("health_check_timeout = %s, want %s", got.HealthCheckTimeout, cfg.HealthCheckTimeout)
	}
	if len(got.Routes) != len(cfg.Routes) {
		t.Fatalf("routes = %d, want %d", len(got.Routes), len(cfg.Routes))
	}
	if len(got.Providers) != len(cfg.Providers) {
		t.Fatalf("providers = %d, want %d", len(got.Providers), len(cfg.Providers))
	}
	if got.Routes[1].UpstreamModel != cfg.Routes[1].UpstreamModel {
		t.Fatalf("route upstream_model = %q, want %q", got.Routes[1].UpstreamModel, cfg.Routes[1].UpstreamModel)
	}
	if got.Providers[1].APIKey != cfg.Providers[1].APIKey {
		t.Fatalf("provider api_key = %q, want %q", got.Providers[1].APIKey, cfg.Providers[1].APIKey)
	}
}
