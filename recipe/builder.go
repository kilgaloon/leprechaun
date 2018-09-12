package recipe

import (
	"io/ioutil"
	"log"
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
func Build(file string) Recipe {
	r := Recipe{}

	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("Unable to open recipe: " + file)
	}

	error := yaml.Unmarshal(data, &r)
	if error != nil {
		log.Fatalf("Unable to unmarshal yaml: %s", error)
	}

	switch r.Definition {
	case "schedule":
		r.StartAt = recipe.ScheduleToTime(r.Schedule)
	}

	return r
}
