package hcloud

import (
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hetznercloud/hcloud-go/v2/hcloud/exp/mockutil"
)

func TestStorageBoxClientGetSnapshot(t *testing.T) {
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
					"created": "2025-08-21T00:00:00Z",
					"storage_box": 42
				}
			}`,
			},
		})

		storageBox := &StorageBox{ID: 42}

		storageBoxSnapshot, _, err := client.StorageBox.GetSnapshotByID(ctx, storageBox, 13)
		require.NoError(t, err)
		require.NotNil(t, storageBoxSnapshot)
		require.NotNil(t, storageBoxSnapshot.Description)

		assert.Equal(t, int64(42), storageBoxSnapshot.ID)
		assert.Equal(t, "my-resource", storageBoxSnapshot.Name)
		assert.Equal(t, "This describes my resource", *storageBoxSnapshot.Description)
		assert.Equal(t, uint64(0), storageBoxSnapshot.Stats.Size)
		assert.Equal(t, time.Date(2025, 8, 21, 00, 00, 0, 0, time.UTC), storageBoxSnapshot.Created)
		assert.Equal(t, "prod", storageBoxSnapshot.Labels["environment"])
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

	t.Run("GetSnapshot (NotFound)", func(t *testing.T) {
		ctx, server, client := makeTestUtils(t)

		server.Expect([]mockutil.Request{
			{
				Method: "GET", Path: "/storage_boxes/42/snapshots/13",
				Status: 404,
				JSONRaw: `{
					"error": {
						"code": "not_found",
						"message": "The resource you requested could not be found."
					}
				}`,
			},
		})

		storageBox := &StorageBox{ID: 42}

		storageBoxSnapshot, resp, err := client.StorageBox.GetSnapshotByID(ctx, storageBox, 13)
		require.NotNil(t, resp)
		require.NoError(t, err)
		require.Nil(t, storageBoxSnapshot)

		assert.Equal(t, 404, resp.StatusCode)
	})
}

func TestStorageBoxClientListSnapshot(t *testing.T) {
	ctx, server, client := makeTestUtils(t)

	server.Expect([]mockutil.Request{
		{
			Method: "GET", Path: "/storage_boxes/42/snapshots?label_selector=environment%3Dprod",
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

	opts := StorageBoxSnapshotListOpts{
		LabelSelector: "environment=prod",
	}
	snapshots, err := client.StorageBox.AllSnapshotsWithOpts(ctx, storageBox, opts)
	require.NoError(t, err)
	require.Len(t, snapshots, 1)

	assert.Equal(t, int64(42), snapshots[0].ID)
	assert.Equal(t, "my-resource", snapshots[0].Name)
}

func TestStorageBoxClientCreateSnapshot(t *testing.T) {
	ctx, server, client := makeTestUtils(t)

	t.Run("CreateSnapshot (With Description)", func(t *testing.T) {
		server.Expect([]mockutil.Request{
			{
				Method: "POST", Path: "/storage_boxes/42/snapshots",
				Status: 201,
				Want: func(t *testing.T, r *http.Request) {
					body, err := io.ReadAll(r.Body)
					require.NoError(t, err)

					assert.JSONEq(t, `{ "description": "Test Snapshot" }`, string(body))
				},
				JSONRaw: `{
				"snapshot": { "id": 14, "storage_box": 42 },
				"action": { "id": 13 }
			}`,
			},
		})

		storageBox := &StorageBox{ID: 42}

		opts := StorageBoxSnapshotCreateOpts{
			Ptr("Test Snapshot"),
		}
		result, _, err := client.StorageBox.CreateSnapshot(ctx, storageBox, opts)
		require.NoError(t, err)
		require.NotNil(t, result.Action)
		require.NotNil(t, result.Snapshot)
	})

	t.Run("CreateSnapshot (Without Description)", func(t *testing.T) {
		server.Expect([]mockutil.Request{
			{
				Method: "POST", Path: "/storage_boxes/42/snapshots",
				Status: 201,
				Want: func(t *testing.T, r *http.Request) {
					body, err := io.ReadAll(r.Body)
					require.NoError(t, err)

					assert.JSONEq(t, `{}`, string(body))
				},
				JSONRaw: `{
				"snapshot": { "id": 14, "storage_box": 42 },
				"action": { "id": 13 }
			}`,
			},
		})

		storageBox := &StorageBox{ID: 42}

		opts := StorageBoxSnapshotCreateOpts{}
		result, _, err := client.StorageBox.CreateSnapshot(ctx, storageBox, opts)
		require.NoError(t, err)
		require.NotNil(t, result.Action)
		require.NotNil(t, result.Snapshot)
	})
}

func TestStorageBoxClientUpdateSnapshot(t *testing.T) {
	ctx, server, client := makeTestUtils(t)

	server.Expect([]mockutil.Request{
		{
			Method: "PUT", Path: "/storage_boxes/42/snapshots/13",
			Status: 200,
			Want: func(t *testing.T, r *http.Request) {
				body, err := io.ReadAll(r.Body)
				require.NoError(t, err)

				assert.JSONEq(t, `{ "labels": { "environment": "prod" } }`, string(body))
			},
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
						"environment": "prod"
					},
					"created": "2016-01-30T23:55:00+00:00",
					"storage_box": 42
				}
			}`,
		},
	})

	storageBoxSnapshot := &StorageBoxSnapshot{
		ID: 13,
		StorageBox: &StorageBox{
			ID: 42,
		},
	}

	opts := StorageBoxSnapshotUpdateOpts{
		Labels: map[string]string{
			"environment": "prod",
		},
	}
	storageBoxSnapshot, _, err := client.StorageBox.UpdateSnapshot(ctx, storageBoxSnapshot, opts)
	require.NoError(t, err)
	require.NotNil(t, storageBoxSnapshot)

	assert.Equal(t, "prod", storageBoxSnapshot.Labels["environment"])
}

func TestStorageBoxClientDeleteSnapshot(t *testing.T) {
	ctx, server, client := makeTestUtils(t)

	server.Expect([]mockutil.Request{
		{
			Method: "DELETE", Path: "/storage_boxes/42/snapshots/13",
			Status:  200,
			JSONRaw: `{ "action": { "id": 5 } }`,
		},
	})

	storageBoxSnapshot := &StorageBoxSnapshot{
		ID: 13,
		StorageBox: &StorageBox{
			ID: 42,
		},
	}

	action, _, err := client.StorageBox.DeleteSnapshot(ctx, storageBoxSnapshot)
	require.NoError(t, err)
	require.NotNil(t, action)
}
