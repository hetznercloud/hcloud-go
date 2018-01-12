package hcloud

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/hetznercloud/hcloud-go/hcloud/schema"
)

func TestImageClient(t *testing.T) {
	t.Run("GetByID", func(t *testing.T) {
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
		image, _, err := env.Client.Image.GetByID(ctx, 1)
		if err != nil {
			t.Fatal(err)
		}
		if image == nil {
			t.Fatal("no image")
		}
		if image.ID != 1 {
			t.Errorf("unexpected image ID: %v", image.ID)
		}

		t.Run("via Get", func(t *testing.T) {
			image, _, err := env.Client.Image.Get(ctx, "1")
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
	})

	t.Run("GetByID (not found)", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/images/1", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(schema.ErrorResponse{
				Error: schema.Error{
					Code: ErrorCodeNotFound,
				},
			})
		})

		ctx := context.Background()
		image, _, err := env.Client.Image.GetByID(ctx, 1)
		if err != nil {
			t.Fatal(err)
		}
		if image != nil {
			t.Fatal("expected no image")
		}
	})

	t.Run("GetByName", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/images", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.RawQuery != "name=my+image" {
				t.Fatal("missing name query")
			}
			json.NewEncoder(w).Encode(schema.ImageListResponse{
				Images: []schema.Image{
					{
						ID: 1,
					},
				},
			})
		})

		ctx := context.Background()
		image, _, err := env.Client.Image.GetByName(ctx, "my image")
		if err != nil {
			t.Fatal(err)
		}
		if image == nil {
			t.Fatal("no image")
		}
		if image.ID != 1 {
			t.Errorf("unexpected image ID: %v", image.ID)
		}

		t.Run("via Get", func(t *testing.T) {
			image, _, err := env.Client.Image.Get(ctx, "my image")
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
	})

	t.Run("GetByName (not found)", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/images", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.RawQuery != "name=my+image" {
				t.Fatal("missing name query")
			}
			json.NewEncoder(w).Encode(schema.ImageListResponse{
				Images: []schema.Image{},
			})
		})

		ctx := context.Background()
		image, _, err := env.Client.Image.GetByName(ctx, "my image")
		if err != nil {
			t.Fatal(err)
		}
		if image != nil {
			t.Fatal("unexpected image")
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
			if types, ok := r.URL.Query()["type"]; !ok {
				t.Errorf("expected type to be set; got %q", types)
			} else if len(types) != 2 || types[0] != "snapshot" || types[1] != "system" {
				t.Errorf("expected type to be ['snapshot','system']; got %q", types)
			}
			if boundTo, ok := r.URL.Query()["bound_to"]; !ok {
				t.Errorf("expected bound_to to be set; got %q", boundTo)
			} else if len(boundTo) != 2 || boundTo[0] != "1" || boundTo[1] != "2" {
				t.Errorf("expected bound_to to be ['1','2']; got %q", boundTo)
			}
			if sort, ok := r.URL.Query()["sort"]; !ok {
				t.Errorf("expected sort to be set; got %q", sort)
			} else if len(sort) != 2 ||
				sort[0] != "id:asc" ||
				sort[1] != "name:desc" {
				t.Errorf("expected sort to be ['id:asc', 'name:desc']; got %q", sort)
			}
			json.NewEncoder(w).Encode(schema.ImageListResponse{
				Images: []schema.Image{
					{ID: 1},
					{ID: 2},
				},
			})
		})

		opts := ImageListOpts{
			Types: []ImageType{
				ImageTypeSnapshot,
				ImageTypeSystem,
			},
			BoundTo: []*Server{
				{ID: 1},
				{ID: 2},
			},
		}
		opts.Page = 2
		opts.PerPage = 50
		opts.
			Sort("id", Asc).
			Sort("name", Desc)

		ctx := context.Background()
		page := env.Client.Image.List(ctx, opts)
		if page.GoTo(2) || page.Err() != nil {
			t.Fatalf("unexpected error or resource not exhausted on page.GoTo(2): %v", page.Err())
		}
		if len(page.Content()) != 2 {
			t.Fatal("expected 2 images")
		}
	})

	t.Run("All", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/images", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(struct {
				Images []schema.Image `json:"images"`
				Meta   schema.Meta    `json:"meta"`
			}{
				Images: []schema.Image{
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
		images, err := env.Client.Image.All(ctx)
		if err != nil {
			t.Fatalf("Image.List failed: %s", err)
		}
		if len(images) != 3 {
			t.Fatalf("expected 3 images; got %d", len(images))
		}
		if images[0].ID != 1 || images[1].ID != 2 || images[2].ID != 3 {
			t.Errorf("unexpected images")
		}
	})

	t.Run("Delete", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/images/1", func(w http.ResponseWriter, r *http.Request) {
			return
		})

		var (
			ctx   = context.Background()
			image = &Image{ID: 1}
		)
		_, err := env.Client.Image.Delete(ctx, image)
		if err != nil {
			t.Fatalf("Image.Delete failed: %s", err)
		}
	})
}

func TestImageClientUpdate(t *testing.T) {
	var (
		ctx   = context.Background()
		image = &Image{ID: 1}
	)

	t.Run("description and type", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/images/1", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "PUT" {
				t.Error("expected PUT")
			}
			var reqBody schema.ImageUpdateRequest
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				t.Fatal(err)
			}
			if reqBody.Description == nil || *reqBody.Description != "test" {
				t.Errorf("unexpected description: %v", reqBody.Description)
			}
			if reqBody.Type == nil || *reqBody.Type != "snapshot" {
				t.Errorf("unexpected type: %v", reqBody.Type)
			}
			json.NewEncoder(w).Encode(schema.ImageUpdateResponse{
				Image: schema.Image{
					ID: 1,
				},
			})
		})

		opts := ImageUpdateOpts{
			Description: String("test"),
			Type:        ImageTypeSnapshot,
		}
		updatedImage, _, err := env.Client.Image.Update(ctx, image, opts)
		if err != nil {
			t.Fatal(err)
		}

		if updatedImage.ID != 1 {
			t.Errorf("unexpected image ID: %v", updatedImage.ID)
		}
	})

	t.Run("no updates", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/images/1", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "PUT" {
				t.Error("expected PUT")
			}
			var reqBody schema.ImageUpdateRequest
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				t.Fatal(err)
			}
			if reqBody.Description != nil {
				t.Errorf("unexpected no description, but got: %v", reqBody.Description)
			}
			if reqBody.Type != nil {
				t.Errorf("unexpected no type, but got: %v", reqBody.Type)
			}
			json.NewEncoder(w).Encode(schema.ImageUpdateResponse{
				Image: schema.Image{
					ID: 1,
				},
			})
		})

		opts := ImageUpdateOpts{}
		updatedImage, _, err := env.Client.Image.Update(ctx, image, opts)
		if err != nil {
			t.Fatal(err)
		}

		if updatedImage.ID != 1 {
			t.Errorf("unexpected image ID: %v", updatedImage.ID)
		}
	})
}
