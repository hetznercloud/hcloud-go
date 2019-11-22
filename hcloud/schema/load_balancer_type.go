package schema

// LoadBalancerType defines the schema of a LoadBalancer type.
type LoadBalancerType struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Requests    int    `json:"requests"`
	Services    int    `json:"services"`
}

// LoadBalancerTypeListResponse defines the schema of the response when
// listing LoadBalancer types.
type LoadBalancerTypeListResponse struct {
	LoadBalancerTypes []LoadBalancerType `json:"load_balancer_types"`
}

// LoadBalancerTypeGetResponse defines the schema of the response when
// retrieving a single LoadBalancer type.
type LoadBalancerTypeGetResponse struct {
	LoadBalancerType LoadBalancerType `json:"load_balancer_type"`
}
