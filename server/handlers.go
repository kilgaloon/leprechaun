package server

import (
	"net/http"
)

const (
	// WebhookEndpoint defines endpoint where webhook is
	WebhookEndpoint = "/hook"
)

func (server Server) webhook(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query()["id"][0]
	// find recipe with that id
	server.FindInPool(key)
}
