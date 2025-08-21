package schema

import "time"

type StorageBoxSnapshot struct {
	ID          int64                   `json:"id"`
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Stats       StorageBoxSnapshotStats `json:"stats"`
	IsAutomatic bool                    `json:"is_automatic"`
	Labels      map[string]string       `json:"labels"`
	Created     time.Time               `json:"created"`
	StorageBox  int64                   `json:"storage_box"`
}

type StorageBoxSnapshotStats struct {
	Size           int64 `json:"size"`
	SizeFilesystem int64 `json:"size_filesystem"`
}

type StorageBoxSnapshotGetResponse struct {
	Snapshot StorageBoxSnapshot `json:"snapshot"`
}

type StorageBoxSnapshotListResponse struct {
	Snapshots []StorageBoxSnapshot `json:"snapshots"`
}

type StorageBoxSnapshotCreateRequest struct {
	Description string `json:"description"`
}

type StorageBoxSnapshotCreateResponse struct {
	Snapshot StorageBoxSnapshot `json:"snapshot"`
	Action   Action             `json:"action"`
}

// TODO: Both are required by the spec? Why?
type StorageBoxSnapshotUpdateRequest struct {
	Description string            `json:"description"`
	Labels      map[string]string `json:"labels"`
}

type StorageBoxSnapshotUpdateResponse struct {
	Snapshot StorageBoxSnapshot `json:"snapshot"`
}
