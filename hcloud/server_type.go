package hcloud

import "github.com/hetznercloud/hcloud-go/hcloud/schema"

// ServerType represents a server type in the Hetzner Cloud.
type ServerType struct {
	ID          int
	Name        string
	Description string
	Cores       int
	Memory      float32
	Disk        int
	StorageType StorageType
}

// ServerTypeFromSchema converts a schema.ServerType to a ServerType.
func ServerTypeFromSchema(s schema.ServerType) ServerType {
	return ServerType{
		ID:          s.ID,
		Name:        s.Name,
		Description: s.Description,
		Cores:       s.Cores,
		Memory:      s.Memory,
		Disk:        s.Disk,
		StorageType: StorageType(s.StorageType),
	}
}

// StorageType specifies the type of storage.
type StorageType string

const (
	// StorageTypeLocal is the type for local storage.
	StorageTypeLocal StorageType = "local"

	// StorageTypeCeph is the type for remote storage.
	StorageTypeCeph = "ceph"
)
