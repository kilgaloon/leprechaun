package socket

import "errors"

// Command is closure that will be called to execute command
type Command func(args ...string) ([][]string, error)

// Registrator is agent that will be registered
// with regi
type Registrator struct {
	Name     string
	Commands map[string]Command
}

// Command Set command to registrator
func (r *Registrator) Command(name string, command Command) {
	r.Commands[name] = command
}

// Call specified command
func (r Registrator) Call(name string, args ...string) ([][]string, error) {
	if command, exist := r.Commands[name]; exist {
		return command(args...)
	}

	return nil, errors.New("Command does not exists, or it's not registered")
}

// CreateRegistrator create registrator struct
// to be pushed to Socket registry
func CreateRegistrator(name string) *Registrator {
	return &Registrator{
		Name:     name,
		Commands: make(map[string]Command),
	}
}
