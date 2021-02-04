package hcloud

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"testing"

	"github.com/hetznercloud/hcloud-go/hcloud/schema"
)

func TestNetworkClientGetByID(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/networks/1", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(schema.NetworkGetResponse{
			Network: schema.Network{
				ID: 1,
			},
		})
	})
	ctx := context.Background()

	network, _, err := env.Client.Network.GetByID(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
	if network == nil {
		t.Fatal("no network")
	}
	if network.ID != 1 {
		t.Errorf("unexpected network ID: %v", network.ID)
	}

	t.Run("called via Get", func(t *testing.T) {
		network, _, err := env.Client.Network.Get(ctx, "1")
		if err != nil {
			t.Fatal(err)
		}
		if network == nil {
			t.Fatal("no network")
		}
		if network.ID != 1 {
			t.Errorf("unexpected network ID: %v", network.ID)
		}
	})
}

func TestNetworkClientGetByIDNotFound(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/networks/1", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(schema.ErrorResponse{
			Error: schema.Error{
				Code: string(ErrorCodeNotFound),
			},
		})
	})

	ctx := context.Background()
	network, _, err := env.Client.Network.GetByID(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
	if network != nil {
		t.Fatal("expected no network")
	}
}

func TestNetworkClientGetByName(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/networks", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RawQuery != "name=mynet" {
			t.Fatal("missing name query")
		}
		json.NewEncoder(w).Encode(schema.NetworkListResponse{
			Networks: []schema.Network{
				{
					ID:   1,
					Name: "mynet",
				},
			},
		})
	})
	ctx := context.Background()

	network, _, err := env.Client.Network.GetByName(ctx, "mynet")
	if err != nil {
		t.Fatal(err)
	}
	if network == nil {
		t.Fatal("no network")
	}
	if network.ID != 1 {
		t.Errorf("unexpected network ID: %v", network.ID)
	}

	t.Run("via Get", func(t *testing.T) {
		network, _, err := env.Client.Network.Get(ctx, "mynet")
		if err != nil {
			t.Fatal(err)
		}
		if network == nil {
			t.Fatal("no network")
		}
		if network.ID != 1 {
			t.Errorf("unexpected network ID: %v", network.ID)
		}
	})
}

func TestNetworkClientGetByNameNotFound(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/networks", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RawQuery != "name=mynet" {
			t.Fatal("missing name query")
		}
		json.NewEncoder(w).Encode(schema.NetworkListResponse{
			Networks: []schema.Network{},
		})
	})

	ctx := context.Background()
	network, _, err := env.Client.Network.GetByName(ctx, "mynet")
	if err != nil {
		t.Fatal(err)
	}
	if network != nil {
		t.Fatal("unexpected network")
	}
}

func TestNetworkClientGetByNameEmpty(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	ctx := context.Background()
	network, _, err := env.Client.Network.GetByName(ctx, "")
	if err != nil {
		t.Fatal(err)
	}
	if network != nil {
		t.Fatal("unexpected network")
	}
}

func TestNetworkCreate(t *testing.T) {
	var (
		ctx           = context.Background()
		_, ipRange, _ = net.ParseCIDR("10.0.1.0/24")
	)

	t.Run("missing required field name", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		opts := NetworkCreateOpts{}
		_, _, err := env.Client.Network.Create(ctx, opts)
		if err == nil || err.Error() != "missing name" {
			t.Fatalf("Network.Create should fail with \"missing name\" but failed with %s", err)
		}
	})

	t.Run("missing required field ip range", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		opts := NetworkCreateOpts{
			Name: "my-network",
		}
		_, _, err := env.Client.Network.Create(ctx, opts)
		if err == nil || err.Error() != "missing IP range" {
			t.Fatalf("Network.Create should fail with \"missing IP range\" but failed with %s", err)
		}
	})

	t.Run("required fields", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/networks", func(w http.ResponseWriter, r *http.Request) {
			var reqBody schema.NetworkCreateRequest
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				t.Fatal(err)
			}
			if reqBody.Name != "my-network" {
				t.Errorf("unexpected Name: %v", reqBody.Name)
			}
			if reqBody.IPRange != "10.0.1.0/24" {
				t.Errorf("unexpected IPRange: %v", reqBody.Name)
			}
			json.NewEncoder(w).Encode(schema.NetworkCreateResponse{
				Network: schema.Network{
					ID: 1,
				},
			})
		})
		opts := NetworkCreateOpts{
			Name:    "my-network",
			IPRange: ipRange,
		}
		_, _, err := env.Client.Network.Create(ctx, opts)
		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestNetworkDelete(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/networks/1", func(w http.ResponseWriter, r *http.Request) {})

	var (
		ctx     = context.Background()
		network = &Network{ID: 1}
	)
	_, err := env.Client.Network.Delete(ctx, network)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNetworkClientUpdate(t *testing.T) {
	var (
		ctx     = context.Background()
		network = &Network{ID: 1}
	)

	t.Run("update name", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/networks/1", func(w http.ResponseWriter, r *http.Request) {
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
			json.NewEncoder(w).Encode(schema.NetworkUpdateResponse{
				Network: schema.Network{
					ID: 1,
				},
			})
		})

		opts := NetworkUpdateOpts{
			Name: "test",
		}
		updatedNetwork, _, err := env.Client.Network.Update(ctx, network, opts)
		if err != nil {
			t.Fatal(err)
		}

		if updatedNetwork.ID != 1 {
			t.Errorf("unexpected network ID: %v", updatedNetwork.ID)
		}
	})

	t.Run("update labels", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/networks/1", func(w http.ResponseWriter, r *http.Request) {
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
			json.NewEncoder(w).Encode(schema.NetworkUpdateResponse{
				Network: schema.Network{
					ID: 1,
				},
			})
		})

		opts := NetworkUpdateOpts{
			Labels: map[string]string{"key": "value"},
		}
		updatedNetwork, _, err := env.Client.Network.Update(ctx, network, opts)
		if err != nil {
			t.Fatal(err)
		}

		if updatedNetwork.ID != 1 {
			t.Errorf("unexpected network ID: %v", updatedNetwork.ID)
		}
	})

	t.Run("no updates", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/networks/1", func(w http.ResponseWriter, r *http.Request) {
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
			json.NewEncoder(w).Encode(schema.NetworkUpdateResponse{
				Network: schema.Network{
					ID: 1,
				},
			})
		})

		opts := NetworkUpdateOpts{}
		updatedNetwork, _, err := env.Client.Network.Update(ctx, network, opts)
		if err != nil {
			t.Fatal(err)
		}

		if updatedNetwork.ID != 1 {
			t.Errorf("unexpected network ID: %v", updatedNetwork.ID)
		}
	})
}

