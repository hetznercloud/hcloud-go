package hcloud

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/hetznercloud/hcloud-go/hcloud/schema"
)

func TestFloatingIPFromSchema(t *testing.T) {
	data := []byte(`{
		"id": 4711,
		"description": "Web Frontend",
		"ip": "131.232.99.1",
		"type": "ipv4",
		"server": 42,
		"dns_ptr": "fip01.example.com",
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
	if floatingIP.Description != "Web Frontend" {
		t.Errorf("unexpected description: %v", floatingIP.Description)
	}
	if floatingIP.IP != "131.232.99.1" {
		t.Errorf("unexpected IP: %v", floatingIP.IP)
	}
	if floatingIP.Type != FloatingIPTypeIPv4 {
		t.Errorf("unexpected type: %v", floatingIP.Type)
	}
	if floatingIP.Server == nil || floatingIP.Server.ID != 42 {
		t.Errorf("unexpected server: %v", floatingIP.Server)
	}
	if floatingIP.DNSPtr == nil || floatingIP.DNSPtr[floatingIP.IP] != "fip01.example.com" {
		t.Errorf("unexpected DNS ptr: %v", floatingIP.DNSPtr)
	}
	if floatingIP.HomeLocation == nil || floatingIP.HomeLocation.ID != 1 {
		t.Errorf("unexpected home location: %v", floatingIP.HomeLocation)
	}
}

func TestFloatingIPClientGet(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/floating_ips/1", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(schema.FloatingIPGetResponse{
			FloatingIP: schema.FloatingIP{
				ID: 1,
			},
		})
	})

	ctx := context.Background()
	floatingIP, _, err := env.Client.FloatingIP.Get(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
	if floatingIP == nil {
		t.Fatal("no Floating IP")
	}
	if floatingIP.ID != 1 {
		t.Errorf("unexpected ID: %v", floatingIP.ID)
	}
}

func TestFloatingIPClientList(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/floating_ips", func(w http.ResponseWriter, r *http.Request) {
		if page := r.URL.Query().Get("page"); page != "2" {
			t.Errorf("expected page 2; got %q", page)
		}
		if perPage := r.URL.Query().Get("per_page"); perPage != "50" {
			t.Errorf("expected per_page 50; got %q", perPage)
		}
		json.NewEncoder(w).Encode(schema.FloatingIPListResponse{
			FloatingIPs: []schema.FloatingIP{
				{ID: 1},
				{ID: 2},
			},
		})
	})

	opts := FloatingIPListOpts{}
	opts.Page = 2
	opts.PerPage = 50

	ctx := context.Background()
	floatingIPs, _, err := env.Client.FloatingIP.List(ctx, opts)
	if err != nil {
		t.Fatal(err)
	}
	if len(floatingIPs) != 2 {
		t.Fatal("expected 2 Floating IPs")
	}
}
