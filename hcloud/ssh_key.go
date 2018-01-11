package hcloud

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"

	"github.com/hetznercloud/hcloud-go/hcloud/schema"
)

// SSHKey represents a SSH key in the Hetzner Cloud.
type SSHKey struct {
	ID          int
	Name        string
	Fingerprint string
	PublicKey   string
}

// SSHKeyClient is a client for the SSH keys API.
type SSHKeyClient struct {
	client *Client
}

// GetByID retrieves a SSH key by its ID.
func (c *SSHKeyClient) GetByID(ctx context.Context, id int) (*SSHKey, *Response, error) {
	req, err := c.client.NewRequest(ctx, "GET", fmt.Sprintf("/ssh_keys/%d", id), nil)
	if err != nil {
		return nil, nil, err
	}

	var body schema.SSHKeyGetResponse
	resp, err := c.client.Do(req, &body)
	if err != nil {
		if IsError(err, ErrorCodeNotFound) {
			return nil, resp, nil
		}
		return nil, nil, err
	}
	return SSHKeyFromSchema(body.SSHKey), resp, nil
}

// GetByName retrieves a SSH key by its name.
func (c *SSHKeyClient) GetByName(ctx context.Context, name string) (*SSHKey, *Response, error) {
	path := "/ssh_keys?name=" + url.QueryEscape(name)
	req, err := c.client.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	var body schema.SSHKeyListResponse
	resp, err := c.client.Do(req, &body)
	if err != nil {
		return nil, nil, err
	}

	if len(body.SSHKeys) == 0 {
		return nil, resp, nil
	}
	return SSHKeyFromSchema(body.SSHKeys[0]), resp, nil
}

// Get retrieves a SSH key by its ID if the input can be parsed as an integer, otherwise it retrieves a SSH key by its name.
func (c *SSHKeyClient) Get(ctx context.Context, idOrName string) (*SSHKey, *Response, error) {
	if id, err := strconv.Atoi(idOrName); err == nil {
		return c.GetByID(ctx, int(id))
	}
	return c.GetByName(ctx, idOrName)
}

// SSHKeyPage serves as accessor of the SSH keys API pagination.
type SSHKeyPage struct {
	Page
	content []*SSHKey
}

// Content contains the content of the current page.
func (p *SSHKeyPage) Content() []*SSHKey {
	return p.content
}

// SSHKeyListOpts specifies options for listing SSH keys.
type SSHKeyListOpts struct {
	ListOpts
}

// List returns an accessor to control the SSH keys API pagination.
func (c *SSHKeyClient) List(ctx context.Context, opts SSHKeyListOpts) *SSHKeyPage {
	page := &SSHKeyPage{}
	page.pageGetter = pageGetter(func(start, end int) (resp *Response, exhausted bool, err error) {
		allSSHKeys := []*SSHKey{}
		if opts.PerPage == 0 {
			opts.PerPage = 50
		}

		resp, exhausted, err = c.client.all(func(page int) (*Response, error) {
			opts.Page = page
			sshKeys, resp, err := c.list(ctx, opts)
			if err != nil {
				return resp, err
			}
			allSSHKeys = append(allSSHKeys, sshKeys...)
			return resp, nil
		}, start, end)
		page.content = allSSHKeys
		return
	})
	return page
}

// list returns a list of SSH keys for a specific page.
func (c *SSHKeyClient) list(ctx context.Context, opts SSHKeyListOpts) ([]*SSHKey, *Response, error) {
	path := "/ssh_keys?" + opts.URLValues().Encode()
	req, err := c.client.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	var body schema.SSHKeyListResponse
	resp, err := c.client.Do(req, &body)
	if err != nil {
		return nil, nil, err
	}
	sshKeys := make([]*SSHKey, 0, len(body.SSHKeys))
	for _, s := range body.SSHKeys {
		sshKeys = append(sshKeys, SSHKeyFromSchema(s))
	}
	return sshKeys, resp, nil
}

// All returns all SSH keys.
func (c *SSHKeyClient) All(ctx context.Context) ([]*SSHKey, error) {
	opts := SSHKeyListOpts{}
	opts.PerPage = 50
	page := c.List(ctx, opts)
	if page.All(); page.Err() != nil {
		return nil, page.Err()
	}
	return page.Content(), nil
}

// SSHKeyCreateOpts specifies parameters for creating a SSH key.
type SSHKeyCreateOpts struct {
	Name      string
	PublicKey string
}

// Validate checks if options are valid.
func (o SSHKeyCreateOpts) Validate() error {
	if o.Name == "" {
		return errors.New("missing name")
	}
	if o.PublicKey == "" {
		return errors.New("missing public key")
	}
	return nil
}

// Create creates a new SSH key with the given options.
func (c *SSHKeyClient) Create(ctx context.Context, opts SSHKeyCreateOpts) (*SSHKey, *Response, error) {
	if err := opts.Validate(); err != nil {
		return nil, nil, err
	}

	reqBodyData, err := json.Marshal(schema.SSHKeyCreateRequest{
		Name:      opts.Name,
		PublicKey: opts.PublicKey,
	})
	if err != nil {
		return nil, nil, err
	}

	req, err := c.client.NewRequest(ctx, "POST", "/ssh_keys", bytes.NewReader(reqBodyData))
	if err != nil {
		return nil, nil, err
	}

	var respBody schema.SSHKeyCreateResponse
	resp, err := c.client.Do(req, &respBody)
	if err != nil {
		return nil, resp, err
	}
	return SSHKeyFromSchema(respBody.SSHKey), resp, nil
}

// Delete deletes a SSH key.
func (c *SSHKeyClient) Delete(ctx context.Context, sshKey *SSHKey) (*Response, error) {
	req, err := c.client.NewRequest(ctx, "DELETE", fmt.Sprintf("/ssh_keys/%d", sshKey.ID), nil)
	if err != nil {
		return nil, err
	}
	return c.client.Do(req, nil)
}

// SSHKeyUpdateOpts specifies options for updating a SSH key.
type SSHKeyUpdateOpts struct {
	Name string
}

// Update updates a SSH key.
func (c *SSHKeyClient) Update(ctx context.Context, sshKey *SSHKey, opts SSHKeyUpdateOpts) (*SSHKey, *Response, error) {
	reqBody := schema.SSHKeyUpdateRequest{
		Name: opts.Name,
	}
	reqBodyData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, nil, err
	}

	path := fmt.Sprintf("/ssh_keys/%d", sshKey.ID)
	req, err := c.client.NewRequest(ctx, "PUT", path, bytes.NewReader(reqBodyData))
	if err != nil {
		return nil, nil, err
	}

	respBody := schema.SSHKeyUpdateResponse{}
	resp, err := c.client.Do(req, &respBody)
	if err != nil {
		return nil, resp, err
	}
	return SSHKeyFromSchema(respBody.SSHKey), resp, nil
}
