package cache_test

import (
    "fmt"

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
                var ready chan bool

                BeforeEach(func() {
                    c = cache.NewCache()

                    dir = &cache.Directive{
                        Key:        "cacher:register:tests",
                        Queuer:     queue.NewQueue(1),
                        Source:     make(chan interface{}),
                        Closer:     make(chan bool),
                    }

                    ready = make(chan bool)

                    go func() {
                        c.Register(dir)
                        ready <- true
                    }()
                })

                It("should do nothing", func() {
                    <- ready
                    close(ready)

                    fmt.Println(c.Keys())
                    Expect(len(c.Keys())).To(Equal(1))

                    exited := make(chan bool)

                    go func() {
                        c.Register(dir)
                        exited <- true
                    }()

                    <- exited

                    fmt.Println(c.Keys())
                    Expect(len(c.Keys())).To(Equal(1))

                    close(exited)
                })
            })

            Context("with a unique key", func() {
                var exited chan bool

                BeforeEach(func() {
                    c = cache.NewCache()

                    dir = &cache.Directive{
                        Key:        "cacher:register:tests",
                        Queuer:     queue.NewQueue(2),
                        Source:     make(chan interface{}),
                        Closer:     make(chan bool),
                    }

                    exited = make(chan bool)

                    go func() {
                        c.Register(dir)

                        exited <- true
                    }()
                })

                Context("when stopper channel closes", func() {
                    It("stops pumping data and exits", func() {
                        close(dir.Closer)
                        value := <- exited
                        Expect(value).To(Equal(true))
                        close(exited)
                    })
                })

                Context("when there isn't any data avaiable from source", func() {
                    It("doesn't forward anything to the queuer", func() {
                        Expect(dir.Queuer.Count()).To(Equal(0))
                        close(dir.Closer)
                        <- exited
                        close(exited)
                    })
                })

                Context("when there is data available from source channel", func() {
                    BeforeEach(func() {
                        dir.Source <- 0
                    })

                    It("forwards the data to the queuer", func() {
                        Expect(dir.Queuer.Count()).To(Equal(1))
                        close(dir.Closer)
                        <- exited
                        close(exited)
                    })
                })
            })
        })
    })
})
