package cron

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

func (c Cron) cmdinfo(w http.ResponseWriter, r *http.Request) {
	recipeQueueNum := strconv.Itoa(len(c.Service.Entries()))

	resp := api.InfoResponse{
		Status:         c.GetStatus().String(),
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

func (c Cron) cmdpause(w http.ResponseWriter, r *http.Request) {
	resp := api.MessageResponse{}

	c.Pause()
	if c.GetStatus() == daemon.Paused {
		w.WriteHeader(http.StatusOK)
		resp.Message = "Cron paused"
	}

	j, err := json.Marshal(resp)
	if err != nil {
		log.Fatal(err)
	}

	w.Write(j)
}

func (c Cron) cmdstart(w http.ResponseWriter, r *http.Request) {
	resp := api.MessageResponse{}

	if c.GetStatus() == daemon.Started {
		w.WriteHeader(http.StatusExpectationFailed)
		resp.Message = "Cron already started"

		return
	}

	go c.Start()
	// update methodology
	for {
		if c.GetStatus() == daemon.Started {
			w.WriteHeader(http.StatusOK)
			resp.Message = "Cron started"

			break
		}
	}

	j, err := json.Marshal(resp)
	if err != nil {
		raven.CaptureError(err, nil)
		log.Fatal(err)
	}

	w.Write(j)
}

func (c Cron) cmdstop(w http.ResponseWriter, r *http.Request) {
	resp := api.MessageResponse{}

	c.Stop()
	if c.GetStatus() == daemon.Stopped {
		w.WriteHeader(http.StatusOK)
		resp.Message = "Cron stopped"
	}

	j, err := json.Marshal(resp)
	if err != nil {
		raven.CaptureError(err, nil)
		log.Fatal(err)
	}

	w.Write(j)
}
