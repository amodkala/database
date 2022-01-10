package raft

import (
	"context"
	"math/rand"
	"time"

	"github.com/amodkala/raft/proto"
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

func (cm *CM) startElection() {
	cm.Lock()
	electionTerm := cm.currentTerm
	cm.lastReset = time.Now()
	cm.votedFor = cm.self
	cm.Unlock()

	votes := 1

	for _, peer := range cm.peers {

		go func(peer proto.RaftClient) {
			req := &proto.RequestVoteRequest{
				Term:        electionTerm,
				CandidateId: cm.self,
			}
			if res, err := peer.RequestVote(context.Background(), req); err == nil {
				cm.Lock()
				defer cm.Unlock()

				if cm.state != "candidate" {
					return
				}

				if res.Term > electionTerm {
					cm.becomeFollower(res.Term)
					return
				}

				if res.Term == electionTerm {
					if res.VoteGranted {
						votes += 1
						if votes > len(cm.peers)/2 {
							cm.becomeLeader()
							return
						}
					}

				}
			}
		}(peer)
	}

	go cm.startElectionTimer()
}
