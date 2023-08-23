package tasks

import (
	"html/template"
	"net/http"

	"github.com/reymom/go-calendar-frontend/internal/routes/helpers"
	"github.com/reymom/go-calendar-tutorial/pkg/model"
	"github.com/urfave/negroni/v3"
)

const TaskPrefixURL string = "/calendar/tasks"

type Router struct {
	templates       *template.Template
	tasksDao        model.TasksDao
	currentTemplate *template.Template
}

func NewRouteGenerator(templates *template.Template, tasksDao model.TasksDao) *Router {
	return &Router{templates: templates, tasksDao: tasksDao}
}

func (r *Router) GenerateRoutes(mux *http.ServeMux) error {
	mux.Handle(TaskPrefixURL, negroni.New(negroni.Wrap(http.HandlerFunc(r.tasksListHandler))))
	mux.Handle(TaskPrefixURL+"/create", negroni.New(negroni.Wrap(http.HandlerFunc(r.submitTaskHandler))))
	mux.Handle(TaskPrefixURL+"/remove", negroni.New(negroni.Wrap(http.HandlerFunc(r.removeTaskHandler))))
	mux.Handle(TaskPrefixURL+"/complete", negroni.New(negroni.Wrap(http.HandlerFunc(r.setCompleteTaskHandler))))
	return nil
}

func (r *Router) setCurrentTemplate(templateName string) error {
	r.currentTemplate = r.templates.Lookup(tasksListTemplateName)
	if r.currentTemplate == nil {
		return helpers.ErrTemplateEmpty
	}
	return nil
}
