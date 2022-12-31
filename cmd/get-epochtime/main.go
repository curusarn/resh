package main

import (
	"fmt"
	"time"
)

// Small utility to get epochtime in millisecond precision
// Doesn't check arguments
// Exits with status 1 on error
func main() {
	fmt.Printf("%s", timeToEpochTime(time.Now()))
}

func timeToEpochTime(t time.Time) string {
	return fmt.Sprintf("%.2f", float64(t.UnixMilli())/1000)
}
