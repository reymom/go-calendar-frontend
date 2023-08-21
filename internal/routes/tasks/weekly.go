package tasks

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/reymom/go-calendar-frontend/internal/routes/helpers"
	"github.com/reymom/go-calendar-tutorial/pkg/model"
	"github.com/rs/zerolog/log"
)

func (r *Router) handleWeeklyCalendarMode(w http.ResponseWriter, req *http.Request) {
	var (
		year          int
		week          int
		includesToday bool
	)

	now := time.Now()
	if req.FormValue("year") == "" || req.FormValue("week") == "" {
		year = now.Year()
		_, week = now.ISOWeek()
		includesToday = true
	} else {
		year, _ = strconv.Atoi(req.FormValue("year"))
		week, _ = strconv.Atoi(req.FormValue("week"))
		_, currentWeek := now.ISOWeek()
		if (now.Year() == year) && (currentWeek == week) {
			includesToday = true
		}
	}
	firstDate, _ := helpers.GetFirstDayOfISOWeek(year, week)
	month := firstDate.Month()

	filter := model.NewWeeklyFilter(uint(week), uint(year))
	weeklyTasks, e := r.tasksDao.ListTasks(context.Background(), filter)
	if e != nil {
		log.Err(e).Msgf("could not get tasks for iso week %d", week)
	}

	tasks := make(map[int][]model.Task, 7)
	for _, task := range weeklyTasks {
		if _, ok := tasks[task.StartsAt.Day()]; !ok {
			tasks[task.StartsAt.Day()] = []model.Task{task}
			continue
		}
		tasks[task.StartsAt.Day()] = append(tasks[task.StartsAt.Day()], task)
	}

	wDays, e := getWeekDaysForWeek(year, week)
	if e != nil {
		log.Err(e).Msg("could not get week days")
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}

	m := struct {
		Tasks         map[int][]model.Task
		Week          int
		WeekDays      []weekDays
		Common        commonViewFields
		IncludesToday bool
		Month         time.Month
	}{
		Tasks:    tasks,
		Week:     week,
		WeekDays: wDays,
		Common: commonViewFields{
			Mode:    filterModeIdWeek,
			Year:    year,
			IsToday: includesToday,
		},
		IncludesToday: includesToday,
		Month:         month,
	}
	w.Header().Set("Content-Type", "text/html")
	e = r.currentTemplate.Execute(w, m)
	if e != nil {
		log.Err(e).Msgf("error while executing the %s template", r.currentTemplate.Name())
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}
}

func getWeekDaysForWeek(year, week int) ([]weekDays, error) {
	firstDate, e := helpers.GetFirstDayOfISOWeek(year, week)
	if e != nil {
		return nil, e
	}

	month := firstDate.Month()
	nDays := make([]day, 0, 7)
	for i := 0; i < 7; i++ {
		nDays = append(nDays, day{
			Month: month,
			Day:   firstDate.Day(),
		})
		*firstDate = firstDate.AddDate(0, 0, 1)
	}

	return []weekDays{{
		Week: week,
		Days: nDays,
	}}, nil
}
