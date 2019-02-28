package daemon

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/kilgaloon/leprechaun/api"
)

// InfoResponse defines how response looks like when
// info for running daemon is returned
type InfoResponse struct {
	PID        int
	ConfigPath string
	PidPath    string
	Debug      bool
}

// GetInfo display info for running daemon
func (d *Daemon) GetInfo() *InfoResponse {
	r, err := api.HTTPClient.Get(api.RevealEndpoint("/{agent}", api.Cmd("daemon")))
	if err != nil {
		log.Fatal(err)
	}

	if r.StatusCode != 200 {
		fmt.Println("No such command")
		return &InfoResponse{}
	}

	defer r.Body.Close()

	resp := &InfoResponse{}
	err = json.NewDecoder(r.Body).Decode(resp)
	if err != nil {
		log.Fatal(err)
	}

	return resp
}
