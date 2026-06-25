package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/lumberjack.v2"
)

func init() { register(Files, newFileBackend) }

// fileBackend implements Backend by routing each level to its own rotating
// JSON file under the logs/ directory.
//
// Files created (rotated daily by lumberjack):
//
//	logs/debug-YYYY-MM-DD.log  — Debug entries only
//	logs/info-YYYY-MM-DD.log   — Info entries only
//	logs/warn-YYYY-MM-DD.log   — Warn entries only
//	logs/error-YYYY-MM-DD.log  — Error and Fatal entries
type fileBackend struct {
	debug *zap.Logger
	info  *zap.Logger
	warn  *zap.Logger
	err   *zap.Logger // handles both Error and Fatal
}

func newFileBackend(_ string) Backend {
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		panic(fmt.Sprintf("fileBackend: cannot create log dir: %v", err))
	}

	today := time.Now().Format("2006-01-02")

	enc := zapcore.NewJSONEncoder(fileEncoderConfig())

	return &fileBackend{
		debug: buildFileLogger(enc, logDir, "debug", today, exactLevel(zapcore.DebugLevel)),
		info:  buildFileLogger(enc, logDir, "info", today, exactLevel(zapcore.InfoLevel)),
		warn:  buildFileLogger(enc, logDir, "warn", today, exactLevel(zapcore.WarnLevel)),
		// Error file captures Error and Fatal (>= Error)
		err: buildFileLogger(enc, logDir, "error", today, zap.LevelEnablerFunc(func(l zapcore.Level) bool {
			return l >= zapcore.ErrorLevel
		})),
	}
}

func (f *fileBackend) Debug(msg string, fields ...zap.Field) { f.debug.Debug(msg, fields...) }
func (f *fileBackend) Info(msg string, fields ...zap.Field)  { f.info.Info(msg, fields...) }
func (f *fileBackend) Warn(msg string, fields ...zap.Field)  { f.warn.Warn(msg, fields...) }
func (f *fileBackend) Error(msg string, fields ...zap.Field) { f.err.Error(msg, fields...) }

// Fatal writes at Fatal level without calling os.Exit — Logger.Fatal owns the exit.
func (f *fileBackend) Fatal(msg string, fields ...zap.Field) { f.err.Fatal(msg, fields...) }

func (f *fileBackend) With(fields ...zap.Field) Backend {
	return &fileBackend{
		debug: f.debug.With(fields...),
		info:  f.info.With(fields...),
		warn:  f.warn.With(fields...),
		err:   f.err.With(fields...),
	}
}

func (f *fileBackend) Sync() error {
	for _, l := range []*zap.Logger{f.debug, f.info, f.warn, f.err} {
		if err := l.Sync(); err != nil {
			return err
		}
	}
	return nil
}

// buildFileLogger creates a *zap.Logger that writes JSON to a rotating file.
// zap.WithFatalHook(WriteThenNoop) prevents the internal zap Fatal from calling
// os.Exit — Logger.Fatal handles the exit after all backends have written.
func buildFileLogger(enc zapcore.Encoder, dir, name, date string, enabler zapcore.LevelEnabler) *zap.Logger {
	path := filepath.Join(dir, fmt.Sprintf("%s-%s.log", name, date))
	w := &lumberjack.Logger{
		Filename:   path,
		MaxSize:    50, // MB
		MaxBackups: 7,
		MaxAge:     30, // days
		Compress:   true,
	}
	core := zapcore.NewCore(enc, zapcore.AddSync(w), enabler)
	return zap.New(core,
		zap.AddStacktrace(zapcore.ErrorLevel),
		zap.WithFatalHook(zapcore.WriteThenNoop),
	)
}

// exactLevel returns a LevelEnabler that passes only entries at exactly lvl.
func exactLevel(lvl zapcore.Level) zap.LevelEnablerFunc {
	return func(l zapcore.Level) bool { return l == lvl }
}

func fileEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      zapcore.OmitKey, // caller is injected by Logger.dispatch
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
	}
}
