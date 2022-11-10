package hcloud

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hetznercloud/hcloud-go/hcloud/schema"
	"github.com/stretchr/testify/assert"
)

func TestPrimaryIPClient(t *testing.T) {
	t.Run("GetByID", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/primary_ips/1", func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(schema.PrimaryIPGetResult{
				PrimaryIP: schema.PrimaryIP{
					ID: 1,
				},
			})
		})

		ctx := context.Background()
		primaryIP, _, err := env.Client.PrimaryIP.GetByID(ctx, 1)
		if err != nil {
			t.Fatal(err)
		}
		if primaryIP == nil {
			t.Fatal("no primary_ip")
		}
		if primaryIP.ID != 1 {
			t.Errorf("unexpected primary_ip ID: %v", primaryIP.ID)
		}

		t.Run("via Get", func(t *testing.T) {
			primaryIP, _, err := env.Client.PrimaryIP.Get(ctx, "1")
			if err != nil {
				t.Fatal(err)
			}
			if primaryIP == nil {
				t.Fatal("no primary_ip")
			}
			if primaryIP.ID != 1 {
				t.Errorf("unexpected primary_ip ID: %v", primaryIP.ID)
			}
		})
	})

	t.Run("GetByID (not found)", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/primary_ips/1", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(schema.ErrorResponse{
				Error: schema.Error{
					Code: string(ErrorCodeNotFound),
				},
			})
		})

		ctx := context.Background()
		primaryIP, _, err := env.Client.PrimaryIP.GetByID(ctx, 1)
		if err != nil {
			t.Fatal(err)
		}
		if primaryIP != nil {
			t.Fatal("expected no primary_ip")
		}
	})

	t.Run("GetByName", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/primary_ips", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.RawQuery != "name=fsn1-dc8" {
				t.Fatal("missing name query")
			}
			json.NewEncoder(w).Encode(schema.PrimaryIPListResult{
				PrimaryIPs: []schema.PrimaryIP{
					{
						ID: 1,
					},
				},
			})
		})

		ctx := context.Background()
		primaryIP, _, err := env.Client.PrimaryIP.GetByName(ctx, "fsn1-dc8")
		if err != nil {
			t.Fatal(err)
		}
		if primaryIP == nil {
			t.Fatal("no primary_ip")
		}
		if primaryIP.ID != 1 {
			t.Errorf("unexpected primary_ip ID: %v", primaryIP.ID)
		}

		t.Run("via Get", func(t *testing.T) {
			primaryIP, _, err := env.Client.PrimaryIP.Get(ctx, "fsn1-dc8")
			if err != nil {
				t.Fatal(err)
			}
			if primaryIP == nil {
				t.Fatal("no primary_ip")
			}
			if primaryIP.ID != 1 {
				t.Errorf("unexpected primary_ip ID: %v", primaryIP.ID)
			}
		})
	})

	t.Run("GetByName (not found)", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/primary_ips", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.RawQuery != "name=fsn1-dc8" {
				t.Fatal("missing name query")
			}
			json.NewEncoder(w).Encode(schema.PrimaryIPListResult{
				PrimaryIPs: []schema.PrimaryIP{},
			})
		})

		ctx := context.Background()
		primaryIP, _, err := env.Client.PrimaryIP.GetByName(ctx, "fsn1-dc8")
		if err != nil {
			t.Fatal(err)
		}
		if primaryIP != nil {
			t.Fatal("unexpected primary_ip")
		}
	})

	t.Run("GetByName (empty)", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		ctx := context.Background()
		primaryIP, _, err := env.Client.PrimaryIP.GetByName(ctx, "")
		if err != nil {
			t.Fatal(err)
		}
		if primaryIP != nil {
			t.Fatal("unexpected primary_ip")
		}
	})

	t.Run("GetByIP", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/primary_ips", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.RawQuery != "ip=127.0.0.1" {
				t.Fatal("missing name query")
			}
			json.NewEncoder(w).Encode(schema.PrimaryIPListResult{
				PrimaryIPs: []schema.PrimaryIP{
					{
						ID: 1,
					},
				},
			})
		})

		ctx := context.Background()
		primaryIP, _, err := env.Client.PrimaryIP.GetByIP(ctx, "127.0.0.1")
		if err != nil {
			t.Fatal(err)
		}
		if primaryIP == nil {
			t.Fatal("no primary_ip")
		}
		if primaryIP.ID != 1 {
			t.Errorf("unexpected primary_ip ID: %v", primaryIP.ID)
		}
	})

	t.Run("List", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/primary_ips", func(w http.ResponseWriter, r *http.Request) {
			if page := r.URL.Query().Get("page"); page != "2" {
				t.Errorf("expected page 2; got %q", page)
			}
			if perPage := r.URL.Query().Get("per_page"); perPage != "50" {
				t.Errorf("expected per_page 50; got %q", perPage)
			}
			if name := r.URL.Query().Get("name"); name != "nbg1-dc3" {
				t.Errorf("expected name nbg1-dc3; got %q", name)
			}
			json.NewEncoder(w).Encode(schema.PrimaryIPListResult{
				PrimaryIPs: []schema.PrimaryIP{
					{ID: 1},
					{ID: 2},
				},
			})
		})

		opts := PrimaryIPListOpts{}
		opts.Page = 2
		opts.PerPage = 50
		opts.Name = "nbg1-dc3"

		ctx := context.Background()
		primaryIPs, _, err := env.Client.PrimaryIP.List(ctx, opts)
		if err != nil {
			t.Fatal(err)
		}
		if len(primaryIPs) != 2 {
			t.Fatal("expected 2 primary_ips")
		}
	})

	t.Run("All", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/primary_ips", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(struct {
				PrimaryIPs []PrimaryIP `json:"primary_ips"`
				Meta       schema.Meta `json:"meta"`
			}{
				PrimaryIPs: []PrimaryIP{
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
		primaryIPs, err := env.Client.PrimaryIP.All(ctx)
		if err != nil {
			t.Fatalf("PrimaryIP.List failed: %s", err)
		}
		if len(primaryIPs) != 3 {
			t.Fatalf("expected 3 primary_ips; got %d", len(primaryIPs))
		}
		if primaryIPs[0].ID != 1 || primaryIPs[1].ID != 2 || primaryIPs[2].ID != 3 {
			t.Errorf("unexpected primary_ips")
		}
	})
	t.Run("Create", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/primary_ips", func(w http.ResponseWriter, r *http.Request) {
			var reqBody PrimaryIPCreateOpts
			if r.Method != "POST" {
				t.Error("expected POST")
			}
			w.Header().Set("Content-Type", "application/json")
			expectedReqBody := PrimaryIPCreateOpts{
				Name:         "my-primary-ip",
				Type:         PrimaryIPTypeIPv4,
				AssigneeType: "server",
				Datacenter:   "fsn-dc14",
				Labels: func() map[string]string {
					labels := map[string]string{"key": "value"}
					return labels
				}(),
			}
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				t.Fatal(err)
			}
			if !cmp.Equal(expectedReqBody, reqBody) {
				t.Log(cmp.Diff(expectedReqBody, reqBody))
				t.Error("unexpected request body")
			}
			json.NewEncoder(w).Encode(PrimaryIPCreateResult{
				PrimaryIP: &PrimaryIP{ID: 1},
				Action:    &Action{ID: 14},
			})
		})

		ctx := context.Background()
		opts := PrimaryIPCreateOpts{
			Name:         "my-primary-ip",
			Type:         PrimaryIPTypeIPv4,
			AssigneeType: "server",
			Labels:       map[string]string{"key": "value"},
			Datacenter:   "fsn-dc14",
		}

		result, resp, err := env.Client.PrimaryIP.Create(ctx, opts)
		assert.NoError(t, err)
		assert.NotNil(t, resp, "no response returned")
		assert.NotNil(t, result.PrimaryIP, "no primary IP returned")
		assert.NotNil(t, result.Action, "no action returned")
	})
	t.Run("Update", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/primary_ips/1", func(w http.ResponseWriter, r *http.Request) {
			var reqBody PrimaryIPUpdateOpts
			if r.Method != "PUT" {
				t.Error("expected PUT")
			}
			w.Header().Set("Content-Type", "application/json")
			autoDelete := true
			expectedReqBody := PrimaryIPUpdateOpts{
				Name:       "my-primary-ip",
				AutoDelete: &autoDelete,
				Labels: func() *map[string]string {
					labels := map[string]string{"key": "value"}
					return &labels
				}(),
			}
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				t.Fatal(err)
			}
			if !cmp.Equal(expectedReqBody, reqBody) {
				t.Log(cmp.Diff(expectedReqBody, reqBody))
				t.Error("unexpected request body")
			}
			json.NewEncoder(w).Encode(schema.PrimaryIPUpdateResult{
				PrimaryIP: schema.PrimaryIP{ID: 1, IP: "2001:db8::/64"},
			})
		})

		ctx := context.Background()
		labels := map[string]string{"key": "value"}
		autoDelete := true
		opts := PrimaryIPUpdateOpts{
			Name:       "my-primary-ip",
			AutoDelete: &autoDelete,
			Labels:     &labels,
		}

		primaryIP := PrimaryIP{ID: 1, IP: net.ParseIP("2001:db8::")}
		result, resp, err := env.Client.PrimaryIP.Update(ctx, &primaryIP, opts)
		assert.NoError(t, err)
		assert.NotNil(t, resp, "no response returned")
		if result.ID != 1 {
			t.Errorf("unexpected Primary IP ID: %v", result.ID)
		}
		assert.Equal(t, primaryIP.IP, result.IP, "parsed the wrong IP")
	})
	t.Run("Assign", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/primary_ips/1/actions/assign", func(w http.ResponseWriter, r *http.Request) {
			var reqBody PrimaryIPAssignOpts
			if r.Method != "POST" {
				t.Error("expected POST")
			}
			w.Header().Set("Content-Type", "application/json")
			expectedReqBody := PrimaryIPAssignOpts{
				AssigneeType: "server",
				AssigneeID:   1,
				ID:           1,
			}
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				t.Fatal(err)
			}
			if !cmp.Equal(expectedReqBody, reqBody) {
				t.Log(cmp.Diff(expectedReqBody, reqBody))
				t.Error("unexpected request body")
			}
			json.NewEncoder(w).Encode(PrimaryIPAssignResult{
				Action: schema.Action{ID: 1},
			})
		})

		ctx := context.Background()
		opts := PrimaryIPAssignOpts{
			AssigneeType: "server",
			AssigneeID:   1,
			ID:           1,
		}

		action, resp, err := env.Client.PrimaryIP.Assign(ctx, opts)
		assert.NoError(t, err)
		assert.NotNil(t, resp, "no response returned")
		assert.NotNil(t, action, "no action returned")
	})
	t.Run("Unassign", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/primary_ips/1/actions/unassign", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "POST" {
				t.Error("expected POST")
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(PrimaryIPAssignResult{
				Action: schema.Action{ID: 1},
			})
		})

		ctx := context.Background()

		action, resp, err := env.Client.PrimaryIP.Unassign(ctx, 1)
		assert.NoError(t, err)
		assert.NotNil(t, resp, "no response returned")
		assert.NotNil(t, action, "no action returned")
	})
	t.Run("ChangeDNSPtr", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/primary_ips/1/actions/change_dns_ptr", func(w http.ResponseWriter, r *http.Request) {
			var reqBody PrimaryIPChangeDNSPtrOpts
			if r.Method != "POST" {
				t.Error("expected POST")
			}
			w.Header().Set("Content-Type", "application/json")
			expectedReqBody := PrimaryIPChangeDNSPtrOpts{
				ID: 1,
			}
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				t.Fatal(err)
			}
			if !cmp.Equal(expectedReqBody, reqBody) {
				t.Log(cmp.Diff(expectedReqBody, reqBody))
				t.Error("unexpected request body")
			}
			json.NewEncoder(w).Encode(PrimaryIPChangeDNSPtrResult{
				Action: schema.Action{ID: 1},
			})
		})

		ctx := context.Background()
		opts := PrimaryIPChangeDNSPtrOpts{
			ID: 1,
		}

		action, resp, err := env.Client.PrimaryIP.ChangeDNSPtr(ctx, opts)
		assert.NoError(t, err)
		assert.NotNil(t, resp, "no response returned")
		assert.NotNil(t, action, "no action returned")
	})
	t.Run("ChangeProtection", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/primary_ips/1/actions/change_protection", func(w http.ResponseWriter, r *http.Request) {
			var reqBody PrimaryIPChangeProtectionOpts
			if r.Method != "POST" {
				t.Error("expected POST")
			}
			w.Header().Set("Content-Type", "application/json")
			expectedReqBody := PrimaryIPChangeProtectionOpts{
				ID:     1,
				Delete: true,
			}
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				t.Fatal(err)
			}
			if !cmp.Equal(expectedReqBody, reqBody) {
				t.Log(cmp.Diff(expectedReqBody, reqBody))
				t.Error("unexpected request body")
			}
			json.NewEncoder(w).Encode(PrimaryIPChangeProtectionResult{
				Action: schema.Action{ID: 1},
			})
		})

		ctx := context.Background()
		opts := PrimaryIPChangeProtectionOpts{
			ID:     1,
			Delete: true,
		}

		action, resp, err := env.Client.PrimaryIP.ChangeProtection(ctx, opts)
		assert.NoError(t, err)
		assert.NotNil(t, resp, "no response returned")
		assert.NotNil(t, action, "no action returned")
	})
}
