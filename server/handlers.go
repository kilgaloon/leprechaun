package server

import (
	"net/http"
)

const (
	// WebhookEndpoint defines endpoint where webhook is
	WebhookEndpoint = "/hook"
	// PingEndpoint defines endpoints for healthcheck
	PingEndpoint = "/ping"
)

func (server Server) webhook(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query()["id"][0]
	// find recipe with that id
	server.FindInPool(key)
}

func (server Server) ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	_, err := w.Write([]byte("PONG"))
	if err != nil {
		server.Agent.GetLogs().Error("%s", err)
	}
}
