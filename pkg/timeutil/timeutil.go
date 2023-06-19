package timeutil

import (
	"time"

	"github.com/Songmu/flextime"
)

var (
	JST = time.FixedZone("Asia/Tokyo", 9*60*60)
)

func Now() time.Time {
	// TZ environment variable is set, but also set in code.
	return flextime.Now().In(JST)
}
