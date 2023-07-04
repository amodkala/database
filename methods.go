package raft

import (
    "fmt"
	"log"
	"net"

	"github.com/amodkala/raft/proto"
	"google.golang.org/grpc"
)

func (cm *CM) Start(addr string) error {
	cm.mu.Lock()
	cm.self = addr
    fmt.Println(cm.self, len(cm.peers))
	cm.mu.Unlock()

    lis, err := net.Listen("tcp", cm.self)
	if err != nil {
		log.Fatalf("failed to listen -> %v", err)
	}
	server := grpc.NewServer()
	proto.RegisterRaftServer(server, cm)

	log.Printf("starting consensus module at %s\n", addr)
	cm.becomeFollower(0)
	server.Serve(lis)

    return nil
}

//
// Replicate sends commands to the local consensus model,
// and if it is the leader it replicates the command across
// all peers
//
func (cm *CM) Replicate(entries ...Entry) (string, error) {
	cm.mu.Lock()

    if cm.state != "leader" {
        return cm.leader, fmt.Errorf("Node is not leader")
    }

    for _, entry := range entries {
        log.Printf("%s adding entry %v to log\n", cm.self, entry)
        cm.log = append(cm.log, &proto.Entry{
            Term:  cm.currentTerm,
            Key:   entry.Key,
            Value: entry.Value,
        })
    }

	cm.mu.Unlock()

    // TODO entries ready to be committed will appear here via cm.CommitChan

    return "", nil
}

func (cm *CM) addPeer(addr string) error {
    conn, err := grpc.Dial(addr)
    if err != nil {
        return fmt.Errorf("failed to connect to gRPC channel: %w", err)
    }
    client := proto.NewRaftClient(conn)

    cm.mu.Lock()
    defer cm.mu.Unlock()

    cm.peers = append(cm.peers, client)
    cm.nextIndex = append(cm.nextIndex, cm.lastApplied + 1)
    cm.matchIndex = append(cm.matchIndex, 0)
	log.Printf("%s added peer %s\n", cm.self, addr)

    return nil
}
