package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func runInit(args []string) error {
	fs := flag.NewFlagSet("init", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	configPath := fs.String("config", defaultConfigPath(mustExecutablePath()), "path to TOML config")
	force := fs.Bool("force", false, "overwrite existing file")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if !*force {
		if _, err := os.Stat(*configPath); err == nil {
			return fmt.Errorf("%s already exists; rerun with --force to overwrite", *configPath)
		}
	}

	if err := os.MkdirAll(filepath.Dir(*configPath), 0o700); err != nil {
		return err
	}

	return os.WriteFile(*configPath, []byte(defaultConfigTemplate), 0o600)
}

const defaultConfigTemplate = `listen = "127.0.0.1:8080"
mode = "round_robin"
failure_threshold = 1
cooldown = "20s"
health_check_interval = "15s"
health_check_timeout = "3s"

[[routes]]
prefix = "/codex"
type = "openai"

# Add providers in the dashboard or uncomment examples below.
#
# [[providers]]
# name = "provider-1"
# base_url = "https://example-provider-1.com"
# api_key = "sk-your-key-1"
#
# [[providers]]
# name = "provider-2"
# base_url = "https://example-provider-2.com"
# api_key = "sk-your-key-2"
`
