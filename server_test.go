package hcloud

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestServerUnmarshalJSON(t *testing.T) {
	data := []byte(`{
		"id": 1,
		"name": "server.example.com",
		"status": "running",
		"created": "2017-08-16T17:29:14+00:00",
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
		"server_type": {
			"id": 2
		},
		"outgoing_traffic": 123456,
		"ingoing_traffic": 7891011,
		"included_traffic": 654321,
		"backup_window": "22-02",
		"rescue_enabled": true,
		"iso": {
			"id": 4711,
			"name": "FreeBSD-11.0-RELEASE-amd64-dvd1",
			"description": "FreeBSD 11.0 x64",
			"type": "public"
		}
	}`)

	var v Server
	if err := json.Unmarshal(data, &v); err != nil {
		t.Fatal(err)
	}

	if v.ID != 1 {
		t.Errorf("unexpected ID: %v", v.ID)
	}
	if v.Name != "server.example.com" {
		t.Errorf("unexpected name: %v", v.Name)
	}
	if v.Status != ServerStatusRunning {
		t.Errorf("unexpected status: %v", v.Status)
	}
	if !v.Created.Equal(time.Date(2017, 8, 16, 17, 29, 14, 0, time.UTC)) {
		t.Errorf("unexpected created date: %v", v.Created)
	}
	if v.PublicNet.IPv4.IP != "1.2.3.4" {
		t.Errorf("unexpected public net IPv4 IP: %v", v.PublicNet.IPv4.IP)
	}
	if v.ServerType.ID != 2 {
		t.Errorf("unexpected server type ID: %v", v.ServerType.ID)
	}
	if v.IncludedTraffic != 654321 {
		t.Errorf("unexpected included traffic: %v", v.IncludedTraffic)
	}
	if v.OutgoingTraffic != 123456 {
		t.Errorf("unexpected outgoing traffic: %v", v.OutgoingTraffic)
	}
	if v.IngoingTraffic != 7891011 {
		t.Errorf("unexpected ingoing traffic: %v", v.IngoingTraffic)
	}
	if v.BackupWindow != "22-02" {
		t.Errorf("unexpected backup window: %v", v.BackupWindow)
	}
	if !v.RescueEnabled {
		t.Errorf("unexpected rescue enabled state: %v", v.RescueEnabled)
	}
	if v.ISO == nil || v.ISO.ID != 4711 {
		t.Errorf("unexpected ISO: %v", v.ISO)
	}
}

func TestServerUnmarshalJSONNoTraffic(t *testing.T) {
	data := []byte(`{
		"outgoing_traffic": null,
		"ingoing_traffic": null
	}`)

	var v Server
	if err := json.Unmarshal(data, &v); err != nil {
		t.Fatal(err)
	}

	if v.OutgoingTraffic != 0 {
		t.Errorf("unexpected outgoing traffic: %v", v.OutgoingTraffic)
	}
	if v.IngoingTraffic != 0 {
		t.Errorf("unexpected ingoing traffic: %v", v.IngoingTraffic)
	}
}

func TestServerPublicNetUnmarshalJSON(t *testing.T) {
	data := []byte(`{
		"ipv4": {
			"ip": "1.2.3.4",
			"blocked": false,
			"dns_ptr": "server.example.com"
		},
		"ipv6": {
        		"ip": "2a01:4f8:1c19:1403::/64",
        		"blocked": false,
        		"dns_ptr": []
      		},
      		"floating_ips": [4]
	}`)

	var v ServerPublicNet
	if err := json.Unmarshal(data, &v); err != nil {
		t.Fatal(err)
	}

	if v.IPv4.IP != "1.2.3.4" {
		t.Errorf("unexpected IPv4 IP: %v", v.IPv4.IP)
	}
	if v.IPv6.IP != "2a01:4f8:1c19:1403::/64" {
		t.Errorf("unexpected IPv6 IP: %v", v.IPv6.IP)
	}
	if len(v.FloatingIPs) != 1 || v.FloatingIPs[0].ID != 4 {
		t.Errorf("unexpected Floating IPs: %v", v.FloatingIPs)
	}
}

func TestServerPublicNetIPv4UnmarshalJSON(t *testing.T) {
	data := []byte(`{
		"ip": "1.2.3.4",
		"blocked": true,
		"dns_ptr": "server.example.com"
	}`)

	var v ServerPublicNetIPv4
	if err := json.Unmarshal(data, &v); err != nil {
		t.Fatal(err)
	}

	if v.IP != "1.2.3.4" {
		t.Errorf("unexpected IP: %v", v.IP)
	}
	if !v.Blocked {
		t.Errorf("unexpected blocked state: %v", v.Blocked)
	}
	if v.DNSPtr != "server.example.com" {
		t.Errorf("unexpected DNS ptr: %v", v.DNSPtr)
	}
}

