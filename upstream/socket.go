package upstream

import (
    "io"
    "net"
)

func NewSocket() Upstreamer {
    s := &Socket{
        state:      make(chan int, 1),
        listener:   make(chan net.Listener, 1),
    }

    s.state <- CLOSED
    s.listener <- nil

    return s
}

// An upstream interface for receiving data from peer services on the same
// system that Informix is running on via a UNIX style socket.
type Socket struct {
    state       chan int

    // Upstream listener and connection store.
    listener    chan net.Listener

    downstream  io.Writer
}

// Implement the Upstreamer interface.
func (s *Socket) State() int {
    st := <- s.state
    s.state <- st

    return st
}

// Implement the Upstreamer interface.
func (s *Socket) Open(address string, downstream io.Writer) error {
    if s.State() == OPEN {
        return nil
    }

    // Open underlying UNIX socket at address, which in this case should
    // be a filesytem path.
    l, err := net.Listen("unix", address)
    if err != nil {
        return err
    }

    <- s.state
    s.state <- OPEN

    <- s.listener
    s.listener <- l

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
        l := <- s.listener
        s.listener <- l

        if l == nil {
            return
        }

        conn, err := l.Accept()

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

    <- s.state
    s.state <- CLOSED

    l := <- s.listener
    s.listener <- nil

    if err := l.Close(); err != nil {
        return err
    }

    return nil
}
