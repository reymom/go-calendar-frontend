package tasks

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/reymom/go-calendar-tutorial/pkg/model"
	"github.com/rs/zerolog/log"
)

func (r *Router) handleDailyCalendarMode(w http.ResponseWriter, req *http.Request) {
	day, _ := strconv.Atoi(req.FormValue("day"))
	month, _ := strconv.Atoi(req.FormValue("month"))
	year, _ := strconv.Atoi(req.FormValue("year"))

	filter := model.NewDailyFilter(uint(day), time.Month(month), uint(year))
	dailyTasks, e := r.tasksDao.ListTasks(req.Context(), filter)
	if e != nil {
		log.Err(e).Msgf("could not get tasks for day %d, month %s, year %d",
			day, time.Month(month).String(), year)
	}
	loc, _ := time.LoadLocation("Local")
	tasks := make(map[int][]model.Task, 3)
	for _, task := range dailyTasks {
		if _, ok := tasks[task.StartsAt.In(loc).Hour()]; !ok {
			tasks[task.StartsAt.In(loc).Hour()] = []model.Task{task}
			continue
		}
		tasks[task.StartsAt.In(loc).Hour()] = append(tasks[task.StartsAt.In(loc).Hour()], task)
	}

	var isToday bool
	if time.Now().Day() == day {
		isToday = true
	}

	_, offset := time.Now().In(loc).Zone()
	location := fmt.Sprintf("GMT%+d", offset/3600)
	m := struct {
		Tasks    map[int][]model.Task
		Common   commonViewFields
		Month    time.Month
		Day      int
		WeekDay  time.Weekday
		Location string
		Creation bool
	}{
		Tasks: tasks,
		Common: commonViewFields{
			Mode:    filterModeIdDay,
			Year:    year,
			IsToday: isToday,
		},
		Month:    time.Month(month),
		Day:      day,
		WeekDay:  time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local).Weekday(),
		Location: location,
		Creation: false,
	}
	w.Header().Set("Content-Type", "text/html")
	e = r.currentTemplate.Execute(w, m)
	if e != nil {
		log.Err(e).Msgf("error while executing the %s template", r.currentTemplate.Name())
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}
}
