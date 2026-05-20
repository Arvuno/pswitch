package metrics

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

const (
	retentionWindow = 7 * 24 * time.Hour
	hourLayout      = "2006-01-02T15:00:00Z07:00"
	dayLayout       = "2006-01-02"
)

type Usage struct {
	RequestCount int64
	FailureCount int64
	InputTokens  int64
	OutputTokens int64
	TotalTokens  int64
}

type WindowSnapshot struct {
	RequestCount int64 `json:"request_count"`
	FailureCount int64 `json:"failure_count"`
	InputTokens  int64 `json:"input_tokens"`
	OutputTokens int64 `json:"output_tokens"`
	TotalTokens  int64 `json:"total_tokens"`
}

type ProviderSnapshot struct {
	Name         string         `json:"name"`
	Totals       WindowSnapshot `json:"totals"`
	Last24h      WindowSnapshot `json:"last_24h"`
	Last7d       WindowSnapshot `json:"last_7d"`
	TotalTokens  int64          `json:"total_tokens"`
	RequestCount int64          `json:"request_count"`
	FailureCount int64          `json:"failure_count"`
	InputTokens  int64          `json:"input_tokens"`
	OutputTokens int64          `json:"output_tokens"`
}

type ModelSnapshot struct {
	Name         string         `json:"name"`
	Totals       WindowSnapshot `json:"totals"`
	Last24h      WindowSnapshot `json:"last_24h"`
	Last7d       WindowSnapshot `json:"last_7d"`
	TotalTokens  int64          `json:"total_tokens"`
	RequestCount int64          `json:"request_count"`
	FailureCount int64          `json:"failure_count"`
	InputTokens  int64          `json:"input_tokens"`
	OutputTokens int64          `json:"output_tokens"`
}

type OverviewSnapshot struct {
	TotalRequests     int64 `json:"total_requests"`
	TotalFailures     int64 `json:"total_failures"`
	TotalInputTokens  int64 `json:"total_input_tokens"`
	TotalOutputTokens int64 `json:"total_output_tokens"`
	TotalTokens       int64 `json:"total_tokens"`
}

type WindowsSnapshot struct {
	Last24h WindowSnapshot `json:"last_24h"`
	Last7d  WindowSnapshot `json:"last_7d"`
}

type SeriesPoint struct {
	Label        string `json:"label"`
	RequestCount int64  `json:"request_count"`
	FailureCount int64  `json:"failure_count"`
	TotalTokens  int64  `json:"total_tokens"`
}

type SeriesSnapshot struct {
	Hourly24h []SeriesPoint `json:"hourly_24h"`
	Daily7d   []SeriesPoint `json:"daily_7d"`
}

type Snapshot struct {
	Overview  OverviewSnapshot            `json:"overview"`
	Windows   WindowsSnapshot             `json:"windows"`
	Providers map[string]ProviderSnapshot `json:"providers"`
	Models    map[string]ModelSnapshot    `json:"models"`
	Series    SeriesSnapshot              `json:"series"`
}

func (s Snapshot) Provider(name string) ProviderSnapshot {
	return s.Providers[name]
}

func (s Snapshot) Model(name string) ModelSnapshot {
	return s.Models[name]
}

type storeFile struct {
	Overview   usageRecord                       `json:"overview"`
	Providers  map[string]usageRecord            `json:"providers"`
	Models     map[string]usageRecord            `json:"models"`
	Hourly     map[string]usageRecord            `json:"hourly"`
	ByProvider map[string]map[string]usageRecord `json:"by_provider"`
	ByModel    map[string]map[string]usageRecord `json:"by_model"`
}

type usageRecord struct {
	RequestCount int64 `json:"request_count"`
	FailureCount int64 `json:"failure_count"`
	InputTokens  int64 `json:"input_tokens"`
	OutputTokens int64 `json:"output_tokens"`
	TotalTokens  int64 `json:"total_tokens"`
}

type Store struct {
	mu   sync.RWMutex
	path string
	data storeFile
}

func Open(path string) (*Store, error) {
	store := &Store{
		path: path,
		data: storeFile{
			Providers:  make(map[string]usageRecord),
			Models:     make(map[string]usageRecord),
			Hourly:     make(map[string]usageRecord),
			ByProvider: make(map[string]map[string]usageRecord),
			ByModel:    make(map[string]map[string]usageRecord),
		},
	}

	if err := store.load(); err != nil {
		return nil, err
	}
	return store, nil
}

func (s *Store) RecordSuccess(provider string, usage Usage, model string, now time.Time) {
	usage.FailureCount = 0
	s.record(provider, model, usage, now)
}

func (s *Store) RecordFailure(provider, model string, now time.Time) {
	s.record(provider, model, Usage{FailureCount: 1}, now)
}

