package workers

import (
	"bytes"
	"os/exec"
	"strings"
	"time"

	"github.com/kilgaloon/leprechaun/context"
	"github.com/kilgaloon/leprechaun/log"
)

// AsyncMarker is string in step that we use to know
// does that command need to be done async
const AsyncMarker = "->"

// Worker is single worker and information about it
// worker is processing all steps
type Worker struct {
	StartedAt      time.Time
	WorkingOn      string
	Steps          []string
	Context        *context.Context
	Logs           log.Logs
	DoneChan       chan string
	Name           string
	TasksPerformed int
}

// Run starts worker
//
// TODO: Worker should send information back to some channel
// this can be used to write informations to {job}.lock file
func (w *Worker) Run(steps []string) {
	w.Steps = steps

	for _, step := range w.Steps {
		w.Logs.Info("Step %s is in progress... \n", step)
		// replace variables
		parts := strings.Fields(step)

		if parts[0] == AsyncMarker {
			step = w.Context.Transpile(strings.Join(parts[1:], " "))
			go w.workOnStep(step)
		} else {
			step = w.Context.Transpile(step)
			w.workOnStep(step)
		}
	}
}

func (w *Worker) workOnStep(step string) {
	cmd := exec.Command("bash", "-c", step)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	w.WorkingOn = step
	if err != nil {
		w.Logs.Info("Step %s failed to start. Reason: %s \n", step, stderr.String())
		w.WorkingOn = ""
	}

	w.Logs.Info("Step %s finished... \n\n", step)
	// there is output, write it to info
	if len(out.String()) > 0 {
		w.Logs.Info("Step %s -> output: %s", step, out.String())
	}

	w.Done()
}

// Done signals that this worker is done and send his id for cleaner
func (w *Worker) Done() {
	w.TasksPerformed++
	// worker performed all tasks, and can be done
	if w.TasksPerformed == len(w.Steps) {
		w.DoneChan <- w.Name
	}
}
