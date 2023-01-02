package epochtime

import (
	"fmt"
	"time"
)

func TimeToString(t time.Time) string {
	return fmt.Sprintf("%.2f", float64(t.UnixMilli())/1000)
}

func Now() string {
	return TimeToString(time.Now())
}
