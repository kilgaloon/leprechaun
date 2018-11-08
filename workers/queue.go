package workers

import "sync"

// Queue holds list of workers that are in queue
type Queue struct {
	elements []*Worker
	sync.Mutex
}

func (q *Queue) len() int {
	return len(q.elements)
}

func (q *Queue) empty() {
	q.elements = q.elements[:0]
}

func (q *Queue) isEmpty() bool {
	if q.len() < 1 {
		return true
	}

	return false
}

// pop first and remove it from stack
func (q *Queue) pop() *Worker {
	q.Lock()
	defer q.Unlock()

	var w *Worker
	q.elements, w = q.elements[1:], q.elements[0]

	return w
}

func (q *Queue) push(w *Worker) {
	q.Lock()
	// ** TODO ** : Introducing priorities on tasks will need to push worker by priorities and reorder elements slice
	q.elements = append(q.elements, w)
	q.Unlock()
}

// NewQueue creates new queue for workers
// func NewQueue() Queue {
// 	return Queue{
// 		RWMutex: make(sync.RWMutex),
// 	}
// }
