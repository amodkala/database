package raft

import "time"

func (cm *CM) becomeFollower(term int32) {
	cm.Lock()
	cm.state = "follower"
	cm.currentTerm = term
	cm.votedFor = ""
	cm.lastReset = time.Now()
	cm.Unlock()

	go cm.startElectionTimer()
}

func (cm *CM) becomeCandidate() {
	cm.Lock()
	cm.state = "candidate"
	cm.currentTerm += 1
	cm.Unlock()

	cm.startElection()
}

func (cm *CM) becomeLeader() {
	cm.Lock()
	cm.state = "leader"
	cm.leader = cm.self
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
