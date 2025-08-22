package hcloud

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/hetznercloud/hcloud-go/v2/hcloud/exp/ctxutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/schema"
)

// StorageBox represents a storage box in Hetzner.
type StorageBox struct {
	ID             int64
	Username       *string
	Status         StorageBoxStatus
	Name           string
	StorageBoxType *StorageBoxType
	Location       *Location
	AccessSettings StorageBoxAccessSettings
	Server         *string
	System         *string
	Stats          *StorageBoxStats
	Labels         map[string]string
	Protection     StorageBoxProtection
	SnapshotPlan   *StorageBoxSnapshotPlan
	Created        time.Time
}

// StorageBoxAccessSettings represents the access settings of a storage box.
type StorageBoxAccessSettings struct {
	ReachableExternally bool
	SambaEnabled        bool
	SSHEnabled          bool
	WebDAVEnabled       bool
	ZFSEnabled          bool
}

// StorageBoxStats represents the disk usage statistics of a storage box.
type StorageBoxStats struct {
	Size          uint64
	SizeData      uint64
	SizeSnapshots uint64
}

// StorageBoxProtection represents the protection level of a storage box.
type StorageBoxProtection struct {
	Delete bool
}

// StorageBoxSnapshotPlan represents the snapshot plan of a storage box.
type StorageBoxSnapshotPlan struct {
	MaxSnapshots int
	Minute       *int
	Hour         *int
	DayOfWeek    *int
	DayOfMonth   *int
}

// StorageBoxStatus specifies a storage box's status.
type StorageBoxStatus string

const (
	// StorageBoxStatusActive is the status when a storage box is active.
	StorageBoxStatusActive StorageBoxStatus = "active"

	// StorageBoxStatusInitializing is the status when a storage box is initializing.
	StorageBoxStatusInitializing StorageBoxStatus = "initializing"

	// StorageBoxStatusLocked is the status when a storage box is locked.
	StorageBoxStatusLocked StorageBoxStatus = "locked"
)

// StorageBoxClient is a client for the storage box API.
type StorageBoxClient struct {
	client *Client
	Action *ResourceActionClient
}

// GetByID retrieves a storage box by its ID. If the storage box does not exist, nil is returned.
func (c *StorageBoxClient) GetByID(ctx context.Context, id int64) (*StorageBox, *Response, error) {
	const opPath = "/storage_boxes/%d"
	ctx = ctxutil.SetOpPath(ctx, opPath)

	reqPath := fmt.Sprintf(opPath, id)

	respBody, resp, err := getRequest[schema.StorageBoxGetResponse](ctx, c.client, reqPath)
	if err != nil {
		if IsError(err, ErrorCodeNotFound) {
			return nil, resp, nil
		}
		return nil, resp, err
	}

	return StorageBoxFromSchema(respBody.StorageBox), resp, nil
}

// GetByName retrieves a storage box by its name. If the storage box does not exist, nil is returned.
func (c *StorageBoxClient) GetByName(ctx context.Context, name string) (*StorageBox, *Response, error) {
	return firstByName(name, func() ([]*StorageBox, *Response, error) {
		return c.List(ctx, StorageBoxListOpts{Name: name})
	})
}

// Get retrieves a storage box by its ID if the input can be parsed as an integer, otherwise it
// retrieves a storage box by its name. If the storage box does not exist, nil is returned.
func (c *StorageBoxClient) Get(ctx context.Context, idOrName string) (*StorageBox, *Response, error) {
	return getByIDOrName(ctx, c.GetByID, c.GetByName, idOrName)
}

// StorageBoxListOpts specifies options for listing storage boxes.
type StorageBoxListOpts struct {
	ListOpts
	Name string
}

func (l StorageBoxListOpts) values() url.Values {
	vals := l.ListOpts.Values()
	if l.Name != "" {
		vals.Add("name", l.Name)
	}
	return vals
}

