package deprecationutil

import (
	"fmt"
	"time"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// ImageMessage return a deprecation message when the given Image is
// deprecated and whether the given Image is unavailable.
//
// Experimental: `exp` package is experimental, breaking changes may occur within minor releases.
func ImageMessage(image *hcloud.Image) (string, bool) {
	if image.IsDeprecated() {
		if time.Now().After(image.UnavailableAfter()) {
			return fmt.Sprintf(
				"Image %q is unavailable and can no longer be ordered",
				image.Name,
			), true
		}
		return fmt.Sprintf(
			"Image %q is deprecated and will no longer be available for order as of %s",
			image.Name,
			image.UnavailableAfter().Format(time.DateOnly),
		), false
	}

	return "", false
}
