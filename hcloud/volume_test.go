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

func TestVolumeClientGetByID(t *testing.T) {
	t.Run("GetByID", func(t *testing.T) {
		ctx, server, client := makeTestUtils(t)

		server.Expect([]mockutil.Request{
			{
				Method: "GET", Path: "/volumes/1",
				Status: 200,
				JSONRaw: `{
					"volume": {
						"id": 1,
						"name": "volume1",
						"created": "2026-10-27T11:17:32Z",
						"labels": {
							"key": "value"
						},
						"location": {
							"id": 3,
							"name": "hel1",
							"description": "Helsinki DC Park 1",
							"city": "Helsinki",
							"country": "FI",
							"latitude": 60.169855,
							"longitude": 24.938379,
							"network_zone": "eu-central"
						},
						"protection": {
							"delete": true
						},
						"status": "available",
						"size": 10,
						"format": "xfs",
						"server": 2,
						"linux_device": "/dev/disk/by-id/scsi-0HC_Volume_1"
					}
				}`,
			},
		})

		volume, resp, err := client.Volume.GetByID(ctx, 1)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NotNil(t, volume)
		assert.Equal(t, int64(1), volume.ID)
		assert.Equal(t, "volume1", volume.Name)
		assert.Equal(t, "2026-10-27T11:17:32Z", volume.Created.Format(time.RFC3339))
		assert.Equal(t, map[string]string{"key": "value"}, volume.Labels)
		require.NotNil(t, volume.Location)
		assert.Equal(t, int64(3), volume.Location.ID)
		assert.Equal(t, "hel1", volume.Location.Name)
		assert.Equal(t, VolumeProtection{Delete: true}, volume.Protection)
		assert.Equal(t, VolumeStatusAvailable, volume.Status)
		assert.Equal(t, 10, volume.Size)
		assert.Equal(t, new("xfs"), volume.Format)
		require.NotNil(t, volume.Server)
		assert.Equal(t, int64(2), volume.Server.ID)
		assert.Equal(t, "/dev/disk/by-id/scsi-0HC_Volume_1", volume.LinuxDevice)
	})

	t.Run("GetByID (not found)", func(t *testing.T) {
		ctx, server, client := makeTestUtils(t)

		server.Expect([]mockutil.Request{
			{
				Method: "GET", Path: "/volumes/1",
				Status:  404,
				JSONRaw: `{ "error": { "code": "not_found" }}`,
			},
		})

		volume, resp, err := client.Volume.GetByID(ctx, 1)
		require.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Nil(t, volume)
	})

	t.Run("Get", func(t *testing.T) {
		ctx, server, client := makeTestUtils(t)

		server.Expect([]mockutil.Request{
			{
				Method: "GET", Path: "/volumes/1",
				Status:  200,
				JSONRaw: `{ "volume": { "id": 1 }}`,
			},
		})

		volume, resp, err := client.Volume.Get(ctx, "1")
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NotNil(t, volume)
		assert.Equal(t, int64(1), volume.ID)
	})
}

func TestVolumeClientGetByName(t *testing.T) {
	t.Run("GetByName", func(t *testing.T) {
		ctx, server, client := makeTestUtils(t)

		server.Expect([]mockutil.Request{
			{
				Method: "GET", Path: "/volumes?name=my-volume",
				Status:  200,
				JSONRaw: `{ "volumes": [{ "id": 1 }]}`,
			},
		})

		volume, resp, err := client.Volume.GetByName(ctx, "my-volume")
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NotNil(t, volume)
		assert.Equal(t, int64(1), volume.ID)
	})

	t.Run("GetByName (not found)", func(t *testing.T) {
		ctx, server, client := makeTestUtils(t)

		server.Expect([]mockutil.Request{
			{
				Method: "GET", Path: "/volumes?name=my-volume",
				Status:  200,
				JSONRaw: `{ "volumes": []}`,
			},
		})

		volume, resp, err := client.Volume.GetByName(ctx, "my-volume")
		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Nil(t, volume)
	})

	t.Run("GetByName (empty name)", func(t *testing.T) {
		ctx, _, client := makeTestUtils(t)

		volume, resp, err := client.Volume.GetByName(ctx, "")
		require.NoError(t, err)
		require.Nil(t, resp)
		require.Nil(t, volume)
	})

	t.Run("Get", func(t *testing.T) {
		ctx, server, client := makeTestUtils(t)

		server.Expect([]mockutil.Request{
			{
				Method: "GET", Path: "/volumes?name=my-volume",
				Status:  200,
				JSONRaw: `{ "volumes": [{ "id": 1 }]}`,
			},
		})

		volume, resp, err := client.Volume.Get(ctx, "my-volume")
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NotNil(t, volume)
		assert.Equal(t, int64(1), volume.ID)
	})
}

