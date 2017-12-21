package wiot_test

import (
    "os"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"

    . "github.com/andrew-bodine/informix/downstream/wiot"
)

var _ = Describe("wiot", func() {
    Context("NewOptions()", func() {
        Context("with at least one blank string param", func() {
            It("returns nil", func() {
                Expect(NewOptions("", "a", "b", "c")).To(BeNil())
                Expect(NewOptions("a", "", "b", "c")).To(BeNil())
                Expect(NewOptions("a", "b", "", "c")).To(BeNil())
                Expect(NewOptions("a", "b", "c", "")).To(BeNil())
            })
        })

        Context("with all valid params", func() {
            var o *Options

            BeforeEach(func() {
                o = NewOptions("test", "test", "test", "test")
            })

            It("sets broker correctly", func() {
                b := "tcp://test.messaging.internetofthings.ibmcloud.com:1883"
                Expect(o.Broker).To(Equal(b))
            })

            It("sets the client id correctly", func() {
                Expect(o.ClientId).To(Equal("g:test:test:test"))
            })

            It("sets username correctly", func() {
                Expect(o.Username).To(Equal("use-token-auth"))
            })

            It("sets password correctly", func() {
                Expect(o.Password).To(Equal("test"))
            })
        })
    })

    Context("NewOptionsFromEnv()", func() {
        var o *Options
        var backup map[string]string

        BeforeEach(func() {
            backup = make(map[string]string)

            if o := os.Getenv(Org); o != "" {
                backup[Org] = o
            }
            os.Setenv(Org, "envtest")

            if dt := os.Getenv(DeviceType); dt != "" {
                backup[DeviceType] = dt
            }
            os.Setenv(DeviceType, "envtest")

            if di := os.Getenv(DeviceId); di != "" {
                backup[DeviceId] = di
            }
            os.Setenv(DeviceId, "envtest")

            if t := os.Getenv(Token); t != "" {
                backup[Token] = t
            }
            os.Setenv(Token, "envtest")
        })

        AfterEach(func() {
            for k, v := range backup {
                os.Setenv(k, v)
            }
        })

        It("returns corresponding options", func() {
            o = NewOptionsFromEnv()
            Expect(o).NotTo(BeNil())
            b := "tcp://envtest.messaging.internetofthings.ibmcloud.com:1883"
            Expect(o.Broker).To(Equal(b))
            Expect(o.ClientId).To(Equal("g:envtest:envtest:envtest"))
            Expect(o.Username).To(Equal("use-token-auth"))
            Expect(o.Password).To(Equal("envtest"))
        })
    })

    Context("Options", func() {
        var o *Options

        BeforeEach(func() {
            o = NewOptions("org", "type", "id", "token")
        })

        Context("DeviceType()", func() {
            It("returns the correct device type", func() {
                Expect(o.DeviceType()).To(Equal("type"))
            })
        })

        Context("DeviceId()", func() {
            It("returns the correct device id", func() {
                Expect(o.DeviceId()).To(Equal("id"))
            })
        })
    })
})
