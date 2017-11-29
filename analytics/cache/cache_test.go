package cache_test

import (
	"time"

	cache "."

	"github.com/andrew-bodine/informer/analytics/queue"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("cache", func() {
	var c cache.Cacher
	var dir *cache.Directive

	Describe("Cacher", func() {
		Context("Register", func() {
			Context("without a unique key", func() {
				BeforeEach(func() {
					c = cache.NewCache()

					dir = &cache.Directive{
						Key:    "cacher:register:tests",
						Queuer: queue.NewQueue(1),
						Source: make(chan interface{}),
					}

					c.Register(dir)
				})

				It("should do nothing", func() {
					c.Register(dir)

					Expect(len(c.Keys())).To(Equal(1))
					dir.Stop()
				})
			})

			Context("with a unique key", func() {
				BeforeEach(func() {
					c = cache.NewCache()

					dir = &cache.Directive{
						Key:    "cacher:register:tests",
						Queuer: queue.NewQueue(2),
						Source: make(chan interface{}),
					}

					c.Register(dir)
				})

				Context("when directive is stopped", func() {
					It("stops pumping data and exits", func() {
						dir.Stop()
					})
				})

				Context("when there isn't any data avaiable from source", func() {
					It("doesn't forward anything to the queuer", func() {
						Expect(dir.Queuer.Count()).To(Equal(0))
						dir.Stop()
					})
				})

				Context("when there is data available from source channel", func() {
					BeforeEach(func() {
						dir.Source <- 0
					})

					It("forwards the data to the queuer", func() {
						Expect(dir.Queuer.Count()).To(Equal(1))
						dir.Stop()
					})
				})
			})
		})

		Context("Unregister", func() {
			Context("with a registered directive", func() {
				BeforeEach(func() {
					c = cache.NewCache()

					dir = &cache.Directive{
						Key:    "cacher:register:tests",
						Queuer: queue.NewQueue(1),
						Source: make(chan interface{}),
					}

					c.Register(dir)
				})

				It("closes the directive, and cleans up", func() {
					Expect(dir.Queuer.Count()).To(Equal(0))
					c.Unregister(dir)

					timeout := make(chan bool)

					// Try to write to the source channel, this should block
					// if the directive was actually unregistered as there is
					// no receiver.
					go func() {
						select {
						case dir.Source <- 0:
							timeout <- false
						case <-time.After(time.Millisecond):
							timeout <- true
						}
					}()

					value := <-timeout
					Expect(value).To(Equal(true))
				})
			})
			Context("with an unknown directive", func() {})
			Context("with a registered directive closer, closed", func() {})
		})
	})
})
