package raft

import (
	"context"
    "log"
	"time"

    "github.com/amodkala/database/pkg/common"
)

//
// AppendEntries implements a client response to an AppendEntries RPC received from a peer
//
func (cm *CM) AppendEntries(ctx context.Context, req *AppendEntriesRequest) (*AppendEntriesResponse, error) {
	cm.mu.Lock()
    cm.lastReset = time.Now()
    cm.mu.Unlock()

    // start := time.Now()

	res := &AppendEntriesResponse{
		Term:    cm.currentTerm,
		Success: false,
	}

    if req.Term < cm.currentTerm {
        return res, nil
    }

	if req.Term > cm.currentTerm || cm.state != "follower"{
		cm.becomeFollower(req.Term)
	}

    if cm.leader != req.LeaderId {
        cm.leader = req.LeaderId
    }

    if req.PrevLogIndex < uint32(len(cm.log)) && req.PrevLogTerm == cm.log[req.PrevLogIndex].RaftTerm {
        res.Success = true

        logInsertIndex := req.PrevLogIndex + 1
        newEntriesIndex := 0

        for {
            if logInsertIndex >= uint32(len(cm.log)) || newEntriesIndex >= len(req.Entries) {
                break
            }
            if cm.log[logInsertIndex].RaftTerm != req.Entries[newEntriesIndex].RaftTerm {
                break
            }
            logInsertIndex++
            newEntriesIndex++
        }

        if newEntriesIndex < len(req.Entries) {
            log.Printf("term %d -> %s added entries to log %v\n", cm.currentTerm, cm.self, req.Entries)

            newEntries := []common.Entry{}
            for _, entry := range req.Entries[newEntriesIndex:] {
                newEntries = append(newEntries, *entry)
            }
            cm.log = append(cm.log[:logInsertIndex], newEntries...)
        }

        if req.LeaderCommit > cm.commitIndex {
            cm.commitIndex = min(req.LeaderCommit, uint32(len(cm.log)-1))
            if cm.commitIndex > cm.lastApplied {
                // tell client these have been committed
                cm.mu.Lock()
                entries := cm.log[cm.lastApplied+1 : cm.commitIndex+1]
                cm.lastApplied = cm.commitIndex
                cm.mu.Unlock()

                for _, entry := range entries {
                    // again check whether this entry can be committed
                    cm.commitChan <- entry
                }
            }
        }

    }

    // log.Printf(`term %d -> %s responded to heartbeat in %d ms`, cm.currentTerm, cm.self, time.Since(start).Milliseconds())
	return res, nil
}

func min(a, b uint32) uint32 {
	if a > b {
		return b
	}
	return a
}

//
// RequestVote implements a client response to a RequestVote RPC received from a peer
//
func (cm *CM) RequestVote(ctx context.Context, req *RequestVoteRequest) (*RequestVoteResponse, error) {
	cm.mu.Lock()
    lastLogIndex := int32(len(cm.log) - 1)
    lastLogTerm := cm.log[lastLogIndex].Term
    cm.lastReset = time.Now()
	cm.mu.Unlock()

    start := time.Now()

    var res *RequestVoteResponse

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

        res = &RequestVoteResponse{
            Term: cm.currentTerm,
            VoteGranted: true,
        }
	} else {
        res = &RequestVoteResponse{
            Term: cm.currentTerm,
            VoteGranted: false,
        }
    }

    log.Printf(`term %d -> %s responded to vote request in %d ms`, cm.currentTerm, cm.self, time.Since(start).Milliseconds())
	return res, nil
}
