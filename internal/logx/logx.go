package logx

import (
	"os"
	"strings"
	"sync"

	"github.com/mattn/go-isatty"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	mu     sync.RWMutex
	logger = zap.NewNop().Sugar()
)

func Init(colorOverride *bool) func() {
	base := newLogger(resolveColorEnabled(colorOverride, terminalSupportsColor()))
	return Use(base)
}

func Use(base *zap.Logger) func() {
	if base == nil {
		base = zap.NewNop()
	}

	mu.Lock()
	logger = base.Sugar()
	mu.Unlock()

	return func() {
		_ = base.Sync()
	}
}

func Infof(format string, args ...any) {
	L().Infof(format, args...)
}

func Warnf(format string, args ...any) {
	L().Warnf(format, args...)
}

func Debugf(format string, args ...any) {
	L().Debugf(format, args...)
}

func L() *zap.SugaredLogger {
	mu.RLock()
	defer mu.RUnlock()
	return logger
}

func resolveColorEnabled(override *bool, supported bool) bool {
	if override != nil {
		return *override
	}
	return supported
}

func terminalSupportsColor() bool {
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	if strings.EqualFold(os.Getenv("TERM"), "dumb") {
		return false
	}

	fd := os.Stdout.Fd()
	return isatty.IsTerminal(fd) || isatty.IsCygwinTerminal(fd)
}

func newLogger(color bool) *zap.Logger {
	encoderCfg := zapcore.EncoderConfig{
		TimeKey:          "time",
		LevelKey:         "level",
		MessageKey:       "msg",
		EncodeTime:       zapcore.ISO8601TimeEncoder,
		EncodeLevel:      uppercaseLevelEncoder,
		ConsoleSeparator: " ",
	}

	if color {
		encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderCfg),
		zapcore.AddSync(os.Stdout),
		zapcore.DebugLevel,
	)
	return zap.New(core)
}

func uppercaseLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(strings.ToUpper(level.String()))
}
