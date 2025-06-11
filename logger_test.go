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

import "testing"

func TestDefault(t *testing.T) {
	l := Default()
	l.Info("Application started", "version", "1.0.0", "config", "default")
}
func TestInfo(t *testing.T) {
	l := New(
		WithCaller(),
		WithInfoLevel(),
		WithMaxAge(1),
		WithMaxSize(100),
		WithJSONFormat(),
	)
	l.Info("Application started", "version", "1.0.0")
	Info("Hello logger")
}
func TestWithFile(t *testing.T) {
	l := New(
		WithInfoLevel(),
		WithOutputFile("app.log"),
		WithMaxAge(1),
		WithMaxSize(100),
		WithJSONFormat(),
	)
	l.Error("Hello error message", "data", 1)
	l.Info("Hello info message", "data", 1)
	l.Warn("Hello warning message", "data", 1)
	l.Debug("Hello debug message", "data", 1)

	l.Info("Application started", "version", "1.0.0")
}
func TestWithMaxBackups(t *testing.T) {
	l := New(
		WithInfoLevel(),
		WithOutputFile("app.log"),
		WithMaxAge(30),
		WithMaxSize(100),
		WithJSONFormat(),
	)

	l.Error("Hello error message", "data", 1)
	l.Info("Hello info message", "data", 1)
	l.Warn("Hello warning message", "data", 1)
	l.Debug("Hello debug message", "data", 1)
	l.Info("Application started", "version", "1.0.0")
}

func TestWithOptions(t *testing.T) {
	// Create initial logger
	l := New(
		WithInfoLevel(),
		WithOutputFile("test.log"),
	)

	// Apply new options
	newLogger := l.WithOptions(
		WithDebugLevel(),
		WithJSONFormat(),
	)

	// Verify changes were applied
	if newLogger.GetLevel() != LevelDebug {
		t.Error("Debug level not applied")
	}
	if !newLogger.GetConfig().JSONFormat {
		t.Error("JSON format not applied")
	}
	if newLogger.GetConfig().OutputFile != "test.log" {
		t.Error("Original output file not preserved")
	}
}
