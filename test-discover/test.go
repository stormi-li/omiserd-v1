package main

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/stormi-li/omiserd-v1"
)

var redisAddr = "118.25.196.166:3934"
var password = "12982397StrongPassw0rd"

func main() {
	discover := omiserd.NewClient(&redis.Options{Addr: redisAddr, Password: password}, omiserd.Config).NewDiscover()
	// discover.SetReconnectStrategy(func(address string, data map[string]string, err error) bool {
	// 	return err != nil
	// })
	discover.SetDiscoverStrategy(func(serverName string, discover *omiserd.Discover) (string, map[string]string) {
		return discover.DiscoverByWeight(serverName)
	})
	discover.DiscoverAndListen("redis1", func(address string, data map[string]string) {
		fmt.Println(address, data)
	})
}
