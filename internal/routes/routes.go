package routes

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/reymom/go-calendar-frontend/internal/config"
	"github.com/reymom/go-calendar-frontend/internal/routes/tasks"
	"github.com/reymom/go-calendar-tutorial/pkg/psql"
)

func GenerateRoutes(conf *config.Config) (http.Handler, error) {
	t, e := template.New("html-tmpl").Funcs(newFuncMap()).ParseGlob("www/views/pages/*/*.html")
	if e != nil {
		return nil, e
	}

	tasksDao, e := psql.NewTaskDao(&psql.TaskDaoConfig{
		ConnectionStringRead:  conf.CalendarConnectionStringRead,
		ConnectionStringWrite: conf.CalendarConnectionStringWrite,
		MaxReadConnections:    2,
		MaxWriteConnections:   1,
	})
	if e != nil {
		return nil, e
	}

	mux := http.NewServeMux()
	mux.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./www/public/images"))))
	mux.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./www/public/css"))))
	e = tasks.NewRouteGenerator(t, tasksDao).GenerateRoutes(mux)
	if e != nil {
		return nil, e
	}

	mux.Handle("/", http.HandlerFunc(rootHandler))
	return mux, nil
}

func rootHandler(w http.ResponseWriter, req *http.Request) {
	now := time.Now()
	url := fmt.Sprintf("%s?mode=3&month=%d&year=%d", tasks.TaskPrefixURL, now.Month(), now.Year())
	http.Redirect(w, req, url, http.StatusSeeOther)
}
