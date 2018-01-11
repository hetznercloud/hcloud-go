package hcloud

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/hetznercloud/hcloud-go/hcloud/schema"
)

// Datacenter represents a datacenter in the Hetzner Cloud.
type Datacenter struct {
	ID          int
	Name        string
	Description string
	Location    *Location
	ServerTypes DatacenterServerTypes
}

// DatacenterServerTypes represents the server types available and supported in a datacenter.
type DatacenterServerTypes struct {
	Supported []*ServerType
	Available []*ServerType
}

// DatacenterClient is a client for the datacenter API.
type DatacenterClient struct {
	client *Client
}

// GetByID retrieves a datacenter by its ID.
func (c *DatacenterClient) GetByID(ctx context.Context, id int) (*Datacenter, *Response, error) {
	req, err := c.client.NewRequest(ctx, "GET", fmt.Sprintf("/datacenters/%d", id), nil)
	if err != nil {
		return nil, nil, err
	}

	var body schema.DatacenterGetResponse
	resp, err := c.client.Do(req, &body)
	if err != nil {
		if IsError(err, ErrorCodeNotFound) {
			return nil, resp, nil
		}
		return nil, resp, err
	}
	return DatacenterFromSchema(body.Datacenter), resp, nil
}

// GetByName retrieves an datacenter by its name.
func (c *DatacenterClient) GetByName(ctx context.Context, name string) (*Datacenter, *Response, error) {
	path := "/datacenters?name=" + url.QueryEscape(name)
	req, err := c.client.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	var body schema.DatacenterListResponse
	resp, err := c.client.Do(req, &body)
	if err != nil {
		return nil, nil, err
	}

	if len(body.Datacenters) == 0 {
		return nil, resp, nil
	}
	return DatacenterFromSchema(body.Datacenters[0]), resp, nil
}

// Get retrieves a datacenter by its ID if the input can be parsed as an integer, otherwise it retrieves a datacenter by its name.
func (c *DatacenterClient) Get(ctx context.Context, idOrName string) (*Datacenter, *Response, error) {
	if id, err := strconv.Atoi(idOrName); err == nil {
		return c.GetByID(ctx, int(id))
	}
	return c.GetByName(ctx, idOrName)
}

// DatacenterPage serves as accessor of the datacenters API pagination.
type DatacenterPage struct {
	Page
	content []*Datacenter
}

// Content contains the content of the current page.
func (p *DatacenterPage) Content() []*Datacenter {
	return p.content
}

// DatacenterListOpts specifies options for listing datacenters.
type DatacenterListOpts struct {
	ListOpts
}

// List returns an accessor to control the datacenters API pagination.
func (c *DatacenterClient) List(ctx context.Context, opts DatacenterListOpts) *DatacenterPage {
	page := &DatacenterPage{}
	page.pageGetter = pageGetter(func(start, end int) (resp *Response, exhausted bool, err error) {
		allDatacenters := []*Datacenter{}
		if opts.PerPage == 0 {
			opts.PerPage = 50
		}

		resp, exhausted, err = c.client.all(func(page int) (*Response, error) {
			opts.Page = page
			datacenters, resp, err := c.list(ctx, opts)
			if err != nil {
				return resp, err
			}
			allDatacenters = append(allDatacenters, datacenters...)
			return resp, nil
		}, start, end)
		page.content = allDatacenters
		return
	})
	return page
}

// list returns a list of datacenters for a specific page.
func (c *DatacenterClient) list(ctx context.Context, opts DatacenterListOpts) ([]*Datacenter, *Response, error) {
	path := "/datacenters?" + opts.URLValues().Encode()
	req, err := c.client.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	var body schema.DatacenterListResponse
	resp, err := c.client.Do(req, &body)
	if err != nil {
		return nil, nil, err
	}
	datacenters := make([]*Datacenter, 0, len(body.Datacenters))
	for _, i := range body.Datacenters {
		datacenters = append(datacenters, DatacenterFromSchema(i))
	}
	return datacenters, resp, nil
}

// All returns all datacenters.
func (c *DatacenterClient) All(ctx context.Context) ([]*Datacenter, error) {
	opts := DatacenterListOpts{}
	opts.PerPage = 50
	page := c.List(ctx, opts)
	if page.All(); page.Err() != nil {
		return nil, page.Err()
	}
	return page.Content(), nil
}
