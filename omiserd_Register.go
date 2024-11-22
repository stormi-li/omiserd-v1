package omiserd

import (
	"context"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stormi-li/omipc-v1"
)

type Register struct {
	redisClient     *redis.Client
	ServerName      string
	Address         string
	Weight          int
	Data            map[string]string
	namespace       string
	omipcClient     *omipc.Client
	ctx             context.Context
	registerHandler func(register *Register)
	messageHandler  func(command, message string, register *Register)
	close           chan struct{}
}

func (register *Register) Close() {
	if register.close != nil {
		register.close <- struct{}{}
		time.Sleep(100 * time.Millisecond)
		<-register.close
	}

	register.redisClient.Close()
	register.omipcClient.Close()
	log.Println("register server for", register.ServerName+"["+register.Address+"]", "is closed")
}

func (register *Register) SetRegisterHandler(handler func(register *Register)) {
	register.registerHandler = handler
}

func (register *Register) SetMessageHandler(handler func(command, message string, register *Register)) {
	register.messageHandler = handler
}

func (register *Register) Register(weight int) {
	register.Weight = weight
	log.Println("register server for", register.ServerName+"["+register.Address+"]", "is starting")
	go register.registerHandle()
	register.messageHandle()
}

func (register *Register) RegisterAndListen(weight int, handler func(port string)) {
	register.Register(weight)
	for {
		handler(":" + strings.Split(register.Address, ":")[1])
	}
}

func (register *Register) registerHandle() {
	for {
		register.Data["weight"] = strconv.Itoa(register.Weight)
		register.Data["process_id"] = strconv.Itoa(os.Getpid())
		register.Data["host"], _ = os.Hostname()
		register.registerHandler(register)
		jsonStrData := mapToJsonStr(register.Data)
		key := register.namespace + register.ServerName + namespace_separator + register.Address
		register.redisClient.Set(register.ctx, key, jsonStrData, config_expire_time)
		time.Sleep(config_expire_time / 2)
	}
}

func (register *Register) messageHandle() {
	channel := register.namespace + register.ServerName + namespace_separator + register.Address
	register.close = register.omipcClient.Listen(channel, func(message string) bool {
		command, message := splitMessage(message, namespace_separator)
		if command == Command_update_weight {
			weight, err := strconv.Atoi(message)
			if err == nil {
				register.Weight = weight
			}
		}
		register.messageHandler(command, message, register)
		return true
	})
}

func (register *Register) SendMessage(command string, message string) {
	key := register.namespace + register.ServerName + namespace_separator + register.Address
	register.omipcClient.Notify(key, command+namespace_separator+message)
}
