package logger

import (
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger wraps zap.Logger for structured logging
type Logger struct {
	*zap.Logger
}

// NewLogger creates a new zap logger with JSON output
func NewLogger() (*Logger, error) {
	// Try different log directories in order of preference
	logDirs := []string{
		"/var/log/kubelocal",
		filepath.Join(os.Getenv("HOME"), ".local/share/kubelocal/logs"),
		"./logs",
	}

	var logFile string
	for _, dir := range logDirs {
		if err := os.MkdirAll(dir, 0755); err == nil {
			logFile = filepath.Join(dir, "install.log")
			break
		}
	}

	if logFile == "" {
		return nil, fmt.Errorf("failed to create log directory")
	}

	// Configure zap with JSON encoder
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{logFile}
	config.Encoding = "json"
	config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)

	zapLogger, err := config.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build zap logger: %w", err)
	}

	return &Logger{zapLogger}, nil
}

// Close closes the logger and flushes remaining data
func (l *Logger) Close() error {
	return l.Sync()
}
