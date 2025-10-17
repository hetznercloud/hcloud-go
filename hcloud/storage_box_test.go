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

func TestStorageBoxClientGet(t *testing.T) {
	t.Run("GetByID", func(t *testing.T) {
		ctx, server, client := makeTestUtils(t)

		server.Expect([]mockutil.Request{
			{
				Method: "GET", Path: "/storage_boxes/42",
				Status: 200,
				JSONRaw: `
				{
					"storage_box": {
						"id": 42,
						"username": "u123456",
						"status": "active",
						"name": "my-storage-box",
						"storage_box_type": {
							"id": 1,
							"name": "bx11",
							"description": "BX11"
						},
						"location": {
							"id": 1,
							"name": "fsn1",
							"description": "Falkenstein DC Park 1",
							"country": "DE",
							"city": "Falkenstein",
							"latitude": 50.47612,
							"longitude": 12.370071
						},
						"access_settings": {
							"reachable_externally": true,
							"samba_enabled": true,
							"ssh_enabled": true,
							"webdav_enabled": false,
							"zfs_enabled": false
						},
						"server": "u123456.your-storagebox.de",
						"system": "BX",
						"stats": {
							"size": 1073741824,
							"size_data": 536870912,
							"size_snapshots": 268435456
						},
						"labels": {
							"environment": "prod"
						},
						"protection": {
							"delete": true
						},
						"snapshot_plan": {
							"max_snapshots": 10,
							"minute": 0,
							"hour": 2,
							"day_of_week": null,
							"day_of_month": null
						},
						"created": "2023-01-01T12:00:00+00:00"
					}
				}`,
			},
		})

		storageBox, _, err := client.StorageBox.GetByID(ctx, 42)
		require.NoError(t, err)
		require.NotNil(t, storageBox, "no storage box")
		assert.Equal(t, int64(42), storageBox.ID, "unexpected storage box ID")
		assert.Equal(t, "my-storage-box", storageBox.Name, "unexpected storage box name")
		assert.Equal(t, StorageBoxStatusActive, storageBox.Status, "unexpected storage box status")
		assert.Equal(t, "u123456", storageBox.Username, "unexpected storage box username")
		assert.Equal(t, uint64(1073741824), storageBox.Stats.Size, "unexpected storage box size")
	})

	t.Run("GetByID (not found)", func(t *testing.T) {
		ctx, server, client := makeTestUtils(t)

		server.Expect([]mockutil.Request{
			{
				Method: "GET", Path: "/storage_boxes/1",
				Status:  404,
				JSONRaw: `{"error": {"code": "not_found"}}`,
			},
		})

		storageBox, _, err := client.StorageBox.GetByID(ctx, 1)
		require.NoError(t, err)
		assert.Nil(t, storageBox, "expected no storage box")
	})

	t.Run("GetByName", func(t *testing.T) {
		ctx, server, client := makeTestUtils(t)

		server.Expect([]mockutil.Request{
			{
				Method: "GET", Path: "/storage_boxes?name=foobar",
				Status: 200,
				JSONRaw: `{
					"storage_boxes": [{ "id": 1 }]
				}`,
			},
		})

		storageBox, _, err := client.StorageBox.GetByName(ctx, "foobar")
		require.NoError(t, err)
		require.NotNil(t, storageBox, "no storage box")
		assert.Equal(t, int64(1), storageBox.ID, "unexpected storage box ID")
	})
}

func TestStorageBoxClientList(t *testing.T) {
	ctx, server, client := makeTestUtils(t)

	server.Expect([]mockutil.Request{
		{
			Method: "GET", Path: "/storage_boxes?page=2&per_page=50",
			Status: 200,
			JSONRaw: `{
				"storage_boxes": [
					{ "id": 1 },
					{ "id": 2 }
				]
			}`,
		},
	})

	opts := StorageBoxListOpts{}
	opts.Page = 2
	opts.PerPage = 50
	storageBoxes, _, err := client.StorageBox.List(ctx, opts)
	require.NoError(t, err)
	assert.Len(t, storageBoxes, 2, "expected 2 storage boxes")
	assert.Equal(t, int64(1), storageBoxes[0].ID)
	assert.Equal(t, int64(2), storageBoxes[1].ID)
}

