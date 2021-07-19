package hcloud

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

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

}

func TestPlacementGroupClientGebByName(t *testing.T) {

}

func TestPlacementGroupClientGebByNameNotFound(t *testing.T) {

}

func TestPlacementGroupClientGebByNameEmpty(t *testing.T) {

}

func TestPlacementGroupDelete(t *testing.T) {

}

func TestPlacementGroupUpdate(t *testing.T) {

}