func TestNetworkClientChangeIPRange(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/networks/1/actions/change_ip_range", func(w http.ResponseWriter, r *http.Request) {
		var reqBody schema.NetworkActionChangeIPRangeRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatal(err)
		}
		if reqBody.IPRange != "10.0.1.0/24" {
			t.Errorf("unexpected type: %v", reqBody.IPRange)
		}
		json.NewEncoder(w).Encode(schema.NetworkActionChangeIPRangeResponse{
			Action: schema.Action{
				ID: 1,
			},
		})
	})

	ctx := context.Background()
	_, newIPRange, _ := net.ParseCIDR("10.0.1.0/24")
	opts := NetworkChangeIPRangeOpts{
		IPRange: newIPRange,
	}
	action, _, err := env.Client.Network.ChangeIPRange(ctx, &Network{ID: 1}, opts)
	if err != nil {
		t.Fatal(err)
	}
	if action.ID != 1 {
		t.Errorf("unexpected action ID: %d", action.ID)
	}
}

func TestNetworkClientAddSubnet(t *testing.T) {
	t.Run("type server with ip range", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/networks/1/actions/add_subnet", func(w http.ResponseWriter, r *http.Request) {
			var reqBody schema.NetworkActionAddSubnetRequest
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				t.Fatal(err)
			}
			if reqBody.Type != "cloud" {
				t.Errorf("unexpected Type: %v", reqBody.Type)
			}
			if reqBody.IPRange != "10.0.1.0/24" {
				t.Errorf("unexpected IPRange: %v", reqBody.IPRange)
			}
			if reqBody.NetworkZone != "eu-central" {
				t.Errorf("unexpected NetworkZone: %v", reqBody.NetworkZone)
			}
			json.NewEncoder(w).Encode(schema.NetworkActionAddSubnetResponse{
				Action: schema.Action{
					ID: 1,
				},
			})
		})

		ctx := context.Background()
		_, ipRange, _ := net.ParseCIDR("10.0.1.0/24")
		opts := NetworkAddSubnetOpts{
			Subnet: NetworkSubnet{
				Type:        NetworkSubnetTypeCloud,
				IPRange:     ipRange,
				NetworkZone: NetworkZoneEUCentral,
			},
		}
		action, _, err := env.Client.Network.AddSubnet(ctx, &Network{ID: 1}, opts)
		if err != nil {
			t.Fatal(err)
		}
		if action.ID != 1 {
			t.Errorf("unexpected action ID: %d", action.ID)
		}
	})

	t.Run("type server without ip range", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/networks/1/actions/add_subnet", func(w http.ResponseWriter, r *http.Request) {
			var reqBody schema.NetworkActionAddSubnetRequest
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				t.Fatal(err)
			}
			if reqBody.Type != "cloud" {
				t.Errorf("unexpected Type: %v", reqBody.Type)
			}
			if reqBody.IPRange != "" {
				t.Errorf("unexpected IPRange: %v", reqBody.IPRange)
			}
			if reqBody.NetworkZone != "eu-central" {
				t.Errorf("unexpected NetworkZone: %v", reqBody.NetworkZone)
			}
			json.NewEncoder(w).Encode(schema.NetworkActionAddSubnetResponse{
				Action: schema.Action{
					ID: 1,
				},
			})
		})

		ctx := context.Background()
		opts := NetworkAddSubnetOpts{
			Subnet: NetworkSubnet{
				Type:        NetworkSubnetTypeCloud,
				NetworkZone: NetworkZoneEUCentral,
			},
		}
		action, _, err := env.Client.Network.AddSubnet(ctx, &Network{ID: 1}, opts)
		if err != nil {
			t.Fatal(err)
		}
		if action.ID != 1 {
			t.Errorf("unexpected action ID: %d", action.ID)
		}
	})

	t.Run("type vswitch with ip range", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/networks/1/actions/add_subnet", func(w http.ResponseWriter, r *http.Request) {
			var reqBody schema.NetworkActionAddSubnetRequest
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				t.Fatal(err)
			}
			if reqBody.Type != "vswitch" {
				t.Errorf("unexpected Type: %v", reqBody.Type)
			}
			if reqBody.IPRange != "10.0.1.0/24" {
				t.Errorf("unexpected IPRange: %v", reqBody.IPRange)
			}
			if reqBody.NetworkZone != "eu-central" {
				t.Errorf("unexpected NetworkZone: %v", reqBody.NetworkZone)
			}
			if reqBody.VSwitchID != 123 {
				t.Errorf("unexpected VSwitchID: %v", reqBody.VSwitchID)
			}
			json.NewEncoder(w).Encode(schema.NetworkActionAddSubnetResponse{
				Action: schema.Action{
					ID: 1,
				},
			})
		})

		ctx := context.Background()
		_, ipRange, _ := net.ParseCIDR("10.0.1.0/24")
		opts := NetworkAddSubnetOpts{
			Subnet: NetworkSubnet{
				Type:        NetworkSubnetTypeVSwitch,
				IPRange:     ipRange,
				NetworkZone: NetworkZoneEUCentral,
				VSwitchID:   123,
			},
		}
		action, _, err := env.Client.Network.AddSubnet(ctx, &Network{ID: 1}, opts)
		if err != nil {
			t.Fatal(err)
		}
		if action.ID != 1 {
			t.Errorf("unexpected action ID: %d", action.ID)
		}
	})
}

