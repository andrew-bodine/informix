package policy

import (
    linuxproc "github.com/c9s/goprocinfo/linux"

    "github.com/andrew-bodine/informix/downstream"
)

const (
    MinPercentAvailable = 25
    LowAvailableMemory = "Available memory is below 25%."
)

// A built-in Policy implementation that monitors system memory stats.
type Memory struct {
    Downstream  downstream.Downstreamer
}

// Implement the queue.PushHandler interface.
func (m *Memory) AfterPush(data []interface{}) {
    mems := make([]*linuxproc.MemInfo, len(data))

    for i, d := range data {
        mems[i] = d.(*linuxproc.MemInfo)
    }

    if m.Downstream == nil {
        return
    }

    // Memory policy pipeline.
    if m.LessThanAvailable(mems, MinPercentAvailable) {
        m.Downstream.Publish("memory", map[string]interface{}{
            "error": LowAvailableMemory,
        })
    }
}

// LessThanAvailable returns whether or not the percentage of availabe
// memory left on the system is less than the provided percent.
func (m *Memory) LessThanAvailable(data []*linuxproc.MemInfo, percent int) bool {
    available := data[0].MemAvailable * 100 / data[0].MemTotal

    if available < uint64(percent) {
        return true
    }

    return false
}
