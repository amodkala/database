package raft

import (
    "fmt"
	"net"

	"google.golang.org/grpc"

    "github.com/amodkala/database/pkg/common"
    tx "github.com/amodkala/database/pkg/transaction"
)

func (cm *CM) Start(opts ...CMOpts) error {
    lis, err := net.Listen("tcp", cm.self)
	if err != nil {
		return fmt.Errorf("failed to listen -> %v", err)
	}
	server := grpc.NewServer()
	RegisterRaftServer(server, cm)

    for _, opt := range opts {
        opt(cm)
    }

    startTerm := uint32(0)

    cm.log.Write(&common.Entry{
        RaftTerm: startTerm,
    })
	cm.becomeFollower(startTerm)
    
    go func() {
        cm.errChan <- server.Serve(lis)
    }()

    for {
        select {
        case err :=  <-cm.errChan:
            return fmt.Errorf("Consensus Module encountered error -> %v", err)
        default:
        }
    }
}

func (cm *CM) Length() uint32 { return cm.log.Length() }

//
// Replicate sends commands to the local consensus model,
// and if it is the leader it replicates the command across
// all peers
//
func (cm *CM) Replicate(tx tx.Tx) (string, error) {

	cm.mu.Lock()
    defer cm.mu.Unlock()

    if !(cm.state == "leader" && cm.self == cm.leader) {
        return cm.leader, fmt.Errorf("Node is not leader")
    }

    for _, entry := range tx.Entries() {
        entry.RaftTerm = cm.currentTerm
    }

    id := tx.ID()
    if _, ok := cm.txChans[id]; ok {
        return "", fmt.Errorf("transaction id %s is already in use", id)
    }
    txChan := make(chan *common.Entry)
    cm.txChans[id] = txChan

    if err := cm.log.Write(tx.Entries()...); err != nil {
        return "", fmt.Errorf("error writing entries to raft consensus log: %v", err)
    }

    cm.lastApplied += tx.Length()

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
    client := NewRaftClient(conn)

    cm.mu.Lock()
    cm.peers = append(cm.peers, client)
    cm.nextIndex = append(cm.nextIndex, 0) // temporary, changed on election
    cm.matchIndex = append(cm.matchIndex, 0)
    cm.mu.Unlock()

    return nil
}
