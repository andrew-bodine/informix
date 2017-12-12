package daemon

import (
    "fmt"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/andrew-bodine/informix/analytics"
)

func Daemon(args []string) {
    fmt.Println("Informix is starting up.")

    signals := make(chan os.Signal, 1)

    // Ask the OS to notify us on the signals channel, for the following
    // well known os signals.
    signal.Notify(
        signals,
        syscall.SIGHUP,
        syscall.SIGINT,
        syscall.SIGTERM,
        syscall.SIGQUIT,
    )

    // Start builtin analytics and monitoring routine.
    builtin := analytics.NewBuiltin()
    builtin.Run(time.Second * 3)

    http.HandleFunc("/analytics/builtin", builtin.CacheHandler)
    http.ListenAndServe(":80", nil)

    // Wait for a signal.
    <- signals

    // Stop builtin analytics and monitoring routine.
    builtin.Stop()

    fmt.Println("Informix is shutdown.")
}
