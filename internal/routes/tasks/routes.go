package tasks

import (
	"html/template"
	"net/http"

	"github.com/reymom/go-calendar-frontend/internal/routes/helpers"
	"github.com/reymom/go-calendar-tutorial/pkg/model"
	"github.com/urfave/negroni/v3"
)

type Router struct {
	templates       *template.Template
	tasksDao        model.TasksDao
	currentTemplate *template.Template
}

func NewRouteGenerator(templates *template.Template, tasksDao model.TasksDao) *Router {
	return &Router{templates: templates, tasksDao: tasksDao}
}

func (r *Router) GenerateRoutes(mux *http.ServeMux) error {
	mux.Handle("/calendar/tasks", negroni.New(negroni.Wrap(http.HandlerFunc(r.tasksListHandler))))
	return nil
}

func (r *Router) setCurrentTemplate(templateName string) error {
	r.currentTemplate = r.templates.Lookup(tasksListTemplateName)
	if r.currentTemplate == nil {
		return helpers.ErrTemplateEmpty
	}
	return nil
}
