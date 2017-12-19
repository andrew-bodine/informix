package downstream_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestDownstream(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Downstream Suite")
}
