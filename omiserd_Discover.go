package omiserd

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"math/rand"

	"github.com/go-redis/redis/v8"
)

type Discover struct {
	redisClient       *redis.Client
	namespace         string
	ctx               context.Context
	reconnectStrategy func(address string, data map[string]string, err error) bool
	discoverStrategy  func(serverName string, discover *Discover) (string, map[string]string)
}

func (discover *Discover) Close() {
	discover.redisClient.Close()
}

func (discover *Discover) DiscoverByName(serverName string) map[string]map[string]string {
	keys := getKeysByNamespace(discover.redisClient, discover.namespace+serverName)
	res := map[string]map[string]string{}
	for _, key := range keys {
		data, _ := discover.redisClient.Get(discover.ctx, discover.namespace+serverName+namespace_separator+key).Result()
		res[key] = jsonStrToMap(data)
	}
	return res
}

func (discover *Discover) IsAlive(serverName string, address string) bool {
	data := discover.GetData(serverName, address)
	if data == nil || data["weight"] == "0" || data["weight"] == "" {
		return false
	}
	return true
}

func (discover *Discover) GetData(serverName string, address string) map[string]string {
	dataStr, _ := discover.redisClient.Get(discover.ctx, discover.namespace+serverName+namespace_separator+address).Result()
	if dataStr == "" {
		return map[string]string{}
	}
	return jsonStrToMap(dataStr)
}

func (discover *Discover) DiscoverAllServers() map[string]map[string]map[string]string {
	keys := getKeysByNamespace(discover.redisClient, discover.namespace[:len(discover.namespace)-1])
	res := map[string]map[string]map[string]string{}
	for _, key := range keys {
		data, _ := discover.redisClient.Get(discover.ctx, discover.namespace+key).Result()
		parts := split(key)
		if res[parts[0]] == nil {
			res[parts[0]] = map[string]map[string]string{}
		}
		res[parts[0]][parts[1]] = jsonStrToMap(data)
	}
	return res
}

func (discover *Discover) DiscoverByWeight(serverName string) (string, map[string]string) {
	addrs := discover.DiscoverByName(serverName)
	var addressPool []string
	var dataPool []map[string]string
	for name, data := range addrs {
		weight, _ := strconv.Atoi(data["weight"])
		for i := 0; i < weight; i++ {
			addressPool = append(addressPool, name)
			dataPool = append(dataPool, data)
		}
	}
	if len(addressPool) == 0 {
		return "", map[string]string{}
	}
	selectIndex := rand.Intn(len(addressPool))
	return addressPool[selectIndex], dataPool[selectIndex]
}

func (discover *Discover) SetReconnectStrategy(strategy func(address string, data map[string]string, err error) bool) {
	discover.reconnectStrategy = strategy
}

func (discover *Discover) SetDiscoverStrategy(strategy func(serverName string, discover *Discover) (string, map[string]string)) {
	discover.discoverStrategy = strategy
}

func (discover *Discover) DiscoverAndListen(serverName string, connectHandler func(address string, data map[string]string)) {
	address := ""
	for {
		data := discover.GetData(serverName, address)
		var err error
		if len(data) == 0 {
			err = fmt.Errorf("no service found")
		}
		if discover.reconnectStrategy(address, data, err) {
			address, data = discover.discoverStrategy(serverName, discover)
			if address != "" {
				go func() {
					defer func() {
						if r := recover(); r != nil {
							log.Println(r)
						}
					}()
					connectHandler(address, data)
				}()
			}
		}
		time.Sleep(config_expire_time / 2)
	}
}

func split(address string) []string {
	index := strings.Index(address, namespace_separator)
	if index == -1 {
		return nil
	}
	return []string{address[:index], address[index+1:]}
}
