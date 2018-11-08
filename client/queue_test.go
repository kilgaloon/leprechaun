package client

import (
	"os"
	"testing"
)

var (
	fk = New("test", cfgWrap.New("test", *path))
)

func TestBuildQueue(t *testing.T) {
	fk.BuildQueue()

	if len(fk.Queue.Stack) != 3 {
		t.Errorf("Queue stack length expected to be 3, got %d", len(fk.Queue.Stack))
	}

	// reset queue to 0 to test AddToQueue
	q := &fk.Queue
	q.Stack = q.Stack[:0]

	clientInfo(t)
	wg.Done()

}

func TestAddToQueue(t *testing.T) {
	fk.AddToQueue(&fk.Queue.Stack, fk.GetConfig().GetRecipesPath()+"/schedule.yml")
	if len(fk.Queue.Stack) != 1 {
		t.Errorf("Queue stack length expected to be 1, got %d", len(fakeClient.Queue.Stack))
	}

	fk.AddToQueue(&fk.Queue.Stack, fk.GetConfig().GetRecipesPath()+"/hook.yml")
	if len(fk.Queue.Stack) != 1 {
		t.Errorf("Queue stack length expected to be 0, got %d", len(fakeClient.Queue.Stack))
	}

	clientInfo(t)
	wg.Done()

}

func clientInfo(t *testing.T) {
	_, err := fakeClient.clientInfo(os.Stdin)
	if err != nil {
		t.Fail()
	}
}
