package hcloud

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/hetznercloud/hcloud-go/v2/hcloud/schema"
)

func TestServerClientGetByID(t *testing.T) {
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

	server, _, err := env.Client.Server.GetByID(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
	if server == nil {
		t.Fatal("no server")
	}
	if server.ID != 1 {
		t.Errorf("unexpected server ID: %v", server.ID)
	}

	t.Run("called via Get", func(t *testing.T) {
		server, _, err := env.Client.Server.Get(ctx, "1")
		if err != nil {
			t.Fatal(err)
		}
		if server == nil {
			t.Fatal("no server")
		}
		if server.ID != 1 {
			t.Errorf("unexpected server ID: %v", server.ID)
		}
	})
}

func TestServerClientGetByIDNotFound(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers/1", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(schema.ErrorResponse{
			Error: schema.Error{
				Code: string(ErrorCodeNotFound),
			},
		})
	})

	ctx := context.Background()
	server, _, err := env.Client.Server.GetByID(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
	if server != nil {
		t.Fatal("expected no server")
	}
}

func TestServerClientGetByName(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RawQuery != "name=myserver" {
			t.Fatal("missing name query")
		}
		json.NewEncoder(w).Encode(schema.ServerListResponse{
			Servers: []schema.Server{
				{
					ID:   1,
					Name: "myserver",
				},
			},
		})
	})
	ctx := context.Background()

	server, _, err := env.Client.Server.GetByName(ctx, "myserver")
	if err != nil {
		t.Fatal(err)
	}
	if server == nil {
		t.Fatal("no server")
	}
	if server.ID != 1 {
		t.Errorf("unexpected server ID: %v", server.ID)
	}

	t.Run("via Get", func(t *testing.T) {
		server, _, err := env.Client.Server.Get(ctx, "myserver")
		if err != nil {
			t.Fatal(err)
		}
		if server == nil {
			t.Fatal("no server")
		}
		if server.ID != 1 {
			t.Errorf("unexpected server ID: %v", server.ID)
		}
	})
}

func TestServerClientGetByNameNotFound(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RawQuery != "name=myserver" {
			t.Fatal("missing name query")
		}
		json.NewEncoder(w).Encode(schema.ServerListResponse{
			Servers: []schema.Server{},
		})
	})

	ctx := context.Background()
	server, _, err := env.Client.Server.GetByName(ctx, "myserver")
	if err != nil {
		t.Fatal(err)
	}
	if server != nil {
		t.Fatal("unexpected server")
	}
}

func TestServerClientGetByNameEmpty(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	ctx := context.Background()
	server, _, err := env.Client.Server.GetByName(ctx, "")
	if err != nil {
		t.Fatal(err)
	}
	if server != nil {
		t.Fatal("unexpected server")
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

	env.Mux.HandleFunc("/servers", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(struct {
			Servers []schema.Server `json:"servers"`
			Meta    schema.Meta     `json:"meta"`
		}{
			Servers: []schema.Server{
				{ID: 1},
				{ID: 2},
				{ID: 3},
			},
			Meta: schema.Meta{
				Pagination: &schema.MetaPagination{
					Page:         1,
					LastPage:     1,
					PerPage:      3,
					TotalEntries: 3,
				},
			},
		})
	})

	ctx := context.Background()
	servers, err := env.Client.Server.All(ctx)
	if err != nil {
		t.Fatalf("Servers.List failed: %s", err)
	}
	if len(servers) != 3 {
		t.Fatalf("expected 3 servers; got %d", len(servers))
	}
	if servers[0].ID != 1 || servers[1].ID != 2 || servers[2].ID != 3 {
		t.Errorf("unexpected servers")
	}
}

func TestServersAllWithOpts(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers", func(w http.ResponseWriter, r *http.Request) {
		if labelSelector := r.URL.Query().Get("label_selector"); labelSelector != "key=value" {
			t.Errorf("unexpected label selector: %s", labelSelector)
		}
		if name := r.URL.Query().Get("name"); name != "my-server" {
			t.Errorf("unexpected name: %s", name)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(struct {
			Servers []schema.Server `json:"servers"`
			Meta    schema.Meta     `json:"meta"`
		}{
			Servers: []schema.Server{
				{ID: 1},
				{ID: 2},
				{ID: 3},
			},
			Meta: schema.Meta{
				Pagination: &schema.MetaPagination{
					Page:         1,
					LastPage:     1,
					PerPage:      3,
					TotalEntries: 3,
				},
			},
		})
	})

	ctx := context.Background()
	opts := ServerListOpts{ListOpts: ListOpts{LabelSelector: "key=value"}, Name: "my-server"}
	servers, err := env.Client.Server.AllWithOpts(ctx, opts)
	if err != nil {
		t.Fatalf("Servers.List failed: %s", err)
	}
	if len(servers) != 3 {
		t.Fatalf("expected 3 servers; got %d", len(servers))
	}
	if servers[0].ID != 1 || servers[1].ID != 2 || servers[2].ID != 3 {
		t.Errorf("unexpected servers")
	}
}

