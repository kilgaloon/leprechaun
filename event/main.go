package event

import (
	"sync"

	"github.com/kilgaloon/leprechaun/log"
)

// EventHandler Create new struct of handler for global usage
var EventHandler *Handler

type eventClosure func()

// Handler handles events
type Handler struct {
	events       map[string][]eventClosure
	eventChannel chan string
	mu           sync.Mutex
	log.Logs
}

// Subscribe for event and trigger callback
func (handler *Handler) Subscribe(event string, callback eventClosure) {
	handler.mu.Lock()
	handler.events[event] = append(handler.events[event], callback)
	handler.mu.Unlock()
}

// Dispatch an event
func (handler *Handler) Dispatch(event string) {
	handler.eventChannel <- event
}

// Listen listens for events
func (handler *Handler) listen() {
	go func() {
		for {
			select {
			case event := <-handler.eventChannel:
				handler.mu.Lock()
				if events, subscribed := handler.events[event]; subscribed {
					for _, trigger := range events {
						trigger()
					}

					handler.Info("Event %s dispatched", event)
				}
				handler.mu.Unlock()

			}
		}
	}()

}

// NewHandler creates new handler
func NewHandler(log log.Logs) *Handler {
	handler := &Handler{}
	handler.events = make(map[string][]eventClosure)
	handler.eventChannel = make(chan string)
	handler.Logs = log

	handler.listen()

	return handler
}
