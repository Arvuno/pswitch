package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"pswitch/internal/config"
	"pswitch/internal/logx"
	"pswitch/internal/metrics"
	"pswitch/internal/pool"
)

func TestHandlerFailsOverToSecondProvider(t *testing.T) {
	firstReceived := make(chan *http.Request, 1)
	secondReceived := make(chan *http.Request, 1)

	first := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		firstReceived <- r.Clone(r.Context())
		http.Error(w, "broken", http.StatusServiceUnavailable)
	}))
	defer first.Close()

	second := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secondReceived <- r.Clone(r.Context())
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"ok":true}`)
	}))
	defer second.Close()

	p, err := pool.New(config.Config{
		Mode:             "sequential",
		FailureThreshold: 1,
		Cooldown:         time.Second,
		Providers: []config.Provider{
			{Name: "first", BaseURL: first.URL, APIKey: "k1"},
			{Name: "second", BaseURL: second.URL, APIKey: "k2"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	h := NewHandler(p, Options{})
	srv := httptest.NewServer(h)
	defer srv.Close()

	reqBody, _ := json.Marshal(map[string]any{"model": "gpt-5.4", "input": "hello"})
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, srv.URL+"/v1/responses", bytes.NewReader(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if got, want := resp.StatusCode, http.StatusOK; got != want {
		t.Fatalf("status = %d, want %d", got, want)
	}
	body, _ := io.ReadAll(resp.Body)
	if string(body) != `{"ok":true}` {
		t.Fatalf("body = %q, want %q", string(body), `{"ok":true}`)
	}

	select {
	case r := <-firstReceived:
		if got, want := r.Header.Get("Authorization"), "Bearer k1"; got != want {
			t.Fatalf("first auth = %q, want %q", got, want)
		}
	default:
		t.Fatal("first provider was not called")
	}

	select {
	case r := <-secondReceived:
		if got, want := r.Header.Get("Authorization"), "Bearer k2"; got != want {
			t.Fatalf("second auth = %q, want %q", got, want)
		}
	default:
		t.Fatal("second provider was not called")
	}
}

func TestHandlerLogsUsageFromJSONResponse(t *testing.T) {
	var logBuf bytes.Buffer
	logger := testLogger(&logBuf)
	restore := logx.Use(logger)
	defer restore()

	upstreamServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"id":"resp_1","usage":{"input_tokens":12,"output_tokens":34,"total_tokens":46}}`)
	}))
	defer upstreamServer.Close()

	p, err := pool.New(config.Config{
		Mode:             "sequential",
		FailureThreshold: 1,
		Cooldown:         time.Second,
		Providers: []config.Provider{
			{Name: "only", BaseURL: upstreamServer.URL, APIKey: "k1"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	h := NewHandler(p, Options{})
	srv := httptest.NewServer(h)
	defer srv.Close()

	reqBody, _ := json.Marshal(map[string]any{"model": "gpt-5.4", "input": "hello"})
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, srv.URL+"/v1/responses", bytes.NewReader(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()

	if got := logBuf.String(); !strings.Contains(got, "input_tokens=12") || !strings.Contains(got, "output_tokens=34") || !strings.Contains(got, "total_tokens=46") {
		t.Fatalf("log output = %q, want token usage fields", got)
	}
}

func TestHandlerLogsUsageFromStreamingResponseCompletedEvent(t *testing.T) {
	var logBuf bytes.Buffer
	logger := testLogger(&logBuf)
	restore := logx.Use(logger)
	defer restore()

	upstreamServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, "event: response.completed\n")
		_, _ = io.WriteString(w, "data: {\"type\":\"response.completed\",\"response\":{\"usage\":{\"input_tokens\":21,\"output_tokens\":13,\"total_tokens\":34}}}\n\n")
	}))
	defer upstreamServer.Close()

	p, err := pool.New(config.Config{
		Mode:             "sequential",
		FailureThreshold: 1,
		Cooldown:         time.Second,
		Providers: []config.Provider{
			{Name: "only", BaseURL: upstreamServer.URL, APIKey: "k1"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	h := NewHandler(p, Options{})
	srv := httptest.NewServer(h)
	defer srv.Close()

	reqBody, _ := json.Marshal(map[string]any{"model": "gpt-5.4", "input": "hello", "stream": true})
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, srv.URL+"/v1/responses", bytes.NewReader(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()

	if got := waitForLog(&logBuf, 500*time.Millisecond); !strings.Contains(got, "input_tokens=21") || !strings.Contains(got, "output_tokens=13") || !strings.Contains(got, "total_tokens=34") {
		t.Fatalf("log output = %q, want streaming token usage fields", got)
	}
}

func TestHandlerRecordsMetricsOnSuccessAndFailure(t *testing.T) {
	metricsStore, err := metrics.Open(t.TempDir() + "/metrics.json")
	if err != nil {
		t.Fatal(err)
	}

	first := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "broken", http.StatusServiceUnavailable)
	}))
	defer first.Close()

	second := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"id":"resp_1","usage":{"input_tokens":12,"output_tokens":34,"total_tokens":46}}`)
	}))
	defer second.Close()

	p, err := pool.New(config.Config{
		Mode:             "sequential",
		FailureThreshold: 1,
		Cooldown:         time.Second,
		Providers: []config.Provider{
			{Name: "first", BaseURL: first.URL, APIKey: "k1"},
			{Name: "second", BaseURL: second.URL, APIKey: "k2"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	h := NewHandler(p, Options{Metrics: metricsStore})
	srv := httptest.NewServer(h)
	defer srv.Close()

	reqBody, _ := json.Marshal(map[string]any{"model": "gpt-5.4", "input": "hello"})
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, srv.URL+"/v1/responses", bytes.NewReader(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()

	snapshot := metricsStore.Snapshot(time.Now())
	if got, want := snapshot.Overview.TotalFailures, int64(1); got != want {
		t.Fatalf("total failures = %d, want %d", got, want)
	}
	if got, want := snapshot.Overview.TotalRequests, int64(1); got != want {
		t.Fatalf("total requests = %d, want %d", got, want)
	}
	if got, want := snapshot.Provider("second").TotalTokens, int64(46); got != want {
		t.Fatalf("second provider total tokens = %d, want %d", got, want)
	}
	if got, want := snapshot.Model("gpt-5.4").TotalTokens, int64(46); got != want {
		t.Fatalf("gpt-5.4 total tokens = %d, want %d", got, want)
	}
	if got, want := snapshot.Model("gpt-5.4").FailureCount, int64(1); got != want {
		t.Fatalf("gpt-5.4 failures = %d, want %d", got, want)
	}
}

func TestHandlerRecordsStreamingUsageEvenIfClientWriteFailsAfterUsageArrives(t *testing.T) {
	var logBuf bytes.Buffer
	logger := testLogger(&logBuf)
	restore := logx.Use(logger)
	defer restore()

	metricsStore, err := metrics.Open(t.TempDir() + "/metrics.json")
	if err != nil {
		t.Fatal(err)
	}

	upstreamServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, "event: response.completed\n")
		_, _ = io.WriteString(w, "data: {\"type\":\"response.completed\",\"response\":{\"usage\":{\"input_tokens\":21,\"output_tokens\":13,\"total_tokens\":34}}}\n\n")
	}))
	defer upstreamServer.Close()

	p, err := pool.New(config.Config{
		Mode:             "sequential",
		FailureThreshold: 1,
		Cooldown:         time.Second,
		Providers: []config.Provider{
			{Name: "only", BaseURL: upstreamServer.URL, APIKey: "k1"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	h := NewHandler(p, Options{Metrics: metricsStore})
	reqBody, _ := json.Marshal(map[string]any{"model": "gpt-5.4", "input": "hello", "stream": true})
	req := httptest.NewRequest(http.MethodPost, "/v1/responses", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	writer := &failingResponseWriter{header: make(http.Header), err: errors.New("client disconnected")}
	h.ServeHTTP(writer, req)

	snapshot := metricsStore.Snapshot(time.Now())
	if got, want := snapshot.Overview.TotalRequests, int64(1); got != want {
		t.Fatalf("total requests = %d, want %d", got, want)
	}
	if got, want := snapshot.Model("gpt-5.4").TotalTokens, int64(34); got != want {
		t.Fatalf("gpt-5.4 total tokens = %d, want %d", got, want)
	}
	if got := waitForLog(&logBuf, 500*time.Millisecond); !strings.Contains(got, "total_tokens=34") {
		t.Fatalf("log output = %q, want usage logged before stream write failure", got)
	}
}

func testLogger(buf *bytes.Buffer) *zap.Logger {
	encoderCfg := zap.NewDevelopmentEncoderConfig()
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderCfg),
		zapcore.AddSync(buf),
		zapcore.DebugLevel,
	)
	return zap.New(core)
}

func waitForLog(buf *bytes.Buffer, timeout time.Duration) string {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		got := buf.String()
		if strings.Contains(got, "total_tokens=") {
			return got
		}
		time.Sleep(10 * time.Millisecond)
	}
	return buf.String()
}

type failingResponseWriter struct {
	header http.Header
	status int
	err    error
}

func (w *failingResponseWriter) Header() http.Header {
	return w.header
}

func (w *failingResponseWriter) WriteHeader(status int) {
	w.status = status
}

func (w *failingResponseWriter) Write([]byte) (int, error) {
	return 0, w.err
}

func (w *failingResponseWriter) Flush() {}
