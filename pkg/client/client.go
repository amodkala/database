package client

import (
    "github.com/amodkala/db/raft"
)

struct Client {
    cm raft.CM
    addr string
}

func New() *Client {
    return &Client {
        cm: raft.New()
    }
}

func (c *Client) Start() { c.cm.Start(c.addr) }

func (c *Client) Write(key, value string) {}

func (c *Client) Read(key string) {}

func (c *Client) Delete(key string) {}
