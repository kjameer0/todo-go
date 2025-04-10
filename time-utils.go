package main

import (
	"fmt"
	"time"
)

func monthDayYear(t time.Time) string {
	y, m, d := t.Date()
	return fmt.Sprintf("%v %v %v %v", t.Weekday().String(), m.String()[:3], d, y)
}

func addDayToDate(t time.Time, days int) time.Time {
	if days <= 0 {
		days = 0
	}
	midnight := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
	aheadTime := midnight.Add(time.Hour * 24 * time.Duration(days))
	return aheadTime
}
