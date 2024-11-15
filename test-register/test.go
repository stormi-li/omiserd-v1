package main

import (
	"fmt"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/stormi-li/omiserd-v1"
)

var redisAddr = "118.25.196.166:3934"
var password = "12982397StrongPassw0rd"

func main() {
	omiC := omiserd.NewClient(&redis.Options{Addr: redisAddr, Password: password}, omiserd.Config)
	r := omiC.NewRegister("mysql", "118.25.196.166:3933")
	r.Data["username"] = "root"
	r.Data["database"] = "USER"
	r.Data["password"] = "12982397StrongPassw0rd"
	r.Register(1)
	r = omiC.NewRegister("redis", "118.25.196.166:6379")
	count := 0
	r.SetRegisterHandler(func(register *omiserd.Register) {
		register.Data["count"] = strconv.Itoa(count)
		count++
	})
	r.SetMessageHandler(func(command, message string, register *omiserd.Register) {
		fmt.Println(command, message)
	})
	r.Register(1)
	select {}
}
