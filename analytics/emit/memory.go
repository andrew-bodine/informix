package emit

import (
    "io/ioutil"
    "strconv"
    "strings"
)

const (
    MEMORY = "memory"
)

// An Emitter implementation for memory information in linux.
type memory struct {

    // Keep a reference to the most recent memory stats, this is
    // useful for following the temporal compression best practice.
    last    chan map[string]int
}

// Don't expose the Memory struct because of the channel wrapping the
// last memory stats. Due to the write before ready nature of a channel,
// we need to put something in at creation.
func Memory() *memory {
    m := &memory{
        last:   make(chan map[string]int, 1),
    }

    m.last <- nil

    return m
}

// Implement the Generator interface.
func (m *memory) Generate() interface{} {
    next, err := m.ParseProcMeminfo()
    if err != nil {
        return nil
    }

    prev := <- m.last
    m.last <- next

    // Temporal compression.
    if next["MemAvailable"] == prev["MemAvailable"] {
        return nil
    }

    return next
}

// ParseProcMeminfo reads the current Linux memory stats from /proc/meminfo
// and returns a map representation of the contents.
func (m *memory) ParseProcMeminfo() (map[string]int, error) {
    bs, err := ioutil.ReadFile("/proc/meminfo")
    if err != nil {
        return nil, err
    }

    meminfo := make(map[string]int)

    for _, line := range strings.Split(string(bs), "\n") {
        parts := strings.Split(line, ":")
        if len(parts) != 2 {
            continue
        }

        parts[1] = strings.Trim(parts[1], " ")
        parts[1] = strings.Split(parts[1], " ")[0]

        value, _ := strconv.Atoi(string(parts[1]))
        meminfo[parts[0]] = value
    }

    return meminfo, nil
}