func TestServerPublicNetIPv6UnmarshalJSON(t *testing.T) {
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

	var v ServerPublicNetIPv6
	if err := json.Unmarshal(data, &v); err != nil {
		t.Fatal(err)
	}

	if v.IP != "2a01:4f8:1c11:3400::/64" {
		t.Errorf("unexpected IP: %v", v.IP)
	}
	if !v.Blocked {
		t.Errorf("unexpected blocked state: %v", v.Blocked)
	}
	if len(v.DNSPtr) != 1 {
		t.Errorf("unexpected DNS ptr: %v", v.DNSPtr)
	}
}

func TestServerPublicNetIPv6DNSPtrUnmarshalJSON(t *testing.T) {
	data := []byte(`{
		"ip": "2a01:4f8:1c11:3400::1/64",
		"dns_ptr": "server01.example.com"
	}`)

	var v ServerPublicNetIPv6DNSPtr
	if err := json.Unmarshal(data, &v); err != nil {
		t.Fatal(err)
	}

	if v.IP != "2a01:4f8:1c11:3400::1/64" {
		t.Errorf("unexpected IP: %v", v.IP)
	}
	if v.DNSPtr != "server01.example.com" {
		t.Errorf("unexpected DNS ptr: %v", v.DNSPtr)
	}
}

func TestServersGet(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers/1", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{
			"server": {
				"id": 1
			}
		}`)
	})

	ctx := context.Background()
	server, _, err := env.Client.Server.Get(ctx, 1)
	if err != nil {
		t.Fatalf("Servers.Get failed: %s", err)
	}
	if server == nil {
		t.Fatal("no server")
	}
	if server.ID != 1 {
		t.Errorf("unexpected server ID: %v", server.ID)
	}
}

func TestServersList(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers", func(w http.ResponseWriter, r *http.Request) {
		if page := r.URL.Query().Get("page"); page != "2" {
			t.Errorf("expected page 2; got %q", page)
		}
		if perPage := r.URL.Query().Get("per_page"); perPage != "50" {
			t.Errorf("expected per_page 50; got %q", perPage)
		}
		fmt.Fprint(w, `{
			"servers": [
				{
					"id": 1
				},
				{
					"id": 2
				}
			]
		}`)
	})

	opts := ServerListOpts{}
	opts.Page = 2
	opts.PerPage = 50

	ctx := context.Background()
	servers, _, err := env.Client.Server.List(ctx, opts)
	if err != nil {
		t.Fatalf("Servers.List failed: %s", err)
	}
	if len(servers) != 2 {
		t.Fatal("expected 2 servers")
	}
}

func TestServersListAll(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch page := r.URL.Query().Get("page"); page {
		case "", "1":
			fmt.Fprint(w, `{
				"servers": [
					{
						"id": 1
					}
				],
				"meta": {
					"pagination": {
						"page": 1,
						"per_page": 1,
						"previous_page": null,
						"next_page": 2,
						"last_page": 3,
						"total_entries": 3
					}
				}
			}`)
		case "2":
			fmt.Fprint(w, `{
				"servers": [
					{
						"id": 2
					}
				],
				"meta": {
					"pagination": {
						"page": 2,
						"per_page": 1,
						"previous_page": 1,
						"next_page": 3,
						"last_page": 3,
						"total_entries": 3
					}
				}
			}`)
		case "3":
			fmt.Fprint(w, `{
				"servers": [
					{
						"id": 3
					}
				],
				"meta": {
					"pagination": {
						"page": 3,
						"per_page": 1,
						"previous_page": 2,
						"next_page": null,
						"last_page": 3,
						"total_entries": 3
					}
				}
			}`)
		default:
			panic("bad page")
		}
	})

	ctx := context.Background()
	servers, err := env.Client.Server.ListAll(ctx)
	if err != nil {
		t.Fatalf("Servers.List failed: %s", err)
	}
	if len(servers) != 3 {
		t.Fatalf("expected 3 servers; got %d", len(servers))
	}
	if servers[0].ID != 1 {
		t.Errorf("")
	}
}

func TestServersCreate(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{
			"server": {
				"id": 1
			}
		}`)
	})

	ctx := context.Background()
	result, _, err := env.Client.Server.Create(ctx, ServerCreateOpts{
		Name:       "test",
		ServerType: ServerType{ID: 1},
		Image:      Image{ID: 2},
	})
	if err != nil {
		t.Fatalf("Server.Create failed: %s", err)
	}
	if result.Server == nil {
		t.Fatal("no server")
	}
	if result.Server.ID != 1 {
		t.Errorf("unexpected server ID: %v", result.Server.ID)
	}
}

func TestServersDelete(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers/1", func(w http.ResponseWriter, r *http.Request) {
		return
	})

	ctx := context.Background()
	_, err := env.Client.Server.Delete(ctx, 1)
	if err != nil {
		t.Fatalf("Server.Delete failed: %s", err)
	}
}
