package hcloud

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/hetznercloud/hcloud-go/hcloud/schema"
)

func TestLoadBalancerClientGetByID(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/load_balancers/1", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(schema.LoadBalancerGetResponse{
			LoadBalancer: schema.LoadBalancer{
				ID: 1,
			},
		})
	})
	ctx := context.Background()

	loadBalancer, _, err := env.Client.LoadBalancer.GetByID(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
	if loadBalancer == nil {
		t.Fatal("no loadBalancer")
	}
	if loadBalancer.ID != 1 {
		t.Errorf("unexpected loadBalancer ID: %v", loadBalancer.ID)
	}

	t.Run("called via Get", func(t *testing.T) {
		loadBalancer, _, err := env.Client.LoadBalancer.Get(ctx, "1")
		if err != nil {
			t.Fatal(err)
		}
		if loadBalancer == nil {
			t.Fatal("no loadBalancer")
		}
		if loadBalancer.ID != 1 {
			t.Errorf("unexpected loadBalancer ID: %v", loadBalancer.ID)
		}
	})
}

func TestLoadBalancerClientGetByIDNotFound(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/load_balancers/1", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(schema.ErrorResponse{
			Error: schema.Error{
				Code: string(ErrorCodeNotFound),
			},
		})
	})

	ctx := context.Background()
	loadBalancer, _, err := env.Client.LoadBalancer.GetByID(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
	if loadBalancer != nil {
		t.Fatal("expected no loadBalancer")
	}
}

func TestLoadBalancerClientGetByName(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/load_balancers", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RawQuery != "name=mylb" {
			t.Fatal("missing name query")
		}
		json.NewEncoder(w).Encode(schema.LoadBalancerListResponse{
			LoadBalancers: []schema.LoadBalancer{
				{
					ID:   1,
					Name: "mylb",
				},
			},
		})
	})
	ctx := context.Background()

	loadBalancer, _, err := env.Client.LoadBalancer.GetByName(ctx, "mylb")
	if err != nil {
		t.Fatal(err)
	}
	if loadBalancer == nil {
		t.Fatal("no loadBalancer")
	}
	if loadBalancer.ID != 1 {
		t.Errorf("unexpected loadBalancer ID: %v", loadBalancer.ID)
	}

	t.Run("via Get", func(t *testing.T) {
		loadBalancer, _, err := env.Client.LoadBalancer.Get(ctx, "mylb")
		if err != nil {
			t.Fatal(err)
		}
		if loadBalancer == nil {
			t.Fatal("no loadBalancer")
		}
		if loadBalancer.ID != 1 {
			t.Errorf("unexpected loadBalancer ID: %v", loadBalancer.ID)
		}
	})
}

func TestLoadBalancerClientGetByNameNotFound(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/load_balancers", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RawQuery != "name=mylb" {
			t.Fatal("missing name query")
		}
		json.NewEncoder(w).Encode(schema.LoadBalancerListResponse{
			LoadBalancers: []schema.LoadBalancer{},
		})
	})

	ctx := context.Background()
	loadBalancer, _, err := env.Client.LoadBalancer.GetByName(ctx, "mylb")
	if err != nil {
		t.Fatal(err)
	}
	if loadBalancer != nil {
		t.Fatal("unexpected loadBalancer")
	}
}

