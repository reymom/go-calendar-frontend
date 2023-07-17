package tasks

import (
	"html/template"
	"net/http"

	"github.com/reymom/go-calendar-tutorial/pkg/model"
	"github.com/urfave/negroni/v3"
)

type Router struct {
	templates *template.Template
	tasksDao  model.TasksDao
}

func NewRouteGenerator(templates *template.Template, tasksDao model.TasksDao) *Router {
	return &Router{templates: templates, tasksDao: tasksDao}
}

func (r *Router) GenerateRoutes(mux *http.ServeMux) error {
	mux.Handle("calendar/tasks", negroni.New(negroni.Wrap(http.HandlerFunc(r.tasksListHandler))))
	return nil
}
