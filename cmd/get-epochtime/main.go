package main

import (
	"fmt"

	"github.com/curusarn/resh/internal/epochtime"
)

// Small utility to get epochtime in millisecond precision
// Doesn't check arguments
// Exits with status 1 on error
func main() {
	fmt.Printf("%s", epochtime.Now())
}
