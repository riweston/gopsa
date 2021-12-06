package utils

import (
	"time"
)

func DateCalculator(today time.Time) (string, string) {
	mondayIterator := today
	sundayIterator := today
	nextMonth := time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, time.UTC).AddDate(0, 1, -1)

	for mondayIterator.Weekday().String() != "Monday" {
		mondayIterator = mondayIterator.AddDate(0, 0, -1)
	}
	for {
		if sundayIterator.Weekday().String() == "Sunday" {
			break
		}
		if sundayIterator == nextMonth {
			break
		}
		sundayIterator = sundayIterator.AddDate(0, 0, +1)
	}
	mondayIterator.Format("2006-01-02")
	return mondayIterator.Format("2006-01-02"), sundayIterator.Format("2006-01-02")
}
