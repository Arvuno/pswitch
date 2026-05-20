package anthropic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"
)

type Options struct {
	Model         string
	UpstreamModel string
	Upstream      http.Handler
}

type Handler struct {
	model         string
	upstreamModel string
	upstream      http.Handler
}

func NewHandler(opts Options) *Handler {
	return &Handler{
		model:         strings.TrimSpace(opts.Model),
		upstreamModel: strings.TrimSpace(opts.UpstreamModel),
		upstream:      opts.Upstream,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet && r.URL.Path == "/v1/models":
		h.handleModels(w)
	case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/v1/models/"):
		h.handleModel(w, strings.TrimPrefix(r.URL.Path, "/v1/models/"))
	case r.Method == http.MethodPost && r.URL.Path == "/v1/messages":
		h.handleMessages(w, r)
	case r.Method == http.MethodPost && r.URL.Path == "/v1/messages/count_tokens":
		h.handleCountTokens(w, r)
	default:
		writeError(w, http.StatusNotFound, "not found")
	}
}

func (h *Handler) handleModels(w http.ResponseWriter) {
	model := h.publicModel("")
	writeJSON(w, http.StatusOK, map[string]any{
		"data": []map[string]any{
			{
				"id":           model,
				"type":         "model",
				"display_name": model,
				"created_at":   time.Now().UTC().Format(time.RFC3339),
			},
		},
		"first_id": model,
		"last_id":  model,
		"has_more": false,
	})
}

func (h *Handler) handleModel(w http.ResponseWriter, id string) {
	model := h.publicModel(id)
	writeJSON(w, http.StatusOK, map[string]any{
		"id":           model,
		"type":         "model",
		"display_name": model,
		"created_at":   time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *Handler) handleCountTokens(w http.ResponseWriter, r *http.Request) {
	req, err := decodeMessagesRequest(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"input_tokens": estimateTokens(buildTranscript(req.System, req.Messages)),
	})
}

func (h *Handler) handleMessages(w http.ResponseWriter, r *http.Request) {
	req, err := decodeMessagesRequest(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	openAIResp, err := h.invokeUpstream(r.Context(), req)
	if err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}

	if req.Stream {
		h.writeStream(w, req, openAIResp)
		return
	}

	writeJSON(w, http.StatusOK, h.buildMessageResponse(req, openAIResp))
}

func (h *Handler) invokeUpstream(ctx context.Context, req messagesRequest) (openAIResponse, error) {
	payload := map[string]any{
		"model":             h.upstreamModel,
		"input":             buildTranscript(req.System, req.Messages),
		"max_output_tokens": req.MaxTokens,
		"stream":            false,
	}
	if system := contentText(req.System); system != "" {
		payload["instructions"] = system
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return openAIResponse{}, err
	}

	upstreamReq, err := http.NewRequestWithContext(ctx, http.MethodPost, "/v1/responses", bytes.NewReader(body))
	if err != nil {
		return openAIResponse{}, err
	}
	upstreamReq.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	h.upstream.ServeHTTP(rec, upstreamReq)

	if rec.Code < 200 || rec.Code >= 300 {
		return openAIResponse{}, fmt.Errorf("upstream status %d: %s", rec.Code, strings.TrimSpace(rec.Body.String()))
	}

	var resp openAIResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		return openAIResponse{}, err
	}
	return resp, nil
}

func (h *Handler) buildMessageResponse(req messagesRequest, openAIResp openAIResponse) map[string]any {
	return map[string]any{
		"id":            anthropicMessageID(),
		"type":          "message",
		"role":          "assistant",
		"model":         h.publicModel(req.Model),
		"content":       []map[string]any{{"type": "text", "text": openAIResp.OutputText()}},
		"stop_reason":   openAIResp.StopReason(),
		"stop_sequence": nil,
		"usage": map[string]any{
			"input_tokens":  openAIResp.Usage.InputTokens,
			"output_tokens": openAIResp.Usage.OutputTokens,
		},
	}
}

func (h *Handler) writeStream(w http.ResponseWriter, req messagesRequest, openAIResp openAIResponse) {
	w.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.WriteHeader(http.StatusOK)

	messageID := anthropicMessageID()
	model := h.publicModel(req.Model)
	text := openAIResp.OutputText()
	flusher, _ := w.(http.Flusher)

	writeEvent(w, "message_start", map[string]any{
		"type": "message_start",
		"message": map[string]any{
			"id":            messageID,
			"type":          "message",
			"role":          "assistant",
			"model":         model,
			"content":       []any{},
			"stop_reason":   nil,
			"stop_sequence": nil,
			"usage": map[string]any{
				"input_tokens":  openAIResp.Usage.InputTokens,
				"output_tokens": 0,
			},
		},
	})
	flush(flusher)

	writeEvent(w, "content_block_start", map[string]any{
		"type":          "content_block_start",
		"index":         0,
		"content_block": map[string]any{"type": "text", "text": ""},
	})
	flush(flusher)

	for _, chunk := range splitChunks(text, 256) {
		writeEvent(w, "content_block_delta", map[string]any{
			"type":  "content_block_delta",
			"index": 0,
			"delta": map[string]any{
				"type": "text_delta",
				"text": chunk,
			},
		})
		flush(flusher)
	}

	writeEvent(w, "content_block_stop", map[string]any{
		"type":  "content_block_stop",
		"index": 0,
	})
	flush(flusher)

	writeEvent(w, "message_delta", map[string]any{
		"type": "message_delta",
		"delta": map[string]any{
			"stop_reason":   openAIResp.StopReason(),
			"stop_sequence": nil,
		},
		"usage": map[string]any{
			"input_tokens":  openAIResp.Usage.InputTokens,
			"output_tokens": openAIResp.Usage.OutputTokens,
		},
	})
	flush(flusher)

	writeEvent(w, "message_stop", map[string]any{"type": "message_stop"})
	flush(flusher)
}

