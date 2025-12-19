package hcloud

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/hetznercloud/hcloud-go/v2/hcloud/exp/mockutil"
)

func TestZoneRRSetNameAndType(t *testing.T) {
	t.Run("with name and type", func(t *testing.T) {
		gotName, gotType, err := (&ZoneRRSet{Name: "www", Type: "A"}).nameAndType()
		require.NoError(t, err)
		require.Equal(t, "www", gotName)
		require.Equal(t, ZoneRRSetTypeA, gotType)
	})
	t.Run("with id", func(t *testing.T) {
		gotName, gotType, err := (&ZoneRRSet{ID: "www/A"}).nameAndType()
		require.NoError(t, err)
		require.Equal(t, "www", gotName)
		require.Equal(t, ZoneRRSetTypeA, gotType)
	})
	t.Run("with invalid id", func(t *testing.T) {
		gotName, gotType, err := (&ZoneRRSet{ID: "www"}).nameAndType()
		require.EqualError(t, err, "invalid value 'www' for field [ID] in [*hcloud.ZoneRRSet]")
		require.Empty(t, gotName)
		require.Equal(t, ZoneRRSetType(""), gotType)
	})
	t.Run("missing", func(t *testing.T) {
		gotName, gotType, err := (&ZoneRRSet{Name: "www"}).nameAndType()
		require.EqualError(t, err, "missing required together fields [Name, Type] in [*hcloud.ZoneRRSet]")
		require.Empty(t, gotName)
		require.Equal(t, ZoneRRSetType(""), gotType)

		gotName, gotType, err = (&ZoneRRSet{Type: "A"}).nameAndType()
		require.EqualError(t, err, "missing required together fields [Name, Type] in [*hcloud.ZoneRRSet]")
		require.Empty(t, gotName)
		require.Equal(t, ZoneRRSetType(""), gotType)

		gotName, gotType, err = (&ZoneRRSet{}).nameAndType()
		require.EqualError(t, err, "missing one of fields [ID, Name] in [*hcloud.ZoneRRSet]")
		require.Empty(t, gotName)
		require.Equal(t, ZoneRRSetType(""), gotType)
	})
}

func TestZoneGetRRSet(t *testing.T) {
	t.Run("GetRRSetByNameAndType", func(t *testing.T) {
		ctx, server, client := makeTestUtils(t)

		server.Expect([]mockutil.Request{
			{
				Method: "GET", Path: "/zones/example.com/rrsets/www/A",
				Status: 200,
				JSONRaw: `{
					"rrset": {
						"zone": 42,
						"id": "www/A",
						"name": "www",
						"type": "A",
						"ttl": 3600,
						"labels": {
							"key": "value"
						},
						"protection": {
							"change": true
						},
						"records": [
							{ "value": "198.51.100.1", "comment": "web server" }
						]
					}
				}`,
			},
		})

		result, resp, err := client.Zone.GetRRSetByNameAndType(ctx,
			&Zone{Name: "example.com"},
			"www", "A",
		)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, int64(42), result.Zone.ID)
		require.Equal(t, "www/A", result.ID)
		require.Equal(t, "www", result.Name)
		require.Equal(t, ZoneRRSetTypeA, result.Type)
		require.Equal(t, Ptr(3600), result.TTL)
		require.Equal(t, map[string]string{"key": "value"}, result.Labels)
		require.True(t, result.Protection.Change)
		require.Equal(t, []ZoneRRSetRecord{
			{Value: "198.51.100.1", Comment: "web server"},
		}, result.Records)
	})

	t.Run("GetRRSetByID", func(t *testing.T) {
		ctx, server, client := makeTestUtils(t)

		server.Expect([]mockutil.Request{
			{
				Method: "GET", Path: "/zones/example.com/rrsets/www/A",
				Status: 200,
				JSONRaw: `{
					"rrset": { "zone": 42, "id": "www/A", "name": "www", "type": "A" }
				}`,
			},
		})

		result, resp, err := client.Zone.GetRRSetByID(ctx,
			&Zone{Name: "example.com"},
			"www/A",
		)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, int64(42), result.Zone.ID)
		require.Equal(t, "www/A", result.ID)
		require.Equal(t, "www", result.Name)
	})
}

