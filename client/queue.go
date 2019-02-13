package client

import (
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/kilgaloon/leprechaun/recipe"
	schedule "github.com/kilgaloon/leprechaun/recipe/schedule"
)

// Queue stack for pulling out recipes
type Queue struct {
	Stack []recipe.Recipe
}

// BuildQueue takes all recipes and put them in queue
func (client *Client) BuildQueue() {
	client.GetMutex().Lock()
	defer client.GetMutex().Unlock()

	q := Queue{}

	files, err := ioutil.ReadDir(client.GetConfig().GetRecipesPath())
	if err != nil {
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

	client.Queue = q
}

// AddToQueue takes freshly created recipes and add them to queue
func (client *Client) AddToQueue(stack *[]recipe.Recipe, path string) {
	client.GetMutex().Lock()
	defer client.GetMutex().Unlock()

	if filepath.Ext(path) == ".yml" {
		r, err := recipe.Build(path)
		if err != nil {
			client.Error(err.Error())
		}

		if r.Definition == "schedule" {
			*stack = append(*stack, r)
		}
	}
}

// ProcessQueue queue
func (client *Client) ProcessQueue() {
	now := time.Now()
	compare := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, time.UTC)

	for index := range client.Queue.Stack {
		r := &client.Queue.Stack[index]
		// If recipe had some errors
		// don't run it again
		if r.Err != nil {
			continue
		}

		go func(r *recipe.Recipe) {
			client.GetMutex().Lock()
			defer client.GetMutex().Unlock()

			// if client is stopped reschedule recipe but don't run it
			if client.stopped {
				r.StartAt = schedule.ScheduleToTime(r.Schedule)
			} else {
				if compare.Equal(r.StartAt) {
					worker, err := client.CreateWorker(r)
					if err == nil {
						client.Lock()
						client.Info("%s file is in progress... \n", r.Name)
						// worker takeover steps and works on then
						worker.Run()
						// then proceed with unlock
						client.Unlock()
						// schedule recipe for next execution
						r.StartAt = schedule.ScheduleToTime(r.Schedule)
					}
				}
			}
		}(r)
	}
}
