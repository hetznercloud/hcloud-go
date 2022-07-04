package hcloud

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
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
	PrivateNet      []ServerPrivateNet
	ServerType      *ServerType
	Datacenter      *Datacenter
	IncludedTraffic uint64
	OutgoingTraffic uint64
	IngoingTraffic  uint64
	BackupWindow    string
	RescueEnabled   bool
	Locked          bool
	ISO             *ISO
	Image           *Image
	Protection      ServerProtection
	Labels          map[string]string
	Volumes         []*Volume
	PrimaryDiskSize int
	PlacementGroup  *PlacementGroup
}

// ServerProtection represents the protection level of a server.
type ServerProtection struct {
	Delete, Rebuild bool
}

// ServerStatus specifies a server's status.
type ServerStatus string

const (
	// ServerStatusInitializing is the status when a server is initializing.
	ServerStatusInitializing ServerStatus = "initializing"

	// ServerStatusOff is the status when a server is off.
	ServerStatusOff ServerStatus = "off"

	// ServerStatusRunning is the status when a server is running.
	ServerStatusRunning ServerStatus = "running"

	// ServerStatusStarting is the status when a server is being started.
	ServerStatusStarting ServerStatus = "starting"

	// ServerStatusStopping is the status when a server is being stopped.
	ServerStatusStopping ServerStatus = "stopping"

	// ServerStatusMigrating is the status when a server is being migrated.
	ServerStatusMigrating ServerStatus = "migrating"

	// ServerStatusRebuilding is the status when a server is being rebuilt.
	ServerStatusRebuilding ServerStatus = "rebuilding"

	// ServerStatusDeleting is the status when a server is being deleted.
	ServerStatusDeleting ServerStatus = "deleting"

	// ServerStatusUnknown is the status when a server's state is unknown.
	ServerStatusUnknown ServerStatus = "unknown"
)

// FirewallStatus specifies a Firewall's status.
type FirewallStatus string

const (
	// FirewallStatusPending is the status when a Firewall is pending.
	FirewallStatusPending FirewallStatus = "pending"

	// FirewallStatusApplied is the status when a Firewall is applied.
	FirewallStatusApplied FirewallStatus = "applied"
)

// ServerPublicNet represents a server's public network.
type ServerPublicNet struct {
	IPv4        ServerPublicNetIPv4
	IPv6        ServerPublicNetIPv6
	FloatingIPs []*FloatingIP
	Firewalls   []*ServerFirewallStatus
}

// ServerPublicNetIPv4 represents a server's public IPv4 address.
type ServerPublicNetIPv4 struct {
	ID      int
	IP      net.IP
	Blocked bool
	DNSPtr  string
}

func (n *ServerPublicNetIPv4) IsUnspecified() bool {
	return n.IP == nil || n.IP.Equal(net.IPv4zero)
}

// ServerPublicNetIPv6 represents a Server's public IPv6 network and address.
type ServerPublicNetIPv6 struct {
	ID      int
	IP      net.IP
	Network *net.IPNet
	Blocked bool
	DNSPtr  map[string]string
}

func (n *ServerPublicNetIPv6) IsUnspecified() bool {
	return n.IP == nil || n.IP.Equal(net.IPv6unspecified)
}

// ServerPrivateNet defines the schema of a Server's private network information.
type ServerPrivateNet struct {
	Network    *Network
	IP         net.IP
	Aliases    []net.IP
	MACAddress string
}

// DNSPtrForIP returns the reverse dns pointer of the ip address.
func (n *ServerPublicNetIPv6) DNSPtrForIP(ip net.IP) string {
	return n.DNSPtr[ip.String()]
}

// ServerFirewallStatus represents a Firewall and its status on a Server's
// network interface.
type ServerFirewallStatus struct {
	Firewall Firewall
	Status   FirewallStatus
}

// ServerRescueType represents rescue types.
type ServerRescueType string

// List of rescue types.
const (
	ServerRescueTypeLinux32 ServerRescueType = "linux32"
	ServerRescueTypeLinux64 ServerRescueType = "linux64"
)

