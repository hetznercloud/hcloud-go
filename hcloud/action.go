package hcloud

import (
	"context"
	"fmt"
	"time"

	"github.com/hetznercloud/hcloud-go/hcloud/schema"
)

// Action represents an action in the Hetzner Cloud.
type Action struct {
	ID           int
	Status       string
	Command      string
	Progress     int
	Started      time.Time
	Finished     time.Time
	ErrorCode    string
	ErrorMessage string
}

// ActionFromSchema converts a schema.Action to an Action.
func ActionFromSchema(s schema.Action) Action {
	action := Action{
		ID:       s.ID,
		Status:   s.Status,
		Command:  s.Command,
		Progress: s.Progress,
		Started:  s.Started,
		Finished: s.Finished,
	}
	if s.Error != nil {
		action.ErrorCode = s.Error.Code
		action.ErrorMessage = s.Error.Message
	}
	return action
}

// ActionClient is a client for the actions API.
type ActionClient struct {
	client *Client
}

// Get retrieves an action.
func (c *ActionClient) Get(ctx context.Context, id int) (*Action, *Response, error) {
	req, err := c.client.NewRequest(ctx, "GET", fmt.Sprintf("/actions/%d", id), nil)
	if err != nil {
		return nil, nil, err
	}

	var body struct {
		Action schema.Action `json:"action"`
	}
	resp, err := c.client.Do(req, &body)
	if err != nil {
		return nil, nil, err
	}
	action := ActionFromSchema(body.Action)
	return &action, resp, nil
}
