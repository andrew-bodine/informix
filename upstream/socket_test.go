package upstream_test

import (
    "net"
    "os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

    "github.com/andrew-bodine/informix/upstream"
)

const (
    SOCK = "/tmp/informix.sock"
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
                _ = os.Remove(SOCK)
            })

            Context("Open()", func() {
                Context("when Socket is already closed", func() {
                    It("changes state to open", func() {
                        err := soc.Open(SOCK, nil)
                        Expect(err).To(BeNil())

                        Expect(soc.State()).To(Equal(upstream.OPEN))
                    })

                    It("opens the underlying interface", func() {
                        _ = soc.Open(SOCK, nil)

                        // If the socket was actually opened, we expect
                        // the corresponding file to exist.
                        _, err := os.Stat(SOCK)
                        Expect(err).To(BeNil())
                    })

                    It("streams data to the provided channel", func() {
                        downstream := make(chan interface{})
                        defer close(downstream)

                        writer := &chanWriter{
                            downstream: downstream,
                        }

                        _ = soc.Open(SOCK, writer)

                        sent := []byte("testing")
                        con, err := net.Dial("unix", SOCK)
                        Expect(err).To(BeNil())
                        defer con.Close()
                        _, err = con.Write(sent)
                        Expect(err).To(BeNil())

                        received := <- downstream
                        Expect(received).To(Equal(sent))
                    })

                    Context("an error occurs", func() {
                        BeforeEach(func() {
                            f, _ := os.Create(SOCK)
                            defer f.Close()
                        })

                        It("returns the error", func() {
                            err := soc.Open(SOCK, nil)
                            Expect(err).ToNot(BeNil())
                        })

                        It("doesn't change state", func() {
                            _ = soc.Open(SOCK, nil)
                            Expect(soc.State()).To(Equal(upstream.CLOSED))
                        })
                    })
                })

                Context("when Socket is already open", func() {
                    BeforeEach(func() {
                        _ = soc.Open(SOCK, nil)
                    })

                    It("doesn't do anything", func() {
                        before, err := os.Stat(SOCK)
                        Expect(err).To(BeNil())

                        err = soc.Open(SOCK, nil)
                        Expect(err).To(BeNil())
                        Expect(soc.State()).To(Equal(upstream.OPEN))

                        after, err := os.Stat(SOCK)
                        Expect(err).To(BeNil())
                        Expect(before).To(Equal(after))
                    })
                })
            })

            Context("Close()", func() {
                Context("when Socket is already open", func() {
                    BeforeEach(func() {
                        _ = soc.Open(SOCK, nil)
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
                        _, err := os.Stat(SOCK)
                        Expect(err).ToNot(BeNil())
                    })

                    It("stops streaming data downstream", func() {
                        _ = soc.Close()
                        soc = upstream.NewSocket()

                        downstream := make(chan interface{})
                        defer close(downstream)

                        writer := &chanWriter{
                            downstream: downstream,
                        }

                        _ = soc.Open(SOCK, writer)

                        con, _ := net.Dial("unix", SOCK)
                        defer con.Close()

                        // Should close the underlying interface, so we expect
                        // an error on the next write.
                        _ = soc.Close()

                        _, _ = con.Write([]byte("testing"))
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
                        f, _ := os.Create(SOCK)
                        defer f.Close()

                        err := soc.Close()
                        Expect(err).To(BeNil())

                        _, err = os.Stat(SOCK)
                        Expect(err).To(BeNil())
                    })
                })
            })
        })

        Context("Benchmark", func() {
            Context("more than one socket connection", func() {
                var num = 1
                var conns []net.Conn
                var downstream chan interface{}

                BeforeEach(func() {
                    soc = upstream.NewSocket()

                    downstream = make(chan interface{})

                    writer := &chanWriter{
                        downstream:   downstream,
                    }

                    _ = soc.Open(SOCK, writer)

                    for i := 1; i <= num; i++ {
                        c, err := net.Dial("unix", SOCK)
                        Expect(err).To(BeNil())

                        conns = append(conns, c)
                    }
                })

                AfterEach(func() {
                    _ = soc.Close()

                    close(downstream)

                    for _, c := range conns {
                        c.Close()
                    }
                })

                Context("all open at the same time, then close at the same time", func() {
                    It("processes all the data", func() {
                        done := make(chan bool)
                        defer close(done)

                        // Pump the channel, in runtime there will be another
                        // goroutine doing this. We want to avoid block sends
                        // and receives here for the test and the stream()
                        // routine.
                        go func() {
                            for i := 1; i <= num; i++ {
                                <- downstream
                            }

                            done <- true
                        }()

                        sent := []byte("testing")
                        for _, c := range conns {
                            _, _ = c.Write(sent)
                        }

                        <- done
                    })
                })
            })
        })

    })
})
