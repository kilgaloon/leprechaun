package server

import (
	"net/http"
	"os"

	"github.com/kilgaloon/leprechaun/log"
)

// Agent holds instance of Server
var Agent *Server

// Server instance
type Server struct {
	Config *Config
	Logs   log.Logs
	Queue
	HTTP *http.Server
}

// CreateAgent new server
// Creating new agent will enable usage of Agent variable globally for packages
// that use this package
func CreateAgent(iniPath *string) *Server {
	server := &Server{}
	// load configurations for server
	server.Config = readConfig(*iniPath)
	server.HTTP = &http.Server{Addr: ":" + server.Config.port}

	Agent = server

	return Agent
}

// Start server that will receive webhooks
func (server *Server) Start() {
	// build queue for server
	server.BuildQueue()
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
}

// Stop http server
func (server *Server) Stop() os.Signal {
	server.Logs.Info("Shutting down server")
	if err := server.HTTP.Shutdown(nil); err != nil {
		panic(err)
	}

	return os.Interrupt
}
