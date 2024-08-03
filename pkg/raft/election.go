package raft

import (
	"context"
	"log"
	"math/rand"
	"time"
)

func (cm *CM) startElectionTimer() {
	cm.mu.Lock()
	termStarted := cm.currentTerm
	cm.lastReset = time.Now()
	cm.mu.Unlock()

	timerLength := newDuration(500)

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

	return time.Duration(min+rand.Intn(min)) * time.Millisecond
}

func (cm *CM) startElection() {
	cm.mu.Lock()
	electionTerm := cm.currentTerm
	cm.lastReset = time.Now()
	cm.votedFor = cm.self
	peers := len(cm.peers)
	cm.mu.Unlock()

	if peers == 0 {
		cm.becomeLeader()
		return
	}

	votes := 1

	voteChan := make(chan bool, len(cm.peers))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Ensure cancellation at the end of the function

	for id, peer := range cm.peers {

		go func(ctx context.Context, id int, peer RaftClient, voteChan chan bool) {
			defer func() {
				// Close channel only when all goroutines are done
				if r := recover(); r != nil && r == "send on closed channel" {
					log.Println("Channel was closed, vote result not sent")
				}
			}()

			cm.mu.Lock()
			lastLogIndex := cm.log.Length() - 1
			lastLogEntries, err := cm.log.Read(lastLogIndex)
			if err != nil {
				cm.errChan <- err
				return
			}
			lastLogEntry := lastLogEntries[0]
			lastLogTerm := lastLogEntry.RaftTerm
			cm.mu.Unlock()

			req := &RequestVoteRequest{
				Term:         electionTerm,
				CandidateId:  cm.self,
				LastLogIndex: lastLogIndex,
				LastLogTerm:  lastLogTerm,
			}

			res, err := peer.RequestVote(ctx, req)
			if err != nil {
				log.Printf("%v", err)
				voteChan <- false
				cm.errChan <- err
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
			case <-ctx.Done():
				log.Printf("Context cancelled, stopping vote request")
				return
			case voteChan <- res.Term == electionTerm && res.VoteGranted:
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