func TestStorageBoxClientAll(t *testing.T) {
	ctx, server, client := makeTestUtils(t)

	server.Expect([]mockutil.Request{
		{
			Method: "GET", Path: "/storage_boxes?page=1&per_page=2",
			Status: 200,
			JSONRaw: `{
				"storage_boxes": [
					{ "id": 1 },
					{ "id": 2 }
				],
				"meta": {"pagination": {"page": 1, "last_page": 2, "per_page": 2, "next_page": 2, "previous_page": null, "total_entries": 3}}
			}`,
		},
		{
			Method: "GET", Path: "/storage_boxes?page=2&per_page=2",
			Status: 200,
			JSONRaw: `{
				"storage_boxes": [
					{ "id": 3 }
				],
				"meta": {"pagination": {"page": 2, "last_page": 2, "per_page": 2, "next_page": null, "previous_page": 1, "total_entries": 3}}
			}`,
		},
	})

	storageBoxes, err := client.StorageBox.AllWithOpts(ctx, StorageBoxListOpts{
		ListOpts: ListOpts{PerPage: 2},
	})
	require.NoError(t, err, "fetching all storage boxes failed")
	assert.Len(t, storageBoxes, 3, "expected 3 storage boxes")
	assert.Equal(t, int64(1), storageBoxes[0].ID)
	assert.Equal(t, int64(2), storageBoxes[1].ID)
	assert.Equal(t, int64(3), storageBoxes[2].ID)
}

func TestStorageBoxClientCreate(t *testing.T) {
	ctx, server, client := makeTestUtils(t)

	server.Expect([]mockutil.Request{
		{
			Method: "POST", Path: "/storage_boxes",
			Want: func(t *testing.T, r *http.Request) {
				bodyBytes, err := io.ReadAll(r.Body)
				require.NoError(t, err)
				assert.JSONEq(t, `{
					"name": "my-new-storage-box",
					"storage_box_type": 1,
					"location": "fsn1",
					"password": "secretpassword123",
					"labels": {"env": "test"},
					"ssh_keys": ["ssh-rsa AAAAB3NzaC1yc2E..."],
					"access_settings": {
						"reachable_externally": true,
						"ssh_enabled": false
					}
				}`, string(bodyBytes))
			},
			Status: 201,
			JSONRaw: `{
				"action": { "id": 1 },
				"storage_box": {
					"id": 42,
					"status": "initializing",
					"username": null,
					"server": null,
					"system": null
				}
			}`,
		},
	})

	opts := StorageBoxCreateOpts{
		Name:           "my-new-storage-box",
		StorageBoxType: &StorageBoxType{ID: 1},
		Location:       &Location{Name: "fsn1"},
		Password:       "secretpassword123",
		Labels:         map[string]string{"env": "test"},
		SSHKeys:        []string{"ssh-rsa AAAAB3NzaC1yc2E..."},
		AccessSettings: &StorageBoxCreateOptsAccessSettings{
			ReachableExternally: Ptr(true),
			SSHEnabled:          Ptr(false),
		},
	}

	result, _, err := client.StorageBox.Create(ctx, opts)
	require.NoError(t, err)
	require.NotNil(t, result.Action, "no action returned")
	require.NotNil(t, result.StorageBox, "no storage box returned")

	assert.Equal(t, int64(1), result.Action.ID, "unexpected action ID")
	assert.Equal(t, int64(42), result.StorageBox.ID, "unexpected storage box ID")
}

func TestStorageBoxClientUpdate(t *testing.T) {
	ctx, server, client := makeTestUtils(t)

	server.Expect([]mockutil.Request{
		{
			Method: "PUT", Path: "/storage_boxes/42",
			Want: func(t *testing.T, r *http.Request) {
				bodyBytes, err := io.ReadAll(r.Body)
				require.NoError(t, err)
				assert.JSONEq(t, `{
					"name": "updated-storage-box",
					"labels": {"env": "prod"}
				}`, string(bodyBytes))
			},
			Status: 200,
			JSONRaw: `{
				"storage_box": { "id": 42 }
			}`,
		},
	})

	storageBox := &StorageBox{ID: 42}
	opts := StorageBoxUpdateOpts{
		Name:   "updated-storage-box",
		Labels: map[string]string{"env": "prod"},
	}

	updatedStorageBox, _, err := client.StorageBox.Update(ctx, storageBox, opts)
	require.NoError(t, err)
	require.NotNil(t, updatedStorageBox, "no storage box returned")
	assert.Equal(t, int64(42), updatedStorageBox.ID, "unexpected storage box ID")
}

func TestStorageBoxClientDelete(t *testing.T) {
	ctx, server, client := makeTestUtils(t)

	server.Expect([]mockutil.Request{
		{
			Method: "DELETE", Path: "/storage_boxes/42",
			Status: 200,
			JSONRaw: `{
				"action": { "id": 1 }
			}`,
		},
	})

	storageBox := &StorageBox{ID: 42}

	result, _, err := client.StorageBox.Delete(ctx, storageBox)
	require.NoError(t, err)
	require.NotNil(t, result.Action, "no action returned")
	assert.Equal(t, int64(1), result.Action.ID, "unexpected action ID")
}

