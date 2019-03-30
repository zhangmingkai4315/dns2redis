package dns2redis

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

// RedisHandler handle redis connection
// and save status
type RedisHandler struct {
	serverAndHost string
	client        *redis.Client
	Next          plugin.Handler
}

// ServeDNS implements the plugin.Handler interface.
func (rh RedisHandler) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}
	ip := state.IP()
	queryName := state.QName()
	err := rh.client.Set(queryName, ip, time.Duration(time.Minute*5)).Err()
	if err != nil {
		log.Printf("error save key: %s", err)
	}
	rh.client.IncrBy("queryCounter", 1)
	return rh.Next.ServeDNS(ctx, w, r)
}

// Name implements the Handler interface.
func (rh RedisHandler) Name() string { return "dns2redis" }
