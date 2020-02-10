package hcloud

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/hetznercloud/hcloud-go/hcloud/schema"
)

func TestLoadBalancerTypeClient(t *testing.T) {
	t.Run("GetByID", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/load_balancer_types/1", func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(schema.LoadBalancerTypeGetResponse{
				LoadBalancerType: schema.LoadBalancerType{
					ID: 1,
				},
			})
		})

		ctx := context.Background()
		loadBalancerType, _, err := env.Client.LoadBalancerType.GetByID(ctx, 1)
		if err != nil {
			t.Fatal(err)
		}
		if loadBalancerType == nil {
			t.Fatal("no load balancer type")
		}
		if loadBalancerType.ID != 1 {
			t.Errorf("unexpected loadBalancerType ID: %v", loadBalancerType.ID)
		}

		t.Run("via Get", func(t *testing.T) {
			loadBalancerType, _, err := env.Client.LoadBalancerType.Get(ctx, "1")
			if err != nil {
				t.Fatal(err)
			}
			if loadBalancerType == nil {
				t.Fatal("no load balancer type")
			}
			if loadBalancerType.ID != 1 {
				t.Errorf("unexpected loadBalancerType ID: %v", loadBalancerType.ID)
			}
		})
	})

	t.Run("GetByID (not found)", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/load_balancer_types/1", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(schema.ErrorResponse{
				Error: schema.Error{
					Code: string(ErrorCodeNotFound),
				},
			})
		})

		ctx := context.Background()
		loadBalancerType, _, err := env.Client.LoadBalancerType.GetByID(ctx, 1)
		if err != nil {
			t.Fatal(err)
		}
		if loadBalancerType != nil {
			t.Fatal("expected no loadBalancerType")
		}
	})

	t.Run("GetByName", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/load_balancer_types", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.RawQuery != "name=lb1" {
				t.Fatal("missing name query")
			}
			json.NewEncoder(w).Encode(schema.LoadBalancerTypeListResponse{
				LoadBalancerTypes: []schema.LoadBalancerType{
					{
						ID: 1,
					},
				},
			})
		})

		ctx := context.Background()
		loadBalancerType, _, err := env.Client.LoadBalancerType.GetByName(ctx, "lb1")
		if err != nil {
			t.Fatal(err)
		}
		if loadBalancerType == nil {
			t.Fatal("no loadBalancerType")
		}
		if loadBalancerType.ID != 1 {
			t.Errorf("unexpected loadBalancerType ID: %v", loadBalancerType.ID)
		}

		t.Run("via Get", func(t *testing.T) {
			loadBalancerType, _, err := env.Client.LoadBalancerType.Get(ctx, "lb1")
			if err != nil {
				t.Fatal(err)
			}
			if loadBalancerType == nil {
				t.Fatal("no loadBalancerType")
			}
			if loadBalancerType.ID != 1 {
				t.Errorf("unexpected loadBalancerType ID: %v", loadBalancerType.ID)
			}
		})
	})

	t.Run("GetByName (not found)", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/load_balancer_types", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.RawQuery != "name=lb1" {
				t.Fatal("missing name query")
			}
			json.NewEncoder(w).Encode(schema.LoadBalancerTypeListResponse{
				LoadBalancerTypes: []schema.LoadBalancerType{},
			})
		})

		ctx := context.Background()
		loadBalancerType, _, err := env.Client.LoadBalancerType.GetByName(ctx, "lb1")
		if err != nil {
			t.Fatal(err)
		}
		if loadBalancerType != nil {
			t.Fatal("unexpected loadBalancerType")
		}
	})

	t.Run("List", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/load_balancer_types", func(w http.ResponseWriter, r *http.Request) {
			if page := r.URL.Query().Get("page"); page != "2" {
				t.Errorf("expected page 2; got %q", page)
			}
			if perPage := r.URL.Query().Get("per_page"); perPage != "50" {
				t.Errorf("expected per_page 50; got %q", perPage)
			}
			json.NewEncoder(w).Encode(schema.LoadBalancerTypeListResponse{
				LoadBalancerTypes: []schema.LoadBalancerType{
					{ID: 1},
					{ID: 2},
				},
			})
		})

		opts := LoadBalancerTypeListOpts{}
		opts.Page = 2
		opts.PerPage = 50

		ctx := context.Background()
		loadBalancerTypes, _, err := env.Client.LoadBalancerType.List(ctx, opts)
		if err != nil {
			t.Fatal(err)
		}
		if len(loadBalancerTypes) != 2 {
			t.Fatal("expected 2 loadBalancerTypes")
		}
	})

	t.Run("All", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/load_balancer_types", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(struct {
				LoadBalancerTypes []schema.LoadBalancerType `json:"load_balancer_types"`
				Meta              schema.Meta               `json:"meta"`
			}{
				LoadBalancerTypes: []schema.LoadBalancerType{
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
		loadBalancerTypes, err := env.Client.LoadBalancerType.All(ctx)
		if err != nil {
			t.Fatalf("LoadBalancerType.List failed: %s", err)
		}
		if len(loadBalancerTypes) != 3 {
			t.Fatalf("expected 3 load balancer types; got %d", len(loadBalancerTypes))
		}
		if loadBalancerTypes[0].ID != 1 || loadBalancerTypes[1].ID != 2 || loadBalancerTypes[2].ID != 3 {
			t.Errorf("unexpected load balancer types")
		}
	})
}
