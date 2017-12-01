package queue

// Implement the io.Writer interface.
func (q *queue) Write(b []byte) (int, error) {
	q.Push(b)

	return len(b), nil
}
