package agent

import (
	"io"
	"time"
)

// WorkersList is default command for agents
func (d Default) WorkersList(r io.Writer, args ...string) ([][]string, error) {
	workers := d.GetWorkers()
	resp := [][]string{}

	if workers.Size() < 1 {
		resp = [][]string{
			{"No workers currently working!"},
		}

		return resp, nil
	}

	for name, worker := range workers.GetAll() {
		startedAt := worker.StartedAt.Format(time.UnixDate)
		resp = append(resp, []string{name, startedAt, worker.WorkingOn})
	}

	return resp, nil
}

// KillWorker kills worker by provided name
func (d Default) KillWorker(r io.Writer, args ...string) ([][]string, error) {
	workers := d.GetWorkers()
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
