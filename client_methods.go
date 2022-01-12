package raft

import (
	"log"
	"net"

	"github.com/amodkala/raft/proto"
	"google.golang.org/grpc"
)

func (cm *CM) Start(addr string) {
	cm.Lock()
	cm.self = addr
	cm.Unlock()

	lis, err := net.Listen("tcp", cm.self)
	if err != nil {
		log.Fatalf("failed to listen -> %v", err)
	}
	server := grpc.NewServer()
	proto.RegisterRaftServer(server, cm)

	cm.becomeFollower(0)
	server.Serve(lis)
}

//
// Replicate sends commands to the local consensus model,
// and if it is the leader it replicates the command across
// all peers
//
func (cm *CM) Replicate(entries []Entry) bool {
	cm.Lock()
	defer cm.Unlock()

	// TODO: support request rerouting to leader
	if cm.state == "leader" {
		for _, entry := range entries {
			cm.log = append(cm.log, &proto.Entry{
				Term:  cm.currentTerm,
				Key:   string(entry.Key),
				Value: string(entry.Value),
			})
		}
		return true
	}
	return false
}
