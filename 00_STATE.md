# 00_STATE.md — pswitch Repository Analysis

## Repository Identity

- **Upstream**: https://github.com/wlynxg/pswitch
- **Fork**: https://github.com/okwn/pswitch (forked 2026-05-22)
- **License**: No license file present (license: null)
- **Archived**: false
- **Language**: Go 1.26
- **Topics**: ai-proxy, anthropic, dashboard, gateway, golang, load-balancer, openai-compatible
- **Default branch**: master
- **Upstream stars**: 2 | **Upstream forks**: 1

## Repository Structure

```
pswitch/
├── cmd/pswitch/          # CLI entry point (main, serve, init)
├── internal/
│   ├── adminui/          # Embedded dashboard UI (Go embed)
│   ├── config/           # TOML/JSON config loading, validation
│   ├── logx/             # Structured logging (uber/zap)
│   ├── metrics/          # Persistent metrics store (JSON)
│   ├── pool/             # Provider pool & load balancing
│   ├── protocol/
│   │   ├── anthropic/    # Anthropic /v1/messages adapter
│   │   └── openai/       # OpenAI-compatible handler
│   ├── runtime/          # Runtime config manager
│   ├── server/           # HTTP router & dispatcher
│   └── upstream/         # HTTP helpers, usage extraction
├── docs/                 # Configuration, usage, troubleshooting docs
├── Dockerfile            # Multi-stage Alpine build
├── Makefile              # build, run, test, init, clean
├── go.mod / go.sum       # Dependencies
└── .github/workflows/
    ├── ci.yml            # go test ./... + docker build
    └── release.yml       # Multi-platform release on tag push
```

## Upstream Activity

- **4 commits** on master: initial commit, defaults/CLI, automated releases+README, Docker publishing workflow
- **0 open issues** (none reported)
- **0 open pull requests** (none open)
- Very early-stage project (~2 days old as of analysis)

## Dependency Inventory

| Dependency | Version | Purpose |
|---|---|---|
| github.com/go-chi/chi/v5 | v5.2.5 | HTTP routing |
| github.com/pelletier/go-toml/v2 | v2.2.4 | TOML config parsing |
| go.uber.org/zap | v1.28.0 | Structured logging |
| github.com/mattn/go-isatty | v0.22 | TTY detection for log coloring |
| github.com/stretchr/testify | v1.8.1 | Test assertions (indirect) |
| go.uber.org/goleak | v1.3.0 | Goroutine leak detection (indirect) |

## Notable Architecture Decisions

1. **No `go install` support** — project uses `make build` + binary output to `./bin/pswitch`; no `go install` target
2. **Config precedence**: `settings.json` (runtime) > TOML config file > built-in defaults
3. **Pool modes**: `sequential`, `round_robin`, `least_failures` — all share the same `Pool` struct with different ordering
4. **Anthropic adapter**: Wraps the OpenAI handler internally, transforms Anthropic `/v1/messages` → OpenAI `/v1/responses`
5. **Metrics persistence**: JSON file written on every `record()` call via atomic rename (temp file)
6. **Admin token**: Optional `PSWITCH_ADMIN_TOKEN` env var; dashboard and API both require it when set
7. **Docker**: Multi-stage build on Alpine 3.22, listens on 8080, stores state in `/data`
8. **GitHub Release automation**: Tags trigger multi-platform builds (linux amd64/arm64, darwin amd64/arm64, windows) + Docker image to ghcr.io

## Issues & Quality Observations

- **License**: Repo has `license: null` — no LICENSE file; author should add one
- **No CONTRIBUTING file**: No contribution guidelines present
- **No Go installed**: `make test` fails because `go` binary is not in PATH; CI uses `actions/setup-go@v5` so local test requires Go 1.26
- **No API docs**: No godoc or generated API documentation
- **Very small codebase**: ~4 commits, 2 contributors, early prototype stage
- **No SECURITY policy**: No SECURITY.md
- **Dashboard is embedded**: Uses Go `//go:embed` for static assets; no separate frontend repo

## CI/CD

- `ci.yml`: Runs `go test ./...`, `go build ./...`, and `docker build -t pswitch:ci .` on every push/PR
- `release.yml`: Builds and uploads archives + checksums + Docker image on version tag push