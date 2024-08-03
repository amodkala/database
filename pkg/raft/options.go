package raft

import "fmt"

type CMOpts func(*CM)

func WithLocalPeers(peers []string) CMOpts {
	return func(cm *CM) {
		for _, peer := range peers {
			if err := cm.addPeer(peer); err != nil {
				fmt.Println(err)
			}
		}
	}
}
