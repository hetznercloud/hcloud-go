package hcloud

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/hetznercloud/hcloud-go/v2/hcloud/schema"
)

func TestFirewallCreateOptsValidate(t *testing.T) {
	testCases := map[string]struct {
		Opts  FirewallCreateOpts
		Valid bool
	}{
		"empty": {
			Opts:  FirewallCreateOpts{},
			Valid: false,
		},
		"all set": {
			Opts: FirewallCreateOpts{
				Name:   "name",
				Labels: map[string]string{},
			},
			Valid: true,
		},
	}
	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			err := testCase.Opts.Validate()
			if err == nil && !testCase.Valid || err != nil && testCase.Valid {
				t.FailNow()
			}
		})
	}
}

func TestFirewallClientGetByID(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/firewalls/1", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(schema.FirewallGetResponse{
			Firewall: schema.Firewall{
				ID: 1,
			},
		})
	})

	ctx := context.Background()

	firewall, _, err := env.Client.Firewall.GetByID(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
	if firewall == nil {
		t.Fatal("no firewall")
	}
	if firewall.ID != 1 {
		t.Errorf("unexpected firewall ID: %v", firewall.ID)
	}

	t.Run("called via Get", func(t *testing.T) {
		firewall, _, err := env.Client.Firewall.Get(ctx, "1")
		if err != nil {
			t.Fatal(err)
		}
		if firewall == nil {
			t.Fatal("no firewall")
		}
		if firewall.ID != 1 {
			t.Errorf("unexpected firewall ID: %v", firewall.ID)
		}
	})
}

func TestFirewallClientGetByIDNotFound(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/firewalls/1", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(schema.ErrorResponse{
			Error: schema.Error{
				Code: string(ErrorCodeNotFound),
			},
		})
	})

	ctx := context.Background()

	firewall, _, err := env.Client.Firewall.GetByID(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
	if firewall != nil {
		t.Fatal("expected no firewall")
	}
}

func TestFirewallClientGetByName(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/firewalls", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RawQuery != "name=myfirewall" {
			t.Fatal("missing name query")
		}
		json.NewEncoder(w).Encode(schema.FirewallListResponse{
			Firewalls: []schema.Firewall{
				{
					ID:   1,
					Name: "myfirewall",
				},
			},
		})
	})

	ctx := context.Background()

	firewall, _, err := env.Client.Firewall.GetByName(ctx, "myfirewall")
	if err != nil {
		t.Fatal(err)
	}
	if firewall == nil {
		t.Fatal("no firewall")
	}
	if firewall.ID != 1 {
		t.Errorf("unexpected firewall ID: %v", firewall.ID)
	}

	t.Run("via Get", func(t *testing.T) {
		firewall, _, err := env.Client.Firewall.Get(ctx, "myfirewall")
		if err != nil {
			t.Fatal(err)
		}
		if firewall == nil {
			t.Fatal("no firewall")
		}
		if firewall.ID != 1 {
			t.Errorf("unexpected firewall ID: %v", firewall.ID)
		}
	})
}

func TestFirewallClientGetByNameNotFound(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/firewalls", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RawQuery != "name=myfirewall" {
			t.Fatal("missing name query")
		}
		json.NewEncoder(w).Encode(schema.FirewallListResponse{
			Firewalls: []schema.Firewall{},
		})
	})

	ctx := context.Background()

	firewall, _, err := env.Client.Firewall.GetByName(ctx, "myfirewall")
	if err != nil {
		t.Fatal(err)
	}
	if firewall != nil {
		t.Fatal("unexpected firewall")
	}
}

func TestFirewallClientGetByNameEmpty(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	ctx := context.Background()

	firewall, _, err := env.Client.Firewall.GetByName(ctx, "")
	if err != nil {
		t.Fatal(err)
	}
	if firewall != nil {
		t.Fatal("unexpected firewall")
	}
}

