package register

import (
	"github.com/go-redis/redis/v8"
	"github.com/stormi-li/omipc-v1"
	omiconst "github.com/stormi-li/omiserd-v1/omiserd_const"
	omiutils "github.com/stormi-li/omiserd-v1/omiserd_utils"
)

type MessageHandler struct {
	ompcClient *omipc.Client
	handlers   []func(command, message string, register *Register)
}

func newMessageHander(opts *redis.Options) *MessageHandler {
	return &MessageHandler{
		ompcClient: omipc.NewClient(opts),
		handlers:   []func(command, message string, register *Register){},
	}
}

func (messageHandler *MessageHandler) AddHandleFunc(handler func(command, message string, register *Register)) {
	messageHandler.handlers = append(messageHandler.handlers, handler)
}

func (messageHandler *MessageHandler) Handle(channel string, register *Register) {
	messageHandler.ompcClient.Listen(channel, 0, func(message string) {
		command, message := omiutils.SplitMessage(message, omiconst.Namespace_separator)
		for _, handler := range messageHandler.handlers {
			handler(command, message, register)
		}
	})
}
