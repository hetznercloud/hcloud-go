package hcloud

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/hcloud-go/v2/hcloud/schema"
)

func TestActionError(t *testing.T) {
	assert.Equal(t,
		"action failed (failed)",
		ActionError{Code: "failed", Message: "action failed"}.Error(),
	)
	assert.Equal(t,
		"action failed (failed, 12345)",
		ActionError{Code: "failed", Message: "action failed", action: &Action{ID: 12345}}.Error(),
	)
}

func TestActionClientGetByID(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/actions/1", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(schema.ActionGetResponse{
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
	action, _, err := env.Client.Action.GetByID(ctx, 1)
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

func TestActionClientGetByIDNotFound(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/actions/1", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(schema.ErrorResponse{
			Error: schema.Error{
				Code: string(ErrorCodeNotFound),
			},
		})
	})

	ctx := context.Background()
	action, _, err := env.Client.Action.GetByID(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
	if action != nil {
		t.Fatal("expected no action")
	}
}

func TestActionClientList(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/actions", func(w http.ResponseWriter, r *http.Request) {
		if page := r.URL.Query().Get("page"); page != "2" {
			t.Errorf("expected page 2; got %q", page)
		}
		if perPage := r.URL.Query().Get("per_page"); perPage != "50" {
			t.Errorf("expected per_page 50; got %q", perPage)
		}

		status := r.URL.Query()["status"]
		if len(status) != 2 {
			t.Errorf("expected status to contain 2 elements; got %q", status)
		} else {
			if status[0] != "running" {
				t.Errorf("expected status[0] to be running; got %q", status[0])
			}
			if status[1] != "error" {
				t.Errorf("expected status[1] to be error; got %q", status[1])
			}
		}

		sort := r.URL.Query()["sort"]
		if len(sort) != 3 {
			t.Errorf("expected sort to contain 3 elements; got %q", sort)
		} else {
			if sort[0] != "status" {
				t.Errorf("expected sort[0] to be status; got %q", sort[0])
			}
			if sort[1] != "progress:desc" {
				t.Errorf("expected sort[1] to be progress:desc; got %q", sort[1])
			}
			if sort[2] != "command:asc" {
				t.Errorf("expected sort[2] to be command:asc; got %q", sort[2])
			}
		}
		_ = json.NewEncoder(w).Encode(schema.ActionListResponse{
			Actions: []schema.Action{
				{ID: 1},
				{ID: 2},
			},
		})
	})

	opts := ActionListOpts{}
	opts.Page = 2
	opts.PerPage = 50
	opts.Status = []ActionStatus{ActionStatusRunning, ActionStatusError}
	opts.Sort = []string{"status", "progress:desc", "command:asc"}

	ctx := context.Background()
	actions, _, err := env.Client.Action.List(ctx, opts)
	if err != nil {
		t.Fatal(err)
	}
	if len(actions) != 2 {
		t.Fatal("expected 2 actions")
	}
}

func TestActionClientAll(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/actions", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(struct {
			Actions []schema.Action `json:"actions"`
			Meta    schema.Meta     `json:"meta"`
		}{
			Actions: []schema.Action{
				{ID: 1},
				{ID: 2},
				{ID: 3},
			},
			Meta: schema.Meta{
				Pagination: &schema.MetaPagination{
					Page:         1,
					LastPage:     1,
					PerPage:      3,
					TotalEntries: 3,
				},
			},
		})
	})

	ctx := context.Background()
	actions, err := env.Client.Action.All(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(actions) != 3 {
		t.Fatalf("expected 3 actions; got %d", len(actions))
	}
	if actions[0].ID != 1 || actions[1].ID != 2 || actions[2].ID != 3 {
		t.Errorf("unexpected actions")
	}
}

func TestResourceActionClientGetByID(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/primary_ips/actions/1", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(schema.ActionGetResponse{
			Action: schema.Action{
				ID:       1,
				Status:   "running",
				Command:  "create_primary_ip",
				Progress: 50,
				Started:  time.Date(2017, 12, 4, 14, 31, 1, 0, time.UTC),
			},
		})
	})

	ctx := context.Background()
	action, _, err := env.Client.PrimaryIP.Action.GetByID(ctx, 1)
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

func TestResourceActionClientGetByIDNotFound(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/primary_ips/actions/1", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(schema.ErrorResponse{
			Error: schema.Error{
				Code: string(ErrorCodeNotFound),
			},
		})
	})

	ctx := context.Background()
	action, _, err := env.Client.PrimaryIP.Action.GetByID(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
	if action != nil {
		t.Fatal("expected no action")
	}
}

func TestResourceActionClientList(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/primary_ips/actions", func(w http.ResponseWriter, r *http.Request) {
		if page := r.URL.Query().Get("page"); page != "2" {
			t.Errorf("expected page 2; got %q", page)
		}
		if perPage := r.URL.Query().Get("per_page"); perPage != "50" {
			t.Errorf("expected per_page 50; got %q", perPage)
		}

		status := r.URL.Query()["status"]
		if len(status) != 2 {
			t.Errorf("expected status to contain 2 elements; got %q", status)
		} else {
			if status[0] != "running" {
				t.Errorf("expected status[0] to be running; got %q", status[0])
			}
			if status[1] != "error" {
				t.Errorf("expected status[1] to be error; got %q", status[1])
			}
		}

		sort := r.URL.Query()["sort"]
		if len(sort) != 3 {
			t.Errorf("expected sort to contain 3 elements; got %q", sort)
		} else {
			if sort[0] != "status" {
				t.Errorf("expected sort[0] to be status; got %q", sort[0])
			}
			if sort[1] != "progress:desc" {
				t.Errorf("expected sort[1] to be progress:desc; got %q", sort[1])
			}
			if sort[2] != "command:asc" {
				t.Errorf("expected sort[2] to be command:asc; got %q", sort[2])
			}
		}
		_ = json.NewEncoder(w).Encode(schema.ActionListResponse{
			Actions: []schema.Action{
				{ID: 1},
				{ID: 2},
			},
		})
	})

	opts := ActionListOpts{}
	opts.Page = 2
	opts.PerPage = 50
	opts.Status = []ActionStatus{ActionStatusRunning, ActionStatusError}
	opts.Sort = []string{"status", "progress:desc", "command:asc"}

	ctx := context.Background()
	actions, _, err := env.Client.PrimaryIP.Action.List(ctx, opts)
	if err != nil {
		t.Fatal(err)
	}
	if len(actions) != 2 {
		t.Fatal("expected 2 actions")
	}
}

func TestResourceActionClientAll(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/primary_ips/actions", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(struct {
			Actions []schema.Action `json:"actions"`
			Meta    schema.Meta     `json:"meta"`
		}{
			Actions: []schema.Action{
				{ID: 1},
				{ID: 2},
				{ID: 3},
			},
			Meta: schema.Meta{
				Pagination: &schema.MetaPagination{
					Page:         1,
					LastPage:     1,
					PerPage:      3,
					TotalEntries: 3,
				},
			},
		})
	})

	ctx := context.Background()
	actions, err := env.Client.PrimaryIP.Action.All(ctx, ActionListOpts{})
	if err != nil {
		t.Fatal(err)
	}
	if len(actions) != 3 {
		t.Fatalf("expected 3 actions; got %d", len(actions))
	}
	if actions[0].ID != 1 || actions[1].ID != 2 || actions[2].ID != 3 {
		t.Errorf("unexpected actions")
	}
}
