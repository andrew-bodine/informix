package policy

import (
    linuxproc "github.com/c9s/goprocinfo/linux"

    "github.com/andrew-bodine/informix/downstream"
)

const (
    MaxPercentUser = 75
    HighUserCPU = "User cpu percentage is above 75%."
)

// A built-in Policy implementation that monitors system cpu stats.
type Processor struct {
    Downstream  downstream.Downstreamer
}

// Implement the queue.PushHandler interface.
func (p *Processor) AfterPush(data []interface{}) {
    procs := make([]*linuxproc.Stat, len(data))

    for i, d := range data {
        procs[i] = d.(*linuxproc.Stat)
    }

    if p.Downstream == nil {
        return
    }

    // Processor policy pipeline.
    if p.MoreThanUser(procs, MaxPercentUser) {
        p.Downstream.Publish("processor", map[string]interface{}{
            "error": HighUserCPU,
        })
    }
}

// MoreThanUser returns whether or not the most recent average cpu usage
// indicates the systems processors have been spending more time on user
// processes than the provided threshold.
func (p *Processor) MoreThanUser(data []*linuxproc.Stat, percent int) bool {
    csa := data[0].CPUStatAll

    total := csa.User + csa.Nice + csa.System + csa.Idle + csa.IOWait
    total += csa.IRQ + csa.SoftIRQ + csa.Steal + csa.Guest + csa.GuestNice

    threshold := total * uint64(percent) / uint64(100)

    if csa.User >= threshold {
        return true
    }

    return false
}
