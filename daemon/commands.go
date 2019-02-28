package daemon

import (
	"encoding/json"
	"log"
	"net/http"
)

func (d *Daemon) daemonInfo(w http.ResponseWriter, r *http.Request) {
	resp := &InfoResponse{
		PID:        d.PID,
		ConfigPath: d.ConfigPath,
		PidPath:    d.PidPath,
		Debug:      d.Debug,
	}

	w.WriteHeader(http.StatusOK)

	j, err := json.Marshal(resp)
	if err != nil {
		log.Fatal(err)
	}

	w.Write(j)

	return
}
