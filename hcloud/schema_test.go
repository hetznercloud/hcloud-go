package hcloud

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/hcloud-go/v2/hcloud/schema"
)

type SchemaTestCase struct {
	Data  string
	Check func(t *testing.T, obj any)
}

func TestActionSchema(t *testing.T) {
	data := []byte(`{
		"id": 1,
		"command": "create_server",
		"status": "success",
		"progress": 100,
		"started": "2016-01-30T23:55:00Z",
		"finished": "2016-01-30T23:56:13Z",
		"resources": [
			{
				"id": 42,
				"type": "server"
			}
		],
		"error": {
			"code": "action_failed",
			"message": "Action failed"
		}
	}`)

	var s schema.Action
	assert.NoError(t, json.Unmarshal(data, &s))

	assert.Equal(t, s, SchemaFromAction(ActionFromSchema(s)))

	a := ActionFromSchema(s)
	assert.Equal(t, a, ActionFromSchema(SchemaFromAction(a)))
}

func TestActionsSchema(t *testing.T) {
	data := []byte(`[
		{
			"id": 13,
			"command": "create_server"
		},
		{
			"id": 14,
			"command": "start_server"
		}
	]`)

	var s []schema.Action
	assert.NoError(t, json.Unmarshal(data, &s))

	assert.Equal(t, s, SchemaFromActions(ActionsFromSchema(s)))

	a := ActionsFromSchema(s)
	assert.Equal(t, a, ActionsFromSchema(SchemaFromActions(a)))
}

func TestFloatingIPSchema(t *testing.T) {
	t.Run("IPv6", func(t *testing.T) {
		data := []byte(`{
			"id": 4711,
			"name": "Web Frontend",
			"description": "Web Frontend",
			"created":"2017-08-16T17:29:14+00:00",
			"ip": "2001:db8::/64",
			"type": "ipv6",
			"server": null,
			"dns_ptr": [],
			"blocked": true,
			"home_location": {
				"id": 1,
				"name": "fsn1",
				"description": "Falkenstein DC Park 1",
				"country": "DE",
				"city": "Falkenstein",
				"latitude": 50.47612,
				"longitude": 12.370071,
				"network_zone": "eu-central"
			},
			"protection": {
				"delete": true
			},
			"labels": {
				"key": "value",
				"key2": "value2"
			}
		}`)

		var s schema.FloatingIP
		assert.NoError(t, json.Unmarshal(data, &s))

		assert.Equal(t, s, SchemaFromFloatingIP(FloatingIPFromSchema(s)))

		ip := FloatingIPFromSchema(s)
		assert.Equal(t, ip, FloatingIPFromSchema(SchemaFromFloatingIP(ip)))
	})

	t.Run("IPv4", func(t *testing.T) {
		data := []byte(`{
			"id": 4711,
			"description": "Web Frontend",
			"ip": "131.232.99.1",
			"type": "ipv4",
			"server": 42,
			"dns_ptr": [{
				"ip": "131.232.99.1",
				"dns_ptr": "fip01.example.com"
			}],
			"blocked": false,
			"home_location": {
				"id": 1,
				"name": "fsn1",
				"description": "Falkenstein DC Park 1",
				"country": "DE",
				"city": "Falkenstein",
				"latitude": 50.47612,
				"longitude": 12.370071
			}
		}`)

		var s schema.FloatingIP
		assert.NoError(t, json.Unmarshal(data, &s))

		assert.Equal(t, s, SchemaFromFloatingIP(FloatingIPFromSchema(s)))

		ip := FloatingIPFromSchema(s)
		assert.Equal(t, ip, FloatingIPFromSchema(SchemaFromFloatingIP(ip)))
	})
}

func TestPrimaryIPSchema(t *testing.T) {
	t.Run("IPv6", func(t *testing.T) {
		data := []byte(`{
			"assignee_id": 17,
			"assignee_type": "server",
			"auto_delete": true,
			"blocked": true,
			"created": "2017-08-16T17:29:14+00:00",
			"datacenter": {
				"description": "Falkenstein DC Park 8",
				"id": 42,
				"location": {
					"city": "Falkenstein",
					"country": "DE",
					"description": "Falkenstein DC Park 1",
					"id": 1,
					"latitude": 50.47612,
					"longitude": 12.370071,
					"name": "fsn1",
					"network_zone": "eu-central"
				},
				"name": "fsn1-dc8",
				"server_types": {
					"available": [],
					"available_for_migration": [],
					"supported": []
				}
			},
			"dns_ptr": [
				{
					"dns_ptr": "server.example.com",
					"ip": "fe80::"
				}
			],
			"id": 4711,
			"ip": "fe80::/64",
			"labels": {
				"key": "value",
				"key2": "value2"
			},
			"name": "Web Frontend",
			"protection": {
				"delete": true
			},
			"type": "ipv6"
        }`)

		var s schema.PrimaryIP
		assert.NoError(t, json.Unmarshal(data, &s))

		assert.Equal(t, s, SchemaFromPrimaryIP(PrimaryIPFromSchema(s)))

		ip := PrimaryIPFromSchema(s)
		assert.Equal(t, ip, PrimaryIPFromSchema(SchemaFromPrimaryIP(ip)))
	})

	t.Run("IPv4", func(t *testing.T) {
		data := []byte(`{
			"assignee_id": 17,
			"assignee_type": "server",
			"auto_delete": true,
			"blocked": true,
			"created": "2017-08-16T17:29:14+00:00",
			"datacenter": {
				"description": "Falkenstein DC Park 8",
				"id": 42,
				"location": {
					"city": "Falkenstein",
					"country": "DE",
					"description": "Falkenstein DC Park 1",
					"id": 1,
					"latitude": 50.47612,
					"longitude": 12.370071,
					"name": "fsn1",
					"network_zone": "eu-central"
				},
				"name": "fsn1-dc8",
				"server_types": {
					"available": [],
					"available_for_migration": [],
					"supported": []
				}
			},
			"dns_ptr": [
				{
					"dns_ptr": "server.example.com",
					"ip": "127.0.0.1"
				}
			],
			"id": 4711,
			"ip": "127.0.0.1",
			"labels": {
				"key": "value",
				"key2": "value2"
			},
			"name": "Web Frontend",
			"protection": {
				"delete": true
			},
			"type": "ipv4"
		}`)

		var s schema.PrimaryIP
		if err := json.Unmarshal(data, &s); err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, s, SchemaFromPrimaryIP(PrimaryIPFromSchema(s)))

		ip := PrimaryIPFromSchema(s)
		assert.Equal(t, ip, PrimaryIPFromSchema(SchemaFromPrimaryIP(ip)))
	})
}

func TestISOSchema(t *testing.T) {
	for testName, tc := range map[string]SchemaTestCase{
		"without deprecation": {
			Data: `{
				"id": 4711,
				"name": "FreeBSD-11.0-RELEASE-amd64-dvd1",
				"description": "FreeBSD 11.0 x64",
				"type": "public",
				"architecture": "x86"
			}`,
		},
		"with deprecation": {
			Data: `{
				"id": 4711,
				"name": "FreeBSD-11.0-RELEASE-amd64-dvd1",
				"description": "FreeBSD 11.0 x64",
				"type": "public",
				"architecture": "x86",
				"deprecation": {
					"announced": "2018-01-28T00:00:00+00:00",
					"unavailable_after": "2018-04-28T00:00:00+00:00"
				},
				"deprecated": "2000-01-01T00:00:00+00:00"
			}`,
			Check: func(t *testing.T, obj any) {
				v := obj.(*ISO)
				expDeprecated, err := time.Parse(time.RFC3339, "2018-04-28T00:00:00+00:00")
				if err != nil {
					t.Fatal(err)
				}
				assert.NotNil(t, v.Deprecation)
				assert.Equal(t, v.Deprecated, expDeprecated)
			},
		},
	} {
		t.Run(testName, func(t *testing.T) {
			data := []byte(tc.Data)
			var s schema.ISO
			assert.NoError(t, json.Unmarshal(data, &s))

			assert.Equal(t, s, SchemaFromISO(ISOFromSchema(s)))

			iso := ISOFromSchema(s)

			if tc.Check != nil {
				tc.Check(t, iso)
			}

			assert.Equal(t, iso, ISOFromSchema(SchemaFromISO(iso)))
		})
	}
}

