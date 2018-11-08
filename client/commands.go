package client

import (
	"io"
	"runtime"
	"strconv"
)

// this section is used for command responders

// cmd: client info
func (client *Client) clientInfo(r io.Writer, args ...string) ([][]string, error) {
	client.GetMutex().Lock()
	defer client.GetMutex().Unlock()

	pid := strconv.Itoa(client.GetPID())
	recipeQueueNum := strconv.Itoa(len(client.Queue.Stack))

	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	alloc := strconv.FormatFloat(float64(mem.Alloc/1024)/1024, 'f', 2, 64)

	resp := [][]string{
		{"PID: " + pid},
		{"Config file: " + client.GetConfig().GetPath()},
		{"Recipes in queue: " + recipeQueueNum},
		{"Memory Allocated: " + alloc + " MiB"},
	}

	return resp, nil
}
