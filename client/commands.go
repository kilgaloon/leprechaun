package client

import (
	"runtime"
	"strconv"
	"time"

	"github.com/kilgaloon/leprechaun/socket"
)

// this section is used for command responders

// cmd: client info
func (c *Client) clientInfo(args ...string) ([][]string, error) {
	pid := strconv.Itoa(c.GetPID())
	num := strconv.Itoa(c.Workers.Size())
	recipeQueueNum := strconv.Itoa(len(c.Queue.Stack))

	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	alloc := strconv.FormatFloat(float64(mem.Alloc/1024)/1024, 'f', 2, 64)

	resp := [][]string{
		{"PID: " + pid},
		{"Config file: " + c.GetConfig().PathToConfig},
		{"Number of workers: " + num},
		{"Recipes in queue: " + recipeQueueNum},
		{"Memory Allocated: " + alloc + " MiB"},
	}

	return resp, nil
}

// cmd: client workers:list
func (c *Client) listWorkers(args ...string) ([][]string, error) {
	resp := [][]string{}

	if c.Workers.Size() < 1 {
		resp = [][]string{
			{"No workers currently working!"},
		}
	}

	for name, worker := range c.Workers.GetAll() {
		startedAt := worker.StartedAt.Format(time.UnixDate)
		resp = append(resp, []string{name, startedAt, worker.WorkingOn})
	}

	return resp, nil
}

// cms: client workers:kill {name}
func (c *Client) killWorker(args ...string) ([][]string, error) {
	resp := [][]string{}

	worker, err := c.Workers.GetByName(args[0])
	if err != nil {
		resp = [][]string{
			{err.Error()},
		}
	} else {
		worker.Kill()
		resp = [][]string{
			{"Worker killed"},
		}
	}

	return resp, nil
}

// RegisterCommandSocket returns Registrator
func (c *Client) RegisterCommandSocket() *socket.Registrator {
	r := socket.CreateRegistrator("client")

	// register commands
	r.Command("info", c.clientInfo)
	r.Command("workers:list", c.listWorkers)
	r.Command("workers:kill", c.killWorker)

	return r
}
