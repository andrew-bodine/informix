package emit_test

import (
    linuxproc "github.com/c9s/goprocinfo/linux"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"

    "github.com/andrew-bodine/informix/analytics/emit"
)

var _ = Describe("emit", func() {
    Context("Processor()", func() {
        It("creates a processor generator", func() {
            proc := emit.Processor()
            Expect(proc).NotTo(BeNil())
        })
    })

    Context("processor", func() {

        // Test the Generator implementation.
        Context("Generator", func() {
            Context("Generate()", func() {
                It("should return the current processor stats", func() {
                    proc := emit.Processor()
                    data := proc.Generate()
                    Expect(data).NotTo(BeNil())
                    info := data.(*linuxproc.Stat)
                    Expect(info.CPUStatAll.User).NotTo(Equal(0))
                })

                Context("when there is an error", func() {})
            })
        })
    })
})
