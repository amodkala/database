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
	address   string
	state     string
	leader    string
	lastReset time.Time

	currentTerm uint32
	votedFor    string
	log         []proto.Entry
	commitIndex uint32
	lastApplied uint32
	nextIndex   []uint32
	matchIndex  []uint32
}

func New(name, addr string) *CM {
	return &CM{
		Mutex:       sync.Mutex{},
		self:        name,
		address:     addr,
		log:         []proto.Entry{},
		commitIndex: 0,
		lastApplied: 0,
		nextIndex:   []uint32{},
		matchIndex:  []uint32{},
	}
}
