package hcloud

import (
	"context"
	"fmt"

	"github.com/hetznercloud/hcloud-go/hcloud/schema"
)

// FloatingIP represents a Floating IP in the Hetzner Cloud.
type FloatingIP struct {
	ID           int
	Description  string
	IP           string
	Type         FloatingIPType
	Server       *Server
	DNSPtr       map[string]string
	HomeLocation *Location
}

// FloatingIPType represents the type of a Floating IP.
type FloatingIPType string

// Floating IP types.
const (
	FloatingIPTypeIPv4 FloatingIPType = "ipv4"
	FloatingIPTypeIPv6                = "ipv6"
)

// FloatingIPFromSchema converts a schema.FloatingIP to a FloatingIP.
func FloatingIPFromSchema(s schema.FloatingIP) *FloatingIP {
	f := &FloatingIP{
		ID:           s.ID,
		IP:           s.IP,
		Type:         FloatingIPType(s.Type),
		HomeLocation: LocationFromSchema(s.HomeLocation),
	}
	if s.Description != nil {
		f.Description = *s.Description
	}
	if s.Server != nil {
		f.Server = &Server{ID: *s.Server}
	}
	if s.DNSPtr.IPv4 != nil {
		f.DNSPtr = map[string]string{
			s.IP: *s.DNSPtr.IPv4,
		}
	} else if s.DNSPtr.IPv6 != nil {
		f.DNSPtr = map[string]string{}
		for _, entry := range s.DNSPtr.IPv6 {
			f.DNSPtr[entry.IP] = entry.DNSPtr
		}
	}
	return f
}

// FloatingIPClient is a client for the Floating IP API.
type FloatingIPClient struct {
	client *Client
}

// Get retrieves a Floating IP.
func (c *FloatingIPClient) Get(ctx context.Context, id int) (*FloatingIP, *Response, error) {
	req, err := c.client.NewRequest(ctx, "GET", fmt.Sprintf("/floating_ips/%d", id), nil)
	if err != nil {
		return nil, nil, err
	}

	var body schema.FloatingIPGetResponse
	resp, err := c.client.Do(req, &body)
	if err != nil {
		return nil, resp, err
	}
	return FloatingIPFromSchema(body.FloatingIP), resp, nil
}

// FloatingIPListOpts specifies options for listing Floating IPs.
type FloatingIPListOpts struct {
	ListOpts
}

// List returns a list of Floating IPs for a specific page.
func (c *FloatingIPClient) List(ctx context.Context, opts FloatingIPListOpts) ([]*FloatingIP, *Response, error) {
	path := "/floating_ips?" + valuesForListOpts(opts.ListOpts).Encode()
	req, err := c.client.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	var body schema.FloatingIPListResponse
	resp, err := c.client.Do(req, &body)
	if err != nil {
		return nil, nil, err
	}
	floatingIPs := make([]*FloatingIP, 0, len(body.FloatingIPs))
	for _, s := range body.FloatingIPs {
		floatingIPs = append(floatingIPs, FloatingIPFromSchema(s))
	}
	return floatingIPs, resp, nil
}

// All returns all Floating IPs.
func (c *FloatingIPClient) All(ctx context.Context) ([]*FloatingIP, error) {
	allFloatingIPs := []*FloatingIP{}

	opts := FloatingIPListOpts{}
	opts.PerPage = 50

	_, err := c.client.all(func(page int) (*Response, error) {
		opts.Page = page
		floatingIPs, resp, err := c.List(ctx, opts)
		if err != nil {
			return resp, err
		}
		allFloatingIPs = append(allFloatingIPs, floatingIPs...)
		return resp, nil
	})
	if err != nil {
		return nil, err
	}

	return allFloatingIPs, nil
}
