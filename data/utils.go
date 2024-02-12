package data

import (
	"errors"
	"time"
)

// returns numbers of years
func GetWorkDuration(start, end time.Time) (int, error) {
	if start == (time.Time{}) {
		return 0, errors.New("Project/work start date is required.")
	} else if end == (time.Time{}) {
		return int(time.Since(start)) / 24 / 365, nil
	}
	return int(end.Sub(start)) / 24 / 365, nil
}