func TestZoneListRRSets(t *testing.T) {
	ctx, server, client := makeTestUtils(t)

	server.Expect([]mockutil.Request{
		{
			Method: "GET", Path: "/zones/example.com/rrsets?name=www&page=2&sort=name&type=A",
			Status: 200,
			JSONRaw: `{
				"rrsets": [
					{ "zone": 42, "id": "www/A", "name": "www", "type": "A" },
					{ "zone": 42, "id": "blog/A", "name": "blog", "type": "A" }
				]
			}`,
		},
	})

	result, resp, err := client.Zone.ListRRSets(ctx,
		&Zone{Name: "example.com"},
		ZoneRRSetListOpts{
			Name: "www",
			Type: []ZoneRRSetType{"A"},
			Sort: []string{"name"},
			ListOpts: ListOpts{
				Page: 2,
			},
		},
	)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, []*ZoneRRSet{
		{Zone: &Zone{ID: 42}, ID: "www/A", Name: "www", Type: "A"},
		{Zone: &Zone{ID: 42}, ID: "blog/A", Name: "blog", Type: "A"},
	}, result)
}

func TestZoneAllRRSetsWithOpts(t *testing.T) {
	ctx, server, client := makeTestUtils(t)

	server.Expect([]mockutil.Request{
		{
			Method: "GET", Path: "/zones/example.com/rrsets?page=1&per_page=50&sort=name&type=A",
			Status: 200,
			JSONRaw: `{
				"rrsets": [
					{ "zone": 42, "id": "www/A", "name": "www", "type": "A" },
					{ "zone": 42, "id": "blog/A", "name": "blog", "type": "A" }
				],
				"meta": { "pagination": { "page": 1, "next_page": 2 }}
			}`,
		},
		{
			Method: "GET", Path: "/zones/example.com/rrsets?page=2&per_page=50&sort=name&type=A",
			Status: 200,
			JSONRaw: `{
				"rrsets": [
					{ "zone": 42, "id": "status/A", "name": "status", "type": "A" },
					{ "zone": 42, "id": "support/A", "name": "support", "type": "A" }
				],
				"meta": { "pagination": { "page": 2 }}
			}`,
		},
	})

	result, err := client.Zone.AllRRSetsWithOpts(ctx,
		&Zone{Name: "example.com"},
		ZoneRRSetListOpts{
			Type: []ZoneRRSetType{"A"},
			Sort: []string{"name"},
		},
	)
	require.NoError(t, err)
	require.Equal(t, []*ZoneRRSet{
		{Zone: &Zone{ID: 42}, ID: "www/A", Name: "www", Type: "A"},
		{Zone: &Zone{ID: 42}, ID: "blog/A", Name: "blog", Type: "A"},
		{Zone: &Zone{ID: 42}, ID: "status/A", Name: "status", Type: "A"},
		{Zone: &Zone{ID: 42}, ID: "support/A", Name: "support", Type: "A"},
	}, result)
}

func TestZoneAllRRSets(t *testing.T) {
	ctx, server, client := makeTestUtils(t)

	server.Expect([]mockutil.Request{
		{
			Method: "GET", Path: "/zones/example.com/rrsets?page=1&per_page=50",
			Status: 200,
			JSONRaw: `{
				"rrsets": [
					{ "zone": 42, "id": "www/A", "name": "www", "type": "A" },
					{ "zone": 42, "id": "blog/A", "name": "blog", "type": "A" }
				],
				"meta": { "pagination": { "page": 1, "next_page": 2 }}
			}`,
		},
		{
			Method: "GET", Path: "/zones/example.com/rrsets?page=2&per_page=50",
			Status: 200,
			JSONRaw: `{
				"rrsets": [
					{ "zone": 42, "id": "status/A", "name": "status", "type": "A" },
					{ "zone": 42, "id": "support/A", "name": "support", "type": "A" }
				],
				"meta": { "pagination": { "page": 2 }}
			}`,
		},
	})

	result, err := client.Zone.AllRRSets(ctx,
		&Zone{Name: "example.com"},
	)
	require.NoError(t, err)
	require.Equal(t, []*ZoneRRSet{
		{Zone: &Zone{ID: 42}, ID: "www/A", Name: "www", Type: "A"},
		{Zone: &Zone{ID: 42}, ID: "blog/A", Name: "blog", Type: "A"},
		{Zone: &Zone{ID: 42}, ID: "status/A", Name: "status", Type: "A"},
		{Zone: &Zone{ID: 42}, ID: "support/A", Name: "support", Type: "A"},
	}, result)
}

