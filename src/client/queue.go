package client

import (
	"../log"
	"../recipe"
	"bytes"
	"os"
	"os/exec"
	"strings"
	"path/filepath"
)

// Queue stack for pulling out recipes
type Queue struct {
	Stack []recipe.Recipe
}

// BuildQueue takes all recipes and put them in queue
func BuildQueue(path string, files []os.FileInfo) Queue {
	q := Queue{}

	for _, file := range files {
		fullFilepath := path + "/" + file.Name()
		recipe := recipe.Build(fullFilepath)

		q.Stack = append(q.Stack, recipe)
	}

	return q
}

// AddToQueue takes freshly created recipes and add them to queue
func AddToQueue(stack *[]recipe.Recipe, path string) {
	if filepath.Ext(path) == ".yml" {
		*stack = append(*stack, recipe.Build(path))
	}
}

// ProcessQueue queue
func ProcessQueue(queue *Queue, client *Client) {
	for index, r := range queue.Stack {
		recipe := &queue.Stack[index]

		if IsLocked(r.Name, client) {
			continue
		}

		if recipe.StartIn == 0 {
			if LockProcess(r.Name, client) {
				log.Logger.Info("%s file is in progress... \n", r.Name)

				for _, step := range r.Steps {
					log.Logger.Info("Recipe %s Step %d is in progress... \n", r.Name, (index + 1))
					// replace variables
					step = CurrentContext.Transpile(step)

					parts := strings.Fields(step)
					head := parts[0]
					parts = parts[1:]
					// if is internal command of Leprechaun
					if len(head) >= 7 && head[0:7] == "internal" {
						err := Resolve(head, parts)
						if err != nil {
							RemoveLock(r.Name, client)
							log.Logger.Info("Recipe %s failed on step %d. Reason: %s \n", r.Name, (index + 1), err)
						}
					} else {
						cmd := exec.Command("bash", "-c", step)

						var out bytes.Buffer
						var stderr bytes.Buffer
						cmd.Stdout = &out
						cmd.Stderr = &stderr

						err := cmd.Run()
						if err != nil {
							RemoveLock(r.Name, client)
							log.Logger.Info("Recipe %s Step %d failed to start. Reason: %s \n", r.Name, (index + 1), stderr.String())
							panic(err)
						}

					}

					log.Logger.Info("Recipe %s Step %d finished... \n\n", r.Name, (index + 1))
					RemoveLock(r.Name, client)
				}

				recipe.StartIn = r.WorkEvery
			} else {
				log.Logger.Info("Failed to set lock on %s recipe", r.Name)
				panic("Failed to set lock")
			}
		} else {
			log.Logger.Info("%s recipe will run in %d minutes \n\n", r.Name, recipe.StartIn)

			recipe.StartIn--
		}
	}
}
