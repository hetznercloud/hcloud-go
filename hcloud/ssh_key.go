package hcloud

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

// SSHKey represents a SSH key in the Hetzner Cloud.
type SSHKey struct {
	ID          int
	Name        string
	Fingerprint string
	PublicKey   string
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *SSHKey) UnmarshalJSON(data []byte) error {
	var v struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		Fingerprint string `json:"fingerprint"`
		PublicKey   string `json:"public_key"`
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	s.ID = v.ID
	s.Name = v.Name
	s.Fingerprint = v.Fingerprint
	s.PublicKey = v.PublicKey

	return nil
}

// SSHKeyClient is a client for the SSH keys API.
type SSHKeyClient struct {
	client *Client
}

// Get retrieves a SSH key.
func (c *SSHKeyClient) Get(ctx context.Context, id int) (*SSHKey, *Response, error) {
	req, err := c.client.NewRequest(ctx, "GET", fmt.Sprintf("/ssh_keys/%d", id), nil)
	if err != nil {
		return nil, nil, err
	}

	var body struct {
		SSHKey *SSHKey `json:"ssh_key"`
	}
	resp, err := c.client.Do(req, &body)
	if err != nil {
		return nil, nil, err
	}
	return body.SSHKey, resp, nil
}

// SSHKeyListOpts specifies options for listing SSH keys.
type SSHKeyListOpts struct {
	ListOpts
}

// List returns a list of SSH keys for a specific page.
func (c *SSHKeyClient) List(ctx context.Context, opts SSHKeyListOpts) ([]*SSHKey, *Response, error) {
	path := "/ssh_keys"
	vals := url.Values{}
	if opts.Page > 0 {
		vals.Add("page", strconv.Itoa(opts.Page))
	}
	if opts.PerPage > 0 {
		vals.Add("per_page", strconv.Itoa(opts.PerPage))
	}
	path += "?" + vals.Encode()

	req, err := c.client.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	var body struct {
		SSHKeys []*SSHKey `json:"ssh_keys"`
	}
	resp, err := c.client.Do(req, &body)
	if err != nil {
		return nil, nil, err
	}
	return body.SSHKeys, resp, nil
}

// ListAll returns all SSH keys by going through all pages.
func (c *SSHKeyClient) ListAll(ctx context.Context) ([]*SSHKey, error) {
	allSSHKeys := []*SSHKey{}

	opts := SSHKeyListOpts{}
	opts.PerPage = 50

	_, err := c.client.all(func(page int) (*Response, error) {
		opts.Page = page
		sshKeys, resp, err := c.List(ctx, opts)
		if err != nil {
			return resp, err
		}
		allSSHKeys = append(allSSHKeys, sshKeys...)
		return resp, nil
	})
	if err != nil {
		return nil, err
	}

	return allSSHKeys, nil
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

	reqBody := struct {
		Name      string `json:"name"`
		PublicKey string `json:"public_key"`
	}{
		Name:      opts.Name,
		PublicKey: opts.PublicKey,
	}
	reqBodyData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, nil, err
	}

	req, err := c.client.NewRequest(ctx, "POST", "/ssh_keys", bytes.NewReader(reqBodyData))
	if err != nil {
		return nil, nil, err
	}

	var respBody struct {
		SSHKey *SSHKey `json:"ssh_key"`
	}
	resp, err := c.client.Do(req, &respBody)
	if err != nil {
		return nil, resp, err
	}
	return respBody.SSHKey, resp, nil
}

// Delete deletes a SSH key.
func (c *SSHKeyClient) Delete(ctx context.Context, id int) (*Response, error) {
	req, err := c.client.NewRequest(ctx, "DELETE", fmt.Sprintf("/ssh_keys/%d", id), nil)
	if err != nil {
		return nil, err
	}
	return c.client.Do(req, nil)
}
