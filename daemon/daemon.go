package daemon

import (
    "fmt"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/andrew-bodine/informix/analytics"
    "github.com/andrew-bodine/informix/downstream"
    "github.com/andrew-bodine/informix/downstream/wiot"
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

    // Downstream target for data or events generated from informix.
    var dStreamer downstream.Downstreamer = nil

    // Create a wiot client to pass into builtin analytics.
    if opts := wiot.NewOptionsFromEnv(); opts != nil {
        ds := wiot.NewClient(opts)

        fmt.Println("Informix downstreaming to", opts.Broker)
        if err := ds.Connect(); err == nil {
            dStreamer = ds
        }
    }

    // Start builtin analytics and monitoring routine.
    builtin := analytics.NewBuiltin(dStreamer)
    builtin.Run(time.Second * 3)

    http.HandleFunc("/analytics/builtin", builtin.CacheHandler)
    http.ListenAndServe(":80", nil)

    // TODO: Serve registration protocol on /var/run/informix.sock

    // Wait for a signal.
    <- signals

    // Stop builtin analytics and monitoring routine.
    builtin.Stop()

    fmt.Println("Informix is shutdown.")
}
