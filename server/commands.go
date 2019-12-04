package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/getsentry/raven-go"
	"github.com/kilgaloon/leprechaun/api"
	"github.com/kilgaloon/leprechaun/daemon"
)

// this section is used for command responders

func (server Server) cmdinfo(w http.ResponseWriter, r *http.Request) {
	recipeQueueNum := strconv.Itoa(len(server.Pool.Stack))

	resp := api.InfoResponse{
		Status:         server.GetStatus().String(),
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

func (server Server) cmdpause(w http.ResponseWriter, r *http.Request) {
	resp := &api.MessageResponse{}

	server.Pause()
	if server.GetStatus() == daemon.Paused {
		w.WriteHeader(http.StatusOK)
		resp.Message = "Server paused"
	}

	j, err := json.Marshal(resp)
	if err != nil {
		log.Fatal(err)
	}

	w.Write(j)
}

func (server Server) cmdstart(w http.ResponseWriter, r *http.Request) {
	resp := &api.MessageResponse{}

	if server.GetStatus() == daemon.Started {
		w.WriteHeader(http.StatusExpectationFailed)
		resp.Message = "Server already started"

		return
	}

	go server.Start()
	if server.GetStatus() == daemon.Started {
		w.WriteHeader(http.StatusOK)
		resp.Message = "Server started"
	}

	j, err := json.Marshal(resp)
	if err != nil {
		raven.CaptureError(err, nil)
		log.Fatal(err)
	}

	w.Write(j)
}

func (server Server) cmdstop(w http.ResponseWriter, r *http.Request) {
	resp := &api.MessageResponse{}

	server.Stop()
	if server.GetStatus() == daemon.Stopped {
		w.WriteHeader(http.StatusOK)
		resp.Message = "Server stopped"
	}

	j, err := json.Marshal(resp)
	if err != nil {
		raven.CaptureError(err, nil)
		log.Fatal(err)
	}

	w.Write(j)
}
