package client

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/kilgaloon/leprechaun/api"
	"github.com/kilgaloon/leprechaun/daemon"
)

// this section is used for command responders

func (client *Client) cmdinfo(w http.ResponseWriter, r *http.Request) {
	recipeQueueNum := strconv.Itoa(len(client.Queue.Stack))

	resp := api.InfoResponse{
		Status:         client.GetStatus().String(),
		RecipesInQueue: recipeQueueNum,
	}

	w.WriteHeader(http.StatusOK)

	j, err := json.Marshal(resp)
	if err != nil {
		log.Fatal(err)
	}

	w.Write(j)

	return
}

func (client *Client) cmdpause(w http.ResponseWriter, r *http.Request) {
	var resp struct {
		Message string
	}

	client.Pause()
	if client.GetStatus() == daemon.Paused {
		w.WriteHeader(http.StatusOK)
		resp.Message = "Client paused"
	} else {
		w.WriteHeader(http.StatusExpectationFailed)
		resp.Message = "Client failed to pause"
	}

	j, err := json.Marshal(resp)
	if err != nil {
		log.Fatal(err)
	}

	w.Write(j)
}

func (client *Client) cmdstart(w http.ResponseWriter, r *http.Request) {
	resp := api.MessageResponse{}

	if client.GetStatus() == daemon.Started {
		w.WriteHeader(http.StatusExpectationFailed)
		resp.Message = "Client already started"

		return
	}

	go client.Start()
	if client.GetStatus() == daemon.Started {
		w.WriteHeader(http.StatusOK)
		resp.Message = "Client started"
	} else {
		w.WriteHeader(http.StatusExpectationFailed)
		resp.Message = "Client failed to start"
	}

	j, err := json.Marshal(resp)
	if err != nil {
		log.Fatal(err)
	}

	w.Write(j)
}

func (client *Client) cmdstop(w http.ResponseWriter, r *http.Request) {
	resp := api.MessageResponse{}

	client.Stop()
	if client.GetStatus() == daemon.Stopped {
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
