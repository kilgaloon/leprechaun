package client

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

func (client Client) cmdinfo(w http.ResponseWriter, r *http.Request) {
	recipeQueueNum := strconv.Itoa(len(client.Queue.Stack))

	resp := api.TableResponse{
		Header: []string{"Status", "Recipes in queue"},
		Columns: [][]string{
			{client.GetStatus().String(), recipeQueueNum},
		},
	}

	w.WriteHeader(http.StatusOK)

	j, err := json.Marshal(resp)
	if err != nil {
		raven.CaptureError(err, nil)
		log.Fatal(err)
	}

	w.Write(j)

	return
}

func (client Client) cmdpause(w http.ResponseWriter, r *http.Request) {
	resp := api.TableResponse{
		Header:  []string{"Message"},
		Columns: [][]string{},
	}

	client.Pause()
	if client.GetStatus() == daemon.Paused {
		w.WriteHeader(http.StatusOK)
		resp.Columns = append(resp.Columns, []string{"Paused"})
	}

	j, err := json.Marshal(resp)
	if err != nil {
		raven.CaptureError(err, nil)
		log.Fatal(err)
	}

	w.Write(j)
}

func (client Client) cmdstart(w http.ResponseWriter, r *http.Request) {
	resp := api.TableResponse{
		Header:  []string{"Message"},
		Columns: [][]string{},
	}

	if client.GetStatus() == daemon.Started {
		w.WriteHeader(http.StatusExpectationFailed)
		resp.Columns = append(resp.Columns, []string{"Client already started"})
	} else {
		go client.Start()

		w.WriteHeader(http.StatusOK)
		resp.Columns = append(resp.Columns, []string{"Client started"})
	}

	j, err := json.Marshal(resp)
	if err != nil {
		raven.CaptureError(err, nil)
		log.Fatal(err)
	}

	w.Write(j)
}

func (client Client) cmdstop(w http.ResponseWriter, r *http.Request) {
	resp := api.TableResponse{
		Header:  []string{"Message"},
		Columns: [][]string{},
	}

	if client.GetStatus() == daemon.Started {
		client.Stop()
		if client.GetStatus() == daemon.Stopped {
			w.WriteHeader(http.StatusOK)
			resp.Columns = append(resp.Columns, []string{"Client stopped"})
		}
	} else {
		w.WriteHeader(http.StatusExpectationFailed)
		resp.Columns = append(resp.Columns, []string{"Client already stopped"})
	}

	j, err := json.Marshal(resp)
	if err != nil {
		raven.CaptureError(err, nil)
		log.Fatal(err)
	}

	w.Write(j)
}
