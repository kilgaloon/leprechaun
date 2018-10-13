package workers

import (
	"testing"

	"github.com/kilgaloon/leprechaun/context"
	"github.com/kilgaloon/leprechaun/log"
	"github.com/kilgaloon/leprechaun/recipe"
)

var (
	workers2 = New(
		1,
		"../tests/var/log/leprechaun/workers.output",
		log.Logs{},
		context.New(),
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

	// try to kill worker
	worker.Kill()

	//os.Remove("../tests/test.txt")

}
