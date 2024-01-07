package main

import (
    "log"
    "os"
    "time"

    "github.com/amodkala/raft"
)

func main() {
    addr := os.Args[1]
    peers := os.Args[2:]
    cm, commitChan := raft.New(addr)

    go cm.Start(raft.WithLocalPeers(peers))

    time.Sleep(4 * time.Second)
    cm.Replicate([]byte("hello"), []byte("world"))

    for {
        entry := <-commitChan
        log.Printf("client of %s got entry %s", addr, string(entry))
    }
}
