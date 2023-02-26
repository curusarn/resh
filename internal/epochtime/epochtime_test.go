package epochtime

import (
	"strconv"
	"testing"
	"time"
)

func TestConversion(t *testing.T) {
	epochTime := "1672702332.64"
	seconds, err := strconv.ParseFloat(epochTime, 64)
	if err != nil {
		t.Fatal("Test setup failed: Failed to convert constant")
	}
	if TimeToString(time.UnixMilli(int64(seconds*1000))) != epochTime {
		t.Fatal("EpochTime changed during conversion")
	}
}
