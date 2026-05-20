package server

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
	pruntime "pswitch/internal/runtime"
)

func TestRouterSupportsCodexAndClaudePrefixes(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/responses":
			w.Header().Set("Content-Type", "application/json")
			_, _ = io.WriteString(w, `{"id":"resp_1","status":"completed","output":[{"type":"message","role":"assistant","content":[{"type":"output_text","text":"hello"}]}],"usage":{"input_tokens":1,"output_tokens":2,"total_tokens":3}}`)
		case "/v1/models":
			w.Header().Set("Content-Type", "application/json")
			_, _ = io.WriteString(w, `{"data":[]}`)
		default:
			http.NotFound(w, r)
		}
	}))
	defer upstream.Close()

	cfg := config.Config{
		Listen:              "127.0.0.1:0",
		Mode:                "round_robin",
		FailureThreshold:    1,
		Cooldown:            time.Second,
		HealthCheckInterval: time.Second,
		HealthCheckTimeout:  time.Second,
		Routes: []config.Route{
			{Prefix: "/codex", Kind: "openai", Enabled: true},
			{Prefix: "/claude", Kind: "anthropic", Model: "claude-sonnet-4-20250514", UpstreamModel: "gpt-5.4", Enabled: true},
		},
		Providers: []config.Provider{
			{Name: "only", BaseURL: upstream.URL, APIKey: "k1", Enabled: true},
		},
	}

	dir := t.TempDir()
	manager, err := pruntime.New(dir+"/settings.json", dir+"/metrics.json", cfg)
	if err != nil {
		t.Fatal(err)
	}

	srv := httptest.NewServer(NewRouter(manager, ""))
	defer srv.Close()

	openAIReqBody, _ := json.Marshal(map[string]any{"model": "gpt-5.4", "input": "hello"})
	openAIResp, err := http.Post(srv.URL+"/codex/v1/responses", "application/json", bytes.NewReader(openAIReqBody))
	if err != nil {
		t.Fatal(err)
	}
	_ = openAIResp.Body.Close()
	if got, want := openAIResp.StatusCode, http.StatusOK; got != want {
		t.Fatalf("codex status = %d, want %d", got, want)
	}

	claudeReqBody := `{"model":"claude-sonnet-4-20250514","max_tokens":128,"messages":[{"role":"user","content":"hello"}]}`
	claudeResp, err := http.Post(srv.URL+"/claude/v1/messages", "application/json", bytes.NewBufferString(claudeReqBody))
	if err != nil {
		t.Fatal(err)
	}
	defer claudeResp.Body.Close()
	if got, want := claudeResp.StatusCode, http.StatusOK; got != want {
		t.Fatalf("claude status = %d, want %d", got, want)
	}
	body, _ := io.ReadAll(claudeResp.Body)
	if !bytes.Contains(body, []byte(`"type":"message"`)) {
		t.Fatalf("claude body = %q", string(body))
	}
}

func TestRouterServesAdminStateAndUI(t *testing.T) {
	cfg := config.Config{
		Listen:              "127.0.0.1:0",
		Mode:                "round_robin",
		FailureThreshold:    1,
		Cooldown:            time.Second,
		HealthCheckInterval: time.Second,
		HealthCheckTimeout:  time.Second,
		Routes: []config.Route{
			{Prefix: "/codex", Kind: "openai", Enabled: true},
		},
		Providers: []config.Provider{
			{Name: "only", BaseURL: "http://127.0.0.1:10001", APIKey: "k1", Enabled: true},
		},
	}

	dir := t.TempDir()
	manager, err := pruntime.New(dir+"/settings.json", dir+"/metrics.json", cfg)
	if err != nil {
		t.Fatal(err)
	}

	srv := httptest.NewServer(NewRouter(manager, ""))
	defer srv.Close()

	uiResp, err := http.Get(srv.URL + "/dashboard/")
	if err != nil {
		t.Fatal(err)
	}
	defer uiResp.Body.Close()
	if got, want := uiResp.StatusCode, http.StatusOK; got != want {
		t.Fatalf("ui status = %d, want %d", got, want)
	}

	stateResp, err := http.Get(srv.URL + "/dashboard/api/state")
	if err != nil {
		t.Fatal(err)
	}
	defer stateResp.Body.Close()
	if got, want := stateResp.StatusCode, http.StatusOK; got != want {
		t.Fatalf("state status = %d, want %d", got, want)
	}
}

