package schema

import "time"

type LoadBalancer struct {
	ID               int                    `json:"id"`
	Name             string                 `json:"name"`
	IPv4             string                 `json:"ipv4"`
	IPv6             string                 `json:"ipv6"`
	Location         Location               `json:"location"`
	LoadBalancerType LoadBalancerType       `json:"load_balancer_type"`
	Protection       LoadBalancerProtection `json:"protection"`
	Labels           map[string]string      `json:"labels"`
	Created          time.Time              `json:"created"`
	Services         []LoadBalancerService  `json:"services"`
	Targets          []LoadBalancerTarget   `json:"targets"`
	Algorithm        struct {
		Type string `json:"type"`
	} `json:"algorithm"`
}

// LoadBalancerProtection represents the protection level of a load balancer.
type LoadBalancerProtection struct {
	Delete bool `json:"delete"`
}

// LoadBalancerService represents a service of a load balancer.
type LoadBalancerService struct {
	Protocol        string                         `json:"protocol"`
	ListenPort      int                            `json:"listen_port"`
	DestinationPort int                            `json:"destination_port"`
	Proxyprotocol   bool                           `json:"proxyprotocol"`
	HTTP            *LoadBalancerServiceHTTP       `json:"http"`
	HealthCheck     LoadBalancerServiceHealthCheck `json:"health_check"`
	Status          string                         `json:"status"`
}

// LoadBalancerServiceHTTP represents the http configuration for a LoadBalancerService
type LoadBalancerServiceHTTP struct {
	CookieName     string `json:"cookie_name"`
	CookieLifetime int    `json:"cookie_lifetime"`
}

// LoadBalancerServiceHealthCheck represents
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
	Type          string                           `json:"type"`
	Server        *LoadBalancerTargetServer        `json:"server"`
	LabelSelector *LoadBalancerTargetLabelSelector `json:"label_selector"`
}
type LoadBalancerTargetServer struct {
	ID int `json:"id"`
}
type LoadBalancerTargetLabelSelector struct {
	Selector string `json:"selector"`
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
	Type          string                           `json:"type"`
	Server        *LoadBalancerTargetServer        `json:"server"`
	LabelSelector *LoadBalancerTargetLabelSelector `json:"label_selector,omitempty"`
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

// LoadBalancerActionChangeProtectionRequest defines the schema of the request to
// change the resource protection of a load balancer.
type LoadBalancerActionChangeProtectionRequest struct {
	Delete *bool `json:"delete,omitempty"`
}

// LoadBalancerActionChangeProtectionResponse defines the schema of the response when
// changing the resource protection of a load balancer.
type LoadBalancerActionChangeProtectionResponse struct {
	Action Action `json:"action"`
}
