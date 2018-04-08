package recipe

import (
	"log"
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

// Recipe struct
type Recipe struct {
	Name string
	StartIn int
	WorkEvery int
	Steps []string
}

// Build recipe for use
func Build(file string) (Recipe) {
	r := Recipe{}

	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("Unable to open recipe: " + file)
	}

	error := yaml.Unmarshal(data, &r)
	if error != nil {
		log.Fatalf("Unable to unmarshal yaml: %s", error)
	}

	return r
}