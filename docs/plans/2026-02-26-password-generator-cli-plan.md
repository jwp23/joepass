# Password Generator CLI Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Build a personal CLI tool (`passgen`) that generates random passwords and prints them to stdout for piping to clipboard utilities.

**Architecture:** Single Go module with two source files — `main.go` for flag parsing and output, `generator.go` for password generation logic. Uses only the standard library (`crypto/rand`, `flag`). GitHub Actions CI with linting and vulnerability checks.

**Tech Stack:** Go (stdlib only), GitHub Actions, golangci-lint, govulncheck

---

### Task 1: Initialize Go Module

**Files:**
- Create: `go.mod`

**Step 1: Initialize the module**

Run:
```bash
go mod init github.com/jwp23/password-generator-cli
```

Expected: `go.mod` created with module name and Go version.

**Step 2: Commit**

```bash
git add go.mod
git commit -m "chore: initialize Go module"
```

---

### Task 2: Write Generator — Character Pool Building

**Files:**
- Create: `generator.go`
- Create: `generator_test.go`

**Step 1: Write the failing test for default character pool**

In `generator_test.go`:

```go
package main

import (
	"strings"
	"testing"
)

func TestBuildPool_Defaults(t *testing.T) {
	opts := Options{
		Length:      20,
		NoUpper:    false,
		NoDigits:   false,
		NoSpecial:  false,
		Special:    "",
		NoAmbiguous: false,
	}
	pool := buildPool(opts)

	// Should contain all character types
	for _, c := range "abcxyz" {
		if !strings.ContainsRune(pool, c) {
			t.Errorf("pool missing lowercase %c", c)
		}
	}
	for _, c := range "ABCXYZ" {
		if !strings.ContainsRune(pool, c) {
			t.Errorf("pool missing uppercase %c", c)
		}
	}
	for _, c := range "0123456789" {
		if !strings.ContainsRune(pool, c) {
			t.Errorf("pool missing digit %c", c)
		}
	}
	if !strings.ContainsRune(pool, '!') || !strings.ContainsRune(pool, '@') {
		t.Error("pool missing default special characters")
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test -run TestBuildPool_Defaults -v`

Expected: FAIL — `buildPool` not defined.

**Step 3: Write minimal implementation**

In `generator.go`:

```go
package main

const (
	lowercaseChars = "abcdefghijklmnopqrstuvwxyz"
	uppercaseChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digitChars     = "0123456789"
	defaultSpecial = "!@#$%^&*()-_=+[]{}|;:,.<>?"
	ambiguousChars = "0OIl1"
)

// Options holds configuration for password generation.
type Options struct {
	Length      int
	NoUpper     bool
	NoDigits    bool
	NoSpecial   bool
	Special     string
	NoAmbiguous bool
}

func buildPool(opts Options) string {
	pool := lowercaseChars

	if !opts.NoUpper {
		pool += uppercaseChars
	}
	if !opts.NoDigits {
		pool += digitChars
	}
	if !opts.NoSpecial {
		if opts.Special != "" {
			pool += opts.Special
		} else {
			pool += defaultSpecial
		}
	}

	if opts.NoAmbiguous {
		var filtered []byte
		for i := 0; i < len(pool); i++ {
			ambiguous := false
			for j := 0; j < len(ambiguousChars); j++ {
				if pool[i] == ambiguousChars[j] {
					ambiguous = true
					break
				}
			}
			if !ambiguous {
				filtered = append(filtered, pool[i])
			}
		}
		pool = string(filtered)
	}

	return pool
}
```

**Step 4: Run test to verify it passes**

Run: `go test -run TestBuildPool_Defaults -v`

Expected: PASS

**Step 5: Commit**

```bash
git add generator.go generator_test.go
git commit -m "feat: add character pool building with Options struct"
```

---

### Task 3: Write Generator — Pool Exclusion Tests

**Files:**
- Modify: `generator_test.go`

**Step 1: Write failing tests for exclusion flags**

Append to `generator_test.go`:

```go
func TestBuildPool_NoUpper(t *testing.T) {
	opts := Options{NoUpper: true}
	pool := buildPool(opts)
	for _, c := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		if strings.ContainsRune(pool, c) {
			t.Errorf("pool should not contain uppercase %c", c)
		}
	}
}

func TestBuildPool_NoDigits(t *testing.T) {
	opts := Options{NoDigits: true}
	pool := buildPool(opts)
	for _, c := range "0123456789" {
		if strings.ContainsRune(pool, c) {
			t.Errorf("pool should not contain digit %c", c)
		}
	}
}

func TestBuildPool_NoSpecial(t *testing.T) {
	opts := Options{NoSpecial: true}
	pool := buildPool(opts)
	for _, c := range defaultSpecial {
		if strings.ContainsRune(pool, c) {
			t.Errorf("pool should not contain special %c", c)
		}
	}
}

func TestBuildPool_CustomSpecial(t *testing.T) {
	opts := Options{Special: "!@#"}
	pool := buildPool(opts)
	if !strings.ContainsRune(pool, '!') {
		t.Error("pool missing custom special !")
	}
	if strings.ContainsRune(pool, '^') {
		t.Error("pool should not contain ^ when custom special is set")
	}
}

func TestBuildPool_NoAmbiguous(t *testing.T) {
	opts := Options{NoAmbiguous: true}
	pool := buildPool(opts)
	for _, c := range ambiguousChars {
		if strings.ContainsRune(pool, c) {
			t.Errorf("pool should not contain ambiguous char %c", c)
		}
	}
}
```

**Step 2: Run tests to verify they pass**

Run: `go test -run TestBuildPool -v`

Expected: All PASS (implementation already handles these cases).

**Step 3: Commit**

```bash
git add generator_test.go
git commit -m "test: add pool exclusion and custom special character tests"
```

---

### Task 4: Write Generator — Password Generation Function

**Files:**
- Modify: `generator_test.go`
- Modify: `generator.go`

**Step 1: Write the failing test for Generate**

Append to `generator_test.go`:

```go
func TestGenerate_DefaultLength(t *testing.T) {
	opts := Options{Length: 20}
	pw, err := Generate(opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pw) != 20 {
		t.Errorf("expected length 20, got %d", len(pw))
	}
}

func TestGenerate_CustomLength(t *testing.T) {
	opts := Options{Length: 32}
	pw, err := Generate(opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pw) != 32 {
		t.Errorf("expected length 32, got %d", len(pw))
	}
}

func TestGenerate_OnlyUsesPoolChars(t *testing.T) {
	opts := Options{Length: 100, NoUpper: true, NoDigits: true, NoSpecial: true}
	pw, err := Generate(opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, c := range pw {
		if c < 'a' || c > 'z' {
			t.Errorf("password contains non-lowercase char: %c", c)
		}
	}
}
```

**Step 2: Run tests to verify they fail**

Run: `go test -run TestGenerate -v`

Expected: FAIL — `Generate` not defined.

**Step 3: Write minimal implementation**

Append to `generator.go`:

```go
import (
	"crypto/rand"
	"fmt"
	"math/big"
)
```

Add the `Generate` function:

```go
// Generate creates a random password based on the given options.
func Generate(opts Options) (string, error) {
	pool := buildPool(opts)
	if len(pool) == 0 {
		return "", fmt.Errorf("no characters available: all character types are disabled")
	}
	if opts.Length <= 0 {
		return "", fmt.Errorf("invalid length: %d", opts.Length)
	}

	result := make([]byte, opts.Length)
	poolLen := big.NewInt(int64(len(pool)))

	for i := 0; i < opts.Length; i++ {
		n, err := rand.Int(rand.Reader, poolLen)
		if err != nil {
			return "", fmt.Errorf("failed to generate random number: %w", err)
		}
		result[i] = pool[n.Int64()]
	}

	return string(result), nil
}
```

Note: the `import` block should be merged with any existing imports at the top of `generator.go`. The final file should have a single `import` block with `"crypto/rand"`, `"fmt"`, and `"math/big"`.

**Step 4: Run tests to verify they pass**

Run: `go test -run TestGenerate -v`

Expected: PASS

**Step 5: Commit**

```bash
git add generator.go generator_test.go
git commit -m "feat: add Generate function using crypto/rand"
```

---

### Task 5: Write Generator — Error Cases

**Files:**
- Modify: `generator_test.go`

**Step 1: Write the failing tests for error cases**

Append to `generator_test.go`:

```go
func TestGenerate_EmptyPool(t *testing.T) {
	opts := Options{
		Length:    20,
		NoUpper:   true,
		NoDigits:  true,
		NoSpecial: true,
	}
	// Lowercase is always included by buildPool, so we need to test
	// the error path by testing Generate with length 0 instead.
	opts2 := Options{Length: 0}
	_, err := Generate(opts2)
	if err == nil {
		t.Error("expected error for length 0")
	}
}

func TestGenerate_NegativeLength(t *testing.T) {
	opts := Options{Length: -5}
	_, err := Generate(opts)
	if err == nil {
		t.Error("expected error for negative length")
	}
}
```

**Step 2: Run tests to verify they pass**

Run: `go test -v`

Expected: All PASS (implementation already handles these cases).

**Step 3: Commit**

```bash
git add generator_test.go
git commit -m "test: add error case tests for invalid length"
```

---

### Task 6: Write Main — Flag Parsing and Output

**Files:**
- Create: `main.go`

**Step 1: Write main.go**

```go
package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	length := flag.Int("length", 20, "password length")
	special := flag.String("special", "", "allowed special characters (replaces defaults)")
	noUpper := flag.Bool("no-upper", false, "exclude uppercase letters")
	noDigits := flag.Bool("no-digits", false, "exclude digits")
	noSpecial := flag.Bool("no-special", false, "exclude special characters")
	noAmbiguous := flag.Bool("no-ambiguous", false, "exclude ambiguous characters (0OIl1)")

	flag.Parse()

	opts := Options{
		Length:      *length,
		NoUpper:     *noUpper,
		NoDigits:    *noDigits,
		NoSpecial:   *noSpecial,
		Special:     *special,
		NoAmbiguous: *noAmbiguous,
	}

	pw, err := Generate(opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(pw)
}
```

**Step 2: Build and test manually**

Run:
```bash
go build -o passgen .
./passgen
./passgen -length 32
./passgen -no-special
./passgen -special '!@#'
./passgen -no-ambiguous
./passgen -length 0
```

Expected:
- Default: 20-character password with mixed characters
- `-length 32`: 32-character password
- `-no-special`: password without special characters
- `-special '!@#'`: password with only `!`, `@`, `#` as specials
- `-no-ambiguous`: password without `0`, `O`, `I`, `l`, `1`
- `-length 0`: error message to stderr, exit code 1

**Step 3: Commit**

```bash
git add main.go
git commit -m "feat: add main with flag parsing and password output"
```

---

### Task 7: CI/CD — GitHub Actions Workflow

**Files:**
- Create: `.github/workflows/ci.yml`

**Step 1: Write the workflow file**

```yaml
name: CI

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod

      - name: Vet
        run: go vet ./...

      - name: Lint
        uses: golangci/golangci-lint-action@v6
        with:
          version: latest

      - name: Vulnerability check
        run: |
          go install golang.org/x/vuln/cmd/govulncheck@latest
          govulncheck ./...

      - name: Test
        run: go test -v ./...
```

**Step 2: Commit**

```bash
git add .github/workflows/ci.yml
git commit -m "ci: add GitHub Actions with vet, lint, govulncheck, and tests"
```

---

### Task 8: README

**Files:**
- Create: `README.md`

**Step 1: Write README**

```markdown
# passgen

A CLI tool for generating random passwords.

## Install

```bash
go install github.com/jwp23/password-generator-cli@latest
```

## Usage

```bash
# Generate a 20-character password (default)
passgen

# Copy to clipboard (macOS)
passgen | pbcopy

# Custom length
passgen -length 32

# Specify allowed special characters
passgen -special '!@#$%'

# Exclude character types
passgen -no-upper
passgen -no-digits
passgen -no-special

# Exclude ambiguous characters (0, O, I, l, 1)
passgen -no-ambiguous
```

## Defaults

| Setting | Value |
|---------|-------|
| Length | 20 |
| Lowercase | enabled |
| Uppercase | enabled |
| Digits | enabled |
| Special | `!@#$%^&*()-_=+[]{}|;:,.<>?` |
```

**Step 2: Commit**

```bash
git add README.md
git commit -m "docs: add README with install and usage instructions"
```

---

### Task 9: Final Verification

**Step 1: Run all tests**

Run: `go test -v ./...`

Expected: All tests pass.

**Step 2: Run vet and build**

Run:
```bash
go vet ./...
go build -o passgen .
```

Expected: No warnings, clean build.

**Step 3: Manual smoke test**

Run:
```bash
./passgen
./passgen -length 5
./passgen -no-special -no-digits
./passgen | wc -c
```

Expected: Correct output for each case. `wc -c` should show 21 (20 chars + newline).
