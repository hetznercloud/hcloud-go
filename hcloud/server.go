package hcloud

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/hetznercloud/hcloud-go/hcloud/schema"
)

// Server represents a server in the Hetzner Cloud.
type Server struct {
	ID              int
	Name            string
	Status          ServerStatus
	Created         time.Time
	PublicNet       ServerPublicNet
	ServerType      *ServerType
	IncludedTraffic uint64
	OutgoingTraffic uint64
	IngoingTraffic  uint64
	BackupWindow    string
	RescueEnabled   bool
	ISO             *ISO
}

// ServerFromSchema converts a schema.Server to a Server.
func ServerFromSchema(s schema.Server) *Server {
	server := &Server{
		ID:              s.ID,
		Name:            s.Name,
		Status:          ServerStatus(s.Status),
		Created:         s.Created,
		PublicNet:       ServerPublicNetFromSchema(s.PublicNet),
		ServerType:      ServerTypeFromSchema(s.ServerType),
		IncludedTraffic: s.IncludedTraffic,
		OutgoingTraffic: s.OutgoingTraffic,
		IngoingTraffic:  s.IngoingTraffic,
		BackupWindow:    s.BackupWindow,
		RescueEnabled:   s.RescueEnabled,
	}
	if s.ISO != nil {
		server.ISO = ISOFromSchema(*s.ISO)
	}
	return server
}

// ServerStatus specifies a server's status.
type ServerStatus string

const (
	// ServerStatusInitializing is the status when a server is initializing.
	ServerStatusInitializing ServerStatus = "initializing"

	// ServerStatusOff is the status when a server is off.
	ServerStatusOff = "off"

	// ServerStatusRunning is the status when a server is running.
	ServerStatusRunning = "running"
)

// ServerPublicNet represents a server's public network.
type ServerPublicNet struct {
	IPv4        ServerPublicNetIPv4
	IPv6        ServerPublicNetIPv6
	FloatingIPs []*FloatingIP
}

// ServerPublicNetFromSchema converts a schema.ServerPublicNet to a ServerPublicNet.
func ServerPublicNetFromSchema(s schema.ServerPublicNet) ServerPublicNet {
	publicNet := ServerPublicNet{
		IPv4: ServerPublicNetIPv4FromSchema(s.IPv4),
		IPv6: ServerPublicNetIPv6FromSchema(s.IPv6),
	}
	for _, id := range s.FloatingIPs {
		publicNet.FloatingIPs = append(publicNet.FloatingIPs, &FloatingIP{ID: id})
	}
	return publicNet
}

// ServerPublicNetIPv4 represents a server's public IPv4 network.
type ServerPublicNetIPv4 struct {
	IP      string
	Blocked bool
	DNSPtr  string
}

// ServerPublicNetIPv4FromSchema converts a schema.ServerPublicNetIPv4 to
// a ServerPublicNetIPv4.
func ServerPublicNetIPv4FromSchema(s schema.ServerPublicNetIPv4) ServerPublicNetIPv4 {
	return ServerPublicNetIPv4{
		IP:      s.IP,
		Blocked: s.Blocked,
		DNSPtr:  s.DNSPtr,
	}
}

// ServerPublicNetIPv6 represents a server's public IPv6 network.
type ServerPublicNetIPv6 struct {
	IP      string
	Blocked bool
	DNSPtr  []ServerPublicNetIPv6DNSPtr
}

// ServerPublicNetIPv6FromSchema converts a schema.ServerPublicNetIPv6 to
// a ServerPublicNetIPv6.
func ServerPublicNetIPv6FromSchema(s schema.ServerPublicNetIPv6) ServerPublicNetIPv6 {
	ipv6 := ServerPublicNetIPv6{
		IP:      s.IP,
		Blocked: s.Blocked,
	}
	for _, dnsPtr := range s.DNSPtr {
		ipv6.DNSPtr = append(ipv6.DNSPtr, ServerPublicNetIPv6DNSPtrFromSchema(dnsPtr))
	}
	return ipv6
}

// ServerPublicNetIPv6DNSPtr represents a server's public IPv6 reverse DNS.
type ServerPublicNetIPv6DNSPtr struct {
	IP     string
	DNSPtr string
}

// ServerPublicNetIPv6DNSPtrFromSchema converts a schema.ServerPublicNetIPv6DNSPtr
// to a ServerPublicNetIPv6DNSPtr.
func ServerPublicNetIPv6DNSPtrFromSchema(s schema.ServerPublicNetIPv6DNSPtr) ServerPublicNetIPv6DNSPtr {
	return ServerPublicNetIPv6DNSPtr{
		IP:     s.IP,
		DNSPtr: s.DNSPtr,
	}
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

	var body schema.ServerGetResponse
	resp, err := c.client.Do(req, &body)
	if err != nil {
		return nil, nil, err
	}
	return ServerFromSchema(body.Server), resp, nil
}

// ServerListOpts specifies options for listing servers.
type ServerListOpts struct {
	ListOpts
}

