package event

import (
	"testing"

	"github.com/kilgaloon/leprechaun/config"
	"github.com/kilgaloon/leprechaun/log"
)

var (
	cfgWrap = config.NewConfigs()
	cfg     = cfgWrap.New("test", "../tests/configs/config_regular.ini")
	logger  = log.Logs{
		ErrorLog: cfg.GetErrorLog(),
		InfoLog:  cfg.GetInfoLog(),
	}
	eventHandler = NewHandler(logger)
)

func TestListen(t *testing.T) {
	go eventHandler.listen()
}

func TestSubscribe(t *testing.T) {
	// it doesnt do nothing, we just want to see
	// is this event subscribed
	eventHandler.Subscribe("test", func() {})

	if len(eventHandler.events) < 1 {
		t.Fatalf("Expected number of events 1 we got %d", len(eventHandler.events))
	}
}

func TestDispatch(t *testing.T) {
	eventHandler.Dispatch("test")
}
