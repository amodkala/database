package raft

import (
	"log"
	"time"
)

func (cm *CM) becomeFollower(term int32) {
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
    cm.log = append(cm.log, Entry{
        Term: cm.currentTerm,
        Message: []byte{},
    })
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
