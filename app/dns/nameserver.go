package dns

import (
	"context"

	"v2ray.com/core/common/net"
)

type NameServerInterface interface {
	QueryIP(ctx context.Context, domain string) ([]net.IP, error)
}

type localNameServer struct {
	resolver net.Resolver
}

func (s *localNameServer) QueryIP(ctx context.Context, domain string) ([]net.IP, error) {
	ipAddr, err := s.resolver.LookupIPAddr(ctx, domain)
	if err != nil {
		return nil, err
	}
	var ips []net.IP
	for _, addr := range ipAddr {
		ips = append(ips, addr.IP)
	}
	return ips, nil
}

func NewLocalNameServer() *localNameServer {
	return &localNameServer{
		resolver: net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, server string) (net.Conn, error) {
				// Calling Dial here is scary -- we have to be sure not to
				// dial a name that will require a DNS lookup, or Dial will
				// call back here to translate it. The DNS config parser has
				// already checked that all the cfg.servers are IP
				// addresses, which Dial will use without a DNS lookup.
				var c net.Conn
				var err error
				var d net.Dialer
				c, err = d.DialContext(ctx, "tcp", server)
				if err != nil {
					return nil, err
				}
				return c, nil
			},
		},
	}
}