func TestServersCreateWithSSHKeys(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers", func(w http.ResponseWriter, r *http.Request) {
		var reqBody schema.ServerCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatal(err)
		}
		if len(reqBody.SSHKeys) != 2 || reqBody.SSHKeys[0] != 1 || reqBody.SSHKeys[1] != 2 {
			t.Errorf("unexpected SSH keys: %v", reqBody.SSHKeys)
		}
		json.NewEncoder(w).Encode(schema.ServerCreateResponse{
			Server: schema.Server{
				ID: 1,
			},
			NextActions: []schema.Action{
				{ID: 2},
			},
		})
	})

	ctx := context.Background()
	result, _, err := env.Client.Server.Create(ctx, ServerCreateOpts{
		Name:       "test",
		ServerType: &ServerType{ID: 1},
		Image:      &Image{ID: 2},
		SSHKeys: []*SSHKey{
			{ID: 1},
			{ID: 2},
		},
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
	if result.RootPassword != "" {
		t.Errorf("expected no root password, got: %v", result.RootPassword)
	}
	if len(result.NextActions) != 1 || result.NextActions[0].ID != 2 {
		t.Errorf("unexpected next actions: %v", result.NextActions)
	}
}

func TestServersCreateWithoutSSHKeys(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers", func(w http.ResponseWriter, r *http.Request) {
		var reqBody schema.ServerCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatal(err)
		}
		if len(reqBody.SSHKeys) != 0 {
			t.Errorf("expected no SSH keys, but got %v", reqBody.SSHKeys)
		}
		json.NewEncoder(w).Encode(schema.ServerCreateResponse{
			Server: schema.Server{
				ID: 1,
			},
			RootPassword: Ptr("test"),
		})
	})

	ctx := context.Background()
	result, _, err := env.Client.Server.Create(ctx, ServerCreateOpts{
		Name:       "test",
		ServerType: &ServerType{ID: 1},
		Image:      &Image{ID: 2},
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
	if result.RootPassword != "test" {
		t.Errorf("unexpected root password: %v", result.RootPassword)
	}
}

func TestServersCreateWithVolumes(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers", func(w http.ResponseWriter, r *http.Request) {
		var reqBody schema.ServerCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatal(err)
		}
		if len(reqBody.Volumes) != 2 || reqBody.Volumes[0] != 1 || reqBody.Volumes[1] != 2 {
			t.Errorf("unexpected Volumes: %v", reqBody.Volumes)
		}
		if reqBody.Automount == nil || !*reqBody.Automount {
			t.Errorf("unexpected Automount: %v", reqBody.Automount)
		}
		json.NewEncoder(w).Encode(schema.ServerCreateResponse{
			Server: schema.Server{
				ID: 1,
			},
			NextActions: []schema.Action{
				{ID: 2},
			},
		})
	})

	ctx := context.Background()
	result, _, err := env.Client.Server.Create(ctx, ServerCreateOpts{
		Name:       "test",
		ServerType: &ServerType{ID: 1},
		Image:      &Image{ID: 2},
		Volumes: []*Volume{
			{ID: 1},
			{ID: 2},
		},
		Automount: Ptr(true),
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Server == nil {
		t.Fatal("no server")
	}
	if result.Server.ID != 1 {
		t.Errorf("unexpected server ID: %v", result.Server.ID)
	}
	if len(result.NextActions) != 1 || result.NextActions[0].ID != 2 {
		t.Errorf("unexpected next actions: %v", result.NextActions)
	}
}

func TestServersCreateWithNetworks(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers", func(w http.ResponseWriter, r *http.Request) {
		var reqBody schema.ServerCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatal(err)
		}
		if len(reqBody.Networks) != 2 || reqBody.Networks[0] != 1 || reqBody.Networks[1] != 2 {
			t.Errorf("unexpected Networks: %v", reqBody.Networks)
		}
		json.NewEncoder(w).Encode(schema.ServerCreateResponse{
			Server: schema.Server{
				ID: 1,
			},
			NextActions: []schema.Action{
				{ID: 2},
			},
		})
	})

	ctx := context.Background()
	result, _, err := env.Client.Server.Create(ctx, ServerCreateOpts{
		Name:       "test",
		ServerType: &ServerType{ID: 1},
		Image:      &Image{ID: 2},
		Networks: []*Network{
			{ID: 1},
			{ID: 2},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Server == nil {
		t.Fatal("no server")
	}
	if result.Server.ID != 1 {
		t.Errorf("unexpected server ID: %v", result.Server.ID)
	}
	if len(result.NextActions) != 1 || result.NextActions[0].ID != 2 {
		t.Errorf("unexpected next actions: %v", result.NextActions)
	}
}

func TestServersCreateWithPrivateNetworkOnly(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers", func(w http.ResponseWriter, r *http.Request) {
		var reqBody schema.ServerCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatal(err)
		}
		if len(reqBody.Networks) != 2 || reqBody.Networks[0] != 1 || reqBody.Networks[1] != 2 {
			t.Errorf("unexpected Networks: %v", reqBody.Networks)
		}
		if reqBody.PublicNet.EnableIPv4 != false {
			t.Errorf("unexpected PublicNet.EnableIPv4: %v", reqBody.PublicNet.EnableIPv4)
		}
		if reqBody.PublicNet.EnableIPv6 != false {
			t.Errorf("unexpected PublicNet.EnableIPv6: %v", reqBody.PublicNet.EnableIPv6)
		}
		if reqBody.PublicNet.IPv4ID != 0 {
			t.Errorf("unexpected PublicNet.IPv4: %v", reqBody.PublicNet.IPv4ID)
		}
		if reqBody.PublicNet.IPv6ID != 0 {
			t.Errorf("unexpected PublicNet.IPv6: %v", reqBody.PublicNet.IPv6ID)
		}
		json.NewEncoder(w).Encode(schema.ServerCreateResponse{
			Server: schema.Server{
				ID: 1,
			},
			NextActions: []schema.Action{
				{ID: 2},
			},
		})
	})

	ctx := context.Background()
	result, _, err := env.Client.Server.Create(ctx, ServerCreateOpts{
		Name:       "test",
		ServerType: &ServerType{ID: 1},
		Image:      &Image{ID: 2},
		Networks: []*Network{
			{ID: 1},
			{ID: 2},
		},
		PublicNet: &ServerCreatePublicNet{
			EnableIPv4: false,
			EnableIPv6: false,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Server == nil {
		t.Fatal("no server")
	}
	if result.Server.ID != 1 {
		t.Errorf("unexpected server ID: %v", result.Server.ID)
	}
	if len(result.NextActions) != 1 || result.NextActions[0].ID != 2 {
		t.Errorf("unexpected next actions: %v", result.NextActions)
	}
}

func TestServersCreateWithIPv6Only(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers", func(w http.ResponseWriter, r *http.Request) {
		var reqBody schema.ServerCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatal(err)
		}
		if reqBody.PublicNet.EnableIPv4 != false {
			t.Errorf("unexpected PublicNet.EnableIPv4: %v", reqBody.PublicNet.EnableIPv4)
		}
		if reqBody.PublicNet.EnableIPv6 != true {
			t.Errorf("unexpected PublicNet.EnableIPv6: %v", reqBody.PublicNet.EnableIPv6)
		}
		json.NewEncoder(w).Encode(schema.ServerCreateResponse{
			Server: schema.Server{
				ID: 1,
			},
			NextActions: []schema.Action{
				{ID: 2},
			},
		})
	})

	ctx := context.Background()
	result, _, err := env.Client.Server.Create(ctx, ServerCreateOpts{
		Name:       "test",
		ServerType: &ServerType{ID: 1},
		Image:      &Image{ID: 2},
		PublicNet:  &ServerCreatePublicNet{EnableIPv4: false, EnableIPv6: true},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Server == nil {
		t.Fatal("no server")
	}
	if result.Server.ID != 1 {
		t.Errorf("unexpected server ID: %v", result.Server.ID)
	}
	if len(result.NextActions) != 1 || result.NextActions[0].ID != 2 {
		t.Errorf("unexpected next actions: %v", result.NextActions)
	}
}

func TestServersCreateWithDefaultPublicNet(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers", func(w http.ResponseWriter, r *http.Request) {
		var reqBody schema.ServerCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatal(err)
		}
		if reqBody.PublicNet != nil {
			t.Errorf("unexpected PublicNet: %v", reqBody.PublicNet)
		}
		json.NewEncoder(w).Encode(schema.ServerCreateResponse{
			Server: schema.Server{
				ID: 1,
			},
			NextActions: []schema.Action{
				{ID: 2},
			},
		})
	})

	ctx := context.Background()
	result, _, err := env.Client.Server.Create(ctx, ServerCreateOpts{
		Name:       "test",
		ServerType: &ServerType{ID: 1},
		Image:      &Image{ID: 2},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Server == nil {
		t.Fatal("no server")
	}
	if result.Server.ID != 1 {
		t.Errorf("unexpected server ID: %v", result.Server.ID)
	}
	if len(result.NextActions) != 1 || result.NextActions[0].ID != 2 {
		t.Errorf("unexpected next actions: %v", result.NextActions)
	}
}

func TestServersCreateWithDatacenterID(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers", func(w http.ResponseWriter, r *http.Request) {
		var reqBody schema.ServerCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatal(err)
		}
		if reqBody.Datacenter != "1" {
			t.Errorf("unexpected datacenter: %v", reqBody.Datacenter)
		}
		json.NewEncoder(w).Encode(schema.ServerCreateResponse{
			Server: schema.Server{
				ID: 1,
			},
		})
	})

	ctx := context.Background()
	result, _, err := env.Client.Server.Create(ctx, ServerCreateOpts{
		Name:       "test",
		ServerType: &ServerType{ID: 1},
		Image:      &Image{ID: 2},
		Datacenter: &Datacenter{ID: 1},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Server == nil {
		t.Fatal("no server")
	}
}

func TestServersCreateWithDatacenterName(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers", func(w http.ResponseWriter, r *http.Request) {
		var reqBody schema.ServerCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatal(err)
		}
		if reqBody.Datacenter != "dc1" {
			t.Errorf("unexpected datacenter: %v", reqBody.Datacenter)
		}
		json.NewEncoder(w).Encode(schema.ServerCreateResponse{
			Server: schema.Server{
				ID: 1,
			},
		})
	})

	ctx := context.Background()
	result, _, err := env.Client.Server.Create(ctx, ServerCreateOpts{
		Name:       "test",
		ServerType: &ServerType{ID: 1},
		Image:      &Image{ID: 2},
		Datacenter: &Datacenter{Name: "dc1"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Server == nil {
		t.Fatal("no server")
	}
}

func TestServersCreateWithLocationID(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers", func(w http.ResponseWriter, r *http.Request) {
		var reqBody schema.ServerCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatal(err)
		}
		if reqBody.Location != "1" {
			t.Errorf("unexpected location: %v", reqBody.Location)
		}
		json.NewEncoder(w).Encode(schema.ServerCreateResponse{
			Server: schema.Server{
				ID: 1,
			},
		})
	})

	ctx := context.Background()
	result, _, err := env.Client.Server.Create(ctx, ServerCreateOpts{
		Name:       "test",
		ServerType: &ServerType{ID: 1},
		Image:      &Image{ID: 2},
		Location:   &Location{ID: 1},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Server == nil {
		t.Fatal("no server")
	}
}

func TestServersCreateWithLocationName(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers", func(w http.ResponseWriter, r *http.Request) {
		var reqBody schema.ServerCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatal(err)
		}
		if reqBody.Location != "loc1" {
			t.Errorf("unexpected location: %v", reqBody.Location)
		}
		json.NewEncoder(w).Encode(schema.ServerCreateResponse{
			Server: schema.Server{
				ID: 1,
			},
		})
	})

	ctx := context.Background()
	result, _, err := env.Client.Server.Create(ctx, ServerCreateOpts{
		Name:       "test",
		ServerType: &ServerType{ID: 1},
		Image:      &Image{ID: 2},
		Location:   &Location{Name: "loc1"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Server == nil {
		t.Fatal("no server")
	}
}

func TestServersCreateWithUserData(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers", func(w http.ResponseWriter, r *http.Request) {
		var reqBody schema.ServerCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatal(err)
		}
		if reqBody.UserData != "---user data---" {
			t.Errorf("unexpected userdata: %v", reqBody.UserData)
		}
		json.NewEncoder(w).Encode(schema.ServerCreateResponse{
			Server: schema.Server{
				ID: 1,
			},
		})
	})

	ctx := context.Background()
	result, _, err := env.Client.Server.Create(ctx, ServerCreateOpts{
		Name:       "test",
		ServerType: &ServerType{ID: 1},
		Image:      &Image{ID: 2},
		UserData:   "---user data---",
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Server == nil {
		t.Fatal("no server")
	}
}

func TestServersCreateWithFirewalls(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers", func(w http.ResponseWriter, r *http.Request) {
		var reqBody schema.ServerCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatal(err)
		}
		if len(reqBody.Firewalls) != 2 || reqBody.Firewalls[0].Firewall != 1 || reqBody.Firewalls[1].Firewall != 2 {
			t.Errorf("unexpected Firewalls: %v", reqBody.Firewalls)
		}
		json.NewEncoder(w).Encode(schema.ServerCreateResponse{
			Server: schema.Server{
				ID: 1,
			},
			NextActions: []schema.Action{
				{ID: 2},
			},
		})
	})

	ctx := context.Background()
	result, _, err := env.Client.Server.Create(ctx, ServerCreateOpts{
		Name:       "test",
		ServerType: &ServerType{ID: 1},
		Image:      &Image{ID: 2},
		Firewalls: []*ServerCreateFirewall{
			{Firewall: Firewall{ID: 1}},
			{Firewall: Firewall{ID: 2}},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Server == nil {
		t.Fatal("no server")
	}
	if result.Server.ID != 1 {
		t.Errorf("unexpected server ID: %v", result.Server.ID)
	}
	if len(result.NextActions) != 1 || result.NextActions[0].ID != 2 {
		t.Errorf("unexpected next actions: %v", result.NextActions)
	}
}

func TestServersCreateWithLabels(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers", func(w http.ResponseWriter, r *http.Request) {
		var reqBody schema.ServerCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatal(err)
		}
		if len(reqBody.SSHKeys) != 0 {
			t.Errorf("expected no SSH keys, but got %v", reqBody.SSHKeys)
		}
		if reqBody.Labels == nil || (*reqBody.Labels)["key"] != "value" {
			t.Errorf("unexpected labels in request: %v", reqBody.Labels)
		}
		json.NewEncoder(w).Encode(schema.ServerCreateResponse{
			Server: schema.Server{
				ID: 1,
			},
			RootPassword: Ptr("test"),
		})
	})

	ctx := context.Background()
	result, _, err := env.Client.Server.Create(ctx, ServerCreateOpts{
		Name:       "test",
		ServerType: &ServerType{ID: 1},
		Image:      &Image{ID: 2},
		Labels:     map[string]string{"key": "value"},
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

func TestServersCreateWithoutStarting(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers", func(w http.ResponseWriter, r *http.Request) {
		var reqBody schema.ServerCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatal(err)
		}
		if reqBody.StartAfterCreate == nil || *reqBody.StartAfterCreate {
			t.Errorf("unexpected value for start_after_create: %v", reqBody.StartAfterCreate)
		}
		json.NewEncoder(w).Encode(schema.ServerCreateResponse{
			Server: schema.Server{
				ID: 1,
			},
		})
	})

	ctx := context.Background()
	result, _, err := env.Client.Server.Create(ctx, ServerCreateOpts{
		Name:             "test",
		ServerType:       &ServerType{ID: 1},
		Image:            &Image{ID: 2},
		StartAfterCreate: Ptr(false),
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Server == nil {
		t.Fatal("no server")
	}
}

func TestServersCreateFailIfNoIPsAndNetworksAssigned(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()
	ctx := context.Background()
	_, _, err := env.Client.Server.Create(ctx, ServerCreateOpts{
		Name:       "test",
		ServerType: &ServerType{ID: 1},
		Image:      &Image{ID: 2},
		PublicNet: &ServerCreatePublicNet{
			EnableIPv4: false,
			EnableIPv6: false,
		},
	})
	if err == nil {
		t.Fatal(err)
	}
}

func TestServerCreateWithPlacementGroup(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers", func(w http.ResponseWriter, r *http.Request) {
		var reqBody schema.ServerCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatal(err)
		}
		if reqBody.PlacementGroup != 123 {
			t.Errorf("unexpected placement group id %d", reqBody.PlacementGroup)
		}
		json.NewEncoder(w).Encode(schema.ServerCreateResponse{
			Server: schema.Server{
				ID: 1,
				PlacementGroup: &schema.PlacementGroup{
					ID: 123,
				},
			},
			NextActions: []schema.Action{
				{ID: 2},
			},
		})
	})

	ctx := context.Background()
	result, _, err := env.Client.Server.Create(ctx, ServerCreateOpts{
		Name:           "test",
		ServerType:     &ServerType{ID: 1},
		Image:          &Image{ID: 2},
		PlacementGroup: &PlacementGroup{ID: 123},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Server == nil {
		t.Fatal("no server")
	}
	if result.Server.ID != 1 {
		t.Errorf("unexpected server ID: %d", result.Server.ID)
	}
	if len(result.NextActions) != 1 || result.NextActions[0].ID != 2 {
		t.Errorf("unexpected next actions: %v", result.NextActions)
	}
	if result.Server.PlacementGroup.ID != 123 {
		t.Errorf("unexpected placement group ID: %d", result.Server.PlacementGroup.ID)
	}
}

func TestServerCreateWithoutPrimaryIPsButNetwork(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers", func(w http.ResponseWriter, r *http.Request) {
		var reqBody schema.ServerCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatal(err)
		}
		if reqBody.PublicNet.EnableIPv4 != false && reqBody.PublicNet.EnableIPv6 != false {
			t.Errorf("unexpected public net %v", reqBody.PublicNet)
		}
		json.NewEncoder(w).Encode(schema.ServerCreateResponse{
			Server: schema.Server{
				ID: 1,
			},
			NextActions: []schema.Action{
				{ID: 2},
			},
		})
	})

	ctx := context.Background()
	result, _, err := env.Client.Server.Create(ctx, ServerCreateOpts{
		Name:       "test",
		ServerType: &ServerType{ID: 1},
		Image:      &Image{ID: 2},
		Networks: []*Network{
			{ID: 1},
		},
		PublicNet: &ServerCreatePublicNet{
			EnableIPv4: false,
			EnableIPv6: false,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Server == nil {
		t.Fatal("no server")
	}
	if result.Server.ID != 1 {
		t.Errorf("unexpected server ID: %d", result.Server.ID)
	}
	if len(result.NextActions) != 1 || result.NextActions[0].ID != 2 {
		t.Errorf("unexpected next actions: %v", result.NextActions)
	}
}

func TestServersDelete(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers/1", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Error("expected DELETE")
		}
		json.NewEncoder(w).Encode(schema.ServerDeleteResponse{
			Action: schema.Action{
				ID: 2,
			},
		})
	})

	var (
		ctx    = context.Background()
		server = &Server{ID: 1}
	)
	_, err := env.Client.Server.Delete(ctx, server)
	if err != nil {
		t.Fatalf("Server.Delete failed: %s", err)
	}
}

