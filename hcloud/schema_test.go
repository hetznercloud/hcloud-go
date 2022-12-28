package hcloud

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/hcloud-go/hcloud/schema"
)

func TestActionFromSchema(t *testing.T) {
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
	if err := json.Unmarshal(data, &s); err != nil {
		t.Fatal(err)
	}
	action := ActionFromSchema(s)

	if action.ID != 1 {
		t.Errorf("unexpected ID: %v", action.ID)
	}
	if action.Command != "create_server" {
		t.Errorf("unexpected command: %v", action.Command)
	}
	if action.Status != "success" {
		t.Errorf("unexpected status: %v", action.Status)
	}
	if action.Progress != 100 {
		t.Errorf("unexpected progress: %d", action.Progress)
	}
	if !action.Started.Equal(time.Date(2016, 1, 30, 23, 55, 0, 0, time.UTC)) {
		t.Errorf("unexpected started: %v", action.Started)
	}
	if !action.Finished.Equal(time.Date(2016, 1, 30, 23, 56, 13, 0, time.UTC)) {
		t.Errorf("unexpected finished: %v", action.Started)
	}
	if action.ErrorCode != "action_failed" {
		t.Errorf("unexpected error code: %v", action.ErrorCode)
	}
	if action.ErrorMessage != "Action failed" {
		t.Errorf("unexpected error message: %v", action.ErrorMessage)
	}
	if len(action.Resources) == 1 {
		if action.Resources[0].ID != 42 {
			t.Errorf("unexpected id in resources[0].ID: %v", action.Resources[0].ID)
		}
		if action.Resources[0].Type != ActionResourceTypeServer {
			t.Errorf("unexpected type in resources[0].Type: %v", action.Resources[0].Type)
		}
	} else {
		t.Errorf("unexpected number of resources")
	}
}

func TestActionsFromSchema(t *testing.T) {
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
	if err := json.Unmarshal(data, &s); err != nil {
		t.Fatal(err)
	}
	actions := ActionsFromSchema(s)
	if len(actions) != 2 || actions[0].ID != 13 || actions[1].ID != 14 {
		t.Fatal("unexpected actions")
	}
}

func TestFloatingIPFromSchema(t *testing.T) {
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
		if err := json.Unmarshal(data, &s); err != nil {
			t.Fatal(err)
		}
		floatingIP := FloatingIPFromSchema(s)

		if floatingIP.ID != 4711 {
			t.Errorf("unexpected ID: %v", floatingIP.ID)
		}
		if !floatingIP.Blocked {
			t.Errorf("unexpected value for Blocked: %v", floatingIP.Blocked)
		}
		if floatingIP.Name != "Web Frontend" {
			t.Errorf("unexpected name: %v", floatingIP.Name)
		}
		if floatingIP.Description != "Web Frontend" {
			t.Errorf("unexpected description: %v", floatingIP.Description)
		}
		if floatingIP.IP.String() != "2001:db8::" {
			t.Errorf("unexpected IP: %v", floatingIP.IP)
		}
		if floatingIP.Type != FloatingIPTypeIPv6 {
			t.Errorf("unexpected Type: %v", floatingIP.Type)
		}
		if floatingIP.Server != nil {
			t.Errorf("unexpected Server: %v", floatingIP.Server)
		}
		if floatingIP.DNSPtr == nil || floatingIP.DNSPtrForIP(floatingIP.IP) != "" {
			t.Errorf("unexpected DNS ptr: %v", floatingIP.DNSPtr)
		}
		if floatingIP.HomeLocation == nil || floatingIP.HomeLocation.ID != 1 {
			t.Errorf("unexpected home location: %v", floatingIP.HomeLocation)
		}
		if !floatingIP.Protection.Delete {
			t.Errorf("unexpected Protection.Delete: %v", floatingIP.Protection.Delete)
		}
		if floatingIP.Labels["key"] != "value" || floatingIP.Labels["key2"] != "value2" {
			t.Errorf("unexpected Labels: %v", floatingIP.Labels)
		}
		if !floatingIP.Created.Equal(time.Date(2017, 8, 16, 17, 29, 14, 0, time.UTC)) {
			t.Errorf("unexpected created date: %v", floatingIP.Created)
		}
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
		if err := json.Unmarshal(data, &s); err != nil {
			t.Fatal(err)
		}
		floatingIP := FloatingIPFromSchema(s)

		if floatingIP.ID != 4711 {
			t.Errorf("unexpected ID: %v", floatingIP.ID)
		}
		if floatingIP.Blocked {
			t.Errorf("unexpected value for Blocked: %v", floatingIP.Blocked)
		}
		if floatingIP.Description != "Web Frontend" {
			t.Errorf("unexpected description: %v", floatingIP.Description)
		}
		if floatingIP.IP.String() != "131.232.99.1" {
			t.Errorf("unexpected IP: %v", floatingIP.IP)
		}
		if floatingIP.Type != FloatingIPTypeIPv4 {
			t.Errorf("unexpected type: %v", floatingIP.Type)
		}
		if floatingIP.Server == nil || floatingIP.Server.ID != 42 {
			t.Errorf("unexpected server: %v", floatingIP.Server)
		}
		if floatingIP.DNSPtr == nil || floatingIP.DNSPtrForIP(floatingIP.IP) != "fip01.example.com" {
			t.Errorf("unexpected DNS ptr: %v", floatingIP.DNSPtr)
		}
		if floatingIP.HomeLocation == nil || floatingIP.HomeLocation.ID != 1 {
			t.Errorf("unexpected home location: %v", floatingIP.HomeLocation)
		}
	})
}

func TestPrimaryIPFromSchema(t *testing.T) {
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
		if err := json.Unmarshal(data, &s); err != nil {
			t.Fatal(err)
		}
		primaryIP := PrimaryIPFromSchema(s)

		if primaryIP.ID != 4711 {
			t.Errorf("unexpected ID: %v", primaryIP.ID)
		}
		if !primaryIP.Blocked {
			t.Errorf("unexpected value for Blocked: %v", primaryIP.Blocked)
		}
		if !primaryIP.AutoDelete {
			t.Errorf("unexpected value for AutoDelete: %v", primaryIP.AutoDelete)
		}
		if primaryIP.Name != "Web Frontend" {
			t.Errorf("unexpected name: %v", primaryIP.Name)
		}

		if primaryIP.IP.String() != "fe80::" {
			t.Errorf("unexpected IP: %v", primaryIP.IP)
		}
		if primaryIP.Type != PrimaryIPTypeIPv6 {
			t.Errorf("unexpected Type: %v", primaryIP.Type)
		}
		if primaryIP.AssigneeType != "server" {
			t.Errorf("unexpected AssigneeType: %v", primaryIP.AssigneeType)
		}
		if primaryIP.AssigneeID != 17 {
			t.Errorf("unexpected AssigneeID: %v", primaryIP.AssigneeID)
		}
		dnsPTR, err := primaryIP.GetDNSPtrForIP(primaryIP.IP)
		if err != nil {
			t.Fatal(err)
		}
		if primaryIP.DNSPtr == nil || dnsPTR == "" {
			t.Errorf("unexpected DNS ptr: %v", primaryIP.DNSPtr)
		}
		if primaryIP.Datacenter.Name != "fsn1-dc8" {
			t.Errorf("unexpected datacenter: %v", primaryIP.Datacenter)
		}
		if !primaryIP.Protection.Delete {
			t.Errorf("unexpected Protection.Delete: %v", primaryIP.Protection.Delete)
		}
		if primaryIP.Labels["key"] != "value" || primaryIP.Labels["key2"] != "value2" {
			t.Errorf("unexpected Labels: %v", primaryIP.Labels)
		}
		if !primaryIP.Created.Equal(time.Date(2017, 8, 16, 17, 29, 14, 0, time.UTC)) {
			t.Errorf("unexpected created date: %v", primaryIP.Created)
		}
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
		primaryIP := PrimaryIPFromSchema(s)

		if primaryIP.ID != 4711 {
			t.Errorf("unexpected ID: %v", primaryIP.ID)
		}
		if !primaryIP.Blocked {
			t.Errorf("unexpected value for Blocked: %v", primaryIP.Blocked)
		}
		if !primaryIP.AutoDelete {
			t.Errorf("unexpected value for AutoDelete: %v", primaryIP.AutoDelete)
		}
		if primaryIP.Name != "Web Frontend" {
			t.Errorf("unexpected name: %v", primaryIP.Name)
		}

		if primaryIP.IP.String() != "127.0.0.1" {
			t.Errorf("unexpected IP: %v", primaryIP.IP)
		}
		if primaryIP.Type != PrimaryIPTypeIPv4 {
			t.Errorf("unexpected Type: %v", primaryIP.Type)
		}
		if primaryIP.AssigneeType != "server" {
			t.Errorf("unexpected AssigneeType: %v", primaryIP.AssigneeType)
		}
		if primaryIP.AssigneeID != 17 {
			t.Errorf("unexpected AssigneeID: %v", primaryIP.AssigneeID)
		}
		dnsPTR, err := primaryIP.GetDNSPtrForIP(primaryIP.IP)
		if err != nil {
			t.Fatal(err)
		}
		if primaryIP.DNSPtr == nil || dnsPTR == "" {
			t.Errorf("unexpected DNS ptr: %v", primaryIP.DNSPtr)
		}
		if primaryIP.Datacenter.Name != "fsn1-dc8" {
			t.Errorf("unexpected datacenter: %v", primaryIP.Datacenter)
		}
		if !primaryIP.Protection.Delete {
			t.Errorf("unexpected Protection.Delete: %v", primaryIP.Protection.Delete)
		}
		if primaryIP.Labels["key"] != "value" || primaryIP.Labels["key2"] != "value2" {
			t.Errorf("unexpected Labels: %v", primaryIP.Labels)
		}
		if !primaryIP.Created.Equal(time.Date(2017, 8, 16, 17, 29, 14, 0, time.UTC)) {
			t.Errorf("unexpected created date: %v", primaryIP.Created)
		}
	})
}

