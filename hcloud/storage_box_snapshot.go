package hcloud

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/hetznercloud/hcloud-go/v2/hcloud/exp/ctxutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/schema"
)

type StorageBoxSnapshot struct {
	ID          int64
	Name        string
	Description string
	Stats       *StorageBoxSnapshotStats
	IsAutomatic bool
	Labels      map[string]string
	Created     time.Time
	StorageBox  *StorageBox
}

type StorageBoxSnapshotStats struct {
	Size           int64
	SizeFilesystem int64
}

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

func (c *StorageBoxClient) GetSnapshotByName(
	ctx context.Context,
	storageBox *StorageBox,
	name string,
) (*StorageBoxSnapshot, *Response, error) {
	return firstByName(name, func() ([]*StorageBoxSnapshot, *Response, error) {
		return c.ListSnapshots(ctx, storageBox, StorageBoxSnapshotListOpts{Name: name})
	})
}

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

type StorageBoxSnapshotListOpts struct {
	LabelSelector string
	Name          string
}

func (o StorageBoxSnapshotListOpts) values() url.Values {
	values := url.Values{}
	if o.LabelSelector != "" {
		values.Set("label_selector", o.LabelSelector)
	}
	if o.Name != "" {
		values.Set("name", o.Name)
	}
	return values
}

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

// Implementing this as an alias as other integrations depend on this.
func (c *StorageBoxClient) AllSnapshotsWithOpts(
	ctx context.Context,
	storageBox *StorageBox,
	opts StorageBoxSnapshotListOpts,
) ([]*StorageBoxSnapshot, error) {
	snapshots, _, err := c.ListSnapshots(ctx, storageBox, opts)
	if err != nil {
		return nil, err
	}
	return snapshots, nil
}

// Implementing this as an alias as other integrations depend on this.
func (c *StorageBoxClient) AllSnapshots(
	ctx context.Context,
	storageBox *StorageBox,
) ([]*StorageBoxSnapshot, error) {
	opts := StorageBoxSnapshotListOpts{}
	snapshots, _, err := c.ListSnapshots(ctx, storageBox, opts)
	if err != nil {
		return nil, err
	}
	return snapshots, nil
}

type StorageBoxSnapshotCreateOpts struct {
	Description string
}

type StorageBoxSnapshotCreateResult struct {
	Snapshot *StorageBoxSnapshot
	Action   *Action
}

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

type StorageBoxSnapshotUpdateOpts struct {
	Description string
	Labels      map[string]string
}

func (c *StorageBoxClient) UpdateSnapshot(
	ctx context.Context,
	storageBox *StorageBox,
	snapshot *StorageBoxSnapshot,
	opts StorageBoxSnapshotUpdateOpts,
) (*StorageBoxSnapshot, *Response, error) {
	const opPath = "/storage_boxes/%d/snapshots/%d"
	ctx = ctxutil.SetOpPath(ctx, opPath)

	reqPath := fmt.Sprintf(opPath, storageBox.ID, snapshot.ID)
	reqBody := SchemaFromStorageBoxSnapshotUpdateOpts(opts)

	respBody, resp, err := putRequest[schema.StorageBoxSnapshotUpdateResponse](ctx, c.client, reqPath, reqBody)
	if err != nil {
		return nil, resp, err
	}

	updatedSnapshot := StorageBoxSnapshotFromSchema(respBody.Snapshot)

	return updatedSnapshot, resp, nil
}

func (c *StorageBoxClient) DeleteSnapshot(
	ctx context.Context,
	storageBox *StorageBox,
	snapshot *StorageBoxSnapshot,
) (*Action, *Response, error) {
	const opPath = "/storage_boxes/%d/snapshots/%d"
	ctx = ctxutil.SetOpPath(ctx, opPath)

	reqPath := fmt.Sprintf(opPath, storageBox.ID, snapshot.ID)

	respBody, resp, err := deleteRequest[schema.ActionGetResponse](ctx, c.client, reqPath)

	action := ActionFromSchema(respBody.Action)

	return action, resp, err
}
