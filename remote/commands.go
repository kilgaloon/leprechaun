package remote

import (
	"encoding/json"
	"net/http"

	"github.com/kilgaloon/leprechaun/api"
	"github.com/kilgaloon/leprechaun/daemon"
)

func (rem *Remote) cmdstart(w http.ResponseWriter, r *http.Request) {
	resp := api.TableResponse{
		Header:  []string{"Message"},
		Columns: [][]string{},
	}

	if rem.GetStatus() == daemon.Started {
		w.WriteHeader(http.StatusExpectationFailed)
		resp.Columns = append(resp.Columns, []string{"Client already started"})
	} else {
		go rem.Start()

		w.WriteHeader(http.StatusOK)
		resp.Columns = append(resp.Columns, []string{"Client started"})
	}

	j, err := json.Marshal(resp)
	if err != nil {
		rem.Error(err.Error())
	}

	w.Write(j)
}

func (rem *Remote) cmdstop(w http.ResponseWriter, r *http.Request) {
	resp := api.TableResponse{
		Header:  []string{"Message"},
		Columns: [][]string{},
	}

	if rem.GetStatus() == daemon.Started {
		err := rem.ln.Close()
		if err != nil {
			w.WriteHeader(http.StatusExpectationFailed)
			resp.Columns = append(resp.Columns, []string{"Client fail to stop"})
			rem.Error(err.Error())
		}

		rem.SetStatus(daemon.Stopped)
		if rem.GetStatus() == daemon.Stopped {
			w.WriteHeader(http.StatusOK)
			resp.Columns = append(resp.Columns, []string{"Client stopped"})
		}

		j, err := json.Marshal(resp)
		if err != nil {
			rem.Error(err.Error())
		}

		w.Write(j)
	} else {
		w.WriteHeader(http.StatusExpectationFailed)
		resp.Columns = append(resp.Columns, []string{"Client can't be stopped"})
	}
}
