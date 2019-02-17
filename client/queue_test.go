package client

import (
	"testing"
)

var (
	fk  = New("test", cfgWrap.New("test", *path), false)
	fk2 = New("test", cfgWrap.New("test", *path), false)
)

func TestBuildQueue(t *testing.T) {
	fk.BuildQueue()

	if len(fk.Queue.Stack) != 4 {
		t.Errorf("Queue stack length expected to be 4, got %d", len(fk.Queue.Stack))
	}

	// reset queue to 0 to test AddToQueue
	q := &fk.Queue
	q.Stack = q.Stack[:0]

}

func TestAddToQueue(t *testing.T) {
	fk.AddToQueue(&fk.Queue.Stack, fk.GetConfig().GetRecipesPath()+"/schedule.yml")
	if len(fk.Queue.Stack) != 1 {
		t.Errorf("Queue stack length expected to be 1, got %d", len(fk.Queue.Stack))
	}

	fk.AddToQueue(&fk.Queue.Stack, fk.GetConfig().GetRecipesPath()+"/hook.yml")
	if len(fk.Queue.Stack) != 1 {
		t.Errorf("Queue stack length expected to be 0, got %d", len(fk.Queue.Stack))
	}

	q := &fk.Queue
	q.Stack = q.Stack[:0]

}

func TestProcessQueueNotStoppedClient(t *testing.T) {
	fk2.AddToQueue(&fk2.Queue.Stack, fk2.GetConfig().GetRecipesPath()+"/schedule.yml")
	if len(fk2.Queue.Stack) != 1 {
		t.Errorf("Queue stack length expected to be 1, got %d", len(fk2.Queue.Stack))
	}

	fk2.stopped = false
	fk2.ProcessQueue()

	workers := fk2.GetAllWorkers()
	if len(workers) > 0 {
		t.Errorf("ProcessQueue pushed recipe to queue even if client is stopped!")
	}
}
