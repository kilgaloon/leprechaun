package workers

import (
	"sync"
	"testing"

	"github.com/kilgaloon/leprechaun/config"
	"github.com/kilgaloon/leprechaun/context"
	"github.com/kilgaloon/leprechaun/log"
	"github.com/kilgaloon/leprechaun/recipe"
)

var (
	configs                 = config.NewConfigs()
	ConfigWithSettings      = configs.New("test", "../tests/configs/config_regular.ini")
	ConfigWithQueueSettings = configs.New("test", "../tests/configs/config_test_queue.ini")
	workers2                = New(
		ConfigWithSettings,
		log.Logs{},
		context.New(),
		new(sync.RWMutex),
	)
	r, err       = recipe.Build("../tests/etc/leprechaun/recipes/schedule.yml")
	worker, errr = workers2.CreateWorker(&r)
)

func TestRun(t *testing.T) {
	//steps := []string{"echo 'test output to file' > ../tests/test.txt"}
	//var wg sync.WaitGroup
	//wg.Add(1)

	worker.Run()

	// //steps = []string{"-> echo 'test output to file' > ../tests/test.txt"}
	// go worker.Run()
	//wg.Wait()

	//os.Remove("../tests/test.txt")

}

func TestQueue(t *testing.T) {
	workers2.Queue.empty()

	if !workers2.Queue.isEmpty() {
		t.Fatalf("Queue expected to be empty")
	}

	workers2.Queue.push(worker)

	if workers2.Queue.isEmpty() {
		t.Fatalf("Queue should not be empty")
	}

	w := workers2.Queue.pop()
	if w == nil {
		t.Fail()
	}
}
