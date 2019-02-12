package server

import (
	con "context"
	"net/http"
	"strconv"
	"strings"

	"github.com/kilgaloon/leprechaun/agent"
	"github.com/kilgaloon/leprechaun/config"
	"github.com/mholt/certmagic"
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
	server.Info("Server started")
	// register server to command socket

	if server.isTLS() {
		certmagic.Agreed = true
		certmagic.Email = server.GetConfig().GetNotificationsEmail()
		certmagic.CA = certmagic.LetsEncryptStagingCA

		if err := certmagic.HTTPS(server.GetConfig().GetServerDomain(), server.HTTP.Handler); err != nil {
			server.Error("Httpserver: ListenAndServe() error: %s", err)
		}
	} else {
		if err := server.HTTP.ListenAndServe(); err != nil {
			server.Error("Httpserver: ListenAndServe() error: %s", err)
		}
	}

}

func (server *Server) registerHandles() {
	mux := http.NewServeMux()
	mux.HandleFunc(WebhookEndpoint, server.webhook)
	mux.HandleFunc(PingEndpoint, server.ping)

	server.HTTP.Handler = mux
}

// Stop http server
func (server *Server) Stop(args ...string) ([][]string, error) {
	server.Info("Shutting down server")
	if err := server.HTTP.Shutdown(con.Background()); err != nil {
		return [][]string{}, err
	}

	return [][]string{{"Server shutdown"}}, nil
}

// RegisterAPIHandles to be used in socket communication
// If you want to takeover default commands from agent
// call DefaultCommands from Agent which is same command
func (server *Server) RegisterAPIHandles() map[string]func(w http.ResponseWriter, r *http.Request) {
	cmds := make(map[string]func(w http.ResponseWriter, r *http.Request))

	// this function merge both maps and inject default commands from agent
	return cmds
}

func (server Server) isTLS() bool {
	return strings.Contains(server.GetConfig().Domain, "https://")
}
