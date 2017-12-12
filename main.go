package main

import (
    "os"

    "github.com/andrew-bodine/informix/daemon"
)

func main() {
    daemon.Daemon(os.Args)
}
