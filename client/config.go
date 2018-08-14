package client

import (
	"github.com/kilgaloon/leprechaun/log"
	"gopkg.in/ini.v1"
	"os"
	"strings"
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
		panic(err)
	}

	c := &Config{}
	c.errorLog = cfg.Section("").Key("error_log").String()
	c.infoLog = cfg.Section("").Key("info_log").String()
	c.recipesPath = cfg.Section("").Key("recipes_path").String()

	variables := cfg.Section("variables").Keys()
	for _, variable := range variables {
		CurrentContext.DefineVar(variable.Name(), variable.String())
	}
	// insert environment variables in our context
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		CurrentContext.DefineVar(pair[0], pair[1])
	}

	log.Logger.ErrorLog = c.errorLog
	log.Logger.InfoLog = c.infoLog

	return c
}
