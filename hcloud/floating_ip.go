package hcloud

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/url"

	"github.com/hetznercloud/hcloud-go/hcloud/schema"
)

// FloatingIP represents a Floating IP in the Hetzner Cloud.
type FloatingIP struct {
	ID           int
	Description  string
	IP           net.IP
	Network      *net.IPNet
	Type         FloatingIPType
	Server       *Server
	DNSPtr       map[string]string
	HomeLocation *Location
	Blocked      bool
	Protection   FloatingIPProtection
	Labels       map[string]string
}

// DNSPtrForIP returns the reverse DNS pointer of the IP address.
func (f *FloatingIP) DNSPtrForIP(ip net.IP) string {
	return f.DNSPtr[ip.String()]
}

// FloatingIPProtection represents the protection level of a Floating IP.
type FloatingIPProtection struct {
	Delete bool
}

// FloatingIPType represents the type of a Floating IP.
type FloatingIPType string

// Floating IP types.
const (
	FloatingIPTypeIPv4 FloatingIPType = "ipv4"
	FloatingIPTypeIPv6 FloatingIPType = "ipv6"
)

// FloatingIPClient is a client for the Floating IP API.
type FloatingIPClient struct {
	client *Client
}

// GetByID retrieves a Floating IP by its ID. If the Floating IP does not exist,
// nil is returned.
func (c *FloatingIPClient) GetByID(ctx context.Context, id int) (*FloatingIP, *Response, error) {
	req, err := c.client.NewRequest(ctx, "GET", fmt.Sprintf("/floating_ips/%d", id), nil)
	if err != nil {
		return nil, nil, err
	}

	var body schema.FloatingIPGetResponse
	resp, err := c.client.Do(req, &body)
	if err != nil {
		if IsError(err, ErrorCodeNotFound) {
			return nil, resp, nil
		}
		return nil, resp, err
	}
	return FloatingIPFromSchema(body.FloatingIP), resp, nil
}

// FloatingIPListOpts specifies options for listing Floating IPs.
type FloatingIPListOpts struct {
	ListOpts
}

func (l *FloatingIPListOpts) values() url.Values {
	vals := valuesForListOpts(l.ListOpts)

	return vals
}

// List returns a list of Floating IPs for a specific page.
func (c *FloatingIPClient) List(ctx context.Context, opts FloatingIPListOpts) ([]*FloatingIP, *Response, error) {
	path := "/floating_ips?" + opts.values().Encode()
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
	return c.AllWithOpts(ctx, FloatingIPListOpts{ListOpts{PerPage: 50}})
}

// AllWithOpts returns all Floating IPs for the given options.
func (c *FloatingIPClient) AllWithOpts(ctx context.Context, opts FloatingIPListOpts) ([]*FloatingIP, error) {
	allFloatingIPs := []*FloatingIP{}

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

// FloatingIPCreateOpts specifies options for creating a Floating IP.
type FloatingIPCreateOpts struct {
	Type         FloatingIPType
	HomeLocation *Location
	Server       *Server
	Description  *string
	Labels       map[string]string
}

// Validate checks if options are valid.
func (o FloatingIPCreateOpts) Validate() error {
	switch o.Type {
	case FloatingIPTypeIPv4, FloatingIPTypeIPv6:
		break
	default:
		return errors.New("missing or invalid type")
	}
	if o.HomeLocation == nil && o.Server == nil {
		return errors.New("one of home location or server is required")
	}
	return nil
}

// FloatingIPCreateResult is the result of creating a Floating IP.
type FloatingIPCreateResult struct {
	FloatingIP *FloatingIP
	Action     *Action
}

// Create creates a Floating IP.
func (c *FloatingIPClient) Create(ctx context.Context, opts FloatingIPCreateOpts) (FloatingIPCreateResult, *Response, error) {
	if err := opts.Validate(); err != nil {
		return FloatingIPCreateResult{}, nil, err
	}

	reqBody := schema.FloatingIPCreateRequest{
		Type:        string(opts.Type),
		Description: opts.Description,
	}
	if opts.HomeLocation != nil {
		reqBody.HomeLocation = String(opts.HomeLocation.Name)
	}
	if opts.Server != nil {
		reqBody.Server = Int(opts.Server.ID)
	}
	if opts.Labels != nil {
		reqBody.Labels = &opts.Labels
	}
	reqBodyData, err := json.Marshal(reqBody)
	if err != nil {
		return FloatingIPCreateResult{}, nil, err
	}

	req, err := c.client.NewRequest(ctx, "POST", "/floating_ips", bytes.NewReader(reqBodyData))
	if err != nil {
		return FloatingIPCreateResult{}, nil, err
	}

	var respBody schema.FloatingIPCreateResponse
	resp, err := c.client.Do(req, &respBody)
	if err != nil {
		return FloatingIPCreateResult{}, resp, err
	}
	var action *Action
	if respBody.Action != nil {
		action = ActionFromSchema(*respBody.Action)
	}
	return FloatingIPCreateResult{
		FloatingIP: FloatingIPFromSchema(respBody.FloatingIP),
		Action:     action,
	}, resp, nil
}

// Delete deletes a Floating IP.
func (c *FloatingIPClient) Delete(ctx context.Context, floatingIP *FloatingIP) (*Response, error) {
	req, err := c.client.NewRequest(ctx, "DELETE", fmt.Sprintf("/floating_ips/%d", floatingIP.ID), nil)
	if err != nil {
		return nil, err
	}
	return c.client.Do(req, nil)
}

// FloatingIPUpdateOpts specifies options for updating a Floating IP.
type FloatingIPUpdateOpts struct {
	Description string
	Labels      map[string]string
}

// Update updates a Floating IP.
func (c *FloatingIPClient) Update(ctx context.Context, floatingIP *FloatingIP, opts FloatingIPUpdateOpts) (*FloatingIP, *Response, error) {
	reqBody := schema.FloatingIPUpdateRequest{
		Description: opts.Description,
	}
	if opts.Labels != nil {
		reqBody.Labels = &opts.Labels
	}
	reqBodyData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, nil, err
	}

	path := fmt.Sprintf("/floating_ips/%d", floatingIP.ID)
	req, err := c.client.NewRequest(ctx, "PUT", path, bytes.NewReader(reqBodyData))
	if err != nil {
		return nil, nil, err
	}

	respBody := schema.FloatingIPUpdateResponse{}
	resp, err := c.client.Do(req, &respBody)
	if err != nil {
		return nil, resp, err
	}
	return FloatingIPFromSchema(respBody.FloatingIP), resp, nil
}