// changeDNSPtr changes or resets the reverse DNS pointer for a IP address.
// Pass a nil ptr to reset the reverse DNS pointer to its default value.
func (s *Server) changeDNSPtr(ctx context.Context, client *Client, ip net.IP, ptr *string) (*Action, *Response, error) {
	reqBody := schema.ServerActionChangeDNSPtrRequest{
		IP:     ip.String(),
		DNSPtr: ptr,
	}
	reqBodyData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, nil, err
	}

	path := fmt.Sprintf("/servers/%d/actions/change_dns_ptr", s.ID)
	req, err := client.NewRequest(ctx, "POST", path, bytes.NewReader(reqBodyData))
	if err != nil {
		return nil, nil, err
	}

	respBody := schema.ServerActionChangeDNSPtrResponse{}
	resp, err := client.Do(req, &respBody)
	if err != nil {
		return nil, resp, err
	}
	return ActionFromSchema(respBody.Action), resp, nil
}

// GetDNSPtrForIP searches for the dns assigned to the given IP address.
// It returns an error if there is no dns set for the given IP address.
func (s *Server) GetDNSPtrForIP(ip net.IP) (string, error) {
	if net.IP.Equal(s.PublicNet.IPv4.IP, ip) {
		return s.PublicNet.IPv4.DNSPtr, nil
	} else if dns, ok := s.PublicNet.IPv6.DNSPtr[ip.String()]; ok {
		return dns, nil
	}

	return "", DNSNotFoundError{ip}
}

// ServerClient is a client for the servers API.
type ServerClient struct {
	client *Client
}

// GetByID retrieves a server by its ID. If the server does not exist, nil is returned.
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

// GetByName retrieves a server by its name. If the server does not exist, nil is returned.
func (c *ServerClient) GetByName(ctx context.Context, name string) (*Server, *Response, error) {
	if name == "" {
		return nil, nil, nil
	}
	servers, response, err := c.List(ctx, ServerListOpts{Name: name})
	if len(servers) == 0 {
		return nil, response, err
	}
	return servers[0], response, err
}

// Get retrieves a server by its ID if the input can be parsed as an integer, otherwise it
// retrieves a server by its name. If the server does not exist, nil is returned.
func (c *ServerClient) Get(ctx context.Context, idOrName string) (*Server, *Response, error) {
	if id, err := strconv.Atoi(idOrName); err == nil {
		return c.GetByID(ctx, int(id))
	}
	return c.GetByName(ctx, idOrName)
}

// ServerListOpts specifies options for listing servers.
type ServerListOpts struct {
	ListOpts
	Name   string
	Status []ServerStatus
	Sort   []string
}

func (l ServerListOpts) values() url.Values {
	vals := l.ListOpts.values()
	if l.Name != "" {
		vals.Add("name", l.Name)
	}
	for _, status := range l.Status {
		vals.Add("status", string(status))
	}
	for _, sort := range l.Sort {
		vals.Add("sort", sort)
	}
	return vals
}

// List returns a list of servers for a specific page.
//
// Please note that filters specified in opts are not taken into account
// when their value corresponds to their zero value or when they are empty.
func (c *ServerClient) List(ctx context.Context, opts ServerListOpts) ([]*Server, *Response, error) {
	path := "/servers?" + opts.values().Encode()
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
	return c.AllWithOpts(ctx, ServerListOpts{ListOpts: ListOpts{PerPage: 50}})
}

