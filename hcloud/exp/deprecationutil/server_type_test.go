package deprecationutil

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestServerTypeMessage(t *testing.T) {
	past := time.Now().UTC().AddDate(0, 0, -14)
	future := time.Now().UTC().AddDate(0, 0, 14)

	Deprecated := hcloud.DeprecatableResource{Deprecation: &hcloud.DeprecationInfo{
		Announced:        past,
		UnavailableAfter: future,
	}}
	Unavailable := hcloud.DeprecatableResource{Deprecation: &hcloud.DeprecationInfo{
		Announced:        past,
		UnavailableAfter: past,
	}}

	t.Run("not deprecated", func(t *testing.T) {
		o := &hcloud.ServerType{Name: "cx22"}

		message, isUnavailable := ServerTypeMessage(o, "")
		assert.Empty(t, message)
		assert.False(t, isUnavailable)
	})

	t.Run("deprecated backward compatible", func(t *testing.T) {
		o := &hcloud.ServerType{
			Name:                 "cx22",
			DeprecatableResource: Deprecated,
		}

		message, isUnavailable := ServerTypeMessage(o, "")
		assert.Equal(t, fmt.Sprintf(`Server Type "cx22" is deprecated in all locations and will no longer be available for order as of %s`, future.Format(time.DateOnly)), message)
		assert.False(t, isUnavailable)
	})

	t.Run("unavailable backward compatible", func(t *testing.T) {
		o := &hcloud.ServerType{
			Name:                 "cx22",
			DeprecatableResource: Unavailable,
		}

		message, isUnavailable := ServerTypeMessage(o, "")
		assert.Equal(t, `Server Type "cx22" is unavailable in all locations and can no longer be ordered`, message)
		assert.True(t, isUnavailable)
	})

	t.Run("not deprecated in any locations with given location", func(t *testing.T) {
		o := &hcloud.ServerType{
			Name: "cx22",
			Locations: []hcloud.ServerTypeLocation{
				{
					Location: &hcloud.Location{Name: "fsn1"},
				},
				{
					Location: &hcloud.Location{Name: "nbg1"},
				},
			},
		}

		message, isUnavailable := ServerTypeMessage(o, "fsn1")
		assert.Empty(t, message)
		assert.False(t, isUnavailable)
	})

	t.Run("deprecated in some locations with given location", func(t *testing.T) {
		o := &hcloud.ServerType{
			Name: "cx22",
			Locations: []hcloud.ServerTypeLocation{
				{
					Location:             &hcloud.Location{Name: "fsn1"},
					DeprecatableResource: Deprecated,
				},
				{
					Location: &hcloud.Location{Name: "nbg1"},
				},
			},
		}

		message, isUnavailable := ServerTypeMessage(o, "fsn1")
		assert.Equal(t, fmt.Sprintf(`Server Type "cx22" is deprecated in "fsn1" and will no longer be available for order as of %s`, future.Format(time.DateOnly)), message)
		assert.False(t, isUnavailable)
	})

	t.Run("unavailable in some locations with given location", func(t *testing.T) {
		o := &hcloud.ServerType{
			Name: "cx22",
			Locations: []hcloud.ServerTypeLocation{
				{
					Location:             &hcloud.Location{Name: "fsn1"},
					DeprecatableResource: Unavailable,
				},
				{
					Location: &hcloud.Location{Name: "nbg1"},
				},
			},
		}

		deprecation, isUnavailable := ServerTypeMessage(o, "fsn1")
		assert.Equal(t, `Server Type "cx22" is unavailable in "fsn1" and can no longer be ordered`, deprecation)
		assert.True(t, isUnavailable)
	})

	t.Run("deprecated in all locations with given location", func(t *testing.T) {
		o := &hcloud.ServerType{
			Name: "cx22",
			Locations: []hcloud.ServerTypeLocation{
				{
					Location:             &hcloud.Location{Name: "fsn1"},
					DeprecatableResource: Deprecated,
				},
				{
					Location:             &hcloud.Location{Name: "nbg1"},
					DeprecatableResource: Deprecated,
				},
			},
		}

		message, isUnavailable := ServerTypeMessage(o, "fsn1")
		assert.Equal(t, fmt.Sprintf(`Server Type "cx22" is deprecated in "fsn1" and will no longer be available for order as of %s`, future.Format(time.DateOnly)), message)
		assert.False(t, isUnavailable)
	})

	t.Run("deprecated in all locations with given unknown location", func(t *testing.T) {
		o := &hcloud.ServerType{
			Name: "cx22",
			Locations: []hcloud.ServerTypeLocation{
				{
					Location:             &hcloud.Location{Name: "fsn1"},
					DeprecatableResource: Deprecated,
				},
				{
					Location:             &hcloud.Location{Name: "nbg1"},
					DeprecatableResource: Deprecated,
				},
			},
		}

		message, isUnavailable := ServerTypeMessage(o, "hel1")
		assert.Empty(t, message)
		assert.False(t, isUnavailable)
	})

	t.Run("deprecated in some locations", func(t *testing.T) {
		o := &hcloud.ServerType{
			Name: "cx22",
			Locations: []hcloud.ServerTypeLocation{
				{
					Location:             &hcloud.Location{Name: "fsn1"},
					DeprecatableResource: Deprecated,
				},
				{
					Location: &hcloud.Location{Name: "nbg1"},
				},
			},
		}

		message, isUnavailable := ServerTypeMessage(o, "")
		assert.Empty(t, message)
		assert.False(t, isUnavailable)
	})

	t.Run("deprecated in all locations", func(t *testing.T) {
		o := &hcloud.ServerType{
			Name: "cx22",
			Locations: []hcloud.ServerTypeLocation{
				{
					Location:             &hcloud.Location{Name: "fsn1"},
					DeprecatableResource: Deprecated,
				},
				{
					Location:             &hcloud.Location{Name: "nbg1"},
					DeprecatableResource: Deprecated,
				},
			},
		}

		message, isUnavailable := ServerTypeMessage(o, "")
		assert.Equal(t, `Server Type "cx22" is deprecated in all locations (fsn1,nbg1) and will no longer be available for order`, message)
		assert.False(t, isUnavailable)
	})

	t.Run("deprecated in all locations and unavailable in some locations", func(t *testing.T) {
		o := &hcloud.ServerType{
			Name: "cx22",
			Locations: []hcloud.ServerTypeLocation{
				{
					Location:             &hcloud.Location{Name: "fsn1"},
					DeprecatableResource: Deprecated,
				},
				{
					Location:             &hcloud.Location{Name: "nbg1"},
					DeprecatableResource: Unavailable,
				},
			},
		}

		message, isUnavailable := ServerTypeMessage(o, "")
		assert.Equal(t, `Server Type "cx22" is deprecated in all locations (fsn1,nbg1) and can no longer be ordered some locations (nbg1)`, message)
		assert.False(t, isUnavailable)
	})

	t.Run("unavailable in all locations", func(t *testing.T) {
		o := &hcloud.ServerType{
			Name: "cx22",
			Locations: []hcloud.ServerTypeLocation{
				{
					Location:             &hcloud.Location{Name: "fsn1"},
					DeprecatableResource: Unavailable,
				},
				{
					Location:             &hcloud.Location{Name: "nbg1"},
					DeprecatableResource: Unavailable,
				},
			},
		}

		message, isUnavailable := ServerTypeMessage(o, "")
		assert.Equal(t, `Server Type "cx22" is unavailable in all locations (fsn1,nbg1) and can no longer be ordered`, message)
		assert.True(t, isUnavailable)
	})
}
