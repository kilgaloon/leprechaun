package workers

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/kilgaloon/leprechaun/context"
	"github.com/kilgaloon/leprechaun/log"
	"github.com/kilgaloon/leprechaun/recipe"
)

// AsyncMarker is string in step that we use to know
// does that command need to be done async
const AsyncMarker = "->"

// Worker is single worker and information about it
// worker is processing all steps
type Worker struct {
	StartedAt      time.Time
	WorkingOn      string
	Context        *context.Context
	Logs           log.Logs
	DoneChan       chan string
	ErrorChan      chan Worker
	TasksPerformed int
	Cmd            map[string]*exec.Cmd
	steps          []*Cmd
	Stdout         *os.File
	Recipe         *recipe.Recipe
	Err            error
	mu             *sync.RWMutex
}

// Run starts worker
func (w *Worker) Run() {
	w.StartedAt = time.Now()

	for i, step := range w.Recipe.GetSteps() {
		w.Logs.Info("Step %s is in progress... \n", step)

		if len(step) < 1 {
			continue
		}

		parts := strings.Fields(step)
		if parts[0] == AsyncMarker {
			step = w.Context.Transpile(strings.Join(parts[1:], " "))
			go w.workOnStep(i, step)
		} else {
			w.workOnStep(i, step)
		}
	}
}

func (w *Worker) workOnStep(i int, step string) {
	w.mu.Lock()

	var in bytes.Buffer
	if len(w.steps) > 0 {
		ps := w.steps[i-1]
		if ps.pipe {
			in = ps.Stdout
		}
	}

	cmd, err := NewCmd(step, &in)
	if err != nil {
		w.Logs.Error(err.Error())
	}

	w.WorkingOn = step
	w.Cmd[step] = cmd.cmd
	w.mu.Unlock()

	err = cmd.Run()

	w.mu.Lock()
	w.steps = append(w.steps, cmd)
	w.mu.Unlock()

	if err != nil {
		w.mu.Lock()

		w.Err = err
		w.Recipe.Err = err

		w.ErrorChan <- *w

		w.Logs.Error(w.Err.Error())

		w.mu.Unlock()
		return
	}

	w.Logs.Info("Step %s finished... \n\n", step)
	// command finished executing
	// delete it, and let it rest in pepperonies
	w.mu.Lock()
	delete(w.Cmd, step)
	w.mu.Unlock()

	w.Done()
}

// Kill all commands that worker is working on
func (w *Worker) Kill() {
	w.mu.Lock()
	defer w.mu.Unlock()

	for step, cmd := range w.Cmd {
		if err := cmd.Process.Kill(); err != nil {
			w.Logs.Error("Failed to kill process on step %s: %s", step, err)
			w.Err = err

			return
		}
	}

	w.DoneChan <- w.Recipe.GetName()
}

// Done signals that this worker is done and send his id for cleaner
func (w *Worker) Done() {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.TasksPerformed++
	// worker performed all tasks, and can be done
	if w.TasksPerformed == len(w.Recipe.GetSteps()) {
		w.DoneChan <- w.Recipe.GetName()
	}
}
