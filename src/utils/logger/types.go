package logger

import "go.uber.org/zap"

type BackendName string

const (
	Files   BackendName = "files"
	Console BackendName = "console"
)

type Backend interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field) // must NOT call os.Exit — Logger.Fatal owns it
	With(fields ...zap.Field) Backend
	Sync() error
}

type namedBackend struct {
	name    BackendName
	backend Backend
}

type Logger struct {
	backends   []namedBackend
	callerSkip int
}
