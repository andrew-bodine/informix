package queue_test

import (
	queue "."

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("queue", func() {
	var q queue.Queuer

	Context("NewQueue", func() {
		It("creates a queue size 1 and sets inital count to 0", func() {
			q = queue.NewQueue(1)
			Expect(q.Size()).To(Equal(1))
			Expect(q.Count()).To(Equal(0))
		})
	})

	Context("Queue", func() {

		// Test the Queuer implementation.
		Context("Queuer", func() {
			Context("Copy()", func() {
				Context("when queue is empty", func() {
					BeforeEach(func() {
						q = queue.NewQueue(1)
					})

					It("returns and empty slice", func() {
						items := q.Copy()
						Expect(len(items)).To(Equal(0))
					})
				})

				Context("when queue isn't empty", func() {
					var expected []interface{} = []interface{}{0, 1, 2}

					BeforeEach(func() {
						q = queue.NewQueue(len(expected))

						for _, v := range expected {
							q.Push(v)
						}
					})

					It("returns an accurate copy", func() {
						actual := q.Copy()
						Expect(expected).To(Equal(actual))
					})
				})
			})

			Context("Push()", func() {
				Context("when queue is empty", func() {
					BeforeEach(func() {
						q = queue.NewQueue(1)
					})

					It("appends to the queue", func() {
						q.Push(0)
						Expect(q.Size()).To(Equal(1))
						Expect(q.Count()).To(Equal(1))
					})
				})

				Context("when queue is not empty, and not full", func() {
					BeforeEach(func() {
						q = queue.NewQueue(2)
						q.Push(0)
					})

					It("appends to the queue", func() {
						q.Push(0)
						Expect(q.Size()).To(Equal(2))
						Expect(q.Count()).To(Equal(2))
					})
				})

				Context("when queue is full", func() {
					BeforeEach(func() {
						q = queue.NewQueue(1)
						q.Push(0)
					})

					It("appends to the queue", func() {
						q.Push(0)
						Expect(q.Size()).To(Equal(1))
						Expect(q.Count()).To(Equal(1))
					})
				})
			})

			Context("Drain()", func() {
				Context("when queue is empty", func() {
					BeforeEach(func() {
						q = queue.NewQueue(1)
					})

					It("doesn't change anything", func() {
						q.Drain()
						Expect(q.Size()).To(Equal(1))
						Expect(q.Count()).To(Equal(0))
					})
				})

				Context("when queue is not empty", func() {
					BeforeEach(func() {
						q = queue.NewQueue(1)
						q.Push(0)
					})

					It("empties the queue", func() {
						q.Drain()
						Expect(q.Size()).To(Equal(1))
						Expect(q.Count()).To(Equal(0))
					})
				})
			})

			Context("OnPush()", func() {
				var handler queue.PushHandler
				var delegate chan interface{}

				BeforeEach(func() {
					q = queue.NewQueue(5)

					delegate = make(chan interface{})

					handler = &TestPushHandler{
						delegate:	delegate,
					}
				})

				AfterEach(func() {
					close(delegate)
				})

				It("sets the handler as a hook for pushed data", func() {
					q.OnPush(handler)
					q.Push(0)
					<- delegate
					Expect(q.Count()).To(Equal(1))
				})
			})
		})
	})
})