// AllWithOpts returns all servers for the given options.
func (c *ServerClient) AllWithOpts(ctx context.Context, opts ServerListOpts) ([]*Server, error) {
	allServers := []*Server{}

	err := c.client.all(func(page int) (*Response, error) {
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
	Name             string
	ServerType       *ServerType
	Image            *Image
	SSHKeys          []*SSHKey
	Location         *Location
	Datacenter       *Datacenter
	UserData         string
	StartAfterCreate *bool
	Labels           map[string]string
	Automount        *bool
	Volumes          []*Volume
	Networks         []*Network
	Firewalls        []*ServerCreateFirewall
	PlacementGroup   *PlacementGroup
	PublicNet        *ServerCreatePublicNet
}

type ServerCreatePublicNet struct {
	EnableIPv4 bool
	EnableIPv6 bool
	IPv4       *PrimaryIP
	IPv6       *PrimaryIP
}

// ServerCreateFirewall defines which Firewalls to apply when creating a Server.
type ServerCreateFirewall struct {
	Firewall Firewall
}

// Validate checks if options are valid.
func (o ServerCreateOpts) Validate() error {
	if o.Name == "" {
		return errors.New("missing name")
	}
	if o.ServerType == nil || (o.ServerType.ID == 0 && o.ServerType.Name == "") {
		return errors.New("missing server type")
	}
	if o.Image == nil || (o.Image.ID == 0 && o.Image.Name == "") {
		return errors.New("missing image")
	}
	if o.Location != nil && o.Datacenter != nil {
		return errors.New("location and datacenter are mutually exclusive")
	}
	if o.PublicNet != nil {
		if !o.PublicNet.EnableIPv4 && !o.PublicNet.EnableIPv6 &&
			len(o.Networks) == 0 && (o.StartAfterCreate == nil || *o.StartAfterCreate == true) {
			return errors.New("missing networks or StartAfterCreate == false when EnableIPv4 and EnableIPv6 is false")
		}
	}
	return nil
}

// ServerCreateResult is the result of a create server call.
type ServerCreateResult struct {
	Server       *Server
	Action       *Action
	RootPassword string
	NextActions  []*Action
}

// Create creates a new server.
func (c *ServerClient) Create(ctx context.Context, opts ServerCreateOpts) (ServerCreateResult, *Response, error) {
	if err := opts.Validate(); err != nil {
		return ServerCreateResult{}, nil, err
	}

	var reqBody schema.ServerCreateRequest
	reqBody.UserData = opts.UserData
	reqBody.Name = opts.Name
	reqBody.Automount = opts.Automount
	reqBody.StartAfterCreate = opts.StartAfterCreate
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
	if opts.Labels != nil {
		reqBody.Labels = &opts.Labels
	}
	for _, sshKey := range opts.SSHKeys {
		reqBody.SSHKeys = append(reqBody.SSHKeys, sshKey.ID)
	}
	for _, volume := range opts.Volumes {
		reqBody.Volumes = append(reqBody.Volumes, volume.ID)
	}
	for _, network := range opts.Networks {
		reqBody.Networks = append(reqBody.Networks, network.ID)
	}
	for _, firewall := range opts.Firewalls {
		reqBody.Firewalls = append(reqBody.Firewalls, schema.ServerCreateFirewalls{
			Firewall: firewall.Firewall.ID,
		})
	}

	if opts.PublicNet != nil {
		reqBody.PublicNet = &schema.ServerCreatePublicNet{
			EnableIPv4: opts.PublicNet.EnableIPv4,
			EnableIPv6: opts.PublicNet.EnableIPv6,
		}
		if opts.PublicNet.IPv4 != nil {
			reqBody.PublicNet.IPv4ID = opts.PublicNet.IPv4.ID
		}
		if opts.PublicNet.IPv6 != nil {
			reqBody.PublicNet.IPv6ID = opts.PublicNet.IPv6.ID
		}
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
	if opts.PlacementGroup != nil {
		reqBody.PlacementGroup = opts.PlacementGroup.ID
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
		Server:      ServerFromSchema(respBody.Server),
		Action:      ActionFromSchema(respBody.Action),
		NextActions: ActionsFromSchema(respBody.NextActions),
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

// ServerUpdateOpts specifies options for updating a server.
type ServerUpdateOpts struct {
	Name   string
	Labels map[string]string
}

// Update updates a server.
func (c *ServerClient) Update(ctx context.Context, server *Server, opts ServerUpdateOpts) (*Server, *Response, error) {
	reqBody := schema.ServerUpdateRequest{
		Name: opts.Name,
	}
	if opts.Labels != nil {
		reqBody.Labels = &opts.Labels
	}
	reqBodyData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, nil, err
	}

	path := fmt.Sprintf("/servers/%d", server.ID)
	req, err := c.client.NewRequest(ctx, "PUT", path, bytes.NewReader(reqBodyData))
	if err != nil {
		return nil, nil, err
	}

	respBody := schema.ServerUpdateResponse{}
	resp, err := c.client.Do(req, &respBody)
	if err != nil {
		return nil, resp, err
	}
	return ServerFromSchema(respBody.Server), resp, nil
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
	Labels      map[string]string
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
		if opts.Labels != nil {
			reqBody.Labels = &opts.Labels
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

// ServerRebuildOpts specifies options for rebuilding a server.
type ServerRebuildOpts struct {
	Image *Image
}

// Rebuild rebuilds a server.
func (c *ServerClient) Rebuild(ctx context.Context, server *Server, opts ServerRebuildOpts) (*Action, *Response, error) {
	reqBody := schema.ServerActionRebuildRequest{}
	if opts.Image.ID != 0 {
		reqBody.Image = opts.Image.ID
	} else {
		reqBody.Image = opts.Image.Name
	}
	reqBodyData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, nil, err
	}

	path := fmt.Sprintf("/servers/%d/actions/rebuild", server.ID)
	req, err := c.client.NewRequest(ctx, "POST", path, bytes.NewReader(reqBodyData))
	if err != nil {
		return nil, nil, err
	}

	respBody := schema.ServerActionRebuildResponse{}
	resp, err := c.client.Do(req, &respBody)
	if err != nil {
		return nil, resp, err
	}
	return ActionFromSchema(respBody.Action), resp, nil
}

// AttachISO attaches an ISO to a server.
func (c *ServerClient) AttachISO(ctx context.Context, server *Server, iso *ISO) (*Action, *Response, error) {
	reqBody := schema.ServerActionAttachISORequest{}
	if iso.ID != 0 {
		reqBody.ISO = iso.ID
	} else {
		reqBody.ISO = iso.Name
	}
	reqBodyData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, nil, err
	}

	path := fmt.Sprintf("/servers/%d/actions/attach_iso", server.ID)
	req, err := c.client.NewRequest(ctx, "POST", path, bytes.NewReader(reqBodyData))
	if err != nil {
		return nil, nil, err
	}

	respBody := schema.ServerActionAttachISOResponse{}
	resp, err := c.client.Do(req, &respBody)
	if err != nil {
		return nil, resp, err
	}
	return ActionFromSchema(respBody.Action), resp, nil
}

// DetachISO detaches the currently attached ISO from a server.
func (c *ServerClient) DetachISO(ctx context.Context, server *Server) (*Action, *Response, error) {
	path := fmt.Sprintf("/servers/%d/actions/detach_iso", server.ID)
	req, err := c.client.NewRequest(ctx, "POST", path, nil)
	if err != nil {
		return nil, nil, err
	}

	respBody := schema.ServerActionDetachISOResponse{}
	resp, err := c.client.Do(req, &respBody)
	if err != nil {
		return nil, resp, err
	}
	return ActionFromSchema(respBody.Action), resp, nil
}

// EnableBackup enables backup for a server. Pass in an empty backup window to let the
// API pick a window for you. See the API documentation at docs.hetzner.cloud for a list
// of valid backup windows.
func (c *ServerClient) EnableBackup(ctx context.Context, server *Server, window string) (*Action, *Response, error) {
	reqBody := schema.ServerActionEnableBackupRequest{}
	if window != "" {
		reqBody.BackupWindow = String(window)
	}
	reqBodyData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, nil, err
	}

	path := fmt.Sprintf("/servers/%d/actions/enable_backup", server.ID)
	req, err := c.client.NewRequest(ctx, "POST", path, bytes.NewReader(reqBodyData))
	if err != nil {
		return nil, nil, err
	}

	respBody := schema.ServerActionEnableBackupResponse{}
	resp, err := c.client.Do(req, &respBody)
	if err != nil {
		return nil, resp, err
	}
	return ActionFromSchema(respBody.Action), resp, nil
}

// DisableBackup disables backup for a server.
func (c *ServerClient) DisableBackup(ctx context.Context, server *Server) (*Action, *Response, error) {
	path := fmt.Sprintf("/servers/%d/actions/disable_backup", server.ID)
	req, err := c.client.NewRequest(ctx, "POST", path, nil)
	if err != nil {
		return nil, nil, err
	}

	respBody := schema.ServerActionDisableBackupResponse{}
	resp, err := c.client.Do(req, &respBody)
	if err != nil {
		return nil, resp, err
	}
	return ActionFromSchema(respBody.Action), resp, nil
}

// ServerChangeTypeOpts specifies options for changing a server's type.
type ServerChangeTypeOpts struct {
	ServerType  *ServerType // new server type
	UpgradeDisk bool        // whether disk should be upgraded
}

// ChangeType changes a server's type.
func (c *ServerClient) ChangeType(ctx context.Context, server *Server, opts ServerChangeTypeOpts) (*Action, *Response, error) {
	reqBody := schema.ServerActionChangeTypeRequest{
		UpgradeDisk: opts.UpgradeDisk,
	}
	if opts.ServerType.ID != 0 {
		reqBody.ServerType = opts.ServerType.ID
	} else {
		reqBody.ServerType = opts.ServerType.Name
	}
	reqBodyData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, nil, err
	}

	path := fmt.Sprintf("/servers/%d/actions/change_type", server.ID)
	req, err := c.client.NewRequest(ctx, "POST", path, bytes.NewReader(reqBodyData))
	if err != nil {
		return nil, nil, err
	}

	respBody := schema.ServerActionChangeTypeResponse{}
	resp, err := c.client.Do(req, &respBody)
	if err != nil {
		return nil, resp, err
	}
	return ActionFromSchema(respBody.Action), resp, nil
}

// ChangeDNSPtr changes or resets the reverse DNS pointer for a server IP address.
// Pass a nil ptr to reset the reverse DNS pointer to its default value.
func (c *ServerClient) ChangeDNSPtr(ctx context.Context, server *Server, ip string, ptr *string) (*Action, *Response, error) {
	netIP := net.ParseIP(ip)
	if netIP == nil {
		return nil, nil, InvalidIPError{ip}
	}
	return server.changeDNSPtr(ctx, c.client, net.ParseIP(ip), ptr)
}

// ServerChangeProtectionOpts specifies options for changing the resource protection level of a server.
type ServerChangeProtectionOpts struct {
	Rebuild *bool
	Delete  *bool
}

// ChangeProtection changes the resource protection level of a server.
func (c *ServerClient) ChangeProtection(ctx context.Context, server *Server, opts ServerChangeProtectionOpts) (*Action, *Response, error) {
	reqBody := schema.ServerActionChangeProtectionRequest{
		Rebuild: opts.Rebuild,
		Delete:  opts.Delete,
	}
	reqBodyData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, nil, err
	}

	path := fmt.Sprintf("/servers/%d/actions/change_protection", server.ID)
	req, err := c.client.NewRequest(ctx, "POST", path, bytes.NewReader(reqBodyData))
	if err != nil {
		return nil, nil, err
	}

	respBody := schema.ServerActionChangeProtectionResponse{}
	resp, err := c.client.Do(req, &respBody)
	if err != nil {
		return nil, resp, err
	}
	return ActionFromSchema(respBody.Action), resp, err
}

// ServerRequestConsoleResult is the result of requesting a WebSocket VNC console.
type ServerRequestConsoleResult struct {
	Action   *Action
	WSSURL   string
	Password string
}

// RequestConsole requests a WebSocket VNC console.
func (c *ServerClient) RequestConsole(ctx context.Context, server *Server) (ServerRequestConsoleResult, *Response, error) {
	path := fmt.Sprintf("/servers/%d/actions/request_console", server.ID)
	req, err := c.client.NewRequest(ctx, "POST", path, nil)
	if err != nil {
		return ServerRequestConsoleResult{}, nil, err
	}

	respBody := schema.ServerActionRequestConsoleResponse{}
	resp, err := c.client.Do(req, &respBody)
	if err != nil {
		return ServerRequestConsoleResult{}, resp, err
	}
	return ServerRequestConsoleResult{
		Action:   ActionFromSchema(respBody.Action),
		WSSURL:   respBody.WSSURL,
		Password: respBody.Password,
	}, resp, nil
}

// ServerAttachToNetworkOpts specifies options for attaching a server to a network.
type ServerAttachToNetworkOpts struct {
	Network  *Network
	IP       net.IP
	AliasIPs []net.IP
}

// AttachToNetwork attaches a server to a network.
func (c *ServerClient) AttachToNetwork(ctx context.Context, server *Server, opts ServerAttachToNetworkOpts) (*Action, *Response, error) {
	reqBody := schema.ServerActionAttachToNetworkRequest{
		Network: opts.Network.ID,
	}
	if opts.IP != nil {
		reqBody.IP = String(opts.IP.String())
	}
	for _, aliasIP := range opts.AliasIPs {
		reqBody.AliasIPs = append(reqBody.AliasIPs, String(aliasIP.String()))
	}
	reqBodyData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, nil, err
	}

	path := fmt.Sprintf("/servers/%d/actions/attach_to_network", server.ID)
	req, err := c.client.NewRequest(ctx, "POST", path, bytes.NewReader(reqBodyData))
	if err != nil {
		return nil, nil, err
	}

	respBody := schema.ServerActionAttachToNetworkResponse{}
	resp, err := c.client.Do(req, &respBody)
	if err != nil {
		return nil, resp, err
	}
	return ActionFromSchema(respBody.Action), resp, err
}

