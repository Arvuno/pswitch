# 01_REPO_MAP.md — pswitch Codebase Map

## Package Hierarchy

```
pswitch (module root)
└── cmd/pswitch/
    ├── main.go           # Entry point: "init" vs "serve" dispatch
    ├── serve.go          # serve command, flag parsing, HTTP server setup, health loop
    ├── init.go           # "pswitch init" command
    ├── serve_test.go     # Tests for serve flag parsing and config loading
    ├── docker_support_test.go
    └── release_docs_test.go

internal/
├── adminui/
    ├── admin.go          # Dashboard HTTP handler (chi router, API endpoints, static file serving)
    ├── admin_test.go
    └── assets/          # Embedded static assets (index.html, app.js, style.css, icon.svg)

├── config/
    ├── config.go         # Config struct, Load/LoadJSON/Write/WriteJSON, Validate, defaults
    ├── config_test.go
    └── (raw types for TOML unmarshaling)

├── logx/
    ├── logx.go          # zap logger wrapper (Infof, Warnf, Debugf, Init with color toggle)
    └── logx_test.go

├── metrics/
    ├── store.go         # Persistent metrics store (JSON), Snapshot(), RecordSuccess/RecordFailure
    └── store_test.go

├── pool/
    ├── pool.go          # Pool struct, Candidates(), MarkSuccess/MarkFailure, ProbeDue, Mode types
    └── pool_test.go

├── protocol/
    ├── anthropic/
    │   ├── handler.go   # Anthropic /v1/messages → OpenAI /v1/responses adapter, streaming support
    │   └── handler_test.go
    └── openai/
        ├── handler.go   # OpenAI handler, failover loop, streaming, usage extraction
        └── handler_test.go

├── runtime/
    ├── runtime.go       # Manager struct: Config(), Pool(), Metrics(), UpdateConfig(), ProviderStatuses()
    └── runtime_test.go

├── server/
    ├── router.go        # NewRouter (chi), dispatcher ServeHTTP, matchRoute, pathMatches
    └── router_test.go

└── upstream/
    ├── http.go         # CopyHeaders, ReadRequestBody, CaptureResponse, CopyResponseBody, ShouldFailover, JoinPaths
    └── usage.go        # ExtractUsage, ExtractRequestedModel, StreamUsageCollector
```

## Key Types and Interfaces

### config.Config
```go
type Config struct {
    Listen, Mode string
    FailureThreshold int
    Cooldown, HealthCheckInterval, HealthCheckTimeout time.Duration
    Routes []Route
    Providers []Provider
}
```

### pool.Pool
```go
type Pool struct {
    mu sync.RWMutex
    mode Mode
    failureThreshold int
    cooldown time.Duration
    providers []*providerState
    indexByName map[string]int
    cursor int
}
// Modes: ModeSequential, ModeRoundRobin, ModeLeastFailures
```

### metrics.Store
```go
type Store struct {
    mu sync.RWMutex
    path string
    data storeFile  // persisted JSON
}
// RecordSuccess(provider, usage, model, now), RecordFailure(provider, model, now)
// Snapshot(now) → Snapshot struct with Overview, Windows, Providers, Models, Series
```

### runtime.Manager
```go
type Manager struct {
    mu sync.RWMutex
    settingsPath string
    cfg config.Config
    pool *pool.Pool
    metrics *metrics.Store
}
// Config(), Pool(), Snapshot(), Metrics(), MetricsSnapshot(now), UpdateConfig(), ProviderStatuses()
```

### openai.Handler
```go
type Handler struct {
    pool *pool.Pool
    client *http.Client
    mode pool.Mode
    metrics *metrics.Store
}
// ServeHTTP(w, r) — failover loop over pool.Candidates(), forward, record metrics
```

### anthropic.Handler
```go
type Handler struct {
    model, upstreamModel string
    upstream http.Handler  // delegates to openai.Handler
}
// ServeHTTP dispatches: /v1/models, /v1/models/{id}, /v1/messages, /v1/messages/count_tokens
```

## File Statistics

