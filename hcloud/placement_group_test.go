package hcloud

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hetznercloud/hcloud-go/hcloud/schema"
)

func TestPlacementGroupClientGebByID(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	const id = 1

	env.Mux.HandleFunc(fmt.Sprintf("/placement_groups/%v", id), func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(schema.PlacementGroupGetResponse{
			PlacementGroup: schema.PlacementGroup{
				ID: id,
			},
		})
	})

	checkError := func(t *testing.T, placementGroup *PlacementGroup, err error) {
		if err != nil {
			t.Fatal(err)
		}
		if placementGroup == nil {
			t.Fatal("no placment group")
		}
		if placementGroup.ID != id {
			t.Errorf("unexpected placment group ID: %v", placementGroup.ID)
		}
	}

	ctx := context.Background()

	t.Run("called via GetByID", func(t *testing.T) {
		placementGroup, _, err := env.Client.PlacementGroup.GetByID(ctx, 1)
		checkError(t, placementGroup, err)
	})

	t.Run("called via Get", func(t *testing.T) {
		placementGroup, _, err := env.Client.PlacementGroup.Get(ctx, "1")
		checkError(t, placementGroup, err)
	})
}

func TestPlacementGroupClientGebByIDNotFound(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/placement_groups/1", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(schema.ErrorResponse{
			Error: schema.Error{
				Code: string(ErrorCodeNotFound),
			},
		})
	})

	ctx := context.Background()

	placementGroup, _, err := env.Client.PlacementGroup.GetByID(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
	if placementGroup != nil {
		t.Fatal("expected no placment_group")
	}
}

func TestPlacementGroupClientGetByName(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	const (
		id   = 1
		name = "my_placement_group"
	)

	env.Mux.HandleFunc("/placement_groups", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RawQuery != fmt.Sprintf("name=%v", name) {
			t.Fatal("missing name query")
		}
		json.NewEncoder(w).Encode(schema.PlacementGroupListResponse{
			PlacementGroups: []schema.PlacementGroup{
				{
					ID:   id,
					Name: name,
				},
			},
		})
	})

	checkError := func(t *testing.T, placementGroup *PlacementGroup, err error) {
		if err != nil {
			t.Fatal(err)
		}
		if placementGroup == nil {
			t.Fatal("no placment group")
		}
		if placementGroup.ID != id {
			t.Errorf("unexpected placment group ID: %v", placementGroup.ID)
		}
		if placementGroup.Name != name {
			t.Errorf("unexpected placement group Name: %v", placementGroup.Name)
		}
	}

	ctx := context.Background()

	t.Run("called via GetByID", func(t *testing.T) {
		placementGroup, _, err := env.Client.PlacementGroup.GetByName(ctx, name)
		checkError(t, placementGroup, err)
	})

	t.Run("called via Get", func(t *testing.T) {
		placementGroup, _, err := env.Client.PlacementGroup.Get(ctx, name)
		checkError(t, placementGroup, err)
	})
}

func TestPlacementGroupClientGebByNameNotFound(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	const name = "my_placement_group"

	env.Mux.HandleFunc("/placement_groups", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RawQuery != fmt.Sprintf("name=%v", name) {
			t.Fatal("missing name query")
		}
		json.NewEncoder(w).Encode(schema.PlacementGroupListResponse{
			PlacementGroups: []schema.PlacementGroup{},
		})
	})

	ctx := context.Background()

	placementGroup, _, err := env.Client.PlacementGroup.GetByName(ctx, name)
	if err != nil {
		t.Fatal(err)
	}
	if placementGroup != nil {
		t.Fatal("expected no placment_group")
	}
}

func TestPlacementGroupClientGebByNameEmpty(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	ctx := context.Background()

	placementGroup, _, err := env.Client.PlacementGroup.GetByName(ctx, "")
	if err != nil {
		t.Fatal(err)
	}
	if placementGroup != nil {
		t.Fatal("expected no placment_group")
	}
}

func TestPlacementGroupDelete(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	const id = 1

	env.Mux.HandleFunc(fmt.Sprintf("/placement_groups/%v", id), func(w http.ResponseWriter, r *http.Request) {})

	var (
		ctx           = context.Background()
		placmentGroup = &PlacementGroup{ID: id}
	)

	_, err := env.Client.PlacementGroup.Delete(ctx, placmentGroup)
	if err != nil {
		t.Fatal(err)
	}
}

func TestPlacementGroupUpdate(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	const id = 1

	var (
		ctx           = context.Background()
		placmentGroup = &PlacementGroup{ID: id}
		opts          = PlacementGroupUpdateOpts{
			Name:   "test",
			Labels: map[string]string{"key": "value"},
		}
	)

	env.Mux.HandleFunc("/placement_groups/1", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Error("expected PUT")
		}
		var reqBody schema.PlacementGroupUpdateRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatal(err)
		}
		expectedReqBody := schema.PlacementGroupUpdateRequest{
			Name:   &opts.Name,
			Labels: &opts.Labels,
		}
		if !cmp.Equal(expectedReqBody, reqBody) {
			t.Log(cmp.Diff(expectedReqBody, reqBody))
			t.Error("unexpected request body")
		}
		json.NewEncoder(w).Encode(schema.PlacementGroupUpdateResponse{
			PlacementGroup: schema.PlacementGroup{
				ID: id,
			},
		})
	})

	updatedPlacementGroup, _, err := env.Client.PlacementGroup.Update(ctx, placmentGroup, opts)
	if err != nil {
		t.Fatal(err)
	}
	if updatedPlacementGroup == nil {
		t.Fatal("no placment group")
	}
	if updatedPlacementGroup.ID != id {
		t.Errorf("unexpected placment group ID: %v", updatedPlacementGroup.ID)
	}
}
