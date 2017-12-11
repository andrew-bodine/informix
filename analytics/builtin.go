package analytics

import (
    "time"

    "github.com/andrew-bodine/informix/analytics/emit"
    "github.com/andrew-bodine/informix/analytics/queue"
)

type Builtin interface {
    Cache(string) []interface{}

    Run(time.Duration)

    Stop()
}

func NewBuiltin() Builtin {
    b := &builtin{
        timer:      make(chan *time.Timer, 1),
        stop:       make(chan bool, 1),
        generators: map[string]emit.Generator{},
        queuers:    map[string]queue.Queuer{},
    }

    // Initialize synchronization target.
    b.timer <- nil

    // Setup builtin analytic directives.
    b.generators[emit.MEMORY] = emit.Memory()
    b.queuers[emit.MEMORY] = queue.NewQueue(1)

    return b
}

// Batteries included implementation that manages all built-in
// analytics for Informix.
type builtin struct {
    timer       chan *time.Timer
    stop        chan bool

    generators  map[string]emit.Generator
    queuers     map[string]queue.Queuer
}

func (b *builtin) Cache(key string) []interface{} {
    q, exists := b.queuers[key]

    if !exists {
        return []interface{}{}
    }

    return q.Copy()
}

// Implement the Builtin interface.
func (b *builtin) Run(d time.Duration) {
    t := <- b.timer

    // If the timer isn't nil here, it means this Builtin instance is
    // already running, nothing to do.
    if t != nil {
        b.timer <- t
        return
    }

    t = time.NewTimer(d)
    b.timer <- t

    go func(t *time.Timer) {
        defer func() {
            <- b.timer
            b.timer <- nil
        }()

        for {
            select {
            case <- t.C:
                b.refresh()

                break
            case <- b.stop:
                return
            }

            t.Reset(d)
        }
    }(t)
}

// refresh iterates through and caches the builtin analytic generators.
func (b *builtin) refresh() {
    for k, g := range b.generators {
        v := g.Generate()

        if v == nil {
            continue
        }

        b.queuers[k].Push(v)
    }
}

// Implement the Builtin interface.
func (b *builtin) Stop() {
    b.stop <- true

    for _, q := range b.queuers {
        q.Drain()
    }
}
