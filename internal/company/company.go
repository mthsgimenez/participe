package company

import "strings"

type Company struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func (c *Company) Validate() (problems map[string]string) {
	if strings.TrimSpace(c.Name) == "" {
		problems = map[string]string{"name": "name must not be empty"}
		return
	}

	return nil
}
