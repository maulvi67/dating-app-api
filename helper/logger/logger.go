package logger

import (
	"dating-apps/helper/config"
	"os"
	"strings"

	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type Logger interface {
	Info(msg string, keyvals ...interface{})
	Error(msg string, keyvals ...interface{})
	Debug(msg string, keyvals ...interface{})
	Warn(msg string, keyvals ...interface{})
}

type loggerImpl struct {
	kitLogger kitlog.Logger
}

func NewLogger(logConfig *config.LogConfig) Logger {
	return &loggerImpl{
		kitLogger: NewGoKitLog(logConfig),
	}
}

// NewGoKitLog creates a new Go Kit logger configured according to the given logConfig.
func NewGoKitLog(logConfig *config.LogConfig) kitlog.Logger {
	// Create a base logger that writes to stdout.
	baseLogger := kitlog.NewLogfmtLogger(os.Stdout)

	// Set up the allowed log level.
	var filterOption level.Option
	switch strings.ToLower(logConfig.Level) {
	case "debug":
		filterOption = level.AllowDebug()
	case "info":
		filterOption = level.AllowInfo()
	case "warn", "warning":
		filterOption = level.AllowWarn()
	case "error":
		filterOption = level.AllowError()
	default:
		filterOption = level.AllowInfo()
	}

	// Wrap the base logger with a level filter.
	filteredLogger := level.NewFilter(baseLogger, filterOption)

	// Add default key/value pairs for timestamp and caller info.
	return kitlog.With(filteredLogger, "ts", kitlog.DefaultTimestampUTC, "caller", kitlog.DefaultCaller)
}

// Info logs an informational message.
func (l *loggerImpl) Info(msg string, keyvals ...interface{}) {
	args := append([]interface{}{"msg", msg}, keyvals...)
	_ = level.Info(l.kitLogger).Log(args...)
}

// Error logs an error message.
func (l *loggerImpl) Error(msg string, keyvals ...interface{}) {
	args := append([]interface{}{"msg", msg}, keyvals...)
	_ = level.Error(l.kitLogger).Log(args...)
}

// Debug logs a debug message.
func (l *loggerImpl) Debug(msg string, keyvals ...interface{}) {
	args := append([]interface{}{"msg", msg}, keyvals...)
	_ = level.Debug(l.kitLogger).Log(args...)
}

// Warn logs a warning message.
func (l *loggerImpl) Warn(msg string, keyvals ...interface{}) {
	args := append([]interface{}{"msg", msg}, keyvals...)
	_ = level.Warn(l.kitLogger).Log(args...)
}
