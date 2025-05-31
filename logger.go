/*
 *  MIT License
 *
 * Copyright (c) 2025 Jonas Kaninda
 *
 *  Permission is hereby granted, free of charge, to any person obtaining a copy
 *  of this software and associated documentation files (the "Software"), to deal
 *  in the Software without restriction, including without limitation the rights
 *  to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *  copies of the Software, and to permit persons to whom the Software is
 *  furnished to do so, subject to the following conditions:
 *
 *  The above copyright notice and this permission notice shall be included in all
 *  copies or substantial portions of the Software.
 *
 *  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *  AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *  LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 *  SOFTWARE.
 */

package logger

import (
	"io"
	"log/slog"
	"os"
	"sync"

	"gopkg.in/natefinch/lumberjack.v2"
)

// LogLevel represents the severity levels for log messages
type LogLevel string

// Supported log levels constants
const (
	LevelDebug   LogLevel = "debug"   // Debug level for detailed debugging information
	LevelInfo    LogLevel = "info"    // Info level for general operational messages
	LevelWarning LogLevel = "warning" // Warning level for potentially harmful situations
	LevelError   LogLevel = "error"   // Error level for error events
	LevelOff     LogLevel = "off"     // Off level completely disables logging
)

// Config holds all configurable parameters for the logger
type Config struct {
	Level      LogLevel // Minimum log level to output
	OutputFile string   // Path to log file (empty for stdout)
	MaxAgeDays int      // Maximum number of days to retain old log files
	MaxBackups int      // Maximum number of old log files to retain
	MaxSizeMB  int      // Maximum size in megabytes of the log file before rotation
	Compress   bool     // Whether to compress rotated log files
	JSONFormat bool     // Whether to use JSON formatting
	UseDefault bool     // Whether to use slog.Default() format instead of custom handler
}

// Logger is the main logger struct that wraps slog.Logger with additional features
type Logger struct {
	*slog.Logger                    // Embedded slog.Logger for core logging functionality
	config       Config             // Current configuration
	file         *lumberjack.Logger // File writer for log rotation (nil when using stdout)
	mu           sync.Mutex         // Mutex for thread-safe operations
	disabled     bool               // Flag indicating if logging is completely disabled
}

// Option defines the type for configuration functions that modify Logger settings
type Option func(*Config)

// Default returns a new logger using slog.Default() configuration
// Output format: "2025/05/31 09:00:08 INFO Application started version=1.0.0 config=default"
func Default() *Logger {
	return &Logger{
		Logger:   slog.Default(),           // Use Go's standard library default logger
		config:   Config{UseDefault: true}, // Mark as using default format
		disabled: false,
	}
}

// New creates a new Logger instance with customizable options
// Default configuration:
//   - Level: info
//   - MaxAgeDays: 7
//   - MaxBackups: 3
//   - MaxSizeMB: 100
//   - Format: text
//   - Output: stdout
func New(opts ...Option) *Logger {
	// Set default configuration values
	cfg := Config{
		Level:      LevelInfo,
		MaxAgeDays: 7,
		MaxBackups: 3,
		MaxSizeMB:  100,
		UseDefault: false,
	}

	// Apply all provided configuration options
	for _, opt := range opts {
		opt(&cfg)
	}

	// Create logger instance
	l := &Logger{
		config:   cfg,
		disabled: cfg.Level == LevelOff, // Disable if level is "off"
	}

	// Initialize the underlying logger
	l.initLogger()
	return l
}

// initLogger initializes or re-initializes the underlying slog.Logger based on current config
// This method is thread-safe through mutex locking
func (l *Logger) initLogger() {
	l.mu.Lock()
	defer l.mu.Unlock()

	// If logging is disabled, use a discard handler
	if l.disabled {
		l.Logger = slog.New(slog.NewTextHandler(io.Discard, nil))
		return
	}

	// If using default format, use slog's default logger
	if l.config.UseDefault {
		l.Logger = slog.Default()
		return
	}

	// Default to stdout output
	var output io.Writer = os.Stdout

	// Configure file output if specified
	if l.config.OutputFile != "" {
		l.file = &lumberjack.Logger{
			Filename:   l.config.OutputFile, // Log file path
			MaxAge:     l.config.MaxAgeDays, // Days to retain logs
			MaxBackups: l.config.MaxBackups, // Number of backups to keep
			MaxSize:    l.config.MaxSizeMB,  // Max file size in MB
			Compress:   l.config.Compress,   // Whether to compress backups
		}
		output = l.file
	}

	// Create handler options with configured log level
	opts := &slog.HandlerOptions{
		Level: l.toSlogLevel(l.config.Level),
	}

	// Initialize the appropriate handler based on format preference
	if l.config.JSONFormat {
		l.Logger = slog.New(slog.NewJSONHandler(output, opts))
	} else {
		l.Logger = slog.New(slog.NewTextHandler(output, opts))
	}
}

