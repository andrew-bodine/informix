package cache

import (
	"sync"
)

func NewCache() Cacher {
	return &cache{
		directives: make(map[string]*Directive),
	}
}

// An implementation of the Cacher interface.
type cache struct {
	directives map[string]*Directive

	sync.Mutex
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
	c.Lock()
	defer c.Unlock()

	// If we already have this directive, then we are already
	// processing it. We simply ignore this request.
	if c.directives[dir.Key] != nil {
		return
	}

	// Store in cacher directives, and run it.
	c.directives[dir.Key] = dir
	c.directives[dir.Key].Run()
}

// Implement the Cacher interface.
func (c *cache) Unregister(dir *Directive) {
	c.Lock()
	defer c.Unlock()

	// If we don't have a directive matching the provided key,
	// then ignore this request.
	if c.directives[dir.Key] == nil {
		return
	}

	// Explicitly stop the directive, giving it a chance to do any cleanup.
	c.directives[dir.Key].Stop()

	// Clear from cacher directives.
	delete(c.directives, dir.Key)
}
