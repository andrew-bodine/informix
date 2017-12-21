package daemon_test

import (
    "bufio"
    "fmt"
    "io"
    "os"
    "os/exec"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"

    "github.com/andrew-bodine/informix/downstream/wiot"
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

            Context("with Watson IoT Platform environment variables", func() {
                BeforeEach(func() {
                    daemonCmd.Process.Kill()
                    daemonCmd.Wait()

                    daemonCmd = exec.Command("informix")

                    out, _ = daemonCmd.StdoutPipe()

                    testConf := map[string]string{
                        wiot.Org: "testorg",
                        wiot.DeviceType: "testtype",
                        wiot.DeviceId: "testid",
                        wiot.Token: "testtoken",
                    }

                    for k, v := range testConf {
                        daemonCmd.Env = append(daemonCmd.Env,
                            fmt.Sprintf("%s=%s", k, v),
                        )
                    }

                    daemonCmd.Start()
                })

                It("outputs that it is downstreaming", func() {
                    reader := bufio.NewReader(out)
                    reader.ReadString('\n')
                    second, err := reader.ReadString('\n')
                    Expect(err).To(BeNil())
                    Expect(second).To(ContainSubstring("downstreaming"))

                    daemonCmd.Process.Kill()
                })
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
