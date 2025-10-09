package hcloud

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/hetznercloud/hcloud-go/v2/hcloud/exp/ctxutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/schema"
)

// StorageBoxSubaccount represents a subaccount of a Storage Box.
type StorageBoxSubaccount struct {
	ID             int64
	Username       string
	HomeDirectory  string
	Server         string
	AccessSettings *StorageBoxSubaccountAccessSettings
	Description    string
	Labels         map[string]string
	Created        time.Time
	StorageBox     *StorageBox
}

// StorageBoxSubaccountAccessSettings represents the access settings of a Storage Box subaccount.
type StorageBoxSubaccountAccessSettings struct {
	ReachableExternally bool
	Readonly            bool
	SambaEnabled        bool
	SSHEnabled          bool
	WebDAVEnabled       bool
}

// GetSubaccountByID retrieves a Storage Box subaccount by its ID.
func (c *StorageBoxClient) GetSubaccountByID(
	ctx context.Context,
	storageBox *StorageBox,
	id int64,
) (*StorageBoxSubaccount, *Response, error) {
	const opPath = "/storage_boxes/%d/subaccounts/%d"
	ctx = ctxutil.SetOpPath(ctx, opPath)

	reqPath := fmt.Sprintf(opPath, storageBox.ID, id)

	respBody, resp, err := getRequest[schema.StorageBoxSubaccountGetResponse](ctx, c.client, reqPath)
	if err != nil {
		if IsError(err, ErrorCodeNotFound) {
			return nil, resp, nil
		}
		return nil, resp, err
	}

	subaccount := StorageBoxSubaccountFromSchema(respBody.Subaccount)

	return subaccount, resp, nil
}

// StorageBoxSubaccountListOpts represents the options for listing Storage Box subaccounts.
type StorageBoxSubaccountListOpts struct {
	LabelSelector string
	Username      string
	Sort          []string
}

func (o StorageBoxSubaccountListOpts) values() url.Values {
	vals := url.Values{}
	if o.Username != "" {
		vals.Add("username", o.Username)
	}
	if len(o.LabelSelector) > 0 {
		vals.Add("label_selector", o.LabelSelector)
	}
	for _, sort := range o.Sort {
		vals.Add("sort", sort)
	}
	return vals
}

// ListSubaccounts lists all subaccounts of a Storage Box.
func (c *StorageBoxClient) ListSubaccounts(
	ctx context.Context,
	storageBox *StorageBox,
	opts StorageBoxSubaccountListOpts,
) ([]*StorageBoxSubaccount, *Response, error) {
	const opPath = "/storage_boxes/%d/subaccounts"
	ctx = ctxutil.SetOpPath(ctx, opPath)

	reqPath := fmt.Sprintf(opPath, storageBox.ID)
	reqPath = fmt.Sprintf("%s?%s", reqPath, opts.values().Encode())

	respBody, resp, err := getRequest[schema.StorageBoxSubaccountListResponse](ctx, c.client, reqPath)
	if err != nil {
		return nil, resp, err
	}

	return allFromSchemaFunc(respBody.Subaccounts, StorageBoxSubaccountFromSchema), resp, nil
}

// AllSubaccountsWithOpts retrieves all subaccounts of a Storage Box with the given options.
func (c *StorageBoxClient) AllSubaccountsWithOpts(
	ctx context.Context,
	storageBox *StorageBox,
	opts StorageBoxSubaccountListOpts,
) ([]*StorageBoxSubaccount, error) {
	subaccounts, _, err := c.ListSubaccounts(ctx, storageBox, opts)
	return subaccounts, err
}

// AllSubaccounts retrieves all subaccounts of a Storage Box.
func (c *StorageBoxClient) AllSubaccounts(
	ctx context.Context,
	storageBox *StorageBox,
) ([]*StorageBoxSubaccount, error) {
	opts := StorageBoxSubaccountListOpts{}
	subaccounts, _, err := c.ListSubaccounts(ctx, storageBox, opts)
	return subaccounts, err
}

// StorageBoxSubaccountCreateOpts represents the options for creating a Storage Box subaccount.
type StorageBoxSubaccountCreateOpts struct {
	Password       string
	HomeDirectory  *string
	AccessSettings *StorageBoxSubaccountCreateOptsAccessSettings
	Description    *string
	Labels         map[string]string
}

// StorageBoxSubaccountAccessSettingsOpts represents the options for access settings of a Storage Box subaccount.
type StorageBoxSubaccountCreateOptsAccessSettings struct {
	ReachableExternally *bool
	Readonly            *bool
	SambaEnabled        *bool
	SSHEnabled          *bool
	WebDAVEnabled       *bool
}

// StorageBoxSubaccountCreateResult represents the result of creating a Storage Box subaccount.
type StorageBoxSubaccountCreateResult struct {
	Subaccount *StorageBoxSubaccount
	Action     *Action
}

