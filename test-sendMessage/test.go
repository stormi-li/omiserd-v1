package main

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/stormi-li/omiserd-v1"
)

var redisAddr = "118.25.196.166:3934"
var password = "12982397StrongPassw0rd"

func main() {
	omiC := omiserd.NewClient(&redis.Options{Addr: redisAddr, Password: password}, omiserd.Config)
	r := omiC.NewRegister("redis", "118.25.196.166:6379")
	r.SendMessage("cc", "1:1")
	r1 := omiC.NewRegister("redis1", "118.25.196.166:6378")
	r1.SetMessageHandler(func(command, message string, register *omiserd.Register) {
		fmt.Println(command, message)
	})
	r1.Register(1)
	select {}
}