func TestDatacenterSchema(t *testing.T) {
	data := []byte(`{
		"id": 1,
		"name": "fsn1-dc8",
		"description": "Falkenstein 1 DC 8",
		"location": {
			"id": 1,
			"name": "fsn1",
			"description": "Falkenstein DC Park 1",
			"country": "DE",
			"city": "Falkenstein",
			"latitude": 50.47612,
			"longitude": 12.370071,
			"network_zone": "eu-central"
		},
		"server_types": {
			"supported": [
				1,
				1,
				2,
				3
			],
			"available": [
				1,
				1,
				2,
				3
			]
		}
	}`)

	var s schema.Datacenter
	assert.NoError(t, json.Unmarshal(data, &s))

	assert.Equal(t, s, SchemaFromDatacenter(DatacenterFromSchema(s)))

	dc := DatacenterFromSchema(s)
	assert.Equal(t, dc, DatacenterFromSchema(SchemaFromDatacenter(dc)))
}

func TestLocationSchema(t *testing.T) {
	data := []byte(`{
		"id": 1,
		"name": "fsn1",
		"description": "Falkenstein DC Park 1",
		"country": "DE",
		"city": "Falkenstein",
		"latitude": 50.47612,
		"longitude": 12.370071,
		"network_zone": "eu-central"
	}`)

	var s schema.Location
	assert.NoError(t, json.Unmarshal(data, &s))

	assert.Equal(t, s, SchemaFromLocation(LocationFromSchema(s)))

	l := LocationFromSchema(s)
	assert.Equal(t, l, LocationFromSchema(SchemaFromLocation(l)))
}

func TestServerSchema(t *testing.T) {
	data := []byte(`{
		"id": 1,
		"name": "server.example.com",
		"status": "running",
		"created": "2017-08-16T17:29:14+00:00",
		"public_net": {
			"ipv4": null,
			"ipv6": {
				"ip": "2a01:4f8:1c11:3400::/64",
				"blocked": false,
				"dns_ptr": [
					{
						"ip": "2a01:4f8:1c11:3400::1/64",
						"dns_ptr": "server01.example.com"
					}
				]
			}
		},
		"private_net": [
			{
				"network": 4711,
				"ip": "10.0.1.1",
				"aliases": [
					"10.0.1.2"
				]
			}
		],
		"server_type": {
			"id": 2
		},
		"outgoing_traffic": 123456,
		"ingoing_traffic": 7891011,
		"included_traffic": 654321,
		"backup_window": "22-02",
		"rescue_enabled": true,
		"primary_disk_size": 20,
		"image": {
			"id": 4711,
			"type": "system",
			"status": "available",
			"name": "ubuntu16.04-standard-x64",
			"description": "Ubuntu 16.04 Standard 64 bit",
			"image_size": 2.3,
			"disk_size": 10,
			"created": "2017-08-16T17:29:14+00:00",
			"created_from": {
				"id": 1,
				"name": "Server"
			},
			"bound_to": 1,
			"os_flavor": "ubuntu",
			"os_version": "16.04",
			"rapid_deploy": false
		},
		"iso": {
			"id": 4711,
			"name": "FreeBSD-11.0-RELEASE-amd64-dvd1",
			"description": "FreeBSD 11.0 x64",
			"type": "public"
		},
		"datacenter": {
			"id": 1,
			"name": "fsn1-dc8",
			"description": "Falkenstein 1 DC 8",
			"location": {
				"id": 1,
				"name": "fsn1",
				"description": "Falkenstein DC Park 1",
				"country": "DE",
				"city": "Falkenstein",
				"latitude": 50.47612,
				"longitude": 12.370071,
				"network_zone": "eu-central"
			}
		},
		"protection": {
			"delete": true,
			"rebuild": true
		},
		"locked": true,
		"labels": {
			"key": "value",
			"key2": "value2"
		},
		"volumes": [123, 456, 789],
		"placement_group": {
			"created": "2019-01-08T12:10:00+00:00",
			"id": 897,
			"labels": {
			  "key": "value"
			},
			"name": "my Placement Group",
			"servers": [
			  4711,
			  4712
			],
			"type": "spread"
		}
	}`)

	var s schema.Server
	assert.NoError(t, json.Unmarshal(data, &s))

	assert.Equal(t, s, SchemaFromServer(ServerFromSchema(s)))

	iso := ServerFromSchema(s)
	assert.Equal(t, iso, ServerFromSchema(SchemaFromServer(iso)))
}

func TestServerSchemaNoTraffic(t *testing.T) {
	data := []byte(`{
		"public_net": {
			"ipv4": {
				"ip": "1.2.3.4",
				"blocked": false,
				"dns_ptr": "server01.example.com"
			},
			"ipv6": {
				"ip": "2a01:4f8:1c11:3400::/64",
				"blocked": false,
				"dns_ptr": [
					{
						"ip": "2a01:4f8:1c11:3400::1/64",
						"dns_ptr": "server01.example.com"
					}
				]
			}
		},
		"outgoing_traffic": null,
		"ingoing_traffic": null
	}`)

	var s schema.Server
	assert.NoError(t, json.Unmarshal(data, &s))

	assert.Equal(t, s, SchemaFromServer(ServerFromSchema(s)))

	iso := ServerFromSchema(s)
	assert.Equal(t, iso, ServerFromSchema(SchemaFromServer(iso)))
}

func TestServerPublicNetSchema(t *testing.T) {
	data := []byte(`{
		"ipv4": {
			"id": 1,
			"ip": "1.2.3.4",
			"blocked": false,
			"dns_ptr": "server.example.com"
		},
		"ipv6": {
			"id": 2,
			"ip": "2a01:4f8:1c19:1403::/64",
			"blocked": false,
			"dns_ptr": []
		},
		"floating_ips": [4],
		"firewalls": [
			{
				"id": 23,
				"status": "applied"
			}
		]
	}`)

	var s schema.ServerPublicNet
	assert.NoError(t, json.Unmarshal(data, &s))

	assert.Equal(t, s, SchemaFromServerPublicNet(ServerPublicNetFromSchema(s)))

	net := ServerPublicNetFromSchema(s)
	assert.Equal(t, net, ServerPublicNetFromSchema(SchemaFromServerPublicNet(net)))
}

func TestServerPublicNetIPv4Schema(t *testing.T) {
	data := []byte(`{
		"ip": "1.2.3.4",
		"blocked": true,
		"dns_ptr": "server.example.com"
	}`)

	var s schema.ServerPublicNetIPv4
	assert.NoError(t, json.Unmarshal(data, &s))

	assert.Equal(t, s, SchemaFromServerPublicNetIPv4(ServerPublicNetIPv4FromSchema(s)))

	net := ServerPublicNetIPv4FromSchema(s)
	assert.Equal(t, net, ServerPublicNetIPv4FromSchema(SchemaFromServerPublicNetIPv4(net)))
}

func TestServerPublicNetIPv6Schema(t *testing.T) {
	data := []byte(`{
		"ip": "2a01:4f8:1c11:3400::/64",
		"blocked": true,
		"dns_ptr": [
			{
				"ip": "2a01:4f8:1c11:3400::1/64",
				"blocked": "server01.example.com"
			}
		]
	}`)

	var s schema.ServerPublicNetIPv6
	assert.NoError(t, json.Unmarshal(data, &s))

	assert.Equal(t, s, SchemaFromServerPublicNetIPv6(ServerPublicNetIPv6FromSchema(s)))

	net := ServerPublicNetIPv6FromSchema(s)
	assert.Equal(t, net, ServerPublicNetIPv6FromSchema(SchemaFromServerPublicNetIPv6(net)))
}

func TestServerPrivateNetSchema(t *testing.T) {
	data := []byte(`{
		"network": 4711,
		"ip": "10.0.1.1",
		"alias_ips": [
			"10.0.1.2"
		],
		"mac_address": "86:00:ff:2a:7d:e1"
	}`)

	var s schema.ServerPrivateNet
	assert.NoError(t, json.Unmarshal(data, &s))

	assert.Equal(t, s, SchemaFromServerPrivateNet(ServerPrivateNetFromSchema(s)))

	net := ServerPrivateNetFromSchema(s)
	assert.Equal(t, net, ServerPrivateNetFromSchema(SchemaFromServerPrivateNet(net)))
}

func TestServerTypeSchema(t *testing.T) {
	data := []byte(`{
		"id": 1,
		"name": "cx10",
		"description": "description",
		"cores": 4,
		"memory": 1.0,
		"disk": 20,
		"storage_type": "local",
		"cpu_type": "shared",
		"architecture": "x86",
		"deprecation": null,
		"prices": [
			{
				"location": "fsn1",
				"price_hourly": {
					"net": "1",
					"gross": "1.19"
				},
				"price_monthly": {
					"net": "1",
					"gross": "1.19"
				}
			}
		]
	}`)

	var s schema.ServerType
	assert.NoError(t, json.Unmarshal(data, &s))

	assert.Equal(t, s, SchemaFromServerType(ServerTypeFromSchema(s)))

	net := ServerTypeFromSchema(s)
	assert.Equal(t, net, ServerTypeFromSchema(SchemaFromServerType(net)))
}

