package hcloud

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hetznercloud/hcloud-go/v2/hcloud/exp/mockutil"
)

func TestGetSnapshot(t *testing.T) {
	ctx, server, client := makeTestUtils(t)

	t.Run("GetSnapshot (ByID)", func(t *testing.T) {
		server.Expect([]mockutil.Request{
			{
				Method: "GET", Path: "/storage_boxes/42/snapshots/13",
				Status: 200,
				JSONRaw: `{
				"snapshot": {
					"id": 42,
					"name": "my-resource",
					"description": "This describes my resource",
					"stats": {
				  		"size": 0,
				  		"size_filesystem": 0
						},
					"is_automatic": false,
					"labels": {
						"environment": "prod",
						"example.com/my": "label",
						"just-a-key": ""
					},
					"created": "2016-01-30T23:55:00+00:00",
					"storage_box": 42
				}
			}`,
			},
		})

		storageBox := &StorageBox{ID: 42}

		storageBoxSnapshot, _, err := client.StorageBox.GetSnapshotByID(ctx, storageBox, 13)
		require.NoError(t, err)
		require.NotNil(t, storageBox)

		assert.Equal(t, int64(42), storageBoxSnapshot.ID)
		assert.Equal(t, "my-resource", storageBoxSnapshot.Name)
		assert.Equal(t, "This describes my resource", storageBoxSnapshot.Description)
		assert.Equal(t, int64(0), storageBoxSnapshot.Stats.Size)

	})

	t.Run("GetSnapshot (ByName)", func(t *testing.T) {
		ctx, server, client := makeTestUtils(t)

		server.Expect([]mockutil.Request{
			{
				Method: "GET", Path: "/storage_boxes/42/snapshots?name=my-resource",
				Status: 200,
				JSONRaw: `{
				"snapshots": [{
					"id": 42,
					"name": "my-resource",
					"stats": {
				  		"size": 0,
				  		"size_filesystem": 0
						},
					"is_automatic": false,
					"labels": {
						"environment": "prod",
						"example.com/my": "label",
						"just-a-key": ""
					},
					"created": "2016-01-30T23:55:00+00:00",
					"storage_box": 42
				}]
			}`,
			},
		})

		storageBox := &StorageBox{ID: 42}

		storageBoxSnapshot, _, err := client.StorageBox.GetSnapshotByName(ctx, storageBox, "my-resource")
		require.NoError(t, err)
		require.NotNil(t, storageBox)

		assert.Equal(t, int64(42), storageBoxSnapshot.ID)
		assert.Equal(t, "my-resource", storageBoxSnapshot.Name)
	})

	t.Run("GetSnapshot (ByIDOrName)", func(t *testing.T) {
		ctx, server, client := makeTestUtils(t)

		server.Expect([]mockutil.Request{
			{
				Method: "GET", Path: "/storage_boxes/42/snapshots?name=my-resource",
				Status: 200,
				JSONRaw: `{
				"snapshots": [{
					"id": 42,
					"name": "my-resource",
					"stats": {
				  		"size": 0,
				  		"size_filesystem": 0
						},
					"is_automatic": false,
					"labels": {
						"environment": "prod",
						"example.com/my": "label",
						"just-a-key": ""
					},
					"created": "2016-01-30T23:55:00+00:00",
					"storage_box": 42
				}]
			}`,
			},
		})

		storageBox := &StorageBox{ID: 42}

		storageBoxSnapshot, _, err := client.StorageBox.GetSnapshot(ctx, storageBox, "my-resource")
		require.NoError(t, err)
		require.NotNil(t, storageBox)

		assert.Equal(t, int64(42), storageBoxSnapshot.ID)
		assert.Equal(t, "my-resource", storageBoxSnapshot.Name)
	})
}
