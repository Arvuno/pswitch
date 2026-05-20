package adminui

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"pswitch/internal/config"
	"pswitch/internal/metrics"
	pruntime "pswitch/internal/runtime"
)

func TestHandlerServesEmbeddedIndex(t *testing.T) {
	manager := newTestManager(t)

	srv := httptest.NewServer(New(manager, ""))
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}
	body, _ := io.ReadAll(resp.Body)
	if !bytes.Contains(body, []byte("pswitch admin")) {
		t.Fatalf("body = %q, want embedded admin page", string(body))
	}
	if !bytes.Contains(body, []byte("config-tab-button")) {
		t.Fatalf("body = %q, want config tab controls in embedded admin page", string(body))
	}
	if !bytes.Contains(body, []byte("provider-summary-grid")) {
		t.Fatalf("body = %q, want provider summary panel in embedded admin page", string(body))
	}
	if !bytes.Contains(body, []byte("model-usage-grid")) {
		t.Fatalf("body = %q, want model usage panel in embedded admin page", string(body))
	}
	if !bytes.Contains(body, []byte("icon.svg")) {
		t.Fatalf("body = %q, want icon link in embedded admin page", string(body))
	}
	if !bytes.Contains(body, []byte("least_failures")) {
		t.Fatalf("body = %q, want least_failures mode in embedded admin page", string(body))
	}
}

func TestHandlerServesIconAsset(t *testing.T) {
	manager := newTestManager(t)

	srv := httptest.NewServer(New(manager, ""))
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/icon.svg")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}
	if got, want := resp.Header.Get("Content-Type"), "image/svg+xml"; got != want {
		t.Fatalf("content-type = %q, want %q", got, want)
	}
}

func TestHandlerRejectsUnauthorizedStateRequest(t *testing.T) {
	manager := newTestManager(t)

	srv := httptest.NewServer(New(manager, "secret"))
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/api/state")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusUnauthorized)
	}
}

func TestHandlerMetaExposesDashboardPrefix(t *testing.T) {
	manager := newTestManager(t)

	srv := httptest.NewServer(New(manager, ""))
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/api/meta")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	var got struct {
		AdminPrefix string `json:"admin_prefix"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}
	if got.AdminPrefix != "/dashboard" {
		t.Fatalf("admin_prefix = %q, want %q", got.AdminPrefix, "/dashboard")
	}
}

func TestHandlerSavesConfigAndReturnsRestartNotice(t *testing.T) {
	manager := newTestManager(t)

	srv := httptest.NewServer(New(manager, "secret"))
	defer srv.Close()

	payload := map[string]any{
		"listen":                "0.0.0.0:8080",
		"mode":                  "round_robin",
		"failure_threshold":     1,
		"cooldown":              "20s",
		"health_check_interval": "15s",
		"health_check_timeout":  "3s",
		"routes": []map[string]any{
			{"prefix": "/codex", "type": "openai", "enabled": true},
		},
		"providers": []map[string]any{
			{"name": "one", "base_url": "http://127.0.0.1:10001", "api_key": "k1", "enabled": true},
		},
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequest(http.MethodPut, srv.URL+"/api/config", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set(headerAdminToken, "secret")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("status = %d, want %d body=%s", resp.StatusCode, http.StatusOK, string(body))
	}

	var got struct {
		RequiresRestart bool     `json:"requires_restart"`
		Messages        []string `json:"messages"`
		State           struct {
			Config struct {
				Listen string `json:"listen"`
			} `json:"config"`
		} `json:"state"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}
	if !got.RequiresRestart {
		t.Fatal("save should report restart requirement")
	}
	if got.State.Config.Listen != "0.0.0.0:8080" {
		t.Fatalf("listen = %q, want %q", got.State.Config.Listen, "0.0.0.0:8080")
	}
	if len(got.Messages) == 0 {
		t.Fatal("expected restart warning message")
	}
}

func TestHandlerStatsReturnsMetricsAndProviderHealth(t *testing.T) {
	manager := newTestManager(t)

	now := time.Date(2026, 5, 20, 12, 0, 0, 0, time.UTC)
	manager.Metrics().RecordSuccess("one", metricsUsage(12, 34, 46), "gpt-5.4", now)
	manager.Metrics().RecordFailure("one", "gpt-5.4", now)
	manager.Pool().MarkFailure("one", now)

	srv := httptest.NewServer(New(manager, ""))
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/api/stats")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	var got struct {
		Overview struct {
			TotalRequests int64 `json:"total_requests"`
			TotalFailures int64 `json:"total_failures"`
			TotalTokens   int64 `json:"total_tokens"`
		} `json:"overview"`
		Providers []struct {
			Name         string `json:"name"`
			Healthy      bool   `json:"healthy"`
			TotalTokens  int64  `json:"total_tokens"`
			FailureCount int64  `json:"failure_count"`
		} `json:"providers"`
		Models []struct {
			Name        string `json:"name"`
			TotalTokens int64  `json:"total_tokens"`
		} `json:"models"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}
	if got.Overview.TotalRequests != 1 || got.Overview.TotalFailures != 1 || got.Overview.TotalTokens != 46 {
		t.Fatalf("overview = %#v", got.Overview)
	}
	if len(got.Providers) != 1 {
		t.Fatalf("providers = %d, want 1", len(got.Providers))
	}
	if got.Providers[0].Name != "one" || got.Providers[0].Healthy || got.Providers[0].TotalTokens != 46 || got.Providers[0].FailureCount != 1 {
		t.Fatalf("provider stats = %#v", got.Providers[0])
	}
	if len(got.Models) != 1 || got.Models[0].Name != "gpt-5.4" || got.Models[0].TotalTokens != 46 {
		t.Fatalf("model stats = %#v", got.Models)
	}
}

func newTestManager(t *testing.T) *pruntime.Manager {
	t.Helper()

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

	manager, err := pruntime.New(filepath.Join(dir, "settings.json"), filepath.Join(dir, "metrics.json"), cfg)
	if err != nil {
		t.Fatal(err)
	}
	return manager
}

func metricsUsage(input, output, total int64) metrics.Usage {
	return metrics.Usage{
		RequestCount: 1,
		InputTokens:  input,
		OutputTokens: output,
		TotalTokens:  total,
	}
}
