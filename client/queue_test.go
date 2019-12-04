package client

import (
	"testing"

	"github.com/kilgaloon/leprechaun/daemon"
)

func TestBuildQueue(t *testing.T) {
	def.BuildQueue()

	def.Lock()
	if len(def.Queue.Stack) < 7 {
		t.Errorf("Queue stack length expected to be 7, got %d", len(def.Queue.Stack))
	}
	def.Unlock()

}

func TestQueue(t *testing.T) {
	// reset queue to 0 to test AddToQueue
	def.Lock()
	q := def.Queue
	q.Stack = q.Stack[:0]
	def.Unlock()

	def.AddToQueue(def.GetConfig().GetRecipesPath() + "/schedule.yml")

	def.Lock()
	if len(def.Queue.Stack) != 1 {
		t.Errorf("Queue stack length expected to be 1, got %d", len(def.Queue.Stack))
	}
	def.Unlock()

	def.AddToQueue(def.GetConfig().GetRecipesPath() + "/hook.yml")

	def.Lock()
	if len(def.Queue.Stack) != 1 {
		t.Errorf("Queue stack length expected to be 0, got %d", len(def.Queue.Stack))
	}
	def.Unlock()

	def.Pause()
	def.ProcessQueue()

	def.SetStatus(daemon.Started)
	def.ProcessQueue()
}

func TestFindInRecipe(t *testing.T) {
	// reset queue to 0 to test AddToQueue
	if def.FindRecipe("schedule") == nil {
		t.Fatal("Schedule recipe doesn't exist")
	}

	if def.FindRecipe("random_name") != nil {
		t.Fatal("Random name recipe should not exist in recipe queue")
	}
}
