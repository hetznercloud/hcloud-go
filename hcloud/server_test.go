package hcloud

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/hetznercloud/hcloud-go/hcloud/schema"
)

func TestServerFromSchema(t *testing.T) {
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
	if server.PublicNet.IPv4.IP != "1.2.3.4" {
		t.Errorf("unexpected public net IPv4 IP: %v", server.PublicNet.IPv4.IP)
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
	if !server.RescueEnabled {
		t.Errorf("unexpected rescue enabled state: %v", server.RescueEnabled)
	}
	if server.ISO == nil || server.ISO.ID != 4711 {
		t.Errorf("unexpected ISO: %v", server.ISO)
	}
}

func TestServerFromSchemaNoTraffic(t *testing.T) {
	data := []byte(`{
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

	var s schema.ServerPublicNet
	if err := json.Unmarshal(data, &s); err != nil {
		t.Fatal(err)
	}
	publicNet := ServerPublicNetFromSchema(s)

	if publicNet.IPv4.IP != "1.2.3.4" {
		t.Errorf("unexpected IPv4 IP: %v", publicNet.IPv4.IP)
	}
	if publicNet.IPv6.IP != "2a01:4f8:1c19:1403::/64" {
		t.Errorf("unexpected IPv6 IP: %v", publicNet.IPv6.IP)
	}
	if len(publicNet.FloatingIPs) != 1 || publicNet.FloatingIPs[0].ID != 4 {
		t.Errorf("unexpected Floating IPs: %v", publicNet.FloatingIPs)
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

	if ipv4.IP != "1.2.3.4" {
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

	if ipv6.IP != "2a01:4f8:1c11:3400::/64" {
		t.Errorf("unexpected IP: %v", ipv6.IP)
	}
	if !ipv6.Blocked {
		t.Errorf("unexpected blocked state: %v", ipv6.Blocked)
	}
	if len(ipv6.DNSPtr) != 1 {
		t.Errorf("unexpected DNS ptr: %v", ipv6.DNSPtr)
	}
}

func TestServerPublicNetIPv6DNSPtrFromSchema(t *testing.T) {
	data := []byte(`{
		"ip": "2a01:4f8:1c11:3400::1/64",
		"dns_ptr": "server01.example.com"
	}`)

	var s schema.ServerPublicNetIPv6DNSPtr
	if err := json.Unmarshal(data, &s); err != nil {
		t.Fatal(err)
	}
	dnsPtr := ServerPublicNetIPv6DNSPtrFromSchema(s)

	if dnsPtr.IP != "2a01:4f8:1c11:3400::1/64" {
		t.Errorf("unexpected IP: %v", dnsPtr.IP)
	}
	if dnsPtr.DNSPtr != "server01.example.com" {
		t.Errorf("unexpected DNS ptr: %v", dnsPtr.DNSPtr)
	}
}

func TestServerClientGet(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers/1", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(schema.ServerGetResponse{
			Server: schema.Server{
				ID: 1,
			},
		})
	})

	ctx := context.Background()
	server, _, err := env.Client.Server.Get(ctx, 1)
	if err != nil {
		t.Fatal(err)
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
		json.NewEncoder(w).Encode(schema.ServerListResponse{
			Servers: []schema.Server{
				{ID: 1},
				{ID: 2},
			},
		})
	})

	opts := ServerListOpts{}
	opts.Page = 2
	opts.PerPage = 50

	ctx := context.Background()
	servers, _, err := env.Client.Server.List(ctx, opts)
	if err != nil {
		t.Fatal(err)
	}
	if len(servers) != 2 {
		t.Fatal("expected 2 servers")
	}
}

func TestServersAll(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	firstRequest := true
	env.Mux.HandleFunc("/servers", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if firstRequest {
			firstRequest = false
			w.WriteHeader(http.StatusTooManyRequests)
			fmt.Fprint(w, `{
				"error": {
					"code": "limit_reached",
					"message": "ratelimited"
				}
			}`)
			return
		}

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
	servers, err := env.Client.Server.All(ctx)
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
