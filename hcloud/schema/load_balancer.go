package schema

import "time"

type LoadBalancer struct {
	ID               int              `json:"id"`
	Name             string           `json:"name"`
	IPv4             string           `json:"ipv4"`
	IPv6             string           `json:"ipv6"`
	Location         Location         `json:"location"`
	LoadBalancerType LoadBalancerType `json:"load_balancer_type"`
	Protection       struct {
		Delete bool `json:"delete"`
	} `json:"protection"`
	Labels    map[string]string     `json:"labels"`
	Created   time.Time             `json:"created"`
	Services  []LoadBalancerService `json:"services"`
	Targets   []LoadBalancerTarget  `json:"targets"`
	Algorithm struct {
		Type string `json:"type"`
	} `json:"algorithm"`
}

type LoadBalancerService struct {
	Protocol        string                         `json:"protocol"`
	ListenPort      int                            `json:"listen_port"`
	DestinationPort int                            `json:"destination_port"`
	Proxyprotocol   bool                           `json:"proxyprotocol"`
	HTTP            *LoadBalancerServiceHTTP       `json:"http"`
	HealthCheck     LoadBalancerServiceHealthCheck `json:"health_check"`
	Status          string                         `json:"status"`
}
type LoadBalancerServiceHTTP struct {
	CookieName     string `json:"cookie_name"`
	CookieLifetime int    `json:"cookie_lifetime"`
}
type LoadBalancerServiceHealthCheck struct {
	Protocol string                              `json:"protocol"`
	Port     int                                 `json:"port"`
	Interval int                                 `json:"interval"`
	Timeout  int                                 `json:"timeout"`
	Retries  int                                 `json:"retries"`
	HTTP     *LoadBalancerServiceHealthCheckHTTP `json:"http"`
}
type LoadBalancerServiceHealthCheckHTTP struct {
	Domain string `json:"domain"`
	Path   string `json:"path"`
}

type LoadBalancerTarget struct {
	Type   string                    `json:"type"`
	Server *LoadBalancerTargetServer `json:"server"`
}
type LoadBalancerTargetServer struct {
	ID int `json:"id"`
}

// LoadBalancerListResponse defines the schema of the response when
// listing LoadBalancer.
type LoadBalancerListResponse struct {
	LoadBalancers []LoadBalancer `json:"load_balancers"`
}

// LoadBalancerGetResponse defines the schema of the response when
// retrieving a single LoadBalancer.
type LoadBalancerGetResponse struct {
	LoadBalancer LoadBalancer `json:"load_balancer"`
}

// LoadBalancerListResponse defines the schema of the response when
// listing LoadBalancer.
type LoadBalancerTargetRequest struct {
	Type   string  `json:"type"`
	Server *Server `json:"server"`
}

// VolumeActionDetachVolumeResponse defines the schema of the response when
// creating an detach volume action.
type LoadBalancerTargetResponse struct {
	Action Action `json:"action"`
}

// LoadBalancerCreateRequest defines the schema of the request to create a LoadBalancer.
type LoadBalancerCreateRequest struct {
	Name string `json:"name"`
}

// LoadBalancerCreateResponse defines the schema of the response when
// creating a LoadBalancer.
type LoadBalancerCreateResponse struct {
	LoadBalancer LoadBalancer `json:"load_balancer"`
}
