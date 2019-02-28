package agent

import (
	"io"
	"net/http"
	"os"
	"sync"

	"github.com/kilgaloon/leprechaun/config"
	"github.com/kilgaloon/leprechaun/context"
	"github.com/kilgaloon/leprechaun/daemon"
	"github.com/kilgaloon/leprechaun/log"
	"github.com/kilgaloon/leprechaun/workers"
)

// Agent interface defines service that can be started/stop
// that has workers, config, context
type Agent interface {
	GetContext() *context.Context
	GetConfig() *config.AgentConfig
	GetStatus() daemon.ServiceStatus

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
	Name     string
	Config   *config.AgentConfig
	Context  *context.Context
	Stdin    io.Reader
	Stdout   io.Writer
	Debug    bool
	Status   int
	Pipeline chan string

	*sync.RWMutex
	log.Logs
	workers.Workers
}

// GetName returns agent name
func (d Default) GetName() string {
	return d.Name
}

// GetConfig return current config for agent
func (d *Default) GetConfig() *config.AgentConfig {
	return d.Config
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

// IsDebug determines is agent in debug mode
func (d *Default) IsDebug() bool {
	return d.Debug
}

// GetStatus returns status of agent expressed in int
func (d Default) GetStatus() daemon.ServiceStatus {
	return daemon.ServiceStatus(d.Status)
}

// SetPipeline set pipeline create string channel that
// agent will use to send through
func (d *Default) SetPipeline(pipe chan string) {
	d.Pipeline = pipe
}

// Stop agent
func (d *Default) Stop() {
	d.Status = daemon.Stopped
}

// Pause agent
func (d *Default) Pause() {
	d.Status = daemon.Paused
}

// Unpause agent
func (d *Default) Unpause() {
	d.Status = daemon.Started
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
func New(name string, cfg *config.AgentConfig, debug bool) *Default {
	agent := &Default{}
	agent.Config = cfg
	agent.RWMutex = new(sync.RWMutex)
	agent.Logs = log.Logs{
		Debug:    debug,
		ErrorLog: cfg.GetErrorLog(),
		InfoLog:  cfg.GetInfoLog(),
	}

	agent.Context = context.New()
	agent.Workers = workers.New(
		cfg,
		agent.Logs,
		agent.Context,
		agent.RWMutex,
	)

	agent.Stdin = os.Stdin
	agent.Stdout = os.Stdout
	agent.Debug = debug
	agent.Pipeline = make(chan string)

	return agent
}
