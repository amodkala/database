package main

import (
    "flag"
    "fmt"

    "github.com/amodkala/db/pkg/client"
)

var (
    engine = flag.String("engine", "e", "custom", "which engine to use for WAL and storage (default \"custom\")")
    frequency = flag.Int("frequency", "f", 500, "how often Raft cluster leaders send heartbeats to followers in ms (default 500ms)")
)

func main() {
    flag.Parse()

    fmt.Println("nothing here yet :)")
}
