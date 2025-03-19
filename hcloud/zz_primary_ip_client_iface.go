// Code generated by ifacemaker; DO NOT EDIT.

package hcloud

import (
	"context"
)

// IPrimaryIPClient ...
type IPrimaryIPClient interface {
	// GetByID retrieves a Primary IP by its ID. If the Primary IP does not exist, nil is returned.
	GetByID(ctx context.Context, id int64) (*PrimaryIP, *Response, error)
	// GetByIP retrieves a Primary IP by its IP Address. If the Primary IP does not exist, nil is returned.
	GetByIP(ctx context.Context, ip string) (*PrimaryIP, *Response, error)
	// GetByName retrieves a Primary IP by its name. If the Primary IP does not exist, nil is returned.
	GetByName(ctx context.Context, name string) (*PrimaryIP, *Response, error)
	// Get retrieves a Primary IP by its ID if the input can be parsed as an integer, otherwise it
	// retrieves a Primary IP by its name. If the Primary IP does not exist, nil is returned.
	Get(ctx context.Context, idOrName string) (*PrimaryIP, *Response, error)
	// List returns a list of Primary IPs for a specific page.
	//
	// Please note that filters specified in opts are not taken into account
	// when their value corresponds to their zero value or when they are empty.
	List(ctx context.Context, opts PrimaryIPListOpts) ([]*PrimaryIP, *Response, error)
	// All returns all Primary IPs.
	All(ctx context.Context) ([]*PrimaryIP, error)
	// AllWithOpts returns all Primary IPs for the given options.
	AllWithOpts(ctx context.Context, opts PrimaryIPListOpts) ([]*PrimaryIP, error)
	// Create creates a Primary IP.
	Create(ctx context.Context, opts PrimaryIPCreateOpts) (*PrimaryIPCreateResult, *Response, error)
	// Delete deletes a Primary IP.
	Delete(ctx context.Context, primaryIP *PrimaryIP) (*Response, error)
	// Update updates a Primary IP.
	Update(ctx context.Context, primaryIP *PrimaryIP, opts PrimaryIPUpdateOpts) (*PrimaryIP, *Response, error)
	// Assign a Primary IP to a resource.
	Assign(ctx context.Context, opts PrimaryIPAssignOpts) (*Action, *Response, error)
	// Unassign a Primary IP from a resource.
	Unassign(ctx context.Context, id int64) (*Action, *Response, error)
	// ChangeDNSPtr Change the reverse DNS from a Primary IP.
	ChangeDNSPtr(ctx context.Context, opts PrimaryIPChangeDNSPtrOpts) (*Action, *Response, error)
	// ChangeProtection Changes the protection configuration of a Primary IP.
	ChangeProtection(ctx context.Context, opts PrimaryIPChangeProtectionOpts) (*Action, *Response, error)
}
