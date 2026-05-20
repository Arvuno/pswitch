package anthropic

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMessagesTransformsOpenAIResponse(t *testing.T) {
	upstream := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Path, "/v1/responses"; got != want {
			t.Fatalf("path = %q, want %q", got, want)
		}

		var req map[string]any
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatal(err)
		}
		if got, want := req["model"], "gpt-5.4"; got != want {
			t.Fatalf("model = %v, want %v", got, want)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"id":"resp_1","status":"completed","output":[{"type":"message","role":"assistant","content":[{"type":"output_text","text":"hello from upstream"}]}],"usage":{"input_tokens":12,"output_tokens":34,"total_tokens":46}}`)
	})

	h := NewHandler(Options{
		Model:         "claude-sonnet-4-20250514",
		UpstreamModel: "gpt-5.4",
		Upstream:      upstream,
	})

	reqBody := `{"model":"claude-sonnet-4-20250514","max_tokens":1024,"messages":[{"role":"user","content":"hi"}]}`
	req := httptest.NewRequest(http.MethodPost, "/v1/messages", strings.NewReader(reqBody))
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if got, want := rec.Code, http.StatusOK; got != want {
		t.Fatalf("status = %d, want %d", got, want)
	}

	var resp struct {
		Model   string `json:"model"`
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
		Usage struct {
			InputTokens  int64 `json:"input_tokens"`
			OutputTokens int64 `json:"output_tokens"`
		} `json:"usage"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if got, want := resp.Model, "claude-sonnet-4-20250514"; got != want {
		t.Fatalf("model = %q, want %q", got, want)
	}
	if len(resp.Content) != 1 || resp.Content[0].Text != "hello from upstream" {
		t.Fatalf("content = %+v", resp.Content)
	}
	if got, want := resp.Usage.InputTokens, int64(12); got != want {
		t.Fatalf("input_tokens = %d, want %d", got, want)
	}
}

func TestCountTokensReturnsEstimate(t *testing.T) {
	h := NewHandler(Options{
		Model:         "claude-sonnet-4-20250514",
		UpstreamModel: "gpt-5.4",
		Upstream:      http.NotFoundHandler(),
	})

	reqBody := `{"model":"claude-sonnet-4-20250514","messages":[{"role":"user","content":"hello world"}]}`
	req := httptest.NewRequest(http.MethodPost, "/v1/messages/count_tokens", strings.NewReader(reqBody))
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if got, want := rec.Code, http.StatusOK; got != want {
		t.Fatalf("status = %d, want %d", got, want)
	}
	if !strings.Contains(rec.Body.String(), "input_tokens") {
		t.Fatalf("body = %q, want input_tokens field", rec.Body.String())
	}
}

func TestStreamingReturnsAnthropicEvents(t *testing.T) {
	upstream := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"id":"resp_1","status":"completed","output":[{"type":"message","role":"assistant","content":[{"type":"output_text","text":"streamed text"}]}],"usage":{"input_tokens":1,"output_tokens":2,"total_tokens":3}}`)
	})

	h := NewHandler(Options{
		Model:         "claude-sonnet-4-20250514",
		UpstreamModel: "gpt-5.4",
		Upstream:      upstream,
	})

	reqBody := `{"model":"claude-sonnet-4-20250514","max_tokens":1024,"stream":true,"messages":[{"role":"user","content":"hi"}]}`
	req := httptest.NewRequest(http.MethodPost, "/v1/messages", strings.NewReader(reqBody))
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if got := rec.Header().Get("Content-Type"); !strings.Contains(got, "text/event-stream") {
		t.Fatalf("content-type = %q, want event stream", got)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "event: message_start") || !strings.Contains(body, "event: message_stop") {
		t.Fatalf("body = %q, want anthropic stream events", body)
	}
}

func TestModelsListsConfiguredModel(t *testing.T) {
	h := NewHandler(Options{
		Model:         "claude-sonnet-4-20250514",
		UpstreamModel: "gpt-5.4",
		Upstream:      http.NotFoundHandler(),
	})

	req := httptest.NewRequest(http.MethodGet, "/v1/models", nil)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if got, want := rec.Code, http.StatusOK; got != want {
		t.Fatalf("status = %d, want %d", got, want)
	}
	if !bytes.Contains(rec.Body.Bytes(), []byte("claude-sonnet-4-20250514")) {
		t.Fatalf("body = %q, want configured model", rec.Body.String())
	}
}
