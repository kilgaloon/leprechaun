package agent

import (
	"io"
	"net/http"
	"os"
	"sync"

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
	Mu      *sync.RWMutex
	Context *context.Context
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

// GetMutex for agent
func (d Default) GetMutex() *sync.RWMutex {
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

// DefaultAPIHandles to be used in socket communication
// If you want to takeover default commands from agent
// call DefaultCommands from Agent which is same command
func (d *Default) DefaultAPIHandles() map[string]func(w http.ResponseWriter, r *http.Request) {
	cmds := make(map[string]func(w http.ResponseWriter, r *http.Request))

	cmds["workers/list"] = d.WorkersList
	cmds["workers/kill"] = d.KillWorker

	// this function merge both maps and inject default commands from agent
	return cmds
}

// New default client
func New(name string, cfg *config.AgentConfig) *Default {
	agent := &Default{}
	agent.Name = name
	agent.Config = cfg
	agent.Mu = new(sync.RWMutex)
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

	agent.Stdin = os.Stdin
	agent.Stdout = os.Stdout
	agent.Event = event.NewHandler(agent.Logs)

	return agent
}
