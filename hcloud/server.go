package hcloud

import (
	"context"
	"encoding/json"
	"fmt"
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

type ServerStatus string

const (
	ServerStatusInitializing ServerStatus = "initializing"
	ServerStatusOff                       = "off"
	ServerStatusRunning                   = "running"
)

type ServerPublicNet struct {
	IPv4        ServerPublicNetIPv4
	IPv6        ServerPublicNetIPv6
	FloatingIPs []*FloatingIP
}

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

type ServerPublicNetIPv4 struct {
	IP      string
	Blocked bool
	DNSPtr  string
}

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

type ServerPublicNetIPv6 struct {
	IP      string
	Blocked bool
	DNSPtr  []ServerPublicNetIPv6DNSPtr
}

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

type ServerPublicNetIPv6DNSPtr struct {
	IP     string
	DNSPtr string
}

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

type ServerClient struct {
	client *Client
}

func (c *ServerClient) Get(ctx context.Context, id int) (*Server, *Response, error) {
	req, err := c.client.NewRequest(ctx, "GET", fmt.Sprintf("/servers/%d", id))
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

func (c *ServerClient) List(ctx context.Context) ([]*Server, *Response, error) {
	req, err := c.client.NewRequest(ctx, "GET", "/servers")
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
