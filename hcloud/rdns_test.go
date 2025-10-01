package hcloud

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"testing"

	"github.com/hetznercloud/hcloud-go/v2/hcloud/schema"
)

func TestChangeDNSPtr(t *testing.T) {
	var (
		ctx          = context.Background()
		server       = &Server{ID: 1}
		loadBalancer = &LoadBalancer{ID: 1}
		floatingIP   = &FloatingIP{ID: 1}
		dns          = "example.com"
	)

	tests := []struct {
		name       string
		apiURL     string
		IP         net.IP
		DNS        *string
		changeFunc func(env *testEnv, ip net.IP, dns *string) (*Action, *Response, error)
	}{
		{
			name:   "set via server client",
			apiURL: "/servers/1/actions/change_dns_ptr",
			IP:     net.ParseIP("127.0.0.1"),
			DNS:    &dns,
			changeFunc: func(env *testEnv, ip net.IP, dns *string) (*Action, *Response, error) {
				return env.Client.Server.ChangeDNSPtr(ctx, server, ip.String(), dns)
			},
		},
		{
			name:   "reset via server client",
			apiURL: "/servers/1/actions/change_dns_ptr",
			IP:     net.ParseIP("127.0.0.1"),
			DNS:    nil,
			changeFunc: func(env *testEnv, ip net.IP, dns *string) (*Action, *Response, error) {
				return env.Client.Server.ChangeDNSPtr(ctx, server, ip.String(), dns)
			},
		},
		{
			name:   "set server via rdns client",
			apiURL: "/servers/1/actions/change_dns_ptr",
			IP:     net.ParseIP("127.0.0.1"),
			DNS:    &dns,
			changeFunc: func(env *testEnv, ip net.IP, dns *string) (*Action, *Response, error) {
				return env.Client.RDNS.ChangeDNSPtr(ctx, server, ip, dns)
			},
		},
		{
			name:   "reset server via rdns client",
			apiURL: "/servers/1/actions/change_dns_ptr",
			IP:     net.ParseIP("127.0.0.1"),
			DNS:    nil,
			changeFunc: func(env *testEnv, ip net.IP, dns *string) (*Action, *Response, error) {
				return env.Client.RDNS.ChangeDNSPtr(ctx, server, ip, dns)
			},
		},
		{
			name:   "set via load balancer client",
			apiURL: "/load_balancers/1/actions/change_dns_ptr",
			IP:     net.ParseIP("127.0.0.1"),
			DNS:    &dns,
			changeFunc: func(env *testEnv, ip net.IP, dns *string) (*Action, *Response, error) {
				return env.Client.LoadBalancer.ChangeDNSPtr(ctx, loadBalancer, ip.String(), dns)
			},
		},
		{
			name:   "reset via load balancer client",
			apiURL: "/load_balancers/1/actions/change_dns_ptr",
			IP:     net.ParseIP("127.0.0.1"),
			DNS:    nil,
			changeFunc: func(env *testEnv, ip net.IP, dns *string) (*Action, *Response, error) {
				return env.Client.LoadBalancer.ChangeDNSPtr(ctx, loadBalancer, ip.String(), dns)
			},
		},
		{
			name:   "set load balancer via rdns client",
			apiURL: "/load_balancers/1/actions/change_dns_ptr",
			IP:     net.ParseIP("127.0.0.1"),
			DNS:    &dns,
			changeFunc: func(env *testEnv, ip net.IP, dns *string) (*Action, *Response, error) {
				return env.Client.RDNS.ChangeDNSPtr(ctx, loadBalancer, ip, dns)
			},
		},
		{
			name:   "reset load balancer via rdns client",
			apiURL: "/load_balancers/1/actions/change_dns_ptr",
			IP:     net.ParseIP("127.0.0.1"),
			DNS:    nil,
			changeFunc: func(env *testEnv, ip net.IP, dns *string) (*Action, *Response, error) {
				return env.Client.RDNS.ChangeDNSPtr(ctx, loadBalancer, ip, dns)
			},
		},
		{
			name:   "set via floating ip client",
			apiURL: "/floating_ips/1/actions/change_dns_ptr",
			IP:     net.ParseIP("127.0.0.1"),
			DNS:    &dns,
			changeFunc: func(env *testEnv, ip net.IP, dns *string) (*Action, *Response, error) {
				return env.Client.FloatingIP.ChangeDNSPtr(ctx, floatingIP, ip.String(), dns)
			},
		},
		{
			name:   "reset via floating ip client",
			apiURL: "/floating_ips/1/actions/change_dns_ptr",
			IP:     net.ParseIP("127.0.0.1"),
			DNS:    nil,
			changeFunc: func(env *testEnv, ip net.IP, dns *string) (*Action, *Response, error) {
				return env.Client.FloatingIP.ChangeDNSPtr(ctx, floatingIP, ip.String(), dns)
			},
		},
		{
			name:   "set floating ip via rdns client",
			apiURL: "/floating_ips/1/actions/change_dns_ptr",
			IP:     net.ParseIP("127.0.0.1"),
			DNS:    &dns,
			changeFunc: func(env *testEnv, ip net.IP, dns *string) (*Action, *Response, error) {
				return env.Client.RDNS.ChangeDNSPtr(ctx, floatingIP, ip, dns)
			},
		},
		{
			name:   "reset floating ip via rdns client",
			apiURL: "/floating_ips/1/actions/change_dns_ptr",
			IP:     net.ParseIP("127.0.0.1"),
			DNS:    nil,
			changeFunc: func(env *testEnv, ip net.IP, dns *string) (*Action, *Response, error) {
				return env.Client.RDNS.ChangeDNSPtr(ctx, floatingIP, ip, dns)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := newTestEnv()
			defer env.Teardown()

			env.Mux.HandleFunc(tt.apiURL, func(w http.ResponseWriter, r *http.Request) {
				var reqBody schema.ServerActionChangeDNSPtrRequest
				if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
					t.Fatal(err)
				}
				if reqBody.IP != tt.IP.String() {
					t.Errorf("unexpected IP: %v", reqBody.IP)
				}
				if reqBody.DNSPtr != tt.DNS && *reqBody.DNSPtr != *tt.DNS {
					t.Errorf("unexpected DNS ptr: %v", reqBody.DNSPtr)
				}
				json.NewEncoder(w).Encode(schema.ServerActionChangeDNSPtrResponse{
					Action: schema.Action{
						ID: 1,
					},
				})
			})

			action, _, err := tt.changeFunc(&env, tt.IP, tt.DNS)
			if err != nil {
				t.Fatal(err)
			}
			if action.ID != 1 {
				t.Errorf("unexpected action ID: %d", action.ID)
			}
		})
	}
}

