package agent

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/getsentry/raven-go"
	"github.com/kilgaloon/leprechaun/api"
)

// WorkersList is default command for agents
func (d *Default) WorkersList(w http.ResponseWriter, r *http.Request) {
	resp := api.TableResponse{
		Header: []string{"Name", "Started at", "Working on"},
		Columns: [][]string{},
	}

	for name, worker := range d.GetAllWorkers() {
		startedAt := worker.StartedAt.Format(time.UnixDate)
		resp.Columns = append(resp.Columns, []string{name, startedAt, worker.WorkingOn})
	}

	w.WriteHeader(http.StatusOK)
	j, err := json.Marshal(resp)
	if err != nil {
		raven.CaptureError(err, nil)
		log.Fatal(err)
	}

	w.Write(j)
}

// KillWorker kills worker by provided name
func (d *Default) KillWorker(w http.ResponseWriter, r *http.Request) {
	resp := api.TableResponse{
		Header: []string{"message"},
		Columns: [][]string{},
	}

	worker, err := d.GetWorkerByName(r.URL.Query()["args"][0])
	if err != nil {
		resp.Columns = append(resp.Columns, []string{err.Error()})
	} else {
		worker.Kill()
		resp.Columns = append(resp.Columns, []string{"Worker killed"})
	}

	w.WriteHeader(http.StatusOK)

	j, err := json.Marshal(resp)
	if err != nil {
		raven.CaptureError(err, nil)
		log.Fatal(err)
	}

	w.Write(j)
}
