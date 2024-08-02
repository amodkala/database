package main

import (
    "fmt"
    "log"
    "time"

    ts "google.golang.org/protobuf/types/known/timestamppb"

    "github.com/amodkala/raft/pkg/common"
    "github.com/amodkala/raft/pkg/raft"
    tx "github.com/amodkala/raft/pkg/transaction"
)

const (
    name = "example"
    address = "localhost:8080"
)

func main() {

    entries := []*common.Entry{}

    now := ts.Now()
    for i := range 10 {
        entries = append(entries, &common.Entry{
            Key: uint32(i),
            CreateTime: now,
            Value: fmt.Sprintf("entry %d", i),
        })
    }

    cm := raft.New(name, address)

    txId := "test tx"
    transaction := tx.New(txId, entries...)

    go func() {
        if err := cm.Start(); err != nil {
            log.Fatal(err)
        }
    }()    

    time.Sleep(time.Second)

    if leaderId, err := cm.Replicate(transaction); err != nil {
        switch leaderId {
        case "":
            log.Fatal(err)
        default:
            log.Printf("leader is %s", leaderId)
        }
    }
}
