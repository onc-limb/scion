# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

scion is a Go-based CLI utility designed to streamline workflow when working with AI agents like Claude Code, GitHub Copilot, or Cursor.

## Development Commands

### Go Module Initialization
```bash
go mod init github.com/[username]/scion
```

### Build
```bash
go build -o build/scion ./cmd/scion
```

### Run Tests
```bash
go test ./...
go test -v ./...  # verbose output
go test -cover ./...  # with coverage
```

### Run a Single Test
```bash
go test -run TestName ./path/to/package
```

### Format Code
```bash
go fmt ./...
```

### Lint
```bash
golangci-lint run  # if golangci-lint is installed
go vet ./...  # built-in linter
```

### Install Dependencies
```bash
go mod tidy
go mod download
```

## Project Structure

This is a Go CLI application. When implementing features:
- Main entry point should be in `main.go` or `cmd/scion/main.go`
- Use the standard Go project layout if the project grows
- Consider using cobra or urfave/cli for CLI command structure
- Follow Go naming conventions and idioms