func TestServersDeleteWithResult(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers/1", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Error("expected DELETE")
		}
		json.NewEncoder(w).Encode(schema.ServerDeleteResponse{
			Action: schema.Action{
				ID: 2,
			},
		})
	})

	var (
		ctx    = context.Background()
		server = &Server{ID: 1}
	)
	result, _, err := env.Client.Server.DeleteWithResult(ctx, server)
	if err != nil {
		t.Fatalf("Server.Delete failed: %s", err)
	}
	if result.Action.ID != 2 {
		t.Errorf("unexpected action ID: %v", result.Action.ID)
	}
}

func TestServerClientUpdate(t *testing.T) {
	var (
		ctx    = context.Background()
		server = &Server{ID: 1}
	)

	t.Run("update name", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/servers/1", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "PUT" {
				t.Error("expected PUT")
			}
			var reqBody schema.ServerUpdateRequest
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				t.Fatal(err)
			}
			if reqBody.Name != "test" {
				t.Errorf("unexpected name: %v", reqBody.Name)
			}
			json.NewEncoder(w).Encode(schema.ServerUpdateResponse{
				Server: schema.Server{
					ID: 1,
				},
			})
		})

		opts := ServerUpdateOpts{
			Name: "test",
		}
		updatedServer, _, err := env.Client.Server.Update(ctx, server, opts)
		if err != nil {
			t.Fatal(err)
		}

		if updatedServer.ID != 1 {
			t.Errorf("unexpected server ID: %v", updatedServer.ID)
		}
	})

	t.Run("update labels", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/servers/1", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "PUT" {
				t.Error("expected PUT")
			}
			var reqBody schema.ServerUpdateRequest
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				t.Fatal(err)
			}
			if reqBody.Labels == nil || (*reqBody.Labels)["key"] != "value" {
				t.Errorf("unexpected labels in request: %v", reqBody.Labels)
			}
			json.NewEncoder(w).Encode(schema.ServerUpdateResponse{
				Server: schema.Server{
					ID: 1,
				},
			})
		})

		opts := ServerUpdateOpts{
			Labels: map[string]string{"key": "value"},
		}
		updatedServer, _, err := env.Client.Server.Update(ctx, server, opts)
		if err != nil {
			t.Fatal(err)
		}

		if updatedServer.ID != 1 {
			t.Errorf("unexpected server ID: %v", updatedServer.ID)
		}
	})

	t.Run("no updates", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/servers/1", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "PUT" {
				t.Error("expected PUT")
			}
			var reqBody schema.ServerUpdateRequest
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				t.Fatal(err)
			}
			if reqBody.Name != "" {
				t.Errorf("unexpected no name, but got: %v", reqBody.Name)
			}
			json.NewEncoder(w).Encode(schema.ServerUpdateResponse{
				Server: schema.Server{
					ID: 1,
				},
			})
		})

		opts := ServerUpdateOpts{}
		updatedServer, _, err := env.Client.Server.Update(ctx, server, opts)
		if err != nil {
			t.Fatal(err)
		}

		if updatedServer.ID != 1 {
			t.Errorf("unexpected server ID: %v", updatedServer.ID)
		}
	})
}

