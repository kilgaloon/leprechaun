package recipe

import (
	"io/ioutil"
	"sync"
	"time"

	"github.com/getsentry/raven-go"
	scheduler "github.com/kilgaloon/leprechaun/recipe/schedule"
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
	Pattern    string
	Err        error
	*sync.Mutex
}

// GetName returns name of recipe
func (r *Recipe) GetName() string {
	return r.Name
}

// SetStartAt modifies time when recipe will be started
func (r *Recipe) SetStartAt(t time.Time) {
	r.Lock()
	defer r.Unlock()

	r.StartAt = t
}

// GetStartAt get time when recipe will be started
func (r *Recipe) GetStartAt() time.Time {
	return r.StartAt
}

// GetSteps return array of steps that recipe needs to work on
func (r *Recipe) GetSteps() []string {
	return r.Steps
}

// Build recipe for use
func Build(file string) (Recipe, error) {
	r := Recipe{}
	r.Mutex = new(sync.Mutex)

	data, err := ioutil.ReadFile(file)
	if err != nil {
		raven.CaptureError(err, nil)
		return r, err
	}

	error := yaml.Unmarshal(data, &r)
	if error != nil {
		raven.CaptureError(err, nil)
		return r, err
	}

	switch r.Definition {
	case "schedule":
		r.StartAt = scheduler.ScheduleToTime(r.Schedule)
	}

	return r, nil
}