func (h *Handler) publicModel(requested string) string {
	requested = strings.TrimSpace(requested)
	if requested != "" {
		return requested
	}
	if h.model != "" {
		return h.model
	}
	if h.upstreamModel != "" {
		return h.upstreamModel
	}
	return "claude"
}

type messagesRequest struct {
	Model     string           `json:"model"`
	System    any              `json:"system"`
	Messages  []messageContent `json:"messages"`
	MaxTokens int              `json:"max_tokens"`
	Stream    bool             `json:"stream"`
}

type messageContent struct {
	Role    string `json:"role"`
	Content any    `json:"content"`
}

type openAIResponse struct {
	Output []struct {
		Type    string `json:"type"`
		Role    string `json:"role"`
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	} `json:"output"`
	Usage struct {
		InputTokens  int64 `json:"input_tokens"`
		OutputTokens int64 `json:"output_tokens"`
		TotalTokens  int64 `json:"total_tokens"`
	} `json:"usage"`
	Status string `json:"status"`
}

func (r openAIResponse) OutputText() string {
	var parts []string
	for _, item := range r.Output {
		for _, content := range item.Content {
			if strings.TrimSpace(content.Text) != "" {
				parts = append(parts, content.Text)
			}
		}
	}
	return strings.Join(parts, "")
}

func (r openAIResponse) StopReason() string {
	if strings.TrimSpace(r.Status) == "completed" {
		return "end_turn"
	}
	return "stop_sequence"
}

func decodeMessagesRequest(body io.Reader) (messagesRequest, error) {
	var req messagesRequest
	if err := json.NewDecoder(body).Decode(&req); err != nil {
		return messagesRequest{}, fmt.Errorf("decode request: %w", err)
	}
	if strings.TrimSpace(req.Model) == "" {
		return messagesRequest{}, fmt.Errorf("model is required")
	}
	if req.MaxTokens <= 0 {
		req.MaxTokens = 1024
	}
	return req, nil
}

func buildTranscript(system any, messages []messageContent) string {
	parts := make([]string, 0, len(messages)+1)
	if text := contentText(system); text != "" {
		parts = append(parts, "system: "+text)
	}
	for _, message := range messages {
		text := contentText(message.Content)
		if text == "" {
			continue
		}
		role := strings.TrimSpace(message.Role)
		if role == "" {
			role = "user"
		}
		parts = append(parts, role+": "+text)
	}
	return strings.Join(parts, "\n\n")
}

func contentText(content any) string {
	switch value := content.(type) {
	case string:
		return strings.TrimSpace(value)
	case []any:
		var parts []string
		for _, item := range value {
			if text := contentText(item); text != "" {
				parts = append(parts, text)
			}
		}
		return strings.Join(parts, "\n")
	case map[string]any:
		if kind, _ := value["type"].(string); strings.TrimSpace(kind) == "text" {
			if text, _ := value["text"].(string); text != "" {
				return strings.TrimSpace(text)
			}
		}
		return ""
	default:
		return ""
	}
}

func estimateTokens(text string) int {
	const charsPerToken = 4.0
	trimmed := strings.TrimSpace(text)
	if trimmed == "" {
		return 0
	}
	return int(math.Ceil(float64(len(trimmed)) / charsPerToken))
}

func splitChunks(text string, chunkSize int) []string {
	if chunkSize <= 0 || len(text) <= chunkSize {
		return []string{text}
	}

	out := make([]string, 0, (len(text)+chunkSize-1)/chunkSize)
	for len(text) > chunkSize {
		out = append(out, text[:chunkSize])
		text = text[chunkSize:]
	}
	if text != "" {
		out = append(out, text)
	}
	return out
}

func anthropicMessageID() string {
	return "msg_" + time.Now().UTC().Format("20060102150405.000000000")
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]any{
		"type":  "error",
		"error": map[string]any{"type": "invalid_request_error", "message": message},
	})
}

func writeEvent(w http.ResponseWriter, event string, payload any) {
	data, _ := json.Marshal(payload)
	_, _ = io.WriteString(w, "event: "+event+"\n")
	_, _ = io.WriteString(w, "data: ")
	_, _ = w.Write(data)
	_, _ = io.WriteString(w, "\n\n")
}

func flush(flusher http.Flusher) {
	if flusher != nil {
		flusher.Flush()
	}
}
