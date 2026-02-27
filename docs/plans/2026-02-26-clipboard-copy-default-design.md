# Clipboard-Copy Default Behavior

**Date:** 2026-02-26
**Status:** Approved

## Problem

Running `joepass | pbcopy` every time adds friction. The primary use case is copying a generated password to the clipboard, so the tool should do that by default.

## Design

### Behavior

`joepass` generates a password, copies it to the system clipboard, and prints `Make it so! Password go!` to stderr. No password is ever written to stdout.

### Clipboard Strategy

Try clipboard commands in priority order. Use the first binary found on `$PATH`:

| Priority | Command | Platform |
|----------|---------|----------|
| 1 | `pbcopy` | macOS |
| 2 | `wl-copy` | Linux (Wayland) |
| 3 | `xclip -selection clipboard` | Linux (X11) |
| 4 | `xsel --clipboard --input` | Linux (X11 fallback) |

This avoids OS detection and naturally handles the Linux Wayland/X11 split.

### Error Handling

- **No clipboard tool found:** Exit non-zero with: `error: no clipboard tool found. Install one of: pbcopy, wl-copy, xclip, xsel`
- **Clipboard command fails:** Exit non-zero with a description of the failure.
- **Password is never printed to stdout.** This prevents accidental exposure in terminal history, shell logs, or shoulder-surfing.

### No New Flags

YAGNI. No `--print`, `--stdout`, or `--no-copy`.

## Files Changed

| File | Change |
|------|--------|
| `main.go` | Replace `fmt.Println(pw)` with clipboard copy + stderr message |
| `clipboard.go` (new) | `CopyToClipboard(text string) error` with try-in-order logic |
| `README.md` | Remove `joepass \| pbcopy` example, describe new default behavior |

## Testing

- Test `CopyToClipboard` with integration test (skip in CI where clipboard tools are unavailable).
- Existing `generator_test.go` is unaffected since `Generate()` still returns a string.

## Decisions

- **Clipboard-only, no stdout fallback.** Falling back to stdout when no clipboard tool is found creates security risks (password in terminal history, shell logs, stale clipboard from muscle memory). Fail hard instead.
- **Try-in-order over OS detection.** `runtime.GOOS` alone cannot distinguish Wayland from X11 on Linux. Checking `$PATH` for known clipboard binaries is simpler and more portable.
- **Single success message.** `Make it so! Password go!` printed to stderr. Fun, branded, and minimal.
