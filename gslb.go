// Package gslb implements a plugin that returns details about the resolving
// querying it.
package gslb

import (
	"context"
	"net"
	"strconv"

	clog "github.com/coredns/coredns/plugin/pkg/log"
	"github.com/coredns/coredns/request"

	"github.com/miekg/dns"
)

// Define log to be a logger with the plugin name in it. This way we can just use log.Info and
// friends to log.
var log = clog.NewWithPlugin("gslb")

// Gslb is a plugin that returns your IP address, port and the protocol used for connecting
// to CoreDNS.
type Gslb struct{}

// ServeDNS implements the plugin.Handler interface.
func (wh Gslb) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}

	a := new(dns.Msg)
	a.SetReply(r)
	a.Authoritative = true

	ip := getRealRequestIP(state)
	var rr dns.RR
	var rrs dns.RR

	url := getIPNetURL(ipToZoneURL, ip)
	log.Infof("request url=%s", url)

	log.Infof("state IP=%s, QName=%s, QClass=%d, QType=%d, Proto=%s, Family=%d", ip, state.QName(), state.QClass(),
		state.QType(), state.Proto(), state.Family())

	switch state.Family() {
	case 1:
		rr = new(dns.A)
		rr.(*dns.A).Hdr = dns.RR_Header{Name: state.QName(), Rrtype: dns.TypeA, Class: state.QClass()}
		rr.(*dns.A).A = net.ParseIP(ip).To4()

		zoneIP, err := requestURL(url)
		if err != nil {
			log.Errorf("request url for fecth zone ip error=%s", err)
		}
		rrs = new(dns.A)
		rrs.(*dns.A).Hdr = dns.RR_Header{Name: state.QName(), Rrtype: dns.TypeA, Class: state.QClass()}
		rrs.(*dns.A).A = net.ParseIP(zoneIP).To4()
	case 2:
		rr = new(dns.AAAA)
		rr.(*dns.AAAA).Hdr = dns.RR_Header{Name: state.QName(), Rrtype: dns.TypeAAAA, Class: state.QClass()}
		rr.(*dns.AAAA).AAAA = net.ParseIP(ip)

		rrs = new(dns.AAAA)
		rrs.(*dns.AAAA).Hdr = dns.RR_Header{Name: state.QName(), Rrtype: dns.TypeAAAA, Class: state.QClass()}
		rrs.(*dns.AAAA).AAAA = net.ParseIP("fe80::b2d2:77b5:ced3:4e40").To16()
	}

	srv := new(dns.SRV)
	srv.Hdr = dns.RR_Header{Name: "_" + state.Proto() + "." + state.QName(), Rrtype: dns.TypeSRV, Class: state.QClass()}
	if state.QName() == "." {
		srv.Hdr.Name = "_" + state.Proto() + state.QName()
	}
	port, _ := strconv.Atoi(state.Port())
	srv.Port = uint16(port)
	srv.Target = "."

	a.Answer = []dns.RR{rrs}
	a.Extra = []dns.RR{rr, srv}

	w.WriteMsg(a)

	return 0, nil
}

// Name implements the Handler interface.
func (wh Gslb) Name() string { return "Gslb" }


func getRealRequestIP(req request.Request) string {
	r := req.Req

	if opt := r.IsEdns0(); opt != nil {
		// for i := len(r.Extra) - 1; i >= 0; i-- {
		// 	log.Infof("%d.opt=%+v", i, r.Extra[i])
		// 	if r.Extra[i].Header().Rrtype == dns.TypeOPT {
		// 		log.Infof("%d.opt=%s", i, r.Extra[i].(*dns.OPT).String())
		// 	}
		// }
		for _, o := range opt.Option {
			switch v := o.(type) {
			case *dns.EDNS0_SUBNET:
				switch v.Family {
				case 1:
					return v.Address.To4().String()
				case 2:
					return v.Address.String()
				}
			}
		}
		log.Infof("return opt=%v", opt)
		log.Infof("return opt Header=%v", opt.Hdr.String())
	}

	return req.IP()
}