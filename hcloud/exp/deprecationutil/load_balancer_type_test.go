package deprecationutil

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestLoadBalancerTypeMessage(t *testing.T) {
	t.Run("not deprecated", func(t *testing.T) {
		o := &hcloud.LoadBalancerType{Name: "lb11"}

		message, isUnavailable := LoadBalancerTypeMessage(o)
		assert.Empty(t, message)
		assert.False(t, isUnavailable)
	})

	t.Run("deprecated", func(t *testing.T) {
		now := time.Now()
		deprecated := now.AddDate(0, 0, -1)
		unavailable := now.AddDate(0, 6, -1)

		o := &hcloud.LoadBalancerType{
			Name: "lb11",
			DeprecatableResource: hcloud.DeprecatableResource{
				Deprecation: &hcloud.DeprecationInfo{
					Announced:        deprecated,
					UnavailableAfter: unavailable,
				},
			},
		}

		message, isUnavailable := LoadBalancerTypeMessage(o)
		assert.Equal(t, fmt.Sprintf(
			`Load Balancer Type "lb11" is deprecated and will no longer be available for order as of %s`,
			unavailable.Format(time.DateOnly),
		), message)
		assert.False(t, isUnavailable)
	})

	t.Run("unavailable", func(t *testing.T) {
		now := time.Now()
		deprecated := now.AddDate(0, -6, -1)
		unavailable := now.AddDate(0, 0, -1)

		o := &hcloud.LoadBalancerType{
			Name: "lb11",
			DeprecatableResource: hcloud.DeprecatableResource{
				Deprecation: &hcloud.DeprecationInfo{
					Announced:        deprecated,
					UnavailableAfter: unavailable,
				},
			},
		}

		message, isUnavailable := LoadBalancerTypeMessage(o)
		assert.Equal(t, `Load Balancer Type "lb11" is unavailable and can no longer be ordered`, message)
		assert.True(t, isUnavailable)
	})
}
