package hcloud

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

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
		t.Fatal("no load balancer")
	}
	if loadBalancer.ID != 1 {
		t.Errorf("unexpected load balancer ID: %v", loadBalancer.ID)
	}

	t.Run("called via Get", func(t *testing.T) {
		loadBalancer, _, err := env.Client.LoadBalancer.Get(ctx, "1")
		if err != nil {
			t.Fatal(err)
		}
		if loadBalancer == nil {
			t.Fatal("no load balancer")
		}
		if loadBalancer.ID != 1 {
			t.Errorf("unexpected load balancer ID: %v", loadBalancer.ID)
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
		t.Fatal("expected no load balancer")
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
		t.Fatal("no load balancer")
	}
	if loadBalancer.ID != 1 {
		t.Errorf("unexpected load balancer ID: %v", loadBalancer.ID)
	}

	t.Run("via Get", func(t *testing.T) {
		loadBalancer, _, err := env.Client.LoadBalancer.Get(ctx, "mylb")
		if err != nil {
			t.Fatal(err)
		}
		if loadBalancer == nil {
			t.Fatal("no load balancer")
		}
		if loadBalancer.ID != 1 {
			t.Errorf("unexpected load balancer ID: %v", loadBalancer.ID)
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
		t.Fatal("unexpected load balancer")
	}
}

func TestLoadBalancerClientGetByNameEmpty(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	ctx := context.Background()

	loadBalancer, _, err := env.Client.LoadBalancer.GetByName(ctx, "")
	if err != nil {
		t.Fatal(err)
	}
	if loadBalancer != nil {
		t.Fatal("unexpected load balancer")
	}
}

func TestLoadBalancerCreate(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/load_balancers", func(w http.ResponseWriter, r *http.Request) {
		var reqBody schema.LoadBalancerCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatal(err)
		}
		expectedReqBody := schema.LoadBalancerCreateRequest{
			Name:             "load-balancer",
			LoadBalancerType: "lb1",
			Algorithm: &schema.LoadBalancerCreateRequestAlgorithm{
				Type: "round_robin",
			},
			Location: String("fsn1"),
		}
		if !cmp.Equal(expectedReqBody, reqBody) {
			t.Log(cmp.Diff(expectedReqBody, reqBody))
			t.Error("unexpected request body")
		}
		json.NewEncoder(w).Encode(schema.LoadBalancerCreateResponse{
			LoadBalancer: schema.LoadBalancer{ID: 2},
			Action:       schema.Action{ID: 1},
		})
	})

	var (
		ctx       = context.Background()
		lbType    = &LoadBalancerType{Name: "lb1"}
		algorithm = &LoadBalancerAlgorithm{Type: LoadBalancerAlgorithmTypeRoundRobin}
		location  = &Location{Name: "fsn1"}
		opts      = LoadBalancerCreateOpts{
			Name:             "load-balancer",
			LoadBalancerType: lbType,
			Algorithm:        algorithm,
			Location:         location,
		}
	)

	result, _, err := env.Client.LoadBalancer.Create(ctx, opts)
	if err != nil {
		t.Fatal(err)
	}
	if result.Action.ID != 1 {
		t.Errorf("unexpected action ID: %d", result.Action.ID)
	}
	if result.LoadBalancer.ID != 2 {
		t.Errorf("unexpected load balancer ID: %d", result.LoadBalancer.ID)
	}
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
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/load_balancers/1", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Error("expected PUT")
		}
		var reqBody schema.LoadBalancerUpdateRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatal(err)
		}
		expectedReqBody := schema.LoadBalancerUpdateRequest{
			Name: String("test"),
		}
		if !cmp.Equal(expectedReqBody, reqBody) {
			t.Log(cmp.Diff(expectedReqBody, reqBody))
			t.Error("unexpected request body")
		}
		json.NewEncoder(w).Encode(schema.LoadBalancerUpdateResponse{
			LoadBalancer: schema.LoadBalancer{
				ID: 1,
			},
		})
	})

	var (
		ctx          = context.Background()
		loadBalancer = &LoadBalancer{ID: 1}
	)

	opts := LoadBalancerUpdateOpts{
		Name: "test",
	}
	updatedLoadBalancer, _, err := env.Client.LoadBalancer.Update(ctx, loadBalancer, opts)
	if err != nil {
		t.Fatal(err)
	}
	if updatedLoadBalancer.ID != 1 {
		t.Errorf("unexpected load balancer ID: %v", updatedLoadBalancer.ID)
	}
}

