package omiserd

import (
	"context"
	"strconv"

	"github.com/go-redis/redis/v8"
	omiconst "github.com/stormi-li/omiserd-v1/omiserd_const"
	omiutils "github.com/stormi-li/omiserd-v1/omiserd_utils"
)

// Discover 是服务发现的核心结构
type Discover struct {
	redisClient *redis.Client   // Redis 客户端实例
	prefix      string          // 命名空间前缀
	ctx         context.Context // 上下文，用于 Redis 操作
}

// NewDiscover 创建一个 Discover 实例
// 参数：
// - opts: Redis 连接配置
// - prefix: 命名空间前缀
// 返回值：*Discover
func NewDiscover(opts *redis.Options, prefix string) *Discover {
	return &Discover{
		redisClient: redis.NewClient(opts), // 初始化 Redis 客户端
		prefix:      prefix,                // 设置命名空间前缀
		ctx:         context.Background(),  // 默认上下文
	}
}

// Close 关闭 Redis 客户端
func (discover *Discover) Close() {
	discover.redisClient.Close()
}

// Get 获取指定服务名下的所有实例地址
// 参数：
// - serverName: 服务名
// 返回值：[]string（服务实例的地址列表）
func (discover *Discover) Get(serverName string) []string {
	// 使用命名空间工具函数获取所有与服务名相关的键
	return omiutils.GetKeysByNamespace(discover.redisClient, discover.prefix+serverName)
}

// GetByWeight 根据权重获取服务实例地址池
// 参数：
// - serverName: 服务名
// 返回值：[]string（包含地址权重的地址池）
func (discover *Discover) GetByWeight(serverName string) []string {
	addresses := discover.Get(serverName) // 获取服务实例
	var addressPool []string              // 地址池，用于存放按权重分配的地址

	for _, address := range addresses {
		// 获取实例数据
		data := discover.GetData(serverName, address)

		// 提取权重信息，默认为 1
		weight, err := strconv.Atoi(data["weight"])
		if err != nil || weight <= 0 {
			weight = 1
		}

		// 根据权重将地址加入地址池
		for i := 0; i < weight; i++ {
			addressPool = append(addressPool, address)
		}
	}

	return addressPool
}

// GetData 获取某个服务实例的详细数据
// 参数：
// - serverName: 服务名
// - address: 实例地址
// 返回值：map[string]string（实例数据）
func (discover *Discover) GetData(serverName string, address string) map[string]string {
	// 构造 Redis 键名并从 Redis 中获取值
	key := discover.prefix + serverName + omiconst.Namespace_separator + address
	dataStr, err := discover.redisClient.Get(discover.ctx, key).Result()
	if err != nil {
		return map[string]string{}
	}

	// 将 JSON 字符串转为 map
	data := omiutils.JsonStrToMap(dataStr)
	return data
}

// IsAlive 判断某个实例是否可用
// 参数：
// - serverName: 服务名
// - address: 实例地址
// 返回值：bool
func (discover *Discover) IsAlive(serverName string, address string) bool {
	data := discover.GetData(serverName, address)
	if len(data) == 0 {
		return false
	}
	if data["weight"] == "0" {
		return false
	}
	return true
}

// GetAll 获取所有服务及其实例地址
// 返回值：map[string][]string（服务名 -> 实例地址列表的映射）
func (discover *Discover) GetAll() map[string][]string {
	// 获取与当前命名空间相关的所有键
	keys := omiutils.GetKeysByNamespace(discover.redisClient, discover.prefix[:len(discover.prefix)-1])
	result := map[string][]string{} // 存储结果

	for _, key := range keys {
		// 分割键名为服务名和地址
		name, address := omiutils.SplitMessage(key, omiconst.Namespace_separator)

		// 初始化服务名对应的地址列表（如果尚未存在）
		if _, exists := result[name]; !exists {
			result[name] = []string{}
		}

		// 将地址加入对应服务名的地址列表
		result[name] = append(result[name], address)
	}

	return result
}

func (discover *Discover) NewMonitor(serverName string) *Monitor {
	return NewMonitor(serverName, discover)
}
