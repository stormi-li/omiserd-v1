package omiserd

import (
	"github.com/go-redis/redis/v8"
	omiconst "github.com/stormi-li/omiserd-v1/omiserd_const"
)

func NewClient(opts *redis.Options, nodeType omiconst.NodeType) *Client {
	if nodeType == omiconst.Server {
		return newClient(opts, omiconst.Prefix_Server, omiconst.Server)
	}
	if nodeType == omiconst.Web {
		return newClient(opts, omiconst.Prefix_Web, omiconst.Web)
	}
	return newClient(opts, omiconst.Prefix_Config, omiconst.Config)
}
