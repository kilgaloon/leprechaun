package client

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/kilgaloon/leprechaun/log"
)

// Client settings and configurations
type Client struct {
	Config *Config
	Logs   log.Logs
}

// Create new client
func Create(iniPath *string) *Client {
	var client = &Client{}
	// load configurations for server
	client.Config = readConfig(*iniPath)

	return client
}

// Start client
func (client *Client) Start() {
	f, err := os.OpenFile(client.Config.PIDFile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic("Failed to start client, can't save PID")
	}

	pid := strconv.Itoa(os.Getpid())
	_, err = f.WriteString(pid)
	if err != nil {
		panic("Failed to start client, can't save PID")
	}

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

// GetPID gets current PID of client
func (client *Client) GetPID() int {
	PIDFile := client.Config.PIDFile

	data, err := ioutil.ReadFile(PIDFile)
	if err != nil {
		log.Logger.Error("Failed to read pid from .pid file. %s", err)
	}

	pid, err := strconv.Atoi(string(data))
	if err != nil {
		log.Logger.Error("Failed to parse pid from .pid file. %s", err)
	}

	return pid
}

// Stop client
func (client *Client) Stop() (os.Signal, bool) {
	pid := client.GetPID()
	process, err := os.FindProcess(pid)
	if err != nil {
		log.Logger.Error("Can't find process with that PID. %s", err)
	}

	state, err := process.Wait()
	log.Logger.Info("Stopping Leprechaun, please wait...")

	var killed error
	if err != nil {
		killed = process.Kill()
		if killed != nil {
			log.Logger.Error("Can't kill process with that PID. %s", err)
			os.Exit(3)
		} else {
			return syscall.SIGTERM, true
		}
	} else {
		if state.Exited() {
			return syscall.SIGTERM, true
		}
	}

	return os.Interrupt, false
}
