package main

import (
    "flag"
    "fmt"

    "github.com/amodkala/db/pkg/client"
)

var (
    engine string 
    frequency int 
)

func main() {
    flag.StringVar(&engine, "engine", "e", "rocksdb", "which engine to use for WAL and storage (default \"rocksdb\")")
    flag.IntVar(&frequency, "frequency", "f", 500, "how often Raft cluster leaders send heartbeats to followers in ms (default 500ms)")
    flag.Parse()

    fmt.Println("%s, %d", engine, frequency)
}
