package raft

import (
    "sync"
    "time"

    "github.com/amodkala/raft/proto"
)

//
// Entry is the struct that is returned to the client indicating that the
// key/value pair has been replicated across a majority of peers and can be
// stored
//
type Entry struct {
    Key   string
    Value string
}

//
// CM (Consensus Module) is a struct that implements the Raft consensus
// algorithm (https://raft.github.io/raft.pdf) up to and including Section 5
// of the original paper.
//
type CM struct {
    mu sync.Mutex

    self       string
    state      string
    Leader     string
    peers      []proto.RaftClient
    lastReset  time.Time
    CommitChan chan Entry

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
func New() *CM {
    return &CM{
        mu:          sync.Mutex{},
        CommitChan:  make(chan Entry),
        log:         []*proto.Entry{},
        commitIndex: 0,
        lastApplied: 0,
        nextIndex:   []int32,
        matchIndex:  []int32,
    }
}
