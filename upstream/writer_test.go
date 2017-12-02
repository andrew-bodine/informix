package upstream_test

type chanWriter struct {
    downstream  chan interface{}
}

// Implement the io.Writer interface.
func (cw *chanWriter) Write(b []byte) (int, error) {
    cw.downstream <- b

    return len(b), nil
}