func TestZoneCreateRRSet(t *testing.T) {
	ctx, server, client := makeTestUtils(t)

	server.Expect([]mockutil.Request{
		{
			Method: "POST", Path: "/zones/example.com/rrsets",
			Want: func(t *testing.T, r *http.Request) {
				body, err := io.ReadAll(r.Body)
				require.NoError(t, err)
				require.JSONEq(t, `{
					"name": "www",
					"type": "A",
					"ttl": 3600,
					"labels": { "key": "value" },
					"records": [
						{ "value": "198.51.100.1", "comment": "web server" }
					]
				}`, string(body))
			},
			Status: 200,
			JSONRaw: `{
				"rrset": { "zone": 42, "id": "www/A", "name": "www", "type": "A" },
				"action": { "id": 14 }
			}`,
		},
	})

	result, resp, err := client.Zone.CreateRRSet(ctx,
		&Zone{Name: "example.com"},
		ZoneRRSetCreateOpts{
			Name:   "www",
			Type:   ZoneRRSetTypeA,
			TTL:    Ptr(3600),
			Labels: map[string]string{"key": "value"},
			Records: []ZoneRRSetRecord{
				{Value: "198.51.100.1", Comment: "web server"},
			},
		},
	)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, result.RRSet)

	require.Equal(t, int64(42), result.RRSet.Zone.ID)
	require.Equal(t, "www/A", result.RRSet.ID)
	require.Equal(t, "www", result.RRSet.Name)
	require.Equal(t, ZoneRRSetTypeA, result.RRSet.Type)
	require.Equal(t, int64(14), result.Action.ID)
}

func TestZoneUpdateRRSet(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctx, server, client := makeTestUtils(t)

		server.Expect([]mockutil.Request{
			{
				Method: "PUT", Path: "/zones/example.com/rrsets/www/A",
				Want: func(t *testing.T, r *http.Request) {
					body, err := io.ReadAll(r.Body)
					require.NoError(t, err)
					require.JSONEq(t, `{
						"labels": { "key": "value" }
					}`, string(body))
				},
				Status: 200,
				JSONRaw: `{
					"rrset": { "zone": 42, "id": "www/A", "name": "www", "type": "A" }
				}`,
			},
		})

		result, resp, err := client.Zone.UpdateRRSet(ctx,
			&ZoneRRSet{
				Zone: &Zone{Name: "example.com"},
				ID:   "www/A",
			},
			ZoneRRSetUpdateOpts{
				Labels: map[string]string{"key": "value"},
			},
		)
		require.NoError(t, err)
		require.NotNil(t, resp)

		require.Equal(t, int64(42), result.Zone.ID)
		require.Equal(t, "www/A", result.ID)
		require.Equal(t, "www", result.Name)
		require.Equal(t, ZoneRRSetTypeA, result.Type)
	})

	t.Run("success reset labels", func(t *testing.T) {
		ctx, server, client := makeTestUtils(t)

		server.Expect([]mockutil.Request{
			{
				Method: "PUT", Path: "/zones/example.com/rrsets/www/A",
				Want: func(t *testing.T, r *http.Request) {
					body, err := io.ReadAll(r.Body)
					require.NoError(t, err)
					require.JSONEq(t, `{
						"labels": {}
					}`, string(body))
				},
				Status: 200,
				JSONRaw: `{
					"rrset": { "zone": 42, "id": "www/A", "name": "www", "type": "A" }
				}`,
			},
		})

		result, resp, err := client.Zone.UpdateRRSet(ctx,
			&ZoneRRSet{
				Zone: &Zone{Name: "example.com"},
				ID:   "www/A",
			},
			ZoneRRSetUpdateOpts{
				Labels: map[string]string{},
			},
		)
		require.NoError(t, err)
		require.NotNil(t, resp)

		require.Equal(t, int64(42), result.Zone.ID)
		require.Equal(t, "www/A", result.ID)
		require.Equal(t, "www", result.Name)
		require.Equal(t, ZoneRRSetTypeA, result.Type)
	})
}

