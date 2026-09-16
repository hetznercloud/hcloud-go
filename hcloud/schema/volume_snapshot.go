package schema

import "time"

// VolumeSnapshot defines the proposed schema of a block volume snapshot.
type VolumeSnapshot struct {
	ID          int64             `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Status      string            `json:"status"`
	Volume      int64             `json:"volume"`
	Location    Location          `json:"location"`
	Size        int               `json:"size"`
	Labels      map[string]string `json:"labels"`
	Created     time.Time         `json:"created"`
}

// VolumeSnapshotGetResponse defines the response when retrieving a snapshot.
type VolumeSnapshotGetResponse struct {
	Snapshot VolumeSnapshot `json:"snapshot"`
}

// VolumeSnapshotListResponse defines the response when listing snapshots.
type VolumeSnapshotListResponse struct {
	Snapshots []VolumeSnapshot `json:"snapshots"`
}

// VolumeSnapshotCreateRequest defines the request for a snapshot of a volume.
type VolumeSnapshotCreateRequest struct {
	Name        string             `json:"name"`
	Description string             `json:"description,omitempty"`
	Labels      *map[string]string `json:"labels,omitempty"`
}

// VolumeSnapshotCreateResponse defines the response when creating a snapshot.
type VolumeSnapshotCreateResponse struct {
	Snapshot VolumeSnapshot `json:"snapshot"`
	Action   Action         `json:"action"`
}

// VolumeSnapshotDeleteResponse defines the response when deleting a snapshot.
type VolumeSnapshotDeleteResponse struct {
	Action Action `json:"action"`
}
