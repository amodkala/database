package raft

// Convenience struct for cluster peers that should allow for cleaner code
// (cm.nextIndex[peerID] => peer.nextIndex) and more verbose logging (current 
// impl doesn't store peer ids, just indices)
// Should also ease development of cluster membership changes whenever I add 
// that

import (
    "context"

	"google.golang.org/grpc"

    "github.com/amodkala/db/proto"
)

type Peer struct {
    ID string
    grpcClient proto.RaftClient
    NextIndex int
    MatchIndex int
}


func (p Peer) AppendEntries(ctx context.Context, in *proto.AppendEntriesRequest, opts ...grpc.CallOption) (*proto.AppendEntriesResponse, error) {
    return p.grpcClient.AppendEntries(ctx, in, opts...)
}
func (p Peer) RequestVote(ctx context.Context, in *proto.RequestVoteRequest, opts ...grpc.CallOption) (*proto.RequestVoteResponse, error) {
    return p.grpcClient.RequestVote(ctx, in, opts...)
}
