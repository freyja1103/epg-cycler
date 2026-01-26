# AGENTS.md

This document provides essential information for agentic coding assistants working with this Go codebase.

## Build Commands

```bash
# Build the entire project
go build -v ./...

# Build a specific package
go build -v ./cmd/epg-cycler

# Install the binary globally
go install ./cmd/epg-cycler
```

## Lint Commands

Go's standard formatting tools are used for linting:

```bash
# Format all Go files
go fmt ./...

# Vet the code for suspicious constructs
go vet ./...

# Run tests with race detection
go test -race ./...
```

## Test Commands

```bash
# Run all tests
go test -v ./...

# Run a specific test
go test -v ./edcb-api-client -run TestEDCBAPI

# Run tests with coverage
go test -cover ./...

# Run tests with coverage and output to file
go test -coverprofile=coverage.out ./...

# View coverage report in browser
go tool cover -html=coverage.out
```

## Code Style Guidelines

### Imports

1. Standard library imports come first, grouped and alphabetized
2. Third-party imports come second, grouped and alphabetized
3. Local project imports come last, grouped and alphabetized
4. Use meaningful aliases for imports when needed to avoid conflicts

Example:

```go
import (
    "context"
    "fmt"
    "net/http"

    "github.com/some/third-party"
    "gopkg.in/natefinch/lumberjack.v2"

    "github.com/freyja1103/epg-cycler/internal/pkg"
)
```

### Formatting

1. Use `go fmt` for all code formatting
2. Line length should be kept reasonable (preferably under 100 characters)
3. Use tabs for indentation (standard Go practice)
4. Comments should be full sentences with proper punctuation

### Types

1. Use descriptive names for types, variables, and functions
2. Exported identifiers should have comments
3. Prefer structs over maps for structured data
4. Use interfaces to define behavior contracts

### Naming Conventions

1. Use camelCase for variables and functions
2. Use PascalCase for exported identifiers
3. Use UPPER_CASE for constants
4. Acronyms should follow Go conventions (URL, not Url)
5. Receiver names should be short (one or two letters)
6. Interface names should describe behavior (Reader, Writer, Interface suffix acceptable)

### Error Handling

1. Always handle errors explicitly
2. Wrap errors with context using `fmt.Errorf("message: %w", err)`
3. Don't ignore errors with `_`
4. Use `errors.Is()` and `errors.As()` for error comparison
5. Log errors appropriately with context

### Logging

1. Use the `log/slog` package for structured logging
2. Include relevant context in log messages
3. Use appropriate log levels (Debug, Info, Warn, Error)
4. Use `slog.Any()` for complex values

Example:

```go
slog.ErrorContext(ctx, "failed to process request", slog.Any("error", err))
```

### Documentation

1. All exported functions, types, and variables should have comments
2. Comments should explain why, not what
3. Use Godoc style comments
4. Include examples for complex functions

### Testing

1. Place tests in the same package with `_test` suffix
2. Use table-driven tests when appropriate
3. Name test files with `_test.go` suffix
4. Use meaningful test function names: `Test<Object>_<Case>`
5. Test edge cases and error conditions

### Context Usage

1. Pass context to functions that might block or take time
2. Check for context cancellation in long-running operations
3. Use context values sparingly and only for request-scoped data

### Dependencies

1. Minimize external dependencies
2. Pin dependency versions in `go.mod`
3. Regularly update dependencies
4. Use Go modules for dependency management

### Project Structure

1. Use standard Go project layout:
   - `/cmd` - main applications
   - `/internal` - private application code
   - `/pkg` - library code that can be used by external projects
   - `/api` - OpenAPI/Swagger specs, protocol definition files

### API Design

1. Define clean interfaces for external APIs
2. Use context for request-scoped data and cancellation
3. Return concrete types rather than interfaces where possible
4. Handle timeouts and retries appropriately

### Concurrency

1. Use goroutines and channels for concurrent operations
2. Protect shared resources with mutexes or other synchronization primitives
3. Avoid data races
4. Prefer communicating through channels rather than sharing memory

### Performance

1. Profile code before optimizing
2. Avoid premature optimization
3. Use appropriate data structures for the task
4. Consider memory allocation impact

## Project Specifics

This project is an EPG (Electronic Program Guide) cycler for EDCB (EDCB is a TV recording system commonly used in Japan). It manages TV recordings and performs post-recording operations like organizing files and shutting down systems.

Key components:

- EDCB API client for communicating with EDCB systems
- Syobocal API client for accessing Japanese anime program information
- Engine for processing recordings and managing the cycling logic

## Common Patterns

1. Use of interfaces for API clients to enable mocking in tests
2. Structured logging with context
3. Error wrapping with context for better debugging
4. Time parsing for Japanese timezone (Asia/Tokyo)
5. File operations for organizing recorded media
6. OS-specific shutdown commands for Windows and Linux
