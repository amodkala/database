package raft

import (
    "sync"
    "time"

    "github.com/amodkala/raft/proto"
)

type Entry struct {
    Key []byte
    Value []byte
}

//
// CM (Consensus Module) is a struct that implements the Raft consensus
// algorithm (https://raft.github.io/raft.pdf) up to and including Section 5
// of the original paper.
//
type CM struct {
    proto.UnimplementedRaftServer
    mu sync.Mutex

    self       string
    state      string
    leader     string
    peers      []proto.RaftClient
    lastReset  time.Time
    CommitChan chan proto.Entry

    currentTerm int32
    votedFor    string
    log         []*proto.Entry
    commitIndex int32
    lastApplied int32
    nextIndex   []int32
    matchIndex  []int32
}

//
// New initializes a CM with some of its default values, the rest are
// initialized when Start is called
//
func New(commitChan chan proto.Entry, opts ...CMOpts) *CM {
    cm := &CM{
        mu:          sync.Mutex{},
        CommitChan:  commitChan,
        log:         []*proto.Entry{},
        commitIndex: 0,
        lastApplied: 0,
        nextIndex:   []int32{},
        matchIndex:  []int32{},
    }

    for _, opt := range opts {
        opt(cm)
    }

    return cm
}
