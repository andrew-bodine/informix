package upstream

import (
    "io"
    "net"
    "sync"
)

type writerChan struct {
    downstream  chan interface{}
}

// Implement the io.Writer interface.
func (w *writerChan) Write(b []byte) (int, error) {
    w.downstream <- b

    return len(b), nil
}

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

    downstream  *writerChan
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
func (s *Socket) Open(address string, downstream chan interface{}) error {
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

        // Wrap the downstream channel with an io.Writer implementation to
        // make streaming from socket to channel really simple.
        s.downstream = &writerChan{
            downstream: downstream,
        }

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