func TestServerClientPoweron(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers/1/actions/poweron", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(schema.ServerActionPoweronResponse{
			Action: schema.Action{
				ID: 1,
			},
		})
	})

	ctx := context.Background()
	action, _, err := env.Client.Server.Poweron(ctx, &Server{ID: 1})
	if err != nil {
		t.Fatal(err)
	}
	if action.ID != 1 {
		t.Errorf("unexpected action ID: %d", action.ID)
	}
}

func TestServerClientReboot(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers/1/actions/reboot", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(schema.ServerActionRebootResponse{
			Action: schema.Action{
				ID: 1,
			},
		})
	})

	ctx := context.Background()
	action, _, err := env.Client.Server.Reboot(ctx, &Server{ID: 1})
	if err != nil {
		t.Fatal(err)
	}
	if action.ID != 1 {
		t.Errorf("unexpected action ID: %d", action.ID)
	}
}

func TestServerClientReset(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers/1/actions/reset", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(schema.ServerActionResetResponse{
			Action: schema.Action{
				ID: 1,
			},
		})
	})

	ctx := context.Background()
	action, _, err := env.Client.Server.Reset(ctx, &Server{ID: 1})
	if err != nil {
		t.Fatal(err)
	}
	if action.ID != 1 {
		t.Errorf("unexpected action ID: %d", action.ID)
	}
}

func TestServerClientShutdown(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers/1/actions/shutdown", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(schema.ServerActionShutdownResponse{
			Action: schema.Action{
				ID: 1,
			},
		})
	})

	ctx := context.Background()
	action, _, err := env.Client.Server.Shutdown(ctx, &Server{ID: 1})
	if err != nil {
		t.Fatal(err)
	}
	if action.ID != 1 {
		t.Errorf("unexpected action ID: %d", action.ID)
	}
}

