package daemon

import (
	"os"
	"strconv"

	"github.com/kilgaloon/leprechaun/agent"
	"github.com/kilgaloon/leprechaun/api"
)

// Daemon is long living process that serves as middleware
// and access to multiple agents
type Daemon struct {
	PID     int
	PidFile *os.File
	agents  []*agent.Default
	args    []string
	api     api.API
}

// GetPID gets current PID of client
func (d *Daemon) GetPID() int {
	return d.PID
}

// AddAgent push agent to list of agents
func (d *Daemon) AddAgent(a *agent.Default) {
	d.agents = append(d.agents, a)
}

// New create new daemon that is prepared to receive list of agents and sets up pid for usage
func New(pidPath *string, args []string, api api.API, agents []*agent.Default) *Daemon {
	d := new(Daemon)
	f, err := os.OpenFile(*pidPath, os.O_RDWR|os.O_CREATE, 0644)
	d.PidFile = f
	if err != nil {
		panic("Failed to start client, can't save PID, reason: " + err.Error())
	}

	d.PID = os.Getpid()
	pid := strconv.Itoa(d.PID)
	_, err = d.PidFile.WriteString(pid)
	if err != nil {
		panic("Failed to start client, can't save PID")
	}

	d.args = args
	d.api = api

	return d
}
