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
