<div align="center">

# pswitch

### A local multi-provider proxy for Codex-style and Anthropic-style clients, with failover, health recovery, and a built-in admin dashboard

[![Version](https://img.shields.io/github/v/release/wlynxg/pswitch?color=blue&label=version)](https://github.com/wlynxg/pswitch/releases)
[![Platform](https://img.shields.io/badge/platform-Windows%20%7C%20macOS%20%7C%20Linux-lightgrey.svg)](https://github.com/wlynxg/pswitch/releases)
[![Built with Go](https://img.shields.io/badge/built%20with-Go%201.26-00ADD8.svg)](https://go.dev/)
[![Downloads](https://img.shields.io/github/downloads/wlynxg/pswitch/total)](https://github.com/wlynxg/pswitch/releases/latest)

English | [中文](README_ZH.md)

</div>

## Overview

`pswitch` is a lightweight local proxy for routing AI client traffic across multiple upstream providers.

It is designed for setups where you want:

- one stable local endpoint such as `/codex`
- multiple upstream providers behind it
- automatic failover and recovery
- a clean dashboard for traffic, token usage, provider health, and runtime config

By default, `pswitch` starts with a single OpenAI-compatible route on `/codex`. You can add more routes and providers later from the dashboard or config files.

## Screenshot

![pswitch dashboard](docs/assets/dashboard.png)

## Features

- Multiple upstream providers with automatic failover
- Circuit breaking and periodic health recovery probes
- Three routing modes:
  - `round_robin`
  - `sequential`
  - `least_failures`
- OpenAI-compatible routing out of the box
- Optional Anthropic-compatible route adapter
- Persistent dashboard metrics for:
  - requests
  - input / output / total tokens
  - provider failures
  - per-model usage
- Embedded admin dashboard at `/dashboard/`
- Runtime config editing with hot reload where possible
- `settings.json` and `metrics.json` persisted in the current working directory
- Optional admin token protection with `PSWITCH_ADMIN_TOKEN`

## Quick Start

### Build

```bash
make build
```

Binary output:

```bash
./bin/pswitch
```

### Run

```bash
./bin/pswitch
```

Or:

```bash
./bin/pswitch --config ./config.toml
```

Then open:

```text
http://127.0.0.1:8080/dashboard/
```

## Default Behavior

If no config file exists, `pswitch` starts with the built-in default config.

Default startup behavior:

- listen on `0.0.0.0:8080`
- use `round_robin` mode
- expose one route: `/codex`
- start with no preconfigured providers

Default file behavior:

- user config file is startup input only and is never modified by the program
- dashboard-saved runtime config goes to `./settings.json`
- dashboard metrics go to `./metrics.json`
- if `settings.json` exists, it takes precedence on startup

## Example Usage

### Codex / OpenAI-style client

Point your client to:

```text
http://127.0.0.1:8080/codex
```

Example Codex-style config:

```toml
[model_providers.OpenAI]
base_url = "http://127.0.0.1:8080/codex"
wire_api = "responses"
requires_openai_auth = true
```

### Anthropic-style client

If you manually add an Anthropic route, you can point a Claude-style client to it:

```bash
export ANTHROPIC_BASE_URL=http://127.0.0.1:8080/claude
export ANTHROPIC_API_KEY=dummy
```

## Config Example

```toml
listen = "127.0.0.1:8080"
mode = "least_failures"
failure_threshold = 1
cooldown = "20s"
health_check_interval = "15s"
health_check_timeout = "3s"

[[routes]]
prefix = "/codex"
type = "openai"

[[routes]]
prefix = "/claude"
type = "anthropic"
model = "claude-sonnet-4-20250514"
upstream_model = "gpt-5.4"

[[providers]]
name = "provider-a"
base_url = "https://provider-a.example/v1"
api_key = "sk-your-provider-a-key"

[[providers]]
name = "provider-b"
base_url = "https://provider-b.example/v1"
api_key = "sk-your-provider-b-key"
```

## Admin Dashboard

The embedded dashboard is available at:

```text
/dashboard/
```

It provides:

- overview metrics
- 24h / 7d token windows
- provider analytics
- provider health cards
- per-model usage panels
- runtime config editing
- English / Chinese language switch

If `PSWITCH_ADMIN_TOKEN` is set, both the dashboard UI and admin API require it.

## CLI

Run directly:

```bash
pswitch [--config PATH] [--listen ADDR] [--mode sequential|round_robin|least_failures] [--failure-threshold N] [--cooldown DURATION] [--health-check-interval DURATION] [--health-check-timeout DURATION] [--log-color[=true|false]]
```

Generate a starter config:

```bash
./bin/pswitch init
```

## Documentation

- [Configuration](docs/config.md)
- [Usage](docs/usage.md)
- [Logging](docs/logging.md)
- [Troubleshooting](docs/troubleshooting.md)
- [Development](docs/development.md)

## Makefile

- `make build` builds `./bin/pswitch`
- `make run` starts the service
- `make test` runs the test suite
- `make init` generates an example config
- `make clean` removes build artifacts