// CreateSubaccount creates a new subaccount for a Storage Box.
func (c *StorageBoxClient) CreateSubaccount(
	ctx context.Context,
	storageBox *StorageBox,
	opts StorageBoxSubaccountCreateOpts,
) (StorageBoxSubaccountCreateResult, *Response, error) {
	const opPath = "/storage_boxes/%d/subaccounts"
	ctx = ctxutil.SetOpPath(ctx, opPath)

	reqPath := fmt.Sprintf(opPath, storageBox.ID)
	reqBody := SchemaFromStorageBoxSubaccountCreateOpts(opts)

	result := StorageBoxSubaccountCreateResult{}

	respBody, resp, err := postRequest[schema.StorageBoxSubaccountCreateResponse](ctx, c.client, reqPath, reqBody)
	if err != nil {
		return result, resp, err
	}

	result.Action = ActionFromSchema(respBody.Action)
	result.Subaccount = StorageBoxSubaccountFromCreateResponse(respBody.Subaccount)

	return result, resp, nil
}

// StorageBoxSubaccountUpdateOpts represents the options for updating a Storage Box subaccount.
type StorageBoxSubaccountUpdateOpts struct {
	Description *string
	Labels      map[string]string
}

// UpdateSubaccount updates a subaccount of a Storage Box.
func (c *StorageBoxClient) UpdateSubaccount(
	ctx context.Context,
	subaccount *StorageBoxSubaccount,
	opts StorageBoxSubaccountUpdateOpts,
) (*StorageBoxSubaccount, *Response, error) {
	const opPath = "/storage_boxes/%d/subaccounts/%d"
	ctx = ctxutil.SetOpPath(ctx, opPath)

	reqPath := fmt.Sprintf(opPath, subaccount.StorageBox.ID, subaccount.ID)
	reqBody := SchemaFromStorageBoxSubaccountUpdateOpts(opts)

	respBody, resp, err := putRequest[schema.StorageBoxSubaccountUpdateResponse](ctx, c.client, reqPath, reqBody)
	if err != nil {
		return nil, resp, err
	}

	updatedSubaccount := StorageBoxSubaccountFromSchema(respBody.Subaccount)

	return updatedSubaccount, resp, nil
}

// DeleteSubaccount deletes a subaccount from a Storage Box.
func (c *StorageBoxClient) DeleteSubaccount(
	ctx context.Context,
	subaccount *StorageBoxSubaccount,
) (*Action, *Response, error) {
	const opPath = "/storage_boxes/%d/subaccounts/%d"
	ctx = ctxutil.SetOpPath(ctx, opPath)

	reqPath := fmt.Sprintf(opPath, subaccount.StorageBox.ID, subaccount.ID)

	respBody, resp, err := deleteRequest[schema.ActionGetResponse](ctx, c.client, reqPath)
	if err != nil {
		return nil, resp, err
	}

	action := ActionFromSchema(respBody.Action)

	return action, resp, nil
}

// StorageBoxSubaccountResetPasswordOpts represents the options for resetting a Storage Box subaccount's password.
type StorageBoxSubaccountResetPasswordOpts struct {
	Password string
}

// ResetSubaccountPassword resets the password of a Storage Box subaccount.
func (c *StorageBoxClient) ResetSubaccountPassword(
	ctx context.Context,
	subaccount *StorageBoxSubaccount,
	opts StorageBoxSubaccountResetPasswordOpts,
) (*Action, *Response, error) {
	const opPath = "/storage_boxes/%d/subaccounts/%d/actions/reset_subaccount_password"
	ctx = ctxutil.SetOpPath(ctx, opPath)

	reqPath := fmt.Sprintf(opPath, subaccount.StorageBox.ID, subaccount.ID)
	reqBody := SchemaFromStorageBoxSubaccountResetPasswordOpts(opts)

	respBody, resp, err := postRequest[schema.ActionGetResponse](ctx, c.client, reqPath, reqBody)
	if err != nil {
		return nil, resp, err
	}

	return ActionFromSchema(respBody.Action), resp, err
}

// StorageBoxSubaccountAccessSettingsUpdateOpts represents the options for updating access settings of a Storage Box subaccount.
type StorageBoxSubaccountAccessSettingsUpdateOpts struct {
	HomeDirectory       *string
	ReachableExternally *bool
	Readonly            *bool
	SambaEnabled        *bool
	SSHEnabled          *bool
	WebDAVEnabled       *bool
}

// UpdateSubaccountAccessSettings updates the access settings of a Storage Box subaccount.
func (c *StorageBoxClient) UpdateSubaccountAccessSettings(
	ctx context.Context,
	subaccount *StorageBoxSubaccount,
	opts StorageBoxSubaccountAccessSettingsUpdateOpts,
) (*Action, *Response, error) {
	const opPath = "/storage_boxes/%d/subaccounts/%d/actions/update_access_settings"
	ctx = ctxutil.SetOpPath(ctx, opPath)

	reqPath := fmt.Sprintf(opPath, subaccount.StorageBox.ID, subaccount.ID)
	reqBody := SchemaFromStorageBoxSubaccountUpdateAccessSettingsOpts(opts)

	respBody, resp, err := postRequest[schema.ActionGetResponse](ctx, c.client, reqPath, reqBody)
	if err != nil {
		return nil, resp, err
	}

	return ActionFromSchema(respBody.Action), resp, err
}
