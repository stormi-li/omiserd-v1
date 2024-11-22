package main

import (
	"github.com/go-redis/redis/v8"
	"github.com/stormi-li/omiserd-v1"
)

var redisAddr = "118.25.196.166:3934"
var password = "12982397StrongPassw0rd"

func main() {
	omiC := omiserd.NewClient(&redis.Options{Addr: redisAddr, Password: password}, omiserd.Web)
	// r := omiC.NewRegister("web_service", "localhost:8080")
	// r.RegisterAndServe(1, func(port string) {})
	r := omiC.NewRegister("web_service", "localhost:8081")
	r.RegisterAndServe(1, func(port string) {})
	select {}
}
