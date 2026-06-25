package logger

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"go.uber.org/zap"
)

var defaultLog *Logger
var registry = map[BackendName]func(env string) Backend{}

func register(name BackendName, factory func(env string) Backend) {
	registry[name] = factory
}

func Init(env string) {
	names := resolveBackendNames()
	bb := make([]namedBackend, 0, len(names))
	for _, name := range names {
		if factory, ok := registry[name]; ok {
			bb = append(bb, namedBackend{name, factory(env)})
		}
	}
	if len(bb) == 0 {
		if factory, ok := registry[Files]; ok {
			bb = append(bb, namedBackend{Files, factory(env)})
		}
	}
	activeNames := make([]string, len(bb))
	for i, nb := range bb {
		activeNames[i] = string(nb.name)
	}
	defaultLog = &Logger{backends: bb, callerSkip: 2}
	Info("logger initialized",
		zap.String("env", env),
		zap.String("backends", strings.Join(activeNames, ",")),
	)
}

func resolveBackendNames() []BackendName {
	raw := os.Getenv("LOGS_TARGETS")
	if raw == "" {
		return []BackendName{Files}
	}
	parts := strings.Split(raw, ",")
	names := make([]BackendName, 0, len(parts))
	for _, p := range parts {
		switch BackendName(strings.TrimSpace(strings.ToLower(p))) {
		case Files:
			names = append(names, Files)
		case Console:
			names = append(names, Console)
		}
	}
	return names
}

func To(targets ...BackendName) *Logger {
	set := make(map[BackendName]bool, len(targets))
	for _, t := range targets {
		set[t] = true
	}
	filtered := make([]namedBackend, 0, len(targets))
	for _, nb := range defaultLog.backends {
		if set[nb.name] {
			filtered = append(filtered, nb)
		}
	}
	return &Logger{backends: filtered, callerSkip: 2}
}

func (l *Logger) Debug(msg string, fields ...zap.Field) { l.dispatch(msg, fields, Backend.Debug) }
func (l *Logger) Info(msg string, fields ...zap.Field)  { l.dispatch(msg, fields, Backend.Info) }
func (l *Logger) Warn(msg string, fields ...zap.Field)  { l.dispatch(msg, fields, Backend.Warn) }
func (l *Logger) Error(msg string, fields ...zap.Field) { l.dispatch(msg, fields, Backend.Error) }

func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.dispatch(msg, fields, Backend.Fatal)
	l.Sync()
	os.Exit(1)
}

func (l *Logger) Sync() {
	for _, nb := range l.backends {
		_ = nb.backend.Sync()
	}
}

func (l *Logger) dispatch(msg string, fields []zap.Field, fn func(Backend, string, ...zap.Field)) {
	all := make([]zap.Field, 0, len(fields)+1)
	all = append(all, callerField(2+l.callerSkip))
	all = append(all, fields...)
	for _, nb := range l.backends {
		fn(nb.backend, msg, all...)
	}
}

func callerField(skip int) zap.Field {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return zap.Skip()
	}
	if i := strings.LastIndexByte(file, '/'); i >= 0 {
		file = file[i+1:]
	}
	return zap.String("caller", fmt.Sprintf("%s:%d", file, line))
}

func Debug(msg string, fields ...zap.Field) { defaultLog.Debug(msg, fields...) }
func Info(msg string, fields ...zap.Field)  { defaultLog.Info(msg, fields...) }
func Warn(msg string, fields ...zap.Field)  { defaultLog.Warn(msg, fields...) }
func Error(msg string, fields ...zap.Field) { defaultLog.Error(msg, fields...) }
func Fatal(msg string, fields ...zap.Field) { defaultLog.Fatal(msg, fields...) }
func Sync()                                 { defaultLog.Sync() }

func WithRequestID(requestID string) *Logger {
	f := zap.String("request_id", requestID)
	tagged := make([]namedBackend, len(defaultLog.backends))
	for i, nb := range defaultLog.backends {
		tagged[i] = namedBackend{name: nb.name, backend: nb.backend.With(f)}
	}
	return &Logger{backends: tagged, callerSkip: 1}
}
