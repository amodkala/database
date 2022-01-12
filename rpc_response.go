package raft

import (
	"context"
	"log"
	"time"

	"github.com/amodkala/raft/proto"
)

//
// AppendEntries implements a client response to an AppendEntries RPC received from a peer
//
func (cm *CM) AppendEntries(ctx context.Context, req *proto.AppendEntriesRequest) (*proto.AppendEntriesResponse, error) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	res := &proto.AppendEntriesResponse{
		Term:    cm.currentTerm,
		Success: false,
	}

	if req.Term > cm.currentTerm {
		cm.becomeFollower(req.Term)
	}

	if req.Term == cm.currentTerm {
		if cm.leader != req.LeaderId {
			cm.leader = req.LeaderId
		}
		if cm.state != "follower" {
			cm.becomeFollower(req.Term)
		}
		cm.lastReset = time.Now()

		if req.PrevLogIndex == -1 ||
			(req.PrevLogIndex < int32(len(cm.log)) && req.PrevLogTerm == cm.log[req.PrevLogIndex].Term) {
			res.Success = true

			logInsertIndex := req.PrevLogIndex + 1
			newEntriesIndex := 0

			for {
				if logInsertIndex >= int32(len(cm.log)) || newEntriesIndex >= len(req.Entries) {
					break
				}
				if cm.log[logInsertIndex].Term != req.Entries[newEntriesIndex].Term {
					break
				}
				logInsertIndex++
				newEntriesIndex++
			}

			if newEntriesIndex < len(req.Entries) {
				log.Printf("%s added entries to log %v\n", cm.self, req.Entries)
				cm.log = append(cm.log[:logInsertIndex], req.Entries[newEntriesIndex:]...)
			}

			if req.LeaderCommit > cm.commitIndex {
				cm.commitIndex = min(req.LeaderCommit, int32(len(cm.log)-1))
				if cm.commitIndex > cm.lastApplied {
					// tell client these have been committed
					cm.mu.Lock()
					entries := cm.log[cm.lastApplied+1 : cm.commitIndex+1]
					cm.lastApplied = cm.commitIndex
					cm.mu.Unlock()

					for _, entry := range entries {
						result := Entry{Key: entry.Key, Value: entry.Value}
						cm.CommitChan <- result
					}
				}
			}

		}
	}

	return res, nil
}

func min(a, b int32) int32 {
	if a > b {
		return b
	}
	return a
}

//
// RequestVote implements a client response to a RequestVote RPC received from a peer
//
func (cm *CM) RequestVote(ctx context.Context, req *proto.RequestVoteRequest) (*proto.RequestVoteResponse, error) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	var lastLogIndex, lastLogTerm int32

	if len(cm.log) > 0 {
		lastLogIndex = int32(len(cm.log) - 1)
		lastLogTerm = cm.log[lastLogIndex].Term
	} else {
		lastLogIndex = -1
		lastLogTerm = -1
	}

	res := &proto.RequestVoteResponse{
		Term: cm.currentTerm,
	}

	if req.Term > cm.currentTerm {
		cm.becomeFollower(req.Term)
	}

	if req.Term == cm.currentTerm &&
		(cm.votedFor == "" || cm.votedFor == req.CandidateId) &&
		(req.LastLogTerm > lastLogTerm ||
			(req.LastLogTerm == lastLogTerm && req.LastLogIndex >= lastLogIndex)) {
		res.VoteGranted = true
		cm.votedFor = req.CandidateId
		cm.lastReset = time.Now()
	} else {
		res.VoteGranted = false
	}

	return res, nil
}
