package hcloud

import (
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/hetznercloud/hcloud-go/v2/hcloud/exp/mockutil"
)

func TestZoneIDOrName(t *testing.T) {
	t.Run("with id", func(t *testing.T) {
		got, err := (&Zone{ID: 1, Name: "name"}).idOrName()
		require.NoError(t, err)
		require.Equal(t, "1", got)
	})
	t.Run("with name", func(t *testing.T) {
		got, err := (&Zone{Name: "name"}).idOrName()
		require.NoError(t, err)
		require.Equal(t, "name", got)
	})
	t.Run("missing", func(t *testing.T) {
		got, err := (&Zone{}).idOrName()
		require.EqualError(t, err, "missing one of fields [ID, Name] in [*hcloud.Zone]")
		require.Empty(t, got)
	})
}

func TestZoneGet(t *testing.T) {
	t.Run("Get using name", func(t *testing.T) {
		ctx, server, client := makeTestUtils(t)

		server.Expect([]mockutil.Request{
			{
				Method: "GET", Path: "/zones/example.com",
				Status: 200,
				JSONRaw: `{
					"zone": {
						"id": 42,
						"name": "example.com",
						"labels": { "key": "value" },
						"created": "2016-01-30T23:55:00+00:00",
						"mode": "primary",
						"authoritative_nameservers": {
							"assigned": [
								"hydrogen.ns.hetzner.com.",
								"oxygen.ns.hetzner.com.",
								"helium.ns.hetzner.de."
							],
							"delegated": [
								"hydrogen.ns.hetzner.com.",
								"oxygen.ns.hetzner.com.",
								"helium.ns.hetzner.de."
							],
							"delegation_last_check": "2016-01-30T23:55:00+00:00",
							"delegation_status": "valid"
						},
						"primary_nameservers": [
							{ "address": "198.51.100.1", "port": 53 },
							{ "address": "203.0.113.1", "port": 53 }
						],
						"registrar": "hetzner",
						"protection": {
							"delete": true
						},
						"record_count": 4,
						"status": "ok",
						"ttl": 10800
					}
				}`,
			},
		})

		result, resp, err := client.Zone.Get(ctx, "example.com")
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, int64(42), result.ID)
		require.Equal(t, "example.com", result.Name)
		require.Equal(t, ZoneModePrimary, result.Mode)
		require.Equal(t, map[string]string{"key": "value"}, result.Labels)
		require.True(t, result.Protection.Delete)
		require.Equal(t, 4, result.RecordCount)
		require.Equal(t, 10800, result.TTL)
		require.Equal(t, ZoneStatusOk, result.Status)
		require.Equal(t, "hydrogen.ns.hetzner.com.", result.AuthoritativeNameservers.Assigned[0])
		require.Equal(t, "hydrogen.ns.hetzner.com.", result.AuthoritativeNameservers.Delegated[0])
		require.Equal(t, "2016-01-30T23:55:00Z", result.AuthoritativeNameservers.DelegationLastCheck.Format(time.RFC3339))
		require.Equal(t, ZoneDelegationStatusValid, result.AuthoritativeNameservers.DelegationStatus)
		require.Equal(t, ZoneRegistrarHetzner, result.Registrar)
	})

	t.Run("Get using id", func(t *testing.T) {
		ctx, server, client := makeTestUtils(t)

		server.Expect([]mockutil.Request{
			{
				Method: "GET", Path: "/zones/42",
				Status:  200,
				JSONRaw: `{ "zone": { "id": 42, "name": "example.com" } }`,
			},
		})

		result, resp, err := client.Zone.Get(ctx, "42")
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, int64(42), result.ID)
		require.Equal(t, "example.com", result.Name)
	})

	t.Run("GetByID", func(t *testing.T) {
		ctx, server, client := makeTestUtils(t)

		server.Expect([]mockutil.Request{
			{
				Method: "GET", Path: "/zones/42",
				Status:  200,
				JSONRaw: `{ "zone": { "id": 42, "name": "example.com" } }`,
			},
		})

		result, resp, err := client.Zone.GetByID(ctx, 42)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, int64(42), result.ID)
		require.Equal(t, "example.com", result.Name)
	})

	t.Run("GetByName", func(t *testing.T) {
		ctx, server, client := makeTestUtils(t)

		server.Expect([]mockutil.Request{
			{
				Method: "GET", Path: "/zones/example.com",
				Status:  200,
				JSONRaw: `{ "zone": { "id": 42, "name": "example.com" } }`,
			},
		})

		result, resp, err := client.Zone.GetByName(ctx, "example.com")
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, int64(42), result.ID)
		require.Equal(t, "example.com", result.Name)
	})
}

