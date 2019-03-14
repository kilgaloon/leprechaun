package validate

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/kilgaloon/leprechaun/recipe"
)

// CheckRecipe does check on recipe to see is it in valid format
func CheckRecipe(path string) {
	r, err := recipe.Build(path)
	if err != nil {
		log.Fatalf("Invalid format of recipe: %s", err)
	}

	if r.Definition == "" {
		log.Fatal("Definition of recipe not provided")
	}

	jobName := filepath.Base(path)
	if r.GetName() == "" {
		log.Fatal("Name not specified")
	}

	if r.GetName()+".yml" != jobName {
		fmt.Print("It is good practice that you job name is same like file name {name}.yml\n")
	}

	if len(r.GetSteps()) < 1 {
		log.Fatal("Steps definition is either missing or in wrong format")
	}

	checkDefinition(r)

	fmt.Println("Everything good!")
}

func checkDefinition(r recipe.Recipe) {
	switch r.Definition {
	case "schedule":
		if r.Schedule["min"] == 0 && r.Schedule["hour"] == 0 && r.Schedule["day"] == 0 {
			fmt.Println("Your recipe will run on startup, and never again!")
		}

		if len(r.Schedule) < 1 {
			log.Fatal("Schedule definition is either missing or in wrong format")
		}

		for key := range r.Schedule {
			if key != "min" && key != "hour" && key != "day" {
				log.Fatalf("Unknown schedule key %s allowed (min, hour, day)", key)
			}
		}

	case "hook":
		if r.ID == "" {
			log.Fatal("ID of recipe not provided")
		}
	}
}
