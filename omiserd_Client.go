package omiserd

import (
	"github.com/go-redis/redis/v8"
	omiconst "github.com/stormi-li/omiserd-v1/omiserd_const"
	discover "github.com/stormi-li/omiserd-v1/omiserd_discover"
	register "github.com/stormi-li/omiserd-v1/omiserd_register"
)

type Client struct {
	opts     *redis.Options
	prefix   string
	NodeType omiconst.NodeType
}

func newClient(opts *redis.Options, prefix string, nodeType omiconst.NodeType) *Client {
	return &Client{
		opts:     opts,
		prefix:   prefix,
		NodeType: nodeType,
	}
}

func (c *Client) NewRegister(serverName, address string) *register.Register {
	return register.NewRegister(c.opts, serverName, address, c.prefix, c.NodeType)
}

func (c *Client) NewDiscover() *discover.Discover {
	return discover.NewDiscover(c.opts, c.prefix)
}
