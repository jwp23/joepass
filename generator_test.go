package main

import (
	"strings"
	"testing"
)

func TestBuildPoolDefaults(t *testing.T) {
	opts := options{
		Length:      20,
		NoUpper:     false,
		NoDigits:    false,
		NoSpecial:   false,
		Special:     "",
		NoAmbiguous: false,
	}
	pool := buildCharacterPool(opts)

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
	for _, c := range digitChars {
		if !strings.ContainsRune(pool, c) {
			t.Errorf("pool missing digit %c", c)
		}
	}
	if !strings.ContainsRune(pool, '!') || !strings.ContainsRune(pool, '@') {
		t.Error("pool missing default special characters")
	}
}

func TestBuildPoolNoUpper(t *testing.T) {
	opts := options{NoUpper: true}
	pool := buildCharacterPool(opts)
	for _, c := range uppercaseChars {
		if strings.ContainsRune(pool, c) {
			t.Errorf("pool should not contain uppercase %c", c)
		}
	}
}

func TestBuildPoolNoDigits(t *testing.T) {
	opts := options{NoDigits: true}
	pool := buildCharacterPool(opts)
	for _, c := range digitChars {
		if strings.ContainsRune(pool, c) {
			t.Errorf("pool should not contain digit %c", c)
		}
	}
}

func TestBuildPoolNoSpecial(t *testing.T) {
	opts := options{NoSpecial: true}
	pool := buildCharacterPool(opts)
	for _, c := range defaultSpecial {
		if strings.ContainsRune(pool, c) {
			t.Errorf("pool should not contain special %c", c)
		}
	}
}

func TestBuildPoolCustomSpecial(t *testing.T) {
	opts := options{Special: "!@#"}
	pool := buildCharacterPool(opts)
	if !strings.ContainsRune(pool, '!') {
		t.Error("pool missing custom special !")
	}
	if strings.ContainsRune(pool, '^') {
		t.Error("pool should not contain ^ when custom special is set")
	}
}

func TestBuildPoolNoAmbiguous(t *testing.T) {
	opts := options{NoAmbiguous: true}
	pool := buildCharacterPool(opts)
	for _, c := range ambiguousChars {
		if strings.ContainsRune(pool, c) {
			t.Errorf("pool should not contain ambiguous char %c", c)
		}
	}
}

func TestGenerateDefaultLength(t *testing.T) {
	opts := options{Length: 20}
	pw, err := generate(opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pw) != 20 {
		t.Errorf("expected length 20, got %d", len(pw))
	}
}

func TestGenerateCustomLength(t *testing.T) {
	opts := options{Length: 32}
	pw, err := generate(opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pw) != 32 {
		t.Errorf("expected length 32, got %d", len(pw))
	}
}

func TestGenerateOnlyUsesPoolChars(t *testing.T) {
	opts := options{Length: 100, NoUpper: true, NoDigits: true, NoSpecial: true}
	pw, err := generate(opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, c := range pw {
		if c < 'a' || c > 'z' {
			t.Errorf("password contains non-lowercase char: %c", c)
		}
	}
}

func TestGenerateEmptyPool(t *testing.T) {
	// Lowercase is always included by buildCharacterPool, so we need to test
	// the error path by testing generate with length 0 instead.
	opts := options{Length: 0}
	_, err := generate(opts)
	if err == nil {
		t.Error("expected error for length 0")
	}
}

func TestGenerateNegativeLength(t *testing.T) {
	opts := options{Length: -5}
	_, err := generate(opts)
	if err == nil {
		t.Error("expected error for negative length")
	}
}
