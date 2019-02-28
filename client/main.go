package client

import (
	"net/http"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/kilgaloon/leprechaun/agent"
	"github.com/kilgaloon/leprechaun/daemon"

	"github.com/kilgaloon/leprechaun/config"
)

// Client settings and configurations
type Client struct {
	Name string
	*agent.Default
	Queue
}

// New create client as a service
func (client *Client) New(name string, cfg *config.AgentConfig, debug bool) daemon.Service {
	a := agent.New(name, cfg, debug)
	c := &Client{
		name,
		a,
		Queue{},
	}

	return c
}

// GetName returns agent name
func (client Client) GetName() string {
	return client.Name
}

// Start client+
func (client *Client) Start() {
	// build queue
	client.BuildQueue()

	// watch for new recipes
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic("Failed to create watcher")
	}

	defer watcher.Close()
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Create == fsnotify.Create {
					client.AddToQueue(&client.Queue.Stack, event.Name)
				}
			case err := <-watcher.Errors:
				client.Error("error: %s", err)
			}
		}
	}()

	err = watcher.Add(client.GetConfig().GetRecipesPath())
	if err != nil {
		client.Error("%s", err)
	}

	client.Info("Scheduler ready")

	for {
		go client.ProcessQueue()
		client.Info("Scheduler TICK")
		time.Sleep(60 * time.Second)
	}

}

// RegisterAPIHandles to be used in socket communication
// If you want to takeover default commands from agent
// call DefaultCommands from Agent which is same command
func (client *Client) RegisterAPIHandles() map[string]func(w http.ResponseWriter, r *http.Request) {
	cmds := make(map[string]func(w http.ResponseWriter, r *http.Request))

	cmds["info"] = client.clientInfo
	cmds["stop"] = client.cmdstop

	return cmds
}

// Check does client is working on something
// decide this in which status client resides
func (client *Client) isWorking() bool {
	// check does .lock file exists
	if _, err := os.Stat(client.GetConfig().GetLockFile()); os.IsNotExist(err) {
		return false
	}

	return true
}
