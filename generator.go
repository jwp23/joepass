package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
)

const (
	lowercaseChars = "abcdefghijklmnopqrstuvwxyz"
	uppercaseChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digitChars     = "0123456789"
	defaultSpecial = "!@#$%^&*()-_=+[]{}|;:,.<>?"
	ambiguousChars = "0OIl1"
)

type Options struct {
	Length      int
	NoUpper     bool
	NoDigits    bool
	NoSpecial   bool
	Special     string
	NoAmbiguous bool
}

func buildCharacterPool(opts Options) string {
	chars := lowercaseChars

	if !opts.NoUpper {
		chars += uppercaseChars
	}
	if !opts.NoDigits {
		chars += digitChars
	}
	if !opts.NoSpecial {
		special := defaultSpecial
		if opts.Special != "" {
			special = opts.Special
		}
		chars += special
	}

	if opts.NoAmbiguous {
		chars = removeAmbiguous(chars)
	}

	return chars
}

func removeAmbiguous(chars string) string {
	var filtered []byte
	for _, c := range chars {
		if !strings.ContainsRune(ambiguousChars, c) {
			filtered = append(filtered, byte(c))
		}
	}
	return string(filtered)
}

func Generate(opts Options) (string, error) {
	pool := buildCharacterPool(opts)
	if len(pool) == 0 {
		return "", fmt.Errorf("no characters available: all character types are disabled")
	}
	if opts.Length <= 0 {
		return "", fmt.Errorf("invalid length: %d", opts.Length)
	}

	result := make([]byte, opts.Length)
	poolLen := big.NewInt(int64(len(pool)))

	for i := range opts.Length {
		n, err := rand.Int(rand.Reader, poolLen)
		if err != nil {
			return "", fmt.Errorf("failed to generate random number: %w", err)
		}
		result[i] = pool[n.Int64()]
	}

	return string(result), nil
}
