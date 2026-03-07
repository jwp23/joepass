package main

import (
	"os/exec"
	"testing"
)

func TestCopyToClipboardSuccess(t *testing.T) {
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

	err := copyToClipboard("test-password-123")
	if err != nil {
		t.Fatalf("copyToClipboard failed: %v", err)
	}
}

func TestCopyToClipboardNoTool(t *testing.T) {
	// Save original PATH and set to empty to simulate no tools
	t.Setenv("PATH", "")

	err := copyToClipboard("test-password")
	if err == nil {
		t.Fatal("expected error when no clipboard tool is available")
	}
}