// List returns a list of storage boxes for a specific page.
//
// Please note that filters specified in opts are not taken into account
// when their value corresponds to their zero value or when they are empty.
func (c *StorageBoxClient) List(ctx context.Context, opts StorageBoxListOpts) ([]*StorageBox, *Response, error) {
	const opPath = "/storage_boxes?%s"
	ctx = ctxutil.SetOpPath(ctx, opPath)

	reqPath := fmt.Sprintf(opPath, opts.values().Encode())

	respBody, resp, err := getRequest[schema.StorageBoxListResponse](ctx, c.client, reqPath)
	if err != nil {
		return nil, resp, err
	}

	return allFromSchemaFunc(respBody.StorageBoxes, StorageBoxFromSchema), resp, nil
}

// All returns all storage boxes.
func (c *StorageBoxClient) All(ctx context.Context) ([]*StorageBox, error) {
	return c.AllWithOpts(ctx, StorageBoxListOpts{ListOpts: ListOpts{PerPage: 50}})
}

// AllWithOpts returns all storage boxes with the given options.
func (c *StorageBoxClient) AllWithOpts(ctx context.Context, opts StorageBoxListOpts) ([]*StorageBox, error) {
	return iterPages(func(page int) ([]*StorageBox, *Response, error) {
		opts.Page = page
		return c.List(ctx, opts)
	})
}

// StorageBoxCreateOpts specifies parameters for creating a storage box.
type StorageBoxCreateOpts struct {
	Name           string
	StorageBoxType *StorageBoxType
	Location       *Location
	Labels         map[string]string
	Password       string
	// TODO: API only accepts literal public keys
	// - should we accept our structs here and read the public key?
	// - What to do when only the ID is set and not the public key?
	// - Should we tell users that this is not required and they can just create an object with the public key field set?
	SSHKeys        []string
	AccessSettings *StorageBoxCreateOptsAccessSettings
}

// StorageBoxCreateOptsAccessSettings specifies access settings for creating a storage box.
type StorageBoxCreateOptsAccessSettings struct {
	ReachableExternally *bool
	SambaEnabled        *bool
	SSHEnabled          *bool
	WebDAVEnabled       *bool
	ZFSEnabled          *bool
}

// StorageBoxCreateResult is the result of a create storage box operation.
type StorageBoxCreateResult struct {
	StorageBox *StorageBox
	Action     *Action
}

// Create creates a new storage box with the given options.
func (c *StorageBoxClient) Create(ctx context.Context, opts StorageBoxCreateOpts) (StorageBoxCreateResult, *Response, error) {
	const opPath = "/storage_boxes"
	ctx = ctxutil.SetOpPath(ctx, opPath)

	reqBody := SchemaFromStorageBoxCreateOpts(opts)

	result := StorageBoxCreateResult{}

	respBody, resp, err := postRequest[schema.StorageBoxCreateResponse](ctx, c.client, opPath, reqBody)
	if err != nil {
		return result, resp, err
	}

	result.StorageBox = StorageBoxFromSchema(respBody.StorageBox)
	result.Action = ActionFromSchema(respBody.Action)

	return result, resp, nil
}

// StorageBoxUpdateOpts specifies options for updating a storage box.
type StorageBoxUpdateOpts struct {
	Name   string
	Labels map[string]string
}

// Update updates a storage box with the given options.
func (c *StorageBoxClient) Update(ctx context.Context, storageBox *StorageBox, opts StorageBoxUpdateOpts) (*StorageBox, *Response, error) {
	const opPath = "/storage_boxes/%d"
	ctx = ctxutil.SetOpPath(ctx, opPath)

	reqPath := fmt.Sprintf(opPath, storageBox.ID)
	reqBody := SchemaFromStorageBoxUpdateOpts(opts)

	respBody, resp, err := putRequest[schema.StorageBoxUpdateResponse](ctx, c.client, reqPath, reqBody)
	if err != nil {
		return nil, resp, err
	}

	return StorageBoxFromSchema(respBody.StorageBox), resp, nil
}

// Delete deletes a storage box.
func (c *StorageBoxClient) Delete(ctx context.Context, storageBox *StorageBox) (*Action, *Response, error) {
	const opPath = "/storage_boxes/%d"
	ctx = ctxutil.SetOpPath(ctx, opPath)

	reqPath := fmt.Sprintf(opPath, storageBox.ID)

	respBody, resp, err := deleteRequest[schema.ActionGetResponse](ctx, c.client, reqPath)
	if err != nil {
		return nil, resp, err
	}
	return ActionFromSchema(respBody.Action), resp, nil
}

