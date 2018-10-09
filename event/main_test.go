package event

import (
	"testing"
)

var (
	eventHandler = CreateHandler()
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
