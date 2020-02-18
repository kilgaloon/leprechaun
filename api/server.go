package api

import (
	"context"
	"net/http"
)

// API defines socket on which we listen for commands
type API struct {
	HTTP *http.Server
	registeredHandles map[string]func(w http.ResponseWriter, r *http.Request)
}

// Registrator defines interface that help us to register http handler
type Registrator interface {
	GetName() string
	RegisterAPIHandles() map[string]func(w http.ResponseWriter, r *http.Request)
	DefaultAPIHandles() map[string]func(w http.ResponseWriter, r *http.Request)
}

// Register all available http handles
func (a *API) Register(r Registrator) *API {
	for p, f := range r.RegisterAPIHandles() {
		a.RegisterHandle(r.GetName()+"/"+p, f)
	}

	for p, f := range r.DefaultAPIHandles() {
		a.RegisterHandle(r.GetName()+"/"+p, f)
	}

	return a
}

// RegisterHandle append handle before server is started
func (a *API) RegisterHandle(e string, h func(w http.ResponseWriter, r *http.Request)) {
	pattern := "/" + e;

	if _, exist := a.registeredHandles[pattern]; !exist {
		a.registeredHandles[pattern] = h
	}
}

// Start api server
func (a *API) Start() {
	if !IsAPIRunning() {
		mux := http.NewServeMux()
		for pattern, handler := range a.registeredHandles {
			mux.HandleFunc(pattern, handler)
		}

		a.HTTP.Handler = mux
		if err := a.HTTP.ListenAndServe(); err == http.ErrServerClosed {
			a.HTTP = &http.Server{
				Addr: ":11401",
			}

			a.Start()
		}
	}
}

// Stop api server
func (a *API) Stop() error {
	return a.HTTP.Shutdown(context.Background())
}

// New creates new socket
func New() *API {
	api := &API{
		&http.Server{
			Addr: ":11401",
		},
		make(map[string]func(w http.ResponseWriter, r *http.Request)),
	}

	return api
}
