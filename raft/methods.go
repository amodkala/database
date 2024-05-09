package raft

import (
    "fmt"
	"log"
	"net"

	"github.com/amodkala/db/proto"
	"google.golang.org/grpc"
)

func (cm *CM) Start(opts ...CMOpts) error {
    lis, err := net.Listen("tcp", cm.self)
	if err != nil {
		return fmt.Errorf("failed to listen -> %v", err)
	}
	server := grpc.NewServer()
	proto.RegisterRaftServer(server, cm)

    for _, opt := range opts {
        opt(cm)
    }

    startTerm := int32(0)

    cm.log = append(cm.log, Entry{
        Term: startTerm,
        Message: []byte{},
    })
	cm.becomeFollower(startTerm)
    return fmt.Errorf("Consensus Module encountered error -> %v", server.Serve(lis))
}

// Helper function for client to ascertain whether its consensus module is 
// cluster leader
// This should be obsoleted by implementation of redirect
func (cm *CM) isLeader() bool {
    return cm.state == "leader" && cm.self == cm.leader
}

//
// Replicate sends commands to the local consensus model,
// and if it is the leader it replicates the command across
// all peers
//
func (cm *CM) Replicate(entries ...[]byte) (string, error) {
	cm.mu.Lock()
    defer cm.mu.Unlock()

    if !cm.isLeader() {
        return cm.leader, fmt.Errorf("Node is not leader")
    }

    for _, entry := range entries {
        log.Printf("%s adding entry %v to log\n", cm.self, entry)
        cm.log = append(cm.log, Entry{
            Term: cm.currentTerm,
            Message: entry,
        })
    }

    return "", nil
}

func (cm *CM) addPeer(addr string) error {

    retryPolicy := `{
        "methodConfig": [{
            "name": [{"service": "proto.Raft"}],
            "retryPolicy": {
                "MaxAttempts": 5,
                "InitialBackoff": "0.1s",
                "MaxBackoff": "0.1s",
                "BackoffMultiplier": 1,
                "RetryableStatusCodes": ["UNAVAILABLE"]
            }
        }]
    }`

    conn, err := grpc.Dial(
        addr,
        grpc.WithInsecure(),
        grpc.WithDefaultServiceConfig(retryPolicy),
    )
    if err != nil {
        return fmt.Errorf("failed to connect to peer %s: %w", addr, err)
    }
    client := proto.NewRaftClient(conn)

    cm.mu.Lock()
    cm.peers = append(cm.peers, client)
    cm.nextIndex = append(cm.nextIndex, 0) // temporary, changed on election
    cm.matchIndex = append(cm.matchIndex, 0)
    cm.mu.Unlock()

    return nil
}
