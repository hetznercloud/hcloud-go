package serverutils

import (
	"net/netip"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

type IPKind string

const (
	PublicIPv4Kind  IPKind = "public_ipv4"
	PublicIPv6Kind  IPKind = "public_ipv6"
	PrivateIPv4Kind IPKind = "private_ipv4"
)

func FirstAvailableIPByKind(server *hcloud.Server, kind ...IPKind) string {
	mapping := make(map[IPKind]string, 3)

	if !server.PublicNet.IPv4.IsUnspecified() {
		mapping[PublicIPv4Kind] = server.PublicNet.IPv4.IP.String()
	}

	if !server.PublicNet.IPv6.IsUnspecified() {
		publicIPv6Network, ok := netip.AddrFromSlice(server.PublicNet.IPv6.IP)
		if ok {
			mapping[PublicIPv6Kind] = publicIPv6Network.Next().String()
		}
	}

	if len(server.PrivateNet) > 0 {
		mapping[PrivateIPv4Kind] = server.PrivateNet[0].IP.String()
	}

	for _, k := range kind {
		ip, ok := mapping[k]
		if ok {
			return ip
		}
	}

	return ""
}

func FirstAvailableIP(server *hcloud.Server) string {
	return FirstAvailableIPByKind(server, PublicIPv4Kind, PublicIPv6Kind, PrivateIPv4Kind)
}
