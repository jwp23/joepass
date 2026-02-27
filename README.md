# joepass

A CLI tool for generating random passwords.

## Install

```bash
go install github.com/jwp23/joepass@latest
```

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

## Defaults

| Setting | Value |
|---------|-------|
| Length | 20 |
| Lowercase | enabled |
| Uppercase | enabled |
| Digits | enabled |
| Special | `!@#$%^&*()-_=+[]{}|;:,.<>?` |
