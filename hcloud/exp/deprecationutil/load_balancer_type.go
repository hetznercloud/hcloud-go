package deprecationutil

import (
	"fmt"
	"time"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// LoadBalancerTypeMessage return a deprecation message when the given Load Balancer Type is
// deprecated and whether the given Load Balancer Type is unavailable.
//
// Experimental: `exp` package is experimental, breaking changes may occur within minor releases.
func LoadBalancerTypeMessage(lbType *hcloud.LoadBalancerType) (string, bool) {
	if lbType.IsDeprecated() {
		if time.Now().After(lbType.UnavailableAfter()) {
			return fmt.Sprintf(
				"Load Balancer Type %q is unavailable and can no longer be ordered",
				lbType.Name,
			), true
		}
		return fmt.Sprintf(
			"Load Balancer Type %q is deprecated and will no longer be available for order as of %s",
			lbType.Name,
			lbType.UnavailableAfter().Format(time.DateOnly),
		), false
	}

	return "", false
}
