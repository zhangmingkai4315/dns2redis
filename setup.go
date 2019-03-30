package dns2redis

import (
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

func setup(c *caddy.Controller) error {
	redisServerAndPort := defaultAddr
	for c.Next() {
		if c.NextArg() {
			redisServerAndPort = c.Val()
		}
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return RedisHandler{Next: next, serverAndHost: redisServerAndPort}
	})

	return nil
}

