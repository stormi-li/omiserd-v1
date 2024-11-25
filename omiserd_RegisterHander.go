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
		jsonStrData := MapToJsonStr(register.Data)
		key := register.prefix + register.ServerName + Namespace_separator + register.Address
		register.RedisClient.Set(register.ctx, key, jsonStrData, Config_expire_time)
		time.Sleep(Config_expire_time / 2)
	}
}
