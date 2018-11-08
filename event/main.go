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
	events       map[string]eventClosure
	eventChannel chan string
	sync.Mutex
}

// Subscribe for event and trigger callback
func (handler *Handler) Subscribe(event string, callback eventClosure) {
	handler.Lock()
	handler.events[event] = callback
	handler.Unlock()
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
				handler.Lock()
				if trigger, subscribed := handler.events[event]; subscribed {
					trigger()
					log.Logger.Info("Event %s dispatched", event)
				}
				handler.Unlock()

			}
		}
	}()
}

// CreateHandler creates new handler
func CreateHandler() *Handler {
	handler := &Handler{}
	handler.events = make(map[string]eventClosure)
	handler.eventChannel = make(chan string)

	return handler
}

// Start listener
func init() {
	EventHandler = CreateHandler()
	EventHandler.listen()
}
