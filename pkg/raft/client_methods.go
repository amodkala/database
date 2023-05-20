package raft

import (
	"log"
	"net"

	"github.com/amodkala/db/pkg/proto"
	"google.golang.org/grpc"
)

func (cm *CM) Start(addr string) {
	cm.mu.Lock()
	cm.self = addr
	cm.mu.Unlock()

	lis, err := net.Listen("tcp", cm.self)
	if err != nil {
		log.Fatalf("failed to listen -> %v", err)
	}
	server := grpc.NewServer()
	proto.RegisterRaftServer(server, cm)

	cm.becomeFollower(0)
	log.Printf("starting consensus module at %s\n", addr)
	server.Serve(lis)
}

//
// Replicate sends commands to the local consensus model,
// and if it is the leader it replicates the command across
// all peers
//
func (cm *CM) Replicate(entries []Entry) bool {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// TODO: support request rerouting to leader
	if cm.state == "leader" {
		for _, entry := range entries {
			log.Printf("%s adding entry %v to log\n", cm.self, entry)
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
