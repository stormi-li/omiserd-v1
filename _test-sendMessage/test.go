package main

import (
	"github.com/go-redis/redis/v8"
	"github.com/stormi-li/omiserd-v1"
	omiconst "github.com/stormi-li/omiserd-v1/omiserd_const"
)

var redisAddr = "118.25.196.166:3934"
var password = "12982397StrongPassw0rd"

func main() {
	omiC := omiserd.NewClient(&redis.Options{Addr: redisAddr, Password: password}, omiconst.Web)
	r := omiC.NewRegister("web_service", "localhost:8081")
	r.SendMessage(omiconst.Command_update_weight, "3")
}