func TestLoadBalancerClientChangeProtection(t *testing.T) {
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
		expectedReqBody := schema.LoadBalancerActionChangeProtectionRequest{
			Delete: Bool(true),
		}
		if !cmp.Equal(expectedReqBody, reqBody) {
			t.Log(cmp.Diff(expectedReqBody, reqBody))
			t.Error("unexpected request body")
		}
		json.NewEncoder(w).Encode(schema.LoadBalancerActionChangeProtectionResponse{
			Action: schema.Action{
				ID: 1,
			},
		})
	})

	var (
		ctx          = context.Background()
		loadBalancer = &LoadBalancer{ID: 1}
	)

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
}

func TestLoadBalancerClientAddServerTarget(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/load_balancers/1/actions/add_target", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Error("expected POST")
		}
		var reqBody schema.LoadBalancerActionAddTargetRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatal(err)
		}
		expectedReqBody := schema.LoadBalancerActionAddTargetRequest{
			Type: string(LoadBalancerTargetTypeServer),
			Server: &schema.LoadBalancerActionAddTargetRequestServer{
				ID: 1,
			},
			UsePrivateIP: Bool(true),
		}
		if !cmp.Equal(expectedReqBody, reqBody) {
			t.Log(cmp.Diff(expectedReqBody, reqBody))
			t.Error("unexpected request body")
		}
		json.NewEncoder(w).Encode(schema.LoadBalancerActionAddTargetResponse{
			Action: schema.Action{
				ID: 1,
			},
		})
	})

	var (
		ctx          = context.Background()
		loadBalancer = &LoadBalancer{ID: 1}
		server       = &Server{ID: 1}
	)

	opts := LoadBalancerAddServerTargetOpts{
		Server:       server,
		UsePrivateIP: Bool(true),
	}
	action, _, err := env.Client.LoadBalancer.AddServerTarget(ctx, loadBalancer, opts)
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
		var reqBody schema.LoadBalancerActionRemoveTargetRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatal(err)
		}
		expectedReqBody := schema.LoadBalancerActionRemoveTargetRequest{
			Type: string(LoadBalancerTargetTypeServer),
			Server: &schema.LoadBalancerActionRemoveTargetRequestServer{
				ID: 1,
			},
		}
		if !cmp.Equal(expectedReqBody, reqBody) {
			t.Log(cmp.Diff(expectedReqBody, reqBody))
			t.Error("unexpected request body")
		}
		json.NewEncoder(w).Encode(schema.LoadBalancerActionRemoveTargetResponse{
			Action: schema.Action{
				ID: 1,
			},
		})
	})

	var (
		ctx          = context.Background()
		loadBalancer = &LoadBalancer{ID: 1}
		server       = &Server{ID: 1}
	)

	action, _, err := env.Client.LoadBalancer.RemoveServerTarget(ctx, loadBalancer, server)
	if err != nil {
		t.Fatal(err)
	}
	if action.ID != 1 {
		t.Errorf("unexpected action ID: %d", action.ID)
	}
}

