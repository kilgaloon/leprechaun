package cron

import (
	"io/ioutil"

	"github.com/kilgaloon/leprechaun/recipe"
)

func (c *Cron) buildJobs() {
	c.Info("Cron buildJobs started")
	files, err := ioutil.ReadDir(c.GetConfig().GetRecipesPath())
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		fullFilepath := c.GetConfig().GetRecipesPath() + "/" + file.Name()
		recipe, err := recipe.Build(fullFilepath)
		if err != nil {
			c.Error(err.Error())
		}
		// recipes that needs to be pushed to queue
		// needs to be schedule by definition
		if recipe.Definition == "cron" {
			c.Service.AddFunc(recipe.Pattern, func() {
				c.prepareAndRun(recipe)
			})
		}
	}

	c.Info("Cron buildJobs finished")
}

func (c *Cron) prepareAndRun(r recipe.Recipe) {
	worker, err := c.CreateWorker(r)

	if err == nil {
		c.PushToStack(worker)
		worker.Run()
	}
}
