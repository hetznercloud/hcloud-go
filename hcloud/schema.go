package hcloud

import (
	"net"

	"github.com/hetznercloud/hcloud-go/hcloud/schema"
)

// This file provides converter functions to convert models in the
// schema package to models in the hcloud package.

// ActionFromSchema converts a schema.Action to an Action.
func ActionFromSchema(s schema.Action) *Action {
	action := &Action{
		ID:        s.ID,
		Status:    ActionStatus(s.Status),
		Command:   s.Command,
		Progress:  s.Progress,
		Started:   s.Started,
		Resources: []*ResourceReference{},
	}
	if s.Finished != nil {
		action.Finished = *s.Finished
	}
	if s.Error != nil {
		action.ErrorCode = s.Error.Code
		action.ErrorMessage = s.Error.Message
	}
	for _, r := range s.Resources {
		action.Resources = append(action.Resources, &ResourceReference{
			ID:   r.ID,
			Type: ResourceReferenceType(r.Type),
		})
	}
	return action
}

// FloatingIPFromSchema converts a schema.FloatingIP to a FloatingIP.
func FloatingIPFromSchema(s schema.FloatingIP) *FloatingIP {
	f := &FloatingIP{
		ID:           s.ID,
		Type:         FloatingIPType(s.Type),
		HomeLocation: LocationFromSchema(s.HomeLocation),
		Blocked:      s.Blocked,
	}
	if s.Description != nil {
		f.Description = *s.Description
	}
	if s.Server != nil {
		f.Server = &Server{ID: *s.Server}
	}
	if s.Type == "ipv4" {
		f.IP = net.ParseIP(s.IP)
	} else {
		f.IP, f.Network, _ = net.ParseCIDR(s.IP)
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

// ISOFromSchema converts a schema.ISO to an ISO.
func ISOFromSchema(s schema.ISO) *ISO {
	return &ISO{
		ID:          s.ID,
		Name:        s.Name,
		Description: s.Description,
		Type:        ISOType(s.Type),
	}
}

// LocationFromSchema converts a schema.Location to a Location.
func LocationFromSchema(s schema.Location) *Location {
	return &Location{
		ID:          s.ID,
		Name:        s.Name,
		Description: s.Description,
		Country:     s.Country,
		City:        s.City,
		Latitude:    s.Latitude,
		Longitude:   s.Longitude,
	}
}

// DatacenterFromSchema converts a schema.Datacenter to a Datacenter.
func DatacenterFromSchema(s schema.Datacenter) *Datacenter {
	return &Datacenter{
		ID:          s.ID,
		Name:        s.Name,
		Description: s.Description,
		Location:    LocationFromSchema(s.Location),
	}
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
		RescueEnabled:   s.RescueEnabled,
		Datacenter:      DatacenterFromSchema(s.Datacenter),
		Locked:          s.Locked,
	}
	if s.BackupWindow != nil {
		server.BackupWindow = *s.BackupWindow
	}
	if s.OutgoingTraffic != nil {
		server.OutgoingTraffic = *s.OutgoingTraffic
	}
	if s.IngoingTraffic != nil {
		server.IngoingTraffic = *s.IngoingTraffic
	}
	if s.ISO != nil {
		server.ISO = ISOFromSchema(*s.ISO)
	}
	return server
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

// ServerPublicNetIPv4FromSchema converts a schema.ServerPublicNetIPv4 to
// a ServerPublicNetIPv4.
func ServerPublicNetIPv4FromSchema(s schema.ServerPublicNetIPv4) ServerPublicNetIPv4 {
	return ServerPublicNetIPv4{
		IP:      net.ParseIP(s.IP),
		Blocked: s.Blocked,
		DNSPtr:  s.DNSPtr,
	}
}

// ServerPublicNetIPv6FromSchema converts a schema.ServerPublicNetIPv6 to
// a ServerPublicNetIPv6.
func ServerPublicNetIPv6FromSchema(s schema.ServerPublicNetIPv6) ServerPublicNetIPv6 {
	ipv6 := ServerPublicNetIPv6{
		Blocked: s.Blocked,
		DNSPtr:  map[string]string{},
	}
	ipv6.IP, ipv6.Network, _ = net.ParseCIDR(s.IP)

	for _, dnsPtr := range s.DNSPtr {
		ipv6.DNSPtr[dnsPtr.IP] = dnsPtr.DNSPtr
	}
	return ipv6
}

// ServerTypeFromSchema converts a schema.ServerType to a ServerType.
func ServerTypeFromSchema(s schema.ServerType) *ServerType {
	st := &ServerType{
		ID:          s.ID,
		Name:        s.Name,
		Description: s.Description,
		Cores:       s.Cores,
		Memory:      s.Memory,
		Disk:        s.Disk,
		StorageType: StorageType(s.StorageType),
	}

	return st
}

// SSHKeyFromSchema converts a schema.SSHKey to a SSHKey.
func SSHKeyFromSchema(s schema.SSHKey) *SSHKey {
	return &SSHKey{
		ID:          s.ID,
		Name:        s.Name,
		Fingerprint: s.Fingerprint,
		PublicKey:   s.PublicKey,
	}
}

// ImageFromSchema converts a schema.Image to an Image.
func ImageFromSchema(s schema.Image) *Image {
	i := &Image{
		ID:          s.ID,
		Type:        ImageType(s.Type),
		Status:      ImageStatus(s.Status),
		Description: s.Description,
		DiskSize:    s.DiskSize,
		Created:     s.Created,
		RapidDeploy: s.RapidDeploy,
		OSFlavor:    s.OSFlavor,
	}
	if s.Name != nil {
		i.Name = *s.Name
	}
	if s.ImageSize != nil {
		i.ImageSize = *s.ImageSize
	}
	if s.OSVersion != nil {
		i.OSVersion = *s.OSVersion
	}
	if s.CreatedFrom != nil {
		i.CreatedFrom = &Server{
			ID:   s.CreatedFrom.ID,
			Name: s.CreatedFrom.Name,
		}
	}
	if s.BoundTo != nil {
		i.BoundTo = &Server{
			ID: *s.BoundTo,
		}
	}
	return i
}

// PaginationFromSchema converts a schema.MetaPagination to a Pagination.
func PaginationFromSchema(s schema.MetaPagination) Pagination {
	return Pagination{
		Page:         s.Page,
		PerPage:      s.PerPage,
		PreviousPage: s.PreviousPage,
		NextPage:     s.NextPage,
		LastPage:     s.LastPage,
		TotalEntries: s.TotalEntries,
	}
}

// ErrorFromSchema converts a schema.Error to an Error.
func ErrorFromSchema(s schema.Error) Error {
	return Error{
		Code:    ErrorCode(s.Code),
		Message: s.Message,
	}
}
