# 05_PR_CANDIDATES.md — pswitch PR Opportunities

## Quality Gaps Identified

| # | Gap | Impact | Effort | Files |
|---|---|---|---|---|
| 1 | No LICENSE file | Legal risk: users can't legally use the software | Low | root |
| 2 | No CONTRIBUTING.md | Unclear how to contribute | Low | root |
| 3 | No SECURITY.md | No way to report vulnerabilities responsibly | Low | root |
| 4 | No `go install` support | Can't `go install github.com/wlynxg/pswitch@latest` | Low | Makefile, cmd/pswitch/main.go |
| 5 | No request/idle timeouts on upstream HTTP client | Potential resource leaks on slow/unresponsive upstreams | Medium | internal/protocol/openai/handler.go |
| 6 | Hardcoded 5s graceful shutdown | Not configurable for long-running requests | Low | cmd/pswitch/serve.go |
| 7 | No golangci-lint in CI | No linting/static analysis in CI pipeline | Low | .github/workflows/ci.yml |
| 8 | No health check API endpoint | Containers/health checks need `/healthz` or similar | Low | internal/server/router.go, internal/adminui/admin.go |
| 9 | No per-request timeout | Upstream requests can hang indefinitely | Medium | internal/protocol/openai/handler.go |
| 10 | Missing `settings.json` schema validation on load | Corrupted settings.json crashes startup | Low | internal/config/config.go |

---

## Candidate 1: Add LICENSE file (MIT)

**Why**: The GitHub API reports `license: null`. Every open-source project needs a license. MIT is the most common and matches the project's straightforward, lightweight nature.

**Files touched**: `LICENSE` (new file)

**Content**: MIT License text, year 2026, copyright wlynxg.

---

## Candidate 2: Add CONTRIBUTING.md

**Why**: No contribution guidelines exist. Should specify how to run tests, the PR process, code style, and what PRs are welcome.

**Files touched**: `CONTRIBUTING.md` (new file)

**Content**: Dev setup (`make build`, `make test`), PR checklist, testing expectations, go.mod minimum Go version, hot reload note.

---

## Candidate 3: Add SECURITY.md

**Why**: No SECURITY.md. Should provide instructions for reporting vulnerabilities via GitHub Security Advisories or email.

**Files touched**: `SECURITY.md` (new file)

---

## Candidate 4: Add `go install` support

**Why**: Currently only `make build`. Adding `go install` support allows `go install github.com/wlynxg/pswitch/cmd/pswitch@latest`. This is a standard Go idiom and widely expected for Go projects.

**Files touched**: `Makefile` — add `install` target, or alternatively add `cmd/pswitch/install.go` if the project wants a `pswitch install` subcommand.

**Implementation**: Add to Makefile:
```makefile
install:
    go install ./cmd/pswitch
```

Or add `Install` function to `cmd/pswitch/main.go` as a subcommand.

---

## Candidate 5: Add configurable request timeouts for upstream calls

**Why**: The OpenAI handler creates an unconfigured `http.Client{}` (no timeouts). If an upstream provider is unresponsive, requests can hang indefinitely. A timeout should be configurable.

**Files touched**: `internal/protocol/openai/handler.go`, `internal/config/config.go`

**Implementation**: Add `request_timeout` field to config (default 30s), wire into the OpenAI handler's HTTP client. For streaming responses, the timeout needs careful handling (don't cancel mid-stream on timeout, only on connection establishment).

**Risk**: Medium — changing HTTP client behavior could break streaming. But the change is additive with a safe default.

---

## Candidate 6: Add health check endpoint

**Why**: Container orchestrators (Kubernetes, Docker Compose) need a `/healthz` or `/ready` endpoint to determine if the service is alive. Currently no such endpoint exists.

**Files touched**: `internal/server/router.go`, potentially `internal/adminui/admin.go`

**Implementation**: Add `GET /healthz` → 200 OK (no auth required) to the dispatcher. Can also add `GET /ready` which checks if at least one provider is healthy.

---

## Candidate 7: Add golangci-lint to CI

**Why**: The CI currently only runs `go test ./...` and `go build ./...`. Adding golangci-lint would catch common Go bugs, style issues, and static analysis problems pre-commit.

**Files touched**: `.github/workflows/ci.yml`

**Implementation**: Add a `lint` job using `golangci-lint/golangci-lint-action@v6`. Use reasonable defaults (enable `errcheck`, `staticcheck`, `govet`, etc.).

---

## Candidate 8: Add graceful shutdown timeout configuration

**Why**: The 5-second graceful shutdown timeout in `serve.go` is hardcoded. In some deployments, requests take longer. Should be configurable via CLI flag and config.

**Files touched**: `cmd/pswitch/serve.go`, `internal/config/config.go`, `cmd/pswitch/main.go`

---

## Candidate 9: Improve error messages on all-provider-failure

**Why**: When all providers fail, the error returned is `"all providers failed"` which is not descriptive. The actual underlying error (timeout, 401, 429, 500) is lost.

**Files touched**: `internal/protocol/openai/handler.go`

**Implementation**: Collect all errors and surface them in the final error message, e.g. `"all 3 providers failed: provider-a: connection refused; provider-b: 401 unauthorized; provider-c: 500 internal server error"`.

---

## Candidate 10: Validate settings.json on load, provide recovery

**Why**: If `settings.json` is corrupted, the application may fail to start with a cryptic error. Should validate the JSON and fall back to defaults or the TOML config if JSON is corrupt.

**Files touched**: `internal/config/config.go`, `cmd/pswitch/serve.go`

**Implementation**: Wrap JSON load errors in `loadStartupConfig`; on JSON parse error, log a warning and fall back to TOML/default rather than failing entirely.