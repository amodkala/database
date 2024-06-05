package database

import (
    "fmt"
    "log"

    "github.com/amodkala/database/pkg/raft"
    "github.com/amodkala/database/pkg/lsm"
    "github.com/amodkala/database/pkg/wal"
)

// Client is the primary type that will be exposed to users of this
// package/module. 
type Client struct {
    cm *raft.CM
    wal *wal.WAL
    lsm lsm.LSMTree
}

func New(addr string, id uint32) Client {

    cm, _ := raft.New(addr)
    lsm := lsm.New(id)
    wal := wal.New(fmt.Sprintf("%d-log.wal", id))

    return Client{ cm, wal, lsm, }
}

func (c Client) Write(key uint32, value string) {
    if leaderAddr, err := c.cm.Replicate(key, value); leaderAddr != "" && err != nil {
        log.Println("cm is not leader, must redirect")
    }
}

func (c Client) Read(key uint32) (value string, err error)