func (s *Store) Snapshot(now time.Time) Snapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()

	last24Cutoff := truncateHour(now).Add(-23 * time.Hour)
	last7Cutoff := truncateHour(now).Add(-167 * time.Hour)

	snapshot := Snapshot{
		Providers: make(map[string]ProviderSnapshot, len(s.data.Providers)),
		Models:    make(map[string]ModelSnapshot, len(s.data.Models)),
		Series: SeriesSnapshot{
			Hourly24h: buildHourlySeries(s.data.Hourly, truncateHour(now)),
			Daily7d:   buildDailySeries(s.data.Hourly, truncateHour(now)),
		},
	}
	snapshot.Overview = toOverview(s.data.Overview)

	for hour, record := range s.data.Hourly {
		bucketTime, err := time.Parse(hourLayout, hour)
		if err != nil {
			continue
		}
		if !bucketTime.Before(last24Cutoff) {
			snapshot.Windows.Last24h = addWindow(snapshot.Windows.Last24h, record)
		}
		if !bucketTime.Before(last7Cutoff) {
			snapshot.Windows.Last7d = addWindow(snapshot.Windows.Last7d, record)
		}
	}

	for provider, totals := range s.data.Providers {
		item := ProviderSnapshot{
			Name:         provider,
			Totals:       toWindow(totals),
			Last24h:      WindowSnapshot{},
			Last7d:       WindowSnapshot{},
			TotalTokens:  totals.TotalTokens,
			RequestCount: totals.RequestCount,
			FailureCount: totals.FailureCount,
			InputTokens:  totals.InputTokens,
			OutputTokens: totals.OutputTokens,
		}

		for hour, record := range s.data.ByProvider[provider] {
			bucketTime, err := time.Parse(hourLayout, hour)
			if err != nil {
				continue
			}
			if !bucketTime.Before(last24Cutoff) {
				item.Last24h = addWindow(item.Last24h, record)
			}
			if !bucketTime.Before(last7Cutoff) {
				item.Last7d = addWindow(item.Last7d, record)
			}
		}
		snapshot.Providers[provider] = item
	}

	for model, totals := range s.data.Models {
		item := ModelSnapshot{
			Name:         model,
			Totals:       toWindow(totals),
			Last24h:      WindowSnapshot{},
			Last7d:       WindowSnapshot{},
			TotalTokens:  totals.TotalTokens,
			RequestCount: totals.RequestCount,
			FailureCount: totals.FailureCount,
			InputTokens:  totals.InputTokens,
			OutputTokens: totals.OutputTokens,
		}

		for hour, record := range s.data.ByModel[model] {
			bucketTime, err := time.Parse(hourLayout, hour)
			if err != nil {
				continue
			}
			if !bucketTime.Before(last24Cutoff) {
				item.Last24h = addWindow(item.Last24h, record)
			}
			if !bucketTime.Before(last7Cutoff) {
				item.Last7d = addWindow(item.Last7d, record)
			}
		}
		snapshot.Models[model] = item
	}

	return snapshot
}

func (s *Store) record(provider, model string, usage Usage, now time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.pruneLocked(now)

	bucket := truncateHour(now).Format(hourLayout)
	s.data.Overview = addRecord(s.data.Overview, toRecord(usage))
	s.data.Providers[provider] = addRecord(s.data.Providers[provider], toRecord(usage))
	if model != "" {
		s.data.Models[model] = addRecord(s.data.Models[model], toRecord(usage))
	}
	s.data.Hourly[bucket] = addRecord(s.data.Hourly[bucket], toRecord(usage))

	if s.data.ByProvider[provider] == nil {
		s.data.ByProvider[provider] = make(map[string]usageRecord)
	}
	s.data.ByProvider[provider][bucket] = addRecord(s.data.ByProvider[provider][bucket], toRecord(usage))
	if model != "" {
		if s.data.ByModel[model] == nil {
			s.data.ByModel[model] = make(map[string]usageRecord)
		}
		s.data.ByModel[model][bucket] = addRecord(s.data.ByModel[model][bucket], toRecord(usage))
	}

	if err := s.persistLocked(); err != nil {
		return
	}
}

func (s *Store) load() error {
	data, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("read metrics: %w", err)
	}

	var stored storeFile
	if err := json.Unmarshal(data, &stored); err != nil {
		return nil
	}
	if stored.Providers == nil {
		stored.Providers = make(map[string]usageRecord)
	}
	if stored.Models == nil {
		stored.Models = make(map[string]usageRecord)
	}
	if stored.Hourly == nil {
		stored.Hourly = make(map[string]usageRecord)
	}
	if stored.ByProvider == nil {
		stored.ByProvider = make(map[string]map[string]usageRecord)
	}
	if stored.ByModel == nil {
		stored.ByModel = make(map[string]map[string]usageRecord)
	}
	s.data = stored
	return nil
}

