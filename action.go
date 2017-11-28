package hcloud // import "hetzner.cloud/hcloud"

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
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

// UnmarshalJSON implements json.Unmarshaler.
func (a *Action) UnmarshalJSON(data []byte) error {
	var v struct {
		ID       int       `json:"id"`
		Status   string    `json:"status"`
		Command  string    `json:"command"`
		Progress int       `json:"progress"`
		Started  time.Time `json:"started"`
		Finished time.Time `json:"finished"`
		Error    struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	a.ID = v.ID
	a.Status = v.Status
	a.Command = v.Command
	a.Progress = v.Progress
	a.Started = v.Started
	a.Finished = v.Finished
	a.ErrorCode = v.Error.Code
	a.ErrorMessage = v.Error.Message

	return nil
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
		Action *Action `json:"action"`
	}
	resp, err := c.client.Do(req, &body)
	if err != nil {
		return nil, nil, err
	}
	return body.Action, resp, nil
}
