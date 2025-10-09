package hcloud

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/hetznercloud/hcloud-go/v2/hcloud/exp/ctxutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/schema"
)

// StorageBoxSnapshot represents a snapshot of a Storage Box.
type StorageBoxSnapshot struct {
	ID          int64
	Name        string
	Description *string
	Stats       *StorageBoxSnapshotStats
	IsAutomatic bool
	Labels      map[string]string
	Created     time.Time
	StorageBox  *StorageBox
}

// StorageBoxSnapshotStats represents the size of a Storage Box snapshot.
type StorageBoxSnapshotStats struct {
	Size           uint64
	SizeFilesystem uint64
}

// GetSnapshotByID gets a Storage Box Snapshot by its ID.
func (c *StorageBoxClient) GetSnapshotByID(ctx context.Context, storageBox *StorageBox, id int64) (*StorageBoxSnapshot, *Response, error) {
	const optPath = "/storage_boxes/%d/snapshots/%d"
	ctx = ctxutil.SetOpPath(ctx, optPath)

	reqPath := fmt.Sprintf(optPath, storageBox.ID, id)

	respBody, resp, err := getRequest[schema.StorageBoxSnapshotGetResponse](ctx, c.client, reqPath)
	if err != nil {
		if IsError(err, ErrorCodeNotFound) {
			return nil, resp, nil
		}
		return nil, resp, err
	}

	return StorageBoxSnapshotFromSchema(respBody.Snapshot), resp, nil
}

// GetSnapshotByName gets a Storage Box snapshot by its name.
func (c *StorageBoxClient) GetSnapshotByName(
	ctx context.Context,
	storageBox *StorageBox,
	name string,
) (*StorageBoxSnapshot, *Response, error) {
	return firstByName(name, func() ([]*StorageBoxSnapshot, *Response, error) {
		return c.ListSnapshots(ctx, storageBox, StorageBoxSnapshotListOpts{Name: name})
	})
}

// GetSnapshot gets a Storage Box snapshot by its ID or name.
func (c *StorageBoxClient) GetSnapshot(
	ctx context.Context,
	storageBox *StorageBox,
	idOrName string,
) (*StorageBoxSnapshot, *Response, error) {
	return getByIDOrName(
		ctx,
		func(ctx context.Context, id int64) (*StorageBoxSnapshot, *Response, error) {
			return c.GetSnapshotByID(ctx, storageBox, id)
		},
		func(ctx context.Context, name string) (*StorageBoxSnapshot, *Response, error) {
			return c.GetSnapshotByName(ctx, storageBox, name)
		},
		idOrName,
	)
}

// StorageBoxSnapshotListOpts specifies options for listing Storage Box snapshots.
type StorageBoxSnapshotListOpts struct {
	LabelSelector string
	Name          string
	IsAutomatic   *bool
	Sort          []string
}

func (o StorageBoxSnapshotListOpts) values() url.Values {
	vals := url.Values{}
	if o.LabelSelector != "" {
		vals.Set("label_selector", o.LabelSelector)
	}
	if o.Name != "" {
		vals.Set("name", o.Name)
	}
	if o.IsAutomatic != nil {
		vals.Set("is_automatic", fmt.Sprintf("%t", *o.IsAutomatic))
	}
	for _, sort := range o.Sort {
		vals.Add("sort", sort)
	}
	return vals
}

// ListSnapshots lists all snapshots of a Storage Box with the given options.
//
// Pagination is not supported, so this will return all snapshots at once.
func (c *StorageBoxClient) ListSnapshots(
	ctx context.Context,
	storageBox *StorageBox,
	opts StorageBoxSnapshotListOpts,
) ([]*StorageBoxSnapshot, *Response, error) {
	const optPath = "/storage_boxes/%d/snapshots?%s"
	ctx = ctxutil.SetOpPath(ctx, optPath)

	reqPath := fmt.Sprintf(optPath, storageBox.ID, opts.values().Encode())

	respBody, resp, err := getRequest[schema.StorageBoxSnapshotListResponse](ctx, c.client, reqPath)
	if err != nil {
		return nil, resp, err
	}

	return allFromSchemaFunc(respBody.Snapshots, StorageBoxSnapshotFromSchema), resp, nil
}

