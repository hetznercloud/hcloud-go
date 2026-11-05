package deprecationutil

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// ImageMessage return a deprecation message when the given Image is
// deprecated and whether the given Image is unavailable.
//
// Experimental: `exp` package is experimental, breaking changes may occur within minor releases.
func ImageMessage(image *hcloud.Image) (string, bool) {
	return Message(image, "Image", image.Name)
}
