package main

import (
	"context"
	"errors"
	"flag"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"pswitch/internal/config"
	"pswitch/internal/logx"
	pruntime "pswitch/internal/runtime"
	"pswitch/internal/server"
)

type serveArgs struct {
	ConfigPath string
	Listen     string
	Mode       string
	LogColor   *bool
}

type triStateBool struct {
	value bool
	set   bool
}

func (b *triStateBool) String() string {
	if !b.set {
		return ""
	}
	if b.value {
		return "true"
	}
	return "false"
}

func (b *triStateBool) Set(value string) error {
	b.set = true
	switch value {
	case "", "true", "1", "yes", "on":
		b.value = true
		return nil
	case "false", "0", "no", "off":
		b.value = false
		return nil
	default:
		return errors.New("log-color must be true or false")
	}
}

func (b *triStateBool) IsBoolFlag() bool {
	return true
}

func parseServeArgs(args []string) (serveArgs, error) {
	fs := flag.NewFlagSet("serve", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	configPath := fs.String("config", defaultConfigPath(mustExecutablePath()), "path to TOML config")
	overrideListen := fs.String("listen", "", "optional listen address override")
	overrideMode := fs.String("mode", "", "optional mode override: sequential, round_robin, or least_failures")
	var logColor triStateBool
	fs.Var(&logColor, "log-color", "enable or disable colored logs")

	if err := fs.Parse(args); err != nil {
		return serveArgs{}, err
	}

	out := serveArgs{
		ConfigPath: *configPath,
		Listen:     *overrideListen,
		Mode:       *overrideMode,
	}
	if logColor.set {
		out.LogColor = &logColor.value
	}
	return out, nil
}

func runServe(args []string) error {
	parsed, err := parseServeArgs(args)
	if err != nil {
		return err
	}

	stateDir, err := defaultStateDir()
	if err != nil {
		return err
	}
	settingsPath := defaultSettingsPath(stateDir)
	metricsPath := defaultMetricsPath(stateDir)

	cfg, err := loadStartupConfig(parsed.ConfigPath, settingsPath)
	if err != nil {
		return err
	}
	if parsed.Listen != "" {
		cfg.Listen = parsed.Listen
	}
	if parsed.Mode != "" {
		cfg.Mode = parsed.Mode
	}
	if err := cfg.Validate(); err != nil {
		return err
	}

	syncLogs := logx.Init(parsed.LogColor)
	defer syncLogs()

	manager, err := pruntime.New(settingsPath, metricsPath, cfg)
	if err != nil {
		return err
	}
	adminToken, err := resolveAdminToken(cfg.Listen)
	if err != nil {
		return err
	}

	server := &http.Server{
		Addr:              cfg.Listen,
		Handler:           server.NewRouter(manager, adminToken),
		ReadHeaderTimeout: 10 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go healthLoop(ctx, manager)

	listener, err := net.Listen("tcp", cfg.Listen)
	if err != nil {
		return err
	}

	currentCfg := manager.Config()
	logx.Infof("proxy started listen=%s mode=%s providers=%d routes=%d", listener.Addr().String(), currentCfg.Mode, len(currentCfg.Providers), len(currentCfg.Routes))
	logx.Infof("runtime files config=%s settings=%s metrics=%s", parsed.ConfigPath, settingsPath, metricsPath)

	errCh := make(chan error, 1)
	go func() {
		errCh <- server.Serve(listener)
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return server.Shutdown(shutdownCtx)
	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
}

func healthLoop(ctx context.Context, manager *pruntime.Manager) {
	for {
		cfg, providerPool := manager.Snapshot()
		timer := time.NewTimer(cfg.HealthCheckInterval)
		select {
		case <-ctx.Done():
			if !timer.Stop() {
				<-timer.C
			}
			return
		case now := <-timer.C:
			client := &http.Client{Timeout: cfg.HealthCheckTimeout}
			events := providerPool.ProbeDue(ctx, client, now)
			for _, event := range events {
				logx.Infof("health recovered provider=%s", event.Provider)
			}
		}
	}
}

func resolveAdminToken(listen string) (string, error) {
	return strings.TrimSpace(os.Getenv("PSWITCH_ADMIN_TOKEN")), nil
}

func loadStartupConfig(configPath, settingsPath string) (config.Config, error) {
	cfg, err := config.LoadJSON(settingsPath)
	if err == nil {
		return cfg, nil
	}
	if !errors.Is(err, os.ErrNotExist) {
		var pathErr *os.PathError
		if !errors.As(err, &pathErr) || !errors.Is(pathErr.Err, os.ErrNotExist) {
			return config.Config{}, err
		}
	}

	return loadUserConfigOrDefault(configPath)
}

func loadUserConfigOrDefault(path string) (config.Config, error) {
	cfg, err := config.Load(path)
	if err == nil {
		return cfg, nil
	}

	if errors.Is(err, os.ErrNotExist) {
		return config.Default(), nil
	}

	var pathErr *os.PathError
	if errors.As(err, &pathErr) && errors.Is(pathErr.Err, os.ErrNotExist) {
		return config.Default(), nil
	}

	return config.Config{}, err
}

func defaultStateDir() (string, error) {
	return os.Getwd()
}

func defaultConfigPath(executablePath string) string {
	dir := filepath.Dir(executablePath)
	return filepath.Join(dir, "config.toml")
}

func defaultSettingsPath(stateDir string) string {
	return filepath.Join(stateDir, "settings.json")
}

func defaultMetricsPath(stateDir string) string {
	return filepath.Join(stateDir, "metrics.json")
}

func mustExecutablePath() string {
	path, err := os.Executable()
	if err != nil || strings.TrimSpace(path) == "" {
		return filepath.Join(".", "pswitch")
	}
	return path
}
