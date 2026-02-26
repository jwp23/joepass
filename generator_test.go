package main

import (
	"strings"
	"testing"
)

func TestBuildPool_Defaults(t *testing.T) {
	opts := Options{
		Length:      20,
		NoUpper:     false,
		NoDigits:    false,
		NoSpecial:   false,
		Special:     "",
		NoAmbiguous: false,
	}
	pool := buildPool(opts)

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
