package client

import (
	"testing"
)

func TestBuildQueue(t *testing.T) {
	fakeClient.BuildQueue()

	if len(fakeClient.Queue.Stack) != 1 {
		t.Errorf("Queue stack length expected to be 1, got %d", len(fakeClient.Queue.Stack))
	}

	// reset queue to 0 to test AddToQueue
	fakeClient.Queue.Stack = fakeClient.Queue.Stack[:0]
}

func TestAddToQueue(t *testing.T) {
	fakeClient.AddToQueue(&fakeClient.Queue.Stack, fakeClient.Config.RecipesPath+"/schedule.yml")
	if len(fakeClient.Queue.Stack) != 1 {
		t.Errorf("Queue stack length expected to be 1, got %d", len(fakeClient.Queue.Stack))
	}

	fakeClient.AddToQueue(&fakeClient.Queue.Stack, fakeClient.Config.RecipesPath+"/hook.yml")
	if len(fakeClient.Queue.Stack) != 1 {
		t.Errorf("Queue stack length expected to be 0, got %d", len(fakeClient.Queue.Stack))
	}
}

// func TestProcessQueue(t *testing.T) {
// 	now := time.Now()
// 	compare := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, time.UTC)

// 	for index, r := range fakeClient.Queue.Stack {
// 		recipe := &fakeClient.Queue.Stack[index]

// 	}
// }
