package logger

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Logger interface defines the logging methods
type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	LogCommand(cmd *exec.Cmd) error
	Close() error
}

// FileLogger implements Logger interface with file output
type FileLogger struct {
	writer   *bufio.Writer
	file     *os.File
	errorDir string
}

// NewLogger creates a new file logger
func NewLogger() (Logger, error) {
	// Create log directory
	logDir := "/var/log/kubelocal"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	// Create error directory
	errorDir := filepath.Join(logDir, "errors")
	if err := os.MkdirAll(errorDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create error directory: %w", err)
	}

	// Open log file
	logFile := filepath.Join(logDir, "install.log")
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	// Create buffered writer with 1-second flush interval
	writer := bufio.NewWriter(file)

	// Start flush timer
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			writer.Flush()
		}
	}()

	return &FileLogger{
		writer:   writer,
		file:     file,
		errorDir: errorDir,
	}, nil
}

// Close closes the logger and flushes remaining data
func (l *FileLogger) Close() error {
	l.writer.Flush()
	return l.file.Close()
}

// log writes a formatted log message
func (l *FileLogger) log(level, msg string, args ...interface{}) {
	timestamp := time.Now().Format("15:04:05")
	formattedMsg := fmt.Sprintf(msg, args...)
	logEntry := fmt.Sprintf("[%s] [%s] %s\n", timestamp, level, formattedMsg)

	l.writer.WriteString(logEntry)
}

// Debug logs a debug message
func (l *FileLogger) Debug(msg string, args ...interface{}) {
	l.log("DEBUG", msg, args...)
}

// Info logs an info message
func (l *FileLogger) Info(msg string, args ...interface{}) {
	l.log("INFO", msg, args...)
}

// Warn logs a warning message
func (l *FileLogger) Warn(msg string, args ...interface{}) {
	l.log("WARN", msg, args...)
}

// Error logs an error message
func (l *FileLogger) Error(msg string, args ...interface{}) {
	l.log("ERROR", msg, args...)
}

// LogCommand executes a command and logs its output
func (l *FileLogger) LogCommand(cmd *exec.Cmd) error {
	l.Debug("Executing: %s", strings.Join(cmd.Args, " "))

	// Capture stdout and stderr
	output, err := cmd.CombinedOutput()
	outputStr := string(output)

	// Log command output (truncated)
	if len(outputStr) > 0 {
		if len(outputStr) > 500 {
			l.Debug("Command output (truncated): %s", outputStr[:500])
			// Save full output to error file if too long
			l.saveToErrorFile("command_output", outputStr, cmd.Args[0])
		} else {
			l.Debug("Command output: %s", outputStr)
		}
	}

	if err != nil {
		l.Debug("Command failed: %v", err)
		return err
	}

	return nil
}

// saveToErrorFile saves long content to a separate error file
func (l *FileLogger) saveToErrorFile(fileType, content, component string) error {
	timestamp := time.Now().Format("15-04-05")
	filename := fmt.Sprintf("%s_%s_%s.log", timestamp, component, fileType)
	filepath := filepath.Join(l.errorDir, filename)

	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create error file: %w", err)
	}
	defer file.Close()

	// Write content with header
	header := fmt.Sprintf("Error details for %s at %s\n", component, time.Now().Format("2006-01-02 15:04:05"))
	header += strings.Repeat("=", 50) + "\n\n"

	file.WriteString(header)
	file.WriteString(content)

	l.Info("Full %s details saved to: %s", fileType, filepath)
	return nil
}