func TestLoadBalancerCreate(t *testing.T) {
	var (
		ctx = context.Background()
	)

	t.Run("missing required field name", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		opts := LoadBalancerCreateOpts{}
		_, _, err := env.Client.LoadBalancer.Create(ctx, opts)
		if err == nil || err.Error() != "missing name" {
			t.Fatalf("LoadBalancer.Create should fail with \"missing name\" but failed with %s", err)
		}
	})

	t.Run("missing required field load balancer type", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		opts := LoadBalancerCreateOpts{
			Name: "my-loadBalancer",
		}
		_, _, err := env.Client.LoadBalancer.Create(ctx, opts)
		if err == nil || err.Error() != "missing load balancer type" {
			t.Fatalf("LoadBalancer.Create should fail with \"missing load balancer type\" but failed with %s", err)
		}
	})

	t.Run("missing required field location and network zone", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		opts := LoadBalancerCreateOpts{
			Name:             "my-loadBalancer",
			LoadBalancerType: &LoadBalancerType{Name: "lb1"},
		}
		_, _, err := env.Client.LoadBalancer.Create(ctx, opts)
		if err == nil || err.Error() != "one of location and network_zone must be set" {
			t.Fatalf("LoadBalancer.Create should fail with \"one of location and network_zone must be set\" but failed with %s", err)
		}
	})

	t.Run("location and network_zone are mutually exclusive", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		opts := LoadBalancerCreateOpts{
			Name:             "my-loadBalancer",
			LoadBalancerType: &LoadBalancerType{Name: "lb1"},
			Location:         &Location{Name: "fsn1"},
			NetworkZone:      NetworkZoneEUCentral,
		}
		_, _, err := env.Client.LoadBalancer.Create(ctx, opts)
		if err == nil || err.Error() != "location and network_zone are mutually exclusive" {
			t.Fatalf("LoadBalancer.Create should fail with \"location and network_zone are mutually exclusive\" but failed with %s", err)
		}
	})

	t.Run("required fields", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/load_balancers", func(w http.ResponseWriter, r *http.Request) {
			var reqBody schema.LoadBalancerCreateRequest
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				t.Fatal(err)
			}
			if reqBody.Name != "my-load_balancer" {
				t.Errorf("unexpected Name: %v", reqBody.Name)
			}
			if reqBody.LoadBalancerType != "lb1" {
				t.Errorf("unexpected LoadBalancerType: %v", reqBody.LoadBalancerType)
			}
			if reqBody.Location != "fsn1" {
				t.Errorf("unexpected Location: %v", reqBody.Location)
			}
			if reqBody.Algorithm.Type != "round_robin" {
				t.Errorf("unexpected AlgorithmType: %v", reqBody.Algorithm.Type)
			}
			json.NewEncoder(w).Encode(schema.LoadBalancerCreateResponse{
				LoadBalancer: schema.LoadBalancer{
					ID: 1,
				},
				Action: schema.Action{
					ID: 1,
				},
			})
		})
		opts := LoadBalancerCreateOpts{
			Name:             "my-load_balancer",
			LoadBalancerType: &LoadBalancerType{Name: "lb1"},
			Algorithm:        &LoadBalancerAlgorithm{Type: "round_robin"},
			Location:         &Location{Name: "fsn1"},
		}
		_, _, err := env.Client.LoadBalancer.Create(ctx, opts)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("with network_zone", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/load_balancers", func(w http.ResponseWriter, r *http.Request) {
			var reqBody schema.LoadBalancerCreateRequest
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				t.Fatal(err)
			}
			if reqBody.Name != "my-load_balancer" {
				t.Errorf("unexpected Name: %v", reqBody.Name)
			}
			if reqBody.LoadBalancerType != "lb1" {
				t.Errorf("unexpected LoadBalancerType: %v", reqBody.LoadBalancerType)
			}
			if reqBody.NetworkZone != "eu-central" {
				t.Errorf("unexpected NetworkZone: %v", reqBody.NetworkZone)
			}
			if reqBody.Algorithm.Type != "round_robin" {
				t.Errorf("unexpected AlgorithmType: %v", reqBody.Algorithm.Type)
			}
			json.NewEncoder(w).Encode(schema.LoadBalancerCreateResponse{
				LoadBalancer: schema.LoadBalancer{
					ID: 1,
				},
				Action: schema.Action{
					ID: 1,
				},
			})
		})
		opts := LoadBalancerCreateOpts{
			Name:             "my-load_balancer",
			LoadBalancerType: &LoadBalancerType{Name: "lb1"},
			Algorithm:        &LoadBalancerAlgorithm{Type: "round_robin"},
			NetworkZone:      NetworkZoneEUCentral,
		}
		_, _, err := env.Client.LoadBalancer.Create(ctx, opts)
		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestLoadBalancerDelete(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/load_balancers/1", func(w http.ResponseWriter, r *http.Request) {
		return
	})

	var (
		ctx          = context.Background()
		loadBalancer = &LoadBalancer{ID: 1}
	)
	_, err := env.Client.LoadBalancer.Delete(ctx, loadBalancer)
	if err != nil {
		t.Fatal(err)
	}
}

