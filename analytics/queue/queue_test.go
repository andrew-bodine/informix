package queue_test

import (
    queue "."

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("queue", func() {
    var q queue.Queuer

    Context("NewQueue", func() {
        It("creates a Queue size 1 and sets inital count to 0", func() {
            q = queue.NewQueue(1)
            Expect(q.Size()).To(Equal(1))
            Expect(q.Count()).To(Equal(0))
        })
    })

    Context("Queue", func() {
        Context("Push", func() {
            Context("when Queue is empty", func() {
                BeforeEach(func() {
                    q = queue.NewQueue(1)
                })

                It("appends to the Queue", func() {
                    removed := q.Push(0)
                    Expect(removed).To(BeNil())
                    Expect(q.Size()).To(Equal(1))
                    Expect(q.Count()).To(Equal(1))
                })
            })

            Context("when Queue is not empty, and not full", func() {
                BeforeEach(func() {
                    q = queue.NewQueue(2)
                    _ = q.Push(0)
                })

                It("appends to the Queue", func() {
                    removed := q.Push(0)
                    Expect(removed).To(BeNil())
                    Expect(q.Size()).To(Equal(2))
                    Expect(q.Count()).To(Equal(2))
                })
            })

            Context("when Queue is full", func() {
                BeforeEach(func() {
                    q = queue.NewQueue(1)
                    _ = q.Push(0)
                })

                It("appends to the Queue, and returns oldest item", func() {
                    removed := q.Push(0)
                    Expect(removed).ToNot(BeNil())
                    Expect(q.Size()).To(Equal(1))
                    Expect(q.Count()).To(Equal(1))
                })
            })
        })

        Context("MarshalJSON", func() {
            Context("when Queue is empty", func() {
                BeforeEach(func() {
                    q = queue.NewQueue(1)
                })

                It("should return an accurate json string", func() {
                    serialized, err := q.MarshalJSON()
                    Expect(err).To(BeNil())
                    Expect(string(serialized)).To(Equal("\"[]\""))
                })
            })

            Context("when Queue is not empty, and not full", func() {
                BeforeEach(func() {
                    q = queue.NewQueue(2)
                    _ = q.Push(0)
                })

                It("should return an accurate json string", func() {
                    serialized, err := q.MarshalJSON()
                    Expect(err).To(BeNil())
                    Expect(string(serialized)).To(Equal("\"[0]\""))
                })
            })

            Context("when Queue is full", func() {
                BeforeEach(func() {
                    q = queue.NewQueue(1)
                    _ = q.Push(0)
                })

                It("should return an accurate json string", func() {
                    serialized, err := q.MarshalJSON()
                    Expect(err).To(BeNil())
                    Expect(string(serialized)).To(Equal("\"[0]\""))
                })
            })
        })
    })
})
