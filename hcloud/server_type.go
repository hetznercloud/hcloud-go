package hcloud

import (
	"context"
	"fmt"

	"github.com/hetznercloud/hcloud-go/hcloud/schema"
)

// ServerType represents a server type in the Hetzner Cloud.
type ServerType struct {
	ID          int
	Name        string
	Description string
	Cores       int
	Memory      float32
	Disk        int
	StorageType StorageType
}

// StorageType specifies the type of storage.
type StorageType string

const (
	// StorageTypeLocal is the type for local storage.
	StorageTypeLocal StorageType = "local"

	// StorageTypeCeph is the type for remote storage.
	StorageTypeCeph = "ceph"
)

// ServerTypeClient is a client for the server types API.
type ServerTypeClient struct {
	client *Client
}

// Get retrieves a server type.
func (c *ServerTypeClient) Get(ctx context.Context, id int) (*ServerType, *Response, error) {
	req, err := c.client.NewRequest(ctx, "GET", fmt.Sprintf("/server_types/%d", id), nil)
	if err != nil {
		return nil, nil, err
	}

	var body schema.ServerTypeGetResponse
	resp, err := c.client.Do(req, &body)
	if err != nil {
		return nil, nil, err
	}
	return ServerTypeFromSchema(body.ServerType), resp, nil
}

// ServerTypeListOpts specifies options for listing server types.
type ServerTypeListOpts struct {
	ListOpts
}

// List returns a list of server types for a specific page.
func (c *ServerTypeClient) List(ctx context.Context, opts ServerTypeListOpts) ([]*ServerType, *Response, error) {
	path := "/server_types?" + valuesForListOpts(opts.ListOpts).Encode()
	req, err := c.client.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	var body schema.ServerTypeListResponse
	resp, err := c.client.Do(req, &body)
	if err != nil {
		return nil, nil, err
	}
	serverTypes := make([]*ServerType, 0, len(body.ServerTypes))
	for _, s := range body.ServerTypes {
		serverTypes = append(serverTypes, ServerTypeFromSchema(s))
	}
	return serverTypes, resp, nil
}

// All returns all servers.
func (c *ServerTypeClient) All(ctx context.Context) ([]*ServerType, error) {
	allServerTypes := []*ServerType{}

	opts := ServerTypeListOpts{}
	opts.PerPage = 50

	_, err := c.client.all(func(page int) (*Response, error) {
		opts.Page = page
		serverTypes, resp, err := c.List(ctx, opts)
		if err != nil {
			return resp, err
		}
		allServerTypes = append(allServerTypes, serverTypes...)
		return resp, nil
	})
	if err != nil {
		return nil, err
	}

	return allServerTypes, nil
}
