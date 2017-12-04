package hcloud

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/hetznercloud/hcloud-go/hcloud/schema"
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

	var s schema.Action
	if err := json.Unmarshal(data, &s); err != nil {
		t.Fatal(err)
	}
	action := ActionFromSchema(s)

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

func TestActionClientGet(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/actions/1", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(schema.ActionGetResponse{
			Action: schema.Action{
				ID:       1,
				Status:   "running",
				Command:  "create_server",
				Progress: 50,
				Started:  time.Date(2017, 12, 4, 14, 31, 1, 0, time.UTC),
			},
		})
	})

	ctx := context.Background()
	action, _, err := env.Client.Action.Get(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
	if action == nil {
		t.Fatal("no action")
	}
	if action.ID != 1 {
		t.Errorf("unexpected action ID: %v", action.ID)
	}
}
