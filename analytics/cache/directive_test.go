package cache_test

import (
	cache "."

	"github.com/andrew-bodine/informer/analytics/queue"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("cache", func() {
	Context("Directive", func() {
		var dir *cache.Directive
		var src chan interface{}

		Context("Run", func() {
			Context("when it is running", func() {
				BeforeEach(func() {
					src = make(chan interface{})

					dir = &cache.Directive{
						Source: src,
						Queuer: queue.NewQueue(3),
					}

					dir.Run()
					src <- 0
				})

				It("doesn't do anything", func() {
					Expect(dir.Queuer.Count()).To(Equal(1))
					dir.Run()
					Expect(dir.Queuer.Count()).To(Equal(1))
				})
			})

			Context("when it isn't running", func() {
				BeforeEach(func() {
					src = make(chan interface{})

					dir = &cache.Directive{
						Source: src,
						Queuer: queue.NewQueue(5),
					}
				})

				It("starts the directive routine", func() {
					dir.Run()
					src <- 0
					src <- 1
					src <- 2
					Expect(dir.Queuer.Count()).To(Equal(3))
				})
			})
		})

		Context("Stop", func() {
			Context("when it isn't running", func() {
				BeforeEach(func() {
					src = make(chan interface{})

					dir = &cache.Directive{
						Source: src,
						Queuer: queue.NewQueue(5),
					}
				})

				It("doesn't do anything", func() {
					dir.Stop()
				})
			})

			Context("when it is running", func() {
				BeforeEach(func() {
					src = make(chan interface{})

					dir = &cache.Directive{
						Source: src,
						Queuer: queue.NewQueue(5),
					}

					dir.Run()
					src <- 0
					src <- 1
					src <- 2
				})

				It("stops the directive routine", func() {
					Expect(dir.Queuer.Count()).To(Equal(3))
					dir.Stop()
					Expect(dir.Queuer.Count()).To(Equal(3))
				})
			})
		})
	})
})
