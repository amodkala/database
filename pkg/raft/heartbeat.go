package raft

import (
	"context"

	"github.com/amodkala/raft/pkg/common"
)

func (cm *CM) sendHeartbeats() {

	cm.mu.Lock()
	heartbeatTerm := cm.currentTerm
	peers := len(cm.peers)
	cm.mu.Unlock()

	if peers == 0 {
		return
	}

	for id, peer := range cm.peers {

		go func(
			id int,
			peer RaftClient,
		) {
			var entries []*common.Entry

			cm.mu.Lock()
			nextIndex := cm.nextIndex[id]
			prevLogIndex := nextIndex - 1
			prevLogEntries, err := cm.log.Read(prevLogIndex)
			if err != nil {
				cm.errChan <- err
				return
			}
			prevLogEntry := prevLogEntries[0]
			prevLogTerm := prevLogEntry.RaftTerm
			entries, err = cm.log.Read(nextIndex, cm.log.Length()-1)
			if err != nil {
				cm.errChan <- err
				return
			}
			cm.mu.Unlock()

			req := &AppendEntriesRequest{
				Term:         heartbeatTerm,
				LeaderId:     cm.self,
				PrevLogIndex: prevLogIndex,
				PrevLogTerm:  prevLogTerm,
				Entries:      entries,
				LeaderCommit: cm.commitIndex,
			}

			res, err := peer.AppendEntries(context.Background(), req)
			if err != nil {
				cm.errChan <- err
				return
			}

			if res.Term > heartbeatTerm {
				cm.becomeFollower(res.Term)
				return
			}

			cm.mu.Lock()
			if cm.state == "leader" && res.Term == heartbeatTerm && nextIndex == cm.nextIndex[id] {
				if res.Success {
					cm.nextIndex[id] += uint32(len(entries))
					cm.matchIndex[id] = cm.nextIndex[id] - 1

					savedCommitIndex := cm.commitIndex
					for i := cm.commitIndex + 1; i < cm.log.Length(); i++ {
						entries, err := cm.log.Read(i)
						if err != nil {
							cm.errChan <- err
							return
						}
						entry := entries[0]
						if entry.RaftTerm == cm.currentTerm {
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
							entries, err := cm.log.Read(cm.lastApplied+1, cm.commitIndex)
							if err != nil {
								cm.errChan <- err
								return
							}
							cm.lastApplied = cm.commitIndex

							for _, entry := range entries {
								cm.txChans[entry.TxId] <- entry
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