func TestSSHKeySchema(t *testing.T) {
	data := []byte(`{
		"id": 2323,
		"name": "My key",
		"fingerprint": "b7:2f:30:a0:2f:6c:58:6c:21:04:58:61:ba:06:3b:2c",
		"public_key": "ssh-rsa AAAjjk76kgf...Xt",
		"labels": {
			"key": "value",
			"key2": "value2"
		},
		"created":"2017-08-16T17:29:14+00:00"
	}`)

	var s schema.SSHKey
	assert.NoError(t, json.Unmarshal(data, &s))

	assert.Equal(t, s, SchemaFromSSHKey(SSHKeyFromSchema(s)))

	key := SSHKeyFromSchema(s)
	assert.Equal(t, key, SSHKeyFromSchema(SchemaFromSSHKey(key)))
}

func TestVolumeSchema(t *testing.T) {
	data := []byte(`{
		"id": 4711,
		"created": "2016-01-30T23:50:11+00:00",
		"name": "db-storage",
		"status": "creating",
		"server": 2,
		"location": {
			"id": 1,
			"name": "fsn1",
			"description": "Falkenstein DC Park 1",
			"country": "DE",
			"city": "Falkenstein",
			"latitude": 50.47612,
			"longitude": 12.370071
		},
		"size": 42,
		"linux_device":"/dev/disk/by-id/scsi-0HC_volume_1",
		"protection": {
			"delete": true
		},
		"labels": {
			"key": "value",
			"key2": "value2"
		}
	}`)

	var s schema.Volume
	assert.NoError(t, json.Unmarshal(data, &s))

	assert.Equal(t, s, SchemaFromVolume(VolumeFromSchema(s)))

	v := VolumeFromSchema(s)
	assert.Equal(t, v, VolumeFromSchema(SchemaFromVolume(v)))
}

func TestErrorSchema(t *testing.T) {
	testCases := map[string]string{
		"service_error": `{
			"code": "service_error",
			"message": "An error occurred"
		}`,
		"invalid_input": `{
			"code": "invalid_input",
			"message": "invalid input",
			"details": {"fields":[{"name":"broken_field","messages":["is required"]}]}
		}`,
	}

	for name, data := range testCases {
		t.Run(name, func(t *testing.T) {
			var s schema.Error
			assert.NoError(t, json.Unmarshal([]byte(data), &s))

			assert.Equal(t, s, SchemaFromError(ErrorFromSchema(s)))

			e := ErrorFromSchema(s)
			assert.Equal(t, e, ErrorFromSchema(SchemaFromError(e)))
		})
	}
}

func TestPaginationSchema(t *testing.T) {
	data := []byte(`{
		"page": 2,
		"per_page": 25,
		"previous_page": 1,
		"next_page": 3,
		"last_page": 13,
		"total_entries": 322
	}`)

	var s schema.MetaPagination
	assert.NoError(t, json.Unmarshal(data, &s))

	assert.Equal(t, s, SchemaFromPagination(PaginationFromSchema(s)))

	p := PaginationFromSchema(s)
	assert.Equal(t, p, PaginationFromSchema(SchemaFromPagination(p)))
}

func TestImageSchema(t *testing.T) {
	data := []byte(`{
		"id": 4711,
		"type": "system",
		"status": "available",
		"name": "ubuntu16.04-standard-x64",
		"description": "Ubuntu 16.04 Standard 64 bit",
		"image_size": 2.3,
		"disk_size": 10,
		"created": "2016-01-30T23:55:01Z",
		"created_from": {
			"id": 1,
			"name": "my-server1"
		},
		"bound_to": 1,
		"os_flavor": "ubuntu",
		"os_version": "16.04",
		"architecture": "arm",
		"rapid_deploy": false,
		"protection": {
			"delete": true
		},
		"deprecated": "2018-02-28T00:00:00+00:00",
		"deleted": "2016-01-30T23:55:01+00:00",
		"labels": {
			"key": "value",
			"key2": "value2"
		}
	}`)

	var s schema.Image
	assert.NoError(t, json.Unmarshal(data, &s))

	assert.Equal(t, s, SchemaFromImage(ImageFromSchema(s)))

	img := ImageFromSchema(s)
	assert.Equal(t, img, ImageFromSchema(SchemaFromImage(img)))
}

func TestNetworkSchema(t *testing.T) {
	data := []byte(`{
		"id": 4711,
		"name": "mynet",
		"created": "2017-08-16T17:29:14+00:00",
		"ip_range": "10.0.0.0/16",
		"subnets": [
			{
				"type": "server",
				"ip_range": "10.0.1.0/24",
				"network_zone": "eu-central",
				"gateway": "10.0.0.1"
			}
		],
		"routes": [
			{
				"destination": "10.100.1.0/24",
				"gateway": "10.0.1.1"
			}
		],
		"servers": [
			4711
		],
		"protection": {
			"delete": false
		},
		"labels": {}
	}`)

	var s schema.Network
	assert.NoError(t, json.Unmarshal(data, &s))

	assert.Equal(t, s, SchemaFromNetwork(NetworkFromSchema(s)))

	n := NetworkFromSchema(s)
	assert.Equal(t, n, NetworkFromSchema(SchemaFromNetwork(n)))
}

func TestNetworkSubnetSchema(t *testing.T) {
	testCases := map[string]string{
		"server_subnet": `{
			"type": "server",
			"ip_range": "10.0.1.0/24",
			"network_zone": "eu-central",
			"gateway": "10.0.0.1"
		}`,
		"vswitch_subnet": `{
			"type": "vswitch",
			"ip_range": "10.0.1.0/24",
			"network_zone": "eu-central",
			"gateway": "10.0.0.1",
			"vswitch_id": 123
		}`,
	}

	for name, data := range testCases {
		t.Run(name, func(t *testing.T) {
			var s schema.NetworkSubnet
			assert.NoError(t, json.Unmarshal([]byte(data), &s))

			assert.Equal(t, s, SchemaFromNetworkSubnet(NetworkSubnetFromSchema(s)))

			n := NetworkSubnetFromSchema(s)
			assert.Equal(t, n, NetworkSubnetFromSchema(SchemaFromNetworkSubnet(n)))
		})
	}
}

func TestNetworkRouteSchema(t *testing.T) {
	data := []byte(`{
		"destination": "10.100.1.0/24",
		"gateway": "10.0.1.1"
	}`)

	var s schema.NetworkRoute
	assert.NoError(t, json.Unmarshal(data, &s))

	assert.Equal(t, s, SchemaFromNetworkRoute(NetworkRouteFromSchema(s)))

	r := NetworkRouteFromSchema(s)
	assert.Equal(t, r, NetworkRouteFromSchema(SchemaFromNetworkRoute(r)))
}

func TestLoadBalancerTypeSchema(t *testing.T) {
	data := []byte(`{
		"id": 1,
		"name": "lx11",
		"description": "LX11",
		"max_connections": 20000,
		"max_services": 3,
		"max_targets": 25,
		"max_assigned_certificates": 10,
		"deprecated": "2016-01-30T23:50:00+00:00",
		"prices": [
			{
				"location": "fsn1",
				"price_hourly": {
					"net": "1",
					"gross": "1.19"
				},
				"price_monthly": {
					"net": "1",
					"gross": "1.19"
				}
			}
		]
	}`)

	var s schema.LoadBalancerType
	assert.NoError(t, json.Unmarshal(data, &s))

	assert.Equal(t, s, SchemaFromLoadBalancerType(LoadBalancerTypeFromSchema(s)))

	lbt := LoadBalancerTypeFromSchema(s)
	assert.Equal(t, lbt, LoadBalancerTypeFromSchema(SchemaFromLoadBalancerType(lbt)))
}

