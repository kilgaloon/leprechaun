package client

import (
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/getsentry/raven-go"
	"github.com/kilgaloon/leprechaun/daemon"
	"github.com/kilgaloon/leprechaun/recipe"
	schedule "github.com/kilgaloon/leprechaun/recipe/schedule"
)

// Queue stack for pulling out recipes
type Queue struct {
	Stack []recipe.Recipe
}

// BuildQueue takes all recipes and put them in queue
func (client *Client) BuildQueue() {
	client.Info("Scheduler BuildQueue started")
	q := &Queue{}

	files, err := ioutil.ReadDir(client.GetConfig().GetRecipesPath())
	if err != nil {
		raven.CaptureError(err, nil)
		panic(err)
	}

	for _, file := range files {
		fullFilepath := client.GetConfig().GetRecipesPath() + "/" + file.Name()
		recipe, err := recipe.Build(fullFilepath)
		if err != nil {
			client.Error(err.Error())
		}
		// recipes that needs to be pushed to queue
		// needs to be schedule by definition
		if recipe.Definition == "schedule" {
			q.Stack = append(q.Stack, recipe)
		}
	}

	client.Lock()
	client.Queue = q
	client.Unlock()

	client.Info("Scheduler BuildQueue finished")
}

// AddToQueue takes freshly created recipes and add them to queue
func (client *Client) AddToQueue(path string) {
	if filepath.Ext(path) == ".yml" {
		r, err := recipe.Build(path)
		if err != nil {
			client.Error(err.Error())
		}

		if r.Definition == "schedule" {
			client.Lock()
			client.Queue.Stack = append(client.Queue.Stack, r)
			client.Unlock()
		}
	}
}

// FindRecipe in queue
func (client *Client) FindRecipe(name string) *recipe.Recipe {
	client.Lock()
	defer client.Unlock()

	for _, r := range client.Queue.Stack {
		if r.GetName() == name {
			return &r
		}
	}

	return nil
}

// ProcessQueue queue
func (client *Client) ProcessQueue() {
	client.Lock()
	defer client.Unlock()

	now := time.Now()
	compare := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, time.UTC)

	for _, r := range client.Queue.Stack {
		// If recipe had some errors
		// don't run it again
		if r.Err != nil {
			continue
		}

		go func(r recipe.Recipe) {
			// if client is paused reschedule recipe but don't run it
			if client.GetStatus() == daemon.Paused {
				r.SetStartAt(schedule.ScheduleToTime(r.Schedule))
			} else {
				if compare.Equal(r.GetStartAt()) {
					worker, err := client.CreateWorker(r)
					if err == nil {
						client.Info("%s file is in progress... \n", r.GetName())
						// worker takeover steps and works on then
						worker.Run()
						// schedule recipe for next execution
						r.SetStartAt(schedule.ScheduleToTime(r.Schedule))
					}
				}
			}
		}(r)
	}
}
