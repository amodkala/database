package raft

import (
	"log"
	"time"
)

func (cm *CM) becomeFollower(term int32) {
	cm.Lock()
	cm.state = "follower"
	cm.currentTerm = term
	cm.votedFor = ""
	cm.lastReset = time.Now()
	log.Printf("%s became follower in term %d\n", cm.self, cm.currentTerm)
	cm.Unlock()

	go cm.startElectionTimer()
}

func (cm *CM) becomeCandidate() {
	cm.Lock()
	cm.state = "candidate"
	cm.currentTerm += 1
	log.Printf("%s became candidate in term %d\n", cm.self, cm.currentTerm)
	cm.Unlock()

	cm.startElection()
}

func (cm *CM) becomeLeader() {
	cm.Lock()
	cm.state = "leader"
	cm.leader = cm.self
	log.Printf("***** %s became leader in term %d *****\n", cm.self, cm.currentTerm)
	cm.Unlock()

	go func() {
		ticker := time.NewTicker(time.Duration(50) * time.Millisecond)
		defer ticker.Stop()

		for {
			cm.sendHeartbeats()
			<-ticker.C

			cm.Lock()
			if cm.state != "leader" {
				cm.Unlock()
				return
			}
			cm.Unlock()
		}
	}()
}