func TestStorageBoxClientFolders(t *testing.T) {
	ctx, server, client := makeTestUtils(t)

	t.Run("folders", func(t *testing.T) {
		server.Expect([]mockutil.Request{
			{
				Method: "GET", Path: "/storage_boxes/42/folders?",
				Status: 200,
				JSONRaw: `{
					"folders": ["foo", "bar"]
				}`,
			},
		})

		storageBox := &StorageBox{ID: 42}

		result, _, err := client.StorageBox.Folders(ctx, storageBox, StorageBoxFoldersOpts{})
		require.NoError(t, err)
		require.NotNil(t, result.Folders, "no result returned")

		assert.Len(t, result.Folders, 2, "unexpected number of folders")
		assert.Equal(t, "foo", result.Folders[0], "unexpected first folder")
		assert.Equal(t, "bar", result.Folders[1], "unexpected second folder")
	})

	t.Run("sub folders", func(t *testing.T) {
		server.Expect([]mockutil.Request{
			{
				Method: "GET", Path: "/storage_boxes/42/folders?path=%2Ffoo",
				Status: 200,
				JSONRaw: `{
				"folders": ["subfoo", "subbar"]
			}`,
			},
		})

		storageBox := &StorageBox{ID: 42}

		result, _, err := client.StorageBox.Folders(ctx, storageBox, StorageBoxFoldersOpts{Path: "/foo"})
		require.NoError(t, err)
		require.NotNil(t, result.Folders, "no result returned")

		assert.Len(t, result.Folders, 2, "unexpected number of folders")
		assert.Equal(t, "subfoo", result.Folders[0], "unexpected first folder")
	})
}

func TestStorageBoxClientChangeProtection(t *testing.T) {
	ctx, server, client := makeTestUtils(t)

	server.Expect([]mockutil.Request{
		{
			Method: "POST", Path: "/storage_boxes/42/actions/change_protection",
			Status: 201,
			Want: func(t *testing.T, r *http.Request) {
				body, err := io.ReadAll(r.Body)
				require.NoError(t, err, "failed to read request body")

				assert.JSONEq(t, `{ "delete": true }`, string(body))
			},
			JSONRaw: `{ "action": { "id": 13 } }`,
		},
	})

	storageBox := &StorageBox{ID: 42}

	opts := StorageBoxChangeProtectionOpts{Delete: Ptr(true)}
	action, _, err := client.StorageBox.ChangeProtection(ctx, storageBox, opts)
	require.NoError(t, err, "ChangeProtection failed")
	require.NotNil(t, action, "no action returned")

	assert.Equal(t, int64(13), action.ID, "unexpected action ID")
}

func TestStorageBoxResetPassword(t *testing.T) {
	ctx, server, client := makeTestUtils(t)

	server.Expect([]mockutil.Request{
		{
			Method: "POST", Path: "/storage_boxes/42/actions/reset_password",
			Status: 201,
			Want: func(t *testing.T, r *http.Request) {
				body, err := io.ReadAll(r.Body)
				require.NoError(t, err, "failed to read request body")

				assert.JSONEq(t, `{ "password": "newpassword123" }`, string(body))
			},
			JSONRaw: `{ "action": { "id": 13 } }`,
		},
	})

	storageBox := &StorageBox{ID: 42}

	opts := StorageBoxResetPasswordOpts{Password: "newpassword123"}
	action, _, err := client.StorageBox.ResetPassword(ctx, storageBox, opts)
	require.NoError(t, err, "ResetPassword failed")
	require.NotNil(t, action, "no action returned")

	assert.Equal(t, int64(13), action.ID, "unexpected action ID")
}