// ServerDetachFromNetworkOpts specifies options for detaching a server from a network.
type ServerDetachFromNetworkOpts struct {
	Network *Network
}

// DetachFromNetwork detaches a server from a network.
func (c *ServerClient) DetachFromNetwork(ctx context.Context, server *Server, opts ServerDetachFromNetworkOpts) (*Action, *Response, error) {
	reqBody := schema.ServerActionDetachFromNetworkRequest{
		Network: opts.Network.ID,
	}
	reqBodyData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, nil, err
	}

	path := fmt.Sprintf("/servers/%d/actions/detach_from_network", server.ID)
	req, err := c.client.NewRequest(ctx, "POST", path, bytes.NewReader(reqBodyData))
	if err != nil {
		return nil, nil, err
	}

	respBody := schema.ServerActionDetachFromNetworkResponse{}
	resp, err := c.client.Do(req, &respBody)
	if err != nil {
		return nil, resp, err
	}
	return ActionFromSchema(respBody.Action), resp, err
}

// ServerChangeAliasIPsOpts specifies options for changing the alias ips of an already attached network.
type ServerChangeAliasIPsOpts struct {
	Network  *Network
	AliasIPs []net.IP
}

// ChangeAliasIPs changes a server's alias IPs in a network.
func (c *ServerClient) ChangeAliasIPs(ctx context.Context, server *Server, opts ServerChangeAliasIPsOpts) (*Action, *Response, error) {
	reqBody := schema.ServerActionChangeAliasIPsRequest{
		Network:  opts.Network.ID,
		AliasIPs: []string{},
	}
	for _, aliasIP := range opts.AliasIPs {
		reqBody.AliasIPs = append(reqBody.AliasIPs, aliasIP.String())
	}
	reqBodyData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, nil, err
	}
	path := fmt.Sprintf("/servers/%d/actions/change_alias_ips", server.ID)
	req, err := c.client.NewRequest(ctx, "POST", path, bytes.NewReader(reqBodyData))
	if err != nil {
		return nil, nil, err
	}

	respBody := schema.ServerActionDetachFromNetworkResponse{}
	resp, err := c.client.Do(req, &respBody)
	if err != nil {
		return nil, resp, err
	}
	return ActionFromSchema(respBody.Action), resp, err
}

