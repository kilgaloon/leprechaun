package client

import (
	"os"
	"gopkg.in/ini.v1"
	"log"
)

// Config values
type Config struct {
	errorLog string
	recipesPath string
}

func readConfig(path string) (*Config) {
	cfg, err := ini.Load(path)
	if err != nil {
		log.Print("Failed to load ini file");
		os.Exit(1);
	}

	c := &Config{};
	c.errorLog = cfg.Section("").Key("error_log").String()
	c.recipesPath = cfg.Section("").Key("recipes_path").String()

	variables := cfg.Section("variables").Keys()
	for _, variable := range variables {
		CurrentContext.DefineVar(variable.Name(), variable.String())
	}

	return c
}