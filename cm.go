package raft

import (
    "sync"
    "time"

    "github.com/amodkala/raft/proto"
)

type Entry struct {
    Term int32
    Message []byte
}

type Peer struct {
    clientConn proto.RaftClient
    nextIndex int32
    matchIndex int32
}

// CM (Consensus Module) is a struct that implements the Raft consensus
// algorithm (https://raft.github.io/raft.pdf) up to and including Section 5
// of the original paper.
type CM struct {
    proto.UnimplementedRaftServer
    mu sync.Mutex

    // metadata
    self       string
    state      string
    leader     string
    peers      []proto.RaftClient
    lastReset  time.Time
    commitChan chan []byte

    // implementation requirements
    currentTerm int32
    votedFor    string
    log         []Entry
    commitIndex int32
    lastApplied int32
    nextIndex   []int32
    matchIndex  []int32
}

// New initializes a CM with some of its default values, the rest are
// initialized when Start is called
func New(address string, opts ...CMOpts) (*CM, chan []byte) {

    commitChan := make(chan []byte)

    cm := &CM{
        mu:          sync.Mutex{},
        self:        address,
        commitChan:  commitChan,
        log:         []Entry{},
        commitIndex: 0,
        lastApplied: 0,
        nextIndex:   []int32{},
        matchIndex:  []int32{},
    }

    return cm, commitChan
}
