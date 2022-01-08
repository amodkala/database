package raft

//
// CM (Consensus Module) is a struct that implements the Raft consensus
// algorithm (https://raft.github.io/raft.pdf) up to and including Section 5
// of the original paper.
//
type CM struct {
	self string

	currentTerm uint32
	votedFor    string
	log         []string
	commitIndex uint32
	lastApplied uint32
	nextIndex   []uint32
	matchIndex  []uint32
}

func New(addr string) *CM {
	return &CM{
		self: addr,

		currentTerm: 0,
		votedFor:    "",
		log:         []string{},
		commitIndex: 0,
		lastApplied: 0,
		nextIndex:   []uint32{},
		matchIndex:  []uint32{},
	}
}
