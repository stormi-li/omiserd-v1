package omiserd

import "time"

type Monitor struct {
	ServerName string
	Address    string
	Discover   *Discover
}

func NewMonitor(serverName string, d *Discover) *Monitor {
	return &Monitor{
		ServerName: serverName,
		Discover:   d,
	}
}

// Listen 监听服务变化并触发连接
// - interval: 每次检查间隔
// - listenFunc: 监听逻辑，返回新的地址（或空字符串表示无变化）
// - connectFunc: 连接逻辑，处理新的地址
func (monitor *Monitor) ListenAndConnect(interval time.Duration, listenFunc func(serverName, oldAddress string, discover *Discover) string, connectFunc func(address string)) {
	for {
		address := listenFunc(monitor.ServerName, monitor.Address, monitor.Discover)
		if address != "" {
			connectFunc(address)
			monitor.Address = address
		}
		time.Sleep(interval)
	}
}
