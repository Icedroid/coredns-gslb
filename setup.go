package gslb

import (
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"

	"github.com/mholt/caddy"
)

var ipToZoneURL string

func init() {
	caddy.RegisterPlugin("gslb", caddy.Plugin{
		ServerType: "dns",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	c.Next() // 'Gslb'
	if c.NextArg() {
		ipToZoneURL = c.Val()
		if ipToZoneURL == "" {
			return plugin.Error("gslb must have a request host", c.ArgErr())
		}
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return Gslb{}
	})

	return nil
}
