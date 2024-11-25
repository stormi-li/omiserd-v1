package omiserd

import (
	"github.com/go-redis/redis/v8"
)

type Client struct {
	opts     *redis.Options
	prefix   string
	NodeType NodeType
}

func newClient(opts *redis.Options, prefix string, nodeType NodeType) *Client {
	return &Client{
		opts:     opts,
		prefix:   prefix,
		NodeType: nodeType,
	}
}

func (c *Client) NewRegister(serverName, address string) *Register {
	return NewRegister(c.opts, serverName, address, c.prefix, c.NodeType)
}

func (c *Client) NewDiscover() *Discover {
	return NewDiscover(c.opts, c.prefix)
}
