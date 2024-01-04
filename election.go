package raft

import (
	"context"
	"log"
	"math/rand"
    "sync"
	"time"

	"github.com/amodkala/raft/proto"
)

func (cm *CM) startElectionTimer() {
	cm.mu.Lock()
	termStarted := cm.currentTerm
    cm.lastReset = time.Now()
	cm.mu.Unlock()

	timerLength := newDuration()

	ticker := time.NewTicker(time.Duration(10) * time.Millisecond)
	defer ticker.Stop()

	for {
		<-ticker.C

		cm.mu.Lock()
		if cm.state == "leader" || cm.currentTerm != termStarted {
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
	minDuration := 1500
    random := rand.Intn(1500)

	return time.Duration(minDuration + random) * time.Millisecond
}

func (cm *CM) startElection() {
	cm.mu.Lock()
    // on election start, the candidate increments its current term and votes for itself
	electionTerm := cm.currentTerm
	cm.lastReset = time.Now()
	cm.votedFor = cm.self
	cm.mu.Unlock()

	votes := 1

    var wg sync.WaitGroup

	for id, peer := range cm.peers {

        wg.Add(1)

		go func(id int, peer proto.RaftClient) {
            defer wg.Done()

			var lastLogIndex, lastLogTerm int32

			cm.mu.Lock()
            lastLogIndex = int32(len(cm.log) - 1)
            lastLogTerm = cm.log[lastLogIndex].Term
			cm.mu.Unlock()

			req := &proto.RequestVoteRequest{
				Term:         electionTerm,
				CandidateId:  cm.self,
				LastLogIndex: lastLogIndex,
				LastLogTerm:  lastLogTerm,
			}

			res, err := peer.RequestVote(context.Background(), req)
            if err != nil {
                log.Printf("%v", err)
            }

            cm.mu.Lock()

            if cm.state != "candidate" {
                cm.mu.Unlock()
                return
            }

            if res.Term > electionTerm {
                cm.becomeFollower(res.Term)
                cm.mu.Unlock()
                return
            }

            if res.Term == electionTerm && res.VoteGranted {
                votes += 1
            }

            cm.mu.Unlock()
		}(id, peer)
	}

    wg.Wait()

	if votes > len(cm.peers)/2 {
		cm.becomeLeader()
		return
	}

	cm.startElectionTimer()
}
