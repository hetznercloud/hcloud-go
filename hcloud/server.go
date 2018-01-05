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

// ServerPublicNetIPv4 represents a server's public IPv4 network.
type ServerPublicNetIPv4 struct {
	IP      string
	Blocked bool
	DNSPtr  string
}

// ServerPublicNetIPv6 represents a server's public IPv6 network.
type ServerPublicNetIPv6 struct {
	IP      string
	Blocked bool
	DNSPtr  []ServerPublicNetIPv6DNSPtr
}

// ServerPublicNetIPv6DNSPtr represents a server's public IPv6 reverse DNS.
type ServerPublicNetIPv6DNSPtr struct {
	IP     string
	DNSPtr string
}

// ServerRescueType represents rescue types.
type ServerRescueType string

// List of rescue types.
const (
	ServerRescueTypeLinux32   ServerRescueType = "linux32"
	ServerRescueTypeLinux64                    = "linux64"
	ServerRescueTypeFreeBSD64                  = "freebsd64"
)

// ServerClient is a client for the servers API.
type ServerClient struct {
	client *Client
}

// GetByID retrieves a server by its ID.
func (c *ServerClient) GetByID(ctx context.Context, id int) (*Server, *Response, error) {
	req, err := c.client.NewRequest(ctx, "GET", fmt.Sprintf("/servers/%d", id), nil)
	if err != nil {
		return nil, nil, err
	}

	var body schema.ServerGetResponse
	resp, err := c.client.Do(req, &body)
	if err != nil {
		if IsError(err, ErrorCodeNotFound) {
			return nil, resp, nil
		}
		return nil, nil, err
	}
	return ServerFromSchema(body.Server), resp, nil
}

// GetByName retreives a server by its name.
func (c *ServerClient) GetByName(ctx context.Context, name string) (*Server, *Response, error) {
	path := "/servers?name=" + url.QueryEscape(name)
	req, err := c.client.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	var body schema.ServerListResponse
	resp, err := c.client.Do(req, &body)
	if err != nil {
		return nil, nil, err
	}

	if len(body.Servers) == 0 {
		return nil, resp, nil
	}
	return ServerFromSchema(body.Servers[0]), resp, nil
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
	SSHKeys    []*SSHKey
	Location   *Location
	Datacenter *Datacenter
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
	if o.Location != nil && o.Datacenter != nil {
		return errors.New("location and datacenter are mutually exclusive")
	}
	return nil
}

