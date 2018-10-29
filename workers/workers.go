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
	queue
	allowedSize      int
	allowedQueueSize int
	OutputDir        string
	Context          *context.Context
	Logs             log.Logs
	DoneChan         chan string
	ErrorChan        chan *Worker
	*notifier.Notifier
}

// NumOfWorkers returns size of stack/number of workers
func (w Workers) NumOfWorkers() int {
	return len(w.stack)
}

// PushToStack places worker on stack
func (w *Workers) PushToStack(worker *Worker) {
	mu := new(sync.Mutex)

	mu.Lock()
	defer mu.Unlock()

	w.stack[worker.Recipe.Name] = *worker

	return
}

// GetAllWorkers workers from stack
func (w Workers) GetAllWorkers() map[string]Worker {
	return w.stack
}

// GetWorkerByName gets worker by provided name
func (w Workers) GetWorkerByName(name string) (*Worker, error) {
	var worker Worker
	if worker, ok := w.stack[name]; ok {
		return &worker, nil
	}

	return &worker, ErrWorkerNotExist
}

// CreateWorker Create single worker if number is not exceeded and move it to stack
func (w *Workers) CreateWorker(r *recipe.Recipe) (*Worker, error) {
	mu := new(sync.Mutex)

	mu.Lock()
	defer mu.Unlock()

	if _, ok := w.GetWorkerByName(r.Name); ok == nil {
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
	}

	var err error
	worker.Stdout, err = os.Create(w.OutputDir + "/" + worker.Recipe.Name + ".out") // For read access.
	if err != nil {
		w.Logs.Error("%s", err)
	}

	if w.NumOfWorkers() < w.allowedSize {
		// move to stack
		w.PushToStack(worker)

		w.Logs.Info("Worker with NAME: %s created", worker.Recipe.Name)

		return worker, nil
	}

	if w.queue.len() < w.allowedQueueSize {
		w.queue.push(worker)
	} else {
		w.Logs.Error("%s", ErrMaxQueueReached)
		return nil, ErrMaxQueueReached
	}

	w.Logs.Error("%s", ErrMaxReached)
	return nil, ErrMaxReached
}

func (w *Workers) listener() {
	go func() {
		for {
			select {
			case workerName := <-w.DoneChan:
				// When worker is done, check in worker queue is there any to process
				// ** TODO ** : Since we plan to introduce priority now everything is same priority,
				// tasks in queue will need to wait in queue until all higher priority tasks are done
				if w.queue.len() > 0 {
					worker := w.queue.pop()
					w.PushToStack(worker)

					go worker.Run()
				}

				delete(w.stack, workerName)
				w.Logs.Info("Worker with NAME: %s cleaned", workerName)
			case worker := <-w.ErrorChan:
				// send notifications
				go w.NotifyWithOptions(notifications.Options{
					Body: "Your recipe '" + worker.Recipe.Name +
						"' failed on step '" + worker.WorkingOn +
						"' because of error '" + worker.Err.Error() + "'",
				})
				// when worker gets to error, log it
				// and delete it from stack of workers
				// otherwise it will populate stack and pretend to be active
				delete(w.stack, worker.Recipe.Name)
				w.Logs.Error("Worker %s: %s", worker.Recipe.Name, worker.Err)
			}
		}
	}()
}

// New create Workers struct instance
func New(cfg Config, logs log.Logs, ctx *context.Context) *Workers {
	workers := &Workers{
		stack:            make(map[string]Worker),
		allowedSize:      cfg.GetMaxAllowedWorkers(),
		allowedQueueSize: cfg.GetMaxAllowedQueueWorkers(),
		Logs:             logs,
		Context:          ctx,
		DoneChan:         make(chan string),
		ErrorChan:        make(chan *Worker),
		OutputDir:        cfg.GetWorkerOutputDir(),
		Notifier:         notifier.New(cfg, logs),
	}
	// listener listens for varius events coming from workers, currently those are
	// done and errors
	workers.listener()

	return workers
}
