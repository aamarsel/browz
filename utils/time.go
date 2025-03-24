package utils

import (
	"fmt"
	"time"
)

var daysOfWeek = map[time.Weekday]string{
	time.Monday:    "понедельник",
	time.Tuesday:   "вторник",
	time.Wednesday: "среда",
	time.Thursday:  "четверг",
	time.Friday:    "пятница",
	time.Saturday:  "суббота",
	time.Sunday:    "воскресенье",
}

func FormatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	if minutes == 0 {
		return fmt.Sprintf("%d ч", hours)
	}
	return fmt.Sprintf("%d ч %d мин", hours, minutes)
}

func GetDayOfWeek(date time.Time) string {
	return daysOfWeek[date.Weekday()]
}