func TestLoadBalancerClientUpdate(t *testing.T) {
	var (
		ctx          = context.Background()
		loadBalancer = &LoadBalancer{ID: 1}
	)

	t.Run("update name", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/load_balancers/1", func(w http.ResponseWriter, r *http.Request) {
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
			json.NewEncoder(w).Encode(schema.LoadBalancerUpdateResponse{
				LoadBalancer: schema.LoadBalancer{
					ID: 1,
				},
			})
		})

		opts := LoadBalancerUpdateOpts{
			Name: "test",
		}
		updatedLoadBalancer, _, err := env.Client.LoadBalancer.Update(ctx, loadBalancer, opts)
		if err != nil {
			t.Fatal(err)
		}

		if updatedLoadBalancer.ID != 1 {
			t.Errorf("unexpected loadBalancer ID: %v", updatedLoadBalancer.ID)
		}
	})

	t.Run("update labels", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/load_balancers/1", func(w http.ResponseWriter, r *http.Request) {
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
			json.NewEncoder(w).Encode(schema.LoadBalancerUpdateResponse{
				LoadBalancer: schema.LoadBalancer{
					ID: 1,
				},
			})
		})

		opts := LoadBalancerUpdateOpts{
			Labels: map[string]string{"key": "value"},
		}
		updatedLoadBalancer, _, err := env.Client.LoadBalancer.Update(ctx, loadBalancer, opts)
		if err != nil {
			t.Fatal(err)
		}

		if updatedLoadBalancer.ID != 1 {
			t.Errorf("unexpected loadBalancer ID: %v", updatedLoadBalancer.ID)
		}
	})

	t.Run("no updates", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/load_balancers/1", func(w http.ResponseWriter, r *http.Request) {
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
			json.NewEncoder(w).Encode(schema.LoadBalancerUpdateResponse{
				LoadBalancer: schema.LoadBalancer{
					ID: 1,
				},
			})
		})

		opts := LoadBalancerUpdateOpts{}
		updatedLoadBalancer, _, err := env.Client.LoadBalancer.Update(ctx, loadBalancer, opts)
		if err != nil {
			t.Fatal(err)
		}

		if updatedLoadBalancer.ID != 1 {
			t.Errorf("unexpected loadBalancer ID: %v", updatedLoadBalancer.ID)
		}
	})
}

func TestLoadBalancerClientChangeProtection(t *testing.T) {
	var (
		ctx          = context.Background()
		loadBalancer = &LoadBalancer{ID: 1}
	)

	t.Run("enable delete protection", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/load_balancers/1/actions/change_protection", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "POST" {
				t.Error("expected POST")
			}
			var reqBody schema.LoadBalancerActionChangeProtectionRequest
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				t.Fatal(err)
			}
			if reqBody.Delete == nil || *reqBody.Delete != true {
				t.Errorf("unexpected delete: %v", reqBody.Delete)
			}
			json.NewEncoder(w).Encode(schema.LoadBalancerActionChangeProtectionResponse{
				Action: schema.Action{
					ID: 1,
				},
			})
		})

		opts := LoadBalancerChangeProtectionOpts{
			Delete: Bool(true),
		}
		action, _, err := env.Client.LoadBalancer.ChangeProtection(ctx, loadBalancer, opts)
		if err != nil {
			t.Fatal(err)
		}

		if action.ID != 1 {
			t.Errorf("unexpected action ID: %v", action.ID)
		}
	})
}

func TestLoadBalancerClientAddServerTarget(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/load_balancers/1/actions/add_target", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Error("expected POST")
		}
		var reqBody schema.LoadBalancerActionTargetRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatal(err)
		}
		if reqBody.Type != string(LoadBalancerTargetTypeServer) {
			t.Errorf("unexpected type %v", reqBody.Type)
		}
		if reqBody.Server.ID != 1 {
			t.Errorf("unexpected server id %v", reqBody.Server.ID)
		}
		json.NewEncoder(w).Encode(schema.LoadBalancerActionTargetResponse{
			Action: schema.Action{
				ID: 1,
			},
		})
	})

	ctx := context.Background()
	action, _, err := env.Client.LoadBalancer.AddServerTarget(ctx, &LoadBalancer{ID: 1}, &Server{ID: 1})
	if err != nil {
		t.Fatal(err)
	}
	if action.ID != 1 {
		t.Errorf("unexpected action ID: %d", action.ID)
	}
}
func TestLoadBalancerClientAddLabelSelectorTarget(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/load_balancers/1/actions/add_target", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Error("expected POST")
		}
		var reqBody schema.LoadBalancerActionTargetRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatal(err)
		}
		if reqBody.Type != string(LoadBalancerTargetTypeLabelSelector) {
			t.Errorf("unexpected type %v", reqBody.Type)
		}
		if reqBody.LabelSelector.Selector != "key=value" {
			t.Errorf("unexpected LabelSelector %v", reqBody.LabelSelector)
		}
		json.NewEncoder(w).Encode(schema.LoadBalancerActionTargetResponse{
			Action: schema.Action{
				ID: 1,
			},
		})
	})

	ctx := context.Background()
	action, _, err := env.Client.LoadBalancer.AddLabelSelectorTarget(ctx, &LoadBalancer{ID: 1}, "key=value")
	if err != nil {
		t.Fatal(err)
	}
	if action.ID != 1 {
		t.Errorf("unexpected action ID: %d", action.ID)
	}

}

