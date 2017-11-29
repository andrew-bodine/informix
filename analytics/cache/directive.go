package cache

import (
	"sync"

	"github.com/andrew-bodine/informer/analytics/queue"
)

// A Directive is meant as a wrapper type signifying the intent to cache data
// from the source. The lifecycle of a Directive can be represented by the
// following FSM:
//
//         v""""""""""\
// (0)--->(1)        (2)
//         \,,,,,,,,,,^
//
// 1 - Running
// 2 - Stopped

type Directive struct {

	// The implementation assumes that provided keys are unique.
	Key string

	// The backing store implementation meant to store data from the source.
	// If Queuer isn't initialized before registration, an appropriate one
	// will be created and associated here.
	Queuer queue.Queuer

	Source chan interface{}

	// Channel to stop the goroutine consuming the source.
	closer chan bool

	sync.Mutex
}

// TODO:
func (d *Directive) Run() {
	d.Lock()
	defer d.Unlock()

	// If directive is already running, don't start another goroutine.
	if d.closer != nil {
		return
	}

	d.closer = make(chan bool)

	go d.pump()
}

func (d *Directive) pump() {
	for {
		select {
		case data := <-d.Source:
			d.Queuer.Push(data)
			continue
		case <-d.closer:
			return
		}
	}
}

// TODO:
func (d *Directive) Stop() {
	d.Lock()
	defer d.Unlock()

	// If directive isn't running, don't try to stop it.
	if d.closer == nil {
		return
	}

	// Notify pump() goroutine that it should exit.
	d.closer <- true

	close(d.closer)
	d.closer = nil
}