func TestServerClientPoweroff(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers/1/actions/poweroff", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(schema.ServerActionPoweroffResponse{
			Action: schema.Action{
				ID: 1,
			},
		})
	})

	ctx := context.Background()
	action, _, err := env.Client.Server.Poweroff(ctx, &Server{ID: 1})
	if err != nil {
		t.Fatal(err)
	}
	if action.ID != 1 {
		t.Errorf("unexpected action ID: %d", action.ID)
	}
}

func TestServerClientResetPassword(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers/1/actions/reset_password", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(schema.ServerActionResetPasswordResponse{
			Action: schema.Action{
				ID: 1,
			},
			RootPassword: "secret",
		})
	})

	ctx := context.Background()
	result, _, err := env.Client.Server.ResetPassword(ctx, &Server{ID: 1})
	if err != nil {
		t.Fatal(err)
	}
	if result.Action.ID != 1 {
		t.Errorf("unexpected action ID: %d", result.Action.ID)
	}
	if result.RootPassword != "secret" {
		t.Errorf("unexpected root password: %v", result.RootPassword)
	}
}

func TestServerClientCreateImageNoOptions(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers/1/actions/create_image", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(schema.ServerActionCreateImageResponse{
			Action: schema.Action{
				ID: 1,
			},
			Image: schema.Image{
				ID: 1,
			},
		})
	})

	ctx := context.Background()
	result, _, err := env.Client.Server.CreateImage(ctx, &Server{ID: 1}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if result.Action.ID != 1 {
		t.Errorf("unexpected action ID: %d", result.Action.ID)
	}
	if result.Image.ID != 1 {
		t.Errorf("unexpected image ID: %d", result.Image.ID)
	}
}

func TestServerClientCreateImageWithOptions(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers/1/actions/create_image", func(w http.ResponseWriter, r *http.Request) {
		var reqBody schema.ServerActionCreateImageRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatal(err)
		}
		if reqBody.Type == nil || *reqBody.Type != "backup" {
			t.Errorf("unexpected type: %v", reqBody.Type)
		}
		if reqBody.Description == nil || *reqBody.Description != "my backup" {
			t.Errorf("unexpected description: %v", reqBody.Description)
		}
		json.NewEncoder(w).Encode(schema.ServerActionCreateImageResponse{
			Action: schema.Action{
				ID: 1,
			},
			Image: schema.Image{
				ID: 1,
			},
		})
	})

	ctx := context.Background()
	opts := &ServerCreateImageOpts{
		Type:        ImageTypeBackup,
		Description: Ptr("my backup"),
	}
	result, _, err := env.Client.Server.CreateImage(ctx, &Server{ID: 1}, opts)
	if err != nil {
		t.Fatal(err)
	}
	if result.Action.ID != 1 {
		t.Errorf("unexpected action ID: %d", result.Action.ID)
	}
	if result.Image.ID != 1 {
		t.Errorf("unexpected image ID: %d", result.Image.ID)
	}
}

func TestServerClientEnableRescue(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers/1/actions/enable_rescue", func(w http.ResponseWriter, r *http.Request) {
		var reqBody schema.ServerActionEnableRescueRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatal(err)
		}
		if reqBody.Type == nil || *reqBody.Type != "linux64" {
			t.Errorf("unexpected type: %v", reqBody.Type)
		}
		if len(reqBody.SSHKeys) != 2 || reqBody.SSHKeys[0] != 1 || reqBody.SSHKeys[1] != 2 {
			t.Errorf("unexpected SSH keys: %v", reqBody.SSHKeys)
		}
		json.NewEncoder(w).Encode(schema.ServerActionEnableRescueResponse{
			Action: schema.Action{
				ID: 1,
			},
			RootPassword: "test",
		})
	})

	ctx := context.Background()
	opts := ServerEnableRescueOpts{
		Type: ServerRescueTypeLinux64,
		SSHKeys: []*SSHKey{
			{ID: 1},
			{ID: 2},
		},
	}
	result, _, err := env.Client.Server.EnableRescue(ctx, &Server{ID: 1}, opts)
	if err != nil {
		t.Fatal(err)
	}
	if result.Action.ID != 1 {
		t.Errorf("unexpected action ID: %d", result.Action.ID)
	}
	if result.RootPassword != "test" {
		t.Errorf("unexpected root password: %s", result.RootPassword)
	}
}

func TestServerClientDisableRescue(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers/1/actions/disable_rescue", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(schema.ServerActionDisableRescueResponse{
			Action: schema.Action{
				ID: 1,
			},
		})
	})

	ctx := context.Background()
	action, _, err := env.Client.Server.DisableRescue(ctx, &Server{ID: 1})
	if err != nil {
		t.Fatal(err)
	}
	if action.ID != 1 {
		t.Errorf("unexpected action ID: %d", action.ID)
	}
}

func TestServerClientRebuild(t *testing.T) {
	var (
		ctx    = context.Background()
		server = &Server{ID: 1}
	)

	t.Run("with image ID", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/servers/1/actions/rebuild", func(w http.ResponseWriter, r *http.Request) {
			var reqBody schema.ServerActionRebuildRequest
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				t.Fatal(err)
			}
			if id, ok := reqBody.Image.(float64); !ok || id != 1 {
				t.Errorf("unexpected image ID: %v", reqBody.Image)
			}
			json.NewEncoder(w).Encode(schema.ServerActionRebuildResponse{
				Action: schema.Action{
					ID: 1,
				},
			})
		})

		opts := ServerRebuildOpts{
			Image: &Image{ID: 1},
		}
		action, _, err := env.Client.Server.Rebuild(ctx, server, opts)
		if err != nil {
			t.Fatal(err)
		}
		if action.ID != 1 {
			t.Errorf("unexpected action ID: %d", action.ID)
		}
	})

	t.Run("with image name", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/servers/1/actions/rebuild", func(w http.ResponseWriter, r *http.Request) {
			var reqBody schema.ServerActionRebuildRequest
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				t.Fatal(err)
			}
			if name, ok := reqBody.Image.(string); !ok || name != "debian-9" {
				t.Errorf("unexpected image name: %v", reqBody.Image)
			}
			json.NewEncoder(w).Encode(schema.ServerActionRebuildResponse{
				Action: schema.Action{
					ID: 1,
				},
			})
		})

		opts := ServerRebuildOpts{
			Image: &Image{Name: "debian-9"},
		}
		action, _, err := env.Client.Server.Rebuild(ctx, server, opts)
		if err != nil {
			t.Fatal(err)
		}
		if action.ID != 1 {
			t.Errorf("unexpected action ID: %d", action.ID)
		}
	})
}