func TestFirewallCreate(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/firewalls", func(w http.ResponseWriter, r *http.Request) {
		var reqBody schema.FirewallCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatal(err)
		}
		expectedReqBody := schema.FirewallCreateRequest{
			Name: "myfirewall",
			Labels: func() *map[string]string {
				labels := map[string]string{"key": "value"}
				return &labels
			}(),
			ApplyTo: []schema.FirewallResource{
				{
					Type: "server",
					Server: &schema.FirewallResourceServer{
						ID: 2,
					},
				},
			},
		}
		if !cmp.Equal(expectedReqBody, reqBody) {
			t.Log(cmp.Diff(expectedReqBody, reqBody))
			t.Error("unexpected request body")
		}
		json.NewEncoder(w).Encode(schema.FirewallCreateResponse{
			Firewall: schema.Firewall{ID: 1},
		})
	})

	ctx := context.Background()

	opts := FirewallCreateOpts{
		Name:   "myfirewall",
		Labels: map[string]string{"key": "value"},
		ApplyTo: []FirewallResource{
			{
				Type: FirewallResourceTypeServer,
				Server: &FirewallResourceServer{
					ID: 2,
				},
			},
		},
	}
	_, _, err := env.Client.Firewall.Create(ctx, opts)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFirewallCreateValidation(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	ctx := context.Background()
	opts := FirewallCreateOpts{}
	_, _, err := env.Client.Firewall.Create(ctx, opts)
	if err == nil || err.Error() != "missing name" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFirewallDelete(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/firewalls/1", func(w http.ResponseWriter, r *http.Request) {})

	var (
		ctx      = context.Background()
		firewall = &Firewall{ID: 1}
	)

	_, err := env.Client.Firewall.Delete(ctx, firewall)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFirewallClientUpdate(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/firewalls/1", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Error("expected PUT")
		}
		var reqBody schema.FirewallUpdateRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatal(err)
		}
		expectedReqBody := schema.FirewallUpdateRequest{
			Name: Ptr("test"),
			Labels: func() *map[string]string {
				labels := map[string]string{"key": "value"}
				return &labels
			}(),
		}
		if !cmp.Equal(expectedReqBody, reqBody) {
			t.Log(cmp.Diff(expectedReqBody, reqBody))
			t.Error("unexpected request body")
		}
		json.NewEncoder(w).Encode(schema.FirewallUpdateResponse{
			Firewall: schema.Firewall{
				ID: 1,
			},
		})
	})

	var (
		ctx      = context.Background()
		firewall = &Firewall{ID: 1}
	)

	opts := FirewallUpdateOpts{
		Name:   "test",
		Labels: map[string]string{"key": "value"},
	}
	updatedFirewall, _, err := env.Client.Firewall.Update(ctx, firewall, opts)
	if err != nil {
		t.Fatal(err)
	}
	if updatedFirewall.ID != 1 {
		t.Errorf("unexpected firewall ID: %v", updatedFirewall.ID)
	}
}

func TestFirewallSetRules(t *testing.T) {
	description := "allow icmp out"

	tests := []struct {
		name            string
		expectedReqBody schema.FirewallActionSetRulesRequest
		opts            FirewallSetRulesOpts
	}{
		{
			name: "direction in",
			expectedReqBody: schema.FirewallActionSetRulesRequest{
				Rules: []schema.FirewallRuleReq{
					{
						Direction: "in",
						SourceIPs: []string{"10.0.0.5/32", "10.0.0.6/32"},
						Protocol:  "icmp",
					},
				},
			},
			opts: FirewallSetRulesOpts{
				Rules: []FirewallRule{
					{
						Direction: FirewallRuleDirectionIn,
						SourceIPs: []net.IPNet{
							{
								IP:   net.ParseIP("10.0.0.5"),
								Mask: net.CIDRMask(32, 32),
							},
							{
								IP:   net.ParseIP("10.0.0.6"),
								Mask: net.CIDRMask(32, 32),
							},
						},
						Protocol: FirewallRuleProtocolICMP,
					},
				},
			},
		},
		{
			name: "direction out",
			expectedReqBody: schema.FirewallActionSetRulesRequest{
				Rules: []schema.FirewallRuleReq{
					{
						Direction:      "out",
						DestinationIPs: []string{"10.0.0.5/32", "10.0.0.6/32"},
						Protocol:       "icmp",
					},
				},
			},
			opts: FirewallSetRulesOpts{
				Rules: []FirewallRule{
					{
						Direction: FirewallRuleDirectionOut,
						DestinationIPs: []net.IPNet{
							{
								IP:   net.ParseIP("10.0.0.5"),
								Mask: net.CIDRMask(32, 32),
							},
							{
								IP:   net.ParseIP("10.0.0.6"),
								Mask: net.CIDRMask(32, 32),
							},
						},
						Protocol: FirewallRuleProtocolICMP,
					},
				},
			},
		},
		{
			name: "empty",
			expectedReqBody: schema.FirewallActionSetRulesRequest{
				Rules: []schema.FirewallRuleReq{},
			},
			opts: FirewallSetRulesOpts{
				Rules: []FirewallRule{},
			},
		},
		{
			name: "description",
			expectedReqBody: schema.FirewallActionSetRulesRequest{
				Rules: []schema.FirewallRuleReq{
					{
						Direction:      "out",
						DestinationIPs: []string{"10.0.0.5/32", "10.0.0.6/32"},
						Protocol:       "icmp",
						Description:    &description,
					},
				},
			},
			opts: FirewallSetRulesOpts{
				Rules: []FirewallRule{
					{
						Direction: FirewallRuleDirectionOut,
						DestinationIPs: []net.IPNet{
							{
								IP:   net.ParseIP("10.0.0.5"),
								Mask: net.CIDRMask(32, 32),
							},
							{
								IP:   net.ParseIP("10.0.0.6"),
								Mask: net.CIDRMask(32, 32),
							},
						},
						Protocol:    FirewallRuleProtocolICMP,
						Description: &description,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := newTestEnv()
			defer env.Teardown()

			env.Mux.HandleFunc("/firewalls/1/actions/set_rules", func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" {
					t.Error("expected POST")
				}
				var reqBody schema.FirewallActionSetRulesRequest
				if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
					t.Fatal(err)
				}
				if !cmp.Equal(tt.expectedReqBody, reqBody) {
					t.Log(cmp.Diff(tt.expectedReqBody, reqBody))
					t.Error("unexpected request body")
				}
				json.NewEncoder(w).Encode(schema.FirewallActionSetRulesResponse{
					Actions: []schema.Action{
						{
							ID: 1,
						},
					},
				})
			})

			var (
				ctx      = context.Background()
				firewall = &Firewall{ID: 1}
			)

			actions, _, err := env.Client.Firewall.SetRules(ctx, firewall, tt.opts)
			if err != nil {
				t.Fatal(err)
			}
			if len(actions) != 1 || actions[0].ID != 1 {
				t.Errorf("unexpected actions: %v", actions)
			}
		})
	}
}

