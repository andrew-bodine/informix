package upstream_test

import (
    "os"

	upstream "."

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("upstream", func() {
    var soc upstream.Upstreamer

    Context("NewSocket()", func() {
        It("creates a Upstreamer in the closed state", func() {
            soc = upstream.NewSocket()
            Expect(soc).ToNot(BeNil())

            _, ok := soc.(upstream.Upstreamer)
            Expect(ok).To(Equal(true))

            Expect(soc.State()).To(Equal(upstream.CLOSED))
        })
    })

    Context("Socket", func() {

        // Test the Upstreamer implementation.
        Context("Upstreamer", func() {
            BeforeEach(func() {
                soc = upstream.NewSocket()
            })

            AfterEach(func() {
                _ = os.Remove(upstream.SOCK)
            })

            Context("Open()", func() {
                Context("when Socket is already closed", func() {
                    It("changes state to open", func() {
                        err := soc.Open()
                        Expect(err).To(BeNil())

                        Expect(soc.State()).To(Equal(upstream.OPEN))
                    })

                    It("opens the underlying interface", func() {
                        _ = soc.Open()

                        // If the socket was actually opened, we expect
                        // the corresponding file to exist.
                        _, err := os.Stat(upstream.SOCK)
                        Expect(err).To(BeNil())
                    })

                    Context("an error occurs", func() {
                        BeforeEach(func() {
                            f, _ := os.Create(upstream.SOCK)
                            defer f.Close()
                        })

                        It("returns the error", func() {
                            err := soc.Open()
                            Expect(err).ToNot(BeNil())
                        })

                        It("doesn't change state", func() {
                            _ = soc.Open()
                            Expect(soc.State()).To(Equal(upstream.CLOSED))
                        })
                    })
                })

                Context("when Socket is already open", func() {
                    BeforeEach(func() {
                        _ = soc.Open()
                    })

                    It("doesn't do anything", func() {
                        before, err := os.Stat(upstream.SOCK)
                        Expect(err).To(BeNil())

                        err = soc.Open()
                        Expect(err).To(BeNil())
                        Expect(soc.State()).To(Equal(upstream.OPEN))

                        after, err := os.Stat(upstream.SOCK)
                        Expect(err).To(BeNil())
                        Expect(before).To(Equal(after))
                    })
                })
            })

            Context("Close()", func() {
                Context("when Socket is already open", func() {
                    BeforeEach(func() {
                        _ = soc.Open()
                    })

                    It("changes state to closed", func() {
                        err := soc.Close()
                        Expect(err).To(BeNil())

                        Expect(soc.State()).To(Equal(upstream.CLOSED))
                    })

                    It("closes the underlying interface", func() {
                        _ = soc.Close()

                        // If the socket was actually closed, we expect
                        // the corresponding file to not exist.
                        _, err := os.Stat(upstream.SOCK)
                        Expect(err).ToNot(BeNil())
                    })

                    // TODO: Re-enable this test, once you discover what kinds
                    // of errors can happen when closing a net.Listener so you
                    // can simulate that.
                    //
                    // Context("an error occurs", func() {
                    //     BeforeEach(func() {
                    //         // TODO: Simulate error.
                    //     })
                    //
                    //     It("returns the error", func() {
                    //         err = soc.Close()
                    //         Expect(err).ToNot(BeNil())
                    //         Expect(soc.State()).To(Equal(upstream.OPEN))
                    //     })
                    // })
                })

                Context("when Socket is already closed", func() {
                    It("doesn't do anything", func() {
                        f, _ := os.Create(upstream.SOCK)
                        defer f.Close()

                        err := soc.Close()
                        Expect(err).To(BeNil())

                        _, err = os.Stat(upstream.SOCK)
                        Expect(err).To(BeNil())
                    })
                })
            })
        })
    })
})
