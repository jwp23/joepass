# Clipboard-Copy Default Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Make joepass copy generated passwords to the system clipboard by default instead of printing to stdout.

**Architecture:** New `clipboard.go` file with a `CopyToClipboard` function that tries clipboard binaries in priority order (pbcopy, wl-copy, xclip, xsel). `main.go` calls this instead of `fmt.Println`. Password never touches stdout; success message goes to stderr.

**Tech Stack:** Go stdlib (`os/exec`), no new dependencies.

---

### Task 1: Create clipboard.go with CopyToClipboard

**Files:**
- Create: `clipboard.go`
- Create: `clipboard_test.go`

**Step 1: Write the failing test**

Create `clipboard_test.go`:

```go
package main

import (
	"os/exec"
	"testing"
)

func TestCopyToClipboard_Success(t *testing.T) {
	// Skip if no clipboard tool is available (e.g., CI)
	candidates := []string{"pbcopy", "wl-copy", "xclip", "xsel"}
	found := false
	for _, c := range candidates {
		if _, err := exec.LookPath(c); err == nil {
			found = true
			break
		}
	}
	if !found {
		t.Skip("no clipboard tool available")
	}

	err := CopyToClipboard("test-password-123")
	if err != nil {
		t.Fatalf("CopyToClipboard failed: %v", err)
	}
}

func TestCopyToClipboard_NoTool(t *testing.T) {
	// Save original PATH and set to empty to simulate no tools
	t.Setenv("PATH", "")

	err := CopyToClipboard("test-password")
	if err == nil {
		t.Fatal("expected error when no clipboard tool is available")
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test -run TestCopyToClipboard -v ./...`
Expected: FAIL — `CopyToClipboard` is not defined.

**Step 3: Write the implementation**

Create `clipboard.go`:

```go
package main

import (
	"fmt"
	"os/exec"
	"strings"
)

type clipboardCmd struct {
	bin  string
	args []string
}

var clipboardCandidates = []clipboardCmd{
	{bin: "pbcopy"},
	{bin: "wl-copy"},
	{bin: "xclip", args: []string{"-selection", "clipboard"}},
	{bin: "xsel", args: []string{"--clipboard", "--input"}},
}

func CopyToClipboard(text string) error {
	for _, c := range clipboardCandidates {
		path, err := exec.LookPath(c.bin)
		if err != nil {
			continue
		}

		cmd := exec.Command(path, c.args...)
		cmd.Stdin = strings.NewReader(text)

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("clipboard command %s failed: %w", c.bin, err)
		}

		return nil
	}

	return fmt.Errorf("no clipboard tool found. Install one of: pbcopy, wl-copy, xclip, xsel")
}
```

**Step 4: Run tests to verify they pass**

Run: `go test -run TestCopyToClipboard -v ./...`
Expected: PASS (TestCopyToClipboard_Success passes on macOS, skips in CI; TestCopyToClipboard_NoTool passes everywhere).

**Step 5: Run full test suite**

Run: `go test -v ./...`
Expected: All existing tests still pass.

**Step 6: Run vet and lint**

Run: `go vet ./...`
Expected: No issues.

**Step 7: Commit**

```bash
git add clipboard.go clipboard_test.go
git commit -m "Add CopyToClipboard with try-in-order clipboard strategy"
```

---

### Task 2: Update main.go to use clipboard

**Files:**
- Modify: `main.go:29-35` (replace `fmt.Println(pw)` with clipboard copy + stderr message)

**Step 1: Write the change**

In `main.go`, replace:

```go
	fmt.Println(pw)
```

with:

```go
	if err := CopyToClipboard(pw); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintln(os.Stderr, "Make it so! Password go!")
```

Also remove `"fmt"` from the import if it becomes unused. Check: `fmt` is still used for `fmt.Fprintf` and `fmt.Fprintln`, so keep it.

**Step 2: Build and smoke test**

Run: `go build -o joepass .`
Expected: Compiles without errors.

Run: `./joepass`
Expected: Prints `Make it so! Password go!` to terminal. Password is on clipboard — verify with `pbpaste`.

**Step 3: Run full test suite**

Run: `go test -v ./...`
Expected: All tests pass.

**Step 4: Run vet and lint**

Run: `go vet ./...`
Expected: No issues.

**Step 5: Commit**

```bash
git add main.go
git commit -m "Use clipboard for password output instead of stdout"
```

---

### Task 3: Update README.md

**Files:**
- Modify: `README.md`

**Step 1: Update the README**

Replace the Usage section to reflect the new default behavior. Key changes:
- Remove the `joepass | pbcopy` example (no longer needed)
- Add a note that password is copied to clipboard by default
- Keep all flag examples as-is

New Usage section:

```markdown
## Usage

```bash
# Generate a 20-character password (copied to clipboard)
joepass

# Custom length
joepass --length 32

# Specify allowed special characters
joepass --special '!@#$%'

# Exclude character types
joepass --no-upper
joepass --no-digits
joepass --no-special

# Exclude ambiguous characters (0, O, I, l, 1)
joepass --no-ambiguous
```
```

**Step 2: Verify markdown renders correctly**

Skim the file to make sure formatting is correct.

**Step 3: Commit**

```bash
git add README.md
git commit -m "Update README for clipboard-copy default behavior"
```

---

### Task 4: Final verification

**Step 1: Run full test suite**

Run: `go test -v ./...`
Expected: All tests pass.

**Step 2: Run vet**

Run: `go vet ./...`
Expected: No issues.

**Step 3: End-to-end smoke test**

Run: `go build -o joepass . && ./joepass`
Expected output (on stderr): `Make it so! Password go!`
Verify: `pbpaste` shows a 20-character random password.

Run: `./joepass --length 32 && pbpaste | wc -c`
Expected: 32 characters (plus newline = 33 from wc).

Run: `./joepass --no-special --no-digits && pbpaste`
Expected: Only letters in clipboard.