func TestServerClientRebuildWithResult(t *testing.T) {
	var (
		ctx    = context.Background()
		server = &Server{ID: 1}
	)

	t.Run("with image ID", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/servers/1/actions/rebuild", func(w http.ResponseWriter, r *http.Request) {
			var reqBody schema.ServerActionRebuildRequest
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				t.Fatal(err)
			}
			if id, ok := reqBody.Image.(float64); !ok || id != 1 {
				t.Errorf("unexpected image ID: %v", reqBody.Image)
			}
			json.NewEncoder(w).Encode(schema.ServerActionRebuildResponse{
				Action: schema.Action{
					ID: 1,
				},
				RootPassword: Ptr("hetzner"),
			})
		})

		opts := ServerRebuildOpts{
			Image: &Image{ID: 1},
		}
		result, _, err := env.Client.Server.RebuildWithResult(ctx, server, opts)
		if err != nil {
			t.Fatal(err)
		}
		if result.Action.ID != 1 {
			t.Errorf("unexpected action ID: %d", result.Action.ID)
		}
		if result.RootPassword != "hetzner" {
			t.Errorf("unexpected root password: %s", result.RootPassword)
		}
	})

	t.Run("with image name", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/servers/1/actions/rebuild", func(w http.ResponseWriter, r *http.Request) {
			var reqBody schema.ServerActionRebuildRequest
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				t.Fatal(err)
			}
			if name, ok := reqBody.Image.(string); !ok || name != "debian-9" {
				t.Errorf("unexpected image name: %v", reqBody.Image)
			}
			json.NewEncoder(w).Encode(schema.ServerActionRebuildResponse{
				Action: schema.Action{
					ID: 1,
				},
				RootPassword: nil,
			})
		})

		opts := ServerRebuildOpts{
			Image: &Image{Name: "debian-9"},
		}
		result, _, err := env.Client.Server.RebuildWithResult(ctx, server, opts)
		if err != nil {
			t.Fatal(err)
		}
		if result.Action.ID != 1 {
			t.Errorf("unexpected action ID: %d", result.Action.ID)
		}
		if result.RootPassword != "" {
			t.Errorf("unexpected root password: %s", result.RootPassword)
		}
	})
}

func TestServerClientAttachISO(t *testing.T) {
	var (
		ctx    = context.Background()
		server = &Server{ID: 1}
	)

	t.Run("with ISO ID", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/servers/1/actions/attach_iso", func(w http.ResponseWriter, r *http.Request) {
			var reqBody schema.ServerActionAttachISORequest
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				t.Fatal(err)
			}
			if id, ok := reqBody.ISO.(float64); !ok || id != 1 {
				t.Errorf("unexpected ISO ID: %v", reqBody.ISO)
			}
			json.NewEncoder(w).Encode(schema.ServerActionAttachISOResponse{
				Action: schema.Action{
					ID: 1,
				},
			})
		})

		iso := &ISO{ID: 1}
		action, _, err := env.Client.Server.AttachISO(ctx, server, iso)
		if err != nil {
			t.Fatal(err)
		}
		if action.ID != 1 {
			t.Errorf("unexpected action ID: %d", action.ID)
		}
	})

	t.Run("with ISO name", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/servers/1/actions/attach_iso", func(w http.ResponseWriter, r *http.Request) {
			var reqBody schema.ServerActionAttachISORequest
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				t.Fatal(err)
			}
			if name, ok := reqBody.ISO.(string); !ok || name != "debian.iso" {
				t.Errorf("unexpected ISO name: %v", reqBody.ISO)
			}
			json.NewEncoder(w).Encode(schema.ServerActionAttachISOResponse{
				Action: schema.Action{
					ID: 1,
				},
			})
		})

		iso := &ISO{Name: "debian.iso"}
		action, _, err := env.Client.Server.AttachISO(ctx, server, iso)
		if err != nil {
			t.Fatal(err)
		}
		if action.ID != 1 {
			t.Errorf("unexpected action ID: %d", action.ID)
		}
	})
}

func TestServerClientDetachISO(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	var (
		ctx    = context.Background()
		server = &Server{ID: 1}
	)

	env.Mux.HandleFunc("/servers/1/actions/detach_iso", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(schema.ServerActionDetachISOResponse{
			Action: schema.Action{
				ID: 1,
			},
		})
	})

	action, _, err := env.Client.Server.DetachISO(ctx, server)
	if err != nil {
		t.Fatal(err)
	}
	if action.ID != 1 {
		t.Errorf("unexpected action ID: %d", action.ID)
	}
}

func TestServerClientEnableBackup(t *testing.T) {
	var (
		ctx    = context.Background()
		server = &Server{ID: 1}
	)

	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers/1/actions/enable_backup", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(schema.ServerActionEnableBackupResponse{
			Action: schema.Action{
				ID: 1,
			},
		})
	})

	action, _, err := env.Client.Server.EnableBackup(ctx, server, "")
	if err != nil {
		t.Fatal(err)
	}
	if action.ID != 1 {
		t.Errorf("unexpected action ID: %d", action.ID)
	}
}

func TestServerClientDisableBackup(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	var (
		ctx    = context.Background()
		server = &Server{ID: 1}
	)

	env.Mux.HandleFunc("/servers/1/actions/disable_backup", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(schema.ServerActionDisableBackupResponse{
			Action: schema.Action{
				ID: 1,
			},
		})
	})

	action, _, err := env.Client.Server.DisableBackup(ctx, server)
	if err != nil {
		t.Fatal(err)
	}
	if action.ID != 1 {
		t.Errorf("unexpected action ID: %d", action.ID)
	}
}

func TestServerClientChangeType(t *testing.T) {
	var (
		ctx    = context.Background()
		server = &Server{ID: 1}
	)

	t.Run("with server type ID", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/servers/1/actions/change_type", func(w http.ResponseWriter, r *http.Request) {
			var reqBody schema.ServerActionChangeTypeRequest
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				t.Fatal(err)
			}
			if id, ok := reqBody.ServerType.(float64); !ok || id != 1 {
				t.Errorf("unexpected server type ID: %v", reqBody.ServerType)
			}
			if !reqBody.UpgradeDisk {
				t.Error("expected to upgrade disk")
			}
			json.NewEncoder(w).Encode(schema.ServerActionChangeTypeResponse{
				Action: schema.Action{
					ID: 1,
				},
			})
		})

		opts := ServerChangeTypeOpts{
			ServerType:  &ServerType{ID: 1},
			UpgradeDisk: true,
		}
		action, _, err := env.Client.Server.ChangeType(ctx, server, opts)
		if err != nil {
			t.Fatal(err)
		}
		if action.ID != 1 {
			t.Errorf("unexpected action ID: %d", action.ID)
		}
	})

	t.Run("with server type name", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/servers/1/actions/change_type", func(w http.ResponseWriter, r *http.Request) {
			var reqBody schema.ServerActionChangeTypeRequest
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				t.Fatal(err)
			}
			if name, ok := reqBody.ServerType.(string); !ok || name != "type" {
				t.Errorf("unexpected server type name: %v", reqBody.ServerType)
			}
			if !reqBody.UpgradeDisk {
				t.Error("expected to upgrade disk")
			}
			json.NewEncoder(w).Encode(schema.ServerActionChangeTypeResponse{
				Action: schema.Action{
					ID: 1,
				},
			})
		})

		opts := ServerChangeTypeOpts{
			ServerType:  &ServerType{Name: "type"},
			UpgradeDisk: true,
		}
		action, _, err := env.Client.Server.ChangeType(ctx, server, opts)
		if err != nil {
			t.Fatal(err)
		}
		if action.ID != 1 {
			t.Errorf("unexpected action ID: %d", action.ID)
		}
	})
}

