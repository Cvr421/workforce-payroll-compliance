package clock

import (
	"time"
)

var nowFunc = time.Now

func Now() time.Time {
	return nowFunc()
}

func SetNow(f func() time.Time) {
	nowFunc = f
}

func ResetNow() {
	nowFunc = time.Now
}

func Today() time.Time {
	return Now().Truncate(24 * time.Hour)
}