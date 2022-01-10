package raft

import "time"

func (cm *CM) becomeFollower(term uint32) {
	cm.Lock()
	cm.state = "follower"
	cm.currentTerm = term
	cm.votedFor = ""
	cm.lastReset = time.Now()
	cm.Unlock()

	cm.startElectionTimer()
}

func (cm *CM) becomeCandidate() {
	cm.Lock()
	cm.state = "candidate"
	cm.currentTerm += 1
	cm.Unlock()

	cm.startElection()
}

func (cm *CM) becomeLeader() {}