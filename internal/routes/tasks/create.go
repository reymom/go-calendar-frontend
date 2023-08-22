package tasks

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/reymom/go-calendar-tutorial/pkg/model"
	"github.com/rs/zerolog/log"
)

func (r *Router) handleCreationDayTask(w http.ResponseWriter, req *http.Request) {
	day, _ := strconv.Atoi(req.FormValue("day"))
	month, _ := strconv.Atoi(req.FormValue("month"))
	year, _ := strconv.Atoi(req.FormValue("year"))
	colorId, _ := strconv.Atoi(req.FormValue("color"))

	var isToday bool
	if time.Now().Day() == day {
		isToday = true
	}
	loc, _ := time.LoadLocation("Local")
	_, offset := time.Now().In(loc).Zone()
	location := fmt.Sprintf("GMT%+d", offset/3600)

	m := struct {
		Common   commonViewFields
		Month    time.Month
		Day      int
		WeekDay  time.Weekday
		Location string
		Creation bool
		Color    model.ColorId
	}{
		Common: commonViewFields{
			Mode:    filterModeIdDay,
			Year:    year,
			IsToday: isToday,
		},
		Month:    time.Month(month),
		Day:      day,
		WeekDay:  time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local).Weekday(),
		Location: location,
		Creation: true,
		Color:    model.ColorId(colorId),
	}
	w.Header().Set("Content-Type", "text/html")
	e := r.currentTemplate.Execute(w, m)
	if e != nil {
		log.Err(e).Msgf("error while executing the %s template", r.currentTemplate.Name())
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}
}

func (r *Router) taskSubmitTaskHandler(w http.ResponseWriter, req *http.Request) {
	var (
		dayStr   = req.FormValue("day")
		monthStr = req.FormValue("month")
		yearStr  = req.FormValue("year")
	)
	redirectUrl := fmt.Sprintf("%s?mode=1&day=%s&month=%s&year=%s&creation=0",
		TaskPrefixURL, dayStr, monthStr, yearStr)

	color, _ := strconv.Atoi(req.FormValue("color"))
	colorId := model.ColorId(color)

	day, _ := strconv.Atoi(dayStr)
	month, _ := strconv.Atoi(monthStr)
	year, _ := strconv.Atoi(yearStr)

	startHour := req.FormValue("startsAt")
	hourMinute := strings.Split(startHour, ":")
	hour, _ := strconv.Atoi(hourMinute[0])
	minute, _ := strconv.Atoi(hourMinute[1])
	loc, _ := time.LoadLocation("Local")
	startsAt := time.Date(year, time.Month(month), day, hour, minute, 0, 0, time.UTC).In(loc)

	finishHour := req.FormValue("finishesAt")
	hourMinute = strings.Split(finishHour, ":")
	hour, _ = strconv.Atoi(hourMinute[0])
	minute, _ = strconv.Atoi(hourMinute[1])
	finishesAt := time.Date(year, time.Month(month), day, hour, minute, 0, 0, time.UTC).In(loc)

	task, e := r.tasksDao.CreateTask(req.Context(), model.AddableTask{
		Name:        req.FormValue("name"),
		Description: req.FormValue("description"),
		StartsAt:    startsAt,
		FinishesAt:  finishesAt,
		Priority:    model.PriorityTypeIdMiddle,
		Color:       colorId,
	})
	if e != nil {
		log.Err(e).Msg("error creating task")
		http.Redirect(w, req, redirectUrl, http.StatusBadRequest)
		return
	}

	log.Warn().Msgf("task = %+v\n", task)

	http.Redirect(w, req, redirectUrl, http.StatusCreated)
}
