package time

import (
	"math"
	"time"
)

var asiaJakartaTz, _ = time.LoadLocation("Asia/Jakarta")

func Now() time.Time {
	return time.Now().In(asiaJakartaTz)
}

func Date(year int, month time.Month, day, hour, min, sec, nsec int) time.Time {
	return time.Date(year, month, day, hour, min, sec, nsec, asiaJakartaTz)
}

func Parse(layout, value string) (time.Time, error) {
	return time.ParseInLocation(layout, value, asiaJakartaTz)
}

func UnixMilli(msec int64) time.Time {
	return time.UnixMilli(msec).In(asiaJakartaTz)
}

func InLocAsiaJakarta(t time.Time) time.Time {
	return t.In(asiaJakartaTz)
}

func GetTimeForDisplay(value uint64) string {
	if value != 0 {
		return UnixMilli(int64(value)).Format(time.RFC3339)
	}
	return "-"
}

func GetBackOffDuration(attempt int) time.Duration {
	return time.Duration(math.Pow10(attempt-1)) * time.Second
}

func GetStartingTimeForDate(t time.Time) time.Time {
	year, month, day := t.Date()
	return Date(year, month, day, 0, 0, 0, 0)
}

func GetEndingTimeForDate(t time.Time) time.Time {
	year, month, day := t.Date()
	dateTime := Date(year, month, day, 0, 0, 0, 0)
	return dateTime.AddDate(0, 0, 1).Add(-1 * time.Millisecond)
}
