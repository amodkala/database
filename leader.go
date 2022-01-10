package raft

import (
	"context"

	"github.com/amodkala/raft/proto"
)

func (cm *CM) sendHeartbeats() {

	cm.Lock()
	heartbeatTerm := cm.currentTerm
	cm.Unlock()

	req := &proto.AppendEntriesRequest{
		Term:     heartbeatTerm,
		LeaderId: cm.self,
	}

	for _, peer := range cm.peers {
		go func(peer proto.RaftClient) {
			if res, err := peer.AppendEntries(context.Background(), req); err == nil {
				cm.Lock()
				defer cm.Unlock()
				if res.Term > heartbeatTerm {
					cm.becomeFollower(res.Term)
					return
				}
			}
		}(peer)
	}
}
