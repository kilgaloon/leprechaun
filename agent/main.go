package agent

import (
	"io"
	"os"
	"sync"

	"github.com/kilgaloon/leprechaun/api"
	"github.com/kilgaloon/leprechaun/config"
	"github.com/kilgaloon/leprechaun/context"
	"github.com/kilgaloon/leprechaun/event"
	"github.com/kilgaloon/leprechaun/log"
	"github.com/kilgaloon/leprechaun/workers"
)

// Agent interface defines service that can be started/stop
// that has workers, config, context and logs
type Agent interface {
	GetName() string
	GetContext() *context.Context
	GetConfig() *config.AgentConfig
	GetLogs() log.Logs
	GetSocket() *api.Socket

	DefaultCommands(map[string]api.Command) map[string]api.Command

	StandardIO
}

// StandardIO of agent
type StandardIO interface {
	StandardInput
	StandardOutput
}

// StandardInput holds everything for input
type StandardInput interface {
	GetStdin() io.Reader
	SetStdin(r io.Reader)
	io.Reader
}

// StandardOutput holds everything for output
type StandardOutput interface {
	GetStdout() io.Writer
	SetStdout(w io.Writer)
	io.Writer
}

// Default represents default agent
type Default struct {
	Name    string
	PID     int
	Config  *config.AgentConfig
	Mu      *sync.Mutex
	Context *context.Context
	Socket  *api.Socket
	Stdin   io.Reader
	Stdout  io.Writer
	Event   *event.Handler

	log.Logs
	workers.Workers
}

// GetName returns name of the client
func (d Default) GetName() string {
	return d.Name
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

func (d Default) Write(p []byte) (n int, err error) {
	return d.GetStdout().Write(p)
}

func (d Default) Read(p []byte) (n int, err error) {
	return d.GetStdin().Read(p)
}

// GetStdout get agent standard output that can be written to
func (d Default) GetStdout() io.Writer {
	return d.Stdout
}

// GetStdin get agent standard input that can be read from to
func (d Default) GetStdin() io.Reader {
	return d.Stdin
}

// SetStdin ability to change standard input for agent
func (d *Default) SetStdin(r io.Reader) {
	d.Stdin = r
}

// SetStdout ability to change standard input for agent
func (d *Default) SetStdout(w io.Writer) {
	d.Stdout = w
}

// DefaultCommands merge 2 maps into one
// it usability is if some of the agents
// wants to takeover default commands
func (d Default) DefaultCommands(commands map[string]api.Command) map[string]api.Command {
	d.GetMutex().Lock()
	defer d.GetMutex().Unlock()

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
		cfg,
		agent.Logs,
		agent.Context,
		agent.Mu,
	)
	agent.Socket = api.New(cfg.GetCommandSocket())
	agent.Stdin = os.Stdin
	agent.Stdout = os.Stdout
	agent.Event = event.NewHandler(agent.Logs)

	return agent
}
