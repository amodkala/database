package raft

import (
	"context"
    "log"

	"github.com/amodkala/raft/proto"
)

func (cm *CM) sendHeartbeats() {

	cm.mu.Lock()
	heartbeatTerm := cm.currentTerm
	cm.mu.Unlock()

	for id, peer := range cm.peers {
		go func(id int, peer proto.RaftClient) {
			cm.mu.Lock()
			nextIndex := cm.nextIndex[id]
			prevLogIndex := nextIndex - 1
			prevLogTerm := cm.log[prevLogIndex].Term
			entries := []*proto.Entry{}
            cm.mu.Unlock()

            for _, entry := range cm.log[nextIndex:] {
                entries = append(entries, &proto.Entry{
                    Term: entry.Term,
                    Message: entry.Message,
                })
            }

			req := &proto.AppendEntriesRequest{
				Term:         heartbeatTerm,
				LeaderId:     cm.self,
				PrevLogIndex: prevLogIndex,
				PrevLogTerm:  prevLogTerm,
				Entries:      entries,
				LeaderCommit: cm.commitIndex,
			}

			res, err := peer.AppendEntries(context.Background(), req)
            if err != nil {
                log.Printf("%v", err)
            }

            if res.Term > heartbeatTerm {
                cm.becomeFollower(res.Term)
                return
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
							cm.mu.Lock()
							entries := cm.log[cm.lastApplied+1 : cm.commitIndex+1]
							cm.lastApplied = cm.commitIndex
							cm.mu.Unlock()

							for _, entry := range entries {
                                if len(entry.Message) > 0 {
                                    cm.commitChan <- entry.Message
                                }
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
