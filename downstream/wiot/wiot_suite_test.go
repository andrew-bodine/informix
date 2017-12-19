package wiot_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestWiot(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Wiot Suite")
}
