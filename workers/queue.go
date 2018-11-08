package workers

import "sync"

// Queue holds list of workers that are in queue
type Queue struct {
	elements []*Worker
	mu       *sync.Mutex
}

func (q Queue) len() int {
	return len(q.elements)
}

func (q *Queue) empty() {
	q.elements = q.elements[:0]
}

func (q Queue) isEmpty() bool {
	if q.len() < 1 {
		return true
	}

	return false
}

// pop first and remove it from stack
func (q *Queue) pop() *Worker {
	q.mu.Lock()
	defer q.mu.Unlock()

	var w *Worker
	q.elements, w = q.elements[1:], q.elements[0]

	return w
}

func (q *Queue) push(w *Worker) {
	q.mu.Lock()
	// ** TODO ** : Introducing priorities on tasks will need to push worker by priorities and reorder elements slice
	q.elements = append(q.elements, w)
	q.mu.Unlock()
}

// NewQueue creates new queue for workers
func NewQueue() Queue {
	return Queue{
		mu: new(sync.Mutex),
	}
}