type StorageBoxFoldersResult struct {
	Folders []string
}

type StorageBoxFoldersOpts struct {
	Path string
}

func (o StorageBoxFoldersOpts) values() url.Values {
	vals := url.Values{}
	if o.Path != "" {
		vals.Add("path", o.Path)
	}
	return vals
}

// Folders lists folders in a storage box.
func (c *StorageBoxClient) Folders(ctx context.Context, storageBox *StorageBox, opts StorageBoxFoldersOpts) (StorageBoxFoldersResult, *Response, error) {
	const opPath = "/storage_boxes/%d/folders"
	ctx = ctxutil.SetOpPath(ctx, opPath)

	reqPath := fmt.Sprintf(opPath, storageBox.ID)

	if len(opts.values()) > 0 {
		reqPath += "?" + opts.values().Encode()
	}

	result := StorageBoxFoldersResult{}

	respBody, resp, err := getRequest[schema.StorageBoxFoldersResponse](ctx, c.client, reqPath)
	if err != nil {
		return result, resp, err
	}

	result.Folders = respBody.Folders

	return result, resp, nil
}

// StorageBoxChangeProtectionOpts specifies options for changing the protection level of a storage box.
type StorageBoxChangeProtectionOpts struct {
	Delete bool
}

// ChangeProtection changes the protection level of a storage box.
func (c *StorageBoxClient) ChangeProtection(ctx context.Context, storageBox *StorageBox, opts StorageBoxChangeProtectionOpts) (*Action, *Response, error) {
	const opPath = "/storage_boxes/%d/change_protection"
	ctx = ctxutil.SetOpPath(ctx, opPath)

	reqPath := fmt.Sprintf(opPath, storageBox.ID)
	reqBody := SchemaFromStorageBoxChangeProtectionOpts(opts)

	respBody, resp, err := postRequest[schema.ActionGetResponse](ctx, c.client, reqPath, reqBody)
	if err != nil {
		return nil, resp, err
	}

	action := ActionFromSchema(respBody.Action)

	return action, resp, nil
}

// StorageBoxChangeTypeOpts specifies options for changing the type of a storage box.
type StorageBoxChangeTypeOpts struct {
	StorageBoxType *StorageBoxType
}

// ChangeType changes the type of a storage box.
func (c *StorageBoxClient) ChangeType(ctx context.Context, storageBox *StorageBox, opts StorageBoxChangeTypeOpts) (*Action, *Response, error) {
	const opPath = "/storage_boxes/%d/change_type"
	ctx = ctxutil.SetOpPath(ctx, opPath)

	reqPath := fmt.Sprintf(opPath, storageBox.ID)
	reqBody := SchemaFromStorageBoxChangeTypeOpts(opts)

	respBody, resp, err := postRequest[schema.ActionGetResponse](ctx, c.client, reqPath, reqBody)
	if err != nil {
		return nil, resp, err
	}

	action := ActionFromSchema(respBody.Action)

	return action, resp, nil
}

// StorageBoxResetPasswordOpts specifies options for resetting the password of a storage box.
type StorageBoxResetPasswordOpts struct {
	Password string
}

// ResetPassword resets the password of a storage box.
func (c *StorageBoxClient) ResetPassword(ctx context.Context, storageBox *StorageBox, opts StorageBoxResetPasswordOpts) (*Action, *Response, error) {
	const opPath = "/storage_boxes/%d/reset_password"
	ctx = ctxutil.SetOpPath(ctx, opPath)

	reqPath := fmt.Sprintf(opPath, storageBox.ID)
	reqBody := SchemaFromStorageBoxResetPasswordOpts(opts)

	respBody, resp, err := postRequest[schema.ActionGetResponse](ctx, c.client, reqPath, reqBody)
	if err != nil {
		return nil, resp, err
	}

	action := ActionFromSchema(respBody.Action)

	return action, resp, nil
}

