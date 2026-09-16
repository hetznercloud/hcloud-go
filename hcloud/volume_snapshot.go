package hcloud

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/hetznercloud/hcloud-go/v2/hcloud/exp/ctxutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/schema"
)

// VolumeSnapshot represents a point-in-time snapshot of a block [Volume].
//
// This type is an API proposal. The corresponding Cloud API endpoints do not
// exist publicly yet.
type VolumeSnapshot struct {
	ID          int64
	Name        string
	Description string
	Status      VolumeSnapshotStatus
	Volume      *Volume
	Location    *Location
	Size        int
	Labels      map[string]string
	Created     time.Time
}

func (o *VolumeSnapshot) pathID() (string, error) {
	if o.ID == 0 {
		return "", missingField(o, "ID")
	}
	return strconv.FormatInt(o.ID, 10), nil
}

// VolumeSnapshotStatus is the lifecycle status of a volume snapshot.
type VolumeSnapshotStatus string

const (
	VolumeSnapshotStatusCreating  VolumeSnapshotStatus = "creating"
	VolumeSnapshotStatusAvailable VolumeSnapshotStatus = "available"
)

// GetSnapshotByID returns a volume snapshot by ID. A missing snapshot returns nil.
func (c *VolumeClient) GetSnapshotByID(ctx context.Context, id int64) (*VolumeSnapshot, *Response, error) {
	const opPath = "/volume_snapshots/%d"
	ctx = ctxutil.SetOpPath(ctx, opPath)

	respBody, resp, err := getRequest[schema.VolumeSnapshotGetResponse](ctx, c.client, fmt.Sprintf(opPath, id))
	if err != nil {
		if IsError(err, ErrorCodeNotFound) {
			return nil, resp, nil
		}
		return nil, resp, err
	}
	return VolumeSnapshotFromSchema(respBody.Snapshot), resp, nil
}

// GetSnapshotByName returns the first volume snapshot with name. A missing snapshot returns nil.
func (c *VolumeClient) GetSnapshotByName(ctx context.Context, name string) (*VolumeSnapshot, *Response, error) {
	return firstByName(name, func() ([]*VolumeSnapshot, *Response, error) {
		return c.ListSnapshots(ctx, VolumeSnapshotListOpts{Name: name})
	})
}

// VolumeSnapshotListOpts specifies filters for listing volume snapshots.
type VolumeSnapshotListOpts struct {
	ListOpts
	Name          string
	Volume        *Volume
	Status        []VolumeSnapshotStatus
	LabelSelector string
	Sort          []string
}

func (o VolumeSnapshotListOpts) Values() url.Values {
	values := o.ListOpts.Values()
	if o.Name != "" {
		values.Set("name", o.Name)
	}
	if o.Volume != nil {
		values.Set("volume", strconv.FormatInt(o.Volume.ID, 10))
	}
	for _, status := range o.Status {
		values.Add("status", string(status))
	}
	if o.LabelSelector != "" {
		values.Set("label_selector", o.LabelSelector)
	}
	for _, sort := range o.Sort {
		values.Add("sort", sort)
	}
	return values
}

// ListSnapshots returns one page of block volume snapshots.
func (c *VolumeClient) ListSnapshots(ctx context.Context, opts VolumeSnapshotListOpts) ([]*VolumeSnapshot, *Response, error) {
	const opPath = "/volume_snapshots?%s"
	ctx = ctxutil.SetOpPath(ctx, opPath)

	respBody, resp, err := getRequest[schema.VolumeSnapshotListResponse](ctx, c.client, fmt.Sprintf(opPath, opts.Values().Encode()))
	if err != nil {
		return nil, resp, err
	}
	return allFromSchemaFunc(respBody.Snapshots, VolumeSnapshotFromSchema), resp, nil
}

// AllSnapshots returns all block volume snapshots.
func (c *VolumeClient) AllSnapshots(ctx context.Context) ([]*VolumeSnapshot, error) {
	return c.AllSnapshotsWithOpts(ctx, VolumeSnapshotListOpts{})
}

// AllSnapshotsWithOpts returns all block volume snapshots matching opts.
func (c *VolumeClient) AllSnapshotsWithOpts(ctx context.Context, opts VolumeSnapshotListOpts) ([]*VolumeSnapshot, error) {
	if opts.ListOpts.PerPage == 0 {
		opts.ListOpts.PerPage = 50
	}
	return iterPages(func(page int) ([]*VolumeSnapshot, *Response, error) {
		opts.Page = page
		return c.ListSnapshots(ctx, opts)
	})
}

// VolumeSnapshotCreateOpts specifies a point-in-time snapshot request.
type VolumeSnapshotCreateOpts struct {
	Name        string
	Description string
	Labels      map[string]string
}

// Validate checks if snapshot creation options are valid.
func (o VolumeSnapshotCreateOpts) Validate() error {
	if o.Name == "" {
		return missingField(o, "Name")
	}
	return nil
}

// VolumeSnapshotCreateResult is the result of creating a snapshot.
type VolumeSnapshotCreateResult struct {
	Snapshot *VolumeSnapshot
	Action   *Action
}

// CreateSnapshot creates a point-in-time snapshot of volume.
func (c *VolumeClient) CreateSnapshot(ctx context.Context, volume *Volume, opts VolumeSnapshotCreateOpts) (VolumeSnapshotCreateResult, *Response, error) {
	const opPath = "/volumes/%d/actions/create_snapshot"
	ctx = ctxutil.SetOpPath(ctx, opPath)

	result := VolumeSnapshotCreateResult{}
	if volume == nil || volume.ID == 0 {
		return result, nil, missingField(opts, "Volume")
	}
	if err := opts.Validate(); err != nil {
		return result, nil, err
	}
	reqBody := schema.VolumeSnapshotCreateRequest{
		Name:        opts.Name,
		Description: opts.Description,
	}
	if opts.Labels != nil {
		reqBody.Labels = &opts.Labels
	}
	respBody, resp, err := postRequest[schema.VolumeSnapshotCreateResponse](ctx, c.client, fmt.Sprintf(opPath, volume.ID), reqBody)
	if err != nil {
		return result, resp, err
	}
	result.Snapshot = VolumeSnapshotFromSchema(respBody.Snapshot)
	result.Action = ActionFromSchema(respBody.Action)
	return result, resp, nil
}

// DeleteSnapshot deletes a block volume snapshot.
func (c *VolumeClient) DeleteSnapshot(ctx context.Context, snapshot *VolumeSnapshot) (*Action, *Response, error) {
	const opPath = "/volume_snapshots/%d"
	ctx = ctxutil.SetOpPath(ctx, opPath)

	if snapshot == nil || snapshot.ID == 0 {
		return nil, nil, missingField(snapshot, "ID")
	}
	respBody, resp, err := deleteRequest[schema.VolumeSnapshotDeleteResponse](ctx, c.client, fmt.Sprintf(opPath, snapshot.ID))
	if err != nil {
		return nil, resp, err
	}
	return ActionFromSchema(respBody.Action), resp, nil
}

// VolumeSnapshotFromSchema converts a schema snapshot to its public representation.
func VolumeSnapshotFromSchema(s schema.VolumeSnapshot) *VolumeSnapshot {
	return &VolumeSnapshot{
		ID:          s.ID,
		Name:        s.Name,
		Description: s.Description,
		Status:      VolumeSnapshotStatus(s.Status),
		Volume:      &Volume{ID: s.Volume},
		Location:    LocationFromSchema(s.Location),
		Size:        s.Size,
		Labels:      s.Labels,
		Created:     s.Created,
	}
}
