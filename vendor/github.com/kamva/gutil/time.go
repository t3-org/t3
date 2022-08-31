package gutil

import (
	"fmt"
	"time"
)

func FormatDateStuckTogether(t time.Time) string {
	return fmt.Sprintf("%d%d%d%d%d%d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
}
