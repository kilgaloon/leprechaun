package server

import (
	"os"
	"strings"

	"github.com/kilgaloon/leprechaun/client"
	"github.com/kilgaloon/leprechaun/log"
	"gopkg.in/ini.v1"
)

// Config values
type Config struct {
	errorLog    string
	infoLog     string
	port        string
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
	c.port = cfg.Section("").Key("port").String()
	c.recipesPath = cfg.Section("").Key("recipes_path").String()

	variables := cfg.Section("variables").Keys()
	for _, variable := range variables {
		client.CurrentContext.DefineVar(variable.Name(), variable.String())
	}
	// insert environment variables in our context
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		client.CurrentContext.DefineVar(pair[0], pair[1])
	}

	log.Logger.ErrorLog = c.errorLog
	log.Logger.InfoLog = c.infoLog

	return c
}
