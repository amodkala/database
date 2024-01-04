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
    cm.lastReset = time.Now()
    cm.mu.Unlock()

	res := &proto.AppendEntriesResponse{
		Term:    cm.currentTerm,
		Success: false,
	}

	if req.Term > cm.currentTerm {
		cm.becomeFollower(req.Term)
	}

    if cm.leader != req.LeaderId {
        cm.leader = req.LeaderId
    }
    if cm.state != "follower" {
        cm.becomeFollower(req.Term)
    }

    if req.PrevLogIndex < int32(len(cm.log)) && req.PrevLogTerm == cm.log[req.PrevLogIndex].Term {
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

            newEntries := []Entry{}
            for _, entry := range req.Entries[newEntriesIndex:] {
                newEntries = append(newEntries, Entry{
                    Term: entry.Term,
                    Message: entry.Message,
                })
            }
            cm.log = append(cm.log[:logInsertIndex], newEntries...)
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
                    cm.commitChan <- entry.Message
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
    lastLogIndex := int32(len(cm.log) - 1)
    lastLogTerm := cm.log[lastLogIndex].Term
    cm.lastReset = time.Now()
	cm.mu.Unlock()

    var res *proto.RequestVoteResponse

    // candidate is a term ahead, become follower automatically
	if req.Term > cm.currentTerm {
		cm.becomeFollower(req.Term)
	}

    // performs the following checks for whether to grant the candidate a vote:
    // 1. 
	if (req.Term >= cm.currentTerm) &&
        (cm.votedFor == "" || cm.votedFor == req.CandidateId) &&
		(req.LastLogTerm > lastLogTerm ||
			(req.LastLogTerm == lastLogTerm && req.LastLogIndex >= lastLogIndex)) {

        cm.mu.Lock()
		cm.votedFor = req.CandidateId
        cm.currentTerm = req.Term
        cm.mu.Unlock()

        res = &proto.RequestVoteResponse{
            Term: cm.currentTerm,
            VoteGranted: true,
        }
	} else {
        res = &proto.RequestVoteResponse{
            Term: cm.currentTerm,
            VoteGranted: false,
        }
    }

	return res, nil
}