func TestConfigSaveUpdatesRoutesForNewRequests(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/responses":
			w.Header().Set("Content-Type", "application/json")
			_, _ = io.WriteString(w, `{"id":"resp_1","status":"completed","output":[{"type":"message","role":"assistant","content":[{"type":"output_text","text":"hello"}]}],"usage":{"input_tokens":1,"output_tokens":2,"total_tokens":3}}`)
		case "/v1/models":
			w.Header().Set("Content-Type", "application/json")
			_, _ = io.WriteString(w, `{"data":[]}`)
		default:
			http.NotFound(w, r)
		}
	}))
	defer upstream.Close()

	cfg := config.Config{
		Listen:              "127.0.0.1:0",
		Mode:                "round_robin",
		FailureThreshold:    1,
		Cooldown:            time.Second,
		HealthCheckInterval: time.Second,
		HealthCheckTimeout:  time.Second,
		Routes: []config.Route{
			{Prefix: "/codex", Kind: "openai", Enabled: true},
		},
		Providers: []config.Provider{
			{Name: "only", BaseURL: upstream.URL, APIKey: "k1", Enabled: true},
		},
	}

	configPath := t.TempDir() + "/config.toml"
	if err := config.Write(configPath, cfg); err != nil {
		t.Fatal(err)
	}
	manager, err := pruntime.New(filepath.Join(filepath.Dir(configPath), "settings.json"), filepath.Join(filepath.Dir(configPath), "metrics.json"), cfg)
	if err != nil {
		t.Fatal(err)
	}

	srv := httptest.NewServer(NewRouter(manager, ""))
	defer srv.Close()

	reqBody, _ := json.Marshal(map[string]any{"model": "gpt-5.4", "input": "hello"})
	resp, err := http.Post(srv.URL+"/codex/v1/responses", "application/json", bytes.NewReader(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	_ = resp.Body.Close()
	if got, want := resp.StatusCode, http.StatusOK; got != want {
		t.Fatalf("before save status = %d, want %d", got, want)
	}

	savePayload := map[string]any{
		"listen":                "127.0.0.1:0",
		"mode":                  "round_robin",
		"failure_threshold":     1,
		"cooldown":              "1s",
		"health_check_interval": "1s",
		"health_check_timeout":  "1s",
		"routes": []map[string]any{
			{"prefix": "/openai", "type": "openai", "enabled": true},
		},
		"providers": []map[string]any{
			{"name": "only", "base_url": upstream.URL, "api_key": "k1", "enabled": true},
		},
	}
	saveBody, _ := json.Marshal(savePayload)
	saveReq, err := http.NewRequest(http.MethodPut, srv.URL+"/dashboard/api/config", bytes.NewReader(saveBody))
	if err != nil {
		t.Fatal(err)
	}
	saveReq.Header.Set("Content-Type", "application/json")
	saveResp, err := http.DefaultClient.Do(saveReq)
	if err != nil {
		t.Fatal(err)
	}
	_ = saveResp.Body.Close()
	if got, want := saveResp.StatusCode, http.StatusOK; got != want {
		t.Fatalf("save status = %d, want %d", got, want)
	}

	oldResp, err := http.Post(srv.URL+"/codex/v1/responses", "application/json", bytes.NewReader(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	_ = oldResp.Body.Close()
	if got, want := oldResp.StatusCode, http.StatusNotFound; got != want {
		t.Fatalf("old route status = %d, want %d", got, want)
	}

	newResp, err := http.Post(srv.URL+"/openai/v1/responses", "application/json", bytes.NewReader(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	_ = newResp.Body.Close()
	if got, want := newResp.StatusCode, http.StatusOK; got != want {
		t.Fatalf("new route status = %d, want %d", got, want)
	}
}

func TestRouterRecordsUsageIntoDashboardStats(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/responses":
			w.Header().Set("Content-Type", "application/json")
			_, _ = io.WriteString(w, `{"id":"resp_1","usage":{"input_tokens":12,"output_tokens":34,"total_tokens":46}}`)
		case "/v1/models":
			w.Header().Set("Content-Type", "application/json")
			_, _ = io.WriteString(w, `{"data":[]}`)
		default:
			http.NotFound(w, r)
		}
	}))
	defer upstream.Close()

	cfg := config.Config{
		Listen:              "127.0.0.1:0",
		Mode:                "round_robin",
		FailureThreshold:    1,
		Cooldown:            time.Second,
		HealthCheckInterval: time.Second,
		HealthCheckTimeout:  time.Second,
		Routes: []config.Route{
			{Prefix: "/codex", Kind: "openai", Enabled: true},
		},
		Providers: []config.Provider{
			{Name: "only", BaseURL: upstream.URL, APIKey: "k1", Enabled: true},
		},
	}

	configPath := t.TempDir() + "/config.toml"
	if err := config.Write(configPath, cfg); err != nil {
		t.Fatal(err)
	}

	manager, err := pruntime.New(filepath.Join(filepath.Dir(configPath), "settings.json"), filepath.Join(filepath.Dir(configPath), "metrics.json"), cfg)
	if err != nil {
		t.Fatal(err)
	}

	srv := httptest.NewServer(NewRouter(manager, ""))
	defer srv.Close()

	reqBody, _ := json.Marshal(map[string]any{"model": "gpt-5.4", "input": "hello"})
	resp, err := http.Post(srv.URL+"/codex/v1/responses", "application/json", bytes.NewReader(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()

	statsResp, err := http.Get(srv.URL + "/dashboard/api/stats")
	if err != nil {
		t.Fatal(err)
	}
	defer statsResp.Body.Close()

	var got struct {
		Overview struct {
			TotalRequests int64 `json:"total_requests"`
			TotalTokens   int64 `json:"total_tokens"`
		} `json:"overview"`
	}
	if err := json.NewDecoder(statsResp.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}
	if got.Overview.TotalRequests != 1 || got.Overview.TotalTokens != 46 {
		t.Fatalf("stats overview = %#v", got.Overview)
	}
}
