package tasks

import (
	"time"

	"github.com/reymom/go-calendar-frontend/internal/routes/helpers"
)

func getFirstDayOfISOWeek(year int, week int) (*time.Time, error) {
	date := firstDayOfISOWeek(year, week, time.UTC)

	// sanity check
	isoYear, isoWeek := date.ISOWeek()
	if year != isoYear {
		return nil, helpers.ErrIncorrectIsoYear
	}
	if week != isoWeek {
		return nil, helpers.ErrIncorrectIsoWeek
	}

	return &date, nil
}

func firstDayOfISOWeek(year int, week int, timezone *time.Location) time.Time {
	date := time.Date(year, 0, 0, 0, 0, 0, 0, timezone)
	isoYear, isoWeek := date.ISOWeek()
	for date.Weekday() != time.Monday { // iterate back to Monday
		date = date.AddDate(0, 0, -1)
		isoYear, isoWeek = date.ISOWeek()
	}
	for isoYear < year { // iterate forward to the first day of the first week
		date = date.AddDate(0, 0, 1)
		isoYear, isoWeek = date.ISOWeek()
	}
	for isoWeek < week { // iterate forward to the first day of the given week
		date = date.AddDate(0, 0, 1)
		isoYear, isoWeek = date.ISOWeek()
	}
	return date
}
