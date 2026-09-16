# Block volume snapshot API proposal

This proposal supplies the missing provider contract required by CSI
`CREATE_DELETE_SNAPSHOT` implementations. It deliberately does not reuse
server images or Storage Box snapshots: neither resource represents a
point-in-time copy that can restore a Hetzner Cloud block volume.

## Resource

`volume_snapshot` has these fields:

- `id`, `name`, `description`, `labels`, and `created`, following existing
  resource conventions;
- `status`: `creating` or `available`;
- `volume`: the source volume ID, retained even if the source volume is later
  deleted;
- `location`: the snapshot's restore location;
- `size`: the minimum restored volume size in GB.

## Endpoints

| Method | Path | Behavior |
| --- | --- | --- |
| `POST` | `/volumes/{id}/actions/create_snapshot` | Starts a crash-consistent point-in-time snapshot and returns `snapshot` plus `action`. Repeating a request with an existing name in the same project returns `uniqueness_error`. |
| `GET` | `/volume_snapshots/{id}` | Returns one snapshot. |
| `GET` | `/volume_snapshots` | Lists snapshots with normal pagination and `name`, `volume`, `status`, `label_selector`, and `sort` filters. |
| `DELETE` | `/volume_snapshots/{id}` | Starts deletion and returns `action`; deleting an absent snapshot returns `not_found`. |
| `POST` | `/volumes` | Accepts optional `snapshot`; the new volume is restored from that snapshot in its location. `size` must be at least the snapshot size. |

The create action reaches success only after snapshot data is durable. A
snapshot becomes restorable when `status` is `available`. The source volume
may stay attached while a snapshot is created; filesystem/application
quiescing remains the caller's responsibility.
