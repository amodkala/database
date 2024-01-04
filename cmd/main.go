package main

import (
    "os"
    "github.com/amodkala/raft"
)

func main() {
    addr := os.Args[1]
    peers := os.Args[2:]
    cm, _ := raft.New(addr)
    cm.Start(raft.WithLocalPeers(peers))
}
