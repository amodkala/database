package raft

import (
	"log"
	"time"

    "github.com/amodkala/database/pkg/common"
)

func (cm *CM) becomeFollower(term uint32) {
	cm.mu.Lock()
	cm.state = "follower"
	cm.currentTerm = term
	cm.votedFor = ""
	cm.mu.Unlock()

	log.Printf("term %d -> %s became follower\n", cm.currentTerm, cm.self)
	go cm.startElectionTimer()
}

func (cm *CM) becomeCandidate() {

	cm.mu.Lock()
    cm.currentTerm += 1
    if cm.state != "candidate" {
        cm.state = "candidate"
    }
    log.Printf("term %d -> %s became candidate\n", cm.currentTerm, cm.self)
	cm.mu.Unlock()

	cm.startElection()
}

func (cm *CM) becomeLeader() {
	cm.mu.Lock()
	cm.state = "leader"
	cm.leader = cm.self
	log.Printf("***** term %d -> %s became leader *****\n", cm.currentTerm, cm.self)
    cm.log.Write(&common.Entry{
        RaftTerm: cm.currentTerm,
    })
    for i := range cm.nextIndex {
        cm.nextIndex[i] = cm.lastApplied + 1
    }
	cm.mu.Unlock()

    ticker := time.NewTicker(time.Duration(50) * time.Millisecond)
    defer ticker.Stop()

    for {
        <-ticker.C

        cm.mu.Lock()
        if cm.state != "leader" {
            cm.mu.Unlock()
            return
        }
        cm.mu.Unlock()

        cm.sendHeartbeats()
    }
}
