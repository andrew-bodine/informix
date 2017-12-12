package analytics_test

import (
    "net/http"
)

// An implementation of http.ResponseWriter for testing purposes.
type MockResponseWriter struct {
    Buf     string
}

// Implement the http.ResponseWriter interface.
func (m *MockResponseWriter) Header() http.Header {
    return nil
}

// Implement the http.ResponseWriter interface.
func (m *MockResponseWriter) Write(bs []byte) (int, error) {
    m.Buf = string(bs)

    return len(bs), nil
}

// Implement the http.ResponseWriter interface.
func (m *MockResponseWriter) WriteHeader(h int) {}