// ServerMetricType is the type of available metrics for servers.
type ServerMetricType string

// Available types of server metrics. See Hetzner Cloud API documentation for
// details.
const (
	ServerMetricCPU     ServerMetricType = "cpu"
	ServerMetricDisk    ServerMetricType = "disk"
	ServerMetricNetwork ServerMetricType = "network"
)

// ServerGetMetricsOpts configures the call to get metrics for a Server.
type ServerGetMetricsOpts struct {
	Types []ServerMetricType
	Start time.Time
	End   time.Time
	Step  int
}

func (o *ServerGetMetricsOpts) addQueryParams(req *http.Request) error {
	query := req.URL.Query()

	if len(o.Types) == 0 {
		return fmt.Errorf("no metric types specified")
	}
	for _, typ := range o.Types {
		query.Add("type", string(typ))
	}

	if o.Start.IsZero() {
		return fmt.Errorf("no start time specified")
	}
	query.Add("start", o.Start.Format(time.RFC3339))

	if o.End.IsZero() {
		return fmt.Errorf("no end time specified")
	}
	query.Add("end", o.End.Format(time.RFC3339))

	if o.Step > 0 {
		query.Add("step", strconv.Itoa(o.Step))
	}
	req.URL.RawQuery = query.Encode()

	return nil
}

