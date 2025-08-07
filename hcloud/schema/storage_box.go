package schema

import "time"

// StorageBox defines the schema of a storage box.
type StorageBox struct {
	ID             int64                    `json:"id"`
	Username       *string                  `json:"username,omitempty"`
	Status         string                   `json:"status"`
	Name           string                   `json:"name"`
	StorageBoxType StorageBoxType           `json:"storage_box_type"`
	Location       Location                 `json:"location"`
	AccessSettings StorageBoxAccessSettings `json:"access_settings"`
	Server         *string                  `json:"server"`
	System         *string                  `json:"system"`
	Stats          *StorageBoxStats         `json:"stats"`
	Labels         map[string]string        `json:"labels"`
	Protection     StorageBoxProtection     `json:"protection"`
	SnapshotPlan   *StorageBoxSnapshotPlan  `json:"snapshot_plan"`
	Created        time.Time                `json:"created"`
}

// StorageBoxAccessSettings defines the schema of a storage box's access settings.
type StorageBoxAccessSettings struct {
	ReachableExternally bool `json:"reachable_externally"`
	SambaEnabled        bool `json:"samba_enabled"`
	SSHEnabled          bool `json:"ssh_enabled"`
	WebDAVEnabled       bool `json:"webdav_enabled"`
	ZFSEnabled          bool `json:"zfs_enabled"`
}

// StorageBoxStats defines the schema of a storage box's disk usage statistics.
type StorageBoxStats struct {
	Size          int64 `json:"size"`
	SizeData      int64 `json:"size_data"`
	SizeSnapshots int64 `json:"size_snapshots"`
}

// StorageBoxProtection defines the schema of a storage box's resource protection.
type StorageBoxProtection struct {
	Delete bool `json:"delete"`
}

// StorageBoxSnapshotPlan defines the schema of a storage box's snapshot plan.
type StorageBoxSnapshotPlan struct {
	MaxSnapshots int  `json:"max_snapshots"`
	Minute       *int `json:"minute,omitempty"`
	Hour         *int `json:"hour,omitempty"`
	DayOfWeek    *int `json:"day_of_week,omitempty"`
	DayOfMonth   *int `json:"day_of_month,omitempty"`
}

// StorageBoxGetResponse defines the schema of the response when
// retrieving a single storage box.
type StorageBoxGetResponse struct {
	StorageBox StorageBox `json:"storage_box"`
}

// StorageBoxListResponse defines the schema of the response when
// listing storage boxes.
type StorageBoxListResponse struct {
	StorageBoxes []StorageBox `json:"storage_boxes"`
}

// StorageBoxCreateRequest defines the schema for the request to
// create a storage box.
type StorageBoxCreateRequest struct {
	Name           string                                `json:"name"`
	StorageBoxType IDOrName                              `json:"storage_box_type"`
	Location       string                                `json:"location"`
	Labels         *map[string]string                    `json:"labels,omitempty"`
	Password       string                                `json:"password"`
	SSHKeys        []string                              `json:"ssh_keys,omitempty"`
	AccessSettings StorageBoxCreateRequestAccessSettings `json:"access_settings,omitempty"`
}

type StorageBoxCreateRequestAccessSettings struct {
	ReachableExternally *bool `json:"reachable_externally,omitempty"`
	SambaEnabled        *bool `json:"samba_enabled,omitempty"`
	SSHEnabled          *bool `json:"ssh_enabled,omitempty"`
	WebDAVEnabled       *bool `json:"webdav_enabled,omitempty"`
	ZFSEnabled          *bool `json:"zfs_enabled,omitempty"`
}

// StorageBoxCreateResponse defines the schema of the response when
// creating a storage box.
type StorageBoxCreateResponse struct {
	// StorageBox StorageBox `json:"storage_box"` // TODO: Currently not returned from API
	Action Action `json:"action"`
}

// StorageBoxUpdateRequest defines the schema of the request to update a storage box.
type StorageBoxUpdateRequest struct {
	Name   string            `json:"name,omitempty"`
	Labels map[string]string `json:"labels,omitempty"`
}

// StorageBoxUpdateResponse defines the schema of the response when updating a storage box.
type StorageBoxUpdateResponse struct {
	StorageBox StorageBox `json:"storage_box"`
}
