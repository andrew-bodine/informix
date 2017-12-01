package queue

import (
	"encoding/json"
)

// A Queuer is responsible for maintaining a store of any type of objects
// in a FIFO manner.
type Queuer interface {

	// Size reports the current capacity of the queue.
	Size() int

	// TODO: Resize(int)

	// Count reports the number of items in the queue, this is never
	// bigger than Size.
	Count() int

	// Push appends an item to the end of the queue.
	Push(interface{})

	// Drain clears all items from the queue.
	Drain()

	json.Marshaler
}
