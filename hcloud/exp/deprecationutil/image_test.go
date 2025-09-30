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
		assert.Equal(t, "", message)
		assert.False(t, isUnavailable)
	})

	t.Run("deprecated", func(t *testing.T) {
		deprecated := time.Now().UTC().AddDate(0, 0, -1)

		o := &hcloud.Image{Name: "debian-13", Deprecated: deprecated}

		message, isUnavailable := ImageMessage(o)
		assert.Equal(t, fmt.Sprintf(`Image "debian-13" is deprecated and will no longer be available for order as of %s`, deprecated.AddDate(0, 3, 0).Format(time.DateOnly)), message)
		assert.False(t, isUnavailable)
	})

	t.Run("unavailable", func(t *testing.T) {
		deprecated := time.Now().UTC().AddDate(0, -3, -1)

		o := &hcloud.Image{Name: "debian-13", Deprecated: deprecated}

		message, isUnavailable := ImageMessage(o)
		assert.Equal(t, `Image "debian-13" is unavailable and can no longer be ordered`, message)
		assert.True(t, isUnavailable)
	})
}
