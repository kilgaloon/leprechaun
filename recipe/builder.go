package recipe

import (
	"io/ioutil"
	"time"

	"github.com/kilgaloon/leprechaun/recipe/schedule"
	"gopkg.in/yaml.v2"
)

// Recipe struct
type Recipe struct {
	ID         string
	Name       string
	Definition string
	StartAt    time.Time
	Schedule   map[string]int
	Steps      []string
}

// Build recipe for use
func Build(file string) (Recipe, error) {
	r := Recipe{}

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return r, err
	}

	error := yaml.Unmarshal(data, &r)
	if error != nil {
		return r, err
	}

	switch r.Definition {
	case "schedule":
		r.StartAt = recipe.ScheduleToTime(r.Schedule)
	}

	return r, nil
}
