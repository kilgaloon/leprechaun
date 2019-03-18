package server

import (
	con "context"
	"net/http"
	"strconv"
	"strings"

	"github.com/kilgaloon/leprechaun/agent"
	"github.com/kilgaloon/leprechaun/config"
	"github.com/kilgaloon/leprechaun/daemon"
	"github.com/mholt/certmagic"
)

// Agent holds instance of server client
var Agent *Server

// Server instance
type Server struct {
	Name string
	*agent.Default
	Pool
	HTTP *http.Server
}

// New create server
// Creating new agent will enable usage of Agent variable globally for packages
// that use this package
func (server *Server) New(name string, cfg *config.AgentConfig, debug bool) daemon.Service {
	s := &Server{
		name,
		agent.New(name, cfg, debug),
		Pool{},
		&http.Server{Addr: ":" + strconv.Itoa(cfg.GetPort())},
	}

	Agent = s

	return s
}

//GetName returns server name
func (server *Server) GetName() string {
	return server.Name
}

// Start server that will receive webhooks
func (server *Server) Start() {
	server.BuildPool()
	// register all routes
	server.registerHandles()
	server.SetStatus(daemon.Started)

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
	server.Lock()
	defer server.Unlock()

	mux := http.NewServeMux()
	mux.HandleFunc(WebhookEndpoint, server.webhook)
	mux.HandleFunc(PingEndpoint, server.ping)

	server.HTTP.Handler = mux
}

// Stop http server
func (server *Server) Stop() {
	server.Info("Shutting down server")
	if err := server.HTTP.Shutdown(con.Background()); err != nil {
		server.Error(err.Error())
	}

	server.SetStatus(daemon.Stopped)
}

// RegisterAPIHandles to be used in socket communication
// If you want to takeover default commands from agent
// call DefaultCommands from Agent which is same command
func (server *Server) RegisterAPIHandles() map[string]func(w http.ResponseWriter, r *http.Request) {
	cmds := make(map[string]func(w http.ResponseWriter, r *http.Request))

	cmds["info"] = server.cmdinfo
	cmds["stop"] = server.cmdstop
	cmds["start"] = server.cmdstart
	cmds["pause"] = server.cmdpause

	// this function merge both maps and inject default commands from agent
	return cmds
}

func (server *Server) isTLS() bool {
	server.Lock()
	defer server.Unlock()

	return strings.Contains(server.GetConfig().Domain, "https://")
}
