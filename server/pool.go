package server

import (
	"io/ioutil"
	"time"

	"github.com/kilgaloon/leprechaun/log"
	"github.com/kilgaloon/leprechaun/recipe"
)

// Pool stack for pulling out recipes
type Pool struct {
	Stack map[string]*recipe.Recipe
}

// BuildPool takes all recipes and put them in pool
func (server *Server) BuildPool() {
	server.GetMutex().Lock()
	defer server.GetMutex().Unlock()

	q := Pool{}
	q.Stack = make(map[string]*recipe.Recipe)

	files, err := ioutil.ReadDir(server.GetConfig().GetRecipesPath())
	if err != nil {
		server.GetLogs().Error("%s", err)
	}

	for _, file := range files {
		fullFilepath := server.GetConfig().GetRecipesPath() + "/" + file.Name()
		recipe, err := recipe.Build(fullFilepath)
		if err != nil {
			server.GetLogs().Error(err.Error())
		}
		// recipes that needs to be pushed to pool
		// needs to be schedule by definition
		if recipe.Definition == "hook" {
			q.Stack[recipe.ID] = &recipe
		}

	}

	server.Pool = q
}

// FindInPool Find recipe in pool and run it
// **TODO**: Rename this method to something more descriptive
func (server Server) FindInPool(id string) {
	recipe := server.Pool.Stack[id]

	// Recipe has some error, don't execute it
	if recipe.Err != nil {
		return
	}

	log.Logger.Info("%s file is in progress... \n", recipe.Name)

	worker, err := server.CreateWorker(recipe)
	if err != nil {
		// move this worker to queue and retry to work on it
		go server.ProcessRecipe(recipe)
		server.GetLogs().Info("%s", err)
		return
	}

	worker.Run()
}

// ProcessRecipe takes specific recipe and process it
func (server *Server) ProcessRecipe(r *recipe.Recipe) {
	recipe := r

	// Recipe has some error, don't execute it
	if &recipe.Err != nil {
		return
	}

	server.GetLogs().Info("%s file is in progress... \n", recipe.Name)

	worker, err := server.CreateWorker(recipe)
	if err != nil {
		time.Sleep(time.Duration(server.GetConfig().RetryRecipeAfter) * time.Second)
		server.GetLogs().Info("%s, retrying in %d s...", err, server.GetConfig().RetryRecipeAfter)
		// move this worker to queue and work on it when next worker space is available
		go server.ProcessRecipe(recipe)
		return
	}

	worker.Run()
}
