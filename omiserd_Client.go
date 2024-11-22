package omiserd

import (
	"github.com/go-redis/redis/v8"
	discover "github.com/stormi-li/omiserd-v1/omiserd_discover"
	register "github.com/stormi-li/omiserd-v1/omiserd_register"
)

type Client struct {
	opts   *redis.Options
	prefix string
}

func newClient(opts *redis.Options, prefix string) *Client {
	return &Client{
		opts:   opts,
		prefix: prefix,
	}
}

func (c *Client) NewRegister(serverName, address string) *register.Register {
	return register.NewRegister(c.opts, serverName, address, c.prefix)
}

func (c *Client) NewDiscover() *discover.Discover {
	return discover.NewDiscover(c.opts, c.prefix)
}
