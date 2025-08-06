package hcloud

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hetznercloud/hcloud-go/v2/hcloud/exp/mockutil"
)

func TestStorageBoxTypeClientGetByID(t *testing.T) {
	t.Run("GetByID", func(t *testing.T) {
		ctx, server, client := makeTestUtils(t)

		server.Expect([]mockutil.Request{
			{
				Method: "GET", Path: "/storage_box_types/42",
				Status: 200,
				JSONRaw: `
				{
					"storage_box_type": {
						"id": 42,
						"name": "bx11",
						"description": "BX11",
						"snapshot_limit": 10,
						"automatic_snapshot_limit": 10,
						"subaccounts_limit": 200,
						"size": 1073741824,
						"prices": [
							{
								"location": "fsn1",
								"price_hourly": {"net": "1.0000", "gross": "1.1900"},
								"price_monthly": {"net": "1.0000", "gross": "1.1900"},
								"setup_fee": {"net": "1.0000", "gross": "1.1900"}
							}
						],
						"deprecation": {
							"unavailable_after": "2023-09-01T00:00:00+00:00",
							"announced": "2023-06-01T00:00:00+00:00"
						}
					}
				}`,
			},
		})

		storageBoxType, _, err := client.StorageBoxType.GetByID(ctx, 42)
		require.NoError(t, err)
		require.NotNil(t, storageBoxType, "no storage box type")
		assert.Equal(t, int64(42), storageBoxType.ID, "unexpected storage box type ID")
		assert.Equal(t, "bx11", storageBoxType.Name, "unexpected storage box type name")
	})

	t.Run("GetByID (not found)", func(t *testing.T) {
		ctx, server, client := makeTestUtils(t)

		server.Expect([]mockutil.Request{
			{
				Method: "GET", Path: "/storage_box_types/1",
				Status:  404,
				JSONRaw: `{"error": {"code": "not_found"}}`,
			},
		})

		storageBoxType, _, err := client.StorageBoxType.GetByID(ctx, 1)
		require.NoError(t, err)
		assert.Nil(t, storageBoxType, "expected no storage box type")
	})
}

func TestStorageBoxTypeClientList(t *testing.T) {
	ctx, server, client := makeTestUtils(t)

	server.Expect([]mockutil.Request{
		{
			Method: "GET", Path: "/storage_box_types?page=2&per_page=50",
			Status: 200,
			JSONRaw: `{
				"storage_box_types": [
					{"id": 1, "name": "bx11"},
					{"id": 2, "name": "bx21"}
				]
			}`,
		},
	})

	opts := StorageBoxTypeListOpts{}
	opts.Page = 2
	opts.PerPage = 50
	storageBoxTypes, _, err := client.StorageBoxType.List(ctx, opts)
	require.NoError(t, err)
	assert.Len(t, storageBoxTypes, 2, "expected 2 storage box types")
	assert.Equal(t, "bx11", storageBoxTypes[0].Name)
	assert.Equal(t, "bx21", storageBoxTypes[1].Name)
}

func TestStorageBoxTypeClientAll(t *testing.T) {
	ctx, server, client := makeTestUtils(t)

	server.Expect([]mockutil.Request{
		{
			Method: "GET", Path: "/storage_box_types?page=1&per_page=50",
			Status: 200,
			JSONRaw: `{
				"storage_box_types": [
					{"id": 1, "name": "bx11"},
					{"id": 2, "name": "bx21"}
				],
				"meta": {"pagination": {"page": 1, "last_page": 2, "per_page": 2, "next_page": 2, "previous_page": null, "total_entries": 3}}
			}`,
		},
		{
			Method: "GET", Path: "/storage_box_types?page=2&per_page=50",
			Status: 200,
			JSONRaw: `{
				"storage_box_types": [
					{"id": 3, "name": "bx31"}
				],
				"meta": {"pagination": {"page": 2, "last_page": 2, "per_page": 2, "next_page": null, "previous_page": 1, "total_entries": 3}}
			}`,
		},
	})

	storageBoxTypes, err := client.StorageBoxType.All(ctx)
	require.NoError(t, err, "StorageBoxType.List failed")
	assert.Len(t, storageBoxTypes, 3, "expected 3 storage box types")
	assert.Equal(t, int64(1), storageBoxTypes[0].ID)
	assert.Equal(t, int64(2), storageBoxTypes[1].ID)
	assert.Equal(t, int64(3), storageBoxTypes[2].ID)
}
