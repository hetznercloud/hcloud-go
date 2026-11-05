package deprecationutil

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// LoadBalancerTypeMessage return a deprecation message when the given Load Balancer Type is
// deprecated and whether the given Load Balancer Type is unavailable.
//
// Experimental: `exp` package is experimental, breaking changes may occur within minor releases.
func LoadBalancerTypeMessage(loadBalancerType *hcloud.LoadBalancerType) (string, bool) {
	return Message(loadBalancerType, "Load Balancer Type", loadBalancerType.Name)
}
