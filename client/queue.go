package client

import (
	"bytes"
	"io/ioutil"
	"os/exec"
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

	for index, r := range client.Queue.Stack {
		recipe := &client.Queue.Stack[index]

		if IsLocked(r.Name, client) {
			continue
		}

		if compare.Equal(recipe.StartAt) {
			if LockProcess(r.Name, client) {
				client.Logs.Info("%s file is in progress... \n", r.Name)

				event.EventHandler.Dispatch("client:lock")

				for index, step := range r.Steps {
					client.Logs.Info("Recipe %s Step %d is in progress... \n", r.Name, (index + 1))
					// replace variables
					step = client.Context.Transpile(step)

					cmd := exec.Command("bash", "-c", step)

					var out bytes.Buffer
					var stderr bytes.Buffer
					cmd.Stdout = &out
					cmd.Stderr = &stderr

					err := cmd.Run()
					if err != nil {
						client.Logs.Info("Recipe %s Step %d failed to start. Reason: %s \n", r.Name, (index + 1), stderr.String())
					}

					client.Logs.Info("Recipe %s Step %d finished... \n\n", r.Name, (index + 1))
					RemoveLock(r.Name, client)

					event.EventHandler.Dispatch("client:unlock")
				}

				recipe.StartAt = schedule.ScheduleToTime(recipe.Schedule)

			} else {
				client.Logs.Info("Failed to set lock on %s recipe", r.Name)
			}
		}
	}
}