// List returns a list of servers for a specific page.
func (c *ServerClient) List(ctx context.Context, opts ServerListOpts) ([]*Server, *Response, error) {
	path := "/servers?" + valuesForListOpts(opts.ListOpts).Encode()
	req, err := c.client.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	var body schema.ServerListResponse
	resp, err := c.client.Do(req, &body)
	if err != nil {
		return nil, nil, err
	}
	servers := make([]*Server, 0, len(body.Servers))
	for _, s := range body.Servers {
		servers = append(servers, ServerFromSchema(s))
	}
	return servers, resp, nil
}

// All returns all servers.
func (c *ServerClient) All(ctx context.Context) ([]*Server, error) {
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

	var reqBody schema.ServerCreateRequest
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

	var respBody schema.ServerCreateResponse
	resp, err := c.client.Do(req, &respBody)
	if err != nil {
		return ServerCreateResult{}, resp, err
	}
	return ServerCreateResult{
		Server: ServerFromSchema(respBody.Server),
		Action: ActionFromSchema(respBody.Action),
	}, resp, nil
}

// Delete deletes a server.
func (c *ServerClient) Delete(ctx context.Context, id int) (*Response, error) {
	req, err := c.client.NewRequest(ctx, "DELETE", fmt.Sprintf("/servers/%d", id), nil)
	if err != nil {
		return nil, err
	}
	return c.client.Do(req, nil)
}

// Poweron starts a server.
func (c *ServerClient) Poweron(ctx context.Context, server *Server) (*Action, *Response, error) {
	path := fmt.Sprintf("/servers/%d/actions/poweron", server.ID)
	req, err := c.client.NewRequest(ctx, "POST", path, nil)
	if err != nil {
		return nil, nil, err
	}

	respBody := schema.ServerActionPoweronResponse{}
	resp, err := c.client.Do(req, &respBody)
	if err != nil {
		return nil, resp, err
	}
	return ActionFromSchema(respBody.Action), resp, nil
}

// Reboot reboots a server.
func (c *ServerClient) Reboot(ctx context.Context, server *Server) (*Action, *Response, error) {
	path := fmt.Sprintf("/servers/%d/actions/reboot", server.ID)
	req, err := c.client.NewRequest(ctx, "POST", path, nil)
	if err != nil {
		return nil, nil, err
	}

	respBody := schema.ServerActionRebootResponse{}
	resp, err := c.client.Do(req, &respBody)
	if err != nil {
		return nil, resp, err
	}
	return ActionFromSchema(respBody.Action), resp, nil
}

// Reset resets a server.
func (c *ServerClient) Reset(ctx context.Context, server *Server) (*Action, *Response, error) {
	path := fmt.Sprintf("/servers/%d/actions/reset", server.ID)
	req, err := c.client.NewRequest(ctx, "POST", path, nil)
	if err != nil {
		return nil, nil, err
	}

	respBody := schema.ServerActionResetResponse{}
	resp, err := c.client.Do(req, &respBody)
	if err != nil {
		return nil, resp, err
	}
	return ActionFromSchema(respBody.Action), resp, nil
}

// Shutdown shuts down a server.
func (c *ServerClient) Shutdown(ctx context.Context, server *Server) (*Action, *Response, error) {
	path := fmt.Sprintf("/servers/%d/actions/shutdown", server.ID)
	req, err := c.client.NewRequest(ctx, "POST", path, nil)
	if err != nil {
		return nil, nil, err
	}

	respBody := schema.ServerActionShutdownResponse{}
	resp, err := c.client.Do(req, &respBody)
	if err != nil {
		return nil, resp, err
	}
	return ActionFromSchema(respBody.Action), resp, nil
}

// Poweroff stops a server.
func (c *ServerClient) Poweroff(ctx context.Context, server *Server) (*Action, *Response, error) {
	path := fmt.Sprintf("/servers/%d/actions/poweroff", server.ID)
	req, err := c.client.NewRequest(ctx, "POST", path, nil)
	if err != nil {
		return nil, nil, err
	}

	respBody := schema.ServerActionPoweroffResponse{}
	resp, err := c.client.Do(req, &respBody)
	if err != nil {
		return nil, resp, err
	}
	return ActionFromSchema(respBody.Action), resp, nil
}

// ServerResetPasswordResult is the result of resetting a server's password.
type ServerResetPasswordResult struct {
	Action       *Action
	RootPassword string
}

// ResetPassword resets a server's password.
func (c *ServerClient) ResetPassword(ctx context.Context, server *Server) (ServerResetPasswordResult, *Response, error) {
	path := fmt.Sprintf("/servers/%d/actions/reset_password", server.ID)
	req, err := c.client.NewRequest(ctx, "POST", path, nil)
	if err != nil {
		return ServerResetPasswordResult{}, nil, err
	}

	respBody := schema.ServerActionResetPasswordResponse{}
	resp, err := c.client.Do(req, &respBody)
	if err != nil {
		return ServerResetPasswordResult{}, resp, err
	}
	return ServerResetPasswordResult{
		Action:       ActionFromSchema(respBody.Action),
		RootPassword: respBody.RootPassword,
	}, resp, nil
}
