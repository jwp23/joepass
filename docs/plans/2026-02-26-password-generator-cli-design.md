# Password Generator CLI — Design

## Overview

A personal CLI tool (`passgen`) that generates random passwords and prints them to stdout. Built in Go using only the standard library. Designed to pipe to clipboard utilities like `pbcopy`.

## CLI Interface

```
passgen                        # 20 chars, all character types
passgen --length 32            # custom length
passgen --special '!@#$%'      # specify allowed special characters
passgen --no-ambiguous         # exclude 0/O/I/l/1
passgen --no-upper             # disable uppercase
passgen --no-digits            # disable digits
passgen --no-special           # disable special characters
passgen | pbcopy               # pipe to clipboard
```

### Defaults

| Setting | Default |
|---------|---------|
| Length | 20 |
| Lowercase | enabled |
| Uppercase | enabled |
| Digits | enabled |
| Special | enabled |
| Special chars | `!@#$%^&*()-_=+[]{}|;:,.<>?` |

Output is the password string followed by a newline. Nothing else.

## Architecture

### Project Structure

```
password-generator-cli/
├── main.go             # Entry point, flag parsing, output
├── generator.go        # Password generation logic
├── generator_test.go   # Tests
├── go.mod
└── README.md
```

### Character Sets

```
lowercase  = "abcdefghijklmnopqrstuvwxyz"
uppercase  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
digits     = "0123456789"
special    = "!@#$%^&*()-_=+[]{}|;:,.<>?"
ambiguous  = "0OIl1"
```

### Generation Algorithm

1. Build the character pool by concatenating enabled sets.
2. If `--special` is provided, replace the default special set with the user's string.
3. If `--no-ambiguous`, remove ambiguous characters from the pool.
4. Use `crypto/rand` to pick characters uniformly from the pool.
5. Print the result to stdout.

No forced "at least one of each type" logic. Pure random selection from the full pool is simpler and preserves maximum entropy.

### Error Handling

- Empty character pool (all types disabled): print error to stderr, exit 1.
- Invalid length (0 or negative): print error to stderr, exit 1.

## Testing

All tests use Go's built-in `testing` package.

- Default length is 20.
- `--length` flag changes output length.
- `--no-upper`, `--no-digits`, `--no-special` exclude their character sets.
- `--special` replaces the default special characters.
- `--no-ambiguous` removes ambiguous characters.
- Error on empty character pool.
- Error on invalid length.

## CI/CD

GitHub Actions workflow (`.github/workflows/ci.yml`) on push and PR to `main`:

1. `go vet` — catches common mistakes.
2. `golangci-lint` — curated set of linters (staticcheck, errcheck, ineffassign, etc.).
3. `govulncheck` — checks for known vulnerabilities in the Go stdlib.
4. `go test ./...` — runs all tests.

No dependency scanning needed (zero external dependencies).

## Decisions

| Decision | Rationale |
|----------|-----------|
| Go, stdlib only | Single binary, no runtime deps, instant startup |
| `flag` for CLI | Built-in, sufficient for this use case |
| `crypto/rand` | Cryptographically secure randomness |
| No config file | Sensible defaults + flags keeps it simple |
| No forced character variety | Pure random preserves entropy, avoids complexity |
| Moderate CI | `go vet` + `golangci-lint` + `govulncheck` + tests |
