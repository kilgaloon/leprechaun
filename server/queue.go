package server

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"strings"

	"github.com/kilgaloon/leprechaun/client"
	"github.com/kilgaloon/leprechaun/log"
	"github.com/kilgaloon/leprechaun/recipe"
)

// Queue stack for pulling out recipes
type Queue struct {
	Stack map[string]recipe.Recipe
}

// BuildQueue takes all recipes and put them in queue
func (server *Server) BuildQueue() {
	q := Queue{}
	q.Stack = make(map[string]recipe.Recipe)

	files, err := ioutil.ReadDir(server.Config.recipesPath)
	if err != nil {
		server.Logs.Error("%s", err)
	}

	for _, file := range files {
		fullFilepath := server.Config.recipesPath + "/" + file.Name()
		recipe := recipe.Build(fullFilepath)

		// recipes that needs to be pushed to queue
		// needs to be schedule by definition
		if recipe.Definition == "hook" {
			q.Stack[recipe.ID] = recipe
		}

	}

	server.Queue = q
}

// ProcessQueue queue
func (server *Server) ProcessQueue(id string) {
	recipe := server.Queue.Stack[id]

	log.Logger.Info("%s file is in progress... \n", recipe.Name)

	for index, step := range recipe.Steps {
		log.Logger.Info("Recipe %s Step %d is in progress... \n", recipe.Name, (index + 1))
		// replace variables
		step = client.CurrentContext.Transpile(step)

		parts := strings.Fields(step)
		parts = parts[1:]

		cmd := exec.Command("bash", "-c", step)

		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr

		err := cmd.Run()
		if err != nil {
			log.Logger.Info("Recipe %s Step %d failed to start. Reason: %s \n", recipe.Name, (index + 1), stderr.String())
		}

		log.Logger.Info("Recipe %s Step %d finished... \n\n", recipe.Name, (index + 1))

	}
}
