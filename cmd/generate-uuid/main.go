package main

import (
	"fmt"
	"os"

	"github.com/google/uuid"
)

// Small utility to generate UUID's using google/uuid golang package
// Doesn't check arguments
// Exits with status 1 on error
func main() {
	rnd, err := uuid.NewRandom()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: could not get new random source: %v", err)
		os.Exit(1)
	}
	id := rnd.String()
	if id == "" {
		fmt.Fprintf(os.Stderr, "ERROR: got invalid UUID from package")
		os.Exit(1)
	}
	// No newline
	fmt.Print(id)
}
