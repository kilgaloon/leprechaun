package workers

import (
	"bytes"
	"os"
	"sync"
	"time"

	"github.com/kilgaloon/leprechaun/context"
	"github.com/kilgaloon/leprechaun/log"
	"github.com/kilgaloon/leprechaun/recipe"
)

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
	Cmds           map[string]*Cmd
	Stdout         *os.File
	Recipe         recipe.Recipe
	Err            error
	mu             *sync.RWMutex
	Debug          bool
}

// Run starts worker
func (w *Worker) Run() {
	w.StartedAt = time.Now()

	steps := w.Recipe.GetSteps()
	for i, step := range steps {
		w.Logs.Info("Step %s is in progress... \n", step)

		s := Step(step)
		if !s.Validate() {
			return
		}

		step := s.Plain()
		w.mu.Lock()
		var in bytes.Buffer
		if i > 0 {
			prevStep := Step(steps[i-1])
			if val, ok := w.Cmds[prevStep.Plain()]; ok {
				if val.pipe {
					in = val.Stdout
				}
			}
		}

		var cmd *Cmd
		var err error
		if s.IsRemote() {
			cmd, err = NewRemoteCmd(s, &in, w.Context, w.Debug)
		} else {
			shell := w.Context.GetVar("shell").GetValue().(string)
			cmd, err = NewCmd(s, &in, nil, w.Debug, shell)
		}

		if err != nil {
			w.Logs.Error(err.Error())
		}

		w.WorkingOn = step
		w.Cmds[step] = cmd
		w.mu.Unlock()

		// Pipe override Async
		// -> echo "Something" }>
		// will not be executed async because we wan't to pass
		// output to next step, if this task start async then next step
		// will start and output won't be passed to it
		if s.IsAsync() && !s.IsPipe() && s.CanError() {
			go w.workOnStep(cmd)
		} else {
			err = w.workOnStep(cmd)
			// there was error with step and step can't error
			// we break loop of step linear execution
			if err != nil && !s.CanError() {
				break;
			}
		}
	}
}

func (w *Worker) workOnStep(cmd *Cmd) (err error) {
	err = cmd.Run()

	if err != nil {
		w.mu.Lock()

		w.Err = err
		w.Recipe.Err = err

		w.ErrorChan <- *w

		w.Logs.Error(w.Err.Error())

		w.mu.Unlock()
		return
	}

	w.Logs.Info("Step %s finished \n\n", cmd.Step.Plain())

	w.Done()

	return
}

// Kill all commands that worker is working on
func (w *Worker) Kill() {
	w.mu.Lock()
	defer w.mu.Unlock()

	for step, cmd := range w.Cmds {
		if err := cmd.Cmd.Process.Kill(); err != nil {
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
