package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"pswitch/internal/config"
)

func TestParseLogColorOverride(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want *bool
	}{
		{name: "default", args: nil, want: nil},
		{name: "explicit on", args: []string{"--log-color"}, want: boolPtr(true)},
		{name: "explicit off", args: []string{"--log-color=false"}, want: boolPtr(false)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseServeArgs(tt.args)
			if err != nil {
				t.Fatal(err)
			}
			if !sameBoolPtr(got.LogColor, tt.want) {
				t.Fatalf("log color override = %v, want %v", got.LogColor, tt.want)
			}
		})
	}
}

func TestParseServeArgsSupportsGeneralOverrides(t *testing.T) {
	got, err := parseServeArgs([]string{
		"--listen", "0.0.0.0:8080",
		"--mode", "least_failures",
		"--failure-threshold", "3",
		"--cooldown", "45s",
		"--health-check-interval", "30s",
		"--health-check-timeout", "5s",
	})
	if err != nil {
		t.Fatal(err)
	}

	if got.Listen != "0.0.0.0:8080" {
		t.Fatalf("listen = %q, want %q", got.Listen, "0.0.0.0:8080")
	}
	if got.Mode != "least_failures" {
		t.Fatalf("mode = %q, want %q", got.Mode, "least_failures")
	}
	if got.FailureThreshold == nil || *got.FailureThreshold != 3 {
		t.Fatalf("failure threshold = %v, want %d", got.FailureThreshold, 3)
	}
	if got.Cooldown == nil || *got.Cooldown != 45*time.Second {
		t.Fatalf("cooldown = %v, want %s", got.Cooldown, 45*time.Second)
	}
	if got.HealthCheckInterval == nil || *got.HealthCheckInterval != 30*time.Second {
		t.Fatalf("health check interval = %v, want %s", got.HealthCheckInterval, 30*time.Second)
	}
	if got.HealthCheckTimeout == nil || *got.HealthCheckTimeout != 5*time.Second {
		t.Fatalf("health check timeout = %v, want %s", got.HealthCheckTimeout, 5*time.Second)
	}
}

func TestLoadStartupConfigFallsBackToDefaultWhenMissing(t *testing.T) {
	path := t.TempDir() + "/missing.toml"

	cfg, err := loadStartupConfig(path, filepath.Join(t.TempDir(), "settings.json"))
	if err != nil {
		t.Fatal(err)
	}

	if got, want := cfg.Listen, "0.0.0.0:8080"; got != want {
		t.Fatalf("listen = %q, want %q", got, want)
	}
	if got, want := len(cfg.Providers), 0; got != want {
		t.Fatalf("providers = %d, want %d", got, want)
	}
	if got, want := len(cfg.Routes), 1; got != want {
		t.Fatalf("routes = %d, want %d", got, want)
	}
	if got, want := cfg.Routes[0].Prefix, "/codex"; got != want {
		t.Fatalf("route[0].prefix = %q, want %q", got, want)
	}
}

func TestLoadStartupConfigPrefersSettingsJSONWhenPresent(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.toml")
	settingsPath := filepath.Join(dir, "settings.json")

	base := config.Default()
	base.Listen = "127.0.0.1:8080"
	base.Providers = []config.Provider{
		{Name: "from-config", BaseURL: "http://127.0.0.1:10001", APIKey: "k1", Enabled: true},
	}
	if err := config.Write(configPath, base); err != nil {
		t.Fatal(err)
	}

	override := base
	override.Mode = "sequential"
	override.Providers = []config.Provider{
		{Name: "from-settings", BaseURL: "http://127.0.0.1:10002", APIKey: "k2", Enabled: true},
	}
	if err := config.WriteJSON(settingsPath, override); err != nil {
		t.Fatal(err)
	}

	got, err := loadStartupConfig(configPath, settingsPath)
	if err != nil {
		t.Fatal(err)
	}
	if got.Mode != "sequential" {
		t.Fatalf("mode = %q, want %q", got.Mode, "sequential")
	}
	if len(got.Providers) != 1 || got.Providers[0].Name != "from-settings" {
		t.Fatalf("providers = %#v, want settings override", got.Providers)
	}
}

func TestApplyServeOverridesUpdatesGeneralSettings(t *testing.T) {
	cfg := config.Default()

	applyServeOverrides(&cfg, serveArgs{
		Listen:              "0.0.0.0:8080",
		Mode:                "least_failures",
		FailureThreshold:    intPtr(4),
		Cooldown:            durationPtr(50 * time.Second),
		HealthCheckInterval: durationPtr(40 * time.Second),
		HealthCheckTimeout:  durationPtr(6 * time.Second),
	})

	if cfg.Listen != "0.0.0.0:8080" {
		t.Fatalf("listen = %q, want %q", cfg.Listen, "0.0.0.0:8080")
	}
	if cfg.Mode != "least_failures" {
		t.Fatalf("mode = %q, want %q", cfg.Mode, "least_failures")
	}
	if cfg.FailureThreshold != 4 {
		t.Fatalf("failure threshold = %d, want %d", cfg.FailureThreshold, 4)
	}
	if cfg.Cooldown != 50*time.Second {
		t.Fatalf("cooldown = %s, want %s", cfg.Cooldown, 50*time.Second)
	}
	if cfg.HealthCheckInterval != 40*time.Second {
		t.Fatalf("health check interval = %s, want %s", cfg.HealthCheckInterval, 40*time.Second)
	}
	if cfg.HealthCheckTimeout != 6*time.Second {
		t.Fatalf("health check timeout = %s, want %s", cfg.HealthCheckTimeout, 6*time.Second)
	}
}

func TestResolveAdminTokenAllowsNonLoopbackWithoutToken(t *testing.T) {
	t.Setenv("PSWITCH_ADMIN_TOKEN", "")

	got, err := resolveAdminToken("0.0.0.0:8080")
	if err != nil {
		t.Fatal(err)
	}
	if got != "" {
		t.Fatalf("token = %q, want empty", got)
	}
}

func TestDefaultConfigPathUsesBinaryDirectory(t *testing.T) {
	binaryPath := filepath.Join("/tmp", "pswitch-bin", "pswitch")

	got := defaultConfigPath(binaryPath)
	want := filepath.Join("/tmp", "pswitch-bin", "config.toml")
	if got != want {
		t.Fatalf("defaultConfigPath(%q) = %q, want %q", binaryPath, got, want)
	}
}

func TestDefaultStatePathsUseWorkingDirectory(t *testing.T) {
	dir := t.TempDir()
	prev, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.Chdir(prev)
	}()

	stateDir, err := defaultStateDir()
	if err != nil {
		t.Fatal(err)
	}
	if stateDir != dir {
		t.Fatalf("state dir = %q, want %q", stateDir, dir)
	}
	if got, want := defaultSettingsPath(stateDir), filepath.Join(dir, "settings.json"); got != want {
		t.Fatalf("settings path = %q, want %q", got, want)
	}
	if got, want := defaultMetricsPath(stateDir), filepath.Join(dir, "metrics.json"); got != want {
		t.Fatalf("metrics path = %q, want %q", got, want)
	}
}

func boolPtr(v bool) *bool {
	return &v
}

func sameBoolPtr(a, b *bool) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	return *a == *b
}

func intPtr(v int) *int {
	return &v
}

func durationPtr(v time.Duration) *time.Duration {
	return &v
}
