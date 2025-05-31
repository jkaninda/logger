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

	"gopkg.in/natefinch/lumberjack.v2"
)

type LogLevel string

const (
	levelDebug   LogLevel = "debug"
	levelInfo    LogLevel = "info"
	levelWarning LogLevel = "warning"
	levelError   LogLevel = "error"
	levelOff     LogLevel = "off"
)

// Logger is the main logger struct that wraps slog.Logger with additional configuration
type Logger struct {
	logger     *slog.Logger   // Underlying slog logger instance
	level      LogLevel       // Current logging level threshold
	outputFile string         // Path to log file if file logging is enabled
	maxAgeDays int            // Maximum number of days to retain old log files
	jsonFormat bool           // Whether to use JSON formatting
	disabled   bool           // Whether logging is completely disabled
	fileWriter io.WriteCloser // File writer for log rotation
	maxBackups int            // Maximum number of old log files to retain
	maxSize    int            // Maximum size in MB before log rotation
	compress   bool           // Whether to compress rotated log files
}

// ConfigFunc defines the type for configuration functions that modify Logger settings
type ConfigFunc func(*Logger)

// ********* Log Level Configuration Functions **********

// WithDebugLevel sets the log level to debug
func WithDebugLevel() ConfigFunc { return func(l *Logger) { l.level = levelDebug } }

// WithInfoLevel sets the log level to info
func WithInfoLevel() ConfigFunc { return func(l *Logger) { l.level = levelInfo } }

// WithWarningLevel sets the log level to warning
func WithWarningLevel() ConfigFunc { return func(l *Logger) { l.level = levelWarning } }

// WithErrorLevel sets the log level to error
func WithErrorLevel() ConfigFunc { return func(l *Logger) { l.level = levelError } }

// WithLevelOff disables all logging
func WithLevelOff() ConfigFunc { return func(l *Logger) { l.level = levelOff } }

// ************ Output Configuration Functions **********

// WithFile configures file output with the given filename
func WithFile(outputFile string) ConfigFunc {
	return func(l *Logger) { l.outputFile = outputFile }
}

// WithOutPutFile same to WithFile, configures file output with the given filename
func WithOutPutFile(outputFile string) ConfigFunc {
	return func(l *Logger) { l.outputFile = outputFile }
}

// WithMaxAgeInDays sets the maximum number of days to retain old log files
func WithMaxAgeInDays(days int) ConfigFunc {
	return func(l *Logger) { l.maxAgeDays = days }
}

// WithMaxSize sets the maximum size in MB before log rotation occurs
func WithMaxSize(size int) ConfigFunc {
	return func(l *Logger) { l.maxSize = size }
}

// WithCompress enables compression of rotated log files
func WithCompress() ConfigFunc {
	return func(l *Logger) { l.compress = true }
}

// WithMaxBackups sets the maximum number of old log files to retain
func WithMaxBackups(maxBackups int) ConfigFunc {
	return func(l *Logger) { l.maxBackups = maxBackups }
}

// WithJsonFormat enables JSON formatting for log output
func WithJsonFormat() ConfigFunc {
	return func(l *Logger) { l.jsonFormat = true }
}

// ****** Logger Setup Functions ********

// New creates a new Logger instance with optional configurations
// Defaults:
//   - Level: error
//   - Output: stdout
//   - Max age: 7 days
//   - Format: text
func New(configs ...ConfigFunc) *Logger {
	// Initialize with default values
	l := &Logger{
		level:      levelError,
		outputFile: "",
		maxAgeDays: 7,
		logger:     slog.Default(),
	}

	// Apply any provided configurations
	return l.With(configs...)
}
func Default() *Logger {
	// Initialize with default values
	return &Logger{
		level:      levelInfo,
		outputFile: "",
		maxAgeDays: 7,
		maxBackups: 3,
		logger:     slog.Default(),
	}
}

// With applies the given configurations to the Logger and returns the modified Logger
func (l *Logger) With(configs ...ConfigFunc) *Logger {
	// Apply each configuration function
	for _, conf := range configs {
		conf(l)
	}

	// Handle special case where logging is completely disabled
	if l.level == levelOff {
		l.disabled = true
		l.logger = slog.New(slog.NewTextHandler(io.Discard, nil))
		return l
	}

	// Default to stdout output
	var output io.Writer = os.Stdout

	// Configure file output if specified
	if l.outputFile != "" {
		lumberjackWriter := &lumberjack.Logger{
			Filename:   l.outputFile,
			MaxAge:     l.maxAgeDays,
			MaxBackups: l.maxBackups,
			MaxSize:    l.maxSize,
			Compress:   l.compress,
		}
		l.fileWriter = lumberjackWriter
		output = lumberjackWriter
	}

	// Convert our log level to slog's level type
	slogLevel := l.toSlogLevel(l.level)
	opts := &slog.HandlerOptions{Level: slogLevel}

	// Initialize the appropriate handler based on format preference
	if l.jsonFormat {
		l.logger = slog.New(slog.NewJSONHandler(output, opts))
	} else {
		l.logger = slog.New(slog.NewTextHandler(output, opts))
	}

	return l
}

// ******** Logging Methods **********

// Debug logs a message at debug level with optional key-value pairs
func (l *Logger) Debug(msg string, args ...interface{}) {
	if !l.disabled {
		l.logger.Debug(msg, args...)
	}
}

// Info logs a message at info level with optional key-value pairs
func (l *Logger) Info(msg string, args ...interface{}) {
	if !l.disabled {
		l.logger.Info(msg, args...)
	}
}

// Warning logs a message at warning level with optional key-value pairs
func (l *Logger) Warning(msg string, args ...interface{}) {
	if !l.disabled {
		l.logger.Warn(msg, args...)
	}
}

// Error logs a message at error level with optional key-value pairs
func (l *Logger) Error(msg string, args ...interface{}) {
	if !l.disabled {
		l.logger.Error(msg, args...)
	}
}

// Fatal logs a message at error level and exits the program with status 1
func (l *Logger) Fatal(msg string, args ...interface{}) {
	l.logger.Error(msg, args...)
	os.Exit(1)
}

// IsDebugEnabled returns true if debug logging is enabled
func (l *Logger) IsDebugEnabled() bool {
	return l.level == levelDebug && !l.disabled
}

// Close cleans up resources (like file handles) used by the logger
func (l *Logger) Close() error {
	if l.fileWriter != nil {
		return l.fileWriter.Close()
	}
	return nil
}

// toSlogLevel converts our LogLevel type to slog.Level
func (l *Logger) toSlogLevel(level LogLevel) slog.Level {
	switch level {
	case levelDebug:
		return slog.LevelDebug
	case levelInfo:
		return slog.LevelInfo
	case levelWarning:
		return slog.LevelWarn
	case levelError:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