func TestISOFromSchema(t *testing.T) {
	data := []byte(`{
		"id": 4711,
		"name": "FreeBSD-11.0-RELEASE-amd64-dvd1",
		"description": "FreeBSD 11.0 x64",
		"type": "public",
		"deprecated": "2018-02-28T00:00:00+00:00"
	}`)

	var s schema.ISO
	if err := json.Unmarshal(data, &s); err != nil {
		t.Fatal(err)
	}
	iso := ISOFromSchema(s)
	if iso.ID != 4711 {
		t.Errorf("unexpected ID: %v", iso.ID)
	}
	if iso.Name != "FreeBSD-11.0-RELEASE-amd64-dvd1" {
		t.Errorf("unexpected name: %v", iso.Name)
	}
	if iso.Description != "FreeBSD 11.0 x64" {
		t.Errorf("unexpected description: %v", iso.Description)
	}
	if iso.Type != ISOTypePublic {
		t.Errorf("unexpected type: %v", iso.Type)
	}
	if iso.Deprecated.IsZero() {
		t.Errorf("unexpected value for deprecated: %v", iso.Deprecated)
	}
}

func TestDatacenterFromSchema(t *testing.T) {
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
	if err := json.Unmarshal(data, &s); err != nil {
		t.Fatal(err)
	}
	datacenter := DatacenterFromSchema(s)
	if datacenter.ID != 1 {
		t.Errorf("unexpected ID: %v", datacenter.ID)
	}
	if datacenter.Name != "fsn1-dc8" {
		t.Errorf("unexpected Name: %v", datacenter.Name)
	}
	if datacenter.Location == nil || datacenter.Location.ID != 1 {
		t.Errorf("unexpected Location: %v", datacenter.Location)
	}
	if len(datacenter.ServerTypes.Available) != 4 {
		t.Errorf("unexpected ServerTypes.Available (should be 4): %v", len(datacenter.ServerTypes.Available))
	}
	if len(datacenter.ServerTypes.Supported) != 4 {
		t.Errorf("unexpected ServerTypes.Supported length (should be 4): %v", len(datacenter.ServerTypes.Supported))
	}
}

func TestLocationFromSchema(t *testing.T) {
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
	if err := json.Unmarshal(data, &s); err != nil {
		t.Fatal(err)
	}
	location := LocationFromSchema(s)
	if location.ID != 1 {
		t.Errorf("unexpected ID: %v", location.ID)
	}
	if location.Name != "fsn1" {
		t.Errorf("unexpected Name: %v", location.Name)
	}
	if location.Description != "Falkenstein DC Park 1" {
		t.Errorf("unexpected Description: %v", location.Description)
	}
	if location.Country != "DE" {
		t.Errorf("unexpected Country: %v", location.Country)
	}
	if location.City != "Falkenstein" {
		t.Errorf("unexpected City: %v", location.City)
	}
	if location.Latitude != 50.47612 {
		t.Errorf("unexpected Latitude: %v", location.Latitude)
	}
	if location.Longitude != 12.370071 {
		t.Errorf("unexpected Longitude: %v", location.Longitude)
	}
	if location.NetworkZone != "eu-central" {
		t.Errorf("unexpected NetworkZone: %v", location.NetworkZone)
	}
}

func TestServerFromSchema(t *testing.T) {
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
	if err := json.Unmarshal(data, &s); err != nil {
		t.Fatal(err)
	}
	server := ServerFromSchema(s)

	if server.ID != 1 {
		t.Errorf("unexpected ID: %v", server.ID)
	}
	if server.Name != "server.example.com" {
		t.Errorf("unexpected name: %v", server.Name)
	}
	if server.Status != ServerStatusRunning {
		t.Errorf("unexpected status: %v", server.Status)
	}
	if !server.Created.Equal(time.Date(2017, 8, 16, 17, 29, 14, 0, time.UTC)) {
		t.Errorf("unexpected created date: %v", server.Created)
	}
	if !server.PublicNet.IPv4.IsUnspecified() {
		t.Errorf("unexpected public net IPv4: %v", server.PublicNet.IPv4)
	}
	if server.PublicNet.IPv6.IP.String() != "2a01:4f8:1c11:3400::" {
		t.Errorf("unexpected public net IPv6 IP: %v", server.PublicNet.IPv6.IP)
	}
	if server.ServerType.ID != 2 {
		t.Errorf("unexpected server type ID: %v", server.ServerType.ID)
	}
	if server.IncludedTraffic != 654321 {
		t.Errorf("unexpected included traffic: %v", server.IncludedTraffic)
	}
	if server.OutgoingTraffic != 123456 {
		t.Errorf("unexpected outgoing traffic: %v", server.OutgoingTraffic)
	}
	if server.IngoingTraffic != 7891011 {
		t.Errorf("unexpected ingoing traffic: %v", server.IngoingTraffic)
	}
	if server.BackupWindow != "22-02" {
		t.Errorf("unexpected backup window: %v", server.BackupWindow)
	}
	if server.PrimaryDiskSize != 20 {
		t.Errorf("unexpected primary disk size: %v", server.PrimaryDiskSize)
	}
	if !server.RescueEnabled {
		t.Errorf("unexpected rescue enabled state: %v", server.RescueEnabled)
	}
	if server.Image == nil || server.Image.ID != 4711 {
		t.Errorf("unexpected Image: %v", server.Image)
	}
	if server.ISO == nil || server.ISO.ID != 4711 {
		t.Errorf("unexpected ISO: %v", server.ISO)
	}
	if server.Datacenter == nil || server.Datacenter.ID != 1 {
		t.Errorf("unexpected Datacenter: %v", server.Datacenter)
	}
	if !server.Locked {
		t.Errorf("unexpected value for Locked: %v", server.Locked)
	}
	if !server.Protection.Delete {
		t.Errorf("unexpected value for Protection.Delete: %v", server.Protection.Delete)
	}
	if !server.Protection.Rebuild {
		t.Errorf("unexpected value for Protection.Rebuild: %v", server.Protection.Rebuild)
	}
	if server.Labels["key"] != "value" || server.Labels["key2"] != "value2" {
		t.Errorf("unexpected Labels: %v", server.Labels)
	}
	if len(s.Volumes) != 3 {
		t.Errorf("unexpected number of volumes: %v", len(s.Volumes))
	}
	if s.Volumes[0] != 123 || s.Volumes[1] != 456 || s.Volumes[2] != 789 {
		t.Errorf("unexpected volumes: %v", s.Volumes)
	}
	if len(server.PrivateNet) != 1 {
		t.Errorf("unexpected length of PrivateNet: %v", len(server.PrivateNet))
	}
	if server.PrivateNet[0].Network.ID != 4711 {
		t.Errorf("unexpected first private net: %v", server.PrivateNet[0])
	}
	if server.PlacementGroup.ID != 897 {
		t.Errorf("unexpected placement group: %d", server.PlacementGroup.ID)
	}
}

func TestServerFromSchemaNoTraffic(t *testing.T) {
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
	if err := json.Unmarshal(data, &s); err != nil {
		t.Fatal(err)
	}
	server := ServerFromSchema(s)

	if server.OutgoingTraffic != 0 {
		t.Errorf("unexpected outgoing traffic: %v", server.OutgoingTraffic)
	}
	if server.IngoingTraffic != 0 {
		t.Errorf("unexpected ingoing traffic: %v", server.IngoingTraffic)
	}
}

