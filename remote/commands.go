package remote

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/getsentry/raven-go"
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
		raven.CaptureError(err, nil)
		log.Fatal(err)
	}

	w.Write(j)
}
