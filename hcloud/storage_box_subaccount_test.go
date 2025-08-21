package hcloud

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hetznercloud/hcloud-go/v2/hcloud/exp/mockutil"
)

func TestGetSubaccount(t *testing.T) {
	ctx, server, client := makeTestUtils(t)

	t.Run("GetSubaccount (ByID)", func(t *testing.T) {
		server.Expect([]mockutil.Request{
			{
				Method: "GET", Path: "/storage_boxes/42/subaccounts/13",
				Status: 200,
				JSONRaw: `{
					"subaccount": {
						"id": 42,
						"username": "my-user",
						"home_directory": "/home/my-user",
						"server": "my-server",
						"access_settings": {
							"reachable_externally": true,
							"readonly": false,
							"samba_enabled": true,
							"ssh_enabled": false,
							"webdav_enabled": true
						},
						"description": "This describes my subaccount",
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

		subaccount, _, err := client.StorageBox.GetSubaccountByID(ctx, storageBox, 13)
		require.NoError(t, err)
		require.NotNil(t, subaccount)

		assert.Equal(t, int64(42), subaccount.ID)
		assert.Equal(t, "my-user", subaccount.Username)
		assert.Equal(t, "/home/my-user", subaccount.HomeDirectory)
		assert.Equal(t, "my-server", subaccount.Server)
		assert.True(t, subaccount.AccessSettings.ReachableExternally)
		assert.False(t, subaccount.AccessSettings.Readonly)
		assert.True(t, subaccount.AccessSettings.SambaEnabled)
		assert.False(t, subaccount.AccessSettings.SSHEnabled)
		assert.True(t, subaccount.AccessSettings.WebDAVEnabled)
		assert.Equal(t, "This describes my subaccount", subaccount.Description)
		assert.Equal(t, time.Date(2025, 8, 21, 00, 00, 0, 0, time.UTC), subaccount.Created)
		assert.Equal(t, "prod", subaccount.Labels["environment"])
	})

	t.Run("GetSubaccount (not found)", func(t *testing.T) {
		server.Expect([]mockutil.Request{
			{
				Method: "GET", Path: "/storage_boxes/42/subaccounts/13",
				Status: 404,
				JSONRaw: `{
					"error": {
						"code": "not_found",
						"message": "Subaccount not found"
					}
				}`,
			},
		})

		storageBox := &StorageBox{ID: 42}

		subaccount, resp, err := client.StorageBox.GetSubaccountByID(ctx, storageBox, 13)
		require.NoError(t, err)
		require.NotNil(t, resp)

		assert.Nil(t, subaccount)
		assert.Equal(t, 404, resp.StatusCode)
	})
}

func TestListSubaccounts(t *testing.T) {
	ctx, server, client := makeTestUtils(t)

	t.Run("ListSubaccounts", func(t *testing.T) {
		server.Expect([]mockutil.Request{
			{
				Method: "GET", Path: "/storage_boxes/42/subaccounts?label_selector=environment%3Dprod",
				Status: 200,
				JSONRaw: `{
					"subaccounts": [
						{
							"id": 42,
							"username": "my-user",
							"home_directory": "/home/my-user",
							"server": "my-server",
							"access_settings": {
								"reachable_externally": true,
								"readonly": false,
								"samba_enabled": true,
								"ssh_enabled": false,
								"webdav_enabled": true
							},
							"description": "This describes my subaccount",
							"labels": {
								"environment": "prod",
								"example.com/my": "label",
								"just-a-key": ""
							},
							"created": "2025-08-21T00:00:00Z",
							"storage_box": 42
						}
					]
				}`,
			},
		})

		storageBox := &StorageBox{ID: 42}

		opts := StorageBoxSubaccountListOpts{
			LabelSelector: "environment=prod",
		}
		subaccounts, _, err := client.StorageBox.ListSubaccounts(ctx, storageBox, opts)
		require.NoError(t, err)
		require.Len(t, subaccounts, 1)

		subaccount := subaccounts[0]
		assert.Equal(t, int64(42), subaccount.ID)
		assert.Equal(t, "my-user", subaccount.Username)
	})

	t.Run("ListSubaccounts (AllSubaccountsWithOpts)", func(t *testing.T) {
		server.Expect([]mockutil.Request{
			{
				Method: "GET", Path: "/storage_boxes/42/subaccounts?label_selector=environment%3Dprod",
				Status: 200,
				JSONRaw: `{
					"subaccounts": [
						{
							"id": 42,
							"username": "my-user",
							"home_directory": "/home/my-user",
							"server": "my-server",
							"access_settings": {
								"reachable_externally": true,
								"readonly": false,
								"samba_enabled": true,
								"ssh_enabled": false,
								"webdav_enabled": true
							},
							"description": "This describes my subaccount",
							"labels": {
								"environment": "prod",
								"example.com/my": "label",
								"just-a-key": ""
							},
							"created": "2025-08-21T00:00:00Z",
							"storage_box": 42
						}
					]
				}`,
			},
		})

		storageBox := &StorageBox{ID: 42}

		opts := StorageBoxSubaccountListOpts{
			LabelSelector: "environment=prod",
		}
		subaccounts, err := client.StorageBox.AllSubaccountsWithOpts(ctx, storageBox, opts)
		require.NoError(t, err)
		require.Len(t, subaccounts, 1)

		subaccount := subaccounts[0]
		assert.Equal(t, int64(42), subaccount.ID)
		assert.Equal(t, "my-user", subaccount.Username)
	})
}
