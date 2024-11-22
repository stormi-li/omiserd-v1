package register

import (
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stormi-li/omipc-v1"
	omiconst "github.com/stormi-li/omiserd-v1/omiserd_const"
	omiutils "github.com/stormi-li/omiserd-v1/omiserd_utils"
)

type RegisterHandler struct {
	ompcClient *omipc.Client
	Handlers   []func(register *Register)
}

func newRegisterHandler(opts *redis.Options) *RegisterHandler {
	return &RegisterHandler{
		ompcClient: omipc.NewClient(opts),
		Handlers:   []func(register *Register){},
	}
}

func (registerHandler *RegisterHandler) AddHandleFunc(handler func(register *Register)) {
	registerHandler.Handlers = append(registerHandler.Handlers, handler)
}

func (registerHandler *RegisterHandler) Handle(register *Register) {
	for {
		for _, handler := range registerHandler.Handlers {
			handler(register)
		}
		jsonStrData := omiutils.MapToJsonStr(register.Data)
		key := register.prefix + register.ServerName + omiconst.Namespace_separator + register.Address
		register.redisClient.Set(register.ctx, key, jsonStrData, omiconst.Config_expire_time)
		time.Sleep(omiconst.Config_expire_time / 2)
	}
}
