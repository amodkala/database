package raft

import (
	"context"
    "log"
    // "time"

	"github.com/amodkala/database/pkg/proto"
)

func (cm *CM) sendHeartbeats() {

	cm.mu.Lock()
	heartbeatTerm := cm.currentTerm
	cm.mu.Unlock()

	for id, peer := range cm.peers {
		go func(id int, peer proto.RaftClient) {
			entries := []*proto.Entry{}
			cm.mu.Lock()
			nextIndex := cm.nextIndex[id]
			prevLogIndex := nextIndex - 1
			prevLogTerm := cm.log[prevLogIndex].Term
            for _, entry := range cm.log[nextIndex:] {
                entries = append(entries, &proto.Entry{
                    Term: entry.Term,
                    Message: entry.Message,
                })
            }

            // log.Printf(`term %d -> leader %s sending heartbeats with params:
            // peer id: %d nextIndex: %d last log index: %d prevLogIndex: %d entries: %v
            // `, heartbeatTerm, cm.self, id, nextIndex, len(cm.log) - 1, prevLogIndex, entries)
            cm.mu.Unlock()

			req := &proto.AppendEntriesRequest{
				Term:         heartbeatTerm,
				LeaderId:     cm.self,
				PrevLogIndex: prevLogIndex,
				PrevLogTerm:  prevLogTerm,
				Entries:      entries,
				LeaderCommit: cm.commitIndex,
			}

            // start := time.Now()
			res, err := peer.AppendEntries(context.Background(), req)
            // log.Printf(`term %d -> leader %s got heartbeat response from peer %d in %d ms`,
            // cm.currentTerm, cm.self, id, time.Since(start).Milliseconds())
            if err != nil {
                log.Printf("%v", err)
                return
            }

            if res.Term > heartbeatTerm {
                cm.becomeFollower(res.Term)
                return
            }

            cm.mu.Lock()
			if cm.state == "leader" && res.Term == heartbeatTerm && nextIndex == cm.nextIndex[id] {
				if res.Success {
					cm.nextIndex[id] += int32(len(entries))
					cm.matchIndex[id] = cm.nextIndex[id] - 1

                    // log.Printf(`term %d -> leader %s updated peer %d
                    //     nextIndex: %d`, cm.currentTerm, cm.self, id, cm.nextIndex[id])

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
							entries := cm.log[cm.lastApplied+1 : cm.commitIndex+1]
							cm.lastApplied = cm.commitIndex

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
            cm.mu.Unlock()
		}(id, peer)
	}
}
