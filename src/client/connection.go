package client

import (
	"time"
	"../log"
	"io/ioutil"
)

// Client settings and configurations
type Client struct {
	Config *Config
	Logs log.Logs
}

// Start runs server
func Start(iniPath *string) {
	var client = &Client{}
	// load configurations for server
	client.Config = readConfig(*iniPath)
	client.Logs = log.CreateLogs(client.Config.errorLog)

	files, err := ioutil.ReadDir(client.Config.recipesPath)
    if err != nil {
        client.Logs.Error("%s", err)
	}

	q := BuildQueue(client.Config.recipesPath, files)

	for {
		go ProcessQueue(&q, client)

		time.Sleep(60 * time.Second)
	}
	
}