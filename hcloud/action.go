package hcloud

import (
	"context"
	"fmt"
	"net/url"
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
	Resources    []*ActionResource
}

// ActionStatus represents an action's status.
type ActionStatus string

// List of action statuses.
const (
	ActionStatusRunning ActionStatus = "running"
	ActionStatusSuccess              = "success"
	ActionStatusError                = "error"
)

// ActionResource references other resources from an action.
type ActionResource struct {
	ID   int
	Type ActionResourceType
}

// ActionResourceType represents an action's resource reference type.
type ActionResourceType string

// List of action resource reference types.
const (
	ActionResourceTypeServer     ActionResourceType = "server"
	ActionResourceTypeImage                         = "image"
	ActionResourceTypeISO                           = "iso"
	ActionResourceTypeFloatingIP                    = "floating_ip"
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

// ActionPage serves as accessor of the actions API pagination.
type ActionPage struct {
	Page
	content []*Action
}

// Content contains the content of the current page.
func (p *ActionPage) Content() []*Action {
	return p.content
}

// ActionListOpts specifies options for listing actions.
type ActionListOpts struct {
	ListOpts
	Status     []ActionStatus
	Server     *Server
	FloatingIP *FloatingIP
}

// URLValues returns the list opts as url.Values.
func (o ActionListOpts) URLValues() url.Values {
	vals := o.ListOpts.URLValues()
	for _, s := range o.Status {
		vals.Add("status", string(s))
	}
	return vals
}

// List returns an accessor to control the actions API pagination.
func (c *ActionClient) List(ctx context.Context, opts ActionListOpts) *ActionPage {
	page := &ActionPage{}
	page.pageGetter = pageGetter(func(start, end int) (resp *Response, exhausted bool, err error) {
		allActions := []*Action{}
		if opts.PerPage == 0 {
			opts.PerPage = 50
		}

		resp, exhausted, err = c.client.all(func(page int) (*Response, error) {
			opts.Page = page
			actions, resp, err := c.list(ctx, opts)
			if err != nil {
				return resp, err
			}
			allActions = append(allActions, actions...)
			return resp, nil
		}, start, end)
		page.content = allActions
		return
	})
	return page
}

// list returns a list of actions for a specific page.
func (c *ActionClient) list(ctx context.Context, opts ActionListOpts) ([]*Action, *Response, error) {
	path := "/actions?"
	if opts.Server != nil {
		path = fmt.Sprintf("/servers/%d/actions?", opts.Server.ID)
	}
	if opts.FloatingIP != nil {
		path = fmt.Sprintf("/floating_ips/%d/actions?", opts.FloatingIP.ID)
	}
	path = path + opts.URLValues().Encode()
	req, err := c.client.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	var body schema.ActionListResponse
	resp, err := c.client.Do(req, &body)
	if err != nil {
		return nil, nil, err
	}
	actions := make([]*Action, 0, len(body.Actions))
	for _, i := range body.Actions {
		actions = append(actions, ActionFromSchema(i))
	}
	return actions, resp, nil
}

// All returns all actions.
func (c *ActionClient) All(ctx context.Context) ([]*Action, error) {
	opts := ActionListOpts{}
	opts.PerPage = 50
	page := c.List(ctx, opts)
	if page.All(); page.Err() != nil {
		return nil, page.Err()
	}
	return page.Content(), nil
}
