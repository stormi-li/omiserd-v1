# Omiserd 注册中心框架
**作者**: stormi-li  
**Email**: 2785782829@qq.com  
## 简介
**Omiserd** 是一个基于 Redis 的服务注册与发现框架。它允许用户自定义服务的注册和发现逻辑，并且支持远程控制服务节点的状态，例如调整节点权重，从而实现更灵活的服务管理和灰度发布。
## 功能
- **服务注册与发现**：支持动态注册服务和从注册中心发现服务。
- **自定义负载均衡策略**：允许用户自定义服务的负载均衡策略。
- **自定义注册和发现逻辑**：灵活的 API 支持定制化的注册和发现流程。
- **远程控制**：支持通过 Redis 发布订阅机制，远程控制服务节点（如更新服务权重）。
## 教程
```shell
go get github.com/stormi-li/omiserd-v1
```
### 服务注册
```go
package main

import (
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stormi-li/omiserd-v1"
	omiconst "github.com/stormi-li/omiserd-v1/omiserd_const"
	register "github.com/stormi-li/omiserd-v1/omiserd_register"
)

func main() {
	// 初始化 Omiserd 客户端
	omiC := omiserd.NewClient(&redis.Options{Addr: "localhost:6379"}, omiconst.Web)

	// 创建服务注册实例
	r := omiC.NewRegister("web_service", "localhost:8181")

	// 添加消息处理函数
	r.AddMessageHandleFunc(func(command, message string, register *register.Register) {
		fmt.Println("接收到消息:", command, message)
	})

	// 添加注册逻辑处理函数
	r.AddRegisterHandleFunc(func(register *register.Register) {
		register.Data["time"] = time.Now().String() // 动态添加注册时间
	})

	// 权重为 1 ，开始注册服务并运行
	r.RegisterAndServe(1, func(port string) {
		fmt.Println("服务运行在端口:", port)
	})

	select {} 
}
```
### 远程控制
```go
package main

import (
	"github.com/go-redis/redis/v8"
	"github.com/stormi-li/omiserd-v1"
	omiconst "github.com/stormi-li/omiserd-v1/omiserd_const"
)

func main() {
	// 初始化 Omiserd 客户端
	omiC := omiserd.NewClient(&redis.Options{Addr: "localhost:6379"}, omiconst.Web)

	// 创建服务注册实例
	r := omiC.NewRegister("web_service", "localhost:8081")

	// 发送权重更新指令
	r.SendMessage(omiconst.Command_update_weight, "3")
}
```

### 服务发现
```go
package main

import (
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stormi-li/omiserd-v1"
	omiconst "github.com/stormi-li/omiserd-v1/omiserd_const"
	discover "github.com/stormi-li/omiserd-v1/omiserd_discover"
)

func main() {
	// 初始化 Omiserd 客户端
	omiserdC := omiserd.NewClient(&redis.Options{Addr: "localhost:6379"}, omiconst.Web)

	// 创建服务发现实例
	d := omiserdC.NewDiscover()

	// 创建服务监控实例，监听 web_service
	monitor := d.NewMonitor("web_service")

	// 自定义监听处理逻辑
	listenHandleFunc := func(serverName, oldAddress string, discover *discover.Discover) string {
		if !discover.IsAlive(serverName, oldAddress) { // 检查旧实例是否存活
			addresses := discover.GetByWeight(serverName) // 根据权重获取可用地址
			if len(addresses) > 0 {
				return addresses[rand.IntN(len(addresses))] // 随机选择一个地址
			}
		}
		return ""
	}

	// 自定义连接处理逻辑
	connectHandleFunc := func(address string) {
		fmt.Println("连接到服务节点:", address)
	}

	// 开始监听服务并建立连接
	monitor.ListenAndConnect(2*time.Second, listenHandleFunc, connectHandleFunc)
}
```