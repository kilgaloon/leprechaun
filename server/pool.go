package server

import (
	"io/ioutil"
	"time"

	"github.com/kilgaloon/leprechaun/log"
	"github.com/kilgaloon/leprechaun/recipe"
)

// Pool stack for pulling out recipes
type Pool struct {
	Stack map[string]recipe.Recipe
}

// BuildPool takes all recipes and put them in pool
func (server *Server) BuildPool() {
	q := Pool{}
	q.Stack = make(map[string]recipe.Recipe)

	files, err := ioutil.ReadDir(server.Config.RecipesPath)
	if err != nil {
		server.Logs.Error("%s", err)
	}

	for _, file := range files {
		fullFilepath := server.Config.RecipesPath + "/" + file.Name()
		recipe := recipe.Build(fullFilepath)

		// recipes that needs to be pushed to pool
		// needs to be schedule by definition
		if recipe.Definition == "hook" {
			q.Stack[recipe.ID] = recipe
		}

	}

	server.Pool = q
}

// FindInPool Find recipe in pool and run it
func (server *Server) FindInPool(id string) {
	recipe := server.Pool.Stack[id]

	log.Logger.Info("%s file is in progress... \n", recipe.Name)

	// lock mutex
	server.mu.Lock()
	// create worker
	worker, err := server.Workers.CreateWorker(recipe.Name)
	// unlock mutex
	server.mu.Unlock()
	if err != nil {
		// move this worker to queue and retry to work on it
		go server.ProcessRecipe(recipe)
		server.Logs.Info("%s", err)
		return
	}

	worker.Run(recipe.Steps)
}

// ProcessRecipe takes specific recipe and process it
func (server *Server) ProcessRecipe(r recipe.Recipe) {
	recipe := r

	server.Logs.Info("%s file is in progress... \n", recipe.Name)
	// lock mutex
	server.mu.Lock()
	// create worker
	worker, err := server.Workers.CreateWorker(recipe.Name)
	// unlock mutex
	server.mu.Unlock()
	if err != nil {
		time.Sleep(time.Duration(server.Config.RetryRecipeAfter) * time.Second)
		server.Logs.Info("%s, retrying in %d s...", err, server.Config.RetryRecipeAfter)
		// move this worker to queue and work on it when next worker space is available
		go server.ProcessRecipe(recipe)
		return
	}

	worker.Run(recipe.Steps)
}
