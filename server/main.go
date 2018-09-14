package server

import (
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/kilgaloon/leprechaun/config"
	"github.com/kilgaloon/leprechaun/context"
	"github.com/kilgaloon/leprechaun/log"
	"github.com/kilgaloon/leprechaun/workers"
)

// Agent holds instance of Server
var Agent *Server

// Server instance
type Server struct {
	Config *config.ServerConfig
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
func CreateAgent(cfg *config.ServerConfig) *Server {
	server := &Server{}
	// load configurations for server
	server.Config = cfg
	server.Context = context.BuildContext(server)
	server.HTTP = &http.Server{Addr: ":" + strconv.Itoa(server.Config.Port)}
	server.mu = new(sync.Mutex)
	server.Logs = log.Logs{
		ErrorLog: server.Config.ErrorLog,
		InfoLog:  server.Config.InfoLog,
	}
	server.Workers = workers.BuildWorkers(server.Context, cfg.MaxAllowedWorkers, server.Logs)
	

	Agent = server

	return Agent
}

// Start server that will receive webhooks
func (server *Server) Start() {
	// build queue for server
	server.BuildPool()
	// register all routes
	server.registerHandles()
	// listen for port
	server.Logs.Info("Server started")
	if err := server.HTTP.ListenAndServe(); err != nil {
		server.Logs.Error("Httpserver: ListenAndServe() error: %s", err)
	}

}

func (server *Server) registerHandles() {
	http.HandleFunc(WebhookEndpoint, server.webhook)
	http.HandleFunc(PingEndpoint, server.ping)
}

// Stop http server
func (server *Server) Stop() os.Signal {
	server.Logs.Info("Shutting down server")
	if err := server.HTTP.Shutdown(nil); err != nil {
		panic(err)
	}

	return os.Interrupt
}

// GetConfig Gets config for server
func (server Server) GetConfig() *config.ServerConfig {
	return server.Config
}
