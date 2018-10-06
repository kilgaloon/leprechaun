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

	files, err := ioutil.ReadDir(client.Agent.GetConfig().GetRecipesPath())
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		fullFilepath := client.Agent.GetConfig().GetRecipesPath() + "/" + file.Name()
		recipe, err := recipe.Build(fullFilepath)
		if err != nil {
			client.Agent.GetLogs().Error(err.Error())
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
func (client Client) AddToQueue(stack *[]recipe.Recipe, path string) {
	if filepath.Ext(path) == ".yml" {
		r, err := recipe.Build(path)
		if err != nil {
			client.Agent.GetLogs().Error(err.Error())
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
		go func(r *recipe.Recipe) {
			if compare.Equal(r.StartAt) {
				// lock mutex
				client.Agent.GetMutex().Lock()
				// create worker
				worker, err := client.Agent.GetWorkers().CreateWorker(r.Name)
				// unlock mutex
				client.Agent.GetMutex().Unlock()

				if err != nil {
					switch err {
					case client.Agent.GetWorkers().Errors.AllowedWorkersReached:
						client.Agent.GetLogs().Info("%s", err)
						go client.ProcessRecipe(r)
					case client.Agent.GetWorkers().Errors.StillActive:
						client.Agent.GetLogs().Info("Worker with NAME: %s is still active", r.Name)
					}
					// move this worker to queue and work on it when next worker space is available
					go client.ProcessRecipe(r)
					client.Agent.GetLogs().Info("%s", err)
					return
				}

				event.EventHandler.Dispatch("client:lock")
				client.Agent.GetLogs().Info("%s file is in progress... \n", r.Name)
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
	client.Agent.GetMutex().Lock()
	// create worker
	worker, err := client.Agent.GetWorkers().CreateWorker(r.Name)
	// unlock mutex
	client.Agent.GetMutex().Unlock()

	if err != nil {
		switch err {
		case client.Agent.GetWorkers().Errors.AllowedWorkersReached:
			// move this worker to queue and work on it when next worker space is available
			time.Sleep(time.Duration(client.Agent.GetConfig().RetryRecipeAfter) * time.Second)
			client.Agent.GetLogs().Info("%s, retrying in %d s ...", err, client.Agent.GetConfig().RetryRecipeAfter)
			go client.ProcessRecipe(r)
		case client.Agent.GetWorkers().Errors.StillActive:
			client.Agent.GetLogs().Info("Worker with NAME: %s is still active", r.Name)
		}

		return
	}

	event.EventHandler.Dispatch("client:lock")
	client.Agent.GetLogs().Info("%s file is in progress... \n", r.Name)
	// worker takeover steps and works on then
	worker.Run(r.Steps)
	//remove lock on client
	event.EventHandler.Dispatch("client:unlock")
}