func TestLoadBalancerClientRemoveServerTarget(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/load_balancers/1/actions/remove_target", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Error("expected POST")
		}
		var reqBody schema.LoadBalancerActionTargetRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatal(err)
		}
		if reqBody.Type != string(LoadBalancerTargetTypeServer) {
			t.Errorf("unexpected type %v", reqBody.Type)
		}
		if reqBody.Server.ID != 1 {
			t.Errorf("unexpected server id %v", reqBody.Server.ID)
		}
		json.NewEncoder(w).Encode(schema.LoadBalancerActionTargetResponse{
			Action: schema.Action{
				ID: 1,
			},
		})
	})

	ctx := context.Background()
	action, _, err := env.Client.LoadBalancer.RemoveServerTarget(ctx, &LoadBalancer{ID: 1}, &Server{ID: 1})
	if err != nil {
		t.Fatal(err)
	}
	if action.ID != 1 {
		t.Errorf("unexpected action ID: %d", action.ID)
	}
}
func TestLoadBalancerClientRemoveLabelSelectorTarget(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/load_balancers/1/actions/remove_target", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Error("expected POST")
		}
		var reqBody schema.LoadBalancerActionTargetRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatal(err)
		}
		if reqBody.Type != string(LoadBalancerTargetTypeLabelSelector) {
			t.Errorf("unexpected type %v", reqBody.Type)
		}
		if reqBody.LabelSelector.Selector != "key=value" {
			t.Errorf("unexpected LabelSelector %v", reqBody.LabelSelector)
		}
		json.NewEncoder(w).Encode(schema.LoadBalancerActionTargetResponse{
			Action: schema.Action{
				ID: 1,
			},
		})
	})

	ctx := context.Background()
	action, _, err := env.Client.LoadBalancer.RemoveLabelSelectorTarget(ctx, &LoadBalancer{ID: 1}, "key=value")
	if err != nil {
		t.Fatal(err)
	}
	if action.ID != 1 {
		t.Errorf("unexpected action ID: %d", action.ID)
	}

}

