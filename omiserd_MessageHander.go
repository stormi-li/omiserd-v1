package omiserd

import (
	"github.com/go-redis/redis/v8"
	"github.com/stormi-li/omipc-v1"
)

type MessageHandler struct {
	ompcClient  *omipc.Client
	handleFuncs map[string]func(message string)
}

func newMessageHander(opts *redis.Options) *MessageHandler {
	return &MessageHandler{
		ompcClient:  omipc.NewClient(opts),
		handleFuncs: map[string]func(message string){},
	}
}

func (messageHandler *MessageHandler) AddHandleFunc(command string, handleFunc func(message string)) {
	messageHandler.handleFuncs[command] = handleFunc
}

func (messageHandler *MessageHandler) Handle(channel string) {
	messageHandler.ompcClient.Listen(channel, 0, func(message string) {
		command, message := SplitMessage(message, Namespace_separator)
		if handleFunc, ok := messageHandler.handleFuncs[command]; ok {
			handleFunc(message)
		}
	})
}