func TestLoadBalancerSchema(t *testing.T) {
	data := []byte(`{
		"id": 4711,
		"name": "Web Frontend",
		"public_net": {
			"ipv4": {
				"ip": "131.232.99.1",
				"dns_ptr": "example.org"
			},
			"ipv6": {
				"ip": "2001:db8::1",
				"dns_ptr": "example.com"
			}
		},
		"private_net": [
			{
				"network": 4711,
				"ip": "10.0.255.1"
			}
		],
		"location": {
			"id": 1,
			"name": "fsn1",
			"description": "Falkenstein DC Park 1",
			"country": "DE",
			"city": "Falkenstein",
			"latitude": 50.47612,
			"longitude": 12.370071,
			"network_zone": "eu-central"
		},
		"load_balancer_type": {
			"id": 1,
			"name": "lx11",
			"description": "LX11",
			"max_connections": 20000,
			"services": 3,
			"prices": [
				{
					"location": "fsn-1",
					"price_hourly": {
						"net": "1",
						"gross": "1.19"
					},
					"price_monthly": {
						"net": "1",
						"gross": "1.19"
					}
				}
			]
		},
		"outgoing_traffic": 123456,
		"ingoing_traffic": 7891011,
		"included_traffic": 654321,
		"protection": {
			"delete": false
		},
		"labels": {},
		"created": "2016-01-30T23:50:00+00:00",
		"services": [
			{
				"protocol": "http",
				"listen_port": 443,
				"destination_port": 80,
				"proxyprotocol": false,
				"sticky_sessions": false,
				"http": {
					"cookie_name": "HCLBSTICKY",
					"cookie_lifetime": 300,
					"certificates": [
						897
					]
				},
				"health_check": {
					"protocol": "http",
					"port": 4711,
					"interval": 15,
					"timeout": 10,
					"retries": 3,
					"http": {
						"domain": "example.com",
						"path": "/"
					}
				}
			}
		],
		"targets": [
			{
				"type": "server",
				"server": {
					"id": 80
				},
				"label_selector": null,
				"health_status": [
					{
						"listen_port": 443,
						"status": "healthy"
					}
				],
				"use_private_ip": false
			},
			{
				"type": "label_selector",
				"label_selector": {
					"selector": "lbt"
				},
				"targets": [
					{
						"type": "server",
						"server": {
							"id": 80
						},
						"health_status": [
							{
								"listen_port": 443,
								"status": "healthy"
							}
						],
						"use_private_ip": false
					}
				]
			}
		],
		"algorithm": {
			"type": "round_robin"
		}
	}`)

	var s schema.LoadBalancer
	assert.NoError(t, json.Unmarshal(data, &s))

	assert.Equal(t, s, SchemaFromLoadBalancer(LoadBalancerFromSchema(s)))

	lb := LoadBalancerFromSchema(s)
	assert.Equal(t, lb, LoadBalancerFromSchema(SchemaFromLoadBalancer(lb)))
}

func TestLoadBalancerServiceSchema(t *testing.T) {
	data := []byte(`{
		"protocol": "http",
		"listen_port": 443,
		"destination_port": 80,
		"proxyprotocol": false,
		"http": {
			"cookie_name": "HCLBSTICKY",
			"cookie_lifetime": 300,
			"certificates": [
				897
			],
			"redirect_http": true,
			"sticky_sessions": true
		},
		"health_check": {
			"protocol": "http",
			"port": 4711,
			"interval": 15,
			"timeout": 10,
			"retries": 3,
			"http": {
				"domain": "example.com",
				"path": "/",
				"response": "",
				"status_codes":["200","201"],
				"tls": false
			}
		}
	}`)

	var s schema.LoadBalancerService
	assert.NoError(t, json.Unmarshal(data, &s))

	assert.Equal(t, s, SchemaFromLoadBalancerService(LoadBalancerServiceFromSchema(s)))

	lbs := LoadBalancerServiceFromSchema(s)
	assert.Equal(t, lbs, LoadBalancerServiceFromSchema(SchemaFromLoadBalancerService(lbs)))
}

func TestLoadBalancerServiceHealthCheckSchema(t *testing.T) {
	data := []byte(`{
		"protocol": "http",
		"port": 4711,
		"interval": 15,
		"timeout": 10,
		"retries": 3,
		"http": {
			"domain": "example.com",
			"path": "/",
			"response": "",
			"status_codes":["200","201"],
			"tls": false
		}
	}`)

	var s *schema.LoadBalancerServiceHealthCheck
	assert.NoError(t, json.Unmarshal(data, &s))

	assert.Equal(t, s, SchemaFromLoadBalancerServiceHealthCheck(LoadBalancerServiceHealthCheckFromSchema(s)))

	hc := LoadBalancerServiceHealthCheckFromSchema(s)
	assert.Equal(t, hc, LoadBalancerServiceHealthCheckFromSchema(SchemaFromLoadBalancerServiceHealthCheck(hc)))
}

func TestLoadBalancerTargetSchema(t *testing.T) {
	testCases := map[string]string{
		"server target": `{
			"type": "server",
			"server": {
				"id": 80
			},
			"label_selector": null,
			"health_status": [
				{
					"listen_port": 443,
					"status": "healthy"
				}
			],
			"use_private_ip": false
		}`,
		"label_selector target": `{
			"type": "label_selector",
			"label_selector": {
				"selector": "lbt"
			},
			"targets": [
				{
					"type": "server",
					"server": {
						"id": 80
					},
					"health_status": [
						{
							"listen_port": 443,
							"status": "healthy"
						}
					]
				}
			]
		}`,
		"ip target": `{
			"type": "ip",
			"ip": {
				"ip": "1.2.3.4"
			}
		}`,
	}

	for name, data := range testCases {
		t.Run(name, func(t *testing.T) {
			var s schema.LoadBalancerTarget
			assert.NoError(t, json.Unmarshal([]byte(data), &s))

			assert.Equal(t, s, SchemaFromLoadBalancerTarget(LoadBalancerTargetFromSchema(s)))

			lbt := LoadBalancerTargetFromSchema(s)
			assert.Equal(t, lbt, LoadBalancerTargetFromSchema(SchemaFromLoadBalancerTarget(lbt)))
		})
	}
}

func TestLoadBalancerTargetHealthStatusSchema(t *testing.T) {
	data := []byte(`{
		"listen_port": 443,
		"status": "healthy"
	}`)

	var s schema.LoadBalancerTargetHealthStatus
	assert.NoError(t, json.Unmarshal(data, &s))

	assert.Equal(t, s, SchemaFromLoadBalancerTargetHealthStatus(LoadBalancerTargetHealthStatusFromSchema(s)))

	hs := LoadBalancerTargetHealthStatusFromSchema(s)
	assert.Equal(t, hs, LoadBalancerTargetHealthStatusFromSchema(SchemaFromLoadBalancerTargetHealthStatus(hs)))
}

func TestCertificateSchema(t *testing.T) {
	testCases := []struct {
		name     string
		data     string
		expected *Certificate
	}{
		{
			name: "uploaded certificate",
			data: `{
				"id": 897,
				"name": "my website cert",
				"labels": {},
				"type": "uploaded",
				"certificate": "-----BEGIN CERTIFICATE-----\n...",
				"created": "2016-01-30T23:50:00+00:00",
				"not_valid_before": "2016-01-30T23:51:00+00:00",
				"not_valid_after": "2016-01-30T23:55:00+00:00",
				"domain_names": [
					"example.com",
					"webmail.example.com",
					"www.example.com"
				],
				"fingerprint": "03:c7:55:9b:2a:d1:04:17:09:f6:d0:7f:18:34:63:d4:3e:5f",
				"used_by": [
					{"id": 42, "type": "loadbalancer"}
				]
			}`,
			expected: &Certificate{
				ID:             897,
				Name:           "my website cert",
				Labels:         map[string]string{},
				Type:           "uploaded",
				Certificate:    "-----BEGIN CERTIFICATE-----\n...",
				Created:        mustParseTime(t, "2016-01-30T23:50:00+00:00"),
				NotValidBefore: mustParseTime(t, "2016-01-30T23:51:00+00:00"),
				NotValidAfter:  mustParseTime(t, "2016-01-30T23:55:00+00:00"),
				DomainNames:    []string{"example.com", "webmail.example.com", "www.example.com"},
				Fingerprint:    "03:c7:55:9b:2a:d1:04:17:09:f6:d0:7f:18:34:63:d4:3e:5f",
				UsedBy: []CertificateUsedByRef{
					{ID: 42, Type: "loadbalancer"},
				},
			},
		},
		{
			name: "managed certificate",
			data: `{
				"id": 898,
				"name": "managed certificate",
				"labels": {},
				"type": "managed",
				"certificate": "-----BEGIN CERTIFICATE-----\n...",
				"created": "2016-01-30T23:50:00+00:00",
				"not_valid_before": "2016-01-30T23:51:00+00:00",
				"not_valid_after": "2016-01-30T23:55:00+00:00",
				"domain_names": [
					"example.com",
					"webmail.example.com",
					"www.example.com"
				],
				"fingerprint": "03:c7:55:9b:2a:d1:04:17:09:f6:d0:7f:18:34:63:d4:3e:5f",
				"status": {
					"issuance": "completed",
					"renewal": "failed",
					"error": {
						"code": "dns_zone_not_found",
						"message": "DNS zone not found"
					}
				},
				"used_by": [
					{"id": 42, "type": "loadbalancer"}
				]
			}`,
			expected: &Certificate{
				ID:             898,
				Name:           "managed certificate",
				Labels:         map[string]string{},
				Type:           "managed",
				Certificate:    "-----BEGIN CERTIFICATE-----\n...",
				Created:        mustParseTime(t, "2016-01-30T23:50:00+00:00"),
				NotValidBefore: mustParseTime(t, "2016-01-30T23:51:00+00:00"),
				NotValidAfter:  mustParseTime(t, "2016-01-30T23:55:00+00:00"),
				DomainNames:    []string{"example.com", "webmail.example.com", "www.example.com"},
				Fingerprint:    "03:c7:55:9b:2a:d1:04:17:09:f6:d0:7f:18:34:63:d4:3e:5f",
				Status: &CertificateStatus{
					Issuance: CertificateStatusTypeCompleted,
					Renewal:  CertificateStatusTypeFailed,
					Error: &Error{
						Code:    "dns_zone_not_found",
						Message: "DNS zone not found",
					},
				},
				UsedBy: []CertificateUsedByRef{
					{ID: 42, Type: "loadbalancer"},
				},
			},
		},
	}
	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			var s schema.Certificate
			assert.NoError(t, json.Unmarshal([]byte(testCase.data), &s))

			assert.Equal(t, s, SchemaFromCertificate(CertificateFromSchema(s)))

			ct := CertificateFromSchema(s)
			assert.Equal(t, ct, CertificateFromSchema(SchemaFromCertificate(ct)))

			assert.Equal(t, testCase.expected, CertificateFromSchema(s))
			assert.Equal(t, s, SchemaFromCertificate(testCase.expected))
		})
	}
}

