package hcloud

import (
	"encoding/json"
	"testing"
	"time"
)

func TestActionUnmarshalJSON(t *testing.T) {
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

	var v Action
	if err := json.Unmarshal(data, &v); err != nil {
		t.Fatal(err)
	}

	if v.ID != 1 {
		t.Errorf("unexpected ID: %v", v.ID)
	}
	if v.Command != "create_server" {
		t.Errorf("unexpected command: %v", v.Command)
	}
	if v.Status != "success" {
		t.Errorf("unexpected status: %v", v.Status)
	}
	if v.Progress != 100 {
		t.Errorf("unexpected progress: %d", v.Progress)
	}
	if !v.Started.Equal(time.Date(2016, 1, 30, 23, 55, 0, 0, time.UTC)) {
		t.Errorf("unexpected started: %v", v.Started)
	}
	if !v.Finished.Equal(time.Date(2016, 1, 30, 23, 56, 13, 0, time.UTC)) {
		t.Errorf("unexpected finished: %v", v.Started)
	}
	if v.ErrorCode != "action_failed" {
		t.Errorf("unexpected error code: %v", v.ErrorCode)
	}
	if v.ErrorMessage != "Action failed" {
		t.Errorf("unexpected error message: %v", v.ErrorMessage)
	}
}
