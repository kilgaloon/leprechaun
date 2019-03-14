package workers

import (
	"errors"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/kilgaloon/leprechaun/context"
	"github.com/kilgaloon/leprechaun/log"
	"github.com/kilgaloon/leprechaun/notifier"
	"github.com/kilgaloon/leprechaun/notifier/notifications"
	"github.com/kilgaloon/leprechaun/recipe"
)

var (
	// ErrWorkerNotExist is error when worker doesn't exist in stack
	ErrWorkerNotExist = errors.New("No worker with that name")
	// ErrStillActive is error when in some cases when worker is created and
	// worker with that name still exists (working on something)
	// worker get their names from recipe names, so basically some recipe can't be run twice
	ErrStillActive = errors.New("Worker still active")
	// ErrMaxReached is error that says you that no more workers
	// is allowed and this is specified in config
	ErrMaxReached = errors.New("Maximum allowed workers reached, worker moved to queue")
	// ErrMaxQueueReached is error that says you that no more workers
	// is allowed in queue and this is specified in config
	ErrMaxQueueReached = errors.New("Maximum allowed workers in queue reached, worker disposed")
)

// Config defines interface which we use to build workers struct
type Config interface {
	GetMaxAllowedWorkers() int
	GetMaxAllowedQueueWorkers() int
	GetWorkerOutputDir() string
	notifier.Config
}

// Workers hold everything about workers
type Workers struct {
	stack map[string]Worker
	Queue
	allowedSize      int
	allowedQueueSize int
	OutputDir        string
	Context          *context.Context
	Logs             log.Logs
	DoneChan         chan string
	ErrorChan        chan Worker
	*notifier.Notifier
	*sync.RWMutex
}

// NumOfWorkers returns size of stack/number of workers
func (w *Workers) NumOfWorkers() int {
	w.Lock()
	defer w.Unlock()

	return len(w.stack)
}

// PushToStack places worker on stack
func (w *Workers) PushToStack(worker *Worker) {
	w.Lock()
	defer w.Unlock()

	w.stack[worker.Recipe.GetName()] = *worker
}

// GetAllWorkers workers from stack
func (w Workers) GetAllWorkers() map[string]Worker {
	return w.stack
}

// GetWorkerByName gets worker by provided name
func (w *Workers) GetWorkerByName(name string) (*Worker, error) {
	w.Lock()
	defer w.Unlock()

	var worker Worker
	if worker, ok := w.stack[name]; ok {
		return &worker, nil
	}

	return &worker, ErrWorkerNotExist
}

// DeleteWorkerByName Removes worker from stack
func (w *Workers) DeleteWorkerByName(name string) {
	_, err := w.GetWorkerByName(name)
	if err == nil {
		w.Lock()
		delete(w.stack, name)
		w.Unlock()
	}
}

// CreateWorker Create single worker if number is not exceeded and move it to stack
func (w *Workers) CreateWorker(r *recipe.Recipe) (*Worker, error) {
	if _, ok := w.GetWorkerByName(r.GetName()); ok == nil {
		return nil, ErrStillActive
	}

	worker := &Worker{
		StartedAt: time.Now(),
		Context:   w.Context,
		Logs:      w.Logs,
		DoneChan:  w.DoneChan,
		ErrorChan: w.ErrorChan,
		Recipe:    r,
		Cmd:       make(map[string]*exec.Cmd),
		mu:        new(sync.RWMutex),
	}

	var err error
	worker.Stdout, err = os.Create(w.OutputDir + "/" + worker.Recipe.GetName() + ".out") // For read access.
	if err != nil {
		w.Logs.Error("%s", err)
	}

	if w.NumOfWorkers() < w.allowedSize {
		// move to stack
		w.PushToStack(worker)

		w.Logs.Info("Worker with NAME: %s created", worker.Recipe.GetName())

		return worker, nil
	}

	if w.Queue.len() < w.allowedQueueSize {
		w.Queue.push(worker)
	} else {
		w.Logs.Error("%s", ErrMaxQueueReached)
		return nil, ErrMaxQueueReached
	}

	w.Logs.Error("%s", ErrMaxReached)
	return nil, ErrMaxReached
}

func (w Workers) listener() {
	go func() {
		for {
			select {
			case workerName := <-w.DoneChan:
				// When worker is done, check in worker queue is there any to process
				// ** TODO ** : Since we plan to introduce priority now everything is same priority,
				// tasks in queue will need to wait in queue until all higher priority tasks are done
				if !w.Queue.isEmpty() {
					worker := w.Queue.pop()
					w.PushToStack(worker)

					go worker.Run()
				}
				w.DeleteWorkerByName(workerName)
				w.Logs.Info("Worker with NAME: %s cleaned", workerName)
			case worker := <-w.ErrorChan:
				// send notifications
				go w.NotifyWithOptions(notifications.Options{
					Body: "Your recipe '" + worker.Recipe.GetName() +
						"' failed on step '" + worker.WorkingOn +
						"' because of error '" + worker.Err.Error() + "'",
				})
				// when worker gets to error, log it
				// and delete it from stack of workers
				// otherwise it will populate stack and pretend to be active
				w.DeleteWorkerByName(worker.Recipe.GetName())
				w.Logs.Error("Worker %s: %s", worker.Recipe.GetName(), worker.Err)
			}
		}
	}()
}

// New create Workers struct instance
func New(cfg Config, logs log.Logs, ctx *context.Context) Workers {
	workers := Workers{
		stack:            make(map[string]Worker),
		allowedSize:      cfg.GetMaxAllowedWorkers(),
		allowedQueueSize: cfg.GetMaxAllowedQueueWorkers(),
		Logs:             logs,
		Context:          ctx,
		DoneChan:         make(chan string),
		ErrorChan:        make(chan Worker),
		OutputDir:        cfg.GetWorkerOutputDir(),
		Notifier:         notifier.New(cfg, logs),
		RWMutex:          new(sync.RWMutex),
	}
	// listener listens for varius events coming from workers, currently those are
	// done and errors
	workers.listener()

	return workers
}
