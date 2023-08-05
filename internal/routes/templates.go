package routes

import (
	"html/template"

	"github.com/rs/zerolog/log"
)

func newFuncMap() template.FuncMap {
	return template.FuncMap{
		"dict":    passDictInTemplate,
		"FmtWeek": fmtWeek,
	}
}

func passDictInTemplate(values ...interface{}) map[string]interface{} {
	if len(values)%2 != 0 {
		log.Err(nil).Msg("invalid dict call")
	}
	dict := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			log.Err(nil).Msg("dict keys must be strings")
		}
		dict[key] = values[i+1]
	}
	return dict
}

func fmtWeek(isoWeek int) string {
	return ""
}