func TestVolumeClientGet(t *testing.T) {
	t.Run("Get (numeric name)", func(t *testing.T) {
		ctx, server, client := makeTestUtils(t)

		server.Expect([]mockutil.Request{
			{
				Method: "GET", Path: "/volumes/123",
				Status:  404,
				JSONRaw: `{ "error": { "code": "not_found" }}`,
			},
			{
				Method: "GET", Path: "/volumes?name=123",
				Status:  200,
				JSONRaw: `{ "volumes": [{ "id": 1 }]}`,
			},
		})

		volume, resp, err := client.Volume.Get(ctx, "123")
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NotNil(t, volume)
		assert.Equal(t, int64(1), volume.ID)
	})

	t.Run("Get (empty name)", func(t *testing.T) {
		ctx, _, client := makeTestUtils(t)

		volume, resp, err := client.Volume.Get(ctx, "")
		require.NoError(t, err)
		require.Nil(t, resp)
		require.Nil(t, volume)
	})
}

func TestVolumeClientCreate(t *testing.T) {
	t.Run("Create (minimal with server)", func(t *testing.T) {
		ctx, server, client := makeTestUtils(t)

		server.Expect([]mockutil.Request{
			{
				Method: "POST", Path: "/volumes",
				Want: func(t *testing.T, r *http.Request) {
					body, err := io.ReadAll(r.Body)
					require.NoError(t, err)
					assert.JSONEq(t, `{
						"name": "my-volume",
						"size": 42,
						"server": 1
					}`, string(body))
				},
				Status: 201,
				JSONRaw: `{
					"volume": { "id": 1 },
					"action": { "id": 2 },
					"next_actions": [{ "id": 3 }]
				}`,
			},
		})

		result, resp, err := client.Volume.Create(ctx, VolumeCreateOpts{
			Name:   "my-volume",
			Size:   42,
			Server: &Server{ID: 1},
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, int64(1), result.Volume.ID)
		assert.Equal(t, int64(2), result.Action.ID)
		require.Len(t, result.NextActions, 1)
		assert.Equal(t, int64(3), result.NextActions[0].ID)
	})

	t.Run("Create (full with server)", func(t *testing.T) {
		ctx, server, client := makeTestUtils(t)

		server.Expect([]mockutil.Request{
			{
				Method: "POST", Path: "/volumes",
				Want: func(t *testing.T, r *http.Request) {
					body, err := io.ReadAll(r.Body)
					require.NoError(t, err)
					assert.JSONEq(t, `{
						"name": "my-volume",
						"size": 42,
						"server": 1,
						"automount": true,
						"format": "xfs",
						"labels": { "key": "value" }
					}`, string(body))
				},
				Status: 201,
				JSONRaw: `{
					"volume": { "id": 1 },
					"action": { "id": 2 },
					"next_actions": [{ "id": 3 }]
				}`,
			},
		})

		result, resp, err := client.Volume.Create(ctx, VolumeCreateOpts{
			Name:      "my-volume",
			Size:      42,
			Server:    &Server{ID: 1},
			Automount: new(true),
			Format:    new("xfs"),
			Labels:    map[string]string{"key": "value"},
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, int64(1), result.Volume.ID)
		assert.Equal(t, int64(2), result.Action.ID)
		require.Len(t, result.NextActions, 1)
		assert.Equal(t, int64(3), result.NextActions[0].ID)
	})

	t.Run("Create (with location)", func(t *testing.T) {
		ctx, server, client := makeTestUtils(t)

		server.Expect([]mockutil.Request{
			{
				Method: "POST", Path: "/volumes",
				Want: func(t *testing.T, r *http.Request) {
					body, err := io.ReadAll(r.Body)
					require.NoError(t, err)
					assert.JSONEq(t, `{
						"name": "my-volume",
						"size": 42,
						"location": 1,
						"labels": {"key": "value"}
					}`, string(body))
				},
				Status: 201,
				JSONRaw: `{
					"volume": { "id": 1 },
					"action": { "id": 2 },
					"next_actions": []
				}`,
			},
		})

		result, resp, err := client.Volume.Create(ctx, VolumeCreateOpts{
			Name:     "my-volume",
			Size:     42,
			Location: &Location{ID: 1},
			Labels:   map[string]string{"key": "value"},
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, int64(1), result.Volume.ID)
		assert.Equal(t, int64(2), result.Action.ID)
		require.Empty(t, result.NextActions)
	})
}

func TestVolumeClientUpdate(t *testing.T) {
	t.Run("minimal", func(t *testing.T) {
		ctx, server, client := makeTestUtils(t)

		server.Expect([]mockutil.Request{
			{
				Method: "PUT", Path: "/volumes/1",
				Want: func(t *testing.T, r *http.Request) {
					body, err := io.ReadAll(r.Body)
					require.NoError(t, err)
					assert.JSONEq(t, `{}`, string(body))
				},
				Status:  200,
				JSONRaw: `{"volume": {"id": 1}}`,
			},
		})

		volume, resp, err := client.Volume.Update(ctx, &Volume{ID: 1}, VolumeUpdateOpts{})
		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, int64(1), volume.ID)
	})

	t.Run("full", func(t *testing.T) {
		ctx, server, client := makeTestUtils(t)

		server.Expect([]mockutil.Request{
			{
				Method: "PUT", Path: "/volumes/1",
				Want: func(t *testing.T, r *http.Request) {
					body, err := io.ReadAll(r.Body)
					require.NoError(t, err)
					assert.JSONEq(t, `{
						"name": "volume1",
						"labels": {"key": "value"}
					}`, string(body))
				},
				Status:  200,
				JSONRaw: `{"volume": {"id": 1}}`,
			},
		})

		result, resp, err := client.Volume.Update(ctx, &Volume{ID: 1}, VolumeUpdateOpts{
			Name:   "volume1",
			Labels: map[string]string{"key": "value"},
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, int64(1), result.ID)
	})
}

func TestVolumeClientDelete(t *testing.T) {
	ctx, server, client := makeTestUtils(t)

	server.Expect([]mockutil.Request{
		{
			Method: "DELETE", Path: "/volumes/1",
			Status: 204,
		},
	})

	resp, err := client.Volume.Delete(ctx, &Volume{ID: 1})
	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestVolumeClientAttachWithOpts(t *testing.T) {
	t.Run("minimal", func(t *testing.T) {
		ctx, server, client := makeTestUtils(t)

		server.Expect([]mockutil.Request{
			{
				Method: "POST", Path: "/volumes/1/actions/attach",
				Want: func(t *testing.T, r *http.Request) {
					body, err := io.ReadAll(r.Body)
					require.NoError(t, err)
					assert.JSONEq(t, `{
						"server": 2
					}`, string(body))
				},
				Status:  201,
				JSONRaw: `{"action": {"id": 3}}`,
			},
		})

		action, resp, err := client.Volume.AttachWithOpts(ctx, &Volume{ID: 1}, VolumeAttachOpts{
			Server: &Server{ID: 2},
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, int64(3), action.ID)
	})

	t.Run("full", func(t *testing.T) {
		ctx, server, client := makeTestUtils(t)

		server.Expect([]mockutil.Request{
			{
				Method: "POST", Path: "/volumes/1/actions/attach",
				Want: func(t *testing.T, r *http.Request) {
					body, err := io.ReadAll(r.Body)
					require.NoError(t, err)
					assert.JSONEq(t, `{
						"server": 2,
						"automount": false
					}`, string(body))
				},
				Status:  201,
				JSONRaw: `{"action": {"id": 3}}`,
			},
		})

		action, resp, err := client.Volume.AttachWithOpts(ctx, &Volume{ID: 1}, VolumeAttachOpts{
			Server:    &Server{ID: 2},
			Automount: new(false),
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, int64(3), action.ID)
	})
}

func TestVolumeClientAttach(t *testing.T) {
	ctx, server, client := makeTestUtils(t)

	server.Expect([]mockutil.Request{
		{
			Method: "POST", Path: "/volumes/1/actions/attach",
			Want: func(t *testing.T, r *http.Request) {
				body, err := io.ReadAll(r.Body)
				require.NoError(t, err)
				assert.JSONEq(t, `{
					"server": 2
				}`, string(body))
			},
			Status:  201,
			JSONRaw: `{"action": {"id": 3}}`,
		},
	})

	action, resp, err := client.Volume.Attach(ctx, &Volume{ID: 1}, &Server{ID: 2})
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, int64(3), action.ID)
}

func TestVolumeClientDetach(t *testing.T) {
	ctx, server, client := makeTestUtils(t)

	server.Expect([]mockutil.Request{
		{
			Method: "POST", Path: "/volumes/1/actions/detach",
			Want: func(t *testing.T, r *http.Request) {
				body, err := io.ReadAll(r.Body)
				require.NoError(t, err)
				require.Empty(t, body)
			},
			Status:  201,
			JSONRaw: `{"action": {"id": 3}}`,
		},
	})

	action, resp, err := client.Volume.Detach(ctx, &Volume{ID: 1})
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, int64(3), action.ID)
}

func TestVolumeClientChangeProtection(t *testing.T) {
	ctx, server, client := makeTestUtils(t)

	server.Expect([]mockutil.Request{
		{
			Method: "POST", Path: "/volumes/1/actions/change_protection",
			Want: func(t *testing.T, r *http.Request) {
				body, err := io.ReadAll(r.Body)
				require.NoError(t, err)
				assert.JSONEq(t, `{
					"delete": true
				}`, string(body))
			},
			Status:  201,
			JSONRaw: `{"action": {"id": 3}}`,
		},
	})

	action, resp, err := client.Volume.ChangeProtection(ctx, &Volume{ID: 1}, VolumeChangeProtectionOpts{
		Delete: new(true),
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, int64(3), action.ID)
}

func TestVolumeClientResize(t *testing.T) {
	ctx, server, client := makeTestUtils(t)

	server.Expect([]mockutil.Request{
		{
			Method: "POST", Path: "/volumes/1/actions/resize",
			Want: func(t *testing.T, r *http.Request) {
				body, err := io.ReadAll(r.Body)
				require.NoError(t, err)
				assert.JSONEq(t, `{
					"size": 50
				}`, string(body))
			},
			Status:  201,
			JSONRaw: `{"action": {"id": 3}}`,
		},
	})

	action, resp, err := client.Volume.Resize(ctx, &Volume{ID: 1}, 50)
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, int64(3), action.ID)
}
