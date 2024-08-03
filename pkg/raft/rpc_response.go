package raft

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/amodkala/raft/pkg/common"
)

// AppendEntries implements a client response to an AppendEntries RPC received from a peer
func (cm *CM) AppendEntries(ctx context.Context, req *AppendEntriesRequest) (*AppendEntriesResponse, error) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.lastReset = time.Now()

	// start := time.Now()

	res := &AppendEntriesResponse{
		Term:    cm.currentTerm,
		Success: false,
	}

	if req.Term < cm.currentTerm {
		return res, nil
	}

	if req.Term > cm.currentTerm || cm.state != "follower" {
		cm.leader = req.LeaderId
		cm.becomeFollower(req.Term)
	}

	prevLogEntries, err := cm.log.Read(req.PrevLogIndex)
	if err != nil {
		return nil, fmt.Errorf("error reading previous log entries: %v", err)
	}
	prevLogEntry := prevLogEntries[0]

	if req.PrevLogIndex < cm.log.Length() && req.PrevLogTerm == prevLogEntry.RaftTerm {
		res.Success = true

		logInsertIndex := req.PrevLogIndex + 1
		newEntriesIndex := 0

		for {
			if logInsertIndex >= cm.log.Length() || newEntriesIndex >= len(req.Entries) {
				break
			}
			logInsertEntries, err := cm.log.Read(logInsertIndex)
			if err != nil {
				return nil, fmt.Errorf("error reading log insert entries: %v", err)
			}
			logInsertEntry := logInsertEntries[0]
			if logInsertEntry.RaftTerm != req.Entries[newEntriesIndex].RaftTerm {
				break
			}
			logInsertIndex++
			newEntriesIndex++
		}

		if newEntriesIndex < len(req.Entries) {
			log.Printf("term %d -> %s added entries to log %v\n", cm.currentTerm, cm.self, req.Entries)

			newEntries := []*common.Entry{}
			for _, entry := range req.Entries[newEntriesIndex:] {
				newEntries = append(newEntries, entry)
			}
			if err := cm.log.Write(newEntries...); err != nil {
				return nil, fmt.Errorf("error writing new entries to log: %v", err)
			}
		}

		if req.LeaderCommit > cm.commitIndex {
			cm.commitIndex = min(req.LeaderCommit, cm.log.Length()-1)
			if cm.commitIndex > cm.lastApplied {
				// indicate that log entries from index cm.lastApplied + 1 to
				// cm.commitIndex (inclusive) can be committed
				entries, err := cm.log.Read(cm.lastApplied+1, cm.commitIndex)
				if err != nil {
					return nil, fmt.Errorf("error reading entries for commitment: %v", err)
				}
				cm.lastApplied = cm.commitIndex

				for _, entry := range entries {
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

// RequestVote implements a client response to a RequestVote RPC received from a peer
func (cm *CM) RequestVote(ctx context.Context, req *RequestVoteRequest) (*RequestVoteResponse, error) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	lastLogIndex := cm.log.Length() - 1
	lastLogEntries, err := cm.log.Read(lastLogIndex)
	if err != nil {
		return nil, fmt.Errorf("error reading last log entry: %v", err)
	}
	lastLogEntry := lastLogEntries[0]
	lastLogTerm := lastLogEntry.RaftTerm
	cm.lastReset = time.Now()

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
			Term:        cm.currentTerm,
			VoteGranted: true,
		}
	} else {
		res = &RequestVoteResponse{
			Term:        cm.currentTerm,
			VoteGranted: false,
		}
	}

	log.Printf(`term %d -> %s responded to vote request in %d ms`, cm.currentTerm, cm.self, time.Since(start).Milliseconds())
	return res, nil
}
