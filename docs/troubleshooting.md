# Troubleshooting

- `init` writes to `config.toml` in the binary directory by default
- dashboard-edited runtime config is written to `settings.json` in the current working directory
- usage and provider stats are written to `metrics.json` in the current working directory
- `api_key` must be non-empty for each provider
- if one provider fails, `pswitch` will try the next one
- health probes only log when a provider recovers
- if usage is missing in logs, the upstream response did not include token usage
- if `/claude` is enabled, make sure its route has `upstream_model`
