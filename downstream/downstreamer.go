// The downstream package contains necessary tools for interacting with
// a downstream remote.
package downstream

// A Downstreamer provides an interface for managing and sending on a downstream
// client.
type Downstreamer interface {

    // Connect instructs the downstreamer implementation that it should do
    // any necessary initial connection setup.
    Connect() error

    // Publish tells a downstreamer to send data downstream.
    Publish(string, map[string]interface{}) error
}
