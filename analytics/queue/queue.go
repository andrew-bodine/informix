package queue

import (
    "sync"
)

func NewQueue(size int) Queuer {
    return &queue{
        size:   size,
        count:  0,

        head:   nil,
        tail:   nil,
    }
}

type item struct {
    data   interface{}
    next    *item
}

// An implementation of the Queuer interface, backed by memory.
type queue struct {
    size    int
    count   int

    head    *item
    tail    *item

    sync.Mutex
}

func (q *queue) Size() int {
    return q.size
}

func (q *queue) Count() int {
    return q.count
}

// Push appends the item to the end of the queue.
func (q *queue) Push(data interface{}) interface{} {
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
        return nil
    case q.size:
        removed := q.head
        q.tail.next = i
        q.tail = q.tail.next
        q.head = q.head.next
        return removed.data
    default:
        q.tail.next = i
        q.tail = q.tail.next
        q.count++
        return nil
    }
}
