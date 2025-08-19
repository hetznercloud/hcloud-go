package hcloud

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/hetznercloud/hcloud-go/v2/hcloud/exp/ctxutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/schema"
)

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

type StorageBoxSubaccountAccessSettings struct {
	ReachableExternally bool
	Readonly            bool
	SambaEnabled        bool
	SSHEnabled          bool
	WebDAVEnabled       bool
}

type StorageBoxSubaccountAccessSettingsOpts struct {
	ReachableExternally *bool
	Readonly            *bool
	SambaEnabled        *bool
	SSHEnabled          *bool
	WebDAVEnabled       *bool
}

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

type StorageBoxSubaccountListOpts struct {
	LabelSelector string
}

func (o StorageBoxSubaccountListOpts) values() url.Values {
	vals := url.Values{}
	if o.LabelSelector != "" {
		vals.Add("label_selector", o.LabelSelector)
	}
	return vals
}

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

func (c *StorageBoxClient) AllSubaccountsWithOpts(
	ctx context.Context,
	storageBox *StorageBox,
	opts StorageBoxSubaccountListOpts,
) ([]*StorageBoxSubaccount, error) {
	subaccounts, _, err := c.ListSubaccounts(ctx, storageBox, opts)
	return subaccounts, err
}

func (c *StorageBoxClient) AllSubaccounts(
	ctx context.Context,
	storageBox *StorageBox,
) ([]*StorageBoxSubaccount, error) {
	opts := StorageBoxSubaccountListOpts{}
	subaccounts, _, err := c.ListSubaccounts(ctx, storageBox, opts)
	return subaccounts, err
}

type StorageBoxSubaccountCreateOpts struct {
	Password       string
	HomeDirectory  *string
	AccessSettings *StorageBoxSubaccountAccessSettingsOpts
	Description    *string
	Labels         map[string]string
}

type StorageBoxSubaccountCreateResult struct {
	Subaccount *StorageBoxSubaccount
	Action     *Action
}

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

	result.Subaccount = &StorageBoxSubaccount{
		ID: respBody.Subaccount.ID,
		StorageBox: &StorageBox{
			ID: respBody.Subaccount.StorageBox,
		},
	}

	return result, resp, nil
}

type StorageBoxSubaccountUpdateOpts struct {
	Labels      map[string]string
	Description string
}

func (c *StorageBoxClient) UpdateSubaccount(
	ctx context.Context,
	storageBox *StorageBox,
	subaccount *StorageBoxSubaccount,
	opts StorageBoxSubaccountUpdateOpts,
) (*StorageBoxSubaccount, *Response, error) {
	const opPath = "/storage_boxes/%d/subaccounts/%d"
	ctx = ctxutil.SetOpPath(ctx, opPath)

	reqPath := fmt.Sprintf(opPath, storageBox.ID, subaccount.ID)
	reqBody := SchemaFromStorageBoxSubaccountUpdateOpts(opts)

	respBody, resp, err := putRequest[schema.StorageBoxSubaccountUpdateResponse](ctx, c.client, reqPath, reqBody)
	if err != nil {
		return nil, resp, err
	}

	updatedSubaccount := StorageBoxSubaccountFromSchema(respBody.Subaccount)

	return updatedSubaccount, resp, nil
}

func (c *StorageBoxClient) DeleteSubaccount(
	ctx context.Context,
	storageBox *StorageBox,
	subaccount *StorageBoxSubaccount,
) (*Action, *Response, error) {
	const opPath = "/storage_boxes/%d/subaccounts/%d"
	ctx = ctxutil.SetOpPath(ctx, opPath)

	reqPath := fmt.Sprintf(opPath, storageBox.ID, subaccount.ID)

	respBody, resp, err := deleteRequest[schema.ActionGetResponse](ctx, c.client, reqPath)
	if err != nil {
		return nil, resp, err
	}

	action := ActionFromSchema(respBody.Action)

	return action, resp, nil
}
