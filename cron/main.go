package cron

import (
	"net/http"

	"github.com/robfig/cron"

	"github.com/kilgaloon/leprechaun/agent"

	"github.com/kilgaloon/leprechaun/config"
)

// Agent holds instance of Client
var Agent *Cron

// Cron settings and configurations
type Cron struct {
	*agent.Default
	Service *cron.Cron
}

// New create client
// Creating new agent will enable usage of Agent variable globally for packages
// that use this package
func New(name string, cfg *config.AgentConfig, debug bool) *Cron {
	cron := &Cron{
		agent.New(name, cfg, debug),
		cron.New(),
	}

	Agent = cron

	return cron
}

// Start client
func (c *Cron) Start() {
	// build queue
	c.GetMutex().Lock()
	c.buildJobs()
	c.GetMutex().Unlock()

	c.Service.Start()

	c.Event.Dispatch("cron:ready")
}

// RegisterAPIHandles to be used in socket communication
// If you want to takeover default commands from agent
// call DefaultCommands from Agent which is same command
func (c *Cron) RegisterAPIHandles() map[string]func(w http.ResponseWriter, r *http.Request) {
	cmds := make(map[string]func(w http.ResponseWriter, r *http.Request))

	// this function merge both maps and inject default commands from agent
	return cmds
}

// Stop client
func (c *Cron) Stop() {
	c.Service.Stop()
}
