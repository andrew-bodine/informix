package emit

import (
    linuxproc "github.com/c9s/goprocinfo/linux"
)

const (
    PROCESSOR = "processor"
)

func Processor() *processor {
    p := &processor{
        last:   make(chan *linuxproc.Stat, 1),
    }

    p.last <- nil

    return p
}

// An Emitter implementation for processor information in linux.
type processor struct {

    // Keep a reference to the most recent memory stats, this is
    // useful for following the temporal compression best practice.
    last    chan *linuxproc.Stat
}

// Implement the Generator interface.
func (p *processor) Generate() interface{} {
    next, err := linuxproc.ReadStat("/proc/stat")
    if err != nil {
    	return nil
    }

    _ = <- p.last
    p.last <- next

    return next
}
