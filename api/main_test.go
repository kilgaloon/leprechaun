package api

import (
	"errors"
	"sync"
	"testing"
)

type TestAgent struct{}

// Test registering commands
func (ta TestAgent) RegisterCommands() map[string]Command {
	var cmds = make(map[string]Command)

	cmds["test"] = Command{
		Closure: func(args ...string) ([][]string, error) {
			var resp = [][]string{
				[]string{"TEST"},
			}

			return resp, nil
		},
		Definition: Definition{},
	}

	cmds["test_with_error"] = Command{
		Closure: func(args ...string) ([][]string, error) {
			return nil, errors.New("Test error")
		},
		Definition: Definition{},
	}

	return cmds
}

var (
	API   = New("../tests/var/run/leprechaun/.sock")
	Agent = &TestAgent{}
)

func TestRegister(t *testing.T) {
	var wg sync.WaitGroup

	wg.Add(1)

	go API.Register(Agent)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-API.readyChan:
				resp := API.Command("agent test")
				if resp == "" {
					t.Fail()
				}

				resp2 := API.Command("agent test_with_error")
				if resp2 == "" {
					t.Fail()
				}

				resp3 := API.Command("test")
				if resp3 == "" {
					t.Fail()
				}

				resp4 := API.Command("agent not_exist")
				if resp4 == "" {
					t.Fail()
				}

				return
			}

		}
	}()

	wg.Wait()
}

func TestCall(t *testing.T) {
	API.commands = Agent.RegisterCommands()
	r, err := API.Call("test")
	if err != nil {
		t.Fail()
	}

	if r[0][0] != "TEST" {
		t.Fail()
	}

	_, err = API.Call("test_with_error")
	if err == nil {
		t.Fail()
	}
}
