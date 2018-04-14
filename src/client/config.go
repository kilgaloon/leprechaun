package client

import (
	"../log"
	"gopkg.in/ini.v1"
	"os"
)

// Config values
type Config struct {
	errorLog    string
	infoLog     string
	recipesPath string
}

func readConfig(path string) *Config {
	cfg, err := ini.Load(path)
	if err != nil {
		panic("Failed to load ini file")
		os.Exit(1)
	}

	c := &Config{}
	c.errorLog = cfg.Section("").Key("error_log").String()
	c.infoLog = cfg.Section("").Key("info_log").String()
	c.recipesPath = cfg.Section("").Key("recipes_path").String()

	variables := cfg.Section("variables").Keys()
	for _, variable := range variables {
		CurrentContext.DefineVar(variable.Name(), variable.String())
	}

	log.Logger.ErrorLog = c.errorLog
	log.Logger.InfoLog = c.infoLog

	return c
}
