package wiot_test

import (
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"

    mqtt "github.com/eclipse/paho.mqtt.golang"

    . "github.com/andrew-bodine/informix/wiot"
)

func MQTTServerIsUp(broker string) bool {
    opts := mqtt.NewClientOptions().AddBroker(broker)
    client := mqtt.NewClient(opts)

    if token := client.Connect(); token.Wait() && token.Error() != nil {
            return false
    }

    return true
}

const (
    TestBroker= "tcp://localhost:1883"
)

var _ = Describe("wiot", func() {
    Context("NewClient()", func() {
        Context("with invalid options", func() {
            It("return nil", func() {
                c := NewClient(nil)
                Expect(c).To(BeNil())
            })
        })

        Context("with valid options", func() {
            It("returns a client", func() {
                o := NewOptions("test", "test", "test", "test")

                c := NewClient(o)
                Expect(c).NotTo(BeNil())
            })
        })
    })

    Context("client", func() {
        var opts *Options

        BeforeEach(func() {
            opts = NewOptions("test", "test", "test", "test")
        })

        Context("Connect()", func() {
            Context("when there is an error", func() {
                It("returns the error", func() {
                    c := NewClient(opts)

                    // This should error because there isn't an MQTT server
                    // listening at the configured broker address.
                    err := c.Connect()
                    Expect(err).NotTo(BeNil())
                })
            })

            Context("with our live test broker", func() {
                It("connects without error", func() {
                    if !MQTTServerIsUp(TestBroker) {
                        Skip("It seems the test broker is down, skipping.")
                        return
                    }
                    opts.Broker = TestBroker

                    c := NewClient(opts)
                    err := c.Connect()
                    Expect(err).To(BeNil())
                })
            })
        })

        Context("Publish()", func() {
            var payload map[string]interface{}

            BeforeEach(func() {
                payload = map[string]interface{}{
                    "foo":  "bar",
                }
            })

            Context("before it is connected", func() {
                It("returns an error", func() {
                    c := NewClient(opts)

                    err := c.Publish("test", payload)
                    Expect(err).NotTo(BeNil())
                })
            })

            Context("after it is connected", func() {
                Context("with a serialization error", func() {
                    It("returns the error", func() {
                        if !MQTTServerIsUp(TestBroker) {
                            Skip("It seems the test broker is down, skipping.")
                            return
                        }
                        opts.Broker = TestBroker

                        c := NewClient(opts)
                        c.Connect()
                        err := c.Publish("test",
                            map[string]interface{}{
                                "": "foobar",
                            },
                        )
                        Expect(err).NotTo(BeNil())
                    })
                })

                It("publishes messages to the broker without error", func() {
                    if !MQTTServerIsUp(TestBroker) {
                        Skip("It seems the test broker is down, skipping.")
                        return
                    }
                    opts.Broker = TestBroker

                    c := NewClient(opts)
                    c.Connect()

                    topic := "iot-2/type/test/id/test/evt/test/fmt/json"

                    // Get a reference to the underlying MQTT client so we
                    // can listen for the expected message.
                    mqttCli := c.MQTTClient()

                    // Hookup handler for test message.
                    done := make(chan mqtt.Message, 1)

                    handler := func(client mqtt.Client, msg mqtt.Message) {
                        done <- msg
                    }

                    token := mqttCli.Subscribe(topic, 0, handler)
                    token.Wait()
                    Expect(token.Error()).To(BeNil())

                    err := c.Publish("test", payload)
                    Expect(err).To(BeNil())

                    msg := <- done
                    Expect(msg).NotTo(BeNil())
                })
            })
        })
    })
})
