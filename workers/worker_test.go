package workers

import (
	"os"
	"testing"

	"github.com/kilgaloon/leprechaun/context"
	"github.com/kilgaloon/leprechaun/log"
)

var (
	workers2 = New(
		1,
		"../tests/var/log/leprechaun/workers.output",
		log.Logs{},
		context.New(),
	)
	worker, err = workers2.CreateWorker("test")
)

func TestRun(t *testing.T) {
	steps := []string{"echo 'test output to file' > ../tests/test.txt"}
	//var wg sync.WaitGroup
	//wg.Add(1)

	worker.Run(steps)

	steps = []string{"-> echo 'test output to file' > ../tests/test.txt"}
	go worker.Run(steps)
	//wg.Wait()

	// try to kill worker
	worker.Kill()

	os.Remove("../tests/test.txt")

}
