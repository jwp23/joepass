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
