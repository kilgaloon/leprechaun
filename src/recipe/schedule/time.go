package recipe

import (
	"time"
)

const (
	// Day defined how much day have hours
	Day = time.Hour * 24
)

// ScheduleToTime takes map when recipe needs to be scheduled
// and converts it in time
func ScheduleToTime(schedule map[string]int) time.Time {
	now := time.Now()
	min := schedule["min"]
	hour := schedule["hour"]
	day := schedule["day"]

	now = now.Add(time.Minute * time.Duration(min))
	now = now.Add(time.Hour * time.Duration(hour))
	now = now.Add(Day * time.Duration(day))

	return time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, time.UTC)
}
