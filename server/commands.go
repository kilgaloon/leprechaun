package server

import (
	"github.com/kilgaloon/leprechaun/api"
)

// RegisterCommandSocket returns Registrator
func (s *Server) RegisterCommandSocket() *api.Registrator {
	r := api.CreateRegistrator(s)

	// register commands
	r.Command("stop", s.Stop)

	return r
}
 