// Assign assigns a Floating IP to a server.
func (c *FloatingIPClient) Assign(ctx context.Context, floatingIP *FloatingIP, server *Server) (*Action, *Response, error) {
	reqBody := schema.FloatingIPActionAssignRequest{
		Server: server.ID,
	}
	reqBodyData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, nil, err
	}

	path := fmt.Sprintf("/floating_ips/%d/actions/assign", floatingIP.ID)
	req, err := c.client.NewRequest(ctx, "POST", path, bytes.NewReader(reqBodyData))
	if err != nil {
		return nil, nil, err
	}

	var respBody schema.FloatingIPActionAssignResponse
	resp, err := c.client.Do(req, &respBody)
	if err != nil {
		return nil, resp, err
	}
	return ActionFromSchema(respBody.Action), resp, nil
}

// Unassign unassigns a Floating IP from the currently assigned server.
func (c *FloatingIPClient) Unassign(ctx context.Context, floatingIP *FloatingIP) (*Action, *Response, error) {
	var reqBody schema.FloatingIPActionUnassignRequest
	reqBodyData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, nil, err
	}

	path := fmt.Sprintf("/floating_ips/%d/actions/unassign", floatingIP.ID)
	req, err := c.client.NewRequest(ctx, "POST", path, bytes.NewReader(reqBodyData))
	if err != nil {
		return nil, nil, err
	}

	var respBody schema.FloatingIPActionUnassignResponse
	resp, err := c.client.Do(req, &respBody)
	if err != nil {
		return nil, resp, err
	}
	return ActionFromSchema(respBody.Action), resp, nil
}

// ChangeDNSPtr changes or resets the reverse DNS pointer for a Floating IP address.
// Pass a nil ptr to reset the reverse DNS pointer to its default value.
func (c *FloatingIPClient) ChangeDNSPtr(ctx context.Context, floatingIP *FloatingIP, ip string, ptr *string) (*Action, *Response, error) {
	reqBody := schema.FloatingIPActionChangeDNSPtrRequest{
		IP:     ip,
		DNSPtr: ptr,
	}
	reqBodyData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, nil, err
	}

	path := fmt.Sprintf("/floating_ips/%d/actions/change_dns_ptr", floatingIP.ID)
	req, err := c.client.NewRequest(ctx, "POST", path, bytes.NewReader(reqBodyData))
	if err != nil {
		return nil, nil, err
	}

	respBody := schema.FloatingIPActionChangeDNSPtrResponse{}
	resp, err := c.client.Do(req, &respBody)
	if err != nil {
		return nil, resp, err
	}
	return ActionFromSchema(respBody.Action), resp, nil
}

// FloatingIPChangeProtectionOpts specifies options for changing the resource protection level of a Floating IP.
type FloatingIPChangeProtectionOpts struct {
	Delete *bool
}

// ChangeProtection changes the resource protection level of a Floating IP.
func (c *FloatingIPClient) ChangeProtection(ctx context.Context, floatingIP *FloatingIP, opts FloatingIPChangeProtectionOpts) (*Action, *Response, error) {
	reqBody := schema.FloatingIPActionChangeProtectionRequest{
		Delete: opts.Delete,
	}
	reqBodyData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, nil, err
	}

	path := fmt.Sprintf("/floating_ips/%d/actions/change_protection", floatingIP.ID)
	req, err := c.client.NewRequest(ctx, "POST", path, bytes.NewReader(reqBodyData))
	if err != nil {
		return nil, nil, err
	}

	respBody := schema.FloatingIPActionChangeProtectionResponse{}
	resp, err := c.client.Do(req, &respBody)
	if err != nil {
		return nil, resp, err
	}
	return ActionFromSchema(respBody.Action), resp, err
}
