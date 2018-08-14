package client

import (
	"github.com/kilgaloon/leprechaun/log"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"io/ioutil"
	"time"
)

// Client settings and configurations
type Client struct {
	Config *Config
	Logs   log.Logs
}

// Start runs server
func Start(iniPath *string) {
	var client = &Client{}
	// load configurations for server
	client.Config = readConfig(*iniPath)

	files, err := ioutil.ReadDir(client.Config.recipesPath)
	if err != nil {
		client.Logs.Error("%s", err)
	}

	q := BuildQueue(client.Config.recipesPath, files)

	// watch for new recipes
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		client.Logs.Error("Failed to create watcher")
	}

	defer watcher.Close()

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Create == fsnotify.Create {
					AddToQueue(&q.Stack, event.Name)
				}
			case err := <-watcher.Errors:
				fmt.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(client.Config.recipesPath)
	if err != nil {
		fmt.Println(err)
	}

	for {

		go ProcessQueue(&q, client)

		time.Sleep(60 * time.Second)
	}

}