| Path | Lines | Purpose |
|---|---|---|
| internal/protocol/anthropic/handler.go | 408 | Anthropic adapter |
| internal/adminui/admin.go | 389 | Dashboard API |
| internal/metrics/store.go | 462 | Metrics persistence |
| internal/config/config.go | 457 | Config loading/validation |
| internal/protocol/openai/handler.go | 207 | OpenAI proxy handler |
| internal/server/router.go | 93 | HTTP routing & dispatch |
| cmd/pswitch/serve.go | 330 | CLI serve command |
| internal/pool/pool.go | 265 | Load balancer pool |
| internal/runtime/runtime.go | 135 | Runtime manager |
| internal/upstream/http.go | 116 | HTTP utilities |
| internal/upstream/usage.go | 87 | Usage extraction |
| internal/logx/logx.go | ~50 | Logging wrapper |
| docs/ | ~300 | Documentation |

## Test Files

All packages have corresponding `*_test.go` files:
- `cmd/pswitch/serve_test.go`
- `cmd/pswitch/init.go`
- `cmd/pswitch/release_docs_test.go`
- `cmd/pswitch/docker_support_test.go`
- `internal/adminui/admin_test.go`
- `internal/config/config_test.go`
- `internal/logx/logx_test.go`
- `internal/metrics/store_test.go`
- `internal/pool/pool_test.go`
- `internal/protocol/anthropic/handler_test.go`
- `internal/protocol/openai/handler_test.go`
- `internal/runtime/runtime_test.go`
- `internal/server/router_test.go`

## Entry Point Flow

```
pswitch [args]
  └─ main.main()
      ├─ "init" → runInit() in init.go
      └─ else → runServe() in serve.go
             ├─ parseServeArgs() → serveArgs
             ├─ loadStartupConfig() → config.Config (settings.json > TOML > default)
             ├─ applyServeOverrides() → CLI flags
             ├─ cfg.Validate()
             ├─ logx.Init()
             ├─ pruntime.New() → Manager
             ├─ server.NewRouter() → http.Handler
             ├─ net.Listen() on cfg.Listen
             ├─ healthLoop() goroutine → periodic provider probing
             └─ server.Serve()
```

## Configuration Flow

```
settings.json (runtime, JSON)
  OR config.toml (startup, TOML)
  OR built-in defaults
        ↓
  config.LoadJSON / config.Load → Config
        ↓
  cfg.Validate() → error if invalid
        ↓
  runtime.Manager (immutable after creation, replaced on UpdateConfig)
        ↓
  pool.New(cfg) → Pool
        ↓
  pool.Pool → used by openai.Handler and server.dispatcher
```

## Request Routing Flow

```
HTTP Request
    ↓
dispatcher.ServeHTTP(w, r)
    ├─ manager.Snapshot() → config.Config + *Pool
    ├─ matchRoute(cfg.Routes, r.URL.Path) → matched Route
    ├─ openai.NewHandler(pool, options) → *openai.Handler
    ├─ if anthropic route: anthropic.NewHandler(wrapping openai.Handler)
    ├─ if prefix != "/": http.StripPrefix(handler)
    └─ handler.ServeHTTP(w, r)
            ├─ read request body
            ├─ pool.Candidates(mode, now) → healthy providers
            ├─ for each candidate: forward() upstream
            ├─ if failover status: MarkFailure, continue to next
            ├─ if success: MarkSuccess, record usage
            └─ write response to client
```

## Health Check Loop

```
healthLoop(ctx, manager)
    ├─ timer = cfg.HealthCheckInterval
    └─ select:
         ├─ ctx.Done() → return
         └─ timer.C → now
                ├─ providerPool.ProbeDue(ctx, client, now)
                │     → marks unhealthy providers as due for probe
                │     → HTTP GET /v1/models with Bearer token
                │     → MarkSuccess (if 2xx) or MarkFailure
                └─ for each "recovered" event: logx.Infof
```

## Data Persistence

| File | Format | Written By | Read By |
|---|---|---|---|
| settings.json | JSON (TOML-style) | runtime.Manager.UpdateConfig() | loadStartupConfig() |
| metrics.json | JSON | metrics.Store.record() (on every request) | metrics.Store.load() |

Both use atomic rename via `os.CreateTemp` + `os.Rename` for safe writes.