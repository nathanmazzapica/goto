# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`goto` is a CLI tool for bookmarking directories and jumping between them. Users save "markers" (name-to-path mappings) and navigate using `goto <marker>` or via the `tp` shell function.

## Build & Development Commands

```bash
make build          # Build binary for current platform -> bin/goto
make build-all      # Cross-compile for all platforms (linux/darwin/windows, amd64/arm64)
make clean          # Remove local binary
```

## Code Quality (Required for PRs)

```bash
go fmt ./...        # Format code
staticcheck ./...   # Run linter (install: go install honnef.co/go/tools/cmd/staticcheck@latest)
```

## Architecture

**Entry point**: `cmd/goto.go` - parses flags, orchestrates marker operations, manages "previous" marker for recall functionality.

**Marker package** (`internal/marker/`): CRUD operations for markers stored as `name:path` pairs in `~/.config/goto/.markers`.
- `LoadMarkers()` reads markers file
- `SaveMarkers()` writes markers file (0600 permissions)
- `Add()` prevents duplicates, returns `ErrAlreadyExists`
- `Delete()` removes marker and immediately saves

**Data flow**: The binary outputs directory paths to stdout. A shell wrapper function (`tp`) captures this output and runs `cd` since external programs cannot change the shell's working directory.

**Config initialization**: `ensureDotFiles()` in `cmd/goto.go` creates `~/.config/goto/.markers` on startup.
