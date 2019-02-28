package cron

import (
	"net/http"

	"github.com/robfig/cron"

	"github.com/kilgaloon/leprechaun/agent"
	"github.com/kilgaloon/leprechaun/daemon"

	"github.com/kilgaloon/leprechaun/config"
)

// Cron settings and configurations
type Cron struct {
	Name string
	*agent.Default
	Service *cron.Cron
}

// New create client
// Creating new agent will enable usage of Agent variable globally for packages
// that use this package
func (c *Cron) New(name string, cfg *config.AgentConfig, debug bool) daemon.Service {
	cron := &Cron{
		name,
		agent.New(name, cfg, debug),
		cron.New(),
	}

	return cron
}

// GetName returns agent name
func (c Cron) GetName() string {
	return c.Name
}

// Start client
func (c *Cron) Start() {
	c.buildJobs()

	c.Service.Start()

	c.Info("Cron started")
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
	c.Status = daemon.Stopped

	c.Info("Cron service stopped")
}