func TestServerPublicNetFromSchema(t *testing.T) {
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
	if err := json.Unmarshal(data, &s); err != nil {
		t.Fatal(err)
	}
	publicNet := ServerPublicNetFromSchema(s)
	if publicNet.IPv4.ID != 1 {
		t.Errorf("unexpected IPv4 ID: %v", publicNet.IPv4.ID)
	}
	if publicNet.IPv4.IP.String() != "1.2.3.4" {
		t.Errorf("unexpected IPv4 IP: %v", publicNet.IPv4.IP)
	}
	if publicNet.IPv6.ID != 2 {
		t.Errorf("unexpected IPv6 ID: %v", publicNet.IPv6.ID)
	}
	if publicNet.IPv6.Network.String() != "2a01:4f8:1c19:1403::/64" {
		t.Errorf("unexpected IPv6 IP: %v", publicNet.IPv6.IP)
	}
	if len(publicNet.FloatingIPs) != 1 || publicNet.FloatingIPs[0].ID != 4 {
		t.Errorf("unexpected Floating IPs: %v", publicNet.FloatingIPs)
	}
	if len(publicNet.Firewalls) != 1 || publicNet.Firewalls[0].Firewall.ID != 23 || publicNet.Firewalls[0].Status != FirewallStatusApplied {
		t.Errorf("unexpected Firewalls: %v", publicNet.Firewalls)
	}
}

func TestServerPublicNetIPv4FromSchema(t *testing.T) {
	data := []byte(`{
		"ip": "1.2.3.4",
		"blocked": true,
		"dns_ptr": "server.example.com"
	}`)

	var s schema.ServerPublicNetIPv4
	if err := json.Unmarshal(data, &s); err != nil {
		t.Fatal(err)
	}
	ipv4 := ServerPublicNetIPv4FromSchema(s)

	if ipv4.IP.String() != "1.2.3.4" {
		t.Errorf("unexpected IP: %v", ipv4.IP)
	}
	if !ipv4.Blocked {
		t.Errorf("unexpected blocked state: %v", ipv4.Blocked)
	}
	if ipv4.DNSPtr != "server.example.com" {
		t.Errorf("unexpected DNS ptr: %v", ipv4.DNSPtr)
	}
}

func TestServerPublicNetIPv6FromSchema(t *testing.T) {
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
	if err := json.Unmarshal(data, &s); err != nil {
		t.Fatal(err)
	}
	ipv6 := ServerPublicNetIPv6FromSchema(s)

	if ipv6.Network.String() != "2a01:4f8:1c11:3400::/64" {
		t.Errorf("unexpected IP: %v", ipv6.IP)
	}
	if !ipv6.Blocked {
		t.Errorf("unexpected blocked state: %v", ipv6.Blocked)
	}
	if len(ipv6.DNSPtr) != 1 {
		t.Errorf("unexpected DNS ptr: %v", ipv6.DNSPtr)
	}
}

func TestServerPrivateNetFromSchema(t *testing.T) {
	data := []byte(`{
		"network": 4711,
		"ip": "10.0.1.1",
		"alias_ips": [
			"10.0.1.2"
		],
		"mac_address": "86:00:ff:2a:7d:e1"
	}`)

	var s schema.ServerPrivateNet
	if err := json.Unmarshal(data, &s); err != nil {
		t.Fatal(err)
	}
	privateNet := ServerPrivateNetFromSchema(s)

	if privateNet.Network.ID != 4711 {
		t.Errorf("unexpected Network: %v", privateNet.Network)
	}
	if privateNet.IP.String() != "10.0.1.1" {
		t.Errorf("unexpected IP: %v", privateNet.IP)
	}
	if len(privateNet.Aliases) != 1 {
		t.Errorf("unexpected number of alias IPs: %v", len(privateNet.Aliases))
	}
	if privateNet.Aliases[0].String() != "10.0.1.2" {
		t.Errorf("unexpected alias IP: %v", privateNet.Aliases[0])
	}
	if privateNet.MACAddress != "86:00:ff:2a:7d:e1" {
		t.Errorf("unexpected mac address: %v", privateNet.MACAddress)
	}
}

