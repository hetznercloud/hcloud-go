package hcloud

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

// SSHKey is a SSH key in the Hetzner Cloud.
type SSHKey struct {
	ID          int
	Name        string
	Fingerprint string
	PublicKey   string
}

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

type SSHKeyClient struct {
	client *Client
}

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

func (c *SSHKeyClient) List(ctx context.Context) ([]*SSHKey, *Response, error) {
	req, err := c.client.NewRequest(ctx, "GET", "/ssh_keys", nil)
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

type SSHKeyCreateOpts struct {
	Name      string
	PublicKey string
}

func (o SSHKeyCreateOpts) Validate() error {
	if o.Name == "" {
		return errors.New("missing name")
	}
	if o.PublicKey == "" {
		return errors.New("missing public key")
	}
	return nil
}

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

func (c *SSHKeyClient) Delete(ctx context.Context, id int) (*Response, error) {
	req, err := c.client.NewRequest(ctx, "DELETE", fmt.Sprintf("/ssh_keys/%d", id), nil)
	if err != nil {
		return nil, err
	}
	return c.client.Do(req, nil)
}
