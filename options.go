package raft

import (
	"log"

	"github.com/amodkala/raft/proto"
	"google.golang.org/grpc"
)

type RaftOpt func(cm *CM)

//
// WithPeers takes a map[id]ip of peers and creates gRPC clients
// for each, then adds a slot for them where appropriate
//
func WithPeers(peers map[string]string) func(*CM) {
	return func(cm *CM) {
		for id, ip := range peers {
			conn, err := grpc.Dial(ip)
			if err != nil {
				log.Printf("couldn't connect to %s\n", id)
			}
			client := proto.NewRaftClient(conn)

			cm.peers[id] = client
			cm.nextIndex[id] = 0
			cm.matchIndex[id] = 0
		}
	}
}
