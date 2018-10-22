package workers

type queue struct {
	elements []*Worker
}

func (q queue) len() int {
	return len(q.elements)
}

func (q *queue) empty() {
	q.elements = q.elements[:0]
}

func (q queue) isEmpty() bool {
	if q.len() < 1 {
		return true
	}

	return false
}

// pop first and remove it from stack
func (q *queue) pop() *Worker {
	var w *Worker

	q.elements, w = q.elements[1:], q.elements[0]

	return w
}

func (q *queue) push(w *Worker) {
	// ** TODO ** : Introducing priorities on tasks will need to push worker by priorities and reorder elements slice
	q.elements = append(q.elements, w)
}