// ServerCreateResult is the result of a create server call.
type ServerCreateResult struct {
	Server       *Server
	Action       *Action
	RootPassword string
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
	for _, sshKey := range opts.SSHKeys {
		reqBody.SSHKeys = append(reqBody.SSHKeys, sshKey.ID)
	}
	if opts.Location != nil {
		if opts.Location.ID != 0 {
			reqBody.Location = strconv.Itoa(opts.Location.ID)
		} else {
			reqBody.Location = opts.Location.Name
		}
	}
	if opts.Datacenter != nil {
		if opts.Datacenter.ID != 0 {
			reqBody.Datacenter = strconv.Itoa(opts.Datacenter.ID)
		} else {
			reqBody.Datacenter = opts.Datacenter.Name
		}
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
	result := ServerCreateResult{
		Server: ServerFromSchema(respBody.Server),
		Action: ActionFromSchema(respBody.Action),
	}
	if respBody.RootPassword != nil {
		result.RootPassword = *respBody.RootPassword
	}
	return result, resp, nil
}

// Delete deletes a server.
func (c *ServerClient) Delete(ctx context.Context, server *Server) (*Response, error) {
	req, err := c.client.NewRequest(ctx, "DELETE", fmt.Sprintf("/servers/%d", server.ID), nil)
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

// ServerCreateImageOpts specifies options for creating an image from a server.
type ServerCreateImageOpts struct {
	Type        ImageType
	Description *string
}

// Validate checks if options are valid.
func (o ServerCreateImageOpts) Validate() error {
	switch o.Type {
	case ImageTypeSnapshot, ImageTypeBackup:
		break
	case "":
		break
	default:
		return errors.New("invalid type")
	}

	return nil
}

// ServerCreateImageResult is the result of creating an image from a server.
type ServerCreateImageResult struct {
	Action *Action
	Image  *Image
}

// CreateImage creates an image from a server.
func (c *ServerClient) CreateImage(ctx context.Context, server *Server, opts *ServerCreateImageOpts) (ServerCreateImageResult, *Response, error) {
	var reqBody schema.ServerActionCreateImageRequest
	if opts != nil {
		if err := opts.Validate(); err != nil {
			return ServerCreateImageResult{}, nil, fmt.Errorf("invalid options: %s", err)
		}
		if opts.Description != nil {
			reqBody.Description = opts.Description
		}
		if opts.Type != "" {
			reqBody.Type = String(string(opts.Type))
		}
	}
	reqBodyData, err := json.Marshal(reqBody)
	if err != nil {
		return ServerCreateImageResult{}, nil, err
	}

	path := fmt.Sprintf("/servers/%d/actions/create_image", server.ID)
	req, err := c.client.NewRequest(ctx, "POST", path, bytes.NewReader(reqBodyData))
	if err != nil {
		return ServerCreateImageResult{}, nil, err
	}

	respBody := schema.ServerActionCreateImageResponse{}
	resp, err := c.client.Do(req, &respBody)
	if err != nil {
		return ServerCreateImageResult{}, resp, err
	}
	return ServerCreateImageResult{
		Action: ActionFromSchema(respBody.Action),
		Image:  ImageFromSchema(respBody.Image),
	}, resp, nil
}

// ServerEnableRescueOpts specifies options for enabling rescue mode for a server.
type ServerEnableRescueOpts struct {
	Type    ServerRescueType
	SSHKeys []*SSHKey
}

// ServerEnableRescueResult is the result of enabling rescue mode for a server.
type ServerEnableRescueResult struct {
	Action       *Action
	RootPassword string
}

// EnableRescue enables rescue mode for a server.
func (c *ServerClient) EnableRescue(ctx context.Context, server *Server, opts ServerEnableRescueOpts) (ServerEnableRescueResult, *Response, error) {
	reqBody := schema.ServerActionEnableRescueRequest{
		Type: String(string(opts.Type)),
	}
	for _, sshKey := range opts.SSHKeys {
		reqBody.SSHKeys = append(reqBody.SSHKeys, sshKey.ID)
	}
	reqBodyData, err := json.Marshal(reqBody)
	if err != nil {
		return ServerEnableRescueResult{}, nil, err
	}

	path := fmt.Sprintf("/servers/%d/actions/enable_rescue", server.ID)
	req, err := c.client.NewRequest(ctx, "POST", path, bytes.NewReader(reqBodyData))
	if err != nil {
		return ServerEnableRescueResult{}, nil, err
	}

	respBody := schema.ServerActionEnableRescueResponse{}
	resp, err := c.client.Do(req, &respBody)
	if err != nil {
		return ServerEnableRescueResult{}, resp, err
	}
	result := ServerEnableRescueResult{
		Action:       ActionFromSchema(respBody.Action),
		RootPassword: respBody.RootPassword,
	}
	return result, resp, nil
}

// DisableRescue disables rescue mode for a server.
func (c *ServerClient) DisableRescue(ctx context.Context, server *Server) (*Action, *Response, error) {
	path := fmt.Sprintf("/servers/%d/actions/disable_rescue", server.ID)
	req, err := c.client.NewRequest(ctx, "POST", path, nil)
	if err != nil {
		return nil, nil, err
	}

	respBody := schema.ServerActionDisableRescueResponse{}
	resp, err := c.client.Do(req, &respBody)
	if err != nil {
		return nil, resp, err
	}
	return ActionFromSchema(respBody.Action), resp, nil
}
