package daemon

import (
	"encoding/json"
	"log"
	"os"
	"strconv"

	"github.com/kilgaloon/leprechaun/api"
	"github.com/olekukonko/tablewriter"
)

// InfoResponse defines how response looks like when
// info for running daemon is returned
type InfoResponse struct {
	PID        int
	ConfigPath string
	PidPath    string
	Debug      bool
	Memory     string
}

// GetInfo display info for running daemon
func (d *Daemon) GetInfo() *InfoResponse {
	r, err := api.HTTPClient.Get(api.RevealEndpoint("/{agent}/info", api.Cmd("daemon")))
	if err != nil {
		log.Fatal(err)
	}

	defer r.Body.Close()

	resp := &InfoResponse{}
	err = json.NewDecoder(r.Body).Decode(resp)
	if err != nil {
		log.Fatal(err)
	}

	return resp
}

func (d *Daemon) renderInfo() {
	table := tablewriter.NewWriter(os.Stdout)

	resp := d.GetInfo()

	pid := strconv.Itoa(resp.PID)
	debug := "No"
	if resp.Debug {
		debug = "Yes"
	}

	table.SetHeader([]string{"PID", "Config path", "Pid path", "Debug", "Memory"})
	table.Append([]string{pid, resp.ConfigPath, resp.PidPath, debug, resp.Memory})

	table.Render()
}

func (d *Daemon) killDaemon() {
	api.HTTPClient.Get(api.RevealEndpoint("/{agent}/kill", api.Cmd("daemon")))
}