func TestNetworkClientDeleteSubnet(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/networks/1/actions/delete_subnet", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(schema.NetworkActionDeleteSubnetResponse{
			Action: schema.Action{
				ID: 1,
			},
		})
	})

	ctx := context.Background()
	_, ipRange, _ := net.ParseCIDR("10.0.1.0/24")
	opts := NetworkDeleteSubnetOpts{
		Subnet: NetworkSubnet{
			IPRange: ipRange,
		},
	}
	action, _, err := env.Client.Network.DeleteSubnet(ctx, &Network{ID: 1}, opts)
	if err != nil {
		t.Fatal(err)
	}
	if action.ID != 1 {
		t.Errorf("unexpected action ID: %d", action.ID)
	}
}

func TestNetworkClientAddRoute(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/networks/1/actions/add_route", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(schema.NetworkActionAddRouteResponse{
			Action: schema.Action{
				ID: 1,
			},
		})
	})

	ctx := context.Background()
	_, destination, _ := net.ParseCIDR("10.0.1.0/24")
	opts := NetworkAddRouteOpts{
		Route: NetworkRoute{
			Destination: destination,
			Gateway:     net.ParseIP("10.0.1.1"),
		},
	}
	action, _, err := env.Client.Network.AddRoute(ctx, &Network{ID: 1}, opts)
	if err != nil {
		t.Fatal(err)
	}
	if action.ID != 1 {
		t.Errorf("unexpected action ID: %d", action.ID)
	}
}

func TestNetworkClientDeleteRoute(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/networks/1/actions/delete_route", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(schema.NetworkActionDeleteRouteResponse{
			Action: schema.Action{
				ID: 1,
			},
		})
	})

	ctx := context.Background()
	_, destination, _ := net.ParseCIDR("10.0.1.0/24")
	opts := NetworkDeleteRouteOpts{
		Route: NetworkRoute{
			Destination: destination,
			Gateway:     net.ParseIP("10.0.1.1"),
		},
	}
	action, _, err := env.Client.Network.DeleteRoute(ctx, &Network{ID: 1}, opts)
	if err != nil {
		t.Fatal(err)
	}
	if action.ID != 1 {
		t.Errorf("unexpected action ID: %d", action.ID)
	}
}

func TestNetworkClientChangeProtection(t *testing.T) {
	var (
		ctx     = context.Background()
		network = &Network{ID: 1}
	)

	t.Run("enable delete protection", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/networks/1/actions/change_protection", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "POST" {
				t.Error("expected POST")
			}
			var reqBody schema.NetworkActionChangeProtectionRequest
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				t.Fatal(err)
			}
			if reqBody.Delete == nil || *reqBody.Delete != true {
				t.Errorf("unexpected delete: %v", reqBody.Delete)
			}
			json.NewEncoder(w).Encode(schema.NetworkActionChangeProtectionResponse{
				Action: schema.Action{
					ID: 1,
				},
			})
		})

		opts := NetworkChangeProtectionOpts{
			Delete: Bool(true),
		}
		action, _, err := env.Client.Network.ChangeProtection(ctx, network, opts)
		if err != nil {
			t.Fatal(err)
		}

		if action.ID != 1 {
			t.Errorf("unexpected action ID: %v", action.ID)
		}
	})
}
