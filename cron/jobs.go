package cron

import (
	"io/ioutil"

	"github.com/kilgaloon/leprechaun/recipe"
)

func (c *Cron) buildJobs() {
	files, err := ioutil.ReadDir(c.GetConfig().GetRecipesPath())
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		fullFilepath := c.GetConfig().GetRecipesPath() + "/" + file.Name()
		recipe, err := recipe.Build(fullFilepath)
		if err != nil {
			c.GetLogs().Error(err.Error())
		}
		// recipes that needs to be pushed to queue
		// needs to be schedule by definition
		if recipe.Definition == "cron" {
			c.Service.AddFunc(recipe.Pattern, func() {
				worker, err := c.CreateWorker(&recipe)
				if err == nil {
					worker.Run()
				}
			})
		}
	}
}