func TestFirewallApplyToResources(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/firewalls/1/actions/apply_to_resources", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Error("expected POST")
		}
		var reqBody schema.FirewallActionApplyToResourcesRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatal(err)
		}
		expectedReqBody := schema.FirewallActionApplyToResourcesRequest{
			ApplyTo: []schema.FirewallResource{
				{
					Type:   "server",
					Server: &schema.FirewallResourceServer{ID: 5},
				},
				{
					Type:          "label_selector",
					LabelSelector: &schema.FirewallResourceLabelSelector{Selector: "a=b"},
				},
			},
		}
		if !cmp.Equal(expectedReqBody, reqBody) {
			t.Log(cmp.Diff(expectedReqBody, reqBody))
			t.Error("unexpected request body")
		}
		json.NewEncoder(w).Encode(schema.FirewallActionApplyToResourcesResponse{
			Actions: []schema.Action{
				{
					ID: 1,
				},
			},
		})
	})

	var (
		ctx      = context.Background()
		firewall = &Firewall{ID: 1}
	)

	resources := []FirewallResource{
		{
			Type:   FirewallResourceTypeServer,
			Server: &FirewallResourceServer{ID: 5},
		},
		{
			Type:          FirewallResourceTypeLabelSelector,
			LabelSelector: &FirewallResourceLabelSelector{Selector: "a=b"},
		},
	}

	actions, _, err := env.Client.Firewall.ApplyResources(ctx, firewall, resources)
	if err != nil {
		t.Fatal(err)
	}
	if len(actions) != 1 || actions[0].ID != 1 {
		t.Errorf("unexpected actions: %v", actions)
	}
}

func TestFirewallRemoveFromResources(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/firewalls/1/actions/remove_from_resources", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Error("expected POST")
		}
		var reqBody schema.FirewallActionRemoveFromResourcesRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatal(err)
		}
		expectedReqBody := schema.FirewallActionRemoveFromResourcesRequest{
			RemoveFrom: []schema.FirewallResource{
				{
					Type:   "server",
					Server: &schema.FirewallResourceServer{ID: 5},
				},
				{
					Type:          "label_selector",
					LabelSelector: &schema.FirewallResourceLabelSelector{Selector: "a=b"},
				},
			},
		}
		if !cmp.Equal(expectedReqBody, reqBody) {
			t.Log(cmp.Diff(expectedReqBody, reqBody))
			t.Error("unexpected request body")
		}
		json.NewEncoder(w).Encode(schema.FirewallActionRemoveFromResourcesResponse{
			Actions: []schema.Action{
				{
					ID: 1,
				},
			},
		})
	})

	var (
		ctx      = context.Background()
		firewall = &Firewall{ID: 1}
	)

	resources := []FirewallResource{
		{
			Type:   FirewallResourceTypeServer,
			Server: &FirewallResourceServer{ID: 5},
		},
		{
			Type:          FirewallResourceTypeLabelSelector,
			LabelSelector: &FirewallResourceLabelSelector{Selector: "a=b"},
		},
	}

	actions, _, err := env.Client.Firewall.RemoveResources(ctx, firewall, resources)
	if err != nil {
		t.Fatal(err)
	}
	if len(actions) != 1 || actions[0].ID != 1 {
		t.Errorf("unexpected actions: %v", actions)
	}
}
