// Code generated by ifacemaker; DO NOT EDIT.

package hcloud

import (
	"context"
)

// IPricingClient ...
type IPricingClient interface {
	// Get retrieves pricing information.
	Get(ctx context.Context) (Pricing, *Response, error)
}
