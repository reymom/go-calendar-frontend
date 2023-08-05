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

const tasksListTemplateName = "taskList"

func (r *Router) tasksListHandler(w http.ResponseWriter, req *http.Request) {
	e := r.setCurrentTemplate(tasksListTemplateName)
	if e != nil {
		log.Err(e).Msgf("error while looking up the %s template", tasksListTemplateName)
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}

	mode, _ := strconv.Atoi(req.FormValue("mode"))
	switch filterModeId(mode) {
	case filterModeIdDay:
		r.handleDailyCalendarMode(w, req)
	case filterModeIdWeek:
		r.handleMonthlyCalendarMode(w, req)
	case filterModeIdMonth:
		r.handleMonthlyCalendarMode(w, req)
	case filterModeIdYear:
		r.handleYearlyCalendarMode(w, req)
	default:
		e = helpers.ErrUnknownFilterMode
		log.Err(e).Msgf("unknown filterModeId=%d", mode)
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}
}

func (r *Router) handleYearlyCalendarMode(w http.ResponseWriter, req *http.Request) {
	year, _ := strconv.Atoi(req.FormValue("year"))

	nTasks := make(map[time.Month]int, 12)
	for i := time.Month(1); i <= 12; i++ {
		filter := model.NewMonthlyFilter(i, uint(year))
		t, e := r.tasksDao.ListTasks(context.Background(), filter)
		if e != nil {
			log.Err(e).Msgf("could not get tasks for month %s", i.String())
		}
		nTasks[i] = len(t)
	}

	m := struct {
		NTasks map[time.Month]int
		Common commonViewFields
	}{
		NTasks: nTasks,
		Common: commonViewFields{
			Mode:    filterModeIdYear,
			Year:    year,
			IsToday: false,
		},
	}
	w.Header().Set("Content-Type", "text/html")
	e := r.currentTemplate.Execute(w, m)
	if e != nil {
		log.Err(e).Msgf("error while executing the %s template", r.currentTemplate.Name())
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}
}

func (r *Router) handleMonthlyCalendarMode(w http.ResponseWriter, req *http.Request) {
	year, _ := strconv.Atoi(req.FormValue("year"))
	month, _ := strconv.Atoi(req.FormValue("month"))

	var includesToday bool
	now := time.Now()
	if (now.Year() == year) && (now.Month() == time.Month(month)) {
		includesToday = true
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
		}
		tasks[task.StartsAt.Day()] = append(tasks[task.StartsAt.Day()], task)
	}

	wDays, e := getWeekDays(year, month)
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
			IsToday: false,
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

func (r *Router) handleDailyCalendarMode(w http.ResponseWriter, req *http.Request) {
	day, _ := strconv.Atoi(req.FormValue("day"))
	month, _ := strconv.Atoi(req.FormValue("month"))
	year, _ := strconv.Atoi(req.FormValue("year"))

	filter := model.NewDailyFilter(uint(day), time.Month(month), uint(year))
	dailyTasks, e := r.tasksDao.ListTasks(context.Background(), filter)
	if e != nil {
		log.Err(e).Msgf("could not get tasks for day %d, month %s, year %d", day, time.Month(month).String(), year)
	}

	var isToday bool
	if time.Now().Day() == day {
		isToday = true
	}
	m := struct {
		Tasks  []model.Task
		Hours  []int
		Common commonViewFields
		Month  time.Month
		Day    uint
	}{
		Tasks: dailyTasks,
		Hours: nil,
		Common: commonViewFields{
			Mode:    filterModeIdMonth,
			Year:    year,
			IsToday: isToday,
		},
		Month: time.Month(month),
		Day:   uint(day),
	}
	w.Header().Set("Content-Type", "text/html")
	e = r.currentTemplate.Execute(w, m)
	if e != nil {
		log.Err(e).Msgf("error while executing the %s template", r.currentTemplate.Name())
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}
}

func getWeekDays(year, month int) ([]weekDays, error) {
	_, isoWeek := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC).ISOWeek()
	firstWeek, e := getFirstDayOfISOWeek(year, isoWeek)
	if e != nil {
		return nil, e
	}
	_, isoWeek = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC).AddDate(0, 1, -1).ISOWeek()
	lastWeek, e := getFirstDayOfISOWeek(year, isoWeek)
	if e != nil {
		return nil, e
	}
	*lastWeek = lastWeek.AddDate(0, 0, 6)

	var lastDay bool
	nMonth, nWeek, i := 0, 0, 1
	wDays, nDays := make([]weekDays, 0, 5), make([]int, 0, 7)
	day := firstWeek.Day()
	for !lastDay {
		if day == 1 {
			nMonth++
		}

		nDays = append(nDays, day)
		if i == 7 {
			wDays = append(wDays, weekDays{
				week:   nWeek,
				nMonth: nMonth,
				days:   nDays,
			})
			nDays = make([]int, 0, 7)
			nWeek++
			i = 0
		}

		i++
		day++
		if lastWeek.Day() == day {
			lastDay = true
		}
	}

	return wDays, nil
}
