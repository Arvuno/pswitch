package openai

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"pswitch/internal/logx"
	"pswitch/internal/metrics"
	"pswitch/internal/pool"
	"pswitch/internal/upstream"
)

type Options struct {
	Client  *http.Client
	Mode    pool.Mode
	Metrics *metrics.Store
}

type Handler struct {
	pool    *pool.Pool
	client  *http.Client
	mode    pool.Mode
	metrics *metrics.Store
}

func NewHandler(providerPool *pool.Pool, opts Options) *Handler {
	client := opts.Client
	if client == nil {
		client = &http.Client{}
	}
	mode := opts.Mode
	if mode == "" {
		mode = providerPool.Mode()
	}
	return &Handler{
		pool:    providerPool,
		client:  client,
		mode:    mode,
		metrics: opts.Metrics,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := upstream.ReadRequestBody(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("read request body: %v", err), http.StatusBadRequest)
		return
	}

	now := time.Now()
	model := upstream.ExtractRequestedModel(body)
	candidates := h.pool.Candidates(h.mode, now)
	if len(candidates) == 0 {
		http.Error(w, "no available providers", http.StatusServiceUnavailable)
		return
	}

	var lastErr error
	var lastResp *upstream.StoredResponse

	for _, candidate := range candidates {
		resp, err := h.forward(r, candidate, body)
		if err != nil {
			lastErr = err
			logx.Warnf("upstream request failed provider=%s err=%v", candidate.Name, err)
			h.pool.MarkFailure(candidate.Name, time.Now())
			h.recordFailure(candidate.Name, model)
			continue
		}

		if upstream.ShouldFailover(resp.StatusCode) {
			stored, storeErr := upstream.CaptureResponse(resp)
			if storeErr != nil {
				lastErr = storeErr
				logx.Warnf("failed to read upstream failure response provider=%s err=%v", candidate.Name, storeErr)
			} else {
				lastResp = stored
			}
			logx.Warnf("upstream returned failover status provider=%s status=%d", candidate.Name, resp.StatusCode)
			h.pool.MarkFailure(candidate.Name, time.Now())
			h.recordFailure(candidate.Name, model)
			continue
		}

		h.pool.MarkSuccess(candidate.Name, time.Now())
		logx.Debugf("upstream request succeeded provider=%s status=%d", candidate.Name, resp.StatusCode)
		if isStreaming(resp.Header) {
			collector := &upstream.StreamUsageCollector{}
			defer resp.Body.Close()
			upstream.CopyHeaders(w.Header(), resp.Header)
			w.WriteHeader(resp.StatusCode)
			if err := upstream.CopyResponseBody(w, resp.Body, true, collector); err != nil {
				if usage, ok := collector.Usage(); ok {
					logUsage(candidate.Name, usage)
					h.recordSuccess(candidate.Name, model, usage)
				} else {
					logx.Warnf("stream forwarding interrupted before usage arrived provider=%s err=%v", candidate.Name, err)
				}
				return
			}
			if usage, ok := collector.Usage(); ok {
				logUsage(candidate.Name, usage)
				h.recordSuccess(candidate.Name, model, usage)
			} else {
				logx.Warnf("stream completed without usage provider=%s status=%d", candidate.Name, resp.StatusCode)
				h.recordSuccess(candidate.Name, model, upstream.UsageSummary{})
			}
			return
		}

		respBody, readErr := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if readErr != nil {
			http.Error(w, fmt.Sprintf("read upstream body: %v", readErr), http.StatusBadGateway)
			return
		}

		if usage, ok := upstream.ExtractUsage(respBody); ok {
			logUsage(candidate.Name, usage)
			h.recordSuccess(candidate.Name, model, usage)
		} else {
			logx.Debugf("usage unavailable provider=%s status=%d", candidate.Name, resp.StatusCode)
			h.recordSuccess(candidate.Name, model, upstream.UsageSummary{})
		}

		upstream.CopyHeaders(w.Header(), resp.Header)
		w.WriteHeader(resp.StatusCode)
		if _, err := w.Write(respBody); err != nil {
			return
		}
		return
	}

	if lastResp != nil {
		upstream.CopyHeaders(w.Header(), lastResp.Header)
		w.WriteHeader(lastResp.StatusCode)
		_, _ = w.Write(lastResp.Body)
		return
	}
	if lastErr != nil {
		http.Error(w, fmt.Sprintf("all providers failed: %v", lastErr), http.StatusBadGateway)
		return
	}

	http.Error(w, "all providers failed", http.StatusBadGateway)
}

func (h *Handler) forward(original *http.Request, candidate pool.ProviderSnapshot, body []byte) (*http.Response, error) {
	upstreamURL := buildTargetURL(candidate.BaseURL, original.URL.Path, original.URL.RawQuery)
	req, err := http.NewRequestWithContext(original.Context(), original.Method, upstreamURL.String(), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	upstream.CopyHeaders(req.Header, original.Header)
	req.Host = upstreamURL.Host
	req.Header.Del("Authorization")
	req.Header.Set("Authorization", "Bearer "+candidate.APIKey)

	return h.client.Do(req)
}

func (h *Handler) recordSuccess(provider, model string, usage upstream.UsageSummary) {
	if h.metrics == nil {
		return
	}
	h.metrics.RecordSuccess(provider, metrics.Usage{
		RequestCount: 1,
		InputTokens:  usage.InputTokens,
		OutputTokens: usage.OutputTokens,
		TotalTokens:  usage.TotalTokens,
	}, model, time.Now())
}

func (h *Handler) recordFailure(provider, model string) {
	if h.metrics == nil {
		return
	}
	h.metrics.RecordFailure(provider, model, time.Now())
}

func isStreaming(header http.Header) bool {
	return strings.Contains(strings.ToLower(header.Get("Content-Type")), "text/event-stream")
}

func buildTargetURL(base *url.URL, path, rawQuery string) *url.URL {
	target := *base
	target.Path = upstream.JoinPaths(base.Path, path)
	target.RawQuery = rawQuery
	target.Fragment = ""
	return &target
}

func logUsage(provider string, usage upstream.UsageSummary) {
	logx.Infof(
		"usage provider=%s input_tokens=%d output_tokens=%d total_tokens=%d",
		provider,
		usage.InputTokens,
		usage.OutputTokens,
		usage.TotalTokens,
	)
}
