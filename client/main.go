package client

import (
	"net/http"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/getsentry/raven-go"
	"github.com/kilgaloon/leprechaun/agent"
	"github.com/kilgaloon/leprechaun/daemon"

	"github.com/kilgaloon/leprechaun/config"
)

// Agent holds client instance
var Agent *Client

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

	Agent = c

	return c
}

// GetName returns agent name
func (client *Client) GetName() string {
	return client.Name
}

// Start client
func (client *Client) Start() {
	// build queue
	client.BuildQueue()

	// watch for new recipes
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		raven.CaptureError(err, nil)
		panic("Failed to create watcher")
	}

	defer watcher.Close()
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Create == fsnotify.Create {
					client.AddToQueue(event.Name)
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
	client.SetStatus(daemon.Started)

	for {
		if client.GetStatus() == daemon.Stopped {
			continue
		}

		go client.ProcessQueue()
		client.Info("Scheduler TICK")
		time.Sleep(60 * time.Second)
	}

}

// Stop client
func (client *Client) Stop() {
	client.Lock()
	// reset queue
	client.Queue.Stack = client.Queue.Stack[:0]
	client.Unlock()
	// set service status to stopped
	client.SetStatus(daemon.Stopped)
}

// RegisterAPIHandles to be used in http communication
func (client *Client) RegisterAPIHandles() map[string]func(w http.ResponseWriter, r *http.Request) {
	cmds := make(map[string]func(w http.ResponseWriter, r *http.Request))

	cmds["info"] = client.cmdinfo
	cmds["stop"] = client.cmdstop
	cmds["start"] = client.cmdstart
	cmds["pause"] = client.cmdpause

	return cmds
}
