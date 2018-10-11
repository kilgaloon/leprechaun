package server

import (
	con "context"
	"net/http"
	"strconv"

	"github.com/kilgaloon/leprechaun/agent"
	"github.com/kilgaloon/leprechaun/api"
	"github.com/kilgaloon/leprechaun/config"
)

// Agent holds instance of Server
var Agent *Server

// Server instance
type Server struct {
	Agent agent.Agent
	Pool
	HTTP *http.Server
}

// New create server
// Creating new agent will enable usage of Agent variable globally for packages
// that use this package
func New(name string, cfg *config.AgentConfig) *Server {
	server := &Server{
		Agent: agent.New(name, cfg),
		HTTP:  &http.Server{Addr: ":" + strconv.Itoa(cfg.GetPort())},
	}

	Agent = server

	return server
}

// GetName of agent
func (server *Server) GetName() string {
	return server.Agent.GetName()
}

// Start server that will receive webhooks
func (server *Server) Start() {
	// build queue for server
	server.Agent.GetMutex().Lock()
	server.BuildPool()
	server.Agent.GetMutex().Unlock()
	// register all routes
	server.registerHandles()
	// listen for port
	server.Agent.GetLogs().Info("Server started")
	// register server to command socket
	go api.New(server.Agent.GetConfig().GetCommandSocket()).Register(server)
	if err := server.HTTP.ListenAndServe(); err != nil {
		server.Agent.GetLogs().Error("Httpserver: ListenAndServe() error: %s", err)
	}

}

func (server *Server) registerHandles() {
	http.HandleFunc(WebhookEndpoint, server.webhook)
	http.HandleFunc(PingEndpoint, server.ping)
}

// Stop http server
func (server *Server) Stop(args ...string) ([][]string, error) {
	server.Agent.GetLogs().Info("Shutting down server")
	if err := server.HTTP.Shutdown(con.Background()); err != nil {
		return [][]string{}, err
	}

	return [][]string{{"Server shutdown"}}, nil
}

// RegisterCommands to be used in internal communication
func (server Server) RegisterCommands() map[string]api.Command {
	cmds := make(map[string]api.Command)

	return server.Agent.DefaultCommands(cmds)
}
