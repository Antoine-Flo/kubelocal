package logger

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger is just a zap.Logger with a LogCommand method
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
	config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)

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

// LogCommand executes a command and logs its output using zap
func (l *Logger) LogCommand(cmd *exec.Cmd) error {
	l.Debug("Executing command",
		zap.String("command", strings.Join(cmd.Args, " ")))

	output, err := cmd.CombinedOutput()
	if len(output) > 0 {
		if len(output) > 500 {
			l.Debug("Command output (truncated)",
				zap.String("output", string(output[:500])),
				zap.Int("total_length", len(output)))
		} else {
			l.Debug("Command output", zap.String("output", string(output)))
		}
	}

	if err != nil {
		l.Debug("Command failed", zap.Error(err))
		return err
	}

	return nil
}