func TestLoadBalancerAddService(t *testing.T) {
	var (
		ctx          = context.Background()
		loadBalancer = &LoadBalancer{ID: 1}
	)

	t.Run("all fields", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()
		env.Mux.HandleFunc("/load_balancers/1/actions/add_service", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "POST" {
				t.Error("expected POST")
			}
			var reqBody schema.LoadBalancerActionAddServiceRequest
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				t.Fatal(err)
			}
			if reqBody.Protocol != "http" {
				t.Errorf("unexpected Protocol: %v", reqBody.Protocol)
			}
			if reqBody.ListenPort != 4711 {
				t.Errorf("unexpected ListenPort: %v", reqBody.ListenPort)
			}
			if reqBody.DestinationPort != 80 {
				t.Errorf("unexpected DestinationPort: %v", reqBody.DestinationPort)
			}
			if reqBody.HTTP == nil {
				t.Errorf("unexpected HTTP: %v", reqBody.HTTP)
			} else {
				if reqBody.HTTP.CookieName != "HCLBSTICKY" {
					t.Errorf("unexpectedTTP.CookieName: %v", reqBody.HTTP.CookieName)
				}
				if reqBody.HTTP.CookieLifetime != 300 {
					t.Errorf("unexpected HealthCheck.CookieLifetime: %v", reqBody.HTTP.CookieLifetime)
				}
			}
			if reqBody.HealthCheck.Protocol != "http" {
				t.Errorf("unexpected HealthCheck.Protocol: %v", reqBody.HealthCheck.Protocol)
			}
			if reqBody.HealthCheck.Port != 4711 {
				t.Errorf("unexpected HealthCheck.Port: %v", reqBody.HealthCheck.Port)
			}
			if reqBody.HealthCheck.Interval != 15 {
				t.Errorf("unexpected HealthCheck.Interval: %v", reqBody.HealthCheck.Interval)
			}
			if reqBody.HealthCheck.Timeout != 10 {
				t.Errorf("unexpected HealthCheck.Timeout: %v", reqBody.HealthCheck.Timeout)
			}
			if reqBody.HealthCheck.Retries != 3 {
				t.Errorf("unexpected HealthCheck.Retries: %v", reqBody.HealthCheck.Retries)
			}
			if reqBody.HealthCheck.HTTP == nil {
				t.Errorf("unexpected HealthCheck.HTTP: %v", reqBody.HealthCheck.HTTP)
			} else {
				if reqBody.HealthCheck.HTTP.Domain != "example.com" {
					t.Errorf("unexpected HealthCheck.HTTP.Domain: %v", reqBody.HealthCheck.HTTP.Domain)
				}
				if reqBody.HealthCheck.HTTP.Path != "/" {
					t.Errorf("unexpected HealthCheck.HTTP.Path: %v", reqBody.HealthCheck.HTTP.Path)
				}
			}
			json.NewEncoder(w).Encode(schema.LoadBalancerActionUpdateHealthCheckResponse{
				Action: schema.Action{
					ID: 1,
				},
			})
		})

		opts := LoadBalancerAddServiceOpts{
			Protocol:        LoadBalancerServiceProtocolHTTP,
			ListenPort:      4711,
			DestinationPort: 80,
			HTTP: &LoadBalancerServiceHTTP{
				CookieName:     "HCLBSTICKY",
				CookieLifetime: 5 * time.Minute,
			},
			HealthCheck: &LoadBalancerServiceHealthCheck{
				Protocol: "http",
				Port:     4711,
				Interval: 15 * time.Second,
				Timeout:  10 * time.Second,
				Retries:  3,
				HTTP: &LoadBalancerServiceHealthCheckHTTP{
					Domain: "example.com",
					Path:   "/",
				},
			},
		}
		action, _, err := env.Client.LoadBalancer.AddService(ctx, loadBalancer, opts)
		if err != nil {
			t.Fatal(err)
		}

		if action.ID != 1 {
			t.Errorf("unexpected action ID: %v", action.ID)
		}
	})

	t.Run("without health check", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()
		env.Mux.HandleFunc("/load_balancers/1/actions/add_service", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "POST" {
				t.Error("expected POST")
			}
			var reqBody schema.LoadBalancerActionAddServiceRequest
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				t.Fatal(err)
			}
			if reqBody.Protocol != "http" {
				t.Errorf("unexpected Protocol: %v", reqBody.Protocol)
			}
			if reqBody.ListenPort != 4711 {
				t.Errorf("unexpected ListenPort: %v", reqBody.ListenPort)
			}
			if reqBody.DestinationPort != 80 {
				t.Errorf("unexpected DestinationPort: %v", reqBody.DestinationPort)
			}
			if reqBody.HTTP == nil {
				t.Errorf("unexpected HTTP: %v", reqBody.HTTP)
			} else {
				if reqBody.HTTP.CookieName != "HCLBSTICKY" {
					t.Errorf("unexpectedTTP.CookieName: %v", reqBody.HTTP.CookieName)
				}
				if reqBody.HTTP.CookieLifetime != 300 {
					t.Errorf("unexpected HealthCheck.CookieLifetime: %v", reqBody.HTTP.CookieLifetime)
				}
			}
			if reqBody.HealthCheck != nil {
				t.Errorf("unexpected HealthCheck: %v", reqBody.HTTP)
			}
			json.NewEncoder(w).Encode(schema.LoadBalancerActionUpdateHealthCheckResponse{
				Action: schema.Action{
					ID: 1,
				},
			})
		})

		opts := LoadBalancerAddServiceOpts{
			Protocol:        LoadBalancerServiceProtocolHTTP,
			ListenPort:      4711,
			DestinationPort: 80,
			HTTP: &LoadBalancerServiceHTTP{
				CookieName:     "HCLBSTICKY",
				CookieLifetime: 5 * time.Minute,
			},
			HealthCheck: nil,
		}
		action, _, err := env.Client.LoadBalancer.AddService(ctx, loadBalancer, opts)
		if err != nil {
			t.Fatal(err)
		}

		if action.ID != 1 {
			t.Errorf("unexpected action ID: %v", action.ID)
		}
	})

	t.Run("protocol tcp and without http configuration", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()
		env.Mux.HandleFunc("/load_balancers/1/actions/add_service", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "POST" {
				t.Error("expected POST")
			}
			var reqBody schema.LoadBalancerActionAddServiceRequest
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				t.Fatal(err)
			}
			if reqBody.Protocol != "tcp" {
				t.Errorf("unexpected Protocol: %v", reqBody.Protocol)
			}
			if reqBody.ListenPort != 4711 {
				t.Errorf("unexpected ListenPort: %v", reqBody.ListenPort)
			}
			if reqBody.DestinationPort != 80 {
				t.Errorf("unexpected DestinationPort: %v", reqBody.DestinationPort)
			}
			if reqBody.HTTP != nil {
				t.Errorf("unexpected HTTP: %v", reqBody.HTTP)
			}
			if reqBody.HealthCheck.Protocol != "http" {
				t.Errorf("unexpected HealthCheck.Protocol: %v", reqBody.HealthCheck.Protocol)
			}
			if reqBody.HealthCheck.Port != 4711 {
				t.Errorf("unexpected HealthCheck.Port: %v", reqBody.HealthCheck.Port)
			}
			if reqBody.HealthCheck.Interval != 15 {
				t.Errorf("unexpected HealthCheck.Interval: %v", reqBody.HealthCheck.Interval)
			}
			if reqBody.HealthCheck.Timeout != 10 {
				t.Errorf("unexpected HealthCheck.Timeout: %v", reqBody.HealthCheck.Timeout)
			}
			if reqBody.HealthCheck.Retries != 3 {
				t.Errorf("unexpected HealthCheck.Retries: %v", reqBody.HealthCheck.Retries)
			}
			if reqBody.HealthCheck.HTTP == nil {
				t.Errorf("unexpected HealthCheck.HTTP: %v", reqBody.HealthCheck.HTTP)
			} else {
				if reqBody.HealthCheck.HTTP.Domain != "example.com" {
					t.Errorf("unexpected HealthCheck.HTTP.Domain: %v", reqBody.HealthCheck.HTTP.Domain)
				}
				if reqBody.HealthCheck.HTTP.Path != "/" {
					t.Errorf("unexpected HealthCheck.HTTP.Path: %v", reqBody.HealthCheck.HTTP.Path)
				}
			}
			json.NewEncoder(w).Encode(schema.LoadBalancerActionUpdateHealthCheckResponse{
				Action: schema.Action{
					ID: 1,
				},
			})
		})

		opts := LoadBalancerAddServiceOpts{
			Protocol:        LoadBalancerServiceProtocolTCP,
			ListenPort:      4711,
			DestinationPort: 80,
			HTTP:            nil,
			HealthCheck: &LoadBalancerServiceHealthCheck{
				Protocol: "http",
				Port:     4711,
				Interval: 15 * time.Second,
				Timeout:  10 * time.Second,
				Retries:  3,
				HTTP: &LoadBalancerServiceHealthCheckHTTP{
					Domain: "example.com",
					Path:   "/",
				},
			},
		}
		action, _, err := env.Client.LoadBalancer.AddService(ctx, loadBalancer, opts)
		if err != nil {
			t.Fatal(err)
		}

		if action.ID != 1 {
			t.Errorf("unexpected action ID: %v", action.ID)
		}
	})
}

