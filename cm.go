package raft

import (
	"sync"
	"time"

	"github.com/amodkala/raft/proto"
)

//
// CM (Consensus Module) is a struct that implements the Raft consensus
// algorithm (https://raft.github.io/raft.pdf) up to and including Section 5
// of the original paper.
//
type CM struct {
	sync.Mutex

	self      string
	state     string
	leader    string
	peers     map[string]proto.RaftClient
	committed chan string
	lastReset time.Time

	currentTerm int32
	votedFor    string
	log         []*proto.Entry
	commitIndex int32
	lastApplied int32
	nextIndex   map[string]int32
	matchIndex  map[string]int32
}

//
// New initializes a CM with some of its default values, the rest are
// initialized when Start is called
//
func New(out chan string) *CM {
	// TODO: add peer discovery
	// TODO: add distinct id + address fields
	return &CM{
		Mutex:       sync.Mutex{},
		committed:   out,
		log:         []*proto.Entry{},
		commitIndex: 0,
		lastApplied: 0,
		nextIndex:   make(map[string]int32),
		matchIndex:  make(map[string]int32),
	}
}
