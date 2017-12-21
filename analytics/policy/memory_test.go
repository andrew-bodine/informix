package policy_test

import (
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"

    linuxproc "github.com/c9s/goprocinfo/linux"

    "github.com/andrew-bodine/informix/analytics/emit"
    . "github.com/andrew-bodine/informix/analytics/policy"
    "github.com/andrew-bodine/informix/analytics/queue"
)

var _ = Describe("policy", func() {
    var p queue.PushHandler
    var memgen emit.Generator
    var data interface{}

    Context("Memory", func() {
        BeforeEach(func() {
            p = &Memory{}

            memgen = emit.Memory()
        })

        Context("LessThanAvailable()", func() {
            Context("with a safe amount of memory available", func() {
                BeforeEach(func() {
                    data = memgen.Generate()
                })

                It("returns false", func() {
                    i := data.(*linuxproc.MemInfo)

                    mp := p.(*Memory)
                    res := mp.LessThanAvailable([]*linuxproc.MemInfo{i}, MinPercentAvailable)
                    Expect(res).To(Equal(false))
                })
            })

            Context("with an unsafe amount of memory available", func() {
                BeforeEach(func() {
                    data = memgen.Generate()

                    info := data.(*linuxproc.MemInfo)
                    info.MemAvailable = uint64(1024)
                    data = info
                })

                It("returns true", func() {
                    i := data.(*linuxproc.MemInfo)

                    mp := p.(*Memory)
                    res := mp.LessThanAvailable([]*linuxproc.MemInfo{i}, MinPercentAvailable)
                    Expect(res).To(Equal(true))
                })
            })
        })

        // Test the queue.PushHandler implementation.
        Context("queue.PushHandler", func() {
            Context("AfterPush()", func() {
                Context("with a nil client", func() {
                    It("doesn't do anything", func() {
                        p.AfterPush([]interface{}{})
                    })
                })

                Context("with a valid client", func() {
                    BeforeEach(func() {
                        m := p.(*Memory)
                        m.Downstream = &MockDownstreamer{
                            Payloads: make(map[string]map[string]interface{}),
                        }

                        p = m
                    })

                    It("forwards the data downstream", func() {
                        data = memgen.Generate()
                        info := data.(*linuxproc.MemInfo)
                        info.MemAvailable = uint64(1024)
                        data = info

                        p.AfterPush([]interface{}{data})
                        m := p.(*Memory)
                        md := m.Downstream.(*MockDownstreamer)
                        Expect(len(md.Payloads)).To(Equal(1))
                    })
                })
            })
        })
    })
})
