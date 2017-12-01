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

    listener    net.Listener

    downstream  io.Writer
    conns       []net.Conn

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
        conn, err := s.listener.Accept()
        if err != nil || s.State() == CLOSED {

            // Stop all io.Copy blocking operations currently ongoing.
            for _, c := range s.conns {
                c.Close()
            }

            return
        }

        go func(c net.Conn) {
            s.Lock()
            s.conns = append(s.conns, c)
            s.Unlock()

            _, _ = io.Copy(s.downstream, c)
        }(conn)
    }
}

// Implement the Upstreamer interface.
func (s *Socket) Close() error {
    s.Lock()
    if s.state == CLOSED {
        return nil
    }

    s.state = CLOSED
    s.Unlock()

    if err := s.listener.Close(); err != nil {
        return err
    }
    s.listener = nil

    return nil
}
