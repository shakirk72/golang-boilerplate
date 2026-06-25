package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func init() { register(Console, newConsoleBackend) }

// consoleBackend implements Backend by writing to stdout.
// Format: pretty console encoder in "dev", JSON in "qa" / "prod".
type consoleBackend struct {
	log *zap.Logger
}

func newConsoleBackend(env string) Backend {
	cfg := fileEncoderConfig()
	cfg.CallerKey = zapcore.OmitKey // caller injected by Logger.dispatch

	var enc zapcore.Encoder
	if env == "prod" || env == "qa" {
		enc = zapcore.NewJSONEncoder(cfg)
	} else {
		enc = zapcore.NewConsoleEncoder(cfg)
	}

	core := zapcore.NewCore(enc, zapcore.AddSync(os.Stdout), zapcore.DebugLevel)
	log := zap.New(core,
		zap.AddStacktrace(zapcore.ErrorLevel),
		zap.WithFatalHook(zapcore.WriteThenNoop), // Logger.Fatal owns os.Exit
	)
	return &consoleBackend{log: log}
}

func (c *consoleBackend) Debug(msg string, fields ...zap.Field) { c.log.Debug(msg, fields...) }
func (c *consoleBackend) Info(msg string, fields ...zap.Field)  { c.log.Info(msg, fields...) }
func (c *consoleBackend) Warn(msg string, fields ...zap.Field)  { c.log.Warn(msg, fields...) }
func (c *consoleBackend) Error(msg string, fields ...zap.Field) { c.log.Error(msg, fields...) }
func (c *consoleBackend) Fatal(msg string, fields ...zap.Field) { c.log.Fatal(msg, fields...) }

func (c *consoleBackend) With(fields ...zap.Field) Backend {
	return &consoleBackend{log: c.log.With(fields...)}
}

func (c *consoleBackend) Sync() error {
	return c.log.Sync()
}
