package policy_test

import (
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"

    linuxproc "github.com/c9s/goprocinfo/linux"

    "github.com/andrew-bodine/informix/analytics/emit"
    . "github.com/andrew-bodine/informix/analytics/policy"
    "github.com/andrew-bodine/informix/analytics/queue"
)

func NewBadCPUStat(cs linuxproc.CPUStat) linuxproc.CPUStat {
    total := cs.User + cs.Nice + cs.System + cs.Idle + cs.IOWait
    total += cs.IRQ + cs.SoftIRQ + cs.Steal + cs.Guest + cs.GuestNice

    return linuxproc.CPUStat{
        Id:         cs.Id,
        User:       total,
        Nice:       0,
        System:     0,
        Idle:       0,
        IOWait:     0,
        IRQ:        0,
        SoftIRQ:    0,
        Steal:      0,
        Guest:      0,
        GuestNice:  0,
    }
}

var _ = Describe("policy", func() {
    var ph queue.PushHandler
    var procgen emit.Generator
    var data interface{}

    Context("Processor", func() {
        BeforeEach(func() {
            ph = &Processor{}

            procgen = emit.Processor()
        })

        Context("MoreThanUser()", func() {
            BeforeEach(func() {
                data = procgen.Generate()
            })

            Context("with a safe level of cpu usage", func() {
                It("return false", func() {
                    i := data.(*linuxproc.Stat)

                    proc := ph.(*Processor)
                    res := proc.MoreThanUser([]*linuxproc.Stat{i}, MaxPercentUser)
                    Expect(res).To(Equal(false))
                })
            })

            Context("with an unsafe level of cpu usage", func() {
                BeforeEach(func() {
                    i := data.(*linuxproc.Stat)

                    i.CPUStatAll = NewBadCPUStat(i.CPUStatAll)

                    data = i
                })

                It("returns true", func() {
                    i := data.(*linuxproc.Stat)

                    proc := ph.(*Processor)
                    res := proc.MoreThanUser([]*linuxproc.Stat{i}, MaxPercentUser)
                    Expect(res).To(Equal(true))
                })
            })
        })

        // Test the queue.PushHandler implementation.
        Context("queue.PushHandler", func() {
            Context("with a nil client", func() {
                It("doesn't do anything", func() {
                    ph.AfterPush([]interface{}{})
                })
            })

            Context("with a valid client", func() {
                BeforeEach(func() {
                    proc := ph.(*Processor)
                    proc.Downstream = &MockDownstreamer{
                        Payloads: make(map[string]map[string]interface{}),
                    }

                    ph = proc
                })

                It("forwards the data downstream", func() {
                    data = procgen.Generate()
                    info := data.(*linuxproc.Stat)
                    info.CPUStatAll = NewBadCPUStat(info.CPUStatAll)
                    data = info

                    ph.AfterPush([]interface{}{data})
                    proc := ph.(*Processor)
                    md := proc.Downstream.(*MockDownstreamer)
                    Expect(len(md.Payloads)).To(Equal(1))
                })
            })
        })
    })
})
