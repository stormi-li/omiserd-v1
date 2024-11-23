package main

import (
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stormi-li/omiserd-v1"
	omiconst "github.com/stormi-li/omiserd-v1/omiserd_const"
	discover "github.com/stormi-li/omiserd-v1/omiserd_discover"
)

var redisAddr = "118.25.196.166:3934"
var password = "12982397StrongPassw0rd"

func main() {
	omiserdC := omiserd.NewClient(&redis.Options{Addr: redisAddr, Password: password}, omiconst.Web)
	d := omiserdC.NewDiscover()

	monitor := d.NewMonitor("web_service")

	listenHandleFunc := func(serverName, oldAddress string, discover *discover.Discover) string {
		if !discover.IsAlive(serverName, oldAddress) {
			addresses := discover.GetByWeight(serverName)
			if len(addresses) > 0 {
				return addresses[rand.IntN(len(addresses))]
			}
		}
		return ""
	}

	connectHandleFunc := func(address string) { fmt.Println(address) }

	monitor.ListenAndConnect(2*time.Second, listenHandleFunc, connectHandleFunc)
}
