package upstream

import (
    "net"
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
}

// Implement the Upstreamer interface.
func (s *Socket) State() int {
    return s.state
}

// Implement the Upstreamer interface.
func (s *Socket) Open() error {
    if s.state == OPEN {
        return nil
    }

    // Open underlying UNIX socket at SOCK path.
    l, err := net.Listen("unix", SOCK)
    if err != nil {
        return err
    }

    s.listener = l
    s.state = OPEN

    return nil
}

// Implement the Upstreamer interface.
func (s *Socket) Close() error {
    if s.state == CLOSED {
        return nil
    }

    if err := s.listener.Close(); err != nil {
        return err
    }

    s.listener = nil
    s.state = CLOSED

    return nil
}
