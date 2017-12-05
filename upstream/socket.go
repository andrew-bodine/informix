package upstream

import (
    "io"
    "net"
    "sync"
)

func NewSocket() Upstreamer {
    return &Socket{
        state:  CLOSED,
    }
}

// An upstream interface for receiving data from peer services on the same
// system that Informix is running on via a UNIX style socket.
type Socket struct {
    state       int

    // Upstream listener and connection store.
    listener    net.Listener

    downstream  io.Writer

    sync.Mutex
}

// Implement the Upstreamer interface.
func (s *Socket) State() int {
    s.Lock()
    defer s.Unlock()

    return s.state
}

// Implement the Upstreamer interface.
func (s *Socket) Open(address string, downstream io.Writer) error {
    s.Lock()
    defer s.Unlock()

    if s.state == OPEN {
        return nil
    }

    // Open underlying UNIX socket at address, which in this case should
    // be a filesytem path.
    l, err := net.Listen("unix", address)
    if err != nil {
        return err
    }

    s.listener = l
    s.state = OPEN

    if downstream != nil {
        s.downstream = downstream

        go s.stream()
    }

    return nil
}

// stream listens for incoming data on the underlying interface, and
// forwards any data downstream.
func (s *Socket) stream() {
    for {
        s.Lock()
        if s.listener == nil {
            return
        }
        s.Unlock()

        conn, err := s.listener.Accept()

        // If there was an error while listening to the socket, or if at the
        // time of a new connection the synchronized state is closed, then
        // exit immediately.
        if err != nil || s.State() == CLOSED {
            return
        }

        go func(c net.Conn) {
            defer c.Close()

            _, _ = io.Copy(s.downstream, c)
        }(conn)
    }
}

// Implement the Upstreamer interface.
func (s *Socket) Close() error {
    if s.State() == CLOSED {
        return nil
    }

    s.Lock()
    s.state = CLOSED
    s.Unlock()

    if err := s.listener.Close(); err != nil {
        return err
    }
    s.Lock()
    s.listener = nil
    s.Unlock()

    return nil
}
