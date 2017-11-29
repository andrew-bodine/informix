package cache

import (
    "sync"
)

func NewCache() Cacher {
    return &cache{
        directives:     make(map[string]*Directive),
    }
}

// An implementation of the Cacher interface.
type cache struct {
    directives  map[string]*Directive

    sync.Mutex
}

// Implement the Cacher interface.
func (c *cache) Run() {
    c.Lock()
    c.Unlock()
}

// Implement the Cacher interface.
func (c *cache) Keys() []string {
    var keys []string

    for key, _ := range c.directives {
        keys = append(keys, key)
    }

    return keys
}

// Implement the Cacher interface.
func (c *cache) Register(dir *Directive) {

    // If we already have this directive, then we are already
    // processing it. We simply ignore this request.
    if c.directives[dir.Key] != nil {
        return
    }

    // If the caller isn't providing a closer channel, add one here.
    if dir.Closer == nil {
        dir.Closer = make(chan bool)
    }

    // Store in cache directives.
    c.Lock()
    c.directives[dir.Key] = dir
    c.Unlock()

    // Kickoff goroutine to pump things from source to queuer.
    go c.pump(c.directives[dir.Key])
}

// TODO:
func (c *cache) pump(dir *Directive) {
    for {
        select {
        case data := <- dir.Source:
            dir.Queuer.Push(data)
            continue
        case <- dir.Closer:
            c.Lock()
            defer c.Unlock()

            return
        }
    }
}

// Implement the Cacher interface.
func (c *cache) Unregister(dir *Directive) {
    c.Lock()
    defer c.Unlock()
}

// Implement the Cacher interface.
func (c *cache) Shutdown() {
    c.Lock()
    defer c.Unlock()
}
