package timeutil

import (
	"time"

	"github.com/Songmu/flextime"
)

var (
	JST = time.FixedZone("Asia/Tokyo", 9*60*60)

	DURATION_MONTH = time.Hour * 24 * 30
)

func Now() time.Time {
	// TZ environment variable is set, but also set in code.
	return flextime.Now().In(JST)
}

func NewDate(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, JST)
}

func Max(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	} else {
		return b
	}
}
