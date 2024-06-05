package raft

import (
    "sync"
    "time"
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
    commitChan chan Entry

    // implementation requirements
    currentTerm uint32
    votedFor    string
    log         []Entry
    commitIndex uint32
    lastApplied uint32
    nextIndex   []uint32
    matchIndex  []uint32
}

// New initializes a CM with some of its default values, the rest are
// initialized when Start is called
func New(address string, opts ...CMOpts) (*CM, chan Entry) {

    commitChan := make(chan Entry)

    cm := &CM{
        mu:          sync.Mutex{},
        self:        address,
        commitChan:  commitChan,
        log:         []Entry{},
        commitIndex: 0,
        lastApplied: 0,
        nextIndex:   []uint32{},
        matchIndex:  []uint32{},
    }

    return cm, commitChan
}
