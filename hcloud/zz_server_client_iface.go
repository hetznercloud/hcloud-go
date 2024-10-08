// Code generated by ifacemaker; DO NOT EDIT.

package hcloud

import (
	"context"
)

// IServerClient ...
type IServerClient interface {
	// GetByID retrieves a server by its ID. If the server does not exist, nil is returned.
	GetByID(ctx context.Context, id int64) (*Server, *Response, error)
	// GetByName retrieves a server by its name. If the server does not exist, nil is returned.
	GetByName(ctx context.Context, name string) (*Server, *Response, error)
	// Get retrieves a server by its ID if the input can be parsed as an integer, otherwise it
	// retrieves a server by its name. If the server does not exist, nil is returned.
	Get(ctx context.Context, idOrName string) (*Server, *Response, error)
	// List returns a list of servers for a specific page.
	//
	// Please note that filters specified in opts are not taken into account
	// when their value corresponds to their zero value or when they are empty.
	List(ctx context.Context, opts ServerListOpts) ([]*Server, *Response, error)
	// All returns all servers.
	All(ctx context.Context) ([]*Server, error)
	// AllWithOpts returns all servers for the given options.
	AllWithOpts(ctx context.Context, opts ServerListOpts) ([]*Server, error)
	// Create creates a new server.
	Create(ctx context.Context, opts ServerCreateOpts) (ServerCreateResult, *Response, error)
	// Delete deletes a server.
	//
	// Deprecated: Use [ServerClient.DeleteWithResult] instead.
	Delete(ctx context.Context, server *Server) (*Response, error)
	// DeleteWithResult deletes a server and returns the parsed response containing the action.
	DeleteWithResult(ctx context.Context, server *Server) (*ServerDeleteResult, *Response, error)
	// Update updates a server.
	Update(ctx context.Context, server *Server, opts ServerUpdateOpts) (*Server, *Response, error)
	// Poweron starts a server.
	Poweron(ctx context.Context, server *Server) (*Action, *Response, error)
	// Reboot reboots a server.
	Reboot(ctx context.Context, server *Server) (*Action, *Response, error)
	// Reset resets a server.
	Reset(ctx context.Context, server *Server) (*Action, *Response, error)
	// Shutdown shuts down a server.
	Shutdown(ctx context.Context, server *Server) (*Action, *Response, error)
	// Poweroff stops a server.
	Poweroff(ctx context.Context, server *Server) (*Action, *Response, error)
	// ResetPassword resets a server's password.
	ResetPassword(ctx context.Context, server *Server) (ServerResetPasswordResult, *Response, error)
	// CreateImage creates an image from a server.
	CreateImage(ctx context.Context, server *Server, opts *ServerCreateImageOpts) (ServerCreateImageResult, *Response, error)
	// EnableRescue enables rescue mode for a server.
	EnableRescue(ctx context.Context, server *Server, opts ServerEnableRescueOpts) (ServerEnableRescueResult, *Response, error)
	// DisableRescue disables rescue mode for a server.
	DisableRescue(ctx context.Context, server *Server) (*Action, *Response, error)
	// Rebuild rebuilds a server.
	//
	// Deprecated: Use [ServerClient.RebuildWithResult] instead.
	Rebuild(ctx context.Context, server *Server, opts ServerRebuildOpts) (*Action, *Response, error)
	// RebuildWithResult rebuilds a server.
	RebuildWithResult(ctx context.Context, server *Server, opts ServerRebuildOpts) (ServerRebuildResult, *Response, error)
	// AttachISO attaches an ISO to a server.
	AttachISO(ctx context.Context, server *Server, iso *ISO) (*Action, *Response, error)
	// DetachISO detaches the currently attached ISO from a server.
	DetachISO(ctx context.Context, server *Server) (*Action, *Response, error)
	// EnableBackup enables backup for a server.
	// The window parameter is deprecated and will be ignored.
	EnableBackup(ctx context.Context, server *Server, window string) (*Action, *Response, error)
	// DisableBackup disables backup for a server.
	DisableBackup(ctx context.Context, server *Server) (*Action, *Response, error)
	// ChangeType changes a server's type.
	ChangeType(ctx context.Context, server *Server, opts ServerChangeTypeOpts) (*Action, *Response, error)
	// ChangeDNSPtr changes or resets the reverse DNS pointer for a server IP address.
	// Pass a nil ptr to reset the reverse DNS pointer to its default value.
	ChangeDNSPtr(ctx context.Context, server *Server, ip string, ptr *string) (*Action, *Response, error)
	// ChangeProtection changes the resource protection level of a server.
	ChangeProtection(ctx context.Context, server *Server, opts ServerChangeProtectionOpts) (*Action, *Response, error)
	// RequestConsole requests a WebSocket VNC console.
	RequestConsole(ctx context.Context, server *Server) (ServerRequestConsoleResult, *Response, error)
	// AttachToNetwork attaches a server to a network.
	AttachToNetwork(ctx context.Context, server *Server, opts ServerAttachToNetworkOpts) (*Action, *Response, error)
	// DetachFromNetwork detaches a server from a network.
	DetachFromNetwork(ctx context.Context, server *Server, opts ServerDetachFromNetworkOpts) (*Action, *Response, error)
	// ChangeAliasIPs changes a server's alias IPs in a network.
	ChangeAliasIPs(ctx context.Context, server *Server, opts ServerChangeAliasIPsOpts) (*Action, *Response, error)
	// GetMetrics obtains metrics for Server.
	GetMetrics(ctx context.Context, server *Server, opts ServerGetMetricsOpts) (*ServerMetrics, *Response, error)
	AddToPlacementGroup(ctx context.Context, server *Server, placementGroup *PlacementGroup) (*Action, *Response, error)
	RemoveFromPlacementGroup(ctx context.Context, server *Server) (*Action, *Response, error)
}
