package cron

import (
	"io/ioutil"
	"time"

	"github.com/kilgaloon/leprechaun/recipe"
)

func (c *Cron) buildJobs() {
	files, err := ioutil.ReadDir(c.Agent.GetConfig().GetRecipesPath())
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		fullFilepath := c.Agent.GetConfig().GetRecipesPath() + "/" + file.Name()
		recipe, err := recipe.Build(fullFilepath)
		if err != nil {
			c.Agent.GetLogs().Error(err.Error())
		}
		// recipes that needs to be pushed to queue
		// needs to be schedule by definition
		if recipe.Definition == "cron" {
			c.Service.AddFunc(recipe.Pattern, func() {
				// lock mutex
				c.Agent.GetMutex().Lock()
				// create worker
				worker, err := c.Agent.GetWorkers().CreateWorker(recipe.Name)
				// unlock mutex
				c.Agent.GetMutex().Unlock()

				if err != nil {
					switch err {
					case c.Agent.GetWorkers().Errors.AllowedWorkersReached:
						c.Agent.GetLogs().Info("%s", err)
						go c.processRecipe(&recipe)
					case c.Agent.GetWorkers().Errors.StillActive:
						c.Agent.GetLogs().Info("Worker with NAME: %s is still active", recipe.Name)
					}
					// move this worker to queue and work on it when next worker space is available
					go c.processRecipe(&recipe)
					c.Agent.GetLogs().Info("%s", err)
					return
				}

				worker.Run(recipe.Steps)
			})
		}
	}
}

// ProcessRecipe takes specific recipe and process it
func (c *Cron) processRecipe(r *recipe.Recipe) {
	// lock mutex
	c.Agent.GetMutex().Lock()
	// create worker
	worker, err := c.Agent.GetWorkers().CreateWorker(r.Name)
	// unlock mutex
	c.Agent.GetMutex().Unlock()

	if err != nil {
		switch err {
		case c.Agent.GetWorkers().Errors.AllowedWorkersReached:
			// move this worker to queue and work on it when next worker space is available
			time.Sleep(time.Duration(c.Agent.GetConfig().RetryRecipeAfter) * time.Second)
			c.Agent.GetLogs().Info("%s, retrying in %d s ...", err, c.Agent.GetConfig().RetryRecipeAfter)
			go c.processRecipe(r)
		case c.Agent.GetWorkers().Errors.StillActive:
			c.Agent.GetLogs().Info("Worker with NAME: %s is still active", r.Name)
		}

		return
	}

	c.Agent.GetLogs().Info("%s file is in progress... \n", r.Name)
	// worker takeover steps and works on then
	worker.Run(r.Steps)
	//remove lock on client
}
