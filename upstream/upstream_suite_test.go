package upstream_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestUpstream(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Upstream Suite")
}

// NOTE: Thinking to implement a socket first.

// NOTE: Implement an interface allowing Informix to manage all upstream
// interfaces at once.

// NOTE: Then implement network wrappers (https, udp) that utilize the existing
// socket. This way Informix exposes multiple upstream interfaces, which
// hopefully cover the surface area of where peripheral drivers will be
// executing and sending data from.
