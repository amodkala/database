package raft

import (
	"context"
	"time"

	"github.com/amodkala/raft/proto"
)

//
// AppendEntries implements a client response to an AppendEntries RPC received from a peer
//
func (cm *CM) AppendEntries(ctx context.Context, req *proto.AppendEntriesRequest) (*proto.AppendEntriesResponse, error) {
	cm.Lock()
	defer cm.Unlock()

	res := &proto.AppendEntriesResponse{
		Term:    cm.currentTerm,
		Success: false,
	}

	if req.Term > cm.currentTerm {
		cm.becomeFollower(req.Term)
	}

	if req.Term == cm.currentTerm {
		if cm.state != "follower" {
			cm.becomeFollower(req.Term)
		}
		cm.lastReset = time.Now()
		res.Success = true
	}

	return res, nil
}

//
// RequestVote implements a client response to a RequestVote RPC received from a peer
//
func (cm *CM) RequestVote(ctx context.Context, req *proto.RequestVoteRequest) (*proto.RequestVoteResponse, error) {
	cm.Lock()
	defer cm.Unlock()

	res := &proto.RequestVoteResponse{
		Term: cm.currentTerm,
	}

	if req.Term > cm.currentTerm {
		cm.becomeFollower(req.Term)
	}

	if req.Term == cm.currentTerm &&
		(cm.votedFor == "" || cm.votedFor == req.CandidateId) {
		res.VoteGranted = true
		cm.votedFor = req.CandidateId
		cm.lastReset = time.Now()
	} else {
		res.VoteGranted = false
	}

	return res, nil
}
