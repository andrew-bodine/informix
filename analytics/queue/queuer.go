package queue

import (
	"io"
)

// A Queuer is responsible for maintaining a store of any type of objects
// in a FIFO manner.
type Queuer interface {

	// Size reports the current capacity of the queue.
	Size() int

	// Count reports the number of items in the queue.
	Count() int

	// Copy returns a snapshot of current items in the queue.
	Copy() []interface{}

	// Push appends an item to the end of the queue.
	Push(interface{})

	// Drain clears all items from the queue.
	Drain()

	// TODO: Resize(int)

	io.Writer
}
