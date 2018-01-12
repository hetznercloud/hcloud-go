package hcloud

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/hetznercloud/hcloud-go/hcloud/schema"
)

func TestActionClientGetByID(t *testing.T) {
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
		json.NewEncoder(w).Encode(schema.ErrorResponse{
			Error: schema.Error{
				Code: ErrorCodeNotFound,
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
		if status, ok := r.URL.Query()["status"]; !ok {
			t.Errorf("expected status to be set; got %q", status)
		} else if len(status) != 2 || status[0] != "running" || status[1] != "success" {
			t.Errorf("expected status ['running', 'success']; got %q", status)
		}
		if sort, ok := r.URL.Query()["sort"]; !ok {
			t.Errorf("expected sort to be set; got %q", sort)
		} else if len(sort) != 4 ||
			sort[0] != "status:asc" ||
			sort[1] != "id:desc" ||
			sort[2] != "started:asc" ||
			sort[3] != "progress:desc" {
			t.Errorf("expected sort to be ['status:asc', 'id:desc', 'started:asc', 'progress:desc']; got %q", sort)
		}
		json.NewEncoder(w).Encode(schema.ActionListResponse{
			Actions: []schema.Action{
				{ID: 1},
				{ID: 2},
			},
		})
	})

	opts := ActionListOpts{
		Status: []ActionStatus{
			ActionStatusRunning,
			ActionStatusSuccess,
		},
	}
	opts.Page = 2
	opts.PerPage = 50
	opts.
		Sort("status", Asc).
		Sort("id", Desc).
		Sort("started", Asc).
		Sort("progress", Desc)

	ctx := context.Background()
	page := env.Client.Action.List(ctx, opts)
	if page.GoTo(2) || page.Err() != nil {
		t.Fatalf("unexpected error or resource not exhausted on page.GoTo(2): %v", page.Err())
	}
	if len(page.Content()) != 2 {
		t.Fatal("expected 2 actions")
	}
}

func TestActionClientAll(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/actions", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(struct {
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
