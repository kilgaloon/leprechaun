package server

import (
	"io/ioutil"

	"github.com/kilgaloon/leprechaun/recipe"
)

// Pool stack for pulling out recipes
type Pool struct {
	Stack map[string]*recipe.Recipe
}

// BuildPool takes all recipes and put them in pool
func (server *Server) BuildPool() {
	// mostly used for debug
	server.Info("BuildPool started")

	q := Pool{}
	q.Stack = make(map[string]*recipe.Recipe)

	files, err := ioutil.ReadDir(server.GetConfig().GetRecipesPath())
	if err != nil {
		server.Error("%s", err)
	}

	for _, file := range files {
		fullFilepath := server.GetConfig().GetRecipesPath() + "/" + file.Name()
		recipe, err := recipe.Build(fullFilepath)
		if err != nil {
			server.Error("%s", err)
		}
		// recipes that needs to be pushed to pool
		// needs to be schedule by definition
		if recipe.Definition == "hook" {
			q.Stack[recipe.ID] = &recipe
		}

	}

	server.Lock()
	server.Pool = q
	server.Unlock()

	// mostly used for debug
	server.Info("BuildPool finished")
}

// FindInPool Find recipe in pool and run it
// **TODO**: Rename this method to something more descriptive
func (server *Server) FindInPool(id string) {
	server.Lock()
	recipe := server.Pool.Stack[id]
	server.Unlock()

	// Recipe has some error, don't execute it
	if recipe.Err != nil {
		return
	}

	server.Info("%s file is in progress... \n", recipe.GetName())

	worker, err := server.CreateWorker(recipe)
	if err == nil {
		worker.Run()
	}
}
