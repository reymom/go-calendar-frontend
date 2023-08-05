package tasks

import (
	"time"

	"github.com/reymom/go-calendar-tutorial/pkg/model"
)

type commonViewFields struct {
	Mode    filterModeId
	Year    int
	IsToday bool
}

type filterModeId uint8

const (
	filterModeIdDay filterModeId = iota
	filterModeIdWeek
	filterModeIdMonth
	filterModeIdYear
)

type viewFilter struct {
	Year  int
	Month time.Month
	Week  int
	Day   int
}

func (f viewFilter) GetWeekDay() time.Weekday {
	return time.Date(f.Year, f.Month, f.Day, 0, 0, 0, 0, time.UTC).Weekday()
}

type weeklyTasks struct {
	week       int
	dailyTasks map[time.Weekday][]model.Task
}

type weekDays struct {
	week   int
	nMonth int
	days   []int
}
