package tasks

import (
	"context"
	"net/http"
	"time"

	"github.com/reymom/go-calendar-frontend/internal/routes/helpers"
	"github.com/reymom/go-calendar-tutorial/pkg/model"
	"github.com/rs/zerolog/log"
)

const tasksListTemplateName = "taskList"

func (r *Router) tasksListHandler(w http.ResponseWriter, req *http.Request) {
	template := r.templates.Lookup(tasksListTemplateName)
	if template == nil {
		log.Err(helpers.ErrTemplateEmpty).Msgf("error while looking up the %s template", tasksListTemplateName)
		http.Error(w, helpers.ErrTemplateEmpty.Error(), http.StatusInternalServerError)
		return
	}

	var now = time.Now()
	tasks, e := r.tasksDao.ListTasks(context.Background(), model.NewMonthlyFilter(now.Month(), uint(now.Year())))
	if e != nil {
		log.Err(e).Msg("could not get tasks")
	}
	m := struct {
		Tasks []model.Task
	}{
		Tasks: tasks,
	}
	w.Header().Set("Content-Type", "text/html")
	e = template.Execute(w, m)
	if e != nil {
		log.Err(e).Msgf("error while executing the %s template", tasksListTemplateName)
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}
}
