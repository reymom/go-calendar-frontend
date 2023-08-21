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

func (r *Router) handleMonthlyCalendarMode(w http.ResponseWriter, req *http.Request) {
	var (
		year          int
		month         int
		includesToday bool
	)

	now := time.Now()
	if req.FormValue("year") == "" || req.FormValue("month") == "" {
		year = now.Year()
		month = int(now.Month())
		includesToday = true
	} else {
		year, _ = strconv.Atoi(req.FormValue("year"))
		month, _ = strconv.Atoi(req.FormValue("month"))
		if (now.Year() == year) && (now.Month() == time.Month(month)) {
			includesToday = true
		}
	}

	filter := model.NewMonthlyFilter(time.Month(month), uint(year))
	monthlyTasks, e := r.tasksDao.ListTasks(context.Background(), filter)
	if e != nil {
		log.Err(e).Msgf("could not get tasks for month %s", time.Month(month).String())
	}

	tasks := make(map[int][]model.Task, 31)
	for _, task := range monthlyTasks {
		if _, ok := tasks[task.StartsAt.Day()]; !ok {
			tasks[task.StartsAt.Day()] = []model.Task{task}
			continue
		}
		tasks[task.StartsAt.Day()] = append(tasks[task.StartsAt.Day()], task)
	}

	wDays, e := getWeekDaysForMonth(year, month)
	if e != nil {
		log.Err(e).Msg("could not get week days")
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}

	m := struct {
		Tasks         map[int][]model.Task
		WeekDays      []weekDays
		Common        commonViewFields
		IncludesToday bool
		Month         time.Month
	}{
		Tasks:    tasks,
		WeekDays: wDays,
		Common: commonViewFields{
			Mode:    filterModeIdMonth,
			Year:    year,
			IsToday: includesToday,
		},
		IncludesToday: includesToday,
		Month:         time.Month(month),
	}
	w.Header().Set("Content-Type", "text/html")
	e = r.currentTemplate.Execute(w, m)
	if e != nil {
		log.Err(e).Msgf("error while executing the %s template", r.currentTemplate.Name())
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}
}

func getWeekDaysForMonth(year, month int) ([]weekDays, error) {
	_, firstWeek := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC).ISOWeek()
	firstDate, e := helpers.GetFirstDayOfISOWeek(year, firstWeek)
	if e != nil {
		return nil, e
	}
	_, lastWeek := firstDate.AddDate(0, 1, -1).ISOWeek()
	lastDate, e := helpers.GetFirstDayOfISOWeek(year, lastWeek)
	if e != nil {
		return nil, e
	}
	*lastDate = lastDate.AddDate(0, 0, 6)

	wDays, nDays := make([]weekDays, 0, 5), make([]day, 0, 7)

	for firstDate.Before(lastDate.Add(time.Hour)) {
		nDays = append(nDays, day{
			Month: firstDate.Month(),
			Day:   firstDate.Day(),
		})
		if len(nDays) == 7 {
			_, isoWeek := firstDate.ISOWeek()
			wDays = append(wDays, weekDays{
				Week: isoWeek,
				Days: nDays,
			})
			nDays = make([]day, 0, 7)
		}

		*firstDate = firstDate.AddDate(0, 0, 1)
	}

	return wDays, nil
}
