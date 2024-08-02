package raft

import (
    "fmt"
    "sync"
    "time"

    "github.com/amodkala/raft/pkg/common"
    "github.com/amodkala/raft/pkg/wal"
)

// CM (Consensus Module) is a struct that implements the Raft consensus
// algorithm (https://raft.github.io/raft.pdf) up to and including Section 5
// of the original paper.
type CM struct {
    UnimplementedRaftServer
    mu sync.Mutex

    // metadata
    self       string
    state      string
    leader     string
    peers      []RaftClient
    lastReset  time.Time
    txChans map[string]chan *common.Entry
    commitChan chan *common.Entry
    errChan chan error

    // implementation requirements
    currentTerm uint32
    votedFor    string
    log         *wal.WAL
    commitIndex uint32
    lastApplied uint32
    nextIndex   []uint32
    matchIndex  []uint32
}

// New initializes a CM with some of its default values, the rest are
// initialized when Start is called
func New(id, serverAddress string, opts ...CMOpts) *CM {

    log := wal.New(fmt.Sprintf("/var/%s-log.wal", id))

    cm := &CM{
        mu:          sync.Mutex{},
        self:        serverAddress,
        txChans:  make(map[string]chan *common.Entry),
        commitChan: make(chan *common.Entry),
        errChan: make(chan error),
        log:         log,
        commitIndex: 0,
        lastApplied: 0,
        nextIndex:   []uint32{},
        matchIndex:  []uint32{},
    }

    return cm
}
