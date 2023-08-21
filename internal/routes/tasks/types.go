package tasks

import "time"

type commonViewFields struct {
	Mode    filterModeId
	Year    int
	IsToday bool
}

type filterModeId uint8

const (
	filterModeIdNull filterModeId = iota
	filterModeIdDay
	filterModeIdWeek
	filterModeIdMonth
	filterModeIdYear
)

type weekDays struct {
	Week int
	Days []day
}

type day struct {
	Day   int
	Month time.Month
}