// StorageBoxUpdateAccessSettingsOpts specifies options for updating the access settings of a storage box.
type StorageBoxUpdateAccessSettingsOpts struct {
	SambaEnabled        *bool
	SSHEnabled          *bool
	WebDAVEnabled       *bool
	ZFSEnabled          *bool
	ReachableExternally *bool
}

// UpdateAccessSettings updates the access settings of a storage box.
func (c *StorageBoxClient) UpdateAccessSettings(ctx context.Context, storageBox *StorageBox, opts StorageBoxUpdateAccessSettingsOpts) (*Action, *Response, error) {
	const opPath = "/storage_boxes/%d/update_access_settings"
	ctx = ctxutil.SetOpPath(ctx, opPath)

	reqPath := fmt.Sprintf(opPath, storageBox.ID)
	reqBody := SchemaFromStorageBoxUpdateAccessSettingsOpts(opts)

	respBody, resp, err := postRequest[schema.ActionGetResponse](ctx, c.client, reqPath, reqBody)
	if err != nil {
		return nil, resp, err
	}

	action := ActionFromSchema(respBody.Action)

	return action, resp, nil
}

// StorageBoxRollbackSnapshotOpts specifies options for rolling back a storage box to a snapshot.
type StorageBoxRollbackSnapshotOpts struct {
	Snapshot *StorageBoxSnapshot
}

// RollbackSnapshot rolls back a storage box to a snapshot.
func (c *StorageBoxClient) RollbackSnapshot(
	ctx context.Context,
	storageBox *StorageBox,
	opts StorageBoxRollbackSnapshotOpts,
) (*Action, *Response, error) {
	const opPath = "/storage_boxes/%d/actions/rollback_snapshot"
	ctx = ctxutil.SetOpPath(ctx, opPath)

	reqPath := fmt.Sprintf(opPath, storageBox.ID)
	reqBody := SchemaFromStorageBoxRollbackSnapshotOpts(opts)

	respBody, resp, err := postRequest[schema.ActionGetResponse](ctx, c.client, reqPath, reqBody)
	if err != nil {
		return nil, resp, err
	}

	return ActionFromSchema(respBody.Action), resp, nil
}

// StorageBoxEnableSnapshotPlanOpts specifies options for enabling a snapshot plan for a storage box.
type StorageBoxEnableSnapshotPlanOpts struct {
	MaxSnapshots int
	Minute       *int // Null means every minute.
	Hour         *int // Null means every hour.
	DayOfWeek    *int // Null means every day. // TODO: Starts with 1 for Monday and goes to 7 for Sunday
	DayOfMonth   *int // Null means every day.
}

// EnableSnapshotPlan enables a snapshot plan for a storage box.
func (c *StorageBoxClient) EnableSnapshotPlan(
	ctx context.Context,
	storageBox *StorageBox,
	opts StorageBoxEnableSnapshotPlanOpts,
) (*Action, *Response, error) {
	const opPath = "/storage_boxes/%d/actions/enable_snapshot_plan"
	ctx = ctxutil.SetOpPath(ctx, opPath)

	reqPath := fmt.Sprintf(opPath, storageBox.ID)
	reqBody := SchemaFromStorageBoxEnableSnapshotPlan(opts)

	respBody, resp, err := postRequest[schema.ActionGetResponse](ctx, c.client, reqPath, reqBody)
	if err != nil {
		return nil, resp, err
	}

	return ActionFromSchema(respBody.Action), resp, nil
}

// DisableSnapshotPlan disables the snapshot plan for a storage box.
func (c *StorageBoxClient) DisableSnapshotPlan(
	ctx context.Context,
	storageBox *StorageBox,
) (*Action, *Response, error) {
	const opPath = "/storage_boxes/%d/actions/disable_snapshot_plan"
	ctx = ctxutil.SetOpPath(ctx, opPath)

	reqPath := fmt.Sprintf(opPath, storageBox.ID)

	respBody, resp, err := postRequest[schema.ActionGetResponse](ctx, c.client, reqPath, nil)
	if err != nil {
		return nil, resp, err
	}

	return ActionFromSchema(respBody.Action), resp, nil
}