func TestLoadBalancerAddService(t *testing.T) {
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
		expectedReqBody := schema.LoadBalancerActionAddServiceRequest{
			Protocol:        string(LoadBalancerServiceProtocolHTTP),
			ListenPort:      Int(4711),
			DestinationPort: Int(80),
			HTTP: &schema.LoadBalancerActionAddServiceRequestHTTP{
				CookieName:     String("HCLBSTICKY"),
				CookieLifetime: Int(5 * 60),
				RedirectHTTP:   Bool(false),
				StickySessions: Bool(true),
			},
			HealthCheck: &schema.LoadBalancerActionAddServiceRequestHealthCheck{
				Protocol: "http",
				Port:     Int(4711),
				Interval: Int(15),
				Timeout:  Int(10),
				Retries:  Int(3),
				HTTP: &schema.LoadBalancerActionAddServiceRequestHealthCheckHTTP{
					Domain: String("example.com"),
					Path:   String("/"),
				},
			},
		}
		if !cmp.Equal(expectedReqBody, reqBody) {
			t.Log(cmp.Diff(expectedReqBody, reqBody))
			t.Error("unexpected request body")
		}
		json.NewEncoder(w).Encode(schema.LoadBalancerActionAddServiceResponse{
			Action: schema.Action{
				ID: 1,
			},
		})
	})

	var (
		ctx          = context.Background()
		loadBalancer = &LoadBalancer{ID: 1}
	)

	opts := LoadBalancerAddServiceOpts{
		Protocol:        LoadBalancerServiceProtocolHTTP,
		ListenPort:      Int(4711),
		DestinationPort: Int(80),
		HTTP: &LoadBalancerAddServiceOptsHTTP{
			CookieName:     String("HCLBSTICKY"),
			CookieLifetime: Duration(5 * time.Minute),
			RedirectHTTP:   Bool(false),
			StickySessions: Bool(true),
		},
		HealthCheck: &LoadBalancerAddServiceOptsHealthCheck{
			Protocol: "http",
			Port:     Int(4711),
			Interval: Duration(15 * time.Second),
			Timeout:  Duration(10 * time.Second),
			Retries:  Int(3),
			HTTP: &LoadBalancerAddServiceOptsHealthCheckHTTP{
				Domain: String("example.com"),
				Path:   String("/"),
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
}

func TestLoadBalancerUpdateService(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/load_balancers/1/actions/update_service", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Error("expected POST")
		}
		var reqBody schema.LoadBalancerActionUpdateServiceRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatal(err)
		}
		expectedReqBody := schema.LoadBalancerActionUpdateServiceRequest{
			Protocol:        String(string(LoadBalancerServiceProtocolHTTP)),
			ListenPort:      4711,
			DestinationPort: Int(80),
			HTTP: &schema.LoadBalancerActionUpdateServiceRequestHTTP{
				CookieName:     String("HCLBSTICKY"),
				CookieLifetime: Int(5 * 60),
				RedirectHTTP:   Bool(false),
				StickySessions: Bool(true),
			},
			HealthCheck: &schema.LoadBalancerActionUpdateServiceRequestHealthCheck{
				Protocol: String(string(LoadBalancerServiceProtocolHTTP)),
				Port:     Int(4711),
				Interval: Int(15),
				Timeout:  Int(10),
				Retries:  Int(3),
				HTTP: &schema.LoadBalancerActionUpdateServiceRequestHealthCheckHTTP{
					Domain: String("example.com"),
					Path:   String("/"),
				},
			},
		}
		if !cmp.Equal(expectedReqBody, reqBody) {
			t.Log(cmp.Diff(expectedReqBody, reqBody))
			t.Error("unexpected request body")
		}
		json.NewEncoder(w).Encode(schema.LoadBalancerActionUpdateServiceResponse{
			Action: schema.Action{
				ID: 1,
			},
		})
	})

	var (
		ctx          = context.Background()
		loadBalancer = &LoadBalancer{ID: 1}
	)

	opts := LoadBalancerUpdateServiceOpts{
		Protocol:        LoadBalancerServiceProtocolHTTP,
		DestinationPort: Int(80),
		HTTP: &LoadBalancerUpdateServiceOptsHTTP{
			CookieName:     String("HCLBSTICKY"),
			CookieLifetime: Duration(5 * time.Minute),
			RedirectHTTP:   Bool(false),
			StickySessions: Bool(true),
		},
		HealthCheck: &LoadBalancerUpdateServiceOptsHealthCheck{
			Protocol: LoadBalancerServiceProtocolHTTP,
			Port:     Int(4711),
			Interval: Duration(15 * time.Second),
			Timeout:  Duration(10 * time.Second),
			Retries:  Int(3),
			HTTP: &LoadBalancerUpdateServiceOptsHealthCheckHTTP{
				Domain: String("example.com"),
				Path:   String("/"),
			},
		},
	}
	action, _, err := env.Client.LoadBalancer.UpdateService(ctx, loadBalancer, 4711, opts)
	if err != nil {
		t.Fatal(err)
	}
	if action.ID != 1 {
		t.Errorf("unexpected action ID: %v", action.ID)
	}
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
		expectedReqBody := schema.LoadBalancerDeleteServiceRequest{
			ListenPort: 4711,
		}
		if !cmp.Equal(expectedReqBody, reqBody) {
			t.Log(cmp.Diff(expectedReqBody, reqBody))
			t.Error("unexpected request body")
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
		expectedReqBody := schema.LoadBalancerActionChangeAlgorithmRequest{
			Type: string(LoadBalancerAlgorithmTypeRoundRobin),
		}
		if !cmp.Equal(expectedReqBody, reqBody) {
			t.Log(cmp.Diff(expectedReqBody, reqBody))
			t.Error("unexpected request body")
		}
		json.NewEncoder(w).Encode(schema.LoadBalancerActionChangeAlgorithmResponse{
			Action: schema.Action{
				ID: 1,
			},
		})
	})

	var (
		ctx          = context.Background()
		loadBalancer = &LoadBalancer{ID: 1}
	)

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

