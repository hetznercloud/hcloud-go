package deprecationutil

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestImageMessage(t *testing.T) {
	t.Run("not deprecated", func(t *testing.T) {
		o := &hcloud.Image{Name: "debian-13"}

		message, isUnavailable := ImageMessage(o)
		assert.Empty(t, message)
		assert.False(t, isUnavailable)
	})

	t.Run("deprecated", func(t *testing.T) {
		now := time.Now()
		deprecated := now.AddDate(0, 0, -1)
		unavailable := deprecated.AddDate(0, 3, 0)

		o := &hcloud.Image{Name: "debian-13", DeprecatableResource: hcloud.DeprecatableResource{
			Deprecation: &hcloud.DeprecationInfo{
				Announced:        deprecated,
				UnavailableAfter: unavailable,
			},
		}}

		message, isUnavailable := ImageMessage(o)
		assert.Equal(t, fmt.Sprintf(`Image "debian-13" is deprecated and will no longer be available for order as of %s`, deprecated.AddDate(0, 3, 0).Format(time.DateOnly)), message)
		assert.False(t, isUnavailable)
	})

	t.Run("unavailable", func(t *testing.T) {
		now := time.Now()
		deprecated := now.AddDate(0, -3, -1)
		unavailable := now.AddDate(0, 0, -1)

		o := &hcloud.Image{Name: "debian-13", DeprecatableResource: hcloud.DeprecatableResource{
			Deprecation: &hcloud.DeprecationInfo{
				Announced:        deprecated,
				UnavailableAfter: unavailable,
			},
		}}

		message, isUnavailable := ImageMessage(o)
		assert.Equal(t, `Image "debian-13" is unavailable and can no longer be ordered`, message)
		assert.True(t, isUnavailable)
	})
}