func TestPricingSchema(t *testing.T) {
	data := []byte(`{
		"currency": "EUR",
		"vat_rate": "19.00",
		"image": {
			"price_per_gb_month": {
				"net": "1",
				"gross": "1.19"
			}
		},
		"floating_ip": {
			"price_monthly": {
				"net": "1",
				"gross": "1.19"
			}
		},
		"floating_ips": [
			  {
				"prices": [
				  {
					"location": "fsn1",
					"price_monthly": {
					  "gross": "1.19",
					  "net": "1"
					}
				  }
				],
				"type": "ipv4"
			  }
			],
		"primary_ips": [
			{
				"prices": [
				{
					"location": "fsn1",
					"price_hourly": {
					"gross": "1.1900000000000000",
					"net": "1.0000000000"
					},
					"price_monthly": {
					"gross": "1.1900000000000000",
					"net": "1.0000000000"
					}
				}
				],
				"type": "ipv4"
			}
		],
		"traffic": {
			"price_per_tb": {
				"net": "1",
				"gross": "1.19"
			}
		},
		"server_backup": {
			"percentage": "20"
		},
		"server_types": [
			{
				"id": 4,
				"name": "CX11",
				"prices": [
					{
						"location": "fsn1",
						"price_hourly": {
							"net": "1",
							"gross": "1.19"
						},
						"price_monthly": {
							"net": "1",
							"gross": "1.19"
						}
					}
				]
			}
		],
		"load_balancer_types": [
			{
				"id": 4,
				"name": "LX11",
				"prices": [
					{
						"location": "fsn1",
						"price_hourly": {
							"net": "1",
							"gross": "1.19"
						},
						"price_monthly": {
							"net": "1",
							"gross": "1.19"
						}
					}
				]
			}
		],
		"volume": {
			"price_per_gb_month": {
				"net": "1",
				"gross": "1.19"
			}
		}
	}`)

	var s schema.Pricing
	assert.NoError(t, json.Unmarshal(data, &s))

	assert.Equal(t, s, SchemaFromPricing(PricingFromSchema(s)))

	p := PricingFromSchema(s)
	assert.Equal(t, p, PricingFromSchema(SchemaFromPricing(p)))
}

func TestFirewallSchema(t *testing.T) {
	data := []byte(`{
		"id": 897,
		"name": "my firewall",
		"labels": {
			"key": "value",
			"key2": "value2"
		},
		"created": "2016-01-30T23:50:00+00:00",
		"rules": [
			{
			  "direction": "in",
			  "source_ips": [
				"28.239.13.1/32",
				"28.239.14.0/24",
				"ff21:1eac:9a3b:ee58:5ca:990c:8bc9:c03b/128"
			  ],
			  "destination_ips": [
				"28.239.13.1/32",
				"28.239.14.0/24",
				"ff21:1eac:9a3b:ee58:5ca:990c:8bc9:c03b/128"
			  ],
			  "protocol": "tcp",
			  "port": "80",
			  "description": "allow http in"
			}
		],
		"applied_to": [
			{
			 	"server": {
					"id": 42
				},
				"type": "server"
			},
			{
			 	"label_selector": {
					"selector": "a=b"
				},
				"type": "label_selector"
			}
		  ]
	}`)

	var s schema.Firewall
	assert.NoError(t, json.Unmarshal(data, &s))

	assert.Equal(t, s, SchemaFromFirewall(FirewallFromSchema(s)))

	fw := FirewallFromSchema(s)
	assert.Equal(t, fw, FirewallFromSchema(SchemaFromFirewall(fw)))
}

func TestPlacementGroupSchema(t *testing.T) {
	data := []byte(`{
		"created": "2019-01-08T12:10:00+00:00",
		"id": 897,
		"labels": {
		  "key": "value"
		},
		"name": "my Placement Group",
		"servers": [
		  4711,
		  4712
		],
		"type": "spread"
	}`)

	var s schema.PlacementGroup
	assert.NoError(t, json.Unmarshal(data, &s))

	assert.Equal(t, s, SchemaFromPlacementGroup(PlacementGroupFromSchema(s)))

	pg := PlacementGroupFromSchema(s)
	assert.Equal(t, pg, PlacementGroupFromSchema(SchemaFromPlacementGroup(pg)))
}

func TestDeprecationSchema(t *testing.T) {
	testCases := map[string]string{
		"deprecated_resource": `{
			"deprecation": {
				"announced": "2023-06-01T00:00:00+00:00",
				"unavailable_after": "2023-09-01T00:00:00+00:00"
			}
		}`,
		"non_deprecated_resource": `{
			"deprecation": null
		}`,
	}

	for name, data := range testCases {
		t.Run(name, func(t *testing.T) {
			var s *schema.DeprecationInfo
			assert.NoError(t, json.Unmarshal([]byte(data), &s))

			assert.Equal(t, s, SchemaFromDeprecation(DeprecationFromSchema(s)))

			d := DeprecationFromSchema(s)
			assert.Equal(t, d, DeprecationFromSchema(SchemaFromDeprecation(d)))
		})
	}
}

