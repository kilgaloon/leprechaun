package workers

import (
	"bytes"
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
		true,
	)
	r, _              = recipe.Build("../tests/etc/leprechaun/recipes/schedule.yml")
	canErrorRecipe, _ = recipe.Build("../tests/etc/leprechaun/recipes/schedule_canerror.yml")
	worker, errr      = workers2.CreateWorker(r)
	worker2, _        = workers2.CreateWorker(canErrorRecipe)
)

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
		t.Fatalf("No worker poped from queue")
	}
}

func TestWorkerErrorStep(t *testing.T) {
	steps := canErrorRecipe.GetSteps()
	for _, step := range steps {
		s := Step(step)
		if !s.Validate() {
			return
		}

		var cmd *Cmd
		var err error
		var in bytes.Buffer
		cmd, err = NewCmd(s, &in, nil, true, "bash")

		if err != nil {
			t.Fatalf("Creating NewCmd failed")
		}

		// Pipe override Async
		// -> echo "Something" }>
		// will not be executed async because we wan't to pass
		// output to next step, if this task start async then next step
		// will start and output won't be passed to it
		if s.IsAsync() && !s.IsPipe() && s.CanError() {
			go worker2.workOnStep(cmd)
		} else {
			err = worker2.workOnStep(cmd)
			// there was error with step and step can't error
			// we break loop of step linear execution
			if err == nil && s.CanError() {
				t.Fatal(err)
			} else {
				break;
			}
		}
	}
}
