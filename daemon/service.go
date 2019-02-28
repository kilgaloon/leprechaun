package daemon

import (
	"github.com/kilgaloon/leprechaun/api"
	"github.com/kilgaloon/leprechaun/config"
)

const (
	// Started status
	Started = iota + 1
	// Stopped status
	Stopped
	// Paused status
	Paused
)

// ServiceStatus distincts status of service
type ServiceStatus int

func (s ServiceStatus) String() string {
	switch int(s) {
	case 1:
		return "Started"
	case 2:
		return "Stopped"
	case 3:
		return "Paused"
	default:
		return "Unknown"
	}
}

// Service struct define
type Service interface {
	api.Registrator
	GetStatus() ServiceStatus
	Start()
	Stop()
	Pause()
	Unpause()
	IsDebug() bool
	SetPipeline(chan string)
	New(name string, cfg *config.AgentConfig, debug bool) Service
}
