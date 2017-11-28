package queue

import (
    "encoding/json"
)

// Public interface for Queuer implementations.
type Queuer interface {
    Size() int
    Count() int

    Push(interface{}) interface{}

    json.Marshaler
}
