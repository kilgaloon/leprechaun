package client

import (
	"testing"
)

func TestBuildQueue(t *testing.T) {
	Agent.BuildQueue()

	Agent.Lock()
	if len(Agent.Queue.Stack) != 5 {
		t.Errorf("Queue stack length expected to be 5, got %d", len(Agent.Queue.Stack))
	}
	Agent.Unlock()

}

func TestQueue(t *testing.T) {
	// reset queue to 0 to test AddToQueue
	Agent.Lock()
	q := &Agent.Queue
	q.Stack = q.Stack[:0]
	Agent.Unlock()

	Agent.AddToQueue(Agent.GetConfig().GetRecipesPath() + "/schedule.yml")

	Agent.Lock()
	if len(Agent.Queue.Stack) != 1 {
		t.Errorf("Queue stack length expected to be 1, got %d", len(Agent.Queue.Stack))
	}
	Agent.Unlock()

	Agent.AddToQueue(Agent.GetConfig().GetRecipesPath() + "/hook.yml")

	Agent.Lock()
	if len(Agent.Queue.Stack) != 1 {
		t.Errorf("Queue stack length expected to be 0, got %d", len(Agent.Queue.Stack))
	}
	Agent.Unlock()

	Agent.Pause()
	Agent.ProcessQueue()
}
