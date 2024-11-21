package main

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/stormi-li/omiserd-v1"
)

func main() {
	omiserdC := omiserd.NewClient(&redis.Options{Addr: "localhost:6379"}, omiserd.Web)
	d := omiserdC.NewDiscover()
	fmt.Println(d.DiscoverByName("web_service"))
}