func TestServerClientChangeProtection(t *testing.T) {
	var (
		ctx    = context.Background()
		server = &Server{ID: 1}
	)

	t.Run("enable delete and rebuild protection", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/servers/1/actions/change_protection", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "POST" {
				t.Error("expected POST")
			}
			var reqBody schema.ServerActionChangeProtectionRequest
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				t.Fatal(err)
			}
			if reqBody.Delete == nil || *reqBody.Delete != true {
				t.Errorf("unexpected delete: %v", reqBody.Delete)
			}
			if reqBody.Rebuild == nil || *reqBody.Rebuild != true {
				t.Errorf("unexpected rebuild: %v", reqBody.Rebuild)
			}
			json.NewEncoder(w).Encode(schema.ServerActionChangeProtectionResponse{
				Action: schema.Action{
					ID: 1,
				},
			})
		})

		opts := ServerChangeProtectionOpts{
			Delete:  Ptr(true),
			Rebuild: Ptr(true),
		}
		action, _, err := env.Client.Server.ChangeProtection(ctx, server, opts)
		if err != nil {
			t.Fatal(err)
		}

		if action.ID != 1 {
			t.Errorf("unexpected action ID: %v", action.ID)
		}
	})
}

func TestServerClientRequestConsole(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers/1/actions/request_console", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(schema.ServerActionRequestConsoleResponse{
			Action: schema.Action{
				ID: 1,
			},
			WSSURL:   "wss://console.hetzner.cloud/?server_id=1&token=3db32d15-af2f-459c-8bf8-dee1fd05f49c",
			Password: "9MQaTg2VAGI0FIpc10k3UpRXcHj2wQ6x",
		})
	})

	ctx := context.Background()
	result, _, err := env.Client.Server.RequestConsole(ctx, &Server{ID: 1})
	if err != nil {
		t.Fatal(err)
	}
	if result.Action.ID != 1 {
		t.Errorf("unexpected action ID: %d", result.Action.ID)
	}
	if result.WSSURL != "wss://console.hetzner.cloud/?server_id=1&token=3db32d15-af2f-459c-8bf8-dee1fd05f49c" {
		t.Errorf("unexpected WebSocket URL: %v", result.WSSURL)
	}
	if result.Password != "9MQaTg2VAGI0FIpc10k3UpRXcHj2wQ6x" {
		t.Errorf("unexpected password: %v", result.Password)
	}
}

func TestServerClientAttachToNetwork(t *testing.T) {
	var (
		ctx    = context.Background()
		server = &Server{ID: 1}
	)

	t.Run("attach to network", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/servers/1/actions/attach_to_network", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "POST" {
				t.Error("expected POST")
			}
			var reqBody schema.ServerActionAttachToNetworkRequest
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				t.Fatal(err)
			}
			if reqBody.Network != 1 {
				t.Errorf("unexpected Network: %v", reqBody.Network)
			}
			json.NewEncoder(w).Encode(schema.ServerActionAttachToNetworkResponse{
				Action: schema.Action{
					ID: 1,
				},
			})
		})

		opts := ServerAttachToNetworkOpts{
			Network: &Network{ID: 1},
		}
		action, _, err := env.Client.Server.AttachToNetwork(ctx, server, opts)
		if err != nil {
			t.Fatal(err)
		}

		if action.ID != 1 {
			t.Errorf("unexpected action ID: %v", action.ID)
		}
	})

	t.Run("attach to network with additional parameters", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/servers/1/actions/attach_to_network", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "POST" {
				t.Error("expected POST")
			}
			var reqBody schema.ServerActionAttachToNetworkRequest
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				t.Fatal(err)
			}
			if reqBody.Network != 1 {
				t.Errorf("unexpected Network: %v", reqBody.Network)
			}
			if reqBody.IP == nil || *reqBody.IP != "10.0.1.1" {
				t.Errorf("unexpected IP: %v", *reqBody.IP)
			}
			if len(reqBody.AliasIPs) == 0 || *reqBody.AliasIPs[0] != "10.0.1.1" {
				t.Errorf("unexpected AliasIPs: %v", *reqBody.IP)
			}
			json.NewEncoder(w).Encode(schema.ServerActionAttachToNetworkResponse{
				Action: schema.Action{
					ID: 1,
				},
			})
		})
		ip := net.ParseIP("10.0.1.1")
		aliasIPs := []net.IP{
			ip,
		}
		opts := ServerAttachToNetworkOpts{
			Network:  &Network{ID: 1},
			IP:       ip,
			AliasIPs: aliasIPs,
		}
		action, _, err := env.Client.Server.AttachToNetwork(ctx, server, opts)
		if err != nil {
			t.Fatal(err)
		}

		if action.ID != 1 {
			t.Errorf("unexpected action ID: %v", action.ID)
		}
	})
}

func TestServerClientDetachFromNetwork(t *testing.T) {
	var (
		ctx    = context.Background()
		server = &Server{ID: 1}
	)

	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers/1/actions/detach_from_network", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Error("expected POST")
		}
		var reqBody schema.ServerActionDetachFromNetworkRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatal(err)
		}
		if reqBody.Network != 1 {
			t.Errorf("unexpected Network: %v", reqBody.Network)
		}
		json.NewEncoder(w).Encode(schema.ServerActionAttachToNetworkResponse{
			Action: schema.Action{
				ID: 1,
			},
		})
	})

	opts := ServerDetachFromNetworkOpts{
		Network: &Network{ID: 1},
	}
	action, _, err := env.Client.Server.DetachFromNetwork(ctx, server, opts)
	if err != nil {
		t.Fatal(err)
	}

	if action.ID != 1 {
		t.Errorf("unexpected action ID: %v", action.ID)
	}
}

func TestServerClientChangeAliasIP(t *testing.T) {
	var (
		ctx    = context.Background()
		server = &Server{ID: 1}
	)

	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers/1/actions/change_alias_ips", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Error("expected POST")
		}
		var reqBody schema.ServerActionChangeAliasIPsRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatal(err)
		}
		if reqBody.Network != 1 {
			t.Errorf("unexpected Network: %v", reqBody.Network)
		}
		if len(reqBody.AliasIPs) == 0 || reqBody.AliasIPs[0] != "10.0.1.1" {
			t.Errorf("unexpected AliasIPs: %v", reqBody.AliasIPs[0])
		}
		json.NewEncoder(w).Encode(schema.ServerActionAttachToNetworkResponse{
			Action: schema.Action{
				ID: 1,
			},
		})
	})
	ip := net.ParseIP("10.0.1.1")
	aliasIPs := []net.IP{
		ip,
	}
	opts := ServerChangeAliasIPsOpts{
		Network:  &Network{ID: 1},
		AliasIPs: aliasIPs,
	}
	action, _, err := env.Client.Server.ChangeAliasIPs(ctx, server, opts)
	if err != nil {
		t.Fatal(err)
	}

	if action.ID != 1 {
		t.Errorf("unexpected action ID: %v", action.ID)
	}
}

