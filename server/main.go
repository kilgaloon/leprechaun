package server

import (
	con "context"
	"net/http"
	"strconv"
	"sync"

	"github.com/kilgaloon/leprechaun/agent"
	"github.com/kilgaloon/leprechaun/api"
	"github.com/kilgaloon/leprechaun/config"
	"github.com/kilgaloon/leprechaun/context"
	"github.com/kilgaloon/leprechaun/log"
	"github.com/kilgaloon/leprechaun/recipe"
	"github.com/kilgaloon/leprechaun/workers"
)

// Agent holds instance of Server
var Agent *Server

// Server instance
type Server struct {
	Name   string
	Config *config.AgentConfig
	Logs   log.Logs
	Pool
	HTTP    *http.Server
	Context *context.Context
	mu      *sync.Mutex
	Workers *workers.Workers
}

// CreateAgent new server
// Creating new agent will enable usage of Agent variable globally for packages
// that use this package
func CreateAgent(name string, cfg *config.AgentConfig) *Server {
	def := agent.New(name, cfg)
	server := &Server{
		Name:    name,
		Config:  def.GetConfig(),
		Logs:    def.GetLogs(),
		Context: def.GetContext(),
		mu:      def.Mu,
		Workers: def.GetWorkers(),
		HTTP:    &http.Server{Addr: ":" + strconv.Itoa(cfg.GetPort())},
	}

	Agent = server

	return server
}

// GetName returns name of the Agent
func (server Server) GetName() string {
	return server.Name
}

// Start server that will receive webhooks
func (server *Server) Start() {
	// build queue for server
	server.BuildPool()
	// register all routes
	server.registerHandles()
	// listen for port
	server.Logs.Info("Server started")
	// register server to command socket
	go api.BuildSocket(server.Config.CommandSocket).Register(server)
	if err := server.HTTP.ListenAndServe(); err != nil {
		server.Logs.Error("Httpserver: ListenAndServe() error: %s", err)
	}

}

func (server *Server) registerHandles() {
	http.HandleFunc(WebhookEndpoint, server.webhook)
	http.HandleFunc(PingEndpoint, server.ping)
}

// Stop http server
func (server *Server) Stop(args ...string) ([][]string, error) {
	server.Logs.Info("Shutting down server")
	if err := server.HTTP.Shutdown(con.Background()); err != nil {
		return [][]string{}, err
	}

	return [][]string{{"Server shutdown"}}, nil
}

// GetWorkers return workers for agent
func (server *Server) GetWorkers() *workers.Workers {
	return server.Workers
}

// GetRecipeStack returns stack of recipes
func (server *Server) GetRecipeStack() []recipe.Recipe {
	var stack []recipe.Recipe
	for _, recipe := range server.Pool.Stack {
		stack = append(stack, recipe)
	}

	return stack
}

// GetConfig returns configuration for specific agent
func (server *Server) GetConfig() *config.AgentConfig {
	return server.Config
}

// GetContext returns context of agent
func (server *Server) GetContext() *context.Context {
	return server.Context
}

// GetLogs return logs of agent
func (server *Server) GetLogs() log.Logs {
	return server.Logs
}
