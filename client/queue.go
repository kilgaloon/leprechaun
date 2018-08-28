package client

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/kilgaloon/leprechaun/log"
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

	files, err := ioutil.ReadDir(client.Config.recipesPath)
	if err != nil {
		client.Logs.Error("%s", err)
	}

	for _, file := range files {
		fullFilepath := client.Config.recipesPath + "/" + file.Name()
		recipe := recipe.Build(fullFilepath)

		q.Stack = append(q.Stack, recipe)
	}

	client.Queue = q
}

// AddToQueue takes freshly created recipes and add them to queue
func (client Client) AddToQueue(stack *[]recipe.Recipe, path string) {
	if filepath.Ext(path) == ".yml" {
		*stack = append(*stack, recipe.Build(path))
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
				log.Logger.Info("%s file is in progress... \n", r.Name)

				client.LockChannel <- "lock"

				for index, step := range r.Steps {
					log.Logger.Info("Recipe %s Step %d is in progress... \n", r.Name, (index + 1))
					// replace variables
					step = CurrentContext.Transpile(step)

					parts := strings.Fields(step)
					parts = parts[1:]

					cmd := exec.Command("bash", "-c", step)

					var out bytes.Buffer
					var stderr bytes.Buffer
					cmd.Stdout = &out
					cmd.Stderr = &stderr

					err := cmd.Run()
					if err != nil {
						log.Logger.Info("Recipe %s Step %d failed to start. Reason: %s \n", r.Name, (index + 1), stderr.String())
					}

					log.Logger.Info("Recipe %s Step %d finished... \n\n", r.Name, (index + 1))
					RemoveLock(r.Name, client)

					client.LockChannel <- "unlock"
				}

				recipe.StartAt = schedule.ScheduleToTime(recipe.Schedule)

			} else {
				log.Logger.Info("Failed to set lock on %s recipe", r.Name)
			}
		}
	}
}
