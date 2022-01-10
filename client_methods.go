package raft

import (
	"fmt"
	"net"

	"github.com/amodkala/raft/proto"
	"google.golang.org/grpc"
)

func (cm *CM) Start(addr string) error {
	cm.Lock()
	cm.self = addr
	cm.Unlock()

	lis, err := net.Listen("tcp", cm.self)
	if err != nil {
		return fmt.Errorf("failed to listen -> %v", err)
	}
	server := grpc.NewServer()
	proto.RegisterRaftServer(server, cm)

	cm.becomeFollower(0)
	return server.Serve(lis)
}

//
// Replicate sends commands to the local consensus model,
// and if it is the leader it replicates the command across
// all peers
//
func (cm *CM) Replicate(commands []string) bool {
	cm.Lock()
	defer cm.Unlock()

	// TODO: support request rerouting to leader
	if cm.state == "leader" {
		for _, command := range commands {
			cm.log = append(cm.log, &proto.Entry{
				Term:    cm.currentTerm,
				Command: command,
			})
		}
		return true
	}
	return false
}