func TestDNSPtrFromIP(t *testing.T) {
	var (
		server = Server{
			ID: 1,
			PublicNet: ServerPublicNet{
				IPv4: ServerPublicNetIPv4{
					IP:     net.ParseIP("127.0.0.1"),
					DNSPtr: "ipv4.example.com",
				},
				IPv6: ServerPublicNetIPv6{
					DNSPtr: map[string]string{
						"::1": "ipv6.example.com",
					},
				},
			},
		}
		loadBalancer = LoadBalancer{
			ID: 1,
			PublicNet: LoadBalancerPublicNet{
				IPv4: LoadBalancerPublicNetIPv4{
					IP:     net.ParseIP("127.0.0.1"),
					DNSPtr: "ipv4.example.com",
				},
				IPv6: LoadBalancerPublicNetIPv6{
					IP:     net.ParseIP("::1"),
					DNSPtr: "ipv6.example.com",
				},
			},
		}
		floatingIPv4 = FloatingIP{
			ID: 1,
			IP: net.ParseIP("127.0.0.1"),
			DNSPtr: map[string]string{
				"127.0.0.1": "ipv4.example.com",
			},
		}
		floatintIPv6 = FloatingIP{
			ID: 1,
			IP: net.ParseIP("::1"),
			DNSPtr: map[string]string{
				"::1": "ipv6.example.com",
			},
		}
	)

	tests := []struct {
		name string
		IP   net.IP
		DNS  string
		rdns RDNSSupporter
	}{
		{
			name: "server get dns ptr of IPv4",
			IP:   net.ParseIP("127.0.0.1"),
			DNS:  "ipv4.example.com",
			rdns: &server,
		},
		{
			name: "server get dns ptr of IPv6",
			IP:   net.ParseIP("::1"),
			DNS:  "ipv6.example.com",
			rdns: &server,
		},
		{
			name: "load balancer get dns ptr of IPv4",
			IP:   net.ParseIP("127.0.0.1"),
			DNS:  "ipv4.example.com",
			rdns: &loadBalancer,
		},
		{
			name: "load balancer get dns ptr of IPv6",
			IP:   net.ParseIP("::1"),
			DNS:  "ipv6.example.com",
			rdns: &loadBalancer,
		},
		{
			name: "floating ip get dns ptr of IPv4",
			IP:   net.ParseIP("127.0.0.1"),
			DNS:  "ipv4.example.com",
			rdns: &floatingIPv4,
		},
		{
			name: "floating ip get dns ptr of IPv6",
			IP:   net.ParseIP("::1"),
			DNS:  "ipv6.example.com",
			rdns: &floatintIPv6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			receivedDNS, err := tt.rdns.GetDNSPtrForIP(tt.IP)
			if err != nil {
				t.Fatal(err)
			}

			if tt.DNS != receivedDNS {
				t.Errorf("unexpected dns for ip %s: %s", tt.IP.String(), receivedDNS)
			}
		})
	}
}
