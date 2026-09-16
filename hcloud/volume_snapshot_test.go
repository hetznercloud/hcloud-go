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

func TestVolumeClientGetSnapshot(t *testing.T) {
	ctx, server, client := makeTestUtils(t)
	server.Expect([]mockutil.Request{{
		Method: "GET", Path: "/volume_snapshots/13", Status: 200,
		JSONRaw: `{
			"snapshot": {
				"id": 13,
				"name": "pvc-snapshot",
				"description": "backup",
				"status": "available",
				"volume": 42,
				"location": {"id": 1, "name": "fsn1"},
				"size": 10,
				"labels": {"managed-by": "csi-driver"},
				"created": "2026-09-16T00:00:00Z"
			}
		}`,
	}})

	snapshot, _, err := client.Volume.GetSnapshotByID(ctx, 13)
	require.NoError(t, err)
	require.NotNil(t, snapshot)
	assert.Equal(t, int64(13), snapshot.ID)
	assert.Equal(t, VolumeSnapshotStatusAvailable, snapshot.Status)
	assert.Equal(t, int64(42), snapshot.Volume.ID)
	assert.Equal(t, "fsn1", snapshot.Location.Name)
	assert.Equal(t, time.Date(2026, 9, 16, 0, 0, 0, 0, time.UTC), snapshot.Created)
}

func TestVolumeClientListSnapshots(t *testing.T) {
	ctx, server, client := makeTestUtils(t)
	server.Expect([]mockutil.Request{{
		Method: "GET", Path: "/volume_snapshots?label_selector=managed-by%3Dcsi-driver&name=pvc-snapshot&page=1&per_page=50&sort=created%3Adesc&status=available&volume=42", Status: 200,
		JSONRaw: `{ "snapshots": [{ "id": 13 }] }`,
	}})

	snapshots, err := client.Volume.AllSnapshotsWithOpts(ctx, VolumeSnapshotListOpts{
		Name:          "pvc-snapshot",
		Volume:        &Volume{ID: 42},
		Status:        []VolumeSnapshotStatus{VolumeSnapshotStatusAvailable},
		LabelSelector: "managed-by=csi-driver",
		Sort:          []string{"created:desc"},
	})
	require.NoError(t, err)
	require.Len(t, snapshots, 1)
	assert.Equal(t, int64(13), snapshots[0].ID)
}

func TestVolumeClientCreateSnapshot(t *testing.T) {
	ctx, server, client := makeTestUtils(t)
	server.Expect([]mockutil.Request{{
		Method: "POST", Path: "/volumes/42/actions/create_snapshot", Status: 201,
		Want: func(t *testing.T, r *http.Request) {
			body, err := io.ReadAll(r.Body)
			require.NoError(t, err)
			assert.JSONEq(t, `{
				"name": "pvc-snapshot",
				"description": "created by CSI",
				"labels": {"managed-by": "csi-driver"}
			}`, string(body))
		},
		JSONRaw: `{ "snapshot": { "id": 13 }, "action": { "id": 7 } }`,
	}})

	result, _, err := client.Volume.CreateSnapshot(ctx, &Volume{ID: 42}, VolumeSnapshotCreateOpts{
		Name:        "pvc-snapshot",
		Description: "created by CSI",
		Labels:      map[string]string{"managed-by": "csi-driver"},
	})
	require.NoError(t, err)
	assert.Equal(t, int64(13), result.Snapshot.ID)
	assert.Equal(t, int64(7), result.Action.ID)
}

func TestVolumeClientDeleteSnapshot(t *testing.T) {
	ctx, server, client := makeTestUtils(t)
	server.Expect([]mockutil.Request{{
		Method: "DELETE", Path: "/volume_snapshots/13", Status: 200,
		JSONRaw: `{ "action": { "id": 7 } }`,
	}})

	action, _, err := client.Volume.DeleteSnapshot(ctx, &VolumeSnapshot{ID: 13})
	require.NoError(t, err)
	assert.Equal(t, int64(7), action.ID)
}

func TestVolumeClientCreateFromSnapshot(t *testing.T) {
	ctx, server, client := makeTestUtils(t)
	server.Expect([]mockutil.Request{{
		Method: "POST", Path: "/volumes", Status: 201,
		Want: func(t *testing.T, r *http.Request) {
			body, err := io.ReadAll(r.Body)
			require.NoError(t, err)
			assert.JSONEq(t, `{
				"name": "restored-volume",
				"size": 10,
				"snapshot": 13,
				"location": "fsn1"
			}`, string(body))
		},
		JSONRaw: `{ "volume": { "id": 43 }, "action": { "id": 8 }, "next_actions": [] }`,
	}})

	result, _, err := client.Volume.Create(ctx, VolumeCreateOpts{
		Name:     "restored-volume",
		Size:     10,
		Snapshot: &VolumeSnapshot{ID: 13},
		Location: &Location{Name: "fsn1"},
	})
	require.NoError(t, err)
	assert.Equal(t, int64(43), result.Volume.ID)
}
