package main

import (
	"os"
	"path/filepath"
	"testing"

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

func TestLoadStartupConfigFallsBackToDefaultWhenMissing(t *testing.T) {
	path := t.TempDir() + "/missing.toml"

	cfg, err := loadStartupConfig(path, filepath.Join(t.TempDir(), "settings.json"))
	if err != nil {
		t.Fatal(err)
	}

	if got, want := cfg.Listen, "127.0.0.1:8080"; got != want {
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
