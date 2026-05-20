# Configuration

`pswitch` uses a TOML config file.

```toml
listen = "127.0.0.1:8080"
mode = "round_robin"
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

Fields:

- `listen`: local bind address
- `mode`: `sequential` or `round_robin`
- `failure_threshold`: consecutive failures before a provider is circuit-broken
- `cooldown`: how long to wait before probing a broken provider again
- `health_check_interval`: probe interval
- `health_check_timeout`: probe timeout

Admin console:

- The embedded admin UI is served at `/dashboard/`
- `PSWITCH_ADMIN_TOKEN` is optional; if it is set, the admin UI and admin API require it
- The admin API accepts the token in `X-PSwitch-Admin-Token`
- Changes to `mode`, health-check settings, `routes`, and `providers` are hot reloaded after save
- Changes saved from the admin UI are written to `settings.json` in the current working directory
- The original TOML config file is only used as startup input and is not modified by the program
- Changes to `listen` are persisted to `settings.json` but still require a process restart to take effect

Route fields:

- `prefix`: URL prefix, such as `/codex` or `/claude`
- `type`: `openai` or `anthropic`
- `model`: advertised model name for Anthropic clients
- `upstream_model`: actual upstream model to call for Anthropic routes

Each provider needs:

- `name`
- `base_url`
- `api_key`

Use `pswitch init --config ./config.toml` to generate a template.

If no `routes` are configured, `pswitch` falls back to a single OpenAI-compatible route on `/`.