func TestServerGetMetrics(t *testing.T) {
	tests := []struct {
		name        string
		server      *Server
		opts        ServerGetMetricsOpts
		respStatus  int
		respFn      func() schema.ServerGetMetricsResponse
		expected    ServerMetrics
		expectedErr string
	}{
		{
			name:   "cpu metrics",
			server: &Server{ID: 1},
			opts: ServerGetMetricsOpts{
				Types: []ServerMetricType{ServerMetricCPU},
				Start: mustParseTime(t, "2017-01-01T00:00:00Z"),
				End:   mustParseTime(t, "2017-01-01T23:00:00Z"),
			},
			respFn: func() schema.ServerGetMetricsResponse {
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
				}

				return resp
			},
			expected: ServerMetrics{
				Start: mustParseTime(t, "2017-01-01T00:00:00Z"),
				End:   mustParseTime(t, "2017-01-01T23:00:00Z"),
				TimeSeries: map[string][]ServerMetricsValue{
					"cpu": {
						{Timestamp: 1435781470.622, Value: "42"},
						{Timestamp: 1435781471.622, Value: "43"},
					},
				},
			},
		},
		{
			name:   "all metrics",
			server: &Server{ID: 2},
			opts: ServerGetMetricsOpts{
				Types: []ServerMetricType{
					ServerMetricCPU,
					ServerMetricDisk,
					ServerMetricNetwork,
				},
				Start: mustParseTime(t, "2017-01-01T00:00:00Z"),
				End:   mustParseTime(t, "2017-01-01T23:00:00Z"),
			},
			respFn: func() schema.ServerGetMetricsResponse {
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

				return resp
			},
			expected: ServerMetrics{
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
		{
			name:   "missing metrics types",
			server: &Server{ID: 3},
			opts: ServerGetMetricsOpts{
				Start: mustParseTime(t, "2017-01-01T00:00:00Z"),
				End:   mustParseTime(t, "2017-01-01T23:00:00Z"),
			},
			expectedErr: "add query params: no metric types specified",
		},
		{
			name:   "no start time",
			server: &Server{ID: 4},
			opts: ServerGetMetricsOpts{
				Types: []ServerMetricType{ServerMetricCPU},
				End:   mustParseTime(t, "2017-01-01T23:00:00Z"),
			},
			expectedErr: "add query params: no start time specified",
		},
		{
			name:   "no end time",
			server: &Server{ID: 5},
			opts: ServerGetMetricsOpts{
				Types: []ServerMetricType{ServerMetricCPU},
				Start: mustParseTime(t, "2017-01-01T00:00:00Z"),
			},
			expectedErr: "add query params: no end time specified",
		},
		{
			name:   "call to backend API fails",
			server: &Server{ID: 6},
			opts: ServerGetMetricsOpts{
				Types: []ServerMetricType{ServerMetricCPU},
				Start: mustParseTime(t, "2017-01-01T00:00:00Z"),
				End:   mustParseTime(t, "2017-01-01T23:00:00Z"),
			},
			respStatus:  http.StatusInternalServerError,
			expectedErr: "get metrics: hcloud: server responded with status code 500",
		},
		{
			name: "no server passed",
			opts: ServerGetMetricsOpts{
				Types: []ServerMetricType{ServerMetricCPU},
				Start: mustParseTime(t, "2017-01-01T00:00:00Z"),
				End:   mustParseTime(t, "2017-01-01T23:00:00Z"),
			},
			expectedErr: "illegal argument: server is nil",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			env := newTestEnv()
			defer env.Teardown()

			if tt.server != nil {
				path := fmt.Sprintf("/servers/%d/metrics", tt.server.ID)
				env.Mux.HandleFunc(path, func(rw http.ResponseWriter, r *http.Request) {
					if r.Method != "GET" {
						t.Errorf("expected GET; got %s", r.Method)
					}
					opts := serverMetricsOptsFromURL(t, r.URL)
					if !cmp.Equal(tt.opts, opts) {
						t.Errorf("unexpected opts: url: %s\n%v", r.URL.String(), cmp.Diff(tt.opts, opts))
					}

					status := tt.respStatus
					if status == 0 {
						status = http.StatusOK
					}
					rw.WriteHeader(status)

					if tt.respFn != nil {
						resp := tt.respFn()
						if err := json.NewEncoder(rw).Encode(resp); err != nil {
							t.Errorf("failed to encode response: %v", err)
						}
					}
				})
			}

			ctx := context.Background()
			actual, _, err := env.Client.Server.GetMetrics(ctx, tt.server, tt.opts)
			if tt.expectedErr != "" {
				if tt.expectedErr != err.Error() {
					t.Errorf("expected err: %v; got: %v", tt.expectedErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("failed to get server metrics: %v", err)
			}
			if !cmp.Equal(&tt.expected, actual) {
				t.Errorf("Actual metrics did not equal expected: %s", cmp.Diff(&tt.expected, actual))
			}
		})
	}
}

func serverMetricsOptsFromURL(t *testing.T, u *url.URL) ServerGetMetricsOpts {
	var opts ServerGetMetricsOpts

	for k, vs := range u.Query() {
		switch k {
		case "type":
			for _, v := range vs {
				opts.Types = append(opts.Types, ServerMetricType(v))
			}
		case "start":
			if len(vs) != 1 {
				t.Errorf("expected one value for start; got %d: %v", len(vs), vs)
				continue
			}
			v, err := time.Parse(time.RFC3339, vs[0])
			if err != nil {
				t.Errorf("parse start as RFC3339: %v", err)
			}
			opts.Start = v
		case "end":
			if len(vs) != 1 {
				t.Errorf("expected one value for end; got %d: %v", len(vs), vs)
				continue
			}
			v, err := time.Parse(time.RFC3339, vs[0])
			if err != nil {
				t.Errorf("parse end as RFC3339: %v", err)
			}
			opts.End = v
		case "step":
			if len(vs) != 1 {
				t.Errorf("expected one value for step; got %d: %v", len(vs), vs)
				continue
			}
			v, err := strconv.Atoi(vs[0])
			if err != nil {
				t.Errorf("invalid step: %v", err)
			}
			opts.Step = v
		}
	}

	return opts
}

func TestServerAddToPlacementGroup(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	const (
		serverID         = 1
		actionID         = 42
		placementGroupID = 123
	)

	env.Mux.HandleFunc(fmt.Sprintf("/servers/%d/actions/add_to_placement_group", serverID), func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Error("expected POST")
		}
		var reqBody schema.ServerActionAddToPlacementGroupRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatal(err)
		}
		if reqBody.PlacementGroup != placementGroupID {
			t.Errorf("unexpected PlacementGroup: %v", reqBody.PlacementGroup)
		}
		json.NewEncoder(w).Encode(schema.ServerActionAddToPlacementGroupResponse{
			Action: schema.Action{
				ID: actionID,
			},
		})
	})

	var (
		ctx            = context.Background()
		server         = &Server{ID: serverID}
		placementGroup = &PlacementGroup{ID: placementGroupID}
	)

	action, _, err := env.Client.Server.AddToPlacementGroup(ctx, server, placementGroup)
	if err != nil {
		t.Fatal(err)
	}
	if action.ID != actionID {
		t.Errorf("unexpected action ID: %v", action.ID)
	}
}

func TestServerRemoveFromPlacementGroup(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	const (
		serverID = 1
		actionID = 42
	)

	env.Mux.HandleFunc(fmt.Sprintf("/servers/%d/actions/remove_from_placement_group", serverID), func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Error("expected POST")
		}
		json.NewEncoder(w).Encode(schema.ServerActionRemoveFromPlacementGroupResponse{
			Action: schema.Action{
				ID: actionID,
			},
		})
	})

	var (
		ctx    = context.Background()
		server = &Server{ID: serverID}
	)

	action, _, err := env.Client.Server.RemoveFromPlacementGroup(ctx, server)
	if err != nil {
		t.Fatal(err)
	}
	if action.ID != actionID {
		t.Errorf("unexpected action ID: %v", action.ID)
	}
}
