package raft

import (
	"math/rand"
	"time"
)

func (cm *CM) startElectionTimer() {
	cm.Lock()
	termStarted := cm.currentTerm
	cm.Unlock()

	timerLength := newDuration()

	ticker := time.NewTicker(time.Duration(10) * time.Millisecond)
	defer ticker.Stop()

	for {
		<-ticker.C

		cm.Lock()
		if cm.state == "leader" {
			cm.Unlock()
			return
		}

		if cm.currentTerm != termStarted {
			cm.Unlock()
			return
		}

		if time.Since(cm.lastReset) >= timerLength {
			cm.Unlock()
			cm.becomeCandidate()
			return
		}
		cm.Unlock()
	}
}

func newDuration() time.Duration {
	rand.Seed(time.Now().UnixNano())
	minDuration := 500

	return time.Duration(minDuration) * time.Millisecond
}

func (cm *CM) startElection() {}
