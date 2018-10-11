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
	Agent   agent.Agent
	Service *cron.Cron
}

// New create client
// Creating new agent will enable usage of Agent variable globally for packages
// that use this package
func New(name string, cfg *config.AgentConfig) *Cron {
	cron := &Cron{
		Agent:   agent.New(name, cfg),
		Service: cron.New(),
	}

	Agent = cron

	return cron
}

// GetName of agent
func (c *Cron) GetName() string {
	return c.Agent.GetName()
}

// Start client
func (c *Cron) Start() {
	// build queue
	c.Agent.Lock()
	c.buildJobs()
	c.Agent.Unlock()

	// register client to command socket
	go api.New(c.Agent.GetConfig().GetCommandSocket()).Register(c)

	c.Service.Start()
}

// RegisterCommands to be used in socket communication
// If you want to takeover default commands from agent
// call DefaultCommands from Agent which is same command
func (c *Cron) RegisterCommands() map[string]api.Command {
	cmds := make(map[string]api.Command)

	// this function merge both maps and inject default commands from agent
	return c.Agent.DefaultCommands(cmds)
}

// Stop client
func (c *Cron) Stop() {
	c.Service.Stop()
}
