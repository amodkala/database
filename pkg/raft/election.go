package raft

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/amodkala/db/pkg/proto"
)

func (cm *CM) startElectionTimer() {
	cm.mu.Lock()
	termStarted := cm.currentTerm
	cm.mu.Unlock()

	timerLength := newDuration()

	ticker := time.NewTicker(time.Duration(10) * time.Millisecond)
	defer ticker.Stop()

	for {
		<-ticker.C

		cm.mu.Lock()
		if cm.state == "leader" {
			cm.mu.Unlock()
			return
		}

		if cm.currentTerm != termStarted {
			cm.mu.Unlock()
			return
		}

		if time.Since(cm.lastReset) >= timerLength {
			cm.mu.Unlock()
			cm.becomeCandidate()
			return
		}
		cm.mu.Unlock()
	}
}

func newDuration() time.Duration {
	rand.Seed(time.Now().UnixNano())
	minDuration := 500

	return time.Duration(minDuration) * time.Millisecond
}

func (cm *CM) startElection() {
	cm.mu.Lock()
	electionTerm := cm.currentTerm
	cm.lastReset = time.Now()
	cm.votedFor = cm.self
	cm.mu.Unlock()

	votes := 1

	for id, peer := range cm.peers {

		go func(id string, peer proto.RaftClient) {

			var lastLogIndex, lastLogTerm int32

			cm.mu.Lock()
			if len(cm.log) > 0 {
				lastLogIndex = int32(len(cm.log) - 1)
				lastLogTerm = cm.log[lastLogIndex].Term
			} else {
				lastLogIndex = -1
				lastLogTerm = -1
			}

			req := &proto.RequestVoteRequest{
				Term:         electionTerm,
				CandidateId:  cm.self,
				LastLogIndex: lastLogIndex,
				LastLogTerm:  lastLogTerm,
			}
			cm.mu.Unlock()

			if res, err := peer.RequestVote(context.Background(), req); err == nil {
				cm.mu.Lock()
				defer cm.mu.Unlock()

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
						log.Printf("%s got vote from %s\n", cm.self, id)
					}

				}
			}
		}(id, peer)
	}

	if votes > len(cm.peers)/2 {
		cm.becomeLeader()
		return
	}

	go cm.startElectionTimer()
}
