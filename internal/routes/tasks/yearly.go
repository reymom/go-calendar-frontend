package tasks

import (
	"net/http"
	"strconv"
	"time"

	"github.com/reymom/go-calendar-tutorial/pkg/model"
	"github.com/rs/zerolog/log"
)

func (r *Router) handleYearlyCalendarMode(w http.ResponseWriter, req *http.Request) {
	year, _ := strconv.Atoi(req.FormValue("year"))

	nTasks := make(map[time.Month]int, 12)
	for i := time.Month(1); i <= 12; i++ {
		filter := model.NewMonthlyFilter(i, uint(year))
		t, e := r.tasksDao.ListTasks(req.Context(), filter)
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
