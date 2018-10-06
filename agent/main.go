package agent

import (
	"sync"

	"github.com/kilgaloon/leprechaun/api"
	"github.com/kilgaloon/leprechaun/context"
	"github.com/kilgaloon/leprechaun/workers"

	"github.com/kilgaloon/leprechaun/config"
	"github.com/kilgaloon/leprechaun/log"
)

// Agent interface defines service that can be started/stop
// that has workers, config, context and logs
type Agent interface {
	GetName() string
	GetWorkers() *workers.Workers
	GetContext() *context.Context
	GetConfig() *config.AgentConfig
	GetLogs() log.Logs
	GetSocket() *api.Socket
	GetMutex() *sync.Mutex

	SetPID(i int)
	GetPID() int

	DefaultCommands(map[string]api.Command) map[string]api.Command
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
	Socket  *api.Socket
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

// GetSocket returns information about socket
// and commands available for internal communication
func (d Default) GetSocket() *api.Socket {
	return d.Socket
}

// GetMutex for agent
func (d Default) GetMutex() *sync.Mutex {
	return d.Mu
}

// SetPID sets process id for agent
func (d *Default) SetPID(i int) {
	d.PID = i
}

// GetPID sets process id for agent
func (d Default) GetPID() int {
	return d.PID
}

// DefaultCommands merge 2 maps into one
// it usability is if some of the agents
// wants to takeover default commands
func (d Default) DefaultCommands(commands map[string]api.Command) map[string]api.Command {
	cmds := make(map[string]api.Command)

	cmds["workers:list"] = api.Command{
		Closure: d.WorkersList,
		Definition: api.Definition{
			Text:  "List all currently active workers",
			Usage: "{agent} workers:list",
		},
	}

	cmds["workers:kill"] = api.Command{
		Closure: d.KillWorker,
		Definition: api.Definition{
			Text:  "Kills currently active worker by job name",
			Usage: "{agent} workers:kill {job}",
		},
	}

	for name, command := range commands {
		cmds[name] = command
	}

	return cmds
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
	agent.Socket = api.New(cfg.GetCommandSocket())

	return agent
}
