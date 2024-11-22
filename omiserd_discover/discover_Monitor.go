package discover

import "time"

type Monitor struct {
	serverName string
	discover   *Discover
}

func NewMonitor(serverName string, d *Discover) *Monitor {
	return &Monitor{
		serverName: serverName,
		discover:   d,
	}
}

// Listen 监听服务变化并触发连接
// - interval: 每次检查间隔
// - listenFunc: 监听逻辑，返回新的地址（或空字符串表示无变化）
// - connectFunc: 连接逻辑，处理新的地址
func (monitor *Monitor) ListenAndConnect(interval time.Duration, listenFunc func(serverName, oldAddress string, discover *Discover) string, connectFunc func(address string)) {
	oldAddress := ""
	for {
		address := listenFunc(monitor.serverName, oldAddress, monitor.discover)
		if address != "" {
			connectFunc(address)
			oldAddress = address
		}
		time.Sleep(interval)
	}
}