func TestZoneDeleteRRSet(t *testing.T) {
	ctx, server, client := makeTestUtils(t)

	server.Expect([]mockutil.Request{
		{
			Method: "DELETE", Path: "/zones/example.com/rrsets/www/A",
			Status: 200,
			JSONRaw: `{
				"action": { "id": 14 }
			}`,
		},
	})

	result, resp, err := client.Zone.DeleteRRSet(ctx,
		&ZoneRRSet{
			Zone: &Zone{Name: "example.com"},
			ID:   "www/A",
		},
	)
	require.NoError(t, err)
	require.NotNil(t, resp)

	require.Equal(t, int64(14), result.Action.ID)
}

func TestZoneChangeRRSetProtection(t *testing.T) {
	ctx, server, client := makeTestUtils(t)

	server.Expect([]mockutil.Request{
		{
			Method: "POST", Path: "/zones/example.com/rrsets/www/A/actions/change_protection",
			Want: func(t *testing.T, r *http.Request) {
				body, err := io.ReadAll(r.Body)
				require.NoError(t, err)
				require.JSONEq(t, `{
					"change": true
				}`, string(body))
			},
			Status: 200,
			JSONRaw: `{
				"action": { "id": 14 }
			}`,
		},
	})

	result, resp, err := client.Zone.ChangeRRSetProtection(ctx,
		&ZoneRRSet{
			Zone: &Zone{Name: "example.com"},
			ID:   "www/A",
		},
		ZoneRRSetChangeProtectionOpts{
			Change: Ptr(true),
		},
	)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, int64(14), result.ID)
}

func TestZoneChangeRRSetTTL(t *testing.T) {
	ctx, server, client := makeTestUtils(t)

	server.Expect([]mockutil.Request{
		{
			Method: "POST", Path: "/zones/example.com/rrsets/www/A/actions/change_ttl",
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

	result, resp, err := client.Zone.ChangeRRSetTTL(ctx,
		&ZoneRRSet{
			Zone: &Zone{Name: "example.com"},
			ID:   "www/A",
		},
		ZoneRRSetChangeTTLOpts{
			TTL: Ptr(3600),
		},
	)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, int64(14), result.ID)
}

func TestZoneSetRRSetRecords(t *testing.T) {
	ctx, server, client := makeTestUtils(t)

	server.Expect([]mockutil.Request{
		{
			Method: "POST", Path: "/zones/example.com/rrsets/www/A/actions/set_records",
			Want: func(t *testing.T, r *http.Request) {
				body, err := io.ReadAll(r.Body)
				require.NoError(t, err)
				require.JSONEq(t, `{
					"records": [
						{ "value": "34.68.10.234", "comment": "web server 1" },
						{ "value": "34.68.10.235", "comment": "web server 2" },
						{ "value": "52.12.45.3" }
					]
				}`, string(body))
			},
			Status: 200,
			JSONRaw: `{
				"action": { "id": 14 }
			}`,
		},
	})

	result, resp, err := client.Zone.SetRRSetRecords(ctx,
		&ZoneRRSet{
			Zone: &Zone{Name: "example.com"},
			ID:   "www/A",
		},
		ZoneRRSetSetRecordsOpts{
			Records: []ZoneRRSetRecord{
				{Value: "34.68.10.234", Comment: "web server 1"},
				{Value: "34.68.10.235", Comment: "web server 2"},
				{Value: "52.12.45.3"},
			},
		},
	)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, int64(14), result.ID)
}

func TestZoneAddRRSetRecords(t *testing.T) {
	ctx, server, client := makeTestUtils(t)

	server.Expect([]mockutil.Request{
		{
			Method: "POST", Path: "/zones/example.com/rrsets/www/A/actions/add_records",
			Want: func(t *testing.T, r *http.Request) {
				body, err := io.ReadAll(r.Body)
				require.NoError(t, err)
				require.JSONEq(t, `{
					"records": [
						{ "value": "34.68.10.234", "comment": "web server 1" },
						{ "value": "34.68.10.235", "comment": "web server 2" },
						{ "value": "52.12.45.3" }
					]
				}`, string(body))
			},
			Status: 200,
			JSONRaw: `{
				"action": { "id": 14 }
			}`,
		},
	})

	result, resp, err := client.Zone.AddRRSetRecords(ctx,
		&ZoneRRSet{
			Zone: &Zone{Name: "example.com"},
			ID:   "www/A",
		},
		ZoneRRSetAddRecordsOpts{
			Records: []ZoneRRSetRecord{
				{Value: "34.68.10.234", Comment: "web server 1"},
				{Value: "34.68.10.235", Comment: "web server 2"},
				{Value: "52.12.45.3"},
			},
		},
	)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, int64(14), result.ID)
}

