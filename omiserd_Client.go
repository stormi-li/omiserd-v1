package omiserd

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/stormi-li/omipc-v1"
)

type Client struct {
	opts       *redis.Options
	namespace  string
	serverType string
}

func newClient(opts *redis.Options, serverType string, prefix string) *Client {
	return &Client{
		opts:       opts,
		namespace:  prefix,
		serverType: serverType,
	}
}

func (c *Client) NewRegister(serverName, address string) *Register {
	return &Register{
		redisClient:     redis.NewClient(c.opts),
		ServerName:      serverName,
		Address:         address,
		Data:            map[string]string{},
		namespace:       c.namespace,
		ctx:             context.Background(),
		omipcClient:     omipc.NewClient(c.opts),
		registerHandler: func(register *Register) {},
		messageHandler:  func(command, message string, register *Register) {},
	}
}

func (c *Client) NewDiscover() *Discover {
	discover := Discover{
		redisClient: redis.NewClient(c.opts),
		namespace:   c.namespace,
		ctx:         context.Background(),
	}
	discover.SetReconnectStrategy(func(address string, data map[string]string, err error) bool {
		return err != nil || address == "" || data["weight"] == "0" || data["weight"] == ""
	})
	discover.SetDiscoverStrategy(func(serverName string, discover *Discover) (string, map[string]string) {
		return discover.DiscoverByWeight(serverName)
	})
	return &discover
}
