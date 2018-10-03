package agent

import (
	"sync"

	"github.com/kilgaloon/leprechaun/context"
	"github.com/kilgaloon/leprechaun/workers"

	"github.com/kilgaloon/leprechaun/config"
	"github.com/kilgaloon/leprechaun/log"
)

// Agent interface defines client that can be started/stop
// that has workers, config, context and logs
type Agent interface {
	GetName() string
	GetWorkers() *workers.Workers
	GetContext() *context.Context
	GetConfig() *config.AgentConfig
	GetLogs() log.Logs
}

// Default represents default agent
type Default struct {
	Name    string
	PID     int
	Config  *config.AgentConfig
	Logs    log.Logs
	Mu      *sync.Mutex
	Workers *workers.Workers
	Context *context.Context
}

// GetName returns name of the client
func (d Default) GetName() string {
	return d.Name
}

// GetWorkers return instance of workers
func (d Default) GetWorkers() *workers.Workers {
	return d.Workers
}

// GetContext returns context of agent
func (d Default) GetContext() *context.Context {
	return d.Context
}

// GetLogs returns instance of logs
func (d Default) GetLogs() log.Logs {
	return d.Logs
}

// GetConfig return current config for agent
func (d Default) GetConfig() *config.AgentConfig {
	return d.Config
}

// New default client
func New(name string, cfg *config.AgentConfig) *Default {
	agent := &Default{}
	agent.Name = name
	agent.Config = cfg
	agent.Mu = new(sync.Mutex)
	agent.Logs = log.Logs{
		ErrorLog: cfg.GetErrorLog(),
		InfoLog:  cfg.GetInfoLog(),
	}

	agent.Context = context.New()
	agent.Workers = workers.New(
		cfg.GetMaxAllowedWorkers(),
		agent.Logs,
		agent.Context,
	)

	return agent
}
