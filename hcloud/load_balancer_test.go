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

func TestLoadBalancerClientGetByNumericName(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/load_balancers/123", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(schema.ErrorResponse{
			Error: schema.Error{
				Code: string(ErrorCodeNotFound),
			},
		})
	})

	env.Mux.HandleFunc("/load_balancers", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RawQuery != "name=123" {
			t.Fatal("missing name query")
		}
		json.NewEncoder(w).Encode(schema.LoadBalancerListResponse{
			LoadBalancers: []schema.LoadBalancer{
				{
					ID:   1,
					Name: "123",
				},
			},
		})
	})

	ctx := context.Background()

	loadBalancer, _, err := env.Client.LoadBalancer.Get(ctx, "123")
	if err != nil {
		t.Fatal(err)
	}
	if loadBalancer == nil {
		t.Fatal("no load balancer")
	}
	if loadBalancer.ID != 1 {
		t.Errorf("unexpected load balancer ID: %v", loadBalancer.ID)
	}
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
			LoadBalancerType: schema.IDOrName{Name: "lb1"},
			Algorithm: &schema.LoadBalancerCreateRequestAlgorithm{
				Type: "round_robin",
			},
			Location: Ptr("fsn1"),
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

	env.Mux.HandleFunc("/load_balancers/1", func(w http.ResponseWriter, r *http.Request) {})

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
			Name: Ptr("test"),
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
			Delete: Ptr(true),
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
		Delete: Ptr(true),
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
			UsePrivateIP: Ptr(true),
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
		UsePrivateIP: Ptr(true),
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
			ListenPort:      Ptr(4711),
			DestinationPort: Ptr(80),
			HTTP: &schema.LoadBalancerActionAddServiceRequestHTTP{
				CookieName:     Ptr("HCLBSTICKY"),
				CookieLifetime: Ptr(5 * 60),
				RedirectHTTP:   Ptr(false),
				StickySessions: Ptr(true),
			},
			HealthCheck: &schema.LoadBalancerActionAddServiceRequestHealthCheck{
				Protocol: "http",
				Port:     Ptr(4711),
				Interval: Ptr(15),
				Timeout:  Ptr(10),
				Retries:  Ptr(3),
				HTTP: &schema.LoadBalancerActionAddServiceRequestHealthCheckHTTP{
					Domain: Ptr("example.com"),
					Path:   Ptr("/"),
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
		ListenPort:      Ptr(4711),
		DestinationPort: Ptr(80),
		HTTP: &LoadBalancerAddServiceOptsHTTP{
			CookieName:     Ptr("HCLBSTICKY"),
			CookieLifetime: Ptr(5 * time.Minute),
			RedirectHTTP:   Ptr(false),
			StickySessions: Ptr(true),
		},
		HealthCheck: &LoadBalancerAddServiceOptsHealthCheck{
			Protocol: "http",
			Port:     Ptr(4711),
			Interval: Ptr(15 * time.Second),
			Timeout:  Ptr(10 * time.Second),
			Retries:  Ptr(3),
			HTTP: &LoadBalancerAddServiceOptsHealthCheckHTTP{
				Domain: Ptr("example.com"),
				Path:   Ptr("/"),
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
			Protocol:        Ptr(string(LoadBalancerServiceProtocolHTTP)),
			ListenPort:      4711,
			DestinationPort: Ptr(80),
			HTTP: &schema.LoadBalancerActionUpdateServiceRequestHTTP{
				CookieName:     Ptr("HCLBSTICKY"),
				CookieLifetime: Ptr(5 * 60),
				RedirectHTTP:   Ptr(false),
				StickySessions: Ptr(true),
			},
			HealthCheck: &schema.LoadBalancerActionUpdateServiceRequestHealthCheck{
				Protocol: Ptr(string(LoadBalancerServiceProtocolHTTP)),
				Port:     Ptr(4711),
				Interval: Ptr(15),
				Timeout:  Ptr(10),
				Retries:  Ptr(3),
				HTTP: &schema.LoadBalancerActionUpdateServiceRequestHealthCheckHTTP{
					Domain: Ptr("example.com"),
					Path:   Ptr("/"),
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
		DestinationPort: Ptr(80),
		HTTP: &LoadBalancerUpdateServiceOptsHTTP{
			CookieName:     Ptr("HCLBSTICKY"),
			CookieLifetime: Ptr(5 * time.Minute),
			RedirectHTTP:   Ptr(false),
			StickySessions: Ptr(true),
		},
		HealthCheck: &LoadBalancerUpdateServiceOptsHealthCheck{
			Protocol: LoadBalancerServiceProtocolHTTP,
			Port:     Ptr(4711),
			Interval: Ptr(15 * time.Second),
			Timeout:  Ptr(10 * time.Second),
			Retries:  Ptr(3),
			HTTP: &LoadBalancerUpdateServiceOptsHealthCheckHTTP{
				Domain: Ptr("example.com"),
				Path:   Ptr("/"),
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
			IP:      Ptr("10.0.1.1"),
			IPRange: Ptr("10.0.1.0/24"),
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

	_, ipRange, _ := net.ParseCIDR("10.0.1.0/24")
	opts := LoadBalancerAttachToNetworkOpts{
		Network: network,
		IP:      net.ParseIP("10.0.1.1"),
		IPRange: ipRange,
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
			if reqBody.LoadBalancerType.ID != 1 {
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
			if reqBody.LoadBalancerType.Name != "type" {
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
		UsePrivateIP: Ptr(false),
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

func TestLoadBalancerGetMetrics(t *testing.T) {
	tests := []struct {
		name        string
		lb          *LoadBalancer
		opts        LoadBalancerGetMetricsOpts
		respStatus  int
		respFn      func() schema.LoadBalancerGetMetricsResponse
		expected    LoadBalancerMetrics
		expectedErr string
	}{
		{
			name: "all metrics",
			lb:   &LoadBalancer{ID: 2},
			opts: LoadBalancerGetMetricsOpts{
				Types: []LoadBalancerMetricType{
					LoadBalancerMetricOpenConnections,
					LoadBalancerMetricConnectionsPerSecond,
					LoadBalancerMetricRequestsPerSecond,
					LoadBalancerMetricBandwidth,
				},
				Start: mustParseTime(t, "2017-01-01T00:00:00Z"),
				End:   mustParseTime(t, "2017-01-01T23:00:00Z"),
			},
			respFn: func() schema.LoadBalancerGetMetricsResponse {
				var resp schema.LoadBalancerGetMetricsResponse

				resp.Metrics.Start = mustParseTime(t, "2017-01-01T00:00:00Z")
				resp.Metrics.End = mustParseTime(t, "2017-01-01T23:00:00Z")
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

				return resp
			},
			expected: LoadBalancerMetrics{
				Start: mustParseTime(t, "2017-01-01T00:00:00Z"),
				End:   mustParseTime(t, "2017-01-01T23:00:00Z"),
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
		{
			name: "missing metrics types",
			lb:   &LoadBalancer{ID: 3},
			opts: LoadBalancerGetMetricsOpts{
				Start: mustParseTime(t, "2017-01-01T00:00:00Z"),
				End:   mustParseTime(t, "2017-01-01T23:00:00Z"),
			},
			expectedErr: "missing field [Types] in [hcloud.LoadBalancerGetMetricsOpts]",
		},
		{
			name: "no start time",
			lb:   &LoadBalancer{ID: 4},
			opts: LoadBalancerGetMetricsOpts{
				Types: []LoadBalancerMetricType{LoadBalancerMetricBandwidth},
				End:   mustParseTime(t, "2017-01-01T23:00:00Z"),
			},
			expectedErr: "missing field [Start] in [hcloud.LoadBalancerGetMetricsOpts]",
		},
		{
			name: "no end time",
			lb:   &LoadBalancer{ID: 5},
			opts: LoadBalancerGetMetricsOpts{
				Types: []LoadBalancerMetricType{LoadBalancerMetricBandwidth},
				Start: mustParseTime(t, "2017-01-01T00:00:00Z"),
			},
			expectedErr: "missing field [End] in [hcloud.LoadBalancerGetMetricsOpts]",
		},
		{
			name: "call to backend API fails",
			lb:   &LoadBalancer{ID: 6},
			opts: LoadBalancerGetMetricsOpts{
				Types: []LoadBalancerMetricType{LoadBalancerMetricBandwidth},
				Start: mustParseTime(t, "2017-01-01T00:00:00Z"),
				End:   mustParseTime(t, "2017-01-01T23:00:00Z"),
			},
			respStatus:  http.StatusInternalServerError,
			expectedErr: "hcloud: server responded with status code 500",
		},
		{
			name: "no load balancer passed",
			opts: LoadBalancerGetMetricsOpts{
				Types: []LoadBalancerMetricType{LoadBalancerMetricBandwidth},
				Start: mustParseTime(t, "2017-01-01T00:00:00Z"),
				End:   mustParseTime(t, "2017-01-01T23:00:00Z"),
			},
			expectedErr: "invalid argument 'loadBalancer' [*hcloud.LoadBalancer]: empty value '<nil>'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := newTestEnv()
			defer env.Teardown()

			if tt.lb != nil {
				path := fmt.Sprintf("/load_balancers/%d/metrics", tt.lb.ID)
				env.Mux.HandleFunc(path, func(rw http.ResponseWriter, r *http.Request) {
					if r.Method != "GET" {
						t.Errorf("expected GET; got %s", r.Method)
					}
					opts := loadBalancerMetricsOptsFromURL(t, r.URL)
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
			actual, _, err := env.Client.LoadBalancer.GetMetrics(ctx, tt.lb, tt.opts)
			if tt.expectedErr != "" {
				if tt.expectedErr != err.Error() {
					t.Errorf("expected err: %v; got: %v", tt.expectedErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("failed to get load balancer metrics: %v", err)
			}
			if !cmp.Equal(&tt.expected, actual) {
				t.Errorf("Actual metrics did not equal expected: %s", cmp.Diff(&tt.expected, actual))
			}
		})
	}
}

func loadBalancerMetricsOptsFromURL(t *testing.T, u *url.URL) LoadBalancerGetMetricsOpts {
	var opts LoadBalancerGetMetricsOpts

	for k, vs := range u.Query() {
		switch k {
		case "type":
			for _, v := range vs {
				opts.Types = append(opts.Types, LoadBalancerMetricType(v))
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
