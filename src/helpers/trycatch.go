package helpers

import (
	"fmt"
	"runtime/debug"

	"github.com/yourorg/go-mvc-app/src/utils/logger"
	"go.uber.org/zap"
)

// Try executes fn and recovers from any panic, returning it as an error.
// Use this to wrap any block that could panic.
//
//	err := helpers.Try(func() error {
//	    return doRiskyThing()
//	})
func Try(fn func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			stack := string(debug.Stack())
			logger.Error("panic recovered",
				zap.Any("panic", r),
				zap.String("stack", stack),
			)
			err = fmt.Errorf("panic: %v", r)
		}
	}()
	return fn()
}

// TryVoid is like Try but for functions that don't return an error.
func TryVoid(fn func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			stack := string(debug.Stack())
			logger.Error("panic recovered",
				zap.Any("panic", r),
				zap.String("stack", stack),
			)
			err = fmt.Errorf("panic: %v", r)
		}
	}()
	fn()
	return nil
}
