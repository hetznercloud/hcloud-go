package deprecationutil

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestServerTypeWarning(t *testing.T) {
	day := 24 * time.Hour

	past := time.Now().Add(-14 * day).UTC()
	future := time.Now().Add(14 * day).UTC()

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

		warn, warnIsError := ServerTypeWarning(o, "")
		assert.Equal(t, "", warn)
		assert.False(t, warnIsError)
	})

	t.Run("deprecated backward compatible", func(t *testing.T) {
		o := &hcloud.ServerType{
			Name:                 "cx22",
			DeprecatableResource: Deprecated,
		}

		warn, warnIsError := ServerTypeWarning(o, "")
		assert.Equal(t, fmt.Sprintf(`Server Type "cx22" is deprecated in all locations and will no longer be available for order as of %s.`, future.Format(time.DateOnly)), warn)
		assert.False(t, warnIsError)
	})

	t.Run("unavailable backward compatible", func(t *testing.T) {
		o := &hcloud.ServerType{
			Name:                 "cx22",
			DeprecatableResource: Unavailable,
		}

		warn, warnIsError := ServerTypeWarning(o, "")
		assert.Equal(t, `Server Type "cx22" is unavailable in all locations and can no longer be ordered.`, warn)
		assert.True(t, warnIsError)
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

		warn, warnIsError := ServerTypeWarning(o, "fsn1")
		assert.Equal(t, ``, warn)
		assert.False(t, warnIsError)
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

		warn, warnIsError := ServerTypeWarning(o, "fsn1")
		assert.Equal(t, fmt.Sprintf(`Server Type "cx22" is deprecated in "fsn1" and will no longer be available for order as of %s.`, future.Format(time.DateOnly)), warn)
		assert.False(t, warnIsError)
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

		warn, warnIsError := ServerTypeWarning(o, "fsn1")
		assert.Equal(t, `Server Type "cx22" is unavailable in "fsn1" and can no longer be ordered.`, warn)
		assert.True(t, warnIsError)
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

		warn, warnIsError := ServerTypeWarning(o, "fsn1")
		assert.Equal(t, fmt.Sprintf(`Server Type "cx22" is deprecated in "fsn1" and will no longer be available for order as of %s.`, future.Format(time.DateOnly)), warn)
		assert.False(t, warnIsError)
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

		warn, warnIsError := ServerTypeWarning(o, "hel1")
		assert.Equal(t, "", warn)
		assert.False(t, warnIsError)
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

		warn, warnIsError := ServerTypeWarning(o, "")
		assert.Equal(t, "", warn)
		assert.False(t, warnIsError)
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

		warn, warnIsError := ServerTypeWarning(o, "")
		assert.Equal(t, `Server Type "cx22" is deprecated in all locations (fsn1,nbg1) and will no longer be available for order.`, warn)
		assert.False(t, warnIsError)
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

		warn, warnIsError := ServerTypeWarning(o, "")
		assert.Equal(t, `Server Type "cx22" is deprecated in all locations (fsn1,nbg1) and can no longer be ordered some locations (nbg1).`, warn)
		assert.False(t, warnIsError)
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

		warn, warnIsError := ServerTypeWarning(o, "")
		assert.Equal(t, `Server Type "cx22" is unavailable in all locations (fsn1,nbg1) and can no longer be ordered.`, warn)
		assert.True(t, warnIsError)
	})
}
