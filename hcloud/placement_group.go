package hcloud

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/hetznercloud/hcloud-go/hcloud/schema"
)

type PlacementGroup struct {
	ID      int
	Name    string
	Labels  map[string]string
	Created time.Time
	Servers []int
	Type    PlacementGroupType
}

type PlacementGroupType string

const (
	PlacementGroupTypeSpread PlacementGroupType = "spread"
)

// FirewallClient is a client for the Placement Groups API.
type PlacementGroupClient struct {
	client *Client
}

func (c *PlacementGroupClient) GetByID(ctx context.Context, id int) (*PlacementGroup, *Response, error) {
	req, err := c.client.NewRequest(ctx, "GET", fmt.Sprintf("/placement_groups/%d", id), nil)
	if err != nil {
		return nil, nil, err
	}

	var body schema.PlacementGroupGetResponse
	resp, err := c.client.Do(req, &body)
	if err != nil {
		if IsError(err, ErrorCodeNotFound) {
			return nil, resp, nil
		}
		return nil, nil, err
	}
	return PlacementGroupFromSchema(body.PlacementGroup), resp, nil
}

func (c *PlacementGroupClient) GetByName(ctx context.Context, name string) (*PlacementGroup, *Response, error) {
	if name == "" {
		return nil, nil, nil
	}
	placementGroups, response, err := c.List(ctx, PlacementGroupListOpts{Name: name})
	if len(placementGroups) == 0 {
		return nil, response, err
	}
	return placementGroups[0], response, err
}

func (c *PlacementGroupClient) Get(ctx context.Context, idOrName string) (*PlacementGroup, *Response, error) {
	if id, err := strconv.Atoi(idOrName); err == nil {
		return c.GetByID(ctx, int(id))
	}
	return c.GetByName(ctx, idOrName)
}

type PlacementGroupListOpts struct {
	ListOpts
	Name string
	Type PlacementGroupType
}

func (l PlacementGroupListOpts) values() url.Values {
	vals := l.ListOpts.values()
	if l.Name != "" {
		vals.Add("name", l.Name)
	}
	if l.Type != "" {
		vals.Add("type", string(l.Type))
	}
	return vals
}

func (c *PlacementGroupClient) List(ctx context.Context, opts PlacementGroupListOpts) ([]*PlacementGroup, *Response, error) {
	path := "/placement_groups?" + opts.values().Encode()
	req, err := c.client.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	var body schema.PlacementGroupListResponse
	resp, err := c.client.Do(req, &body)
	if err != nil {
		return nil, nil, err
	}
	placementGroups := make([]*PlacementGroup, 0, len(body.PlacementGroups))
	for _, g := range body.PlacementGroups {
		placementGroups = append(placementGroups, PlacementGroupFromSchema(g))
	}
	return placementGroups, resp, nil
}

func (c *PlacementGroupClient) All(ctx context.Context) ([]*PlacementGroup, error) {
	opts := PlacementGroupListOpts{}
	opts.PerPage = 50

	return c.AllWithOpts(ctx, opts)
}

func (c *PlacementGroupClient) AllWithOpts(ctx context.Context, opts PlacementGroupListOpts) ([]*PlacementGroup, error) {
	allPlacmentGroups := []*PlacementGroup{}

	err := c.client.all(func(page int) (*Response, error) {
		opts.Page = page
		placementGroups, resp, err := c.List(ctx, opts)
		if err != nil {
			return resp, err
		}
		allPlacmentGroups = append(allPlacmentGroups, placementGroups...)
		return resp, nil
	})
	if err != nil {
		return nil, err
	}

	return allPlacmentGroups, nil
}

type PlacementGroupUpdateOpts struct {
	Name   string
	Labels map[string]string
}

func (c *PlacementGroupClient) Update(ctx context.Context, placementGroup *PlacementGroup, opts PlacementGroupUpdateOpts) (*PlacementGroup, *Response, error) {
	reqBody := schema.PlacementGroupUpdateRequest{}
	if opts.Name != "" {
		reqBody.Name = &opts.Name
	}
	if opts.Labels != nil {
		reqBody.Labels = &opts.Labels
	}
	reqBodyData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, nil, err
	}

	path := fmt.Sprintf("/placement_groups/%d", placementGroup.ID)
	req, err := c.client.NewRequest(ctx, "PUT", path, bytes.NewReader(reqBodyData))
	if err != nil {
		return nil, nil, err
	}

	respBody := schema.PlacementGroupUpdateResponse{}
	resp, err := c.client.Do(req, &respBody)
	if err != nil {
		return nil, resp, err
	}

	return PlacementGroupFromSchema(respBody.PlacementGroup), resp, nil
}

func (c *PlacementGroupClient) Delete(ctx context.Context, placementGroup *PlacementGroup) (*Response, error) {
	req, err := c.client.NewRequest(ctx, "DELETE", fmt.Sprintf("/placement_groups/%d", placementGroup.ID), nil)
	if err != nil {
		return nil, err
	}
	return c.client.Do(req, nil)
}
