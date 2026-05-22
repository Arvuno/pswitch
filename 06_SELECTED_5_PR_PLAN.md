# 06_SELECTED_5_PR_PLAN.md — Top 5 PR Recommendations for pswitch

## Rationale for Selection

These 5 PRs are selected based on: **low effort, high impact, minimal risk, and alignment with the project's maturity stage** (very early, ~2 days old, needs basic OSS hygiene).

---

## PR #1: Add LICENSE (MIT) + CONTRIBUTING.md + SECURITY.md

**Priority**: P0 — these are blocking for any serious OSS adoption.

### Changes

Create 3 new files at repo root:

**`LICENSE`** (MIT):
```
MIT License

Copyright (c) 2026 wlynxg

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

**`CONTRIBUTING.md`**:
```
# Contributing to pswitch

## Development Setup

```bash
# Build
make build        # outputs ./bin/pswitch

# Test
make test         # runs go test ./...

# Run locally
make run          # build + run with CONFIG=./config.toml
```

## Code Style

- Run `go fmt ./...` before committing
- All packages should have tests
- Use structured logging (from internal/logx) — no direct fmt.Print*

## Submitting Changes

1. Fork the repository
2. Create a feature branch from `master`
3. Make your changes with tests
4. Ensure `make test` passes
5. Open a pull request against `master`

## What to Contribute

- Bug fixes with tests
- Improved documentation
- Additional protocol handlers (e.g., Google AI Studio, Azure OpenAI)
- Performance improvements
- Additional load-balancing modes
```

**`SECURITY.md`**:
```
# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| 0.1.x   | :white_check_mark: |

## Reporting a Vulnerability

If you discover a security vulnerability, please report it via GitHub Security Advisories.

Do NOT open a public issue for security vulnerabilities.

Please include:
- Description of the vulnerability
- Steps to reproduce
- Potential impact
- Any suggested fixes (optional)
```

### Verification

After PR:
```bash
gh api repos/wlynxg/pswitch --jq '.license.spdx_id'  # should return MIT
```

---

## PR #2: Add `go install` support

**Priority**: P1 — standard Go project hygiene.

### Changes

**`Makefile`** — add `install` target:

```makefile
install:
    go install ./cmd/pswitch
```

No changes to Go code needed. The module name is `pswitch` and the main package is at `cmd/pswitch`, which is the standard layout for `go install`.

**Verification**:
```bash
go install ./cmd/pswitch && which pswitch  # should find installed binary
```

---

## PR #3: Add health check endpoint

**Priority**: P1 — needed for container orchestration deployments.

### Changes

**`internal/server/router.go`** — add a health check handler:

```go
// In NewRouter, before the "/*" catchall:
router.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    _, _ = w.Write([]byte("OK"))
})
```

Optionally also add `GET /ready` that checks provider health, but `/healthz` is the minimum for Kubernetes liveness/readiness probes.

**No auth required** for health endpoints.

**Verification**:
```bash
# After building and running:
curl http://localhost:8080/healthz  # should return 200 OK
```

---

## PR #4: Add golangci-lint to CI

**Priority**: P1 — catches bugs pre-commit, improves code quality.

### Changes

**`.github/workflows/ci.yml`** — modify to:

```yaml
jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: "1.26.x"
      - uses: golangci-lint/golangci-lint-action@v6
        with:
          version: latest

  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: "1.26.x"
      - run: go test ./...
      - run: go build ./...
      - name: Docker image
        run: docker build -t pswitch:ci .
```

Also add a `.golangci.yml` at repo root with sensible defaults:

```yaml
linters:
  enable:
    - errcheck
    - staticcheck
    - govet
    - gofmt
    - goimports
    - misspell
    - unconvert
linters-settings:
  errcheck:
    check-type-assertions: true
    check-blank: true
```

**Verification**:
```bash
golangci-lint run ./...  # should pass with no errors
```

---

## PR #5: Add configurable request timeout for upstream HTTP calls

**Priority**: P2 — improves resilience, prevents hanging requests.

### Changes

**`internal/config/config.go`**:

Add `RequestTimeout time.Duration` to `Config` struct with default 30s.

**`internal/protocol/openai/handler.go`**:

```go
// In NewHandler, use timeout client if configured:
client := opts.Client
if client == nil {
    client = &http.Client{
        Timeout: opts.RequestTimeout,  // from config
    }
}
```

**`cmd/pswitch/serve.go`**:

Wire `request_timeout` CLI flag (or reuse `--cooldown` approach with optionalDuration), and pass it through to the handler via `openai.Options`.

**`docs/config.md`**:

Document the new `request_timeout` field.

**Verification**:
```bash
# Test that timeout works:
# - configure a provider that delays responses
# - set request_timeout=1s
# - verify client gets a timeout error rather than hanging
```

---

## Implementation Order

1. **PR #1** (LICENSE/CONTRIBUTING/SECURITY) — independent, no code changes
2. **PR #2** (go install) — one-line Makefile change, independent
3. **PR #3** (healthz) — small code change, independent
4. **PR #4** (golangci-lint) — CI-only change, may surface issues to fix first
5. **PR #5** (request timeout) — touches multiple files, requires testing

---

## Risk Assessment

| PR | Risk | Reasoning |
|---|---|---|
| #1 LICENSE | None | New files only |
| #2 go install | None | Makefile only |
| #3 healthz | None | Additive endpoint, no auth changes |
| #4 golangci-lint | Low | May fail CI if linter finds issues; run locally first |
| #5 timeouts | Medium | HTTP client changes; use safe default (30s), test streaming |