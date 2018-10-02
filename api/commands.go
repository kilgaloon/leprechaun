package api

import (
	"time"
)

// WorkersList is default command for agents
func (r *Registrator) WorkersList(args ...string) ([][]string, error) {
	workers := r.Agent.GetWorkers()
	resp := [][]string{}

	if workers.Size() < 1 {
		resp = [][]string{
			{"No workers currently working!"},
		}
	}

	for name, worker := range workers.GetAll() {
		startedAt := worker.StartedAt.Format(time.UnixDate)
		resp = append(resp, []string{name, startedAt, worker.WorkingOn, worker.Err.Error()})
	}

	return resp, nil
}

// KillWorker kills worker by provided name
func (r *Registrator) KillWorker(args ...string) ([][]string, error) {
	workers := r.Agent.GetWorkers()
	resp := [][]string{}

	worker, err := workers.GetByName(args[0])
	if err != nil {
		resp = [][]string{
			{err.Error()},
		}
	} else {
		worker.Kill()
		resp = [][]string{
			{"Worker killed"},
		}
	}

	return resp, nil
}
