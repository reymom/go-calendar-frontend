package tasks

import (
	"net/http"
	"strconv"

	"github.com/reymom/go-calendar-frontend/internal/routes/helpers"
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
		if req.FormValue("create") == "1" {
			r.handleCreationDayTask(w, req)
			break
		}
		r.handleDailyCalendarMode(w, req)
	case filterModeIdWeek:
		r.handleWeeklyCalendarMode(w, req)
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
