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
	*agent.Default
	Pool
	HTTP *http.Server
}

// New create server
// Creating new agent will enable usage of Agent variable globally for packages
// that use this package
func New(name string, cfg *config.AgentConfig) *Server {
	server := &Server{
		agent.New(name, cfg),
		Pool{},
		&http.Server{Addr: ":" + strconv.Itoa(cfg.GetPort())},
	}

	Agent = server

	return server
}

// Start server that will receive webhooks
func (server *Server) Start() {
	server.BuildPool()
	// register all routes
	server.registerHandles()
	// listen for port
	server.GetLogs().Info("Server started")
	// register server to command socket
	go api.New(server.GetConfig().GetCommandSocket()).Register(server)
	if err := server.HTTP.ListenAndServe(); err != nil {
		server.GetLogs().Error("Httpserver: ListenAndServe() error: %s", err)
	}

}

func (server *Server) registerHandles() {
	http.HandleFunc(WebhookEndpoint, server.webhook)
	http.HandleFunc(PingEndpoint, server.ping)
}

// Stop http server
func (server *Server) Stop(args ...string) ([][]string, error) {
	server.GetLogs().Info("Shutting down server")
	if err := server.HTTP.Shutdown(con.Background()); err != nil {
		return [][]string{}, err
	}

	return [][]string{{"Server shutdown"}}, nil
}

// RegisterCommands to be used in internal communication
func (server Server) RegisterCommands() map[string]api.Command {
	cmds := make(map[string]api.Command)

	return server.DefaultCommands(cmds)
}
