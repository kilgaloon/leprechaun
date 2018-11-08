package cron

import (
	"github.com/robfig/cron"

	"github.com/kilgaloon/leprechaun/agent"

	"github.com/kilgaloon/leprechaun/api"
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
func New(name string, cfg *config.AgentConfig) *Cron {
	cron := &Cron{
		agent.New(name, cfg),
		cron.New(),
	}

	Agent = cron

	return cron
}

// Start client
func (c *Cron) Start() {
	c.GetMutex().Lock()
	defer c.GetMutex().Unlock()
	// build queue
	c.Lock()
	c.buildJobs()
	c.Unlock()

	// register client to command socket
	go api.New(c.GetConfig().GetCommandSocket()).Register(c)

	c.Service.Start()
}

// RegisterCommands to be used in socket communication
// If you want to takeover default commands from agent
// call DefaultCommands from Agent which is same command
func (c *Cron) RegisterCommands() map[string]api.Command {
	cmds := make(map[string]api.Command)

	// this function merge both maps and inject default commands from agent
	return c.DefaultCommands(cmds)
}

// Stop client
func (c *Cron) Stop() {
	c.Service.Stop()
}
