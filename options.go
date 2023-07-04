package raft

type CMOpts func(*CM) *CM

func WithPeers(peers []string) CMOpts {
    return func(cm *CM) *CM {
        for _, peer := range peers {
            cm.addPeer(peer)
        }
        return cm
    }
}
