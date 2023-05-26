package client

import (
    "github.com/amodkala/db/pkg/proto"
    "github.com/amodkala/db/pkg/raft"
)

struct Client {
    cm raft.CM
    commitChan chan proto.Entry
    addr []byte
}

func New(addr []byte) *Client {
    commitChan := make(chan proto.Entry)
    return &Client {
        cm: raft.New(commitChan),
        commitChan: commitChan,
        addr: addr,
    }
}

// Start has the responsibility of setting up filesystem locations for 
// the raft module's log, plus the actual database, using the external lsmtree
// it then has to start the client's embedded raft module
func (c *Client) Start(addr string) { 
    c.cm.Start(c.addr)
}

// Write is responsible for adding new entries to the database by way of
// the client's raft module if it is the leader, or redirecting the request
// to the cluster leader's client if not
func (c *Client) Write(entries []proto.Entry) {}

// Read returns the value associated with a key directly from the database,
// without needing to interact with the raft cluster.
func (c *Client) Read(key string) {}

func (c *Client) Delete(key string) {}
