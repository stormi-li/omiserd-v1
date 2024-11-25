package omiserd

import (
	"github.com/go-redis/redis/v8"
)

func NewClient(opts *redis.Options, nodeType NodeType) *Client {
	if nodeType == Server {
		return newClient(opts, Prefix_Server, Server)
	}
	if nodeType == Web {
		return newClient(opts, Prefix_Web, Web)
	}
	return newClient(opts, Prefix_Config, Config)
}
