package client

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime"
	"strconv"
)

// this section is used for command responders

// cmd: client info
func (client *Client) clientInfo(w http.ResponseWriter, r *http.Request) {
	pid := strconv.Itoa(client.GetPID())
	recipeQueueNum := strconv.Itoa(len(client.Queue.Stack))

	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	alloc := strconv.FormatFloat(float64(mem.Alloc/1024)/1024, 'f', 2, 64)

	resp := struct {
		PID             string
		ConfigFile      string
		RecipesInQueue  string
		MemoryAllocated string
	}{
		PID:             pid,
		ConfigFile:      client.GetConfig().GetPath(),
		RecipesInQueue:  recipeQueueNum,
		MemoryAllocated: alloc + " MiB",
	}

	w.WriteHeader(http.StatusOK)

	j, err := json.Marshal(resp)
	if err != nil {
		log.Fatal(err)
	}

	w.Write(j)

	return
}

func (client *Client) cmdstop(w http.ResponseWriter, r *http.Request) {
	var resp struct {
		Message string
	}

	s := client.Stop()
	if s {
		w.WriteHeader(http.StatusOK)
		resp.Message = "Client stopped"
	} else {
		w.WriteHeader(http.StatusExpectationFailed)
		resp.Message = "Client failed to stop"
	}

	j, err := json.Marshal(resp)
	if err != nil {
		log.Fatal(err)
	}

	w.Write(j)
}
