package dns2redis

import (
	"fmt"
	"net"
	"strings"

	"github.com/go-redis/redis"

	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"

	"github.com/mholt/caddy"
)

const defaultAddr = "localhost:6379"

func init() {
	caddy.RegisterPlugin("dns2redis", caddy.Plugin{
		ServerType: "dns",
		Action:     setup,
	})
}

func isValidateIPAndPort(serverAndHost string) bool {
	info := strings.Split(serverAndHost, ":")
	testInput := net.ParseIP(info[0])
	if info[0] != "localhost" && testInput.To4() == nil {
		return false
	}
	return true
}

func setup(c *caddy.Controller) error {
	redisServerAndPort := defaultAddr
	for c.Next() {
		if c.NextArg() {
			redisServerAndPort = c.Val()
		}
	}
	isValid := isValidateIPAndPort(redisServerAndPort)
	if isValid == false {
		return plugin.Error("dns2redis", fmt.Errorf("invalid ipaddress: %s", redisServerAndPort))
	}
	client := redis.NewClient(&redis.Options{
		Addr:     redisServerAndPort,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	_, err := client.Ping().Result()
	if err != nil {
		return plugin.Error("dns2redis", fmt.Errorf("connect server: %s error:%s", redisServerAndPort, err))
	}
	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return RedisHandler{
			Next:          next,
			client:        client,
			serverAndHost: redisServerAndPort}
	})

	return nil
}
