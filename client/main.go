package client

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/kilgaloon/leprechaun/agent"

	"github.com/fsnotify/fsnotify"
	"github.com/kilgaloon/leprechaun/api"
	"github.com/kilgaloon/leprechaun/config"
	"github.com/kilgaloon/leprechaun/event"
)

// Agent holds instance of Client
var Agent *Client

// Client settings and configurations
type Client struct {
	*agent.Default
	Queue
}

// New create client
// Creating new agent will enable usage of Agent variable globally for packages
// that use this package
func New(name string, cfg *config.AgentConfig) *Client {
	client := &Client{
		agent.New(name, cfg),
		Queue{},
	}

	Agent = client

	return client
}

// Start client
func (client *Client) Start() {
	// remove hanging .lock file
	os.Remove(client.GetConfig().GetLockFile())
	// SetPID of client
	client.SetPID()
	// build queue
	client.Lock()
	client.BuildQueue()
	client.Unlock()

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
				client.GetLogs().Error("error: %s", err)
			}
		}
	}()

	err = watcher.Add(client.GetConfig().GetRecipesPath())
	if err != nil {
		fmt.Println(err)
	}

	event.EventHandler.Dispatch("client:ready")
	// register client to command socket
	go api.New(client.GetConfig().GetCommandSocket()).Register(client)

	for {
		go client.ProcessQueue()
		time.Sleep(60 * time.Second)
	}

}

// RegisterCommands to be used in socket communication
// If you want to takeover default commands from agent
// call DefaultCommands from Agent which is same command
func (client *Client) RegisterCommands() map[string]api.Command {
	cmds := make(map[string]api.Command)

	cmds["info"] = api.Command{
		Closure: client.clientInfo,
		Definition: api.Definition{
			Text:  "Display some basic info about running client",
			Usage: "client info",
		},
	}

	// this function merge both maps and inject default commands from agent
	return client.DefaultCommands(cmds)
}

// SetPID sets current PID of client
func (client *Client) SetPID() {
	f, err := os.OpenFile(client.GetConfig().GetPIDFile(), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic("Failed to start client, can't save PID, reason: " + err.Error())
	}

	client.PID = os.Getpid()
	pid := strconv.Itoa(client.GetPID())
	_, err = f.WriteString(pid)
	if err != nil {
		panic("Failed to start client, can't save PID")
	}
}

// GetPID gets current PID of client
func (client *Client) GetPID() int {
	return client.PID
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

// Lock client to busy state
func (client *Client) Lock() {
	_, err := os.OpenFile(client.GetConfig().GetLockFile(), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic("Failed to lock client in busy state")
	}
}

// Unlock client to busy state
func (client *Client) Unlock() {
	os.Remove(client.GetConfig().GetLockFile())
}

// Stop client
func (client *Client) Stop() os.Signal {
	var answer string
	forceQuit := false
	quit := false

	fmt.Fprintf(client, "Are you sure?(y/N): ")
	fmt.Fscanf(client, "%s", &answer)

	if client.isWorking() && strings.ToUpper(answer) == "Y" {
		answer = ""
		// if user doesn't choose to force quit we will wait for process, otherwise KILL IT
		fmt.Fprintf(client.GetStdout(), "Client is working on something in the background. Force quit? (y/N)")
		fmt.Fscanf(client.GetStdin(), "%s", &answer)

		if strings.ToUpper(answer) == "Y" {
			forceQuit = true
		}
	} else if !client.isWorking() && strings.ToUpper(answer) == "Y" {
		quit = true
	}

	pid := client.GetPID()
	process, err := os.FindProcess(pid)
	if err != nil {
		client.GetLogs().Error("Can't find process with that PID. %s", err)
	}

	// shutdown gracefully
	if quit {
		state, err := process.Wait()
		client.GetLogs().Info("Stopping Leprechaun, please wait...")

		if err == nil {
			if state.Exited() {
				client.Unlock()
				return syscall.SIGTERM
			}
		} else {
			forceQuit = true
		}
	}

	// force quite and terminate everything
	if forceQuit {
		killed := process.Kill()
		if killed != nil {
			client.GetLogs().Error("Can't kill process with that PID. %s", killed)
		} else {
			client.Unlock()
			return syscall.SIGTERM
		}
	}

	return os.Interrupt
}

func init() {
	// subscribe to events for this package
	event.EventHandler.Subscribe("client:lock", func() {
		Agent.Lock()
	})

	event.EventHandler.Subscribe("client:unlock", func() {
		Agent.Unlock()
	})
}
