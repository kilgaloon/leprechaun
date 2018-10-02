package client

import (
	"runtime"
	"strconv"

	"github.com/kilgaloon/leprechaun/api"
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
		{"Config file: " + c.GetConfig().GetPath()},
		{"Number of workers: " + num},
		{"Recipes in queue: " + recipeQueueNum},
		{"Memory Allocated: " + alloc + " MiB"},
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
func (c *Client) RegisterCommandSocket() *api.Registrator {
	r := api.CreateRegistrator(c)

	// register commands
	r.Command("info", c.clientInfo)

	return r
}