// WithOptions creates a new Logger instance with additional configuration options applied
// Returns a new Logger instance, leaving the original unchanged (immutable pattern)
func (l *Logger) WithOptions(opts ...Option) *Logger {
	// Start with current config
	newCfg := l.config

	// Apply all provided options
	for _, opt := range opts {
		opt(&newCfg)
	}

	// Create and return new logger with updated config
	newLogger := &Logger{
		config:   newCfg,
		disabled: newCfg.Level == LevelOff,
	}
	newLogger.initLogger()
	return newLogger
}

// Close releases resources used by the logger (primarily file handles)
// Should be called when the logger is no longer needed
func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

// ******************** Utility Methods ******************** /

// IsDebugEnabled returns true if debug level logging is enabled
func (l *Logger) IsDebugEnabled() bool {
	return l.config.Level == LevelDebug && !l.disabled
}

// IsInfoEnabled returns true if info level logging is enabled
func (l *Logger) IsInfoEnabled() bool {
	return (l.config.Level == LevelDebug || l.config.Level == LevelInfo) && !l.disabled
}

// IsWarningEnabled returns true if warning level logging is enabled
func (l *Logger) IsWarningEnabled() bool {
	return l.config.Level != LevelError && l.config.Level != LevelOff && !l.disabled
}

// IsErrorEnabled returns true if error level logging is enabled
func (l *Logger) IsErrorEnabled() bool {
	return l.config.Level != LevelOff && !l.disabled
}

// GetLevel returns the current log level setting
func (l *Logger) GetLevel() LogLevel {
	return l.config.Level
}

// GetConfig returns a copy of the current logger configuration
func (l *Logger) GetConfig() Config {
	return l.config
}

// *************** Configuration Options ***************/

// WithLevel sets the minimum log level for output
func WithLevel(level LogLevel) Option {
	return func(c *Config) { c.Level = level }
}

// WithOutputFile configures file output with the given path
func WithOutputFile(file string) Option {
	return func(c *Config) { c.OutputFile = file }
}

// WithMaxAge sets maximum days to retain log files
func WithMaxAge(days int) Option {
	return func(c *Config) { c.MaxAgeDays = days }
}

// WithMaxBackups sets maximum number of old log files to keep
func WithMaxBackups(count int) Option {
	return func(c *Config) { c.MaxBackups = count }
}

// WithMaxSize sets maximum log file size in megabytes before rotation
func WithMaxSize(sizeMB int) Option {
	return func(c *Config) { c.MaxSizeMB = sizeMB }
}

// WithCompression enables compression of rotated log files
func WithCompression() Option {
	return func(c *Config) { c.Compress = true }
}

// WithJSONFormat enables JSON formatting for log output
func WithJSONFormat() Option {
	return func(c *Config) { c.JSONFormat = true }
}

// ************ Predefined Level Options ************/

// WithDebugLevel sets log level to debug
func WithDebugLevel() Option { return WithLevel(LevelDebug) }

// WithInfoLevel sets log level to info
func WithInfoLevel() Option { return WithLevel(LevelInfo) }

// WithWarningLevel sets log level to warning
func WithWarningLevel() Option { return WithLevel(LevelWarning) }

// WithErrorLevel sets log level to error
func WithErrorLevel() Option { return WithLevel(LevelError) }

// WithLevelOff completely disables logging output
func WithLevelOff() Option { return WithLevel(LevelOff) }

// WithDefaultFormat configures the logger to use slog.Default() format
func WithDefaultFormat() Option {
	return func(c *Config) { c.UseDefault = true }
}

// toSlogLevel converts our LogLevel type to slog.Level
func (l *Logger) toSlogLevel(level LogLevel) slog.Level {
	switch level {
	case LevelDebug:
		return slog.LevelDebug
	case LevelInfo:
		return slog.LevelInfo
	case LevelWarning:
		return slog.LevelWarn
	case LevelError:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