func TestLoadBalancerDeleteService(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/load_balancers/1/actions/delete_service", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Error("expected POST")
		}
		var reqBody schema.LoadBalancerDeleteServiceRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatal(err)
		}
		if reqBody.ListenPort != 4711 {
			t.Errorf("unexpected ListenPort %v", reqBody.ListenPort)
		}
		json.NewEncoder(w).Encode(schema.LoadBalancerDeleteServiceResponse{
			Action: schema.Action{
				ID: 1,
			},
		})
	})

	var (
		ctx          = context.Background()
		loadBalancer = &LoadBalancer{ID: 1}
	)
	action, _, err := env.Client.LoadBalancer.DeleteService(ctx, loadBalancer, 4711)
	if err != nil {
		t.Fatal(err)
	}
	if action.ID != 1 {
		t.Errorf("unexpected action ID: %d", action.ID)
	}
}

func TestLoadBalancerClientChangeAlgorithm(t *testing.T) {
	var (
		ctx          = context.Background()
		loadBalancer = &LoadBalancer{ID: 1}
	)

	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/load_balancers/1/actions/change_algorithm", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Error("expected POST")
		}
		var reqBody schema.LoadBalancerActionChangeAlgorithmRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatal(err)
		}
		if reqBody.Type != "round_robin" {
			t.Errorf("unexpected type: %v", reqBody.Type)
		}
		json.NewEncoder(w).Encode(schema.LoadBalancerActionChangeAlgorithmResponse{
			Action: schema.Action{
				ID: 1,
			},
		})
	})

	opts := LoadBalancerChangeAlgorithmOpts{
		Type: LoadBalancerAlgorithmTypeRoundRobin,
	}
	action, _, err := env.Client.LoadBalancer.ChangeAlgorithm(ctx, loadBalancer, opts)
	if err != nil {
		t.Fatal(err)
	}

	if action.ID != 1 {
		t.Errorf("unexpected action ID: %v", action.ID)
	}

}

