package emit

import (
    linuxproc "github.com/c9s/goprocinfo/linux"
)

const (
    MEMORY = "memory"
)

// An Emitter implementation for memory information in linux.
type memory struct {
    last    chan *linuxproc.MemInfo
}

// Don't expose the Memory struct because of the channel wrapping the
// last memory stats. Due to the write before ready nature of a channel,
// we need to put something in at creation.
func Memory() *memory {
    m := &memory{
        last:   make(chan *linuxproc.MemInfo, 1),
    }

    m.last <- nil

    return m
}

// Implement the Generator interface.
func (m *memory) Generate() interface{} {
    next, err := linuxproc.ReadMemInfo("/proc/meminfo")
    if err != nil {
        return nil
    }

    _ = <- m.last
    m.last <- next

    return next
}
