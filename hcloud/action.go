package hcloud

import (
	"context"
	"fmt"

	"github.com/hetznercloud/hcloud-go/hcloud/api"
	"github.com/hetznercloud/hcloud-go/hcloud/api/v1"
)

// ActionClient is a client for the actions API.
type ActionClient struct {
	client *Client
}

// Get retrieves an action.
func (c *ActionClient) Get(ctx context.Context, id int) (*api.Action, *Response, error) {
	req, err := c.client.NewRequest(ctx, "GET", fmt.Sprintf("/actions/%d", id), nil)
	if err != nil {
		return nil, nil, err
	}

	var body v1.ActionGetResponse
	resp, err := c.client.Do(req, &body)
	if err != nil {
		return nil, nil, err
	}
	action := &api.Action{}
	if err := v1.Scheme.Convert(&body.Action, action); err != nil {
		return nil, nil, err
	}
	return action, resp, nil
}
