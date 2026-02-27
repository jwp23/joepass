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

func buildPool(opts Options) string {
	pool := lowercaseChars

	if !opts.NoUpper {
		pool += uppercaseChars
	}
	if !opts.NoDigits {
		pool += digitChars
	}
	if !opts.NoSpecial {
		special := defaultSpecial
		if opts.Special != "" {
			special = opts.Special
		}
		pool += special
	}

	if opts.NoAmbiguous {
		var filtered []byte
		for _, c := range pool {
			if !strings.ContainsRune(ambiguousChars, c) {
				filtered = append(filtered, byte(c))
			}
		}
		pool = string(filtered)
	}

	return pool
}

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

	for i := range opts.Length {
		n, err := rand.Int(rand.Reader, poolLen)
		if err != nil {
			return "", fmt.Errorf("failed to generate random number: %w", err)
		}
		result[i] = pool[n.Int64()]
	}

	return string(result), nil
}
