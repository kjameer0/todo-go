package main

import (
	"fmt"
	"time"
)

func monthDayYear(t time.Time) string {
	y, m, d := t.Date()
	return fmt.Sprintf("%v %v %v %v", time.Weekday(d).String(), m.String()[:3], d, y)
}
