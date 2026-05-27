# Usage

## Build

```bash
make build
```

Binary output:

- `./bin/pswitch`

## Download a release

Download the matching archive from [GitHub Releases](https://github.com/wlynxg/pswitch/releases/latest), extract it, and run the binary from the extracted directory.

## Docker

Start with Docker Compose:

```bash
docker compose up -d --build
```

Or:

```bash
docker build -t pswitch .
docker run -d --name pswitch -p 8080:8080 -v "$(pwd)/data:/data" pswitch
```

Or use the published image:

```bash
docker pull ghcr.io/wlynxg/pswitch:latest
docker run -d --name pswitch -p 8080:8080 -v "$(pwd)/data:/data" ghcr.io/wlynxg/pswitch:latest
```

Notes:

- The container working directory is `/data`.
- `settings.json` and `metrics.json` are persisted under `/data`.
- `CMD ["--config", "/data/config.toml"]` means a mounted config file is used when present, but startup still falls back to the built-in default config if the file is missing.

## Initialize config

```bash
make init
```

Or:

```bash
./bin/pswitch init --config ./config.toml
```

## Run

```bash
make run
```

Or:

```bash
./bin/pswitch
./bin/pswitch --config ./config.toml
./bin/pswitch --listen 0.0.0.0:8080 --mode least_failures --failure-threshold 2 --cooldown 30s --health-check-interval 20s --health-check-timeout 5s
```

Options:

- `--listen`
- `--mode`
- `--failure-threshold`
- `--cooldown`
- `--health-check-interval`
- `--health-check-timeout`
- `--log-color=true|false`

Notes:

- Running `./bin/pswitch` starts the service directly; there is no `serve` subcommand anymore.
- If the config file is missing, `pswitch` starts with the built-in default config.
- The default listen address is `0.0.0.0:8080`.
- The default config path is `config.toml` in the binary directory.
- The user config file is read-only from the program's perspective; dashboard saves go to `settings.json` in the current working directory.
- Dashboard metrics are persisted in `metrics.json` in the current working directory.
- If `settings.json` already exists, it takes precedence over the user config file on startup.
- `PSWITCH_ADMIN_TOKEN` is optional. If set, the admin UI and admin API require it.

## Codex

Point Codex to the local proxy:

```toml
[model_providers.OpenAI]
base_url = "http://127.0.0.1:8080/codex"
wire_api = "responses"
requires_openai_auth = true
```

## Claude Code

If you mannually add an Anthropic-compatible route later, you can point Claude Code to it:

```bash
export ANTHROPIC_BASE_URL=http://127.0.0.1:8080/claude
export ANTHROPIC_API_KEY=dummy
```

Notes:

- `/claude` currently exposes `v1/messages`, `v1/messages/count_tokens`, and `v1/models`
- `count_tokens` is an estimate
- `upstream_model` controls which real model is called behind the Claude-compatible route
