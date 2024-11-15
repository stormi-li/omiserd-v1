package omiserd

import "github.com/go-redis/redis/v8"

func NewClient(opts *redis.Options, nodeType NodeType) *Client {
	if nodeType == Server {
		return newClient(opts, string(Server), Prefix_Server)
	}
	if nodeType == Web {
		return newClient(opts, string(Web), Prefix_Web)
	}
	return newClient(opts, string(Config), Prefix_Config)
}