func TestLoadBalancerCreateOptsToSchema(t *testing.T) {
	testCases := map[string]struct {
		Opts    LoadBalancerCreateOpts
		Request schema.LoadBalancerCreateRequest
	}{
		"minimal": {
			Opts: LoadBalancerCreateOpts{
				Name:             "test",
				LoadBalancerType: &LoadBalancerType{Name: "lb11"},
				Algorithm:        &LoadBalancerAlgorithm{Type: LoadBalancerAlgorithmTypeRoundRobin},
				NetworkZone:      NetworkZoneEUCentral,
			},
			Request: schema.LoadBalancerCreateRequest{
				Name:             "test",
				LoadBalancerType: "lb11",
				Algorithm: &schema.LoadBalancerCreateRequestAlgorithm{
					Type: string(LoadBalancerAlgorithmTypeRoundRobin),
				},
				NetworkZone: Ptr(string(NetworkZoneEUCentral)),
			},
		},
		"all set": {
			Opts: LoadBalancerCreateOpts{
				Name:             "test",
				LoadBalancerType: &LoadBalancerType{Name: "lb11"},
				Algorithm:        &LoadBalancerAlgorithm{Type: LoadBalancerAlgorithmTypeRoundRobin},
				NetworkZone:      NetworkZoneEUCentral,
				Labels:           map[string]string{"foo": "bar"},
				PublicInterface:  Ptr(true),
				Network:          &Network{ID: 3},
				Services: []LoadBalancerCreateOptsService{
					{
						Protocol:        LoadBalancerServiceProtocolHTTP,
						DestinationPort: Ptr(80),
						Proxyprotocol:   Ptr(true),
						HTTP: &LoadBalancerCreateOptsServiceHTTP{
							CookieName:     Ptr("keks"),
							CookieLifetime: Ptr(5 * time.Minute),
							RedirectHTTP:   Ptr(true),
							StickySessions: Ptr(true),
							Certificates:   []*Certificate{{ID: 1}, {ID: 2}},
						},
						HealthCheck: &LoadBalancerCreateOptsServiceHealthCheck{
							Protocol: LoadBalancerServiceProtocolHTTP,
							Port:     Ptr(80),
							Interval: Ptr(5 * time.Second),
							Timeout:  Ptr(1 * time.Second),
							Retries:  Ptr(3),
							HTTP: &LoadBalancerCreateOptsServiceHealthCheckHTTP{
								Domain:      Ptr("example.com"),
								Path:        Ptr("/health"),
								Response:    Ptr("ok"),
								StatusCodes: []string{"2??", "3??"},
								TLS:         Ptr(true),
							},
						},
					},
					{
						Protocol:        LoadBalancerServiceProtocolHTTP,
						DestinationPort: Ptr(443),
						Proxyprotocol:   Ptr(true),
						HTTP: &LoadBalancerCreateOptsServiceHTTP{
							CookieName:     Ptr("keks"),
							CookieLifetime: Ptr(5 * time.Minute),
							RedirectHTTP:   Ptr(true),
							StickySessions: Ptr(true),
							Certificates:   []*Certificate{{ID: 1}, {ID: 2}},
						},
						HealthCheck: &LoadBalancerCreateOptsServiceHealthCheck{
							Protocol: LoadBalancerServiceProtocolHTTP,
							Port:     Ptr(443),
							Interval: Ptr(5 * time.Second),
							Timeout:  Ptr(1 * time.Second),
							Retries:  Ptr(3),
							HTTP: &LoadBalancerCreateOptsServiceHealthCheckHTTP{
								Domain:      Ptr("example.com"),
								Path:        Ptr("/health"),
								Response:    Ptr("ok"),
								StatusCodes: []string{"4??", "5??"},
								TLS:         Ptr(true),
							},
						},
					},
				},
				Targets: []LoadBalancerCreateOptsTarget{
					{
						Type: LoadBalancerTargetTypeServer,
						Server: LoadBalancerCreateOptsTargetServer{
							Server: &Server{ID: 5},
						},
					},
					{
						Type: LoadBalancerTargetTypeIP,
						IP:   LoadBalancerCreateOptsTargetIP{IP: "1.2.3.4"},
					},
				},
			},
			Request: schema.LoadBalancerCreateRequest{
				Name:             "test",
				LoadBalancerType: "lb11",
				Algorithm: &schema.LoadBalancerCreateRequestAlgorithm{
					Type: string(LoadBalancerAlgorithmTypeRoundRobin),
				},
				NetworkZone: Ptr(string(NetworkZoneEUCentral)),
				Labels: func() *map[string]string {
					labels := map[string]string{"foo": "bar"}
					return &labels
				}(),
				PublicInterface: Ptr(true),
				Network:         Ptr(int64(3)),
				Services: []schema.LoadBalancerCreateRequestService{
					{
						Protocol:        string(LoadBalancerServiceProtocolHTTP),
						DestinationPort: Ptr(80),
						Proxyprotocol:   Ptr(true),
						HTTP: &schema.LoadBalancerCreateRequestServiceHTTP{
							CookieName:     Ptr("keks"),
							CookieLifetime: Ptr(5 * 60),
							RedirectHTTP:   Ptr(true),
							StickySessions: Ptr(true),
							Certificates:   Ptr([]int64{1, 2}),
						},
						HealthCheck: &schema.LoadBalancerCreateRequestServiceHealthCheck{
							Protocol: string(LoadBalancerServiceProtocolHTTP),
							Port:     Ptr(80),
							Interval: Ptr(5),
							Timeout:  Ptr(1),
							Retries:  Ptr(3),
							HTTP: &schema.LoadBalancerCreateRequestServiceHealthCheckHTTP{
								Domain:      Ptr("example.com"),
								Path:        Ptr("/health"),
								Response:    Ptr("ok"),
								StatusCodes: Ptr([]string{"2??", "3??"}),
								TLS:         Ptr(true),
							},
						},
					},
					{
						Protocol:        string(LoadBalancerServiceProtocolHTTP),
						DestinationPort: Ptr(443),
						Proxyprotocol:   Ptr(true),
						HTTP: &schema.LoadBalancerCreateRequestServiceHTTP{
							CookieName:     Ptr("keks"),
							CookieLifetime: Ptr(5 * 60),
							RedirectHTTP:   Ptr(true),
							StickySessions: Ptr(true),
							Certificates:   Ptr([]int64{1, 2}),
						},
						HealthCheck: &schema.LoadBalancerCreateRequestServiceHealthCheck{
							Protocol: string(LoadBalancerServiceProtocolHTTP),
							Port:     Ptr(443),
							Interval: Ptr(5),
							Timeout:  Ptr(1),
							Retries:  Ptr(3),
							HTTP: &schema.LoadBalancerCreateRequestServiceHealthCheckHTTP{
								Domain:      Ptr("example.com"),
								Path:        Ptr("/health"),
								Response:    Ptr("ok"),
								StatusCodes: Ptr([]string{"4??", "5??"}),
								TLS:         Ptr(true),
							},
						},
					},
				},
				Targets: []schema.LoadBalancerCreateRequestTarget{
					{
						Type: "server",
						Server: &schema.LoadBalancerCreateRequestTargetServer{
							ID: 5,
						},
					},
					{
						Type: "ip",
						IP: &schema.LoadBalancerCreateRequestTargetIP{
							IP: "1.2.3.4",
						},
					},
				},
			},
		},
	}
	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			req := loadBalancerCreateOptsToSchema(testCase.Opts)
			if !cmp.Equal(testCase.Request, req) {
				t.Log(cmp.Diff(testCase.Request, req))
				t.Fail()
			}
		})
	}
}

func TestLoadBalancerAddServiceOptsToSchema(t *testing.T) {
	testCases := map[string]struct {
		Opts    LoadBalancerAddServiceOpts
		Request schema.LoadBalancerActionAddServiceRequest
	}{
		"minimal": {
			Opts: LoadBalancerAddServiceOpts{
				Protocol: LoadBalancerServiceProtocolHTTP,
			},
			Request: schema.LoadBalancerActionAddServiceRequest{
				Protocol: string(LoadBalancerServiceProtocolHTTP),
			},
		},
		"all set": {
			Opts: LoadBalancerAddServiceOpts{
				Protocol:        LoadBalancerServiceProtocolHTTP,
				DestinationPort: Ptr(80),
				Proxyprotocol:   Ptr(true),
				HTTP: &LoadBalancerAddServiceOptsHTTP{
					CookieName:     Ptr("keks"),
					CookieLifetime: Ptr(5 * time.Minute),
					RedirectHTTP:   Ptr(true),
					StickySessions: Ptr(true),
					Certificates:   []*Certificate{{ID: 1}, {ID: 2}},
				},
				HealthCheck: &LoadBalancerAddServiceOptsHealthCheck{
					Protocol: LoadBalancerServiceProtocolHTTP,
					Port:     Ptr(80),
					Interval: Ptr(5 * time.Second),
					Timeout:  Ptr(1 * time.Second),
					Retries:  Ptr(3),
					HTTP: &LoadBalancerAddServiceOptsHealthCheckHTTP{
						Domain:      Ptr("example.com"),
						Path:        Ptr("/health"),
						Response:    Ptr("ok"),
						StatusCodes: []string{"2??", "3??"},
						TLS:         Ptr(true),
					},
				},
			},
			Request: schema.LoadBalancerActionAddServiceRequest{
				Protocol:        string(LoadBalancerServiceProtocolHTTP),
				DestinationPort: Ptr(80),
				Proxyprotocol:   Ptr(true),
				HTTP: &schema.LoadBalancerActionAddServiceRequestHTTP{
					CookieName:     Ptr("keks"),
					CookieLifetime: Ptr(5 * 60),
					RedirectHTTP:   Ptr(true),
					StickySessions: Ptr(true),
					Certificates:   Ptr([]int64{1, 2}),
				},
				HealthCheck: &schema.LoadBalancerActionAddServiceRequestHealthCheck{
					Protocol: string(LoadBalancerServiceProtocolHTTP),
					Port:     Ptr(80),
					Interval: Ptr(5),
					Timeout:  Ptr(1),
					Retries:  Ptr(3),
					HTTP: &schema.LoadBalancerActionAddServiceRequestHealthCheckHTTP{
						Domain:      Ptr("example.com"),
						Path:        Ptr("/health"),
						Response:    Ptr("ok"),
						StatusCodes: Ptr([]string{"2??", "3??"}),
						TLS:         Ptr(true),
					},
				},
			},
		},
		"no health check": {
			Opts: LoadBalancerAddServiceOpts{
				Protocol:        LoadBalancerServiceProtocolHTTP,
				DestinationPort: Ptr(80),
				Proxyprotocol:   Ptr(true),
				HTTP: &LoadBalancerAddServiceOptsHTTP{
					CookieName:     Ptr("keks"),
					CookieLifetime: Ptr(5 * time.Minute),
					RedirectHTTP:   Ptr(true),
					StickySessions: Ptr(true),
					Certificates:   []*Certificate{{ID: 1}, {ID: 2}},
				},
			},
			Request: schema.LoadBalancerActionAddServiceRequest{
				Protocol:        string(LoadBalancerServiceProtocolHTTP),
				DestinationPort: Ptr(80),
				Proxyprotocol:   Ptr(true),
				HTTP: &schema.LoadBalancerActionAddServiceRequestHTTP{
					CookieName:     Ptr("keks"),
					CookieLifetime: Ptr(5 * 60),
					RedirectHTTP:   Ptr(true),
					StickySessions: Ptr(true),
					Certificates:   Ptr([]int64{1, 2}),
				},
				HealthCheck: nil,
			},
		},
	}
	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			req := loadBalancerAddServiceOptsToSchema(testCase.Opts)
			if !cmp.Equal(testCase.Request, req) {
				t.Log(cmp.Diff(testCase.Request, req))
				t.Fail()
			}
		})
	}
}

