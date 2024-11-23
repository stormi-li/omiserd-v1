package main

import (
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stormi-li/omiserd-v1"
	omiconst "github.com/stormi-li/omiserd-v1/omiserd_const"
	register "github.com/stormi-li/omiserd-v1/omiserd_register"
)

var redisAddr = "118.25.196.166:3934"
var password = "12982397StrongPassw0rd"

func main() {
	omiC := omiserd.NewClient(&redis.Options{Addr: redisAddr, Password: password}, omiconst.Web)

	r := omiC.NewRegister("web_service", "localhost:8181")

	r.AddMessageHandleFunc(func(command, message string, register *register.Register) {
		fmt.Println(command, message)
	})

	r.AddRegisterHandleFunc(func(register *register.Register) {
		register.Data["time"] = time.Now().String()
	})

	r.RegisterAndServe(1, func(port string) {})
	select {}
}