func TestLoadBalancerClientAttachToNetwork(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/load_balancers/1/actions/attach_to_network", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Error("expected POST")
		}
		var reqBody schema.LoadBalancerActionAttachToNetworkRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatal(err)
		}
		expectedReqBody := schema.LoadBalancerActionAttachToNetworkRequest{
			Network: 1,
			IP:      String("10.0.1.1"),
		}
		if !cmp.Equal(expectedReqBody, reqBody) {
			t.Log(cmp.Diff(expectedReqBody, reqBody))
			t.Error("unexpected request body")
		}
		json.NewEncoder(w).Encode(schema.LoadBalancerActionAttachToNetworkResponse{
			Action: schema.Action{
				ID: 1,
			},
		})
	})

	var (
		ctx          = context.Background()
		loadBalancer = &LoadBalancer{ID: 1}
		network      = &Network{ID: 1}
	)

	opts := LoadBalancerAttachToNetworkOpts{
		Network: network,
		IP:      net.ParseIP("10.0.1.1"),
	}
	action, _, err := env.Client.LoadBalancer.AttachToNetwork(ctx, loadBalancer, opts)
	if err != nil {
		t.Fatal(err)
	}
	if action.ID != 1 {
		t.Errorf("unexpected action ID: %v", action.ID)
	}
}

func TestLoadBalancerClientDetachFromNetwork(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/load_balancers/1/actions/detach_from_network", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Error("expected POST")
		}
		var reqBody schema.LoadBalancerActionDetachFromNetworkRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatal(err)
		}
		expectedReqBody := schema.LoadBalancerActionDetachFromNetworkRequest{
			Network: 1,
		}
		if !cmp.Equal(expectedReqBody, reqBody) {
			t.Log(cmp.Diff(expectedReqBody, reqBody))
			t.Error("unexpected request body")
		}
		json.NewEncoder(w).Encode(schema.LoadBalancerActionDetachFromNetworkResponse{
			Action: schema.Action{
				ID: 1,
			},
		})
	})

	var (
		ctx          = context.Background()
		loadBalancer = &LoadBalancer{ID: 1}
		network      = &Network{ID: 1}
	)

	opts := LoadBalancerDetachFromNetworkOpts{
		Network: network,
	}
	action, _, err := env.Client.LoadBalancer.DetachFromNetwork(ctx, loadBalancer, opts)
	if err != nil {
		t.Fatal(err)
	}
	if action.ID != 1 {
		t.Errorf("unexpected action ID: %v", action.ID)
	}
}