func TestLoadBalancerUpdateServiceOptsToSchema(t *testing.T) {
	testCases := map[string]struct {
		Opts    LoadBalancerUpdateServiceOpts
		Request schema.LoadBalancerActionUpdateServiceRequest
	}{
		"empty": {
			Opts:    LoadBalancerUpdateServiceOpts{},
			Request: schema.LoadBalancerActionUpdateServiceRequest{},
		},
		"all set": {
			Opts: LoadBalancerUpdateServiceOpts{
				Protocol:        LoadBalancerServiceProtocolHTTP,
				DestinationPort: Ptr(80),
				Proxyprotocol:   Ptr(true),
				HTTP: &LoadBalancerUpdateServiceOptsHTTP{
					CookieName:     Ptr("keks"),
					CookieLifetime: Ptr(5 * time.Minute),
					RedirectHTTP:   Ptr(true),
					StickySessions: Ptr(true),
					Certificates:   []*Certificate{{ID: 1}, {ID: 2}},
				},
				HealthCheck: &LoadBalancerUpdateServiceOptsHealthCheck{
					Protocol: LoadBalancerServiceProtocolHTTP,
					Port:     Ptr(80),
					Interval: Ptr(5 * time.Second),
					Timeout:  Ptr(1 * time.Second),
					Retries:  Ptr(3),
					HTTP: &LoadBalancerUpdateServiceOptsHealthCheckHTTP{
						Domain:      Ptr("example.com"),
						Path:        Ptr("/health"),
						Response:    Ptr("ok"),
						StatusCodes: []string{"2??", "3??"},
						TLS:         Ptr(true),
					},
				},
			},
			Request: schema.LoadBalancerActionUpdateServiceRequest{
				Protocol:        Ptr(string(LoadBalancerServiceProtocolHTTP)),
				DestinationPort: Ptr(80),
				Proxyprotocol:   Ptr(true),
				HTTP: &schema.LoadBalancerActionUpdateServiceRequestHTTP{
					CookieName:     Ptr("keks"),
					CookieLifetime: Ptr(5 * 60),
					RedirectHTTP:   Ptr(true),
					StickySessions: Ptr(true),
					Certificates:   Ptr([]int64{1, 2}),
				},
				HealthCheck: &schema.LoadBalancerActionUpdateServiceRequestHealthCheck{
					Protocol: Ptr(string(LoadBalancerServiceProtocolHTTP)),
					Port:     Ptr(80),
					Interval: Ptr(5),
					Timeout:  Ptr(1),
					Retries:  Ptr(3),
					HTTP: &schema.LoadBalancerActionUpdateServiceRequestHealthCheckHTTP{
						Domain:      Ptr("example.com"),
						Path:        Ptr("/health"),
						Response:    Ptr("ok"),
						StatusCodes: Ptr([]string{"2??", "3??"}),
						TLS:         Ptr(true),
					},
				},
			},
		},
		"no health check": {
			Opts: LoadBalancerUpdateServiceOpts{
				Protocol:        LoadBalancerServiceProtocolHTTP,
				DestinationPort: Ptr(80),
				Proxyprotocol:   Ptr(true),
				HTTP: &LoadBalancerUpdateServiceOptsHTTP{
					CookieName:     Ptr("keks"),
					CookieLifetime: Ptr(5 * time.Minute),
					RedirectHTTP:   Ptr(true),
					StickySessions: Ptr(true),
					Certificates:   []*Certificate{{ID: 1}, {ID: 2}},
				},
			},
			Request: schema.LoadBalancerActionUpdateServiceRequest{
				Protocol:        Ptr(string(LoadBalancerServiceProtocolHTTP)),
				DestinationPort: Ptr(80),
				Proxyprotocol:   Ptr(true),
				HTTP: &schema.LoadBalancerActionUpdateServiceRequestHTTP{
					CookieName:     Ptr("keks"),
					CookieLifetime: Ptr(5 * 60),
					RedirectHTTP:   Ptr(true),
					StickySessions: Ptr(true),
					Certificates:   Ptr([]int64{1, 2}),
				},
				HealthCheck: nil,
			},
		},
	}
	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			req := loadBalancerUpdateServiceOptsToSchema(testCase.Opts)
			if !cmp.Equal(testCase.Request, req) {
				t.Log(cmp.Diff(testCase.Request, req))
				t.Fail()
			}
		})
	}
}

