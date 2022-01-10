package raft

import (
	"context"

	"github.com/amodkala/raft/proto"
)

//
// AppendEntries implements a client response to an AppendEntries RPC received from a peer
//
func (cm *CM) AppendEntries(context.Context, *proto.AppendEntriesRequest) (*proto.AppendEntriesResponse, error) {
	return nil, nil
}

//
// RequestVote implements a client response to a RequestVote RPC received from a peer
//
func (cm *CM) RequestVote(context.Context, *proto.RequestVoteRequest) (*proto.RequestVoteResponse, error) {
	return nil, nil
}
