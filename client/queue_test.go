package client

import (
	"testing"
)

func TestBuildQueue(t *testing.T) {
	mu.Lock()
	defer mu.Unlock()

	Agent.BuildQueue()
	if len(Agent.Queue.Stack) != 5 {
		t.Errorf("Queue stack length expected to be 5, got %d", len(Agent.Queue.Stack))
	}

}

func TestQueue(t *testing.T) {
	mu.Lock()
	defer mu.Unlock()
	// reset queue to 0 to test AddToQueue
	q := &Agent.Queue
	q.Stack = q.Stack[:0]

	Agent.AddToQueue(Agent.GetConfig().GetRecipesPath() + "/schedule.yml")

	if len(Agent.Queue.Stack) != 1 {
		t.Errorf("Queue stack length expected to be 1, got %d", len(Agent.Queue.Stack))
	}

	Agent.AddToQueue(Agent.GetConfig().GetRecipesPath() + "/hook.yml")

	if len(Agent.Queue.Stack) != 1 {
		t.Errorf("Queue stack length expected to be 0, got %d", len(Agent.Queue.Stack))
	}

	Agent.Pause()
	Agent.ProcessQueue()
}