func TestServerTypeFromSchema(t *testing.T) {
	data := []byte(`{
		"id": 1,
		"name": "cx10",
		"description": "description",
		"cores": 4,
		"memory": 1.0,
		"disk": 20,
		"storage_type": "local",
		"cpu_type": "shared",
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
	if err := json.Unmarshal(data, &s); err != nil {
		t.Fatal(err)
	}
	serverType := ServerTypeFromSchema(s)

	if serverType.ID != 1 {
		t.Errorf("unexpected ID: %v", serverType.ID)
	}
	if serverType.Name != "cx10" {
		t.Errorf("unexpected name: %q", serverType.Name)
	}
	if serverType.Description != "description" {
		t.Errorf("unexpected description: %q", serverType.Description)
	}
	if serverType.Cores != 4 {
		t.Errorf("unexpected cores: %v", serverType.Cores)
	}
	if serverType.Memory != 1.0 {
		t.Errorf("unexpected memory: %v", serverType.Memory)
	}
	if serverType.Disk != 20 {
		t.Errorf("unexpected disk: %v", serverType.Disk)
	}
	if serverType.StorageType != StorageTypeLocal {
		t.Errorf("unexpected storage type: %q", serverType.StorageType)
	}
	if serverType.CPUType != CPUTypeShared {
		t.Errorf("unexpected cpu type: %q", serverType.CPUType)
	}
	if len(serverType.Pricings) != 1 {
		t.Errorf("unexpected number of pricings: %d", len(serverType.Pricings))
	} else {
		if serverType.Pricings[0].Location.Name != "fsn1" {
			t.Errorf("unexpected location name: %v", serverType.Pricings[0].Location.Name)
		}
		if serverType.Pricings[0].Hourly.Net != "1" {
			t.Errorf("unexpected hourly net price: %v", serverType.Pricings[0].Hourly.Net)
		}
		if serverType.Pricings[0].Hourly.Gross != "1.19" {
			t.Errorf("unexpected hourly gross price: %v", serverType.Pricings[0].Hourly.Gross)
		}
		if serverType.Pricings[0].Monthly.Net != "1" {
			t.Errorf("unexpected monthly net price: %v", serverType.Pricings[0].Monthly.Net)
		}
		if serverType.Pricings[0].Monthly.Gross != "1.19" {
			t.Errorf("unexpected monthly gross price: %v", serverType.Pricings[0].Monthly.Gross)
		}
	}
}

func TestSSHKeyFromSchema(t *testing.T) {
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
	if err := json.Unmarshal(data, &s); err != nil {
		t.Fatal(err)
	}
	sshKey := SSHKeyFromSchema(s)

	if sshKey.ID != 2323 {
		t.Errorf("unexpected ID: %v", sshKey.ID)
	}
	if sshKey.Name != "My key" {
		t.Errorf("unexpected name: %v", sshKey.Name)
	}
	if sshKey.Fingerprint != "b7:2f:30:a0:2f:6c:58:6c:21:04:58:61:ba:06:3b:2c" {
		t.Errorf("unexpected fingerprint: %v", sshKey.Fingerprint)
	}
	if sshKey.PublicKey != "ssh-rsa AAAjjk76kgf...Xt" {
		t.Errorf("unexpected public key: %v", sshKey.PublicKey)
	}
	if sshKey.Labels["key"] != "value" || sshKey.Labels["key2"] != "value2" {
		t.Errorf("unexpected labels: %v", sshKey.Labels)
	}
	if !sshKey.Created.Equal(time.Date(2017, 8, 16, 17, 29, 14, 0, time.UTC)) {
		t.Errorf("unexpected created date: %v", sshKey.Created)
	}
}

func TestErrorFromSchema(t *testing.T) {
	t.Run("service_error", func(t *testing.T) {
		data := []byte(`{
			"code": "service_error",
			"message": "An error occurred",
			"details": {}
		}`)

		var s schema.Error
		if err := json.Unmarshal(data, &s); err != nil {
			t.Fatal(err)
		}
		err := ErrorFromSchema(s)

		if err.Code != "service_error" {
			t.Errorf("unexpected code: %v", err.Code)
		}
		if err.Message != "An error occurred" {
			t.Errorf("unexpected message: %v", err.Message)
		}
	})

	t.Run("invalid_input", func(t *testing.T) {
		data := []byte(`{
			"code": "invalid_input",
			"message": "invalid input",
			"details": {
				"fields": [
					{
						"name": "broken_field",
						"messages": ["is required"]
					}
				]
			}
		}`)

		var s schema.Error
		if err := json.Unmarshal(data, &s); err != nil {
			t.Fatal(err)
		}
		err := ErrorFromSchema(s)

		if err.Code != "invalid_input" {
			t.Errorf("unexpected Code: %v", err.Code)
		}
		if err.Message != "invalid input" {
			t.Errorf("unexpected Message: %v", err.Message)
		}
		if d, ok := err.Details.(ErrorDetailsInvalidInput); !ok {
			t.Fatalf("unexpected Details type (should be ErrorDetailsInvalidInput): %v", err.Details)
		} else {
			if len(d.Fields) != 1 {
				t.Fatalf("unexpected Details.Fields length (should be 1): %v", d.Fields)
			}
			if d.Fields[0].Name != "broken_field" {
				t.Errorf("unexpected Details.Fields[0].Name: %v", d.Fields[0].Name)
			}
			if len(d.Fields[0].Messages) != 1 {
				t.Fatalf("unexpected Details.Fields[0].Messages length (should be 1): %v", d.Fields[0].Messages)
			}
			if d.Fields[0].Messages[0] != "is required" {
				t.Errorf("unexpected Details.Fields[0].Messages[0]: %v", d.Fields[0].Messages[0])
			}
		}
	})
}

func TestPaginationFromSchema(t *testing.T) {
	data := []byte(`{
		"page": 2,
		"per_page": 25,
		"previous_page": 1,
		"next_page": 3,
		"last_page": 13,
		"total_entries": 322
	}`)

	var s schema.MetaPagination
	if err := json.Unmarshal(data, &s); err != nil {
		t.Fatal(err)
	}
	p := PaginationFromSchema(s)

	if p.Page != 2 {
		t.Errorf("unexpected page: %v", p.Page)
	}
	if p.PerPage != 25 {
		t.Errorf("unexpected per page: %v", p.PerPage)
	}
	if p.PreviousPage != 1 {
		t.Errorf("unexpected previous page: %v", p.PreviousPage)
	}
	if p.NextPage != 3 {
		t.Errorf("unexpected next page: %d", p.NextPage)
	}
	if p.LastPage != 13 {
		t.Errorf("unexpected last page: %d", p.LastPage)
	}
	if p.TotalEntries != 322 {
		t.Errorf("unexpected total entries: %d", p.TotalEntries)
	}
}

func TestImageFromSchema(t *testing.T) {
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
	if err := json.Unmarshal(data, &s); err != nil {
		t.Fatal(err)
	}
	image := ImageFromSchema(s)

	if image.ID != 4711 {
		t.Errorf("unexpected ID: %v", image.ID)
	}
	if image.Type != ImageTypeSystem {
		t.Errorf("unexpected Type: %v", image.Type)
	}
	if image.Status != ImageStatusAvailable {
		t.Errorf("unexpected Status: %v", image.Status)
	}
	if image.Name != "ubuntu16.04-standard-x64" {
		t.Errorf("unexpected Name: %v", image.Name)
	}
	if image.Description != "Ubuntu 16.04 Standard 64 bit" {
		t.Errorf("unexpected Description: %v", image.Description)
	}
	if image.ImageSize != 2.3 {
		t.Errorf("unexpected ImageSize: %v", image.ImageSize)
	}
	if image.DiskSize != 10 {
		t.Errorf("unexpected DiskSize: %v", image.DiskSize)
	}
	if !image.Created.Equal(time.Date(2016, 1, 30, 23, 55, 1, 0, time.UTC)) {
		t.Errorf("unexpected Created: %v", image.Created)
	}
	if !image.Deleted.Equal(time.Date(2016, 1, 30, 23, 55, 1, 0, time.UTC)) {
		t.Errorf("unexpected Deleted: %v", image.Deleted)
	}
	if image.CreatedFrom == nil || image.CreatedFrom.ID != 1 || image.CreatedFrom.Name != "my-server1" {
		t.Errorf("unexpected CreatedFrom: %v", image.CreatedFrom)
	}
	if image.BoundTo == nil || image.BoundTo.ID != 1 {
		t.Errorf("unexpected BoundTo: %v", image.BoundTo)
	}
	if image.OSVersion != "16.04" {
		t.Errorf("unexpected OSVersion: %v", image.OSVersion)
	}
	if image.OSFlavor != "ubuntu" {
		t.Errorf("unexpected OSFlavor: %v", image.OSFlavor)
	}
	if image.RapidDeploy {
		t.Errorf("unexpected RapidDeploy: %v", image.RapidDeploy)
	}
	if !image.Protection.Delete {
		t.Errorf("unexpected Protection.Delete: %v", image.Protection.Delete)
	}
	if image.Deprecated.IsZero() {
		t.Errorf("unexpected value for Deprecated: %v", image.Deprecated)
	}
	if image.Labels["key"] != "value" || image.Labels["key2"] != "value2" {
		t.Errorf("unexpected Labels: %v", image.Labels)
	}
}

func TestVolumeFromSchema(t *testing.T) {
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
	if err := json.Unmarshal(data, &s); err != nil {
		t.Fatal(err)
	}
	volume := VolumeFromSchema(s)
	if volume.ID != 4711 {
		t.Errorf("unexpected ID: %v", volume.ID)
	}
	if volume.Name != "db-storage" {
		t.Errorf("unexpected name: %v", volume.Name)
	}
	if volume.Status != VolumeStatusCreating {
		t.Errorf("unexpected status: %v", volume.Status)
	}
	if !volume.Created.Equal(time.Date(2016, 1, 30, 23, 50, 11, 0, time.UTC)) {
		t.Errorf("unexpected created date: %s", volume.Created)
	}
	if volume.Server == nil {
		t.Error("no server")
	}
	if volume.Server != nil && volume.Server.ID != 2 {
		t.Errorf("unexpected server ID: %v", volume.Server.ID)
	}
	if volume.Location == nil || volume.Location.ID != 1 {
		t.Errorf("unexpected location: %v", volume.Location)
	}
	if volume.Size != 42 {
		t.Errorf("unexpected size: %v", volume.Size)
	}
	if !volume.Protection.Delete {
		t.Errorf("unexpected value for delete protection: %v", volume.Protection.Delete)
	}
	if len(volume.Labels) != 2 {
		t.Errorf("unexpected number of labels: %d", len(volume.Labels))
	}
	if volume.Labels["key"] != "value" || volume.Labels["key2"] != "value2" {
		t.Errorf("unexpected labels: %v", volume.Labels)
	}
}

func TestNetworkFromSchema(t *testing.T) {
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
	if err := json.Unmarshal(data, &s); err != nil {
		t.Fatal(err)
	}
	network := NetworkFromSchema(s)
	if network.ID != 4711 {
		t.Errorf("unexpected ID: %v", network.ID)
	}
	if network.Name != "mynet" {
		t.Errorf("unexpected Name: %v", network.Name)
	}
	if !network.Created.Equal(time.Date(2017, 8, 16, 17, 29, 14, 0, time.UTC)) {
		t.Errorf("unexpected created date: %v", network.Created)
	}
	if network.IPRange.String() != "10.0.0.0/16" {
		t.Errorf("unexpected IPRange: %v", network.IPRange)
	}
	if len(network.Subnets) != 1 {
		t.Errorf("unexpected length of Subnets: %v", len(network.Subnets))
	}
	if len(network.Routes) != 1 {
		t.Errorf("unexpected length of Routes: %v", len(network.Routes))
	}
	if len(network.Servers) != 1 {
		t.Errorf("unexpected length of Servers: %v", len(network.Servers))
	}
	if network.Servers[0].ID != 4711 {
		t.Errorf("unexpected Server ID: %v", network.Servers[0].ID)
	}
	if network.Protection.Delete {
		t.Errorf("unexpected value for delete protection: %v", network.Protection.Delete)
	}
}

func TestNetworkSubnetFromSchema(t *testing.T) {
	t.Run("type server", func(t *testing.T) {
		data := []byte(`{
			"type": "server",
			"ip_range": "10.0.1.0/24",
			"network_zone": "eu-central",
			"gateway": "10.0.0.1"
		}`)
		var s schema.NetworkSubnet
		if err := json.Unmarshal(data, &s); err != nil {
			t.Fatal(err)
		}
		networkSubnet := NetworkSubnetFromSchema(s)
		if networkSubnet.NetworkZone != "eu-central" {
			t.Errorf("unexpected NetworkZone: %v", networkSubnet.NetworkZone)
		}
		if networkSubnet.Type != "server" {
			t.Errorf("unexpected Type: %v", networkSubnet.Type)
		}
		if networkSubnet.IPRange.String() != "10.0.1.0/24" {
			t.Errorf("unexpected IPRange: %v", networkSubnet.IPRange)
		}
		if networkSubnet.Gateway.String() != "10.0.0.1" {
			t.Errorf("unexpected Gateway: %v", networkSubnet.Gateway)
		}
		if networkSubnet.VSwitchID != 0 {
			t.Errorf("unexpected VSwitchID: %v", networkSubnet.VSwitchID)
		}
	})

	t.Run("type vswitch", func(t *testing.T) {
		data := []byte(`{
			"type": "vswitch",
			"ip_range": "10.0.1.0/24",
			"network_zone": "eu-central",
			"gateway": "10.0.0.1",
			"vswitch_id": 123
		}`)
		var s schema.NetworkSubnet
		if err := json.Unmarshal(data, &s); err != nil {
			t.Fatal(err)
		}
		networkSubnet := NetworkSubnetFromSchema(s)
		if networkSubnet.NetworkZone != "eu-central" {
			t.Errorf("unexpected NetworkZone: %v", networkSubnet.NetworkZone)
		}
		if networkSubnet.Type != "vswitch" {
			t.Errorf("unexpected Type: %v", networkSubnet.Type)
		}
		if networkSubnet.IPRange.String() != "10.0.1.0/24" {
			t.Errorf("unexpected IPRange: %v", networkSubnet.IPRange)
		}
		if networkSubnet.Gateway.String() != "10.0.0.1" {
			t.Errorf("unexpected Gateway: %v", networkSubnet.Gateway)
		}
		if networkSubnet.VSwitchID != 123 {
			t.Errorf("unexpected VSwitchID: %v", networkSubnet.VSwitchID)
		}
	})
}

func TestNetworkRouteFromSchema(t *testing.T) {
	data := []byte(`{
		"destination": "10.100.1.0/24",
		"gateway": "10.0.1.1"
	}`)
	var s schema.NetworkRoute
	if err := json.Unmarshal(data, &s); err != nil {
		t.Fatal(err)
	}
	networkRoute := NetworkRouteFromSchema(s)
	if networkRoute.Destination.String() != "10.100.1.0/24" {
		t.Errorf("unexpected Destination: %v", networkRoute.Destination)
	}
	if networkRoute.Gateway.String() != "10.0.1.1" {
		t.Errorf("unexpected Gateway: %v", networkRoute.Gateway)
	}
}

func TestLoadBalancerTypeFromSchema(t *testing.T) {
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
	if err := json.Unmarshal(data, &s); err != nil {
		t.Fatal(err)
	}
	loadBalancerType := LoadBalancerTypeFromSchema(s)
	if loadBalancerType.ID != 1 {
		t.Errorf("unexpected ID: %v", loadBalancerType.ID)
	}
	if loadBalancerType.Name != "lx11" {
		t.Errorf("unexpected Name: %v", loadBalancerType.Name)
	}
	if loadBalancerType.Description != "LX11" {
		t.Errorf("unexpected Description: %v", loadBalancerType.Description)
	}
	if loadBalancerType.MaxConnections != 20000 {
		t.Errorf("unexpected MaxConnections: %v", loadBalancerType.MaxConnections)
	}
	if loadBalancerType.MaxServices != 3 {
		t.Errorf("unexpected MaxServices: %v", loadBalancerType.MaxServices)
	}
	if loadBalancerType.MaxTargets != 25 {
		t.Errorf("unexpected MaxTargets: %v", loadBalancerType.MaxTargets)
	}
	if loadBalancerType.MaxAssignedCertificates != 10 {
		t.Errorf("unexpected MaxAssignedCertificates: %v", loadBalancerType.MaxAssignedCertificates)
	}
	if len(loadBalancerType.Pricings) != 1 {
		t.Errorf("unexpected number of pricings: %d", len(loadBalancerType.Pricings))
	} else {
		if loadBalancerType.Pricings[0].Location.Name != "fsn1" {
			t.Errorf("unexpected location name: %v", loadBalancerType.Pricings[0].Location.Name)
		}
		if loadBalancerType.Pricings[0].Hourly.Net != "1" {
			t.Errorf("unexpected hourly net price: %v", loadBalancerType.Pricings[0].Hourly.Net)
		}
		if loadBalancerType.Pricings[0].Hourly.Gross != "1.19" {
			t.Errorf("unexpected hourly gross price: %v", loadBalancerType.Pricings[0].Hourly.Gross)
		}
		if loadBalancerType.Pricings[0].Monthly.Net != "1" {
			t.Errorf("unexpected monthly net price: %v", loadBalancerType.Pricings[0].Monthly.Net)
		}
		if loadBalancerType.Pricings[0].Monthly.Gross != "1.19" {
			t.Errorf("unexpected monthly gross price: %v", loadBalancerType.Pricings[0].Monthly.Gross)
		}
	}
}

func TestLoadBalancerFromSchema(t *testing.T) {
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
	if err := json.Unmarshal(data, &s); err != nil {
		t.Fatal(err)
	}
	loadBalancer := LoadBalancerFromSchema(s)
	if loadBalancer.ID != 4711 {
		t.Errorf("unexpected ID: %v", loadBalancer.ID)
	}
	if loadBalancer.Name != "Web Frontend" {
		t.Errorf("unexpected Name: %v", loadBalancer.Name)
	}
	if loadBalancer.PublicNet.IPv4.IP.String() != "131.232.99.1" {
		t.Errorf("unexpected IPv4: %v", loadBalancer.PublicNet.IPv4.IP)
	}
	if loadBalancer.PublicNet.IPv4.DNSPtr != "example.org" {
		t.Errorf("unexpected IPv4.DNSPtr: %v", loadBalancer.PublicNet.IPv4.DNSPtr)
	}
	if loadBalancer.PublicNet.IPv6.IP.String() != "2001:db8::1" {
		t.Errorf("unexpected IPv6: %v", loadBalancer.PublicNet.IPv6)
	}
	if loadBalancer.PublicNet.IPv6.DNSPtr != "example.com" {
		t.Errorf("unexpected IPv6.DNSPtr: %v", loadBalancer.PublicNet.IPv6.DNSPtr)
	}
	if len(loadBalancer.PrivateNet) != 1 {
		t.Errorf("unexpected length of PrivateNet: %v", len(loadBalancer.PrivateNet))
	} else {
		if loadBalancer.PrivateNet[0].Network.ID != 4711 {
			t.Errorf("unexpected Network ID: %v", loadBalancer.PrivateNet[0].Network.ID)
		}
		if loadBalancer.PrivateNet[0].IP.String() != "10.0.255.1" {
			t.Errorf("unexpected Network IP: %v", loadBalancer.PrivateNet[0].IP)
		}
	}
	if loadBalancer.Location == nil || loadBalancer.Location.ID != 1 {
		t.Errorf("unexpected Location: %v", loadBalancer.Location)
	}
	if loadBalancer.LoadBalancerType == nil || loadBalancer.LoadBalancerType.ID != 1 {
		t.Errorf("unexpected LoadBalancerType: %v", loadBalancer.LoadBalancerType)
	}
	if loadBalancer.Protection.Delete {
		t.Errorf("unexpected value for delete protection: %v", loadBalancer.Protection.Delete)
	}
	if !loadBalancer.Created.Equal(time.Date(2016, 01, 30, 23, 50, 00, 0, time.UTC)) {
		t.Errorf("unexpected created date: %v", loadBalancer.Created)
	}
	if len(loadBalancer.Services) != 1 {
		t.Errorf("unexpected length of Services: %v", len(loadBalancer.Services))
	}
	if len(loadBalancer.Targets) != 2 {
		t.Errorf("unexpected length of Targets: %v", len(loadBalancer.Targets))
	}
	if loadBalancer.Algorithm.Type != "round_robin" {
		t.Errorf("unexpected Algorithm.Type: %v", loadBalancer.Algorithm.Type)
	}
	if loadBalancer.IncludedTraffic != 654321 {
		t.Errorf("unexpected included traffic: %v", loadBalancer.IncludedTraffic)
	}
	if loadBalancer.OutgoingTraffic != 123456 {
		t.Errorf("unexpected outgoing traffic: %v", loadBalancer.OutgoingTraffic)
	}
	if loadBalancer.IngoingTraffic != 7891011 {
		t.Errorf("unexpected ingoing traffic: %v", loadBalancer.IngoingTraffic)
	}
}

func TestLoadBalancerServiceFromSchema(t *testing.T) {
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
	if err := json.Unmarshal(data, &s); err != nil {
		t.Fatal(err)
	}
	loadBalancerService := LoadBalancerServiceFromSchema(s)
	if loadBalancerService.Protocol != "http" {
		t.Errorf("unexpected Protocol: %v", loadBalancerService.Protocol)
	}
	if loadBalancerService.ListenPort != 443 {
		t.Errorf("unexpected ListenPort: %v", loadBalancerService.ListenPort)
	}
	if loadBalancerService.DestinationPort != 80 {
		t.Errorf("unexpected DestinationPort: %v", loadBalancerService.DestinationPort)
	}
	if loadBalancerService.Proxyprotocol {
		t.Errorf("unexpected ProxyProtocol: %v", loadBalancerService.Proxyprotocol)
	}
	if loadBalancerService.HTTP.CookieName != "HCLBSTICKY" {
		t.Errorf("unexpected HTTP.CookieName: %v", loadBalancerService.HTTP.CookieName)
	}
	if loadBalancerService.HTTP.CookieLifetime.Seconds() != 300 {
		t.Errorf("unexpected HTTP.CookieLifetime: %v", loadBalancerService.HTTP.CookieLifetime.Seconds())
	}
	if loadBalancerService.HTTP.Certificates[0].ID != 897 {
		t.Errorf("unexpected Certificates[0].ID : %v", loadBalancerService.HTTP.Certificates[0].ID)
	}
	if !loadBalancerService.HTTP.RedirectHTTP {
		t.Errorf("unexpected HTTP.RedirectHTTP: %v", loadBalancerService.HTTP.RedirectHTTP)
	}

	if !loadBalancerService.HTTP.StickySessions {
		t.Errorf("unexpected HTTP.StickySessions: %v", loadBalancerService.HTTP.StickySessions)
	}
	if loadBalancerService.HealthCheck.Protocol != "http" {
		t.Errorf("unexpected HealthCheck.Protocol: %v", loadBalancerService.HealthCheck.Protocol)
	}
	if loadBalancerService.HealthCheck.Port != 4711 {
		t.Errorf("unexpected HealthCheck.Port: %v", loadBalancerService.HealthCheck.Port)
	}
	if loadBalancerService.HealthCheck.Interval.Seconds() != 15 {
		t.Errorf("unexpected HealthCheck.Interval: %v", loadBalancerService.HealthCheck.Interval)
	}
	if loadBalancerService.HealthCheck.Timeout.Seconds() != 10 {
		t.Errorf("unexpected HealthCheck.Timeout: %v", loadBalancerService.HealthCheck.Timeout)
	}
	if loadBalancerService.HealthCheck.Retries != 3 {
		t.Errorf("unexpected HealthCheck.Retries: %v", loadBalancerService.HealthCheck.Retries)
	}
	if loadBalancerService.HealthCheck.HTTP.Domain != "example.com" {
		t.Errorf("unexpected HealthCheck.HTTP.Domain: %v", loadBalancerService.HealthCheck.HTTP.Domain)
	}
	if loadBalancerService.HealthCheck.HTTP.Path != "/" {
		t.Errorf("unexpected HealthCheck.HTTP.Path: %v", loadBalancerService.HealthCheck.HTTP.Path)
	}
	if loadBalancerService.HealthCheck.HTTP.Response != "" {
		t.Errorf("unexpected HealthCheck.HTTP.Response: %v", loadBalancerService.HealthCheck.HTTP.Response)
	}
	if loadBalancerService.HealthCheck.HTTP.TLS {
		t.Errorf("unexpected HealthCheck.HTTP.TLS: %v", loadBalancerService.HealthCheck.HTTP.TLS)
	}
	if len(loadBalancerService.HealthCheck.HTTP.StatusCodes) != 2 {
		t.Errorf("unexpected len(HealthCheck.HTTP.StatusCodes): %v", len(loadBalancerService.HealthCheck.HTTP.StatusCodes))
	} else {
		if loadBalancerService.HealthCheck.HTTP.StatusCodes[0] != "200" {
			t.Errorf("unexpected HealthCheck.HTTP.StatusCodes[0]: %v", loadBalancerService.HealthCheck.HTTP.StatusCodes[0])
		}
		if loadBalancerService.HealthCheck.HTTP.StatusCodes[1] != "201" {
			t.Errorf("unexpected HealthCheck.HTTP.StatusCodes[1]: %v", loadBalancerService.HealthCheck.HTTP.StatusCodes[1])
		}
	}
}

func TestLoadBalancerTargetFromSchema(t *testing.T) {
	t.Run("server target", func(t *testing.T) {
		data := []byte(`{
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
	}`)
		var s schema.LoadBalancerTarget
		if err := json.Unmarshal(data, &s); err != nil {
			t.Fatal(err)
		}
		loadBalancerTarget := LoadBalancerTargetFromSchema(s)
		if loadBalancerTarget.Type != "server" {
			t.Errorf("unexpected Type: %v", loadBalancerTarget.Type)
		}
		if loadBalancerTarget.Server == nil || loadBalancerTarget.Server.Server.ID != 80 {
			t.Errorf("unexpected Server: %v", loadBalancerTarget.Server)
		}
		if loadBalancerTarget.LabelSelector != nil {
			t.Errorf("unexpected LabelSelector.Selector: %v", loadBalancerTarget.LabelSelector)
		}
		if loadBalancerTarget.UsePrivateIP {
			t.Errorf("unexpected UsePrivateIP: %v", loadBalancerTarget.UsePrivateIP)
		}
		if len(loadBalancerTarget.HealthStatus) != 1 {
			t.Errorf("unexpected Health Status length: %v", len(loadBalancerTarget.HealthStatus))
		} else {
			if loadBalancerTarget.HealthStatus[0].ListenPort != 443 {
				t.Errorf("unexpected HealthStatus[0].ListenPort: %v", loadBalancerTarget.HealthStatus[0].ListenPort)
			}
			if loadBalancerTarget.HealthStatus[0].Status != LoadBalancerTargetHealthStatusStatusHealthy {
				t.Errorf("unexpected HealthStatus[0].Status: %v", loadBalancerTarget.HealthStatus[0].Status)
			}
		}
	})
	t.Run("label_selector target", func(t *testing.T) {
		data := []byte(`{
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
		}`)
		var s schema.LoadBalancerTarget
		if err := json.Unmarshal(data, &s); err != nil {
			t.Fatal(err)
		}
		loadBalancerTarget := LoadBalancerTargetFromSchema(s)
		if loadBalancerTarget.Type != "label_selector" {
			t.Errorf("unexpected Type: %v", loadBalancerTarget.Type)
		}
		if loadBalancerTarget.LabelSelector == nil || loadBalancerTarget.LabelSelector.Selector != "lbt" {
			t.Errorf("unexpected LabelSelector: %v", loadBalancerTarget.LabelSelector)
		}
		if loadBalancerTarget.Server != nil {
			t.Errorf("unexpected LabelSelector.Server: %v", loadBalancerTarget.Server)
		}
		if len(loadBalancerTarget.Targets) != 1 {
			t.Errorf("unexpected Targets length: %v", len(loadBalancerTarget.Targets))
		} else {
			if loadBalancerTarget.Targets[0].Server == nil || loadBalancerTarget.Targets[0].Server.Server.ID != 80 {
				t.Errorf("unexpected loadBalancerTarget.Targets[0].Server.Server.ID: %v", loadBalancerTarget.Targets[0].Server.Server.ID)
			}
			if len(loadBalancerTarget.Targets[0].HealthStatus) != 1 {
				t.Errorf("unexpected Targets length: %v", len(loadBalancerTarget.Targets[0].HealthStatus))
			} else {
				if loadBalancerTarget.Targets[0].HealthStatus[0].ListenPort != 443 {
					t.Errorf("unexpected HealthStatus[0].ListenPort: %v", loadBalancerTarget.Targets[0].HealthStatus[0].ListenPort)
				}
				if loadBalancerTarget.Targets[0].HealthStatus[0].Status != LoadBalancerTargetHealthStatusStatusHealthy {
					t.Errorf("unexpected HealthStatus[0].Status: %v", loadBalancerTarget.Targets[0].HealthStatus[0].Status)
				}
			}
		}
	})

	t.Run("ip target", func(t *testing.T) {
		var s schema.LoadBalancerTarget

		data := []byte(`{
			"type": "ip",
			"ip": {
				"ip": "1.2.3.4"
			}
		}`)
		if err := json.Unmarshal(data, &s); err != nil {
			t.Fatal(err)
		}
		lbTgt := LoadBalancerTargetFromSchema(s)
		if lbTgt.Type != LoadBalancerTargetTypeIP {
			t.Errorf("unexpected Type: %s", lbTgt.Type)
		}
		if lbTgt.IP.IP != "1.2.3.4" {
			t.Errorf("unexpected IP: %s", lbTgt.IP.IP)
		}
	})
}

func TestCertificateFromSchema(t *testing.T) {
	tests := []struct {
		name     string
		data     string
		expected Certificate
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
			expected: Certificate{
				ID:             897,
				Name:           "my website cert",
				Type:           "uploaded",
				Certificate:    "-----BEGIN CERTIFICATE-----\n...",
				Created:        mustParseTime(t, apiTimestampFormat, "2016-01-30T23:50:00+00:00"),
				NotValidBefore: mustParseTime(t, apiTimestampFormat, "2016-01-30T23:51:00+00:00"),
				NotValidAfter:  mustParseTime(t, apiTimestampFormat, "2016-01-30T23:55:00+00:00"),
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
			expected: Certificate{
				ID:             898,
				Name:           "managed certificate",
				Type:           "managed",
				Certificate:    "-----BEGIN CERTIFICATE-----\n...",
				Created:        mustParseTime(t, apiTimestampFormat, "2016-01-30T23:50:00+00:00"),
				NotValidBefore: mustParseTime(t, apiTimestampFormat, "2016-01-30T23:51:00+00:00"),
				NotValidAfter:  mustParseTime(t, apiTimestampFormat, "2016-01-30T23:55:00+00:00"),
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
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var s schema.Certificate

			err := json.Unmarshal([]byte(tt.data), &s)
			assert.NoError(t, err)
			actual := CertificateFromSchema(s)
			assert.Equal(t, &tt.expected, actual)
		})
	}
}

func TestPricingFromSchema(t *testing.T) {
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
	if err := json.Unmarshal(data, &s); err != nil {
		t.Fatal(err)
	}
	pricing := PricingFromSchema(s)

	if pricing.Image.PerGBMonth.Currency != "EUR" {
		t.Errorf("unexpected Image.PerGBMonth.Currency: %v", pricing.Image.PerGBMonth.Currency)
	}
	if pricing.Image.PerGBMonth.VATRate != "19.00" {
		t.Errorf("unexpected Image.PerGBMonth.VATRate: %v", pricing.Image.PerGBMonth.VATRate)
	}
	if pricing.Image.PerGBMonth.Net != "1" {
		t.Errorf("unexpected Image.PerGBMonth.Net: %v", pricing.Image.PerGBMonth.Net)
	}
	if pricing.Image.PerGBMonth.Gross != "1.19" {
		t.Errorf("unexpected Image.PerGBMonth.Gross: %v", pricing.Image.PerGBMonth.Gross)
	}

	if pricing.FloatingIP.Monthly.Currency != "EUR" {
		t.Errorf("unexpected FloatingIP.Monthly.Currency: %v", pricing.FloatingIP.Monthly.Currency)
	}
	if pricing.FloatingIP.Monthly.VATRate != "19.00" {
		t.Errorf("unexpected FloatingIP.Monthly.VATRate: %v", pricing.FloatingIP.Monthly.VATRate)
	}
	if pricing.FloatingIP.Monthly.Net != "1" {
		t.Errorf("unexpected FloatingIP.Monthly.Net: %v", pricing.FloatingIP.Monthly.Net)
	}
	if pricing.FloatingIP.Monthly.Gross != "1.19" {
		t.Errorf("unexpected FloatingIP.Monthly.Gross: %v", pricing.FloatingIP.Monthly.Gross)
	}

	if len(pricing.FloatingIPs) != 1 {
		t.Errorf("unexpected number of Floating IPs: %d", len(pricing.FloatingIPs))
	} else {
		p := pricing.FloatingIPs[0]

		if p.Type != FloatingIPTypeIPv4 {
			t.Errorf("unexpected .Type: %s", p.Type)
		}
		if len(p.Pricings) != 1 {
			t.Errorf("unexpected number of prices: %d", len(p.Pricings))
		} else {
			if p.Pricings[0].Location.Name != "fsn1" {
				t.Errorf("unexpected Location.Name: %v", p.Pricings[0].Location.Name)
			}
			if p.Pricings[0].Monthly.Currency != "EUR" {
				t.Errorf("unexpected Monthly.Currency: %v", p.Pricings[0].Monthly.Currency)
			}
			if p.Pricings[0].Monthly.VATRate != "19.00" {
				t.Errorf("unexpected Monthly.VATRate: %v", p.Pricings[0].Monthly.VATRate)
			}
			if p.Pricings[0].Monthly.Net != "1" {
				t.Errorf("unexpected Monthly.Net: %v", p.Pricings[0].Monthly.Net)
			}
			if p.Pricings[0].Monthly.Gross != "1.19" {
				t.Errorf("unexpected Monthly.Gross: %v", p.Pricings[0].Monthly.Gross)
			}
		}
	}

	if len(pricing.PrimaryIPs) != 1 {
		t.Errorf("unexpected number of Primary IPs: %d", len(pricing.PrimaryIPs))
	} else {
		ip := pricing.PrimaryIPs[0]

		if ip.Type != "ipv4" {
			t.Errorf("unexpected .Type: %s", ip.Type)
		}
		if len(ip.Pricings) != 1 {
			t.Errorf("unexpected number of prices: %d", len(ip.Pricings))
		} else {
			if ip.Pricings[0].Location != "fsn1" {
				t.Errorf("unexpected Location: %v", ip.Pricings[0].Location)
			}
			if ip.Pricings[0].Monthly.Net != "1.0000000000" {
				t.Errorf("unexpected Monthly.Net: %v", ip.Pricings[0].Monthly.Net)
			}
			if ip.Pricings[0].Monthly.Gross != "1.1900000000000000" {
				t.Errorf("unexpected Monthly.Gross: %v", ip.Pricings[0].Monthly.Gross)
			}
			if ip.Pricings[0].Hourly.Net != "1.0000000000" {
				t.Errorf("unexpected Hourly.Net: %v", ip.Pricings[0].Hourly.Net)
			}
			if ip.Pricings[0].Hourly.Gross != "1.1900000000000000" {
				t.Errorf("unexpected Hourly.Gross: %v", ip.Pricings[0].Hourly.Gross)
			}
		}
	}

	if pricing.Volume.PerGBMonthly.Currency != "EUR" {
		t.Errorf("unexpected Traffic.PerTB.Currency: %v", pricing.Volume.PerGBMonthly.Currency)
	}
	if pricing.Volume.PerGBMonthly.VATRate != "19.00" {
		t.Errorf("unexpected Traffic.PerTB.VATRate: %v", pricing.Volume.PerGBMonthly.VATRate)
	}
	if pricing.Volume.PerGBMonthly.Net != "1" {
		t.Errorf("unexpected Traffic.PerTB.Net: %v", pricing.Volume.PerGBMonthly.Net)
	}
	if pricing.Volume.PerGBMonthly.Gross != "1.19" {
		t.Errorf("unexpected Traffic.PerTB.Gross: %v", pricing.Volume.PerGBMonthly.Gross)
	}

	if pricing.Traffic.PerTB.Currency != "EUR" {
		t.Errorf("unexpected Traffic.PerTB.Currency: %v", pricing.Traffic.PerTB.Currency)
	}
	if pricing.Traffic.PerTB.VATRate != "19.00" {
		t.Errorf("unexpected Traffic.PerTB.VATRate: %v", pricing.Traffic.PerTB.VATRate)
	}
	if pricing.Traffic.PerTB.Net != "1" {
		t.Errorf("unexpected Traffic.PerTB.Net: %v", pricing.Traffic.PerTB.Net)
	}
	if pricing.Traffic.PerTB.Gross != "1.19" {
		t.Errorf("unexpected Traffic.PerTB.Gross: %v", pricing.Traffic.PerTB.Gross)
	}

	if pricing.ServerBackup.Percentage != "20" {
		t.Errorf("unexpected ServerBackup.Percentage: %v", pricing.ServerBackup.Percentage)
	}

	if len(pricing.ServerTypes) != 1 {
		t.Errorf("unexpected number of server types: %d", len(pricing.ServerTypes))
	} else {
		p := pricing.ServerTypes[0]

		if p.ServerType.ID != 4 {
			t.Errorf("unexpected ServerType.ID: %d", p.ServerType.ID)
		}
		if p.ServerType.Name != "CX11" {
			t.Errorf("unexpected ServerType.Name: %v", p.ServerType.Name)
		}

		if len(p.Pricings) != 1 {
			t.Errorf("unexpected number of prices: %d", len(p.Pricings))
		} else {
			if p.Pricings[0].Location.Name != "fsn1" {
				t.Errorf("unexpected Location.Name: %v", p.Pricings[0].Location.Name)
			}

			if p.Pricings[0].Hourly.Currency != "EUR" {
				t.Errorf("unexpected Hourly.Currency: %v", p.Pricings[0].Hourly.Currency)
			}
			if p.Pricings[0].Hourly.VATRate != "19.00" {
				t.Errorf("unexpected Hourly.VATRate: %v", p.Pricings[0].Hourly.VATRate)
			}
			if p.Pricings[0].Hourly.Net != "1" {
				t.Errorf("unexpected Hourly.Net: %v", p.Pricings[0].Hourly.Net)
			}
			if p.Pricings[0].Hourly.Gross != "1.19" {
				t.Errorf("unexpected Hourly.Gross: %v", p.Pricings[0].Hourly.Gross)
			}

			if p.Pricings[0].Monthly.Currency != "EUR" {
				t.Errorf("unexpected Monthly.Currency: %v", p.Pricings[0].Monthly.Currency)
			}
			if p.Pricings[0].Monthly.VATRate != "19.00" {
				t.Errorf("unexpected Monthly.VATRate: %v", p.Pricings[0].Monthly.VATRate)
			}
			if p.Pricings[0].Monthly.Net != "1" {
				t.Errorf("unexpected Monthly.Net: %v", p.Pricings[0].Monthly.Net)
			}
			if p.Pricings[0].Monthly.Gross != "1.19" {
				t.Errorf("unexpected Monthly.Gross: %v", p.Pricings[0].Monthly.Gross)
			}
		}
	}

	if len(pricing.LoadBalancerTypes) != 1 {
		t.Errorf("unexpected number of Load Balancer types: %d", len(pricing.LoadBalancerTypes))
	} else {
		p := pricing.LoadBalancerTypes[0]

		if p.LoadBalancerType.ID != 4 {
			t.Errorf("unexpected LoadBalancerType.ID: %d", p.LoadBalancerType.ID)
		}
		if p.LoadBalancerType.Name != "LX11" {
			t.Errorf("unexpected LoadBalancerType.Name: %v", p.LoadBalancerType.Name)
		}

		if len(p.Pricings) != 1 {
			t.Errorf("unexpected number of prices: %d", len(p.Pricings))
		} else {
			if p.Pricings[0].Location.Name != "fsn1" {
				t.Errorf("unexpected Location.Name: %v", p.Pricings[0].Location.Name)
			}

			if p.Pricings[0].Hourly.Currency != "EUR" {
				t.Errorf("unexpected Hourly.Currency: %v", p.Pricings[0].Hourly.Currency)
			}
			if p.Pricings[0].Hourly.VATRate != "19.00" {
				t.Errorf("unexpected Hourly.VATRate: %v", p.Pricings[0].Hourly.VATRate)
			}
			if p.Pricings[0].Hourly.Net != "1" {
				t.Errorf("unexpected Hourly.Net: %v", p.Pricings[0].Hourly.Net)
			}
			if p.Pricings[0].Hourly.Gross != "1.19" {
				t.Errorf("unexpected Hourly.Gross: %v", p.Pricings[0].Hourly.Gross)
			}

			if p.Pricings[0].Monthly.Currency != "EUR" {
				t.Errorf("unexpected Monthly.Currency: %v", p.Pricings[0].Monthly.Currency)
			}
			if p.Pricings[0].Monthly.VATRate != "19.00" {
				t.Errorf("unexpected Monthly.VATRate: %v", p.Pricings[0].Monthly.VATRate)
			}
			if p.Pricings[0].Monthly.Net != "1" {
				t.Errorf("unexpected Monthly.Net: %v", p.Pricings[0].Monthly.Net)
			}
			if p.Pricings[0].Monthly.Gross != "1.19" {
				t.Errorf("unexpected Monthly.Gross: %v", p.Pricings[0].Monthly.Gross)
			}
		}
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
				Network:         Ptr(3),
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
							Certificates:   Ptr([]int{1, 2}),
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
					Certificates:   Ptr([]int{1, 2}),
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
					Certificates:   Ptr([]int{1, 2}),
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
					Certificates:   Ptr([]int{1, 2}),
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
					Certificates:   Ptr([]int{1, 2}),
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

				resp.Metrics.Start = mustParseTime(t, time.RFC3339, "2017-01-01T00:00:00Z")
				resp.Metrics.End = mustParseTime(t, time.RFC3339, "2017-01-01T23:00:00Z")
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

				resp.Metrics.Start = mustParseTime(t, time.RFC3339, "2017-01-01T00:00:00Z")
				resp.Metrics.End = mustParseTime(t, time.RFC3339, "2017-01-01T23:00:00Z")
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

				resp.Metrics.Start = mustParseTime(t, time.RFC3339, "2017-01-01T00:00:00Z")
				resp.Metrics.End = mustParseTime(t, time.RFC3339, "2017-01-01T23:00:00Z")
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

				resp.Metrics.Start = mustParseTime(t, time.RFC3339, "2017-01-01T00:00:00Z")
				resp.Metrics.End = mustParseTime(t, time.RFC3339, "2017-01-01T23:00:00Z")
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

				resp.Metrics.Start = mustParseTime(t, time.RFC3339, "2017-01-01T00:00:00Z")
				resp.Metrics.End = mustParseTime(t, time.RFC3339, "2017-01-01T23:00:00Z")
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
				Start: mustParseTime(t, time.RFC3339, "2017-01-01T00:00:00Z"),
				End:   mustParseTime(t, time.RFC3339, "2017-01-01T23:00:00Z"),
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

				resp.Metrics.Start = mustParseTime(t, time.RFC3339, "2017-01-01T00:00:00Z")
				resp.Metrics.End = mustParseTime(t, time.RFC3339, "2017-01-01T23:00:00Z")
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

				resp.Metrics.Start = mustParseTime(t, time.RFC3339, "2017-01-01T00:00:00Z")
				resp.Metrics.End = mustParseTime(t, time.RFC3339, "2017-01-01T23:00:00Z")
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

				resp.Metrics.Start = mustParseTime(t, time.RFC3339, "2017-01-01T00:00:00Z")
				resp.Metrics.End = mustParseTime(t, time.RFC3339, "2017-01-01T23:00:00Z")
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

				resp.Metrics.Start = mustParseTime(t, time.RFC3339, "2017-01-01T00:00:00Z")
				resp.Metrics.End = mustParseTime(t, time.RFC3339, "2017-01-01T23:00:00Z")
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

				resp.Metrics.Start = mustParseTime(t, time.RFC3339, "2017-01-01T00:00:00Z")
				resp.Metrics.End = mustParseTime(t, time.RFC3339, "2017-01-01T23:00:00Z")
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
				Start: mustParseTime(t, time.RFC3339, "2017-01-01T00:00:00Z"),
				End:   mustParseTime(t, time.RFC3339, "2017-01-01T23:00:00Z"),
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

func TestFirewallFromSchema(t *testing.T) {
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
	}
`)
	var f schema.Firewall
	if err := json.Unmarshal(data, &f); err != nil {
		t.Fatal(err)
	}
	firewall := FirewallFromSchema(f)

	if firewall.ID != 897 {
		t.Errorf("unexpected ID: %v", firewall.ID)
	}
	if firewall.Name != "my firewall" {
		t.Errorf("unexpected Name: %v", firewall.Name)
	}
	if firewall.Labels["key"] != "value" || firewall.Labels["key2"] != "value2" {
		t.Errorf("unexpected Labels: %v", firewall.Labels)
	}
	if !firewall.Created.Equal(time.Date(2016, 01, 30, 23, 50, 00, 0, time.UTC)) {
		t.Errorf("unexpected Created date: %v", firewall.Created)
	}
	if len(firewall.Rules) != 1 {
		t.Errorf("unexpected Rules count: %d", len(firewall.Rules))
	}
	if firewall.Rules[0].Direction != FirewallRuleDirectionIn {
		t.Errorf("unexpected Rule Direction: %s", firewall.Rules[0].Direction)
	}
	if len(firewall.Rules[0].SourceIPs) != 3 {
		t.Errorf("unexpected Rule SourceIPs count: %d", len(firewall.Rules[0].SourceIPs))
	}
	if len(firewall.Rules[0].DestinationIPs) != 3 {
		t.Errorf("unexpected Rule DestinationIPs count: %d", len(firewall.Rules[0].DestinationIPs))
	}
	if firewall.Rules[0].Protocol != FirewallRuleProtocolTCP {
		t.Errorf("unexpected Rule Protocol: %s", firewall.Rules[0].Protocol)
	}
	if *firewall.Rules[0].Port != "80" {
		t.Errorf("unexpected Rule Port: %s", *firewall.Rules[0].Port)
	}
	if *firewall.Rules[0].Description != "allow http in" {
		t.Errorf("unexpected Rule Description: %s", *firewall.Rules[0].Description)
	}
	if len(firewall.AppliedTo) != 2 {
		t.Errorf("unexpected UsedBy count: %d", len(firewall.AppliedTo))
	}
	if firewall.AppliedTo[0].Type != FirewallResourceTypeServer {
		t.Errorf("unexpected UsedBy Type: %s", firewall.AppliedTo[0].Type)
	}
	if firewall.AppliedTo[0].Server.ID != 42 {
		t.Errorf("unexpected UsedBy Server ID: %d", firewall.AppliedTo[0].Server.ID)
	}
	if firewall.AppliedTo[1].Type != FirewallResourceTypeLabelSelector {
		t.Errorf("unexpected UsedBy Type: %s", firewall.AppliedTo[0].Type)
	}
	if firewall.AppliedTo[1].LabelSelector.Selector != "a=b" {
		t.Errorf("unexpected UsedBy Label Selector: %s", firewall.AppliedTo[1].LabelSelector.Selector)
	}
}

func TestPlacementGroupFromSchema(t *testing.T) {
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
	}
`)

	var g schema.PlacementGroup
	if err := json.Unmarshal(data, &g); err != nil {
		t.Fatal(err)
	}
	placementGroup := PlacementGroupFromSchema(g)
	if placementGroup.ID != 897 {
		t.Errorf("unexpected ID %d", placementGroup.ID)
	}
	if placementGroup.Name != "my Placement Group" {
		t.Errorf("unexpected Name %s", placementGroup.Name)
	}
	if placementGroup.Labels["key"] != "value" {
		t.Errorf("unexpected Labels: %v", placementGroup.Labels)
	}
	if !placementGroup.Created.Equal(time.Date(2019, 01, 8, 12, 10, 00, 0, time.UTC)) {
		t.Errorf("unexpected Created date: %v", placementGroup.Created)
	}
	if len(placementGroup.Servers) != 2 {
		t.Errorf("unexpected Servers %v", placementGroup.Servers)
	}
	if placementGroup.Type != PlacementGroupTypeSpread {
		t.Errorf("unexpected Type %s", placementGroup.Type)
	}
}
