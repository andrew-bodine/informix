package cache

import (
    "github.com/andrew-bodine/informer/analytics/queue"
)

// A Directive is meant as a wrapper type signifying the intent to cache data
// from the source.
type Directive struct {

    // The implementation assumes that provided keys are unique.
    Key         string

    // The backing store implementation meant to store data from the source.
    // If Queuer isn't initialized before registration, an appropriate one
    // will be created and associated here.
    Queuer      queue.Queuer

    Source      chan interface{}

    // The channel by which the caller can stop and dereference the directive.
    Closer      chan bool
}
