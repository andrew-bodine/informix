// The upstream package contains necessaries for managing various system
// level inbound data interfaces.
package upstream

// An Upstreamer manages an upstream interface.
type Upstreamer interface {

    // Upstreamers are designed to have two states. These two states are closed
    // and open, or 0 and 1 respectively.
    State() int

    // Open allows callers to instruct an Upstreamer to open it's underlying
    // upstream interface at the provided address, and stream data to the
    // provided channel.
    Open(string, chan interface{}) error

    // Close allows callers to tell an Upstreamer to close it's underlying
    // upstream interface.
    Close() error
}
