package queue

import (
    "sync"
)

// Returns the default Queuer implementation internal to this package.
func NewQueue(size int) Queuer {
    return &queue{
        size:   size,
        count:  0,

        head:   nil,
        tail:   nil,
    }
}

// The internal queue implementation relys on a linked list of
// these item objects.
type item struct {
    data   interface{}
    next    *item
}

// An implementation of the Queuer interface, backed by memory.
type queue struct {
    size        int
    count       int

    head        *item
    tail        *item

    sync.Mutex
}

// Implement the Queuer interface.
func (q *queue) Size() int {
    return q.size
}

// Implement the Queuer interface.
func (q *queue) Count() int {
    q.Lock()
    defer q.Unlock()

    return q.count
}

// Implement the Queuer interface.
func (q *queue) Push(data interface{}) {
    q.Lock()
    defer q.Unlock()

    i := &item{
        data:   data,
        next:   nil,
    }

    switch q.count {
    case 0:
        q.tail = i
        q.head = q.tail
        q.count++
        break
    case q.size:
        q.tail.next = i
        q.tail = q.tail.next
        q.head = q.head.next
        break
    default:
        q.tail.next = i
        q.tail = q.tail.next
        q.count++
        break
    }
}

// Implement the Queuer interface.
func (q *queue) Drain() {
    q.Lock()
    defer q.Unlock()

    q.head = nil
    q.tail = nil

    q.count = 0
}
