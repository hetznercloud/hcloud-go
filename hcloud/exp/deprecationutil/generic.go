package deprecationutil

import (
	"fmt"
	"time"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// Message return a deprecation message when the given deprecatable resource is
// deprecated and whether the given resource is unavailable.
//
// Experimental: `exp` package is experimental, breaking changes may occur within minor releases.
func Message[T hcloud.Deprecatable](o T, resource string, identifier string) (string, bool) {
	if o.IsDeprecated() {
		if time.Now().After(o.UnavailableAfter()) {
			return fmt.Sprintf(
				"%s %q is unavailable and can no longer be ordered",
				resource,
				identifier,
			), true
		}
		return fmt.Sprintf(
			"%s %q is deprecated and will no longer be available for order as of %s",
			resource,
			identifier,
			o.UnavailableAfter().Format(time.DateOnly),
		), false
	}

	return "", false
}
