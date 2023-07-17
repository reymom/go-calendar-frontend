package routes

import (
	"html/template"
	"net/http"

	"github.com/reymom/go-calendar-frontend/internal/config"
	"github.com/reymom/go-calendar-frontend/internal/routes/tasks"
	"github.com/reymom/go-calendar-tutorial/pkg/psql"
)

func GenerateRoutes(conf *config.Config) (http.Handler, error) {
	mux := http.NewServeMux()
	mux.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./www/public/images"))))
	mux.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./www/public/css"))))

	t := template.Must(template.New("html-tmpl").ParseGlob("www/views/pages/*/*.html"))
	_, e := t.ParseGlob("www/views/pages/*/*.html")
	if e != nil {
		return nil, e
	}

	tasksDao, e := psql.NewTaskDao(&psql.TaskDaoConfig{
		ConnectionStringRead:  conf.CalendarConnectionStringRead,
		ConnectionStringWrite: conf.CalendarConnectionStringWrite,
		MaxReadConnections:    2,
	})
	if e != nil {
		return nil, e
	}

	e = tasks.NewRouteGenerator(t, tasksDao).GenerateRoutes(mux)
	if e != nil {
		return nil, e
	}

	return mux, nil
}
