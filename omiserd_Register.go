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

// Register 是服务注册和消息处理的核心结构
type Register struct {
	RedisClient     *redis.Client     // Redis 客户端实例
	ServerName      string            // 服务名
	Address         string            // 服务地址（包含主机和端口）
	Weight          int               // 服务权重
	Data            map[string]string // 服务的元数据，如权重、主机名等
	NodeType        NodeType
	prefix          string           // 命名空间前缀
	channel         string           // Redis 发布/订阅使用的频道名
	omipcClient     *omipc.Client    // omipc 客户端，用于异步通信
	ctx             context.Context  // 上下文，用于 Redis 操作
	registerHandler *RegisterHandler // 注册处理器，管理服务注册逻辑
	messageHandler  *MessageHandler  // 消息处理器，处理接收到的消息
	startTime       time.Time
}

// NewRegister 创建一个新的 Register 实例
// 参数：
// - opts: Redis 连接配置
// - serverName: 服务名称
// - address: 服务地址（格式为 "host:port"）
// - prefix: 命名空间前缀
// 返回值：*Register
func NewRegister(opts *redis.Options, serverName, address string, prefix string, nodeType NodeType) *Register {
	register := &Register{
		RedisClient:     redis.NewClient(opts), // 初始化 Redis 客户端
		ServerName:      serverName,
		Address:         address,
		Data:            map[string]string{}, // 初始化空元数据
		prefix:          prefix,
		NodeType:        nodeType,
		ctx:             context.Background(),                                // 默认上下文
		omipcClient:     omipc.NewClient(opts),                               // 创建 omipc 客户端
		registerHandler: newRegisterHandler(opts),                            // 创建服务注册处理器
		messageHandler:  newMessageHander(opts),                              // 创建消息处理器
		channel:         prefix + serverName + namespace_separator + address, // 频道名称由前缀、服务名和地址拼接而成
		startTime:       time.Now(),
	}

	// 添加默认的注册逻辑处理函数
	register.AddRegisterHandleFunc("weight", func() string {
		return strconv.Itoa(register.Weight)
	})
	register.AddRegisterHandleFunc("process_id", func() string {
		return strconv.Itoa(os.Getpid())
	})
	register.AddRegisterHandleFunc("host", func() string {
		host, _ := os.Hostname()
		return host
	})
	register.AddRegisterHandleFunc("start_time", func() string {
		return register.startTime.Format("2006-01-02 15:04:05")
	})
	register.AddRegisterHandleFunc("run_time", func() string {
		return time.Since(register.startTime).String()
	})

	// 添加消息权重修改回调函数
	register.AddMessageHandleFunc(Command_update_weight, func(message string) {
		register.AddRegisterHandleFunc("weight", func() string {
			return message
		})
	})

	return register
}

// AddRegisterHandleFunc 添加额外的注册处理函数
func (register *Register) AddRegisterHandleFunc(key string, handleFunc func() string) {
	register.registerHandler.AddHandleFunc(key, handleFunc)
}

// AddMessageHandleFunc 添加额外的消息处理函数
func (register *Register) AddMessageHandleFunc(command string, handleFunc func(message string)) {
	register.messageHandler.AddHandleFunc(command, handleFunc)
}

// RegisterAndServe 启动服务注册并运行服务
// 参数：
// - weight: 服务权重
// - serverFunc: 服务的启动函数，通常是一个 HTTP 或 TCP 服务器
func (register *Register) RegisterAndServe(weight int, serverFunc func(port string)) {
	register.Weight = weight
	log.Printf("%s register server for %s[%s] is starting", string(register.NodeType), register.ServerName, register.Address)

	// 启动服务注册逻辑和消息处理逻辑
	go register.registerHandler.Handle(register)
	go register.messageHandler.Handle(register.channel)

	// 提取端口号并调用服务启动函数
	if parts := strings.Split(register.Address, ":"); len(parts) == 2 {
		serverFunc(":" + parts[1])
	} else {
		log.Fatalf("invalid address format: %s", register.Address)
	}
}

// SendMessage 发送消息到指定频道
// 参数：
// - command: 消息命令
// - message: 消息内容
func (register *Register) SendMessage(command string, message string) {
	register.omipcClient.Notify(register.channel, command+namespace_separator+message)
}
