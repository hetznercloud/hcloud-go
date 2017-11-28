package hcloud // import "hetzner.cloud/hcloud"

import "encoding/json"

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

// UnmarshalJSON implements json.Unmarshaler.
func (s *ServerType) UnmarshalJSON(data []byte) error {
	var serverType struct {
		ID          int     `json:"id"`
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Cores       int     `json:"cores"`
		Memory      float32 `json:"memory"`
		Disk        int     `json:"disk"`
		StorageType string  `json:"storage_type"`
	}

	if err := json.Unmarshal(data, &serverType); err != nil {
		return err
	}

	s.ID = serverType.ID
	s.Name = serverType.Name
	s.Description = serverType.Description
	s.Cores = serverType.Cores
	s.Memory = serverType.Memory
	s.Disk = serverType.Disk
	s.StorageType = StorageType(serverType.StorageType)

	return nil
}

// StorageType specifies the type of storage.
type StorageType string

const (
	StorageTypeLocal StorageType = "local" // Local storage
	StorageTypeCeph              = "ceph"  // Remote storage
)
