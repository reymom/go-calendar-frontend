package config

import (
	"encoding/json"
	"os"
)

const ConfigJsonName = "calendar.json"

var (
	Version   = "UNKNOWN"
	BuildDate = "UNKNOWN"
)

type Config struct {
	CalendarConnectionStringRead  string `json:"calendarConnectionStringRead"`
	CalendarConnectionStringWrite string `json:"calendarConnectionStringWrite"`
}

func GenerateConfig(path string) (*Config, error) {
	b, e := os.ReadFile(path)
	if e != nil {
		return nil, e
	}
	var ret Config
	e = json.Unmarshal(b, &ret)
	if e != nil {
		return nil, e
	}
	return &ret, nil
}