func TestLoadBalancerUpdateHealthCheck(t *testing.T) {
	var (
		ctx          = context.Background()
		loadBalancer = &LoadBalancer{ID: 1}
	)

	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/load_balancers/1/actions/update_health_check", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Error("expected POST")
		}
		var reqBody schema.LoadBalancerActionUpdateHealthCheckRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatal(err)
		}
		if reqBody.ListenPort != 4711 {
			t.Errorf("unexpected ListenPort: %v", reqBody.ListenPort)
		}
		if reqBody.HealthCheck.Protocol != "http" {
			t.Errorf("unexpected HealthCheck.Protocol: %v", reqBody.HealthCheck.Protocol)
		}
		if reqBody.HealthCheck.Port != 4711 {
			t.Errorf("unexpected HealthCheck.Port: %v", reqBody.HealthCheck.Port)
		}
		if reqBody.HealthCheck.Interval != 15 {
			t.Errorf("unexpected HealthCheck.Interval: %v", reqBody.HealthCheck.Interval)
		}
		if reqBody.HealthCheck.Timeout != 10 {
			t.Errorf("unexpected HealthCheck.Timeout: %v", reqBody.HealthCheck.Timeout)
		}
		if reqBody.HealthCheck.Retries != 3 {
			t.Errorf("unexpected HealthCheck.Retries: %v", reqBody.HealthCheck.Retries)
		}
		if reqBody.HealthCheck.HTTP == nil {
			t.Errorf("unexpected HealthCheck.HTTP: %v", reqBody.HealthCheck.HTTP)
		} else {
			if reqBody.HealthCheck.HTTP.Domain != "example.com" {
				t.Errorf("unexpected HealthCheck.HTTP.Domain: %v", reqBody.HealthCheck.HTTP.Domain)
			}
			if reqBody.HealthCheck.HTTP.Path != "/" {
				t.Errorf("unexpected HealthCheck.HTTP.Path: %v", reqBody.HealthCheck.HTTP.Path)
			}
		}
		json.NewEncoder(w).Encode(schema.LoadBalancerActionUpdateHealthCheckResponse{
			Action: schema.Action{
				ID: 1,
			},
		})
	})

	opts := LoadBalancerUpdateHealthCheckOpts{
		ListenPort: 4711,
		HealthCheck: LoadBalancerServiceHealthCheck{
			Protocol: "http",
			Port:     4711,
			Interval: 15 * time.Second,
			Timeout:  10 * time.Second,
			Retries:  3,
			HTTP: &LoadBalancerServiceHealthCheckHTTP{
				Domain: "example.com",
				Path:   "/",
			},
		},
	}
	action, _, err := env.Client.LoadBalancer.UpdateHealthCheck(ctx, loadBalancer, opts)
	if err != nil {
		t.Fatal(err)
	}

	if action.ID != 1 {
		t.Errorf("unexpected action ID: %v", action.ID)
	}
}
