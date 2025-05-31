# LOGGER

[![Tests](https://github.com/jkaninda/logger/actions/workflows/tests.yml/badge.svg)](https://github.com/jkaninda/logger/actions/workflows/tests.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/jkaninda/logger)](https://goreportcard.com/report/github.com/jkaninda/logger)
[![Go](https://img.shields.io/github/go-mod/go-version/jkaninda/logger)](https://go.dev/)
[![Go Reference](https://pkg.go.dev/badge/github.com/jkaninda/logger.svg)](https://pkg.go.dev/github.com/jkaninda/logger)
[![GitHub Release](https://img.shields.io/github/v/release/jkaninda/logger)](https://github.com/jkaninda/logger/releases)

**Logger** provides a configurable logging solution with multiple output options, log levels, and rotation capabilities built on top of Go's slog package.

## Installation

```bash
go get github.com/jkaninda/logger
```
## Usage Example

```go
l := logger.New(
    WithInfoLevel(),
    WithMaxAge(1),
    WithMaxSize(100),
    WithJSONFormat(),
)
l.Info("Application started", "version", "1.0.0")
```
## Default config

```go
	logger := logger.Default()
	logger.Info("Application started", "version", "1.0.0", "config", "default")
```
---

## Contributing

Contributions are welcome!

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to your fork
5. Open a Pull Request



---
## Give a Star! ⭐

⭐ If you find Okapi useful, please consider giving it a star on [GitHub](https://github.com/jkaninda/logger)!


## License

This project is licensed under the MIT License. See the LICENSE file for details.

## Copyright

Copyright (c) 2025 Jonas Kaninda