func TestZoneUpdateRRSetRecords(t *testing.T) {
	ctx, server, client := makeTestUtils(t)

	server.Expect([]mockutil.Request{
		{
			Method: "POST", Path: "/zones/example.com/rrsets/www/A/actions/update_records",
			Want: func(t *testing.T, r *http.Request) {
				body, err := io.ReadAll(r.Body)
				require.NoError(t, err)
				require.JSONEq(t, `{
					"records": [
						{ "value": "34.68.10.234", "comment": "new comment 1" },
						{ "value": "34.68.10.235", "comment": "new comment 2" },
						{ "value": "52.12.45.3", "comment": "" }
					]
				}`, string(body))
			},
			Status: 200,
			JSONRaw: `{
				"action": { "id": 14 }
			}`,
		},
	})

	result, resp, err := client.Zone.UpdateRRSetRecords(ctx,
		&ZoneRRSet{
			Zone: &Zone{Name: "example.com"},
			ID:   "www/A",
		},
		ZoneRRSetUpdateRecordsOpts{
			Records: []ZoneRRSetRecord{
				{Value: "34.68.10.234", Comment: "new comment 1"},
				{Value: "34.68.10.235", Comment: "new comment 2"},
				{Value: "52.12.45.3"}, // Removes comment
			},
		},
	)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, int64(14), result.ID)
}

func TestZoneRemoveRRSetRecords(t *testing.T) {
	ctx, server, client := makeTestUtils(t)

	server.Expect([]mockutil.Request{
		{
			Method: "POST", Path: "/zones/example.com/rrsets/www/A/actions/remove_records",
			Want: func(t *testing.T, r *http.Request) {
				body, err := io.ReadAll(r.Body)
				require.NoError(t, err)
				require.JSONEq(t, `{
					"records": [
						{ "value": "34.68.10.234", "comment": "web server 1" },
						{ "value": "34.68.10.235", "comment": "web server 2" },
						{ "value": "52.12.45.3" }
					]
				}`, string(body))
			},
			Status: 200,
			JSONRaw: `{
				"action": { "id": 14 }
			}`,
		},
	})

	result, resp, err := client.Zone.RemoveRRSetRecords(ctx,
		&ZoneRRSet{
			Zone: &Zone{Name: "example.com"},
			ID:   "www/A",
		},
		ZoneRRSetRemoveRecordsOpts{
			Records: []ZoneRRSetRecord{
				{Value: "34.68.10.234", Comment: "web server 1"},
				{Value: "34.68.10.235", Comment: "web server 2"},
				{Value: "52.12.45.3"},
			},
		},
	)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, int64(14), result.ID)
}
