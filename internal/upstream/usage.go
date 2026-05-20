package upstream

import (
	"bytes"
	"encoding/json"
	"strings"
)

type UsageSummary struct {
	InputTokens  int64 `json:"input_tokens"`
	OutputTokens int64 `json:"output_tokens"`
	TotalTokens  int64 `json:"total_tokens"`
}

type responseEnvelope struct {
	Usage    *UsageSummary `json:"usage"`
	Response *struct {
		Usage *UsageSummary `json:"usage"`
	} `json:"response"`
}

type requestEnvelope struct {
	Model string `json:"model"`
}

func ExtractUsage(body []byte) (UsageSummary, bool) {
	var envelope responseEnvelope
	if err := json.Unmarshal(body, &envelope); err != nil {
		return UsageSummary{}, false
	}
	if envelope.Usage == nil {
		if envelope.Response == nil || envelope.Response.Usage == nil {
			return UsageSummary{}, false
		}
		return *envelope.Response.Usage, true
	}
	return *envelope.Usage, true
}

func ExtractRequestedModel(body []byte) string {
	var envelope requestEnvelope
	if err := json.Unmarshal(body, &envelope); err != nil {
		return ""
	}
	return strings.TrimSpace(envelope.Model)
}

type StreamUsageCollector struct {
	pending []byte
	usage   *UsageSummary
}

func (c *StreamUsageCollector) Write(p []byte) (int, error) {
	c.pending = append(c.pending, p...)
	for {
		idx := bytes.IndexByte(c.pending, '\n')
		if idx == -1 {
			break
		}
		line := strings.TrimRight(string(c.pending[:idx]), "\r")
		c.pending = c.pending[idx+1:]
		c.consumeLine(line)
	}
	return len(p), nil
}

func (c *StreamUsageCollector) Usage() (UsageSummary, bool) {
	if c.usage == nil {
		return UsageSummary{}, false
	}
	return *c.usage, true
}

func (c *StreamUsageCollector) consumeLine(line string) {
	if !strings.HasPrefix(line, "data:") {
		return
	}

	payload := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
	if payload == "" || payload == "[DONE]" {
		return
	}

	if usage, ok := ExtractUsage([]byte(payload)); ok {
		c.usage = &usage
	}
}
