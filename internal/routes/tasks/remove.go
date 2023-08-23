package tasks

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/rs/zerolog/log"
)

func (r *Router) removeTaskHandler(w http.ResponseWriter, req *http.Request) {
	var (
		dayStr   = req.FormValue("day")
		monthStr = req.FormValue("month")
		yearStr  = req.FormValue("year")
		taskId   = req.FormValue("taskId")
	)
	redirectUrl := fmt.Sprintf("%s?mode=1&day=%s&month=%s&year=%s&creation=0",
		TaskPrefixURL, dayStr, monthStr, yearStr)

	e := r.tasksDao.RemoveTask(req.Context(), taskId)
	if e != nil {
		log.Err(e).Msgf("error removing task %s!", taskId)
	}

	http.Redirect(w, req, redirectUrl, http.StatusSeeOther)
}

func (r *Router) setCompleteTaskHandler(w http.ResponseWriter, req *http.Request) {
	var (
		dayStr   = req.FormValue("day")
		monthStr = req.FormValue("month")
		yearStr  = req.FormValue("year")
		taskId   = req.FormValue("taskId")
	)
	redirectUrl := fmt.Sprintf("%s?mode=1&day=%s&month=%s&year=%s&creation=0",
		TaskPrefixURL, dayStr, monthStr, yearStr)

	complete, _ := strconv.ParseBool(req.FormValue("complete"))
	e := r.tasksDao.SetCompleted(req.Context(), taskId, complete)
	if e != nil {
		log.Err(e).Msgf("error setting task %s as complete=%t!", taskId, complete)
	}

	http.Redirect(w, req, redirectUrl, http.StatusSeeOther)
}
