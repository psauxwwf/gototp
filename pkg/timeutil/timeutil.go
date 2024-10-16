package timeutil

import (
	"fmt"
	"time"
)

const (
	timeFormat = "02.01.2006 15:04:05"
)

func Format(t time.Time) string {
	return t.Format(timeFormat)
}

func NowS() string {
	return Format(time.Now())
}

func Parse(s string) (time.Time, error) {
	location, err := time.LoadLocation("Local")
	if err != nil {
		return time.Time{}, fmt.Errorf("error loading location: %w", err)
	}
	t, err := time.ParseInLocation(timeFormat, s, location)
	if err != nil {
		return time.Time{}, fmt.Errorf("error parsing time: %w", err)
	}
	return t, nil
}
