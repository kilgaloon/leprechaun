package server

import (
	"io/ioutil"

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

	server.GetMutex().Lock()
	worker, err := server.CreateWorker(recipe)
	server.GetMutex().Unlock()
	if err == nil {
		worker.Run()
	}
}
