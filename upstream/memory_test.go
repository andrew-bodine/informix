package upstream_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

    "github.com/andrew-bodine/informix/upstream"
)

var _ = Describe("upstream", func() {
    var mem upstream.Upstreamer

    Context("NewMemory()", func() {
        It("creates a Upstreamer in the closed state", func() {
            mem = upstream.NewMemory()
            Expect(mem).ToNot(BeNil())

            _, ok := mem.(upstream.Upstreamer)
            Expect(ok).To(Equal(true))

            Expect(mem.State()).To(Equal(upstream.CLOSED))
        })
    })

    Context("Memory", func() {
        BeforeEach(func() {
            mem = upstream.NewMemory()
        })

        Context("Upstream()", func() {
            var m *upstream.Memory

            Context("when Memory interface is closed", func() {
                It("returns nil", func() {
                    m = mem.(*upstream.Memory)
                    upstream := m.Upstream()
                    Expect(upstream).To(BeNil())
                })
            })

            Context("when Memory interface is open", func() {
                BeforeEach(func() {
                    _ = mem.Open("", nil)
                    m = mem.(*upstream.Memory)
                })

                It("returns the upstream channel interface", func() {
                    m := mem.(*upstream.Memory)
                    upstream := m.Upstream()
                    Expect(upstream).NotTo(BeNil())
                    close(upstream)
                })
            })
        })

        // Test the Upstreamer implementation.
        Context("Upstreamer", func() {
            Context("Open()", func() {
                Context("when Memory interface is closed", func() {
                    It("changes state to open", func() {
                        err := mem.Open("", nil)
                        Expect(err).To(BeNil())

                        Expect(mem.State()).To(Equal(upstream.OPEN))
                    })

                    It("streams data to the provided writer", func() {
                        downstream := make(chan interface{})
                        defer close(downstream)

                        writer := &chanWriter{
                            downstream: downstream,
                        }

                        _ = mem.Open("", writer)

                        m := mem.(*upstream.Memory)
                        up := m.Upstream()

                        sent := []byte("testing")
                        up <- sent

                        received := <- downstream
                        Expect(received).To(Equal(sent))
                    })
                })
                Context("when Memory interface already open", func() {
                    BeforeEach(func() {
                        _ = mem.Open("", nil)
                    })

                    It("doesn't do anything", func() {
                        m := mem.(*upstream.Memory)
                        up := m.Upstream()

                        err := mem.Open("", nil)
                        Expect(err).To(BeNil())
                        Expect(mem.State()).To(Equal(upstream.OPEN))

                        mAfter := mem.(*upstream.Memory)
                        upAfter := mAfter.Upstream()

                        Expect(up).To(Equal(upAfter))
                    })
                })
            })

            Context("Close()", func() {

                Context("when Memory interface is open", func() {
                    It("changes state to closed", func() {
                        _ = mem.Open("", nil)

                        err := mem.Close()
                        Expect(err).To(BeNil())

                        Expect(mem.State()).To(Equal(upstream.CLOSED))
                    })

                    It("closes the underlying upstream channel", func() {
                        _ = mem.Open("", nil)

                        _ = mem.Close()

                        m := mem.(*upstream.Memory)
                        up := m.Upstream()
                        Expect(up).To(BeNil())
                    })

                })

                Context("when Memory interface already closed", func() {
                    It("doesn't do anything", func() {
                        err := mem.Close()
                        Expect(err).To(BeNil())

                        m := mem.(*upstream.Memory)
                        Expect(m.Upstream()).To(BeNil())
                    })
                })
            })
        })

        // TODO: Benchmark tests.
    })
})