func TestLoadBalancerClientEnablePublicInterface(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/load_balancers/1/actions/enable_public_interface", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(schema.LoadBalancerActionEnablePublicInterfaceResponse{
			Action: schema.Action{
				ID: 1,
			},
		})
	})

	var (
		ctx          = context.Background()
		loadBalancer = &LoadBalancer{ID: 1}
	)

	action, _, err := env.Client.LoadBalancer.EnablePublicInterface(ctx, loadBalancer)
	if err != nil {
		t.Fatal(err)
	}
	if action.ID != 1 {
		t.Errorf("unexpected action ID: %d", action.ID)
	}
}

func TestLoadBalancerClientDisablePublicInterface(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/load_balancers/1/actions/disable_public_interface", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(schema.LoadBalancerActionDisablePublicInterfaceResponse{
			Action: schema.Action{
				ID: 1,
			},
		})
	})

	var (
		ctx          = context.Background()
		loadBalancer = &LoadBalancer{ID: 1}
	)

	action, _, err := env.Client.LoadBalancer.DisablePublicInterface(ctx, loadBalancer)
	if err != nil {
		t.Fatal(err)
	}
	if action.ID != 1 {
		t.Errorf("unexpected action ID: %d", action.ID)
	}
}

func TestLoadBalancerClientChangeType(t *testing.T) {
	var (
		ctx          = context.Background()
		loadBalancer = &LoadBalancer{ID: 1}
	)

	t.Run("with Load Balancer type ID", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/load_balancers/1/actions/change_type", func(w http.ResponseWriter, r *http.Request) {
			var reqBody schema.LoadBalancerActionChangeTypeRequest
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				t.Fatal(err)
			}
			if id, ok := reqBody.LoadBalancerType.(float64); !ok || id != 1 {
				t.Errorf("unexpected Load Balancer type ID: %v", reqBody.LoadBalancerType)
			}
			json.NewEncoder(w).Encode(schema.LoadBalancerActionChangeTypeResponse{
				Action: schema.Action{
					ID: 1,
				},
			})
		})

		opts := LoadBalancerChangeTypeOpts{
			LoadBalancerType: &LoadBalancerType{ID: 1},
		}
		action, _, err := env.Client.LoadBalancer.ChangeType(ctx, loadBalancer, opts)
		if err != nil {
			t.Fatal(err)
		}
		if action.ID != 1 {
			t.Errorf("unexpected action ID: %d", action.ID)
		}
	})

	t.Run("with Load Balancer type name", func(t *testing.T) {
		env := newTestEnv()
		defer env.Teardown()

		env.Mux.HandleFunc("/load_balancers/1/actions/change_type", func(w http.ResponseWriter, r *http.Request) {
			var reqBody schema.LoadBalancerActionChangeTypeRequest
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				t.Fatal(err)
			}
			if name, ok := reqBody.LoadBalancerType.(string); !ok || name != "type" {
				t.Errorf("unexpected Load Balancer type name: %v", reqBody.LoadBalancerType)
			}
			json.NewEncoder(w).Encode(schema.LoadBalancerActionChangeTypeResponse{
				Action: schema.Action{
					ID: 1,
				},
			})
		})

		opts := LoadBalancerChangeTypeOpts{
			LoadBalancerType: &LoadBalancerType{Name: "type"},
		}
		action, _, err := env.Client.LoadBalancer.ChangeType(ctx, loadBalancer, opts)
		if err != nil {
			t.Fatal(err)
		}
		if action.ID != 1 {
			t.Errorf("unexpected action ID: %d", action.ID)
		}
	})
}

