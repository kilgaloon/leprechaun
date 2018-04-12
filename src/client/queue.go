package client

import (
	"../recipe"
	"bytes"
	"os"
	"os/exec"
	"strings"
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

// ProcessQueue queue
func ProcessQueue(queue *Queue, client *Client) {
	for index, r := range queue.Stack {
		recipe := &queue.Stack[index]

		if IsLocked(r.Name, client) {
			continue
		}

		if recipe.StartIn == 0 {
			if LockProcess(r.Name, client) {
				client.Logs.Info("%s file is in progress... \n", r.Name)

				for index, step := range r.Steps {
					client.Logs.Info("Recipe %s Step %d is in progress... \n", r.Name, (index + 1))
					// replace variables
					step = CurrentContext.Transpile(step)

					parts := strings.Fields(step)
					head := parts[0]
					parts = parts[1:]
					// if is internal command of Rainbow
					if len(head) >= 7 && head[0:7] == "rainbow" {
						err := Resolve(head, parts)
						if err != nil {
							client.Logs.Info("Recipe %s failed on step %d. Reason: %s \n", r.Name, (index + 1), err)
						}
					} else {
						cmd := exec.Command(head, parts...)

						var out bytes.Buffer
						var stderr bytes.Buffer
						cmd.Stdout = &out
						cmd.Stderr = &stderr

						err := cmd.Run()
						if err != nil {
							client.Logs.Info("Recipe %s Step %d failed to start. Reason: %s \n", r.Name, (index + 1), stderr.String())
						}

					}

					client.Logs.Info("Recipe %s Step %d finished... \n\n", r.Name, (index + 1))
					RemoveLock(r.Name, client)
				}

				recipe.StartIn = r.WorkEvery
			} else {
				client.Logs.Info("Failed to set lock on %s recipe", r.Name)
			}
		} else {
			client.Logs.Info("%s recipe will run in %d minutes \n\n", r.Name, recipe.StartIn)

			recipe.StartIn--
		}
	}
}