func (s *Store) pruneLocked(now time.Time) {
	cutoff := truncateHour(now).Add(-retentionWindow)
	for hour := range s.data.Hourly {
		bucketTime, err := time.Parse(hourLayout, hour)
		if err != nil || bucketTime.Before(cutoff) {
			delete(s.data.Hourly, hour)
		}
	}
	for provider, hourly := range s.data.ByProvider {
		for hour := range hourly {
			bucketTime, err := time.Parse(hourLayout, hour)
			if err != nil || bucketTime.Before(cutoff) {
				delete(hourly, hour)
			}
		}
		if len(hourly) == 0 {
			delete(s.data.ByProvider, provider)
		}
	}
	for model, hourly := range s.data.ByModel {
		for hour := range hourly {
			bucketTime, err := time.Parse(hourLayout, hour)
			if err != nil || bucketTime.Before(cutoff) {
				delete(hourly, hour)
			}
		}
		if len(hourly) == 0 {
			delete(s.data.ByModel, model)
		}
	}
}

func (s *Store) persistLocked() error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o700); err != nil {
		return fmt.Errorf("create metrics dir: %w", err)
	}
	data, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return fmt.Errorf("encode metrics: %w", err)
	}
	tempFile, err := os.CreateTemp(filepath.Dir(s.path), ".pswitch-metrics-*.json")
	if err != nil {
		return fmt.Errorf("create temp metrics: %w", err)
	}
	tempPath := tempFile.Name()
	defer func() {
		_ = os.Remove(tempPath)
	}()
	if _, err := tempFile.Write(data); err != nil {
		_ = tempFile.Close()
		return fmt.Errorf("write temp metrics: %w", err)
	}
	if err := tempFile.Chmod(0o600); err != nil {
		_ = tempFile.Close()
		return fmt.Errorf("chmod temp metrics: %w", err)
	}
	if err := tempFile.Close(); err != nil {
		return fmt.Errorf("close temp metrics: %w", err)
	}
	if err := os.Rename(tempPath, s.path); err != nil {
		return fmt.Errorf("replace metrics: %w", err)
	}
	return nil
}

func toRecord(usage Usage) usageRecord {
	return usageRecord{
		RequestCount: usage.RequestCount,
		FailureCount: usage.FailureCount,
		InputTokens:  usage.InputTokens,
		OutputTokens: usage.OutputTokens,
		TotalTokens:  usage.TotalTokens,
	}
}

func addRecord(base, delta usageRecord) usageRecord {
	base.RequestCount += delta.RequestCount
	base.FailureCount += delta.FailureCount
	base.InputTokens += delta.InputTokens
	base.OutputTokens += delta.OutputTokens
	base.TotalTokens += delta.TotalTokens
	return base
}

func toOverview(record usageRecord) OverviewSnapshot {
	return OverviewSnapshot{
		TotalRequests:     record.RequestCount,
		TotalFailures:     record.FailureCount,
		TotalInputTokens:  record.InputTokens,
		TotalOutputTokens: record.OutputTokens,
		TotalTokens:       record.TotalTokens,
	}
}

func toWindow(record usageRecord) WindowSnapshot {
	return WindowSnapshot{
		RequestCount: record.RequestCount,
		FailureCount: record.FailureCount,
		InputTokens:  record.InputTokens,
		OutputTokens: record.OutputTokens,
		TotalTokens:  record.TotalTokens,
	}
}

func addWindow(base WindowSnapshot, delta usageRecord) WindowSnapshot {
	base.RequestCount += delta.RequestCount
	base.FailureCount += delta.FailureCount
	base.InputTokens += delta.InputTokens
	base.OutputTokens += delta.OutputTokens
	base.TotalTokens += delta.TotalTokens
	return base
}

func truncateHour(t time.Time) time.Time {
	return t.UTC().Truncate(time.Hour)
}

func buildHourlySeries(hourly map[string]usageRecord, end time.Time) []SeriesPoint {
	out := make([]SeriesPoint, 0, 24)
	start := end.Add(-23 * time.Hour)
	for i := 0; i < 24; i++ {
		ts := start.Add(time.Duration(i) * time.Hour)
		record := hourly[ts.Format(hourLayout)]
		out = append(out, SeriesPoint{
			Label:        ts.Format("15:04"),
			RequestCount: record.RequestCount,
			FailureCount: record.FailureCount,
			TotalTokens:  record.TotalTokens,
		})
	}
	return out
}

func buildDailySeries(hourly map[string]usageRecord, end time.Time) []SeriesPoint {
	days := make(map[string]usageRecord, 7)
	start := end.Add(-6 * 24 * time.Hour)
	for hour, record := range hourly {
		ts, err := time.Parse(hourLayout, hour)
		if err != nil || ts.Before(start) || ts.After(end) {
			continue
		}
		key := ts.Format(dayLayout)
		days[key] = addRecord(days[key], record)
	}

	out := make([]SeriesPoint, 0, 7)
	for i := 0; i < 7; i++ {
		day := end.Add(time.Duration(i-6) * 24 * time.Hour)
		record := days[day.Format(dayLayout)]
		out = append(out, SeriesPoint{
			Label:        day.Format("01-02"),
			RequestCount: record.RequestCount,
			FailureCount: record.FailureCount,
			TotalTokens:  record.TotalTokens,
		})
	}

	sort.SliceStable(out, func(i, j int) bool {
		return out[i].Label < out[j].Label
	})
	return out
}
