package upstream

import (
    "io"
    "sync"
)

func NewMemory() Upstreamer {
    return &Memory{
        state:  CLOSED,
    }
}

// An upstream interface for receiving data from peer goroutines in the same
// address space.
type Memory struct {
    state       int

    // Goroutine to pump the upstream channel, thus always allowing upstream
    // peers to send data.
    upstream    chan []byte
    closer      chan bool

    downstream  io.Writer

    sync.Mutex
}

// Implement the Upstreamer interface.
func (m *Memory) State() int {
    m.Lock()
    defer m.Unlock()

    return m.state
}

// Implement the Upstreamer interface.
func (m *Memory) Open(address string, downstream io.Writer) error {
    m.Lock()
    defer m.Unlock()

    if m.state == OPEN {
        return nil
    }

    m.state = OPEN

    m.upstream = make(chan []byte)
    m.closer = make(chan bool, 1)

    if downstream != nil {
        m.downstream = downstream
        go m.stream()
    }

    return nil
}

// stream pumps the upstream channel, and forwards any data downstream.
func (m *Memory) stream() {
    for {
        select {
        case data := <- m.upstream:
            m.downstream.Write(data)
        case <- m.closer:
            m.closer <- true
            return
        }
    }
}

// Upstream returns a reference to the upstream channel that this Memory
// instance is pumping downstream.
func (m *Memory) Upstream() chan []byte {
    m.Lock()
    defer m.Unlock()

    return m.upstream
}

// Implement the Upstreamer interface.
func (m *Memory) Close() error {
    if m.State() == CLOSED {
        return nil
    }

    m.Lock()
    m.state = CLOSED
    m.Unlock()

    // Notify stream goroutine that it should stop streaming data downstream.
    if m.closer != nil {
        m.closer <- true
    }

    // Wait for callback from stream goroutine that it has exited.
    <- m.closer

    m.Lock()
    close(m.closer)
    m.closer = nil

    close(m.upstream)
    m.upstream = nil
    m.Unlock()

    return nil
}