// ServerMetrics contains the metrics requested for a Server.
type ServerMetrics struct {
	Start      time.Time
	End        time.Time
	Step       float64
	TimeSeries map[string][]ServerMetricsValue
}

// ServerMetricsValue represents a single value in a time series of metrics.
type ServerMetricsValue struct {
	Timestamp float64
	Value     string
}

// GetMetrics obtains metrics for Server.
func (c *ServerClient) GetMetrics(ctx context.Context, server *Server, opts ServerGetMetricsOpts) (*ServerMetrics, *Response, error) {
	var respBody schema.ServerGetMetricsResponse

	if server == nil {
		return nil, nil, fmt.Errorf("illegal argument: server is nil")
	}

	path := fmt.Sprintf("/servers/%d/metrics", server.ID)
	req, err := c.client.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("new request: %v", err)
	}
	if err := opts.addQueryParams(req); err != nil {
		return nil, nil, fmt.Errorf("add query params: %v", err)
	}
	resp, err := c.client.Do(req, &respBody)
	if err != nil {
		return nil, nil, fmt.Errorf("get metrics: %v", err)
	}
	ms, err := serverMetricsFromSchema(&respBody)
	if err != nil {
		return nil, nil, fmt.Errorf("convert response body: %v", err)
	}
	return ms, resp, nil
}

