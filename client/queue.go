package client

import (
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/kilgaloon/leprechaun/event"

	"github.com/kilgaloon/leprechaun/recipe"
	schedule "github.com/kilgaloon/leprechaun/recipe/schedule"
)

// Queue stack for pulling out recipes
type Queue struct {
	Stack []recipe.Recipe
}

// BuildQueue takes all recipes and put them in queue
func (client *Client) BuildQueue() {
	q := Queue{}

	files, err := ioutil.ReadDir(client.Config.RecipesPath)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		fullFilepath := client.Config.RecipesPath + "/" + file.Name()
		recipe := recipe.Build(fullFilepath)

		// recipes that needs to be pushed to queue
		// needs to be schedule by definition
		if recipe.Definition == "schedule" {
			q.Stack = append(q.Stack, recipe)
		}

	}

	client.Queue = q
}

// AddToQueue takes freshly created recipes and add them to queue
func (client Client) AddToQueue(stack *[]recipe.Recipe, path string) {
	if filepath.Ext(path) == ".yml" {
		r := recipe.Build(path)
		if r.Definition == "schedule" {
			*stack = append(*stack, recipe.Build(path))
		}
	}
}

// ProcessQueue queue
func (client *Client) ProcessQueue() {
	now := time.Now()
	compare := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, time.UTC)

	for index := range client.Queue.Stack {
		go func(r *recipe.Recipe) {
			if compare.Equal(r.StartAt) {
				// lock mutex
				client.mu.Lock()
				// create worker
				worker, err := client.Workers.CreateWorker(r.Name)
				// unlock mutex
				client.mu.Unlock()

				if err != nil {
					switch err {
					case client.Workers.Errors.AllowedWorkersReached:
						client.Logs.Info("%s", err)
						go client.ProcessRecipe(r)
					case client.Workers.Errors.StillActive:
						client.Logs.Info("Worker with NAME: %s is still active", r.Name)
					}
					// move this worker to queue and work on it when next worker space is available
					go client.ProcessRecipe(r)
					client.Logs.Info("%s", err)
					return
				}

				event.EventHandler.Dispatch("client:lock")
				client.Logs.Info("%s file is in progress... \n", r.Name)
				// worker takeover steps and works on then
				worker.Run(r.Steps)
				// signal that worker is done
				// then proceed with unlock
				event.EventHandler.Dispatch("client:unlock")
				// schedule recipe for next execution
				r.StartAt = schedule.ScheduleToTime(r.Schedule)
			}
		}(&client.Queue.Stack[index])
	}
}

// ProcessRecipe takes specific recipe and process it
func (client *Client) ProcessRecipe(r *recipe.Recipe) {
	// lock mutex
	client.mu.Lock()
	// create worker
	worker, err := client.Workers.CreateWorker(r.Name)
	// unlock mutex
	client.mu.Unlock()

	if err != nil {
		switch err {
		case client.Workers.Errors.AllowedWorkersReached:
			// move this worker to queue and work on it when next worker space is available
			time.Sleep(time.Duration(client.Config.RetryRecipeAfter) * time.Second)
			client.Logs.Info("%s, retrying in %d s ...", err, client.Config.RetryRecipeAfter)
			go client.ProcessRecipe(r)
		case client.Workers.Errors.StillActive:
			client.Logs.Info("Worker with NAME: %s is still active", r.Name)
		}

		return
	}

	event.EventHandler.Dispatch("client:lock")
	client.Logs.Info("%s file is in progress... \n", r.Name)
	// worker takeover steps and works on then
	worker.Run(r.Steps)
	//remove lock on client
	event.EventHandler.Dispatch("client:unlock")
}
