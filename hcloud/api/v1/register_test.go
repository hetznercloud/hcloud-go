package v1

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/hetznercloud/hcloud-go/hcloud/api"
)

func TestActionFromSchema(t *testing.T) {
	data := []byte(`{
		"id": 1,
		"command": "create_server",
		"status": "success",
		"progress": 100,
		"started": "2016-01-30T23:55:00Z",
		"finished": "2016-01-30T23:56:13Z",
		"resources": [
			{
				"id": 42,
				"type": "server"
			}
		],
		"error": {
			"code": "action_failed",
			"message": "Action failed"
		}
	}`)

	var s Action
	if err := json.Unmarshal(data, &s); err != nil {
		t.Fatal(err)
	}

	action := &api.Action{}

	err := Scheme.Convert(&s, action)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if action.ID != 1 {
		t.Errorf("unexpected ID: %v", action.ID)
	}
	if action.Command != "create_server" {
		t.Errorf("unexpected command: %v", action.Command)
	}
	if action.Status != "success" {
		t.Errorf("unexpected status: %v", action.Status)
	}
	if action.Progress != 100 {
		t.Errorf("unexpected progress: %d", action.Progress)
	}
	if !action.Started.Equal(time.Date(2016, 1, 30, 23, 55, 0, 0, time.UTC)) {
		t.Errorf("unexpected started: %v", action.Started)
	}
	if !action.Finished.Equal(time.Date(2016, 1, 30, 23, 56, 13, 0, time.UTC)) {
		t.Errorf("unexpected finished: %v", action.Started)
	}
	if action.ErrorCode != "action_failed" {
		t.Errorf("unexpected error code: %v", action.ErrorCode)
	}
	if action.ErrorMessage != "Action failed" {
		t.Errorf("unexpected error message: %v", action.ErrorMessage)
	}
}