func TestLoadBalancerClientAddLabelSelectorTarget(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/load_balancers/1/actions/add_target", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Error("expected POST")
		}
		var reqBody schema.LoadBalancerActionAddTargetRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatal(err)
		}
		if reqBody.Type != string(LoadBalancerTargetTypeLabelSelector) {
			t.Errorf("unexpected type %v", reqBody.Type)
		}
		if reqBody.LabelSelector.Selector != "key=value" {
			t.Errorf("unexpected LabelSelector %v", reqBody.LabelSelector)
		}
		if *reqBody.UsePrivateIP != false {
			t.Errorf("unexpected UsePrivateIP %v", reqBody.UsePrivateIP)
		}
		json.NewEncoder(w).Encode(schema.LoadBalancerActionAddTargetResponse{
			Action: schema.Action{
				ID: 1,
			},
		})
	})

	ctx := context.Background()
	action, _, err := env.Client.LoadBalancer.AddLabelSelectorTarget(ctx, &LoadBalancer{ID: 1}, LoadBalancerAddLabelSelectorTargetOpts{
		Selector:     "key=value",
		UsePrivateIP: Bool(false),
	})
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
		var reqBody schema.LoadBalancerActionRemoveTargetRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatal(err)
		}
		if reqBody.Type != string(LoadBalancerTargetTypeLabelSelector) {
			t.Errorf("unexpected type %v", reqBody.Type)
		}
		if reqBody.LabelSelector.Selector != "key=value" {
			t.Errorf("unexpected LabelSelector %v", reqBody.LabelSelector)
		}
		json.NewEncoder(w).Encode(schema.LoadBalancerActionRemoveTargetResponse{
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

func TestLoadBalancerClientAddIPTarget(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/load_balancers/1/actions/add_target", func(w http.ResponseWriter, r *http.Request) {
		var reqBody schema.LoadBalancerActionAddTargetRequest

		if r.Method != "POST" {
			t.Error("expected POST")
		}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatal(err)
		}
		if reqBody.Type != string(LoadBalancerTargetTypeIP) {
			t.Errorf("unexpected type %v", reqBody.Type)
		}
		if reqBody.IP.IP != "1.2.3.4" {
			t.Errorf("unexpected IP target %v", reqBody.IP)
		}
		json.NewEncoder(w).Encode(schema.LoadBalancerActionAddTargetResponse{
			Action: schema.Action{
				ID: 1,
			},
		})
	})

	ctx := context.Background()
	action, _, err := env.Client.LoadBalancer.AddIPTarget(ctx, &LoadBalancer{ID: 1}, LoadBalancerAddIPTargetOpts{
		IP: net.ParseIP("1.2.3.4"),
	})
	if err != nil {
		t.Fatal(err)
	}
	if action.ID != 1 {
		t.Errorf("unexpected action ID: %d", action.ID)
	}
}

func TestLoadBalancerClientRemoveIPTarget(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/load_balancers/1/actions/remove_target", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Error("expected POST")
		}
		var reqBody schema.LoadBalancerActionRemoveTargetRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatal(err)
		}
		if reqBody.Type != string(LoadBalancerTargetTypeIP) {
			t.Errorf("unexpected type %v", reqBody.Type)
		}
		if reqBody.IP.IP != "1.2.3.4" {
			t.Errorf("unexpected IP %v", reqBody.IP)
		}
		json.NewEncoder(w).Encode(schema.LoadBalancerActionRemoveTargetResponse{
			Action: schema.Action{
				ID: 1,
			},
		})
	})

	ctx := context.Background()
	action, _, err := env.Client.LoadBalancer.RemoveIPTarget(ctx, &LoadBalancer{ID: 1}, net.ParseIP("1.2.3.4"))
	if err != nil {
		t.Fatal(err)
	}
	if action.ID != 1 {
		t.Errorf("unexpected action ID: %d", action.ID)
	}
}