func TestZoneList(t *testing.T) {
	ctx, server, client := makeTestUtils(t)

	server.Expect([]mockutil.Request{
		{
			Method: "GET", Path: "/zones?mode=primary&name=example.com&page=2&sort=name",
			Status: 200,
			JSONRaw: `{
				"zones": [
					{ "id": 42, "name": "example.com" },
					{ "id": 43, "name": "woop.com" }
				]
			}`,
		},
	})

	result, resp, err := client.Zone.List(ctx, ZoneListOpts{
		Name: "example.com",
		Mode: ZoneModePrimary,
		Sort: []string{"name"},
		ListOpts: ListOpts{
			Page: 2,
		},
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, []*Zone{
		{ID: 42, Name: "example.com"},
		{ID: 43, Name: "woop.com"},
	}, result)
}

func TestZoneAllWithOpts(t *testing.T) {
	ctx, server, client := makeTestUtils(t)

	server.Expect([]mockutil.Request{
		{
			Method: "GET", Path: "/zones?mode=primary&page=1&sort=name",
			Status: 200,
			JSONRaw: `{
				"zones": [
					{ "id": 42, "name": "example.com" },
					{ "id": 43, "name": "woop.com" }
				],
				"meta": { "pagination": { "page": 1, "next_page": 2 }}
			}`,
		},
		{
			Method: "GET", Path: "/zones?mode=primary&page=2&sort=name",
			Status: 200,
			JSONRaw: `{
				"zones": [
					{ "id": 44, "name": "example2.com" },
					{ "id": 45, "name": "woop2.com" }
				],
				"meta": { "pagination": { "page": 2 }}
			}`,
		},
	})

	result, err := client.Zone.AllWithOpts(ctx, ZoneListOpts{
		Mode: ZoneModePrimary,
		Sort: []string{"name"},
	})
	require.NoError(t, err)
	require.Equal(t, []*Zone{
		{ID: 42, Name: "example.com"},
		{ID: 43, Name: "woop.com"},
		{ID: 44, Name: "example2.com"},
		{ID: 45, Name: "woop2.com"},
	}, result)
}

func TestZoneCreate(t *testing.T) {
	t.Run("Primary", func(t *testing.T) {
		ctx, server, client := makeTestUtils(t)

		server.Expect([]mockutil.Request{
			{
				Method: "POST", Path: "/zones",
				Want: func(t *testing.T, r *http.Request) {
					body, err := io.ReadAll(r.Body)
					require.NoError(t, err)
					require.JSONEq(t, `{
						"name": "example.com",
						"mode": "primary",
						"ttl": 10800,
						"labels": { "key": "value" },
						"rrsets": [
							{ "name": "www", "type": "A", "ttl": 3600, "records": [
								{ "value": "78.34.234.13" },
								{ "value": "78.34.234.14", "comment": "Web server" }
							]}
						]
					}`, string(body))
				},
				Status: 200,
				JSONRaw: `{
					"zone": { "id": 42, "name": "example.com" },
					"action": { "id": 14 }
				}`,
			},
		})

		result, resp, err := client.Zone.Create(ctx, ZoneCreateOpts{
			Name:   "example.com",
			Mode:   ZoneModePrimary,
			TTL:    Ptr(10800),
			Labels: map[string]string{"key": "value"},
			RRSets: []ZoneCreateOptsRRSet{
				{Name: "www", Type: "A", TTL: Ptr(3600), Records: []ZoneRRSetRecord{
					{Value: "78.34.234.13"},
					{Value: "78.34.234.14", Comment: "Web server"},
				}},
			},
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NotNil(t, result.Zone)

		require.Equal(t, int64(42), result.Zone.ID)
		require.Equal(t, "example.com", result.Zone.Name)
		require.Equal(t, int64(14), result.Action.ID)
	})

	t.Run("Primary Zonefile", func(t *testing.T) {
		ctx, server, client := makeTestUtils(t)

		server.Expect([]mockutil.Request{
			{
				Method: "POST", Path: "/zones",
				Want: func(t *testing.T, r *http.Request) {
					body, err := io.ReadAll(r.Body)
					require.NoError(t, err)
					require.JSONEq(t, `{
						"name": "example.com",
						"mode": "primary",
						"ttl": 10800,
						"labels": { "key": "value" },
						"zonefile": "content\ncontent"
					}`, string(body))
				},
				Status: 200,
				JSONRaw: `{
					"zone": { "id": 42, "name": "example.com" },
					"action": { "id": 14 }
				}`,
			},
		})

		result, resp, err := client.Zone.Create(ctx, ZoneCreateOpts{
			Name:     "example.com",
			Mode:     ZoneModePrimary,
			TTL:      Ptr(10800),
			Labels:   map[string]string{"key": "value"},
			Zonefile: "content\ncontent",
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NotNil(t, result.Zone)

		require.Equal(t, int64(42), result.Zone.ID)
		require.Equal(t, "example.com", result.Zone.Name)
		require.Equal(t, int64(14), result.Action.ID)
	})

	t.Run("Secondary", func(t *testing.T) {
		ctx, server, client := makeTestUtils(t)

		server.Expect([]mockutil.Request{
			{
				Method: "POST", Path: "/zones",
				Want: func(t *testing.T, r *http.Request) {
					body, err := io.ReadAll(r.Body)
					require.NoError(t, err)
					require.JSONEq(t, `{
						"name": "example.com",
						"mode": "secondary",
						"ttl": 10800,
						"labels": { "key": "value" },
						"primary_nameservers": [
							{ "address": "78.34.234.13" },
							{ "address": "78.34.234.14", "port": 5353 },
							{ "address": "78.34.234.15", "tsig_algorithm": "hmac-sha256", "tsig_key": "c6fadd761a178a5338761650ff0e9e1012018181" }
						]
					}`, string(body))
				},
				Status: 200,
				JSONRaw: `{
					"zone": { "id": 42, "name": "example.com" },
					"action": { "id": 14 }
				}`,
			},
		})

		result, resp, err := client.Zone.Create(ctx, ZoneCreateOpts{
			Name:   "example.com",
			Mode:   ZoneModeSecondary,
			TTL:    Ptr(10800),
			Labels: map[string]string{"key": "value"},
			PrimaryNameservers: []ZoneCreateOptsPrimaryNameserver{
				{Address: "78.34.234.13"},
				{Address: "78.34.234.14", Port: 5353},
				{Address: "78.34.234.15", TSIGAlgorithm: ZoneTSIGAlgorithmHMACSHA256, TSIGKey: "c6fadd761a178a5338761650ff0e9e1012018181"},
			},
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NotNil(t, result.Zone)

		require.Equal(t, int64(42), result.Zone.ID)
		require.Equal(t, "example.com", result.Zone.Name)
		require.Equal(t, int64(14), result.Action.ID)
	})
}

func TestZoneUpdate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctx, server, client := makeTestUtils(t)

		server.Expect([]mockutil.Request{
			{
				Method: "PUT", Path: "/zones/example.com",
				Want: func(t *testing.T, r *http.Request) {
					body, err := io.ReadAll(r.Body)
					require.NoError(t, err)
					require.JSONEq(t, `{
						"labels": { "key": "value" }
					}`, string(body))
				},
				Status: 200,
				JSONRaw: `{
					"zone": { "id": 42, "name": "example.com" }
				}`,
			},
		})

		result, resp, err := client.Zone.Update(ctx,
			&Zone{Name: "example.com"},
			ZoneUpdateOpts{
				Labels: map[string]string{"key": "value"},
			},
		)
		require.NoError(t, err)
		require.NotNil(t, resp)

		require.Equal(t, int64(42), result.ID)
		require.Equal(t, "example.com", result.Name)
	})

	t.Run("success reset labels", func(t *testing.T) {
		ctx, server, client := makeTestUtils(t)

		server.Expect([]mockutil.Request{
			{
				Method: "PUT", Path: "/zones/example.com",
				Want: func(t *testing.T, r *http.Request) {
					body, err := io.ReadAll(r.Body)
					require.NoError(t, err)
					require.JSONEq(t, `{
						"labels": {}
					}`, string(body))
				},
				Status: 200,
				JSONRaw: `{
					"zone": { "id": 42, "name": "example.com" }
				}`,
			},
		})

		result, resp, err := client.Zone.Update(ctx,
			&Zone{Name: "example.com"},
			ZoneUpdateOpts{
				Labels: map[string]string{},
			},
		)
		require.NoError(t, err)
		require.NotNil(t, resp)

		require.Equal(t, int64(42), result.ID)
		require.Equal(t, "example.com", result.Name)
	})
}

func TestZoneDelete(t *testing.T) {
	ctx, server, client := makeTestUtils(t)

	server.Expect([]mockutil.Request{
		{
			Method: "DELETE", Path: "/zones/example.com",
			Status: 200,
			JSONRaw: `{
				"action": { "id": 14 }
			}`,
		},
	})

	result, resp, err := client.Zone.Delete(ctx, &Zone{Name: "example.com"})
	require.NoError(t, err)
	require.NotNil(t, resp)

	require.Equal(t, int64(14), result.Action.ID)
}

func TestZoneExportZonefile(t *testing.T) {
	ctx, server, client := makeTestUtils(t)

	server.Expect([]mockutil.Request{
		{
			Method: "GET", Path: "/zones/example.com/zonefile",
			Status: 200,
			JSONRaw: `{
				"zonefile": "content\n\texample.com"
			}`,
		},
	})

	result, resp, err := client.Zone.ExportZonefile(ctx, &Zone{Name: "example.com"})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, "content\n\texample.com", result.Zonefile)
}

func TestZoneImportZonefile(t *testing.T) {
	ctx, server, client := makeTestUtils(t)

	server.Expect([]mockutil.Request{
		{
			Method: "POST", Path: "/zones/example.com/actions/import_zonefile",
			Want: func(t *testing.T, r *http.Request) {
				body, err := io.ReadAll(r.Body)
				require.NoError(t, err)
				require.JSONEq(t, `{
					"zonefile": "content\n\texample.com"
				}`, string(body))
			},
			Status: 200,
			JSONRaw: `{
				"action": { "id": 14 }
			}`,
		},
	})

	result, resp, err := client.Zone.ImportZonefile(ctx,
		&Zone{Name: "example.com"},
		ZoneImportZonefileOpts{
			Zonefile: "content\n\texample.com",
		},
	)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, int64(14), result.ID)
}

func TestZoneChangeProtection(t *testing.T) {
	ctx, server, client := makeTestUtils(t)

	server.Expect([]mockutil.Request{
		{
			Method: "POST", Path: "/zones/example.com/actions/change_protection",
			Want: func(t *testing.T, r *http.Request) {
				body, err := io.ReadAll(r.Body)
				require.NoError(t, err)
				require.JSONEq(t, `{
					"delete": true
				}`, string(body))
			},
			Status: 200,
			JSONRaw: `{
				"action": { "id": 14 }
			}`,
		},
	})

	result, resp, err := client.Zone.ChangeProtection(ctx,
		&Zone{Name: "example.com"},
		ZoneChangeProtectionOpts{
			Delete: Ptr(true),
		},
	)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, int64(14), result.ID)
}

func TestZoneChangeTTL(t *testing.T) {
	ctx, server, client := makeTestUtils(t)

	server.Expect([]mockutil.Request{
		{
			Method: "POST", Path: "/zones/example.com/actions/change_ttl",
			Want: func(t *testing.T, r *http.Request) {
				body, err := io.ReadAll(r.Body)
				require.NoError(t, err)
				require.JSONEq(t, `{
					"ttl": 3600
				}`, string(body))
			},
			Status: 200,
			JSONRaw: `{
				"action": { "id": 14 }
			}`,
		},
	})

	result, resp, err := client.Zone.ChangeTTL(ctx,
		&Zone{Name: "example.com"},
		ZoneChangeTTLOpts{
			TTL: 3600,
		},
	)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, int64(14), result.ID)
}

func TestZoneChangePrimaryNameservers(t *testing.T) {
	ctx, server, client := makeTestUtils(t)

	server.Expect([]mockutil.Request{
		{
			Method: "POST", Path: "/zones/example.com/actions/change_primary_nameservers",
			Want: func(t *testing.T, r *http.Request) {
				body, err := io.ReadAll(r.Body)
				require.NoError(t, err)
				require.JSONEq(t, `{
					"primary_nameservers": [
						{ "address": "45.68.23.61", "port": 5353 },
						{ "address": "45.68.23.62" }
					]
				}`, string(body))
			},
			Status: 200,
			JSONRaw: `{
				"action": { "id": 14 }
			}`,
		},
	})

	result, resp, err := client.Zone.ChangePrimaryNameservers(ctx,
		&Zone{Name: "example.com"},
		ZoneChangePrimaryNameserversOpts{
			PrimaryNameservers: []ZoneChangePrimaryNameserversOptsPrimaryNameserver{
				{Address: "45.68.23.61", Port: 5353},
				{Address: "45.68.23.62"},
			},
		},
	)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, int64(14), result.ID)
}
