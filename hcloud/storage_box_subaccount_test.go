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

func TestStorageBoxClientGetSubaccount(t *testing.T) {
	ctx, server, client := makeTestUtils(t)

	t.Run("GetSubaccount (ByID)", func(t *testing.T) {
		server.Expect([]mockutil.Request{
			{
				Method: "GET", Path: "/storage_boxes/42/subaccounts/13",
				Status: 200,
				JSONRaw: `{
					"subaccount": {
						"id": 13,
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

		assert.Equal(t, int64(13), subaccount.ID)
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

	t.Run("GetSubaccount (ByUsername)", func(t *testing.T) {
		server.Expect([]mockutil.Request{
			{
				Method: "GET", Path: "/storage_boxes/42/subaccounts?username=my-user",
				Status: 200,
				JSONRaw: `{
					"subaccounts": [
						{
							"id": 13,
							"username": "my-user",
							"storage_box": 42
						}
					]
				}`,
			},
		})

		storageBox := &StorageBox{ID: 42}

		subaccount, resp, err := client.StorageBox.GetSubaccountByUsername(ctx, storageBox, "my-user")
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NotNil(t, subaccount)

		assert.Equal(t, int64(13), subaccount.ID)
		assert.Equal(t, int64(42), subaccount.StorageBox.ID)
	})

	t.Run("GetSubbacount (IDOrName)", func(t *testing.T) {
		server.Expect([]mockutil.Request{
			{
				Method: "GET", Path: "/storage_boxes/42/subaccounts/13",
				Status: 200,
				JSONRaw: `{
					"subaccount": {
						"id": 13,
						"username": "foobar",
						"storage_box": 42
					}
				}`,
			},
			{
				Method: "GET", Path: "/storage_boxes/42/subaccounts?username=my-user",
				Status: 200,
				JSONRaw: `{
					"subaccounts": [
						{
							"id": 14,
							"username": "my-user",
							"storage_box": 42
						}
					]
				}`,
			},
		})

		storageBox := &StorageBox{ID: 42}

		subaccount, resp, err := client.StorageBox.GetSubaccount(ctx, storageBox, "13")
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NotNil(t, subaccount)

		assert.Equal(t, int64(13), subaccount.ID)
		assert.Equal(t, int64(42), subaccount.StorageBox.ID)

		subaccount, resp, err = client.StorageBox.GetSubaccount(ctx, storageBox, "my-user")
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NotNil(t, subaccount)

		assert.Equal(t, int64(14), subaccount.ID)
		assert.Equal(t, int64(42), subaccount.StorageBox.ID)

	})
}

func TestStorageBoxClientListSubaccounts(t *testing.T) {
	ctx, server, client := makeTestUtils(t)

	t.Run("ListSubaccounts", func(t *testing.T) {
		server.Expect([]mockutil.Request{
			{
				Method: "GET", Path: "/storage_boxes/42/subaccounts?label_selector=environment%3Dprod&sort=id%3Aasc",
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
			Sort:          []string{"id:asc"},
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

func TestStorageBoxClientCreateSubaccount(t *testing.T) {
	ctx, server, client := makeTestUtils(t)

	t.Run("CreateSubaccount (full)", func(t *testing.T) {
		server.Expect([]mockutil.Request{
			{
				Method: "POST", Path: "/storage_boxes/42/subaccounts",
				Status: 201,
				Want: func(t *testing.T, r *http.Request) {
					body, err := io.ReadAll(r.Body)
					require.NoError(t, err)

					expectedBody := `{
						"home_directory": "/home/my-user",
						"password": "my-password",
						"access_settings": {
							"reachable_externally": true,
							"readonly": false,
							"ssh_enabled": false,
							"samba_enabled": true,
							"webdav_enabled": true
						},
						"description": "This describes my subaccount",
						"labels": {
							"environment": "prod"
						}
					}`
					assert.JSONEq(t, expectedBody, string(body))
				},
				JSONRaw: `{
					"subaccount": {
						"id": 42,
						"storage_box": 42
					},
					"action": { "id": 12345 }
				}`,
			},
		})

		storageBox := &StorageBox{ID: 42}

		opts := StorageBoxSubaccountCreateOpts{
			HomeDirectory: Ptr("/home/my-user"),
			Password:      "my-password",
			AccessSettings: &StorageBoxSubaccountCreateOptsAccessSettings{
				ReachableExternally: Ptr(true),
				Readonly:            Ptr(false),
				SambaEnabled:        Ptr(true),
				SSHEnabled:          Ptr(false),
				WebDAVEnabled:       Ptr(true),
			},
			Description: Ptr("This describes my subaccount"),
			Labels: map[string]string{
				"environment": "prod",
			},
		}
		result, _, err := client.StorageBox.CreateSubaccount(ctx, storageBox, opts)
		require.NoError(t, err)
		require.NotNil(t, result)

		subaccount := result.Subaccount
		require.NotNil(t, subaccount)

		assert.Equal(t, int64(42), subaccount.ID)
		assert.Equal(t, int64(42), subaccount.StorageBox.ID)
	})
}

func TestStorageBoxClientUpdateSubaccount(t *testing.T) {
	ctx, server, client := makeTestUtils(t)

	t.Run("UpdateSubaccount (full)", func(t *testing.T) {
		server.Expect([]mockutil.Request{
			{
				Method: "PUT", Path: "/storage_boxes/42/subaccounts/13",
				Status: 200,
				Want: func(t *testing.T, r *http.Request) {
					body, err := io.ReadAll(r.Body)
					require.NoError(t, err)

					expectedBody := `{
						"labels": {
							"environment": "prod",
							"example.com/my": "label",
							"just-a-key": ""
						},
						"description": "Updated description"
					}`
					assert.JSONEq(t, expectedBody, string(body))
				},
				JSONRaw: `{
					"subaccount": {
						"id": 13,
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
						"description": "Updated description",
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

		subaccount := &StorageBoxSubaccount{
			ID: 13,
			StorageBox: &StorageBox{
				ID: 42,
			},
		}

		opts := StorageBoxSubaccountUpdateOpts{
			Description: Ptr("Updated description"),
			Labels: map[string]string{
				"environment":    "prod",
				"example.com/my": "label",
				"just-a-key":     "",
			},
		}

		result, _, err := client.StorageBox.UpdateSubaccount(ctx, subaccount, opts)

		require.NoError(t, err)

		assert.Equal(t, int64(13), result.ID)
	})
}

func TestStorageBoxClientDeleteSubaccount(t *testing.T) {
	ctx, server, client := makeTestUtils(t)

	t.Run("DeleteSubaccount", func(t *testing.T) {
		server.Expect([]mockutil.Request{
			{
				Method: "DELETE", Path: "/storage_boxes/42/subaccounts/13",
				Status:  201,
				JSONRaw: `{ "action": { "id": 5 } }`,
			},
		})

		subaccount := &StorageBoxSubaccount{
			ID: 13,
			StorageBox: &StorageBox{
				ID: 42,
			},
		}

		action, resp, err := client.StorageBox.DeleteSubaccount(ctx, subaccount)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NotNil(t, action)
	})
}

func TestStorageBoxClientResetSubaccountPassword(t *testing.T) {
	ctx, server, client := makeTestUtils(t)

	server.Expect([]mockutil.Request{
		{
			Method: "POST", Path: "/storage_boxes/42/subaccounts/13/actions/reset_subaccount_password",
			Status: 201,
			Want: func(t *testing.T, r *http.Request) {
				body, err := io.ReadAll(r.Body)
				require.NoError(t, err)

				assert.JSONEq(t, `{"password":"foobar"}`, string(body))
			},
			JSONRaw: `{ "action": { "id": 5 } }`,
		},
	})

	subaccount := &StorageBoxSubaccount{
		ID: 13,
		StorageBox: &StorageBox{
			ID: 42,
		},
	}

	opts := StorageBoxSubaccountResetPasswordOpts{
		Password: "foobar",
	}
	action, resp, err := client.StorageBox.ResetSubaccountPassword(ctx, subaccount, opts)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, action)
}

func TestStorageBoxSubbacountUpdateAccessSettings(t *testing.T) {
	ctx, server, client := makeTestUtils(t)

	server.Expect([]mockutil.Request{
		{
			Method: "POST", Path: "/storage_boxes/42/subaccounts/13/actions/update_access_settings",
			Status: 201,
			Want: func(t *testing.T, r *http.Request) {
				body, err := io.ReadAll(r.Body)
				require.NoError(t, err)

				expected := `{
					"samba_enabled": false,
					"ssh_enabled": true,
					"webdav_enabled": false,
					"readonly": false,
					"reachable_externally": true
				}`

				assert.JSONEq(t, expected, string(body))
			},
			JSONRaw: `{ "action": { "id": 5 } }`,
		},
	})

	subaccount := &StorageBoxSubaccount{
		ID: 13,
		StorageBox: &StorageBox{
			ID: 42,
		},
	}

	opts := StorageBoxSubaccountAccessSettingsUpdateOpts{
		SambaEnabled:        Ptr(false),
		SSHEnabled:          Ptr(true),
		WebDAVEnabled:       Ptr(false),
		Readonly:            Ptr(false),
		ReachableExternally: Ptr(true),
	}
	action, resp, err := client.StorageBox.UpdateSubaccountAccessSettings(ctx, subaccount, opts)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, action)
}

func TestStorageBoxSubbacountChangeHomeDirectory(t *testing.T) {
	ctx, server, client := makeTestUtils(t)

	server.Expect([]mockutil.Request{
		{
			Method: "POST", Path: "/storage_boxes/42/subaccounts/13/actions/change_home_directory",
			Status: 201,
			Want: func(t *testing.T, r *http.Request) {
				body, err := io.ReadAll(r.Body)
				require.NoError(t, err)

				expected := `{
					"home_directory": "/foobar"
				}`

				assert.JSONEq(t, expected, string(body))
			},
			JSONRaw: `{ "action": { "id": 5 } }`,
		},
	})

	subaccount := &StorageBoxSubaccount{
		ID: 13,
		StorageBox: &StorageBox{
			ID: 42,
		},
	}

	opts := StorageBoxSubaccountChangeHomeDirectoryOpts{
		HomeDirectory: "/foobar",
	}
	action, resp, err := client.StorageBox.ChangeSubaccountHomeDirectory(ctx, subaccount, opts)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, action)
}
