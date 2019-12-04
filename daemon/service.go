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

// StartStop defines service that can be started and stoped
type StartStop interface {
	Start()
	Stop()
}

// Pause defines service that can be paused and unpaused
type Pause interface {
	Pause()
}

// Service struct define
type Service interface {
	api.Registrator
	GetStatus() ServiceStatus
	SetStatus(s int)
	GetConfig() config.AgentConfig
	StartStop
	Pause
	IsDebug() bool
	New(name string, cfg *config.AgentConfig, debug bool) Service
}
