package runtime

import (
	"path/filepath"
	"testing"
	"time"

	"pswitch/internal/config"
	"pswitch/internal/metrics"
)

func TestManagerUpdateConfigSwapsPoolAndPersistsFile(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.toml")
	settingsPath := filepath.Join(dir, "settings.json")
	metricsPath := filepath.Join(dir, "metrics.json")

	initial := config.Config{
		Listen:              "127.0.0.1:8080",
		Mode:                "round_robin",
		FailureThreshold:    1,
		Cooldown:            20 * time.Second,
		HealthCheckInterval: 15 * time.Second,
		HealthCheckTimeout:  3 * time.Second,
		Routes: []config.Route{
			{Prefix: "/codex", Kind: "openai", Enabled: true},
		},
		Providers: []config.Provider{
			{Name: "one", BaseURL: "http://127.0.0.1:10001", APIKey: "k1", Enabled: true},
		},
	}
	if err := config.Write(configPath, initial); err != nil {
		t.Fatal(err)
	}

	manager, err := New(settingsPath, metricsPath, initial)
	if err != nil {
		t.Fatal(err)
	}

	next := initial
	next.Mode = "sequential"
	next.Providers = []config.Provider{
		{Name: "two", BaseURL: "http://127.0.0.1:10002", APIKey: "k2", Enabled: true},
	}

	result, err := manager.UpdateConfig(next)
	if err != nil {
		t.Fatal(err)
	}
	if result.RequiresRestart {
		t.Fatal("update should not require restart")
	}

	current := manager.Config()
	if current.Mode != "sequential" {
		t.Fatalf("mode = %q, want %q", current.Mode, "sequential")
	}
	candidates := manager.Pool().Candidates("", time.Now())
	if len(candidates) != 1 || candidates[0].Name != "two" {
		t.Fatalf("candidates = %#v, want provider two", candidates)
	}

	reloaded, err := config.Load(configPath)
	if err != nil {
		t.Fatal(err)
	}
	if reloaded.Mode != initial.Mode {
		t.Fatalf("user config mode = %q, want %q", reloaded.Mode, initial.Mode)
	}
	if len(reloaded.Providers) != 1 || reloaded.Providers[0].Name != "one" {
		t.Fatalf("user config providers = %#v, want original provider one", reloaded.Providers)
	}

	saved, err := config.LoadJSON(settingsPath)
	if err != nil {
		t.Fatal(err)
	}
	if saved.Mode != "sequential" {
		t.Fatalf("saved mode = %q, want %q", saved.Mode, "sequential")
	}
	if len(saved.Providers) != 1 || saved.Providers[0].Name != "two" {
		t.Fatalf("saved providers = %#v, want provider two", saved.Providers)
	}
}

func TestManagerProviderStatusesExposeHealth(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")

	cfg := config.Config{
		Listen:              "127.0.0.1:8080",
		Mode:                "round_robin",
		FailureThreshold:    1,
		Cooldown:            20 * time.Second,
		HealthCheckInterval: 15 * time.Second,
		HealthCheckTimeout:  3 * time.Second,
		Routes: []config.Route{
			{Prefix: "/codex", Kind: "openai", Enabled: true},
		},
		Providers: []config.Provider{
			{Name: "one", BaseURL: "http://127.0.0.1:10001", APIKey: "k1", Enabled: true},
		},
	}
	if err := config.Write(path, cfg); err != nil {
		t.Fatal(err)
	}

	manager, err := New(filepath.Join(dir, "settings.json"), filepath.Join(dir, "metrics.json"), cfg)
	if err != nil {
		t.Fatal(err)
	}

	now := time.Unix(100, 0)
	manager.Pool().MarkFailure("one", now)

	statuses := manager.ProviderStatuses()
	if len(statuses) != 1 {
		t.Fatalf("statuses = %d, want 1", len(statuses))
	}
	if statuses[0].Name != "one" {
		t.Fatalf("status name = %q, want %q", statuses[0].Name, "one")
	}
	if statuses[0].Healthy {
		t.Fatal("provider should be unhealthy after failure")
	}
	if statuses[0].ConsecutiveFailures != 1 {
		t.Fatalf("failures = %d, want 1", statuses[0].ConsecutiveFailures)
	}
	if statuses[0].NextProbeAt.IsZero() {
		t.Fatal("next probe time should be set")
	}
}

func TestManagerUpdateConfigReportsRestartForListenChange(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")

	cfg := config.Config{
		Listen:              "127.0.0.1:8080",
		Mode:                "round_robin",
		FailureThreshold:    1,
		Cooldown:            20 * time.Second,
		HealthCheckInterval: 15 * time.Second,
		HealthCheckTimeout:  3 * time.Second,
		Routes: []config.Route{
			{Prefix: "/codex", Kind: "openai", Enabled: true},
		},
		Providers: []config.Provider{
			{Name: "one", BaseURL: "http://127.0.0.1:10001", APIKey: "k1", Enabled: true},
		},
	}
	if err := config.Write(path, cfg); err != nil {
		t.Fatal(err)
	}

	manager, err := New(filepath.Join(dir, "settings.json"), filepath.Join(dir, "metrics.json"), cfg)
	if err != nil {
		t.Fatal(err)
	}

	next := cfg
	next.Listen = "0.0.0.0:8080"

	result, err := manager.UpdateConfig(next)
	if err != nil {
		t.Fatal(err)
	}
	if !result.RequiresRestart {
		t.Fatal("listen change should require restart")
	}
	if len(result.Messages) != 1 {
		t.Fatalf("messages = %v, want 1 warning", result.Messages)
	}
}

func TestManagerExposesMetricsSnapshot(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")

	cfg := config.Config{
		Listen:              "127.0.0.1:8080",
		Mode:                "round_robin",
		FailureThreshold:    1,
		Cooldown:            20 * time.Second,
		HealthCheckInterval: 15 * time.Second,
		HealthCheckTimeout:  3 * time.Second,
		Routes: []config.Route{
			{Prefix: "/codex", Kind: "openai", Enabled: true},
		},
		Providers: []config.Provider{
			{Name: "one", BaseURL: "http://127.0.0.1:10001", APIKey: "k1", Enabled: true},
		},
	}
	if err := config.Write(path, cfg); err != nil {
		t.Fatal(err)
	}

	manager, err := New(filepath.Join(dir, "settings.json"), filepath.Join(dir, "metrics.json"), cfg)
	if err != nil {
		t.Fatal(err)
	}

	now := time.Date(2026, 5, 20, 12, 0, 0, 0, time.UTC)
	manager.Metrics().RecordSuccess("one", metrics.Usage{
		RequestCount: 1,
		InputTokens:  10,
		OutputTokens: 20,
		TotalTokens:  30,
	}, "gpt-5.4", now)
	manager.Metrics().RecordFailure("one", "gpt-5.4", now)

	snapshot := manager.MetricsSnapshot(now)
	if got, want := snapshot.Overview.TotalRequests, int64(1); got != want {
		t.Fatalf("total requests = %d, want %d", got, want)
	}
	if got, want := snapshot.Overview.TotalFailures, int64(1); got != want {
		t.Fatalf("total failures = %d, want %d", got, want)
	}
	if got, want := snapshot.Provider("one").TotalTokens, int64(30); got != want {
		t.Fatalf("provider total tokens = %d, want %d", got, want)
	}
	if got, want := snapshot.Model("gpt-5.4").TotalTokens, int64(30); got != want {
		t.Fatalf("model total tokens = %d, want %d", got, want)
	}
}
