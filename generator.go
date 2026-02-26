package main

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
		if opts.Special != "" {
			pool += opts.Special
		} else {
			pool += defaultSpecial
		}
	}

	if opts.NoAmbiguous {
		var filtered []byte
		for i := 0; i < len(pool); i++ {
			ambiguous := false
			for j := 0; j < len(ambiguousChars); j++ {
				if pool[i] == ambiguousChars[j] {
					ambiguous = true
					break
				}
			}
			if !ambiguous {
				filtered = append(filtered, pool[i])
			}
		}
		pool = string(filtered)
	}

	return pool
}
