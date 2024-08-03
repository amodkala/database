package raft

import (
    "fmt"
    "log"
    "net"
    "os"
    "strconv"
    "time"

    "google.golang.org/grpc"

    "github.com/amodkala/raft/pkg/common"
    tx "github.com/amodkala/raft/pkg/transaction"
)

const (
    serviceNameVar = "RAFT_SERVICE_NAME"
    serviceNameDefault = "raft-peers.default.svc.cluster.local"

    retryPolicy = `{
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
)

func (cm *CM) Start(raftAddress string, opts ...CMOpts) error {

    cm.raftAddress = raftAddress

    lis, err := net.Listen("tcp", raftAddress)
    if err != nil {
        return fmt.Errorf("failed to listen -> %v", err)
    }
    server := grpc.NewServer()
    RegisterRaftServer(server, cm)

    for _, opt := range opts {
        opt(cm)
    }

    startTerm := uint32(0)

    if err := cm.log.Write(&common.Entry{
        RaftTerm: startTerm,
    }); err != nil {
        return fmt.Errorf("error while writing initial log entry: %w", err)
    }

    go func() {
        cm.errChan <- server.Serve(lis)
    }()

    peerTicker := time.NewTicker(time.Duration(5) * time.Second)
    defer peerTicker.Stop()

    cm.becomeFollower(startTerm)

    for {
        select {
            // refresh the peer list every 5 seconds
        case <-peerTicker.C:
            if err := cm.refreshPeers(); err != nil {
                cm.errChan <- err
            }
            // listen for any unrecoverable errors on errChan
        case err := <-cm.errChan:
            return fmt.Errorf("Consensus Module encountered error -> %v", err)
        default:
        }
    }
}

func (cm *CM) Leader() string { return cm.leader }

// Replicate sends commands to the local consensus model,
// and if it is the leader it replicates the command across
// all peers
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
    txChan := make(chan *common.Entry)
    cm.txChans[id] = txChan

    entries := tx.Entries()

    if err := cm.log.Write(entries...); err != nil {
        return "", fmt.Errorf("error writing entries to raft consensus log: %v", err)
    }

    cm.lastApplied += tx.Length()

    return "", nil
}

func (cm *CM) refreshPeers() error {
    serviceName, ok := os.LookupEnv(serviceNameVar)
    if !ok {
        serviceName = serviceNameDefault
    }
    _, srvRecords, err := net.LookupSRV("raft", "tcp", serviceName)
    if err != nil {
        return fmt.Errorf("failed to lookup peer SRV records: %w", err)
    }

    peers := make([]string, 0, len(srvRecords))
    for _, srv := range srvRecords {
        peers = append(peers, net.JoinHostPort(srv.Target, strconv.Itoa(int(srv.Port))))
    }

    cm.mu.Lock()
    defer cm.mu.Unlock()

    // Track current peers to identify removed ones
    currentPeers := make(map[string]struct{})
    for _, peer := range peers {
        currentPeers[peer] = struct{}{}
        if _, ok := cm.peerIndexes[peer]; ok {
            // this peer is already accounted for
            continue
        }

        conn, err := grpc.Dial(
            peer,
            grpc.WithInsecure(),
            grpc.WithDefaultServiceConfig(retryPolicy),
        )
        if err != nil {
            log.Printf("failed to connect to peer %s: %v", peer, err)
            continue // Skip this peer but continue with others
        }

        client := NewRaftClient(conn)
        cm.peers = append(cm.peers, client)
        cm.peerIndexes[peer] = len(cm.peers) - 1
        cm.nextIndex = append(cm.nextIndex, 0) // temporary, changed on election
        cm.matchIndex = append(cm.matchIndex, 0)
    }

    // Claude recommended O(n^2) algorithm below. GROSS

    // Remove peers that are no longer in the list
    for peer, index := range cm.peerIndexes {
        if _, ok := currentPeers[peer]; !ok {
            // Close the connection
            // cm.peers[index].Close() // Assuming the RaftClient has a Close method

            // Remove the peer
            cm.peers = append(cm.peers[:index], cm.peers[index+1:]...)
            cm.nextIndex = append(cm.nextIndex[:index], cm.nextIndex[index+1:]...)
            cm.matchIndex = append(cm.matchIndex[:index], cm.matchIndex[index+1:]...)
            delete(cm.peerIndexes, peer)

            // Update indexes for peers that have shifted
            for p, i := range cm.peerIndexes {
                if i > index {
                    cm.peerIndexes[p] = i - 1
                }
            }
        }
    }

    return nil
}
