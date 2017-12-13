package hcloud

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/hetznercloud/hcloud-go/hcloud/schema"
)

func TestServerTypeClient(t *testing.T) {
	t.Run("Get", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/server_types/1", func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(schema.ServerTypeGetResponse{
				ServerType: schema.ServerType{
					ID: 1,
				},
			})
		})

		ctx := context.Background()
		serverType, _, err := env.Client.ServerType.Get(ctx, 1)
		if err != nil {
			t.Fatal(err)
		}
		if serverType == nil {
			t.Fatal("no server type")
		}
		if serverType.ID != 1 {
			t.Errorf("unexpected server type ID: %v", serverType.ID)
		}
	})

	t.Run("List", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/server_types", func(w http.ResponseWriter, r *http.Request) {
			if page := r.URL.Query().Get("page"); page != "2" {
				t.Errorf("expected page 2; got %q", page)
			}
			if perPage := r.URL.Query().Get("per_page"); perPage != "50" {
				t.Errorf("expected per_page 50; got %q", perPage)
			}
			json.NewEncoder(w).Encode(schema.ServerTypeListResponse{
				ServerTypes: []schema.ServerType{
					{ID: 1},
					{ID: 2},
				},
			})
		})

		opts := ServerTypeListOpts{}
		opts.Page = 2
		opts.PerPage = 50

		ctx := context.Background()
		serverTypes, _, err := env.Client.ServerType.List(ctx, opts)
		if err != nil {
			t.Fatal(err)
		}
		if len(serverTypes) != 2 {
			t.Fatal("expected 2 server types")
		}
	})

	t.Run("All", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		firstRequest := true
		env.Mux.HandleFunc("/server_types", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			if firstRequest {
				firstRequest = false
				w.WriteHeader(http.StatusTooManyRequests)
				json.NewEncoder(w).Encode(schema.ErrorResponse{
					Error: schema.Error{
						Code:    "limit_reached",
						Message: "ratelimited",
					},
				})
				return
			}

			switch page := r.URL.Query().Get("page"); page {
			case "", "1":
				fmt.Fprint(w, `{
					"server_types": [
						{
							"id": 1
						}
					],
					"meta": {
						"pagination": {
							"page": 1,
							"per_page": 1,
							"previous_page": null,
							"next_page": 2,
							"last_page": 3,
							"total_entries": 3
						}
					}
				}`)
			case "2":
				fmt.Fprint(w, `{
					"server_types": [
						{
							"id": 2
						}
					],
					"meta": {
						"pagination": {
							"page": 2,
							"per_page": 1,
							"previous_page": 1,
							"next_page": 3,
							"last_page": 3,
							"total_entries": 3
						}
					}
				}`)
			case "3":
				fmt.Fprint(w, `{
					"server_types": [
						{
							"id": 3
						}
					],
					"meta": {
						"pagination": {
							"page": 3,
							"per_page": 1,
							"previous_page": 2,
							"next_page": null,
							"last_page": 3,
							"total_entries": 3
						}
					}
				}`)
			default:
				panic("bad page")
			}
		})

		ctx := context.Background()
		serverTypes, err := env.Client.ServerType.All(ctx)
		if err != nil {
			t.Fatalf("ServerTypes.List failed: %s", err)
		}
		if len(serverTypes) != 3 {
			t.Fatalf("expected 3 server types; got %d", len(serverTypes))
		}
		if serverTypes[0].ID != 1 {
			t.Errorf("")
		}
	})
}
