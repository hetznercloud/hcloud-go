package hcloud

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// Server represents a server in the Hetzner Cloud.
type Server struct {
	ID              int
	Name            string
	Status          ServerStatus
	Created         time.Time
	PublicNet       ServerPublicNet
	ServerType      ServerType
	IncludedTraffic uint64
	OutgoingTraffic uint64
	IngoingTraffic  uint64
	BackupWindow    string
	RescueEnabled   bool
	ISO             *ISO
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *Server) UnmarshalJSON(data []byte) error {
	var v struct {
		ID              int             `json:"id"`
		Name            string          `json:"name"`
		Status          string          `json:"status"`
		Created         time.Time       `json:"created"`
		PublicNet       ServerPublicNet `json:"public_net"`
		ServerType      ServerType      `json:"server_type"`
		IncludedTraffic uint64          `json:"included_traffic"`
		OutgoingTraffic uint64          `json:"outgoing_traffic"`
		IngoingTraffic  uint64          `json:"ingoing_traffic"`
		BackupWindow    string          `json:"backup_window"`
		RescueEnabled   bool            `json:"rescue_enabled"`
		ISO             *ISO            `json:"iso"`
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	s.ID = v.ID
	s.Name = v.Name
	s.Status = ServerStatus(v.Status)
	s.Created = v.Created
	s.PublicNet = v.PublicNet
	s.ServerType = v.ServerType
	s.IncludedTraffic = v.IncludedTraffic
	s.OutgoingTraffic = v.OutgoingTraffic
	s.IngoingTraffic = v.IngoingTraffic
	s.BackupWindow = v.BackupWindow
	s.RescueEnabled = v.RescueEnabled
	s.ISO = v.ISO

	return nil
}

// ServerStatus specifies a server's status.
type ServerStatus string

const (
	ServerStatusInitializing ServerStatus = "initializing" // Server is initializing
	ServerStatusOff                       = "off"          // Server is off
	ServerStatusRunning                   = "running"      // Server is running
)

// ServerPublicNet represents a server's public network.
type ServerPublicNet struct {
	IPv4        ServerPublicNetIPv4
	IPv6        ServerPublicNetIPv6
	FloatingIPs []*FloatingIP
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *ServerPublicNet) UnmarshalJSON(data []byte) error {
	var v struct {
		IPv4        ServerPublicNetIPv4 `json:"ipv4"`
		IPv6        ServerPublicNetIPv6 `json:"ipv6"`
		FloatingIPs []int               `json:"floating_ips"`
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	s.IPv4 = v.IPv4
	s.IPv6 = v.IPv6

	for _, f := range v.FloatingIPs {
		s.FloatingIPs = append(s.FloatingIPs, &FloatingIP{ID: f})
	}

	return nil
}

// ServerPublicNetIPv4 represents a server's public IPv4 network.
type ServerPublicNetIPv4 struct {
	IP      string
	Blocked bool
	DNSPtr  string
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *ServerPublicNetIPv4) UnmarshalJSON(data []byte) error {
	var v struct {
		IP      string `json:"ip"`
		Blocked bool   `json:"blocked"`
		DNSPtr  string `json:"dns_ptr"`
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	s.IP = v.IP
	s.Blocked = v.Blocked
	s.DNSPtr = v.DNSPtr

	return nil
}

// ServerPublicNetIPv6 represents a server's public IPv6 network.
type ServerPublicNetIPv6 struct {
	IP      string
	Blocked bool
	DNSPtr  []ServerPublicNetIPv6DNSPtr
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *ServerPublicNetIPv6) UnmarshalJSON(data []byte) error {
	var v struct {
		IP      string                      `json:"ip"`
		Blocked bool                        `json:"blocked"`
		DNSPtr  []ServerPublicNetIPv6DNSPtr `json:"dns_ptr"`
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	s.IP = v.IP
	s.Blocked = v.Blocked
	s.DNSPtr = v.DNSPtr

	return nil
}

// ServerPublicNetIPv6DNSPtr represents a server's public IPv6 reverse DNS.
type ServerPublicNetIPv6DNSPtr struct {
	IP     string
	DNSPtr string
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *ServerPublicNetIPv6DNSPtr) UnmarshalJSON(data []byte) error {
	var v struct {
		IP     string `json:"ip"`
		DNSPtr string `json:"dns_ptr"`
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	s.IP = v.IP
	s.DNSPtr = v.DNSPtr

	return nil
}

// ServerClient is a client for the servers API.
type ServerClient struct {
	client *Client
}

// Get retrieves a server.
func (c *ServerClient) Get(ctx context.Context, id int) (*Server, *Response, error) {
	req, err := c.client.NewRequest(ctx, "GET", fmt.Sprintf("/servers/%d", id), nil)
	if err != nil {
		return nil, nil, err
	}

	var body struct {
		Server *Server `json:"server"`
	}
	resp, err := c.client.Do(req, &body)
	if err != nil {
		return nil, nil, err
	}
	return body.Server, resp, nil
}

// ServerListOpts specifies options for listing servers.
type ServerListOpts struct {
	ListOpts
}

// List returns a list of servers for a specific page.
func (c *ServerClient) List(ctx context.Context, opts ServerListOpts) ([]*Server, *Response, error) {
	path := "/servers"
	vals := url.Values{}
	if opts.Page > 0 {
		vals.Add("page", strconv.Itoa(opts.Page))
	}
	if opts.PerPage > 0 {
		vals.Add("per_page", strconv.Itoa(opts.PerPage))
	}
	path += "?" + vals.Encode()

	req, err := c.client.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	var body struct {
		Servers []*Server `json:"servers"`
	}
	resp, err := c.client.Do(req, &body)
	if err != nil {
		return nil, nil, err
	}
	return body.Servers, resp, nil
}

// ListAll returns all servers by going through all pages.
func (c *ServerClient) ListAll(ctx context.Context) ([]*Server, error) {
	allServers := []*Server{}

	opts := ServerListOpts{}
	opts.PerPage = 50

	_, err := c.client.all(func(page int) (*Response, error) {
		opts.Page = page
		servers, resp, err := c.List(ctx, opts)
		if err != nil {
			return resp, err
		}
		allServers = append(allServers, servers...)
		return resp, nil
	})
	if err != nil {
		return nil, err
	}

	return allServers, nil
}

// ServerCreateOpts specifies options for creating a new server.
type ServerCreateOpts struct {
	Name       string
	ServerType ServerType
	Image      Image
}

// Validate checks if options are valid.
func (o ServerCreateOpts) Validate() error {
	if o.Name == "" {
		return errors.New("missing name")
	}
	if o.ServerType.ID == 0 && o.ServerType.Name == "" {
		return errors.New("missing server type")
	}
	if o.Image.ID == 0 && o.Image.Name == "" {
		return errors.New("missing image")
	}
	return nil
}

// ServerCreateResult is the result of a create server call.
type ServerCreateResult struct {
	Server *Server
	Action *Action
}

// Create creates a new server.
func (c *ServerClient) Create(ctx context.Context, opts ServerCreateOpts) (ServerCreateResult, *Response, error) {
	if err := opts.Validate(); err != nil {
		return ServerCreateResult{}, nil, err
	}

	var reqBody struct {
		Name       string      `json:"name"`
		ServerType interface{} `json:"server_type"`
		Image      interface{} `json:"image"`
	}
	reqBody.Name = opts.Name
	if opts.ServerType.ID != 0 {
		reqBody.ServerType = opts.ServerType.ID
	} else if opts.ServerType.Name != "" {
		reqBody.ServerType = opts.ServerType.Name
	}
	if opts.Image.ID != 0 {
		reqBody.Image = opts.Image.ID
	} else if opts.Image.Name != "" {
		reqBody.Image = opts.Image.Name
	}
	reqBodyData, err := json.Marshal(reqBody)
	if err != nil {
		return ServerCreateResult{}, nil, err
	}

	req, err := c.client.NewRequest(ctx, "POST", "/servers", bytes.NewReader(reqBodyData))
	if err != nil {
		return ServerCreateResult{}, nil, err
	}

	var respBody struct {
		Server *Server `json:"server"`
		Action *Action `json:"action"`
	}
	resp, err := c.client.Do(req, &respBody)
	if err != nil {
		return ServerCreateResult{}, resp, err
	}
	result := ServerCreateResult{
		Server: respBody.Server,
		Action: respBody.Action,
	}
	return result, resp, nil
}

// Delete deletes a server.
func (c *ServerClient) Delete(ctx context.Context, id int) (*Response, error) {
	req, err := c.client.NewRequest(ctx, "DELETE", fmt.Sprintf("/servers/%d", id), nil)
	if err != nil {
		return nil, err
	}
	return c.client.Do(req, nil)
}