// AllSnapshotsWithOpts lists all snapshots of a Storage Box with the given options.
func (c *StorageBoxClient) AllSnapshotsWithOpts(
	ctx context.Context,
	storageBox *StorageBox,
	opts StorageBoxSnapshotListOpts,
) ([]*StorageBoxSnapshot, error) {
	snapshots, _, err := c.ListSnapshots(ctx, storageBox, opts)
	return snapshots, err
}

// AllSnapshots lists all snapshots of a Storage Box without any options.
func (c *StorageBoxClient) AllSnapshots(
	ctx context.Context,
	storageBox *StorageBox,
) ([]*StorageBoxSnapshot, error) {
	opts := StorageBoxSnapshotListOpts{}
	snapshots, _, err := c.ListSnapshots(ctx, storageBox, opts)
	return snapshots, err
}

// StorageBoxSnapshotCreateOpts specifies options for creating a Storage Box snapshot.
type StorageBoxSnapshotCreateOpts struct {
	Description *string
}

// StorageBoxSnapshotCreateResult represents the result of creating a Storage Box snapshot.
type StorageBoxSnapshotCreateResult struct {
	Snapshot *StorageBoxSnapshot
	Action   *Action
}

// CreateSnapshot creates a new snapshot for the given Storage Box with the provided options.
func (c *StorageBoxClient) CreateSnapshot(
	ctx context.Context,
	storageBox *StorageBox,
	opts StorageBoxSnapshotCreateOpts,
) (StorageBoxSnapshotCreateResult, *Response, error) {
	const opPath = "/storage_boxes/%d/snapshots"
	ctx = ctxutil.SetOpPath(ctx, opPath)

	reqPath := fmt.Sprintf(opPath, storageBox.ID)
	reqBody := SchemaFromStorageBoxSnapshotCreateOpts(opts)

	result := StorageBoxSnapshotCreateResult{}

	respBody, resp, err := postRequest[schema.StorageBoxSnapshotCreateResponse](ctx, c.client, reqPath, reqBody)
	if err != nil {
		return result, resp, err
	}

	result.Snapshot = StorageBoxSnapshotFromSchema(respBody.Snapshot)
	result.Action = ActionFromSchema(respBody.Action)

	return result, resp, err
}

// StorageBoxSnapshotUpdateOpts specifies options for updating a Storage Box snapshot.
type StorageBoxSnapshotUpdateOpts struct {
	Description *string
	Labels      map[string]string
}

// UpdateSnapshot updates the given snapshot of a Storage Box with the provided options.
func (c *StorageBoxClient) UpdateSnapshot(
	ctx context.Context,
	snapshot *StorageBoxSnapshot,
	opts StorageBoxSnapshotUpdateOpts,
) (*StorageBoxSnapshot, *Response, error) {
	const opPath = "/storage_boxes/%d/snapshots/%d"
	ctx = ctxutil.SetOpPath(ctx, opPath)

	reqPath := fmt.Sprintf(opPath, snapshot.StorageBox.ID, snapshot.ID)
	reqBody := SchemaFromStorageBoxSnapshotUpdateOpts(opts)

	respBody, resp, err := putRequest[schema.StorageBoxSnapshotUpdateResponse](ctx, c.client, reqPath, reqBody)
	if err != nil {
		return nil, resp, err
	}

	updatedSnapshot := StorageBoxSnapshotFromSchema(respBody.Snapshot)

	return updatedSnapshot, resp, nil
}

// DeleteSnapshot deletes the given snapshot of a Storage Box.
func (c *StorageBoxClient) DeleteSnapshot(
	ctx context.Context,
	snapshot *StorageBoxSnapshot,
) (*Action, *Response, error) {
	const opPath = "/storage_boxes/%d/snapshots/%d"
	ctx = ctxutil.SetOpPath(ctx, opPath)

	reqPath := fmt.Sprintf(opPath, snapshot.StorageBox.ID, snapshot.ID)

	respBody, resp, err := deleteRequest[schema.ActionGetResponse](ctx, c.client, reqPath)

	action := ActionFromSchema(respBody.Action)

	return action, resp, err
}