func (c *ServerClient) AddToPlacementGroup(ctx context.Context, server *Server, placementGroup *PlacementGroup) (*Action, *Response, error) {
	reqBody := schema.ServerActionAddToPlacementGroupRequest{
		PlacementGroup: placementGroup.ID,
	}
	reqBodyData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, nil, err
	}
	path := fmt.Sprintf("/servers/%d/actions/add_to_placement_group", server.ID)
	req, err := c.client.NewRequest(ctx, "POST", path, bytes.NewReader(reqBodyData))
	if err != nil {
		return nil, nil, err
	}

	respBody := schema.ServerActionAddToPlacementGroupResponse{}
	resp, err := c.client.Do(req, &respBody)
	if err != nil {
		return nil, resp, err
	}
	return ActionFromSchema(respBody.Action), resp, err
}

func (c *ServerClient) RemoveFromPlacementGroup(ctx context.Context, server *Server) (*Action, *Response, error) {
	path := fmt.Sprintf("/servers/%d/actions/remove_from_placement_group", server.ID)
	req, err := c.client.NewRequest(ctx, "POST", path, nil)
	if err != nil {
		return nil, nil, err
	}

	respBody := schema.ServerActionRemoveFromPlacementGroupResponse{}
	resp, err := c.client.Do(req, &respBody)
	if err != nil {
		return nil, resp, err
	}
	return ActionFromSchema(respBody.Action), resp, err
}
