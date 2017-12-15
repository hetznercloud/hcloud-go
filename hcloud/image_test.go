package hcloud

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/hetznercloud/hcloud-go/hcloud/schema"
)

func TestImageClient(t *testing.T) {
	t.Run("Get", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/images/1", func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(schema.ImageGetResponse{
				Image: schema.Image{
					ID: 1,
				},
			})
		})

		ctx := context.Background()
		image, _, err := env.Client.Image.Get(ctx, 1)
		if err != nil {
			t.Fatal(err)
		}
		if image == nil {
			t.Fatal("no image")
		}
		if image.ID != 1 {
			t.Errorf("unexpected image ID: %v", image.ID)
		}
	})

	t.Run("List", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/images", func(w http.ResponseWriter, r *http.Request) {
			if page := r.URL.Query().Get("page"); page != "2" {
				t.Errorf("expected page 2; got %q", page)
			}
			if perPage := r.URL.Query().Get("per_page"); perPage != "50" {
				t.Errorf("expected per_page 50; got %q", perPage)
			}
			json.NewEncoder(w).Encode(schema.ImageListResponse{
				Images: []schema.Image{
					{ID: 1},
					{ID: 2},
				},
			})
		})

		opts := ImageListOpts{}
		opts.Page = 2
		opts.PerPage = 50

		ctx := context.Background()
		images, _, err := env.Client.Image.List(ctx, opts)
		if err != nil {
			t.Fatal(err)
		}
		if len(images) != 2 {
			t.Fatal("expected 2 images")
		}
	})

	t.Run("All", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		firstRequest := true
		env.Mux.HandleFunc("/images", func(w http.ResponseWriter, r *http.Request) {
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
					"images": [
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
					"images": [
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
					"images": [
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
		images, err := env.Client.Image.All(ctx)
		if err != nil {
			t.Fatalf("Image.List failed: %s", err)
		}
		if len(images) != 3 {
			t.Fatalf("expected 3 images; got %d", len(images))
		}
		if images[0].ID != 1 {
			t.Errorf("Expected first image to have an id of 1; got %d", images[0].ID)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/images/1", func(w http.ResponseWriter, r *http.Request) {
			return
		})

		ctx := context.Background()
		_, err := env.Client.Image.Delete(ctx, 1)
		if err != nil {
			t.Fatalf("Image.Delete failed: %s", err)
		}
	})
}
