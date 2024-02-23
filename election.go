package raft

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/amodkala/raft/proto"
)

func (cm *CM) startElectionTimer() {
	cm.mu.Lock()
	termStarted := cm.currentTerm
    cm.lastReset = time.Now()
	cm.mu.Unlock()

	timerLength := newDuration(1500)

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

func newDuration(min int) time.Duration {
	rand.Seed(time.Now().UnixNano())

	return time.Duration(min + rand.Intn(min)) * time.Millisecond
}

func (cm *CM) startElection() {
	cm.mu.Lock()
	electionTerm := cm.currentTerm
	cm.lastReset = time.Now()
	cm.votedFor = cm.self
	cm.mu.Unlock()

	votes := 1

    voteChan := make(chan bool, len(cm.peers))
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel() // Ensure cancellation at the end of the function

	for id, peer := range cm.peers {

		go func(ctx context.Context, id int, peer proto.RaftClient, voteChan chan bool) {
            defer func() {
                // Close channel only when all goroutines are done
                if r := recover(); r != nil && r == "send on closed channel" {
                    log.Println("Channel was closed, vote result not sent")
                }
            }()

			cm.mu.Lock()
            lastLogIndex := int32(len(cm.log) - 1)
            lastLogTerm := cm.log[lastLogIndex].Term
			cm.mu.Unlock()

			req := &proto.RequestVoteRequest{
				Term:         electionTerm,
				CandidateId:  cm.self,
				LastLogIndex: lastLogIndex,
				LastLogTerm:  lastLogTerm,
			}

            start := time.Now()
			res, err := peer.RequestVote(ctx, req)
            log.Printf(`term %d -> candidate %s got vote request response from peer %d in %d ms`,
            cm.currentTerm, cm.self, id, time.Since(start).Milliseconds())
            if err != nil {
                log.Printf("%v", err)
                voteChan <- false
                return
            }

            cm.mu.Lock()
            defer cm.mu.Unlock()

            if cm.state != "candidate" {
                return
            }

            if res.Term > electionTerm {
                cm.becomeFollower(res.Term)
                return
            }

            select {
            case voteChan <- res.Term == electionTerm && res.VoteGranted:
            case <-ctx.Done():
                log.Printf("Context cancelled, stopping vote request")
                return
            }
		}(ctx, id, peer, voteChan)
	}

    // Collect responses until the majority is reached or all votes are counted
    for i := 0; i < len(cm.peers); i++ {
        if vote := <-voteChan; vote {
            votes++
            if votes > len(cm.peers)/2 {
                cm.becomeLeader()
                cancel() // Cancel remaining goroutines
                break
            }
        }
    }

    if cm.state == "candidate" {
        go cm.startElectionTimer()
    }
}
