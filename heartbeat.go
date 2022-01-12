package raft

import (
	"context"

	"github.com/amodkala/raft/proto"
)

func (cm *CM) sendHeartbeats() {

	cm.Lock()
	heartbeatTerm := cm.currentTerm
	cm.Unlock()

	for id, peer := range cm.peers {
		go func(id string, peer proto.RaftClient) {
			cm.Lock()
			nextIndex := cm.nextIndex[id]
			prevLogIndex := nextIndex - 1
			prevLogTerm := cm.log[prevLogIndex].Term
			entries := cm.log[nextIndex:]

			req := &proto.AppendEntriesRequest{
				Term:         heartbeatTerm,
				LeaderId:     cm.self,
				PrevLogIndex: prevLogIndex,
				PrevLogTerm:  prevLogTerm,
				Entries:      entries,
				LeaderCommit: cm.commitIndex,
			}
			cm.Unlock()

			res, err := peer.AppendEntries(context.Background(), req)
			// TODO: handle case where res == nil
			if err == nil {
				cm.Lock()
				defer cm.Unlock()
				if res.Term > heartbeatTerm {
					cm.becomeFollower(res.Term)
					return
				}
			}

			if cm.state == "leader" && res.Term == heartbeatTerm {
				if res.Success {
					cm.nextIndex[id] += int32(len(entries))
					cm.matchIndex[id] = cm.nextIndex[id] - 1

					savedCommitIndex := cm.commitIndex
					for i := cm.commitIndex + 1; i < int32(len(cm.log)); i++ {
						if cm.log[i].Term == cm.currentTerm {
							matchCount := 1
							for id := range cm.peers {
								if cm.matchIndex[id] >= i {
									matchCount++
								}
							}
							if matchCount > len(cm.peers)/2 {
								cm.commitIndex = i
							}
						}
					}
					if savedCommitIndex != cm.commitIndex {
						if cm.commitIndex > cm.lastApplied {
							// tell client these have been committed
							cm.Lock()
							entries := cm.log[cm.lastApplied+1 : cm.commitIndex+1]
							cm.lastApplied = cm.commitIndex
							cm.Unlock()

							for _, entry := range entries {
								result := Entry{Key: []byte(entry.Key), Value: []byte(entry.Value)}
								cm.CommitChan <- result
							}
						}
					}
				} else {
					cm.nextIndex[id] -= 1
				}
			}
		}(id, peer)
	}
}
