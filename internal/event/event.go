package event

import (
	"strings"
	"time"
)

type Event struct {
	Id          int       `json:"id"`
	Description string    `json:"description"`
	Name        string    `json:"name"`
	Date        time.Time `json:"date"`
}

func (e *Event) Validate() (problems map[string]string) {
	problems = map[string]string{}

	if strings.TrimSpace(e.Name) == "" {
		problems["name"] = "name cannot be empty"
	}

	if strings.TrimSpace(e.Description) == "" {
		problems["description"] = "description cannot be empty"
	}

	if e.Date.IsZero() {
		problems["date"] = "date cannot be empty"
	}

	return
}
