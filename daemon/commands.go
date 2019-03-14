package daemon

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime"
	"strconv"

	"github.com/kilgaloon/leprechaun/api"
)

func (d *Daemon) daemonInfo(w http.ResponseWriter, r *http.Request) {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	alloc := strconv.FormatFloat(float64(mem.Alloc/1024)/1024, 'f', 2, 64)

	resp := &InfoResponse{
		PID:        d.PID,
		ConfigPath: d.GetConfigPath(),
		PidPath:    d.GetPidPath(),
		Debug:      d.Debug,
		Memory:     alloc + "MiB",
	}

	w.WriteHeader(http.StatusOK)

	j, err := json.Marshal(resp)
	if err != nil {
		log.Fatal(err)
	}

	w.Write(j)

	return
}

func (d *Daemon) daemonKill(w http.ResponseWriter, r *http.Request) {
	resp := api.MessageResponse{}

	d.Kill()
	if api.IsAPIRunning() {
		resp.Message = "Failed to kill daemon"
	} else {
		resp.Message = "Daemon killed"
	}

	j, err := json.Marshal(resp)
	if err != nil {
		resp.Message = "Daemon killed"
	}

	w.Write(j)
}
