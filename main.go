package main

import (
	"fmt"
	"os"

	flag "github.com/spf13/pflag"
)

func main() {
	length := flag.Int("length", 20, "password length")
	special := flag.String("special", "", "allowed special characters (replaces defaults)")
	noUpper := flag.Bool("no-upper", false, "exclude uppercase letters")
	noDigits := flag.Bool("no-digits", false, "exclude digits")
	noSpecial := flag.Bool("no-special", false, "exclude special characters")
	noAmbiguous := flag.Bool("no-ambiguous", false, "exclude ambiguous characters (0OIl1)")

	flag.Parse()

	opts := Options{
		Length:      *length,
		NoUpper:     *noUpper,
		NoDigits:    *noDigits,
		NoSpecial:   *noSpecial,
		Special:     *special,
		NoAmbiguous: *noAmbiguous,
	}

	pw, err := Generate(opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(pw)
}
