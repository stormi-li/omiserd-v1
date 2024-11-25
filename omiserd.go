package omiserd

import (
	"github.com/go-redis/redis/v8"
)

func NewClient(opts *redis.Options, nodeType NodeType) *Client {
	if nodeType == Server {
		return newClient(opts, prefix_Server, Server)
	}
	if nodeType == Web {
		return newClient(opts, prefix_Web, Web)
	}
	return newClient(opts, prefix_Config, Config)
}
