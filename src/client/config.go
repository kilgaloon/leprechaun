package client

import (
	"os"
	"gopkg.in/ini.v1"
	"log"
)

// Config values
type Config struct {
	host string
	port string
	errorLog string
	recipesPath string
	commandsPath string
}

func readConfig(path string) (*Config) {
	cfg, err := ini.Load(path)
	if err != nil {
		log.Print("Failed to load ini file");
		os.Exit(1);
	}

	c := &Config{};
	c.host = cfg.Section("").Key("host").String()
	c.port = cfg.Section("").Key("port").String()
	c.errorLog = cfg.Section("").Key("error_log").String()
	c.recipesPath = cfg.Section("").Key("recipes_path").String()
	c.commandsPath = cfg.Section("").Key("commands_path").String()

	variables := cfg.Section("variables").Keys()
	for _, variable := range variables {
		CurrentContext.DefineVar(variable.Name(), variable.String())
	}

	return c
}