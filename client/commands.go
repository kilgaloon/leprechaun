package client

import (
	"runtime"
	"strconv"
)

// this section is used for command responders

// cmd: client info
func (client *Client) clientInfo(args ...string) ([][]string, error) {
	pid := strconv.Itoa(client.GetPID())
	num := strconv.Itoa(client.Agent.GetWorkers().Size())
	recipeQueueNum := strconv.Itoa(len(client.Queue.Stack))

	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	alloc := strconv.FormatFloat(float64(mem.Alloc/1024)/1024, 'f', 2, 64)

	resp := [][]string{
		{"PID: " + pid},
		{"Config file: " + client.Agent.GetConfig().GetPath()},
		{"Number of workers: " + num},
		{"Recipes in queue: " + recipeQueueNum},
		{"Memory Allocated: " + alloc + " MiB"},
	}

	return resp, nil
}