func TestStorageBoxUpdateAccessSettings(t *testing.T) {
	ctx, server, client := makeTestUtils(t)
	storageBox := &StorageBox{ID: 42}

	t.Run("UpdateAccessSettings (all)", func(t *testing.T) {
		server.Expect([]mockutil.Request{
			{
				Method: "POST", Path: "/storage_boxes/42/actions/update_access_settings",
				Status: 201,
				Want: func(t *testing.T, r *http.Request) {
					body, err := io.ReadAll(r.Body)
					require.NoError(t, err, "failed to read request body")

					expected := `{
						"samba_enabled": true,
						"ssh_enabled": false,
						"webdav_enabled": true,
						"zfs_enabled": false,
						"reachable_externally": true
					}`

					assert.JSONEq(t, expected, string(body), "unexpected request body")
				},
				JSONRaw: `{ "action": { "id": 13 } }`,
			},
		})

		opts := StorageBoxUpdateAccessSettingsOpts{
			SambaEnabled:        Ptr(true),
			SSHEnabled:          Ptr(false),
			WebDAVEnabled:       Ptr(true),
			ZFSEnabled:          Ptr(false),
			ReachableExternally: Ptr(true),
		}
		action, _, err := client.StorageBox.UpdateAccessSettings(ctx, storageBox, opts)
		require.NoError(t, err, "UpdateAccessSettings failed")
		require.NotNil(t, action, "no action returned")

		assert.Equal(t, int64(13), action.ID, "unexpected action ID")
	})

	t.Run("UpdateAccessSettings (some)", func(t *testing.T) {
		server.Expect([]mockutil.Request{
			{
				Method: "POST", Path: "/storage_boxes/42/actions/update_access_settings",
				Status: 201,
				Want: func(t *testing.T, r *http.Request) {
					body, err := io.ReadAll(r.Body)
					require.NoError(t, err, "failed to read request body")

					expected := `{
						"samba_enabled": true,
						"ssh_enabled": false
					}`

					assert.JSONEq(t, expected, string(body), "unexpected request body")
				},
				JSONRaw: `{ "action": { "id": 13 } }`,
			},
		})

		opts := StorageBoxUpdateAccessSettingsOpts{
			SambaEnabled: Ptr(true),
			SSHEnabled:   Ptr(false),
		}
		action, _, err := client.StorageBox.UpdateAccessSettings(ctx, storageBox, opts)
		require.NoError(t, err, "UpdateAccessSettings failed")
		require.NotNil(t, action, "no action returned")
	})
}

func TestStorageBoxRollbackSnapshot(t *testing.T) {
	ctx, server, client := makeTestUtils(t)

	server.Expect([]mockutil.Request{
		{
			Method: "POST", Path: "/storage_boxes/42/actions/rollback_snapshot",
			Status: 201,
			Want: func(t *testing.T, r *http.Request) {
				body, err := io.ReadAll(r.Body)
				require.NoError(t, err, "failed to read request body")

				assert.JSONEq(t, `{ "snapshot": 10 }`, string(body), "unexpected request body")
			},
			JSONRaw: `{ "action": { "id": 13 } }`,
		},
	})

	storageBox := &StorageBox{ID: 42}

	opts := StorageBoxRollbackSnapshotOpts{
		Snapshot: &StorageBoxSnapshot{ID: 10},
	}
	action, _, err := client.StorageBox.RollbackSnapshot(ctx, storageBox, opts)
	require.NoError(t, err, "RollbackSnapshot failed")
	require.NotNil(t, action, "no action returned")
}

func TestStorageBoxEnableSnapshotPlan(t *testing.T) {
	ctx, server, client := makeTestUtils(t)

	server.Expect([]mockutil.Request{
		{
			Method: "POST", Path: "/storage_boxes/42/actions/enable_snapshot_plan",
			Status: 201,
			Want: func(t *testing.T, r *http.Request) {
				body, err := io.ReadAll(r.Body)
				require.NoError(t, err, "failed to read request body")

				expectedBody := `{
					"max_snapshots": 10,
					"minute": 5,
					"hour": 6,
					"day_of_week": 7,
					"day_of_month": null
				}`
				assert.JSONEq(t, expectedBody, string(body))
			},
			JSONRaw: `{ "action": { "id": 13 } }`,
		},
	})

	storageBox := &StorageBox{ID: 42}

	opts := StorageBoxEnableSnapshotPlanOpts{
		MaxSnapshots: 10,
		Minute:       5,
		Hour:         6,
		DayOfWeek:    Ptr(time.Sunday),
	}
	action, _, err := client.StorageBox.EnableSnapshotPlan(ctx, storageBox, opts)
	require.NoError(t, err, "RollbackSnapshot failed")
	require.NotNil(t, action, "no action returned")
}

func TestStorageBoxDisableSnapshotPlan(t *testing.T) {
	ctx, server, client := makeTestUtils(t)

	server.Expect([]mockutil.Request{
		{
			Method: "POST", Path: "/storage_boxes/42/actions/disable_snapshot_plan",
			Status:  201,
			JSONRaw: `{ "action": { "id": 13 } }`,
		},
	})

	storageBox := &StorageBox{ID: 42}

	action, _, err := client.StorageBox.DisableSnapshotPlan(ctx, storageBox)
	require.NoError(t, err, "RollbackSnapshot failed")
	require.NotNil(t, action, "no action returned")
}
