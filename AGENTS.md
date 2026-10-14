# AGENTS.md

This file guides agentic coding assistants working in this repository.

## Project summary
- `goto` is a Go CLI for bookmarking directories and printing destinations.
- Entry point: `cmd/goto.go`.
- Marker storage: `internal/marker` reads/writes `~/.config/goto/.markers`.
- CLI prints paths to stdout so the `tp` shell function can `cd`.

## Commands
### Build
- `make build` builds `bin/goto` for the host platform.
- `make build-all` cross-compiles for linux/darwin/windows (amd64/arm64).
- `make clean` removes `bin/goto`.
- `go build -o bin/goto ./cmd/goto.go` matches the Makefile target.

### Format
- `go fmt ./...` (required in CI per README).

### Lint
- `staticcheck ./...` (required in CI per README).
- Install with `go install honnef.co/go/tools/cmd/staticcheck@latest` if needed.

### Test
- `go test ./...` runs all tests.
- `go test ./cmd` runs the CLI tests/benchmarks package only.
- Single test (regex):
  - `go test ./cmd -run ^TestMigrateOldMarkers_OldFileDoesNotExist$`
- Single test group:
  - `go test ./cmd -run ^TestMigrateOldMarkers_`
- Benchmarks:
  - `go test ./cmd -bench BenchmarkSortKeys -run ^$`

## Code style
### Formatting
- Use `gofmt`-compatible formatting; run `go fmt ./...` before final output.
- Use tabs for indentation (default Go formatting).

### Imports
- Standard library imports first, then blank line, then module imports.
- Keep imports minimal; remove unused imports.
- Use canonical Go import paths (no relative imports).

### Naming
- Use Go naming conventions (camel case, exported `CamelCase`).
- Error variables should be named `ErrSomething` and exported only if needed.
- Keep function names verbs for actions (`LoadMarkers`, `SaveMarkers`).
- Avoid one-letter names except common loop indices.

### Types and data structures
- Prefer standard library types and minimal custom types unless necessary.
- Marker data is `map[string]string` keyed by marker name.
- Use slices for ordered output and sort with `sort.Strings`.

### Error handling
- Return errors up the call chain; avoid panics.
- Wrap with context using `fmt.Errorf("context: %w", err)` where helpful.
- Use `os.IsNotExist` checks for filesystem conditions.
- CLI commands should print human-readable errors and `os.Exit(1)`.
- Non-fatal conditions in CLI can print warnings and continue.

### File operations
- Marker file lives at `~/.config/goto/.markers`.
- Directory creation uses `os.MkdirAll(path, 0o700)`.
- Marker file permissions are `0o600` (keep sensitive by default).
- Use `#nosec G304` when file paths are not user-controlled, as seen in `cmd/goto.go`.
- Avoid changing permissions or paths without updating related tests.

### CLI behavior
- The CLI prints destination paths to stdout when teleporting.
- `--print` should only output the path, no extra text.
- `--list` prints a table with sorted markers and aligned columns.
- `--recall` uses the `previous` marker and updates it on navigation.
- `--add` and `--delete` operate on the last CLI arg as the marker name.
- Errors should be user-facing and short; avoid stack traces.

### Marker file format
- Each marker is stored as `name:path` on its own line.
- Marker names are unique map keys; later entries overwrite earlier ones.
- Keep parsing with `strings.Cut` and validate the `:` separator.
- When writing, preserve newline-separated pairs with no trailing newline.

### Paths and filesystem
- Always use `filepath.Join` for building paths.
- Use `os.UserHomeDir` to resolve the home directory.
- `LoadMarkers` should tolerate empty files (no markers).
- Use `os.IsNotExist` to distinguish missing files.

### Logging and output
- Use `fmt.Println`/`fmt.Printf` for user output in CLI flows.
- Use `log.Fatalf` only for early startup failures.
- Avoid noisy debug logging; keep stdout clean for shell usage.

### CLI argument handling
- Use `flag` for parsing booleans and call `flag.Parse()` once.
- The last CLI argument is the marker name target.
- Exit after handling `--list`, `--add`, `--delete`, or `--recall`.
- Keep `--print` output clean for shell capture.

### Migration and recall
- `ensureDotFiles` must create config dir and open the markers file.
- `migrateOldMarkers` copies old data only when new file is empty or missing.
- Treat missing old file as a no-op; use `os.IsNotExist` checks.
- `setRecall` updates the `previous` marker before navigation.

### Benchmarks
- Benchmarks live in `cmd/goto_test.go` and use `Benchmark*` naming.
- Use `go test ./cmd -bench BenchmarkSortKeys -run ^$` to skip tests.

## Testing guidance
- Tests for migration logic live in `cmd/goto_test.go`.
- New tests should prefer `t.TempDir()` for filesystem isolation.
- Avoid writing to real user paths in tests.
- Use table-driven tests for multiple similar cases.

## Review checklist
- Run `go fmt ./...` before final output.
- Run `staticcheck ./...` when possible.
- Run `go test ./...` or at least `go test ./cmd`.
- Keep `--print` output path-only.
- Preserve marker file permissions and directory mode.
- Update README when flags or output change.

## Commit messages
- Use Conventional Commit style prefixes (e.g., `feat: ...`, `fix: ...`, `docs: ...`, `chore: ...`, `refactor: ...`, `test: ...`, `perf: ...`, `ci: ...`, `build: ...`, `revert: ...`).
- Follow https://www.conventionalcommits.org/en/v1.0.0/ for structure and examples.

## Security considerations
- Keep marker data local; do not add network calls.
- Avoid writing to arbitrary paths without validation.
- Preserve `#nosec G304` annotations where required.
- Do not log marker paths unless needed for output.

## Repository structure
- `cmd/` contains the CLI entry point and CLI tests/benchmarks.
- `internal/marker/` contains marker CRUD logic and file I/O.
- `bin/` is the build output directory (created by Makefile).

## Dependencies
- Go 1.23 is required (`go.mod`).
- No external runtime dependencies beyond the standard library.
- Keep new dependencies minimal and justify them if added.

## Documentation notes
- Update README if user-facing commands or flags change.
- Keep CLI examples in sync with actual behavior.

## Cursor/Copilot rules
- No `.cursor/rules`, `.cursorrules`, or `.github/copilot-instructions.md` found.
