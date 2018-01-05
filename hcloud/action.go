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
	Status       ActionStatus
	Command      string
	Progress     int
	Started      time.Time
	Finished     time.Time
	ErrorCode    string
	ErrorMessage string
	Resources    []*ResourceReference
}

// ActionStatus represents an action's status.
type ActionStatus string

// List of action statuses.
const (
	ActionStatusRunning ActionStatus = "running"
	ActionStatusSuccess              = "success"
	ActionStatusError                = "error"
)

// ResourceReference references other resources from an action.
type ResourceReference struct {
	ID   int
	Type ResourceReferenceType
}

// ResourceReferenceType represents an action's resource reference type.
type ResourceReferenceType string

// List of action resource reference types.
const (
	ResourceReferenceTypeServer     ResourceReferenceType = "server"
	ResourceReferenceTypeImage                            = "image"
	ResourceReferenceTypeSSHKey                           = "ssh_keys"
	ResourceReferenceTypeFloatingIP                       = "floating_ip"
)

func (a *Action) Error() error {
	if a.ErrorCode != "" && a.ErrorMessage != "" {
		return fmt.Errorf("%s (%s)", a.ErrorMessage, a.ErrorCode)
	}
	return nil
}

// ActionClient is a client for the actions API.
type ActionClient struct {
	client *Client
}

// GetByID retrieves an action by its ID.
func (c *ActionClient) GetByID(ctx context.Context, id int) (*Action, *Response, error) {
	req, err := c.client.NewRequest(ctx, "GET", fmt.Sprintf("/actions/%d", id), nil)
	if err != nil {
		return nil, nil, err
	}

	var body schema.ActionGetResponse
	resp, err := c.client.Do(req, &body)
	if err != nil {
		if IsError(err, ErrorCodeNotFound) {
			return nil, resp, nil
		}
		return nil, nil, err
	}
	return ActionFromSchema(body.Action), resp, nil
}
