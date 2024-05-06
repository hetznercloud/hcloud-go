package hcloud

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/hcloud-go/hcloud/schema"
)

func TestActionClientWatchOverallProgress(t *testing.T) {
	t.Parallel()
	env := newTestEnv()
	defer env.Teardown()

	callCount := 0

	env.Mux.HandleFunc("/actions", func(w http.ResponseWriter, r *http.Request) {
		callCount++
		var actions []schema.Action

		switch callCount {
		case 1:
			actions = []schema.Action{
				{
					ID:       1,
					Status:   "running",
					Progress: 50,
				},
				{
					ID:       2,
					Status:   "running",
					Progress: 50,
				},
			}
		case 2:
			actions = []schema.Action{
				{
					ID:       1,
					Status:   "running",
					Progress: 75,
				},
				{
					ID:       2,
					Status:   "error",
					Progress: 100,
					Error: &schema.ActionError{
						Code:    "action_failed",
						Message: "action failed",
					},
				},
			}
		case 3:
			actions = []schema.Action{
				{
					ID:       1,
					Status:   "success",
					Progress: 100,
				},
			}
		default:
			t.Errorf("unexpected number of calls to the test server: %v", callCount)
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(struct {
			Actions []schema.Action `json:"actions"`
			Meta    schema.Meta     `json:"meta"`
		}{
			Actions: actions,
			Meta: schema.Meta{
				Pagination: &schema.MetaPagination{
					Page:         1,
					LastPage:     1,
					PerPage:      len(actions),
					TotalEntries: len(actions),
				},
			},
		})
	})

	actions := []*Action{
		{
			ID:     1,
			Status: ActionStatusRunning,
		},
		{
			ID:     2,
			Status: ActionStatusRunning,
		},
	}

	ctx := context.Background()
	progressCh, errCh := env.Client.Action.WatchOverallProgress(ctx, actions)
	progressUpdates := []int{}
	errs := []error{}

	moreProgress, moreErrors := true, true

	for moreProgress || moreErrors {
		var progress int
		var err error

		select {
		case progress, moreProgress = <-progressCh:
			if moreProgress {
				progressUpdates = append(progressUpdates, progress)
			}
		case err, moreErrors = <-errCh:
			if moreErrors {
				errs = append(errs, err)
			}
		}
	}

	if len(errs) != 1 {
		t.Fatalf("expected to receive one error: %v", errs)
	}

	err := errs[0]

	if e, ok := errors.Unwrap(err).(ActionError); !ok || e.Code != "action_failed" {
		t.Fatalf("expected hcloud.Error, but got: %#v", err)
	}

	expectedProgressUpdates := []int{25, 62, 100}
	if !reflect.DeepEqual(progressUpdates, expectedProgressUpdates) {
		t.Fatalf("expected progresses %v but received %v", expectedProgressUpdates, progressUpdates)
	}
}

func TestActionClientWatchOverallProgressInvalidID(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	callCount := 0

	env.Mux.HandleFunc("/actions", func(w http.ResponseWriter, r *http.Request) {
		callCount++
		var actions []schema.Action

		switch callCount {
		case 1:
		default:
			t.Errorf("unexpected number of calls to the test server: %v", callCount)
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(struct {
			Actions []schema.Action `json:"actions"`
			Meta    schema.Meta     `json:"meta"`
		}{
			Actions: actions,
			Meta: schema.Meta{
				Pagination: &schema.MetaPagination{
					Page:         1,
					LastPage:     1,
					PerPage:      len(actions),
					TotalEntries: len(actions),
				},
			},
		})
	})

	actions := []*Action{
		{
			ID:     1,
			Status: ActionStatusRunning,
		},
	}

	ctx := context.Background()
	progressCh, errCh := env.Client.Action.WatchOverallProgress(ctx, actions)
	progressUpdates := []int{}
	errs := []error{}

	moreProgress, moreErrors := true, true

	for moreProgress || moreErrors {
		var progress int
		var err error

		select {
		case progress, moreProgress = <-progressCh:
			if moreProgress {
				progressUpdates = append(progressUpdates, progress)
			}
		case err, moreErrors = <-errCh:
			if moreErrors {
				errs = append(errs, err)
			}
		}
	}

	if len(errs) != 1 {
		t.Fatalf("expected to receive one error: %v", errs)
	}

	err := errs[0]

	assert.Equal(t, "actions not found: [1]", err.Error())

	expectedProgressUpdates := []int{}
	if !reflect.DeepEqual(progressUpdates, expectedProgressUpdates) {
		t.Fatalf("expected progresses %v but received %v", expectedProgressUpdates, progressUpdates)
	}
}

func TestActionClientWatchProgress(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	callCount := 0

	env.Mux.HandleFunc("/actions", func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")
		switch callCount {
		case 1:
			_, _ = w.Write([]byte(`{
				"actions": [
					{ "id": 1, "status": "running", "progress": 50 }
				],
				"meta": { "pagination": { "page": 1 }}
			}`))
		case 2:
			w.WriteHeader(http.StatusConflict)
			_, _ = w.Write([]byte(`{
				"error": { 
					"code": "conflict",
					"message": "conflict"
				}
			}`))
			return
		case 3:
			_, _ = w.Write([]byte(`{
				"actions": [
					{ "id": 1, "status": "error", "progress": 100, "error": {
						"code": "action_failed",
						"message": "action failed"
					} }
				],
				"meta": { "pagination": { "page": 1 }}
			}`))
		default:
			t.Errorf("unexpected number of calls to the test server: %v", callCount)
		}
	})
	action := &Action{
		ID:       1,
		Status:   ActionStatusRunning,
		Progress: 0,
	}

	ctx := context.Background()
	progressCh, errCh := env.Client.Action.WatchProgress(ctx, action)
	var (
		progressUpdates []int
		err             error
	)

loop:
	for {
		select {
		case progress := <-progressCh:
			progressUpdates = append(progressUpdates, progress)
		case err = <-errCh:
			break loop
		}
	}

	if err == nil {
		t.Fatal("expected an error")
	}
	if e, ok := err.(ActionError); !ok || e.Code != "action_failed" {
		t.Fatalf("expected hcloud.Error, but got: %#v", err)
	}
	if len(progressUpdates) != 1 || progressUpdates[0] != 50 {
		t.Fatalf("unexpected progress updates: %v", progressUpdates)
	}
}

func TestActionClientWatchProgressError(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/actions", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		_ = json.NewEncoder(w).Encode(schema.ErrorResponse{
			Error: schema.Error{
				Code:    string(ErrorCodeServiceError),
				Message: "service error",
			},
		})
	})

	action := &Action{ID: 1, Status: ActionStatusRunning}
	ctx := context.Background()
	_, errCh := env.Client.Action.WatchProgress(ctx, action)
	if err := <-errCh; err == nil {
		t.Fatal("expected an error")
	}
}

func TestActionClientWatchProgressInvalidID(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	callCount := 0

	env.Mux.HandleFunc("/actions", func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")
		switch callCount {
		case 1:
			_, _ = w.Write([]byte(`{
				"actions": [],
				"meta": { "pagination": { "page": 1 }}
			}`))
		default:
			t.Errorf("unexpected number of calls to the test server: %v", callCount)
		}
	})
	action := &Action{ID: 1, Status: ActionStatusRunning}

	ctx := context.Background()
	progressCh, errCh := env.Client.Action.WatchProgress(ctx, action)
	var (
		progressUpdates []int
		err             error
	)

loop:
	for {
		select {
		case progress := <-progressCh:
			progressUpdates = append(progressUpdates, progress)
		case err = <-errCh:
			break loop
		}
	}

	assert.Equal(t, "actions not found: [1]", err.Error())

	if len(progressUpdates) != 0 {
		t.Fatalf("unexpected progress updates: %v", progressUpdates)
	}
}
