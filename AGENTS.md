# AGENTS.md for qcard

## Build Commands
- `make qcal` or `make` for Linux AMD64 build
- `make linux-arm` for ARM build
- `make darwin` for macOS build
- `make windows` for Windows build
- `make install` to install to /usr/local/bin
- `make clean` to remove binaries

## Lint Commands
- `go vet ./...` for static analysis (fixes deprecation warnings)
- `golangci-lint run` if golangci-lint is installed

## Test Commands
- `go test ./...` to run all tests (none currently exist)
- `go test -run TestName` for single test (none currently exist)

## Code Style Guidelines
- Follow standard Go formatting: use `go fmt`
- Imports: group standard library, blank line, then local/project imports
- Naming: camelCase for variables/functions, PascalCase for types/constants
- Error handling: use log.Fatal for critical errors, checkError for non-fatal
- Types: define structs in PascalCase, fields in camelCase
- No unnecessary comments; code should be self-explanatory
- Use sync.WaitGroup for concurrency as in fetchAbData
- Constants for colors and formats as in defines.go
- Methods on structs use receiver in camelCase, e.g., (e contactStruct) fancyOutput()
- HTTP requests: check errors before deferring resp.Body.Close()

## Common Fixes
- Deprecated `ioutil` replaced with `io` and `os` (e.g., `ioutil.ReadFile` -> `os.ReadFile`)
- Remove unused variables/functions (e.g., unused `err`, `versionLocation`)
- Ensure proper error handling in HTTP requests