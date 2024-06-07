package database

import (
    "fmt"

    "github.com/amodkala/database/pkg/raft"
    "github.com/amodkala/database/pkg/lsm"
    "github.com/amodkala/database/pkg/wal"
)

// Record is a convenience type to assert that users of the database Client
// conform to the expected types for keys and values (uint32 and string,
// respectively.) 
//
// Each Record is almost immediately turned into a common.Entry to form part of
// a Transaction
type Record struct {
    Key uint32
    Value string
}

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

func (c Client) Write(records ...Record) {
}

func (c Client) Read(key uint32) (value string, err error)
