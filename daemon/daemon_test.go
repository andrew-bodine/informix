package daemon_test

import (
    "bufio"
    "io"
    "os"
    "os/exec"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("daemon", func() {
    Context("Daemon()", func() {
        Context("when it is run as a process", func() {
            var daemonCmd *exec.Cmd

            var out io.ReadCloser

            BeforeEach(func() {
                daemonCmd = exec.Command("informix")

                out, _ = daemonCmd.StdoutPipe()

                daemonCmd.Start()
            })

            AfterEach(func() {
                daemonCmd.Wait()
            })

            It("outputs that it has started", func() {
                reader := bufio.NewReader(out)
                first, err := reader.ReadString('\n')
                Expect(err).To(BeNil())
                Expect(first).To(ContainSubstring("starting"))

                daemonCmd.Process.Kill()
            })

            Context("when it receives an interrupt signal", func() {
                It("responds to the signal", func() {
                    err := daemonCmd.Process.Signal(os.Interrupt)
                    Expect(err).To(BeNil())

                    _, err = daemonCmd.Process.Wait()
                    Expect(err).To(BeNil())
                })
            })
        })
    })
})
