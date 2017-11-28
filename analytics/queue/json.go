package queue

import (
    "bytes"
    "fmt"
)

// Implement json.Marshaler interface.
func (q *queue) MarshalJSON() ([]byte, error) {
    b := bytes.NewBufferString(`"[`)

    ptr := q.head

    if ptr != nil {
        b.WriteString(fmt.Sprintf("%v", ptr.data))
        ptr = ptr.next
    }

    for ptr != nil {
        b.WriteString(`,`)
        b.WriteString(fmt.Sprintf("%v", ptr.data))
        ptr = ptr.next
    }

    b.WriteString(`]"`)

    return b.Bytes(), nil
}
