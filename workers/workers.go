package workers

import (
	"errors"
	"os/exec"
	"time"

	"github.com/kilgaloon/leprechaun/context"
	"github.com/kilgaloon/leprechaun/log"
)

// Errors hold possible errors that can happen on worker
type Errors struct {
	StillActive           error
	AllowedWorkersReached error
}

// Workers hold everything about workers
type Workers struct {
	stack       map[string]Worker
	allowedSize int
	Context     *context.Context
	Logs        log.Logs
	DoneChan    chan string
	Errors
}

// Size returns size of stack/number of workers
func (w Workers) Size() int {
	return len(w.stack)
}

// GetAll workers from stack
func (w Workers) GetAll() map[string]Worker {
	return w.stack
}

// GetByName gets worker by provided name
func (w Workers) GetByName(name string) (Worker, error) {
	var worker Worker
	if worker, ok := w.stack[name]; ok {
		return worker, nil
	}

	return worker, errors.New("No worker with that name")
}

// CreateWorker Create single worker if number is not exceeded and move it to stack
func (w *Workers) CreateWorker(name string) (*Worker, error) {
	if _, ok := w.stack[name]; ok {
		return nil, w.Errors.StillActive
	}

	if w.Size() < w.allowedSize {
		worker := &Worker{
			StartedAt: time.Now(),
			Context:   w.Context,
			Logs:      w.Logs,
			DoneChan:  w.DoneChan,
			Name:      name,
			Cmd:       make(map[string]*exec.Cmd),
		}
		// move to stack
		w.stack[worker.Name] = *worker

		w.Logs.Info("Worker with NAME: %s created", worker.Name)

		return worker, nil
	}

	return nil, w.Errors.AllowedWorkersReached
}

// Cleaner recieves signal from DoneChan and clean workers that are done
func (w Workers) Cleaner() {
	go func() {
		for {
			select {
			case workerName := <-w.DoneChan:
				delete(w.stack, workerName)
				w.Logs.Info("Worker with NAME: %s cleaned", workerName)
			}
		}
	}()
}

// BuildWorkers create Workers struct instance
func BuildWorkers(ctx *context.Context, allowedSize int, logs log.Logs) *Workers {
	workers := &Workers{
		stack:       make(map[string]Worker),
		allowedSize: allowedSize,
		Logs:        logs,
		Context:     ctx,
		DoneChan:    make(chan string),
		Errors: Errors{
			StillActive:           errors.New("Worker still active"),
			AllowedWorkersReached: errors.New("Maximum allowed workers reached"),
		},
	}
	// cleaner sits and wait to clean workers that are done with their job
	workers.Cleaner()

	return workers
}
