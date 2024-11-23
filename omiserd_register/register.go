package register

import (
	"context"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/stormi-li/omipc-v1"
	omiconst "github.com/stormi-li/omiserd-v1/omiserd_const"
)

// Register 是服务注册和消息处理的核心结构
type Register struct {
	RedisClient     *redis.Client     // Redis 客户端实例
	ServerName      string            // 服务名
	Address         string            // 服务地址（包含主机和端口）
	Weight          int               // 服务权重
	Data            map[string]string // 服务的元数据，如权重、主机名等
	prefix          string            // 命名空间前缀
	channel         string            // Redis 发布/订阅使用的频道名
	omipcClient     *omipc.Client     // omipc 客户端，用于异步通信
	ctx             context.Context   // 上下文，用于 Redis 操作
	registerHandler *RegisterHandler  // 注册处理器，管理服务注册逻辑
	messageHandler  *MessageHandler   // 消息处理器，处理接收到的消息
}

// NewRegister 创建一个新的 Register 实例
// 参数：
// - opts: Redis 连接配置
// - serverName: 服务名称
// - address: 服务地址（格式为 "host:port"）
// - prefix: 命名空间前缀
// 返回值：*Register
func NewRegister(opts *redis.Options, serverName, address string, prefix string) *Register {
	register := &Register{
		RedisClient:     redis.NewClient(opts), // 初始化 Redis 客户端
		ServerName:      serverName,
		Address:         address,
		Data:            map[string]string{}, // 初始化空元数据
		prefix:          prefix,
		ctx:             context.Background(),     // 默认上下文
		omipcClient:     omipc.NewClient(opts),    // 创建 omipc 客户端
		registerHandler: newRegisterHandler(opts), // 创建服务注册处理器
		messageHandler:  newMessageHander(opts),   // 创建消息处理器
		// 频道名称由前缀、服务名和地址拼接而成
		channel: prefix + serverName + omiconst.Namespace_separator + address,
	}

	// 添加默认的注册逻辑处理函数
	register.registerHandler.AddHandleFunc(func(register *Register) {
		// 保存服务的权重、进程 ID 和主机名到元数据中
		register.Data["weight"] = strconv.Itoa(register.Weight)
		register.Data["process_id"] = strconv.Itoa(os.Getpid())
		if host, err := os.Hostname(); err == nil {
			register.Data["host"] = host
		} else {
			register.Data["host"] = ""
		}
	})

	// 添加默认的消息处理逻辑
	register.messageHandler.AddHandleFunc(func(command, message string, register *Register) {
		// 如果接收到更新权重的命令，则尝试更新权重
		if command == omiconst.Command_update_weight {
			weight, err := strconv.Atoi(message)
			if err != nil {
				log.Printf("failed to parse weight: %v", err)
			} else {
				register.Weight = weight
			}
		}
	})

	return register
}

// AddRegisterHandleFunc 添加额外的注册处理函数
func (register *Register) AddRegisterHandleFunc(handleFunc func(register *Register)) {
	register.registerHandler.AddHandleFunc(handleFunc)
}

// AddMessageHandleFunc 添加额外的消息处理函数
func (register *Register) AddMessageHandleFunc(handleFunc func(command, message string, register *Register)) {
	register.messageHandler.AddHandleFunc(handleFunc)
}

// RegisterAndServe 启动服务注册并运行服务
// 参数：
// - weight: 服务权重
// - serverFunc: 服务的启动函数，通常是一个 HTTP 或 TCP 服务器
func (register *Register) RegisterAndServe(weight int, serverFunc func(port string)) {
	register.Weight = weight
	log.Printf("register server for %s[%s] is starting", register.ServerName, register.Address)

	// 启动服务注册逻辑和消息处理逻辑
	go register.registerHandler.Handle(register)
	go register.messageHandler.Handle(register.channel, register)

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
	register.omipcClient.Notify(register.channel, command+omiconst.Namespace_separator+message)
}
