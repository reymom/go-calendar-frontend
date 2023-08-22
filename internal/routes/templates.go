package routes

import (
	"fmt"
	"html/template"
	"time"

	"github.com/reymom/go-calendar-frontend/internal/routes/helpers"
	"github.com/reymom/go-calendar-frontend/internal/routes/tasks"
	"github.com/reymom/go-calendar-tutorial/pkg/model"
	"github.com/rs/zerolog/log"
)

func newFuncMap() template.FuncMap {
	return template.FuncMap{
		"dict":               passDictInTemplate,
		"FmtWeek":            fmtWeek,
		"CheckHasTwoMonths":  checkHasTwoMonths,
		"CheckIsToday":       checkIsToday,
		"GetTaskColorClass":  getTaskColorClass,
		"FormatTaskInfo":     printTaskInfo,
		"GetYearURL":         getYearUrl,
		"GetCurrentYearURL":  getCurrentYearUrl,
		"IncreaseMonth":      increaseMonth,
		"GetMonthURL":        getMonthURL,
		"GetCurrentMonthURL": getCurrentMonthURL,
		"GetWeekURL":         getWeekURL,
		"GetCurrentWeekURL":  getCurrentWeekURL,
		"GetDayURL":          getDayURL,
		"GetCurrentDayURL":   getCurrentDayURL,
		"LoopHours":          loopHours,
		"LoopInts":           loopInts,
		"sprintf": func(s, i interface{}) string {
			return fmt.Sprintf("%v%v", s, i)
		},
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

func fmtWeek(isoWeek int, year int) string {
	const format string = "January 2"
	location, _ := time.LoadLocation("Local")
	firstDate, _ := helpers.GetFirstDayOfISOWeek(year, isoWeek)
	return fmt.Sprintf("(%s to %s)", firstDate.In(location).Format(format), firstDate.AddDate(0, 0, 6).In(location).Format(format))
}

func checkHasTwoMonths(isoWeek int, year int) bool {
	firstDate, _ := helpers.GetFirstDayOfISOWeek(year, isoWeek)
	return firstDate.Month() != firstDate.AddDate(0, 0, 7).Month()
}

func checkIsToday(day int) bool {
	return time.Now().Day() == day
}

func getTaskColorClass(color model.ColorId) string {
	switch color {
	case model.ColorIdRed:
		return "red-task"
	case model.ColorIdYellow:
		return "yellow-task"
	case model.ColorIdBlue:
		return "blue-task"
	case model.ColorIdOrange:
		return "orange-task"
	case model.ColorIdGreen:
		return "green-task"
	case model.ColorIdViolet:
		return "violet-task"
	case model.ColorIdCyan:
		return "cyan-task"
	case model.ColorIdBlack:
		return "black-task"
	case model.ColorIdWhite:
		return "white-task"
	default:
		return ""
	}
}

func printTaskInfo(starts time.Time, ends time.Time, description string) string {
	location, _ := time.LoadLocation("Local")
	return fmt.Sprintf("From %s \nTo %s \n %s",
		starts.In(location).Format("15:04"),
		ends.In(location).Format("15:04"),
		description,
	)
}

func getYearUrl(year, addYear int) string {
	date := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC).AddDate(addYear, 0, 0)
	return fmt.Sprintf("%s?mode=4&year=%d", tasks.TaskPrefixURL, date.Year())
}

func getCurrentYearUrl() string {
	return fmt.Sprintf("%s?mode=4&year=%d", tasks.TaskPrefixURL, time.Now().Year())
}

func increaseMonth(month time.Month) time.Month {
	if month == time.December {
		return time.January
	}
	return month + 1
}

func getMonthURL(month time.Month, year int, addMonth int) string {
	date := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC).AddDate(0, addMonth, 0)
	return fmt.Sprintf("%s?mode=3&month=%d&year=%d",
		tasks.TaskPrefixURL, date.Month(), date.Year())
}

func getCurrentMonthURL() string {
	date := time.Now()
	return fmt.Sprintf("%s?mode=3&month=%d&year=%d",
		tasks.TaskPrefixURL, date.Month(), date.Year())
}

func getWeekURL(week, year, addWeek int) string {
	outWeek := week + addWeek
	cond1 := (year == 2020 || year == 2026) && week == 53
	cond2 := !(year == 2020 || year == 2026) && week == 52

	if addWeek == 1 && (cond1 || cond2) {
		year += 1
		outWeek = 1
	}

	if week == 1 && addWeek == -1 {
		year -= 1
		if year == 2020 || year == 2026 {
			outWeek = 53
		} else {
			outWeek = 52
		}
	}
	return fmt.Sprintf("%s?mode=2&week=%d&year=%d", tasks.TaskPrefixURL, outWeek, year)
}

func getCurrentWeekURL() string {
	date := time.Now()
	_, week := date.ISOWeek()
	return fmt.Sprintf("%s?mode=2&week=%d&year=%d",
		tasks.TaskPrefixURL, week, date.Year())
}

func getDayURL(day int, month time.Month, year, addDay int) string {
	date := time.Date(year, month, day, 0, 0, 0, 0, time.UTC).AddDate(0, 0, addDay)
	return fmt.Sprintf("%s?mode=1&day=%d&month=%d&year=%d",
		tasks.TaskPrefixURL, date.Day(), date.Month(), date.Year())
}

func getCurrentDayURL() string {
	date := time.Now()
	return fmt.Sprintf("%s?mode=1&month=%d&year=%d&day=%d",
		tasks.TaskPrefixURL, date.Month(), date.Year(), date.Day())
}

func loopHours() []string {
	hours := make([]string, 0, 24)
	for i := 0; i < 24; i++ {
		hours = append(hours, fmt.Sprintf("%02d:00", i))
	}
	return hours
}

func loopInts(lessThan int) []int {
	list := make([]int, 0, lessThan)
	for i := 0; i < lessThan; i++ {
		list = append(list, i)
	}
	return list
}
