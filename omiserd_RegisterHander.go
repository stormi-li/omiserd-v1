package omiserd

import (
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stormi-li/omipc-v1"
)

type RegisterHandler struct {
	ompcClient  *omipc.Client
	handleFuncs map[string]func() string
}

func newRegisterHandler(opts *redis.Options) *RegisterHandler {
	return &RegisterHandler{
		ompcClient:  omipc.NewClient(opts),
		handleFuncs: map[string]func() string{},
	}
}

func (registerHandler *RegisterHandler) AddHandleFunc(key string, handleFunc func() string) {
	registerHandler.handleFuncs[key] = handleFunc
}

func (registerHandler *RegisterHandler) Handle(register *Register) {
	for {
		for key, handleFunc := range registerHandler.handleFuncs {
			register.Data[key] = handleFunc()
		}
		jsonStrData := mapToJsonStr(register.Data)
		key := register.prefix + register.serverName + namespace_separator + register.address
		register.redisClient.Set(register.ctx, key, jsonStrData, config_expire_time)
		time.Sleep(config_expire_time / 2)
	}
}
