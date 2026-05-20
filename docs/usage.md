# Usage

## Build

```bash
make build
```

Binary output:

- `./bin/pswitch`

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
```

Options:

- `--listen`
- `--mode`
- `--log-color=true|false`

Notes:

- Running `./bin/pswitch` starts the service directly; there is no `serve` subcommand anymore.
- If the config file is missing, `pswitch` starts with the built-in default config.
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

If you manually add an Anthropic-compatible route later, you can point Claude Code to it:

```bash
export ANTHROPIC_BASE_URL=http://127.0.0.1:8080/claude
export ANTHROPIC_API_KEY=dummy
```

Notes:

- `/claude` currently exposes `v1/messages`, `v1/messages/count_tokens`, and `v1/models`
- `count_tokens` is an estimate
- `upstream_model` controls which real model is called behind the Claude-compatible route