func TestServerMetricsFromSchema(t *testing.T) {
	tests := []struct {
		name        string
		respFn      func() *schema.ServerGetMetricsResponse
		expected    *ServerMetrics
		expectedErr string
	}{
		{
			name: "values not tuples",
			respFn: func() *schema.ServerGetMetricsResponse {
				var resp schema.ServerGetMetricsResponse

				resp.Metrics.Start = mustParseTime(t, "2017-01-01T00:00:00Z")
				resp.Metrics.End = mustParseTime(t, "2017-01-01T23:00:00Z")
				resp.Metrics.TimeSeries = map[string]schema.ServerTimeSeriesVals{
					"cpu": {
						Values: []interface{}{"some value"},
					},
				}

				return &resp
			},
			expectedErr: "failed to convert value to tuple: some value",
		},
		{
			name: "invalid tuple size",
			respFn: func() *schema.ServerGetMetricsResponse {
				var resp schema.ServerGetMetricsResponse

				resp.Metrics.Start = mustParseTime(t, "2017-01-01T00:00:00Z")
				resp.Metrics.End = mustParseTime(t, "2017-01-01T23:00:00Z")
				resp.Metrics.TimeSeries = map[string]schema.ServerTimeSeriesVals{
					"cpu": {
						Values: []interface{}{
							[]interface{}{1435781471.622, "43", "something else"},
						},
					},
				}

				return &resp
			},
			expectedErr: "invalid tuple size: 3: [1.435781471622e+09 43 something else]",
		},
		{
			name: "invalid time stamp",
			respFn: func() *schema.ServerGetMetricsResponse {
				var resp schema.ServerGetMetricsResponse

				resp.Metrics.Start = mustParseTime(t, "2017-01-01T00:00:00Z")
				resp.Metrics.End = mustParseTime(t, "2017-01-01T23:00:00Z")
				resp.Metrics.TimeSeries = map[string]schema.ServerTimeSeriesVals{
					"cpu": {
						Values: []interface{}{
							[]interface{}{"1435781471.622", "43"},
						},
					},
				}

				return &resp
			},
			expectedErr: "convert to float64: 1435781471.622",
		},
		{
			name: "invalid value",
			respFn: func() *schema.ServerGetMetricsResponse {
				var resp schema.ServerGetMetricsResponse

				resp.Metrics.Start = mustParseTime(t, "2017-01-01T00:00:00Z")
				resp.Metrics.End = mustParseTime(t, "2017-01-01T23:00:00Z")
				resp.Metrics.TimeSeries = map[string]schema.ServerTimeSeriesVals{
					"cpu": {
						Values: []interface{}{
							[]interface{}{1435781471.622, 43},
						},
					},
				}

				return &resp
			},
			expectedErr: "not a string: 43",
		},
		{
			name: "valid response",
			respFn: func() *schema.ServerGetMetricsResponse {
				var resp schema.ServerGetMetricsResponse

				resp.Metrics.Start = mustParseTime(t, "2017-01-01T00:00:00Z")
				resp.Metrics.End = mustParseTime(t, "2017-01-01T23:00:00Z")
				resp.Metrics.TimeSeries = map[string]schema.ServerTimeSeriesVals{
					"cpu": {
						Values: []interface{}{
							[]interface{}{1435781470.622, "42"},
							[]interface{}{1435781471.622, "43"},
						},
					},
					"disk.0.iops.read": {
						Values: []interface{}{
							[]interface{}{1435781480.622, "100"},
							[]interface{}{1435781481.622, "150"},
						},
					},
					"disk.0.iops.write": {
						Values: []interface{}{
							[]interface{}{1435781480.622, "50"},
							[]interface{}{1435781481.622, "55"},
						},
					},
					"network.0.pps.in": {
						Values: []interface{}{
							[]interface{}{1435781490.622, "70"},
							[]interface{}{1435781491.622, "75"},
						},
					},
					"network.0.pps.out": {
						Values: []interface{}{
							[]interface{}{1435781590.622, "60"},
							[]interface{}{1435781591.622, "65"},
						},
					},
				}

				return &resp
			},
			expected: &ServerMetrics{
				Start: mustParseTime(t, "2017-01-01T00:00:00Z"),
				End:   mustParseTime(t, "2017-01-01T23:00:00Z"),
				TimeSeries: map[string][]ServerMetricsValue{
					"cpu": {
						{Timestamp: 1435781470.622, Value: "42"},
						{Timestamp: 1435781471.622, Value: "43"},
					},
					"disk.0.iops.read": {
						{Timestamp: 1435781480.622, Value: "100"},
						{Timestamp: 1435781481.622, Value: "150"},
					},
					"disk.0.iops.write": {
						{Timestamp: 1435781480.622, Value: "50"},
						{Timestamp: 1435781481.622, Value: "55"},
					},
					"network.0.pps.in": {
						{Timestamp: 1435781490.622, Value: "70"},
						{Timestamp: 1435781491.622, Value: "75"},
					},
					"network.0.pps.out": {
						{Timestamp: 1435781590.622, Value: "60"},
						{Timestamp: 1435781591.622, Value: "65"},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			resp := tt.respFn()
			actual, err := serverMetricsFromSchema(resp)
			if err != nil && tt.expectedErr == "" {
				t.Fatalf("expected no error; got: %v", err)
			}
			if err != nil && tt.expectedErr != err.Error() {
				t.Fatalf("expected error: %s; got: %v", tt.expectedErr, err)
			}
			if !cmp.Equal(tt.expected, actual) {
				t.Errorf("unexpected result:\n%s", cmp.Diff(tt.expected, actual))
			}
		})
	}
}

func TestLoadBalancerMetricsFromSchema(t *testing.T) {
	tests := []struct {
		name        string
		respFn      func() *schema.LoadBalancerGetMetricsResponse
		expected    *LoadBalancerMetrics
		expectedErr string
	}{
		{
			name: "values not tuples",
			respFn: func() *schema.LoadBalancerGetMetricsResponse {
				var resp schema.LoadBalancerGetMetricsResponse

				resp.Metrics.Start = mustParseTime(t, "2017-01-01T00:00:00Z")
				resp.Metrics.End = mustParseTime(t, "2017-01-01T23:00:00Z")
				resp.Metrics.TimeSeries = map[string]schema.LoadBalancerTimeSeriesVals{
					"open_connections": {
						Values: []interface{}{"some value"},
					},
				}

				return &resp
			},
			expectedErr: "failed to convert value to tuple: some value",
		},
		{
			name: "invalid tuple size",
			respFn: func() *schema.LoadBalancerGetMetricsResponse {
				var resp schema.LoadBalancerGetMetricsResponse

				resp.Metrics.Start = mustParseTime(t, "2017-01-01T00:00:00Z")
				resp.Metrics.End = mustParseTime(t, "2017-01-01T23:00:00Z")
				resp.Metrics.TimeSeries = map[string]schema.LoadBalancerTimeSeriesVals{
					"open_connections": {
						Values: []interface{}{
							[]interface{}{1435781471.622, "43", "something else"},
						},
					},
				}

				return &resp
			},
			expectedErr: "invalid tuple size: 3: [1.435781471622e+09 43 something else]",
		},
		{
			name: "invalid time stamp",
			respFn: func() *schema.LoadBalancerGetMetricsResponse {
				var resp schema.LoadBalancerGetMetricsResponse

				resp.Metrics.Start = mustParseTime(t, "2017-01-01T00:00:00Z")
				resp.Metrics.End = mustParseTime(t, "2017-01-01T23:00:00Z")
				resp.Metrics.TimeSeries = map[string]schema.LoadBalancerTimeSeriesVals{
					"open_connections": {
						Values: []interface{}{
							[]interface{}{"1435781471.622", "43"},
						},
					},
				}

				return &resp
			},
			expectedErr: "convert to float64: 1435781471.622",
		},
		{
			name: "invalid value",
			respFn: func() *schema.LoadBalancerGetMetricsResponse {
				var resp schema.LoadBalancerGetMetricsResponse

				resp.Metrics.Start = mustParseTime(t, "2017-01-01T00:00:00Z")
				resp.Metrics.End = mustParseTime(t, "2017-01-01T23:00:00Z")
				resp.Metrics.TimeSeries = map[string]schema.LoadBalancerTimeSeriesVals{
					"open_connections": {
						Values: []interface{}{
							[]interface{}{1435781471.622, 43},
						},
					},
				}

				return &resp
			},
			expectedErr: "not a string: 43",
		},
		{
			name: "valid response",
			respFn: func() *schema.LoadBalancerGetMetricsResponse {
				var resp schema.LoadBalancerGetMetricsResponse

				resp.Metrics.Start = mustParseTime(t, "2017-01-01T00:00:00Z")
				resp.Metrics.End = mustParseTime(t, "2017-01-01T23:00:00Z")
				resp.Metrics.TimeSeries = map[string]schema.LoadBalancerTimeSeriesVals{
					"open_connections": {
						Values: []interface{}{
							[]interface{}{1435781470.622, "42"},
							[]interface{}{1435781471.622, "43"},
						},
					},
					"connections_per_second": {
						Values: []interface{}{
							[]interface{}{1435781480.622, "100"},
							[]interface{}{1435781481.622, "150"},
						},
					},
					"requests_per_second": {
						Values: []interface{}{
							[]interface{}{1435781480.622, "50"},
							[]interface{}{1435781481.622, "55"},
						},
					},
					"bandwidth.in": {
						Values: []interface{}{
							[]interface{}{1435781490.622, "70"},
							[]interface{}{1435781491.622, "75"},
						},
					},
					"bandwidth.out": {
						Values: []interface{}{
							[]interface{}{1435781590.622, "60"},
							[]interface{}{1435781591.622, "65"},
						},
					},
				}

				return &resp
			},
			expected: &LoadBalancerMetrics{
				Start: mustParseTime(t, "2017-01-01T00:00:00Z"),
				End:   mustParseTime(t, "2017-01-01T23:00:00Z"),
				TimeSeries: map[string][]LoadBalancerMetricsValue{
					"open_connections": {
						{Timestamp: 1435781470.622, Value: "42"},
						{Timestamp: 1435781471.622, Value: "43"},
					},
					"connections_per_second": {
						{Timestamp: 1435781480.622, Value: "100"},
						{Timestamp: 1435781481.622, Value: "150"},
					},
					"requests_per_second": {
						{Timestamp: 1435781480.622, Value: "50"},
						{Timestamp: 1435781481.622, Value: "55"},
					},
					"bandwidth.in": {
						{Timestamp: 1435781490.622, Value: "70"},
						{Timestamp: 1435781491.622, Value: "75"},
					},
					"bandwidth.out": {
						{Timestamp: 1435781590.622, Value: "60"},
						{Timestamp: 1435781591.622, Value: "65"},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			resp := tt.respFn()
			actual, err := loadBalancerMetricsFromSchema(resp)
			if err != nil && tt.expectedErr == "" {
				t.Fatalf("expected no error; got: %v", err)
			}
			if err != nil && tt.expectedErr != err.Error() {
				t.Fatalf("expected error: %s; got: %v", tt.expectedErr, err)
			}
			if !cmp.Equal(tt.expected, actual) {
				t.Errorf("unexpected result:\n%s", cmp.Diff(tt.expected, actual))
			}
		})
	}
}
