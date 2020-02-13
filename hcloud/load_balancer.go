package hcloud

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"time"

	"github.com/hetznercloud/hcloud-go/hcloud/schema"
)

// LoadBalancer represents a Load Balancer in the Hetzner Cloud.
type LoadBalancer struct {
	ID               int
	Name             string
	IPv4             net.IP
	IPv6             net.IP
	Location         *Location
	LoadBalancerType *LoadBalancerType
	Algorithm        LoadBalancerAlgorithm
	Services         []*LoadBalancerService
	Targets          []LoadBalancerTarget
	Protection       LoadBalancerProtection
	Labels           map[string]string
	Created          time.Time
}

// LoadBalancerService represents a service of a Load Balancer
type LoadBalancerService struct {
	Protocol        LoadBalancerServiceProtocol
	ListenPort      int
	DestinationPort int
	ProxyProtocol   bool
	HTTP            *LoadBalancerServiceHTTP
	HealthCheck     LoadBalancerServiceHealthCheck
}

// LoadBalancerServiceHTTP represents HTTP specific options for a service of a Load Balancer
type LoadBalancerServiceHTTP struct {
	CookieName     string
	CookieLifeTime int
}

// LoadBalancerServiceHealthCheck represents Health Check specific options for a service of a Load Balancer
type LoadBalancerServiceHealthCheck struct {
	Protocol string
	Port     int
	Interval int
	Timeout  int
	Retries  int
	HTTP     *LoadBalancerServiceHealthCheckHTTP
}

// LoadBalancerServiceHealthCheckHTTP represents HTTP specific options for a Health Check of a Load Balancer
type LoadBalancerServiceHealthCheckHTTP struct {
	Domain string
	Path   string
}

// LoadBalancerAlgorithm represents Algorithm option of a Load Balancer
type LoadBalancerAlgorithm struct {
	Type string
}

// LoadBalancerTargetType specifies a load balancer target type.
type LoadBalancerTargetType string

const (
	// LoadBalancerTargetTypeServer is the type when a cloud server should be linked directly.
	LoadBalancerTargetTypeServer LoadBalancerTargetType = "server"
	// LoadBalancerTargetTypeLabelSelector is the type when multiple cloud server should be linked through a label selector.
	LoadBalancerTargetTypeLabelSelector LoadBalancerTargetType = "label_selector"
)

// LoadBalancerServiceProtocol specifies a load balancer service protocol.
type LoadBalancerServiceProtocol string

const (
	// LoadBalancerServiceProtocolTCP is the protocol when the Load Balancer is used as TCP Load Balancer.
	LoadBalancerServiceProtocolTCP LoadBalancerServiceProtocol = "tcp"
	// LoadBalancerServiceProtocolHTTP is the protocol when the Load Balancer is used as HTTP Load Balancer.
	LoadBalancerServiceProtocolHTTP LoadBalancerServiceProtocol = "http"
	// LoadBalancerServiceProtocolHTTPS is the protocol when the Load Balancer is used as HTTP Load Balancer with SSL Termination.
	LoadBalancerServiceProtocolHTTPS LoadBalancerServiceProtocol = "https"
)

// LoadBalancerAlgorithmType specifies a load balancer service protocol.
type LoadBalancerAlgorithmType string

const (
	// LoadBalancerAlgorithmTypeRoundRobin represents a RoundRobin algorithm.
	LoadBalancerAlgorithmTypeRoundRobin LoadBalancerAlgorithmType = "round_robin"
)

// LoadBalancerTarget represents target of a Load Balancer
type LoadBalancerTarget struct {
	Type LoadBalancerTargetType
	*LoadBalancerTargetServer
	*LoadBalancerTargetLabelSelector
	HealthStatus []LoadBalancerTargetHealthStatus
}

// LoadBalancerTargetServer represents server target of a Load Balancer
type LoadBalancerTargetServer struct {
	Server Server
}

// LoadBalancerTargetLabelSelector represents label selector target of a Load Balancer
type LoadBalancerTargetLabelSelector struct {
	LabelSelector struct {
		Selector string
	}
}

// LoadBalancerTargetHealthStatus represents target health status of a Load Balancer
type LoadBalancerTargetHealthStatus struct {
	ListenPort int
	Healthy    bool
}

// LoadBalancerProtection represents the protection level of a Load Balancer.
type LoadBalancerProtection struct {
	Delete bool
}

// LoadBalancerClient is a client for the Load Balancers API.
type LoadBalancerClient struct {
	client *Client
}

// GetByID retrieves a Load Balancer by its ID. If the Load Balancer does not exist, nil is returned.
func (c *LoadBalancerClient) GetByID(ctx context.Context, id int) (*LoadBalancer, *Response, error) {
	req, err := c.client.NewRequest(ctx, "GET", fmt.Sprintf("/load_balancers/%d", id), nil)
	if err != nil {
		return nil, nil, err
	}

	var body schema.LoadBalancerGetResponse
	resp, err := c.client.Do(req, &body)
	if err != nil {
		if IsError(err, ErrorCodeNotFound) {
			return nil, resp, nil
		}
		return nil, nil, err
	}
	return LoadBalancerFromSchema(body.LoadBalancer), resp, nil
}

// GetByName retrieves a Load Balancer by its name. If the Load Balancer does not exist, nil is returned.
func (c *LoadBalancerClient) GetByName(ctx context.Context, name string) (*LoadBalancer, *Response, error) {
	LoadBalancer, response, err := c.List(ctx, LoadBalancerListOpts{Name: name})
	if len(LoadBalancer) == 0 {
		return nil, response, err
	}
	return LoadBalancer[0], response, err
}

// Get retrieves a Load Balancer by its ID if the input can be parsed as an integer, otherwise it
// retrieves a Load Balancer by its name. If the Load Balancer does not exist, nil is returned.
func (c *LoadBalancerClient) Get(ctx context.Context, idOrName string) (*LoadBalancer, *Response, error) {
	if id, err := strconv.Atoi(idOrName); err == nil {
		return c.GetByID(ctx, int(id))
	}
	return c.GetByName(ctx, idOrName)
}

// LoadBalancerListOpts specifies options for listing Load Balancers.
type LoadBalancerListOpts struct {
	ListOpts
	Name string
}

func (l LoadBalancerListOpts) values() url.Values {
	vals := l.ListOpts.values()
	if l.Name != "" {
		vals.Add("name", l.Name)
	}
	return vals
}

// List returns a list of Load Balancers for a specific page.
func (c *LoadBalancerClient) List(ctx context.Context, opts LoadBalancerListOpts) ([]*LoadBalancer, *Response, error) {
	path := "/load_balancers?" + opts.values().Encode()
	req, err := c.client.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	var body schema.LoadBalancerListResponse
	resp, err := c.client.Do(req, &body)
	if err != nil {
		return nil, nil, err
	}
	LoadBalancers := make([]*LoadBalancer, 0, len(body.LoadBalancers))
	for _, s := range body.LoadBalancers {
		LoadBalancers = append(LoadBalancers, LoadBalancerFromSchema(s))
	}
	return LoadBalancers, resp, nil
}

// All returns all Load Balancers.
func (c *LoadBalancerClient) All(ctx context.Context) ([]*LoadBalancer, error) {
	allLoadBalancer := []*LoadBalancer{}

	opts := LoadBalancerListOpts{}
	opts.PerPage = 50

	_, err := c.client.all(func(page int) (*Response, error) {
		opts.Page = page
		LoadBalancer, resp, err := c.List(ctx, opts)
		if err != nil {
			return resp, err
		}
		allLoadBalancer = append(allLoadBalancer, LoadBalancer...)
		return resp, nil
	})
	if err != nil {
		return nil, err
	}

	return allLoadBalancer, nil
}

// AllWithOpts returns all Load Balancers for the given options.
func (c *LoadBalancerClient) AllWithOpts(ctx context.Context, opts LoadBalancerListOpts) ([]*LoadBalancer, error) {
	var allLoadBalancers []*LoadBalancer

	_, err := c.client.all(func(page int) (*Response, error) {
		opts.Page = page
		LoadBalancers, resp, err := c.List(ctx, opts)
		if err != nil {
			return resp, err
		}
		allLoadBalancers = append(allLoadBalancers, LoadBalancers...)
		return resp, nil
	})
	if err != nil {
		return nil, err
	}

	return allLoadBalancers, nil
}

// LoadBalancerUpdateOpts specifies options for updating a Load Balancer.
type LoadBalancerUpdateOpts struct {
	Name   string
	Labels map[string]string
}

// Update updates a Load Balancer.
func (c *LoadBalancerClient) Update(ctx context.Context, loadBalancer *LoadBalancer, opts LoadBalancerUpdateOpts) (*LoadBalancer, *Response, error) {
	reqBody := schema.LoadBalancerUpdateRequest{
		Name: opts.Name,
	}
	if opts.Labels != nil {
		reqBody.Labels = &opts.Labels
	}
	reqBodyData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, nil, err
	}

	path := fmt.Sprintf("/load_balancers/%d", loadBalancer.ID)
	req, err := c.client.NewRequest(ctx, "PUT", path, bytes.NewReader(reqBodyData))
	if err != nil {
		return nil, nil, err
	}

	respBody := schema.LoadBalancerUpdateResponse{}
	resp, err := c.client.Do(req, &respBody)
	if err != nil {
		return nil, resp, err
	}
	return LoadBalancerFromSchema(respBody.LoadBalancer), resp, nil
}

// LoadBalancerCreateOpts specifies options for creating a new Load Balancer.
type LoadBalancerCreateOpts struct {
	Name             string
	LoadBalancerType *LoadBalancerType
	Algorithm        LoadBalancerAlgorithm
	Location         *Location
	NetworkZone      string
}

// Validate checks if options are valid.
func (o LoadBalancerCreateOpts) Validate() error {
	if o.Name == "" {
		return errors.New("missing name")
	}
	if o.LoadBalancerType == nil || (o.LoadBalancerType.ID == 0 && o.LoadBalancerType.Name == "") {
		return errors.New("missing load balancer type")
	}
	if o.Algorithm.Type == "" {
		return errors.New("missing algorithm type")
	}
	if o.Location != nil && o.NetworkZone != "" {
		return errors.New("location and loadBalancer_zone are mutually exclusive")
	}
	return nil
}

// LoadBalancerCreateResult is the result of a create Load Balancer call.
type LoadBalancerCreateResult struct {
	LoadBalancer *LoadBalancer
	Action       *Action
}

// Create creates a new Load Balancer
func (c *LoadBalancerClient) Create(ctx context.Context, opts LoadBalancerCreateOpts) (LoadBalancerCreateResult, *Response, error) {
	if err := opts.Validate(); err != nil {
		return LoadBalancerCreateResult{}, nil, err
	}
	reqBody := schema.LoadBalancerCreateRequest{
		Name: opts.Name,
		Algorithm: schema.LoadBalancerAlgorithm{
			Type: opts.Algorithm.Type,
		},
	}
	if opts.LoadBalancerType.ID != 0 {
		reqBody.LoadBalancerType = opts.LoadBalancerType.ID
	} else if opts.LoadBalancerType.Name != "" {
		reqBody.LoadBalancerType = opts.LoadBalancerType.Name
	}

	if opts.Location != nil {
		if opts.Location.ID != 0 {
			reqBody.Location = strconv.Itoa(opts.Location.ID)
		} else {
			reqBody.Location = opts.Location.Name
		}
	}
	if opts.NetworkZone != "" {
		reqBody.NetworkZone = opts.NetworkZone
	}
	reqBodyData, err := json.Marshal(reqBody)
	if err != nil {
		return LoadBalancerCreateResult{}, nil, err
	}
	req, err := c.client.NewRequest(ctx, "POST", "/load_balancers", bytes.NewReader(reqBodyData))
	if err != nil {
		return LoadBalancerCreateResult{}, nil, err
	}

	respBody := schema.LoadBalancerCreateResponse{}
	resp, err := c.client.Do(req, &respBody)
	if err != nil {
		return LoadBalancerCreateResult{}, resp, err
	}
	return LoadBalancerCreateResult{
		LoadBalancer: LoadBalancerFromSchema(respBody.LoadBalancer),
		Action:       ActionFromSchema(respBody.Action),
	}, resp, nil
}

// Delete deletes a Load Balancer.
func (c *LoadBalancerClient) Delete(ctx context.Context, loadBalancer *LoadBalancer) (*Response, error) {
	req, err := c.client.NewRequest(ctx, "DELETE", fmt.Sprintf("/load_balancers/%d", loadBalancer.ID), nil)
	if err != nil {
		return nil, err
	}
	return c.client.Do(req, nil)
}

// LoadBalancerTargetOpts specifies options for adding or removing targets from a Load Balancer.
type LoadBalancerTargetOpts struct {
	Server        *Server
	LabelSelector string
}

// Validate checks if options are valid.
func (o LoadBalancerTargetOpts) Validate() error {
	if o.Server == nil && o.LabelSelector == "" {
		return errors.New("missing server or label_selector")
	} else if o.Server != nil && o.LabelSelector != "" {
		return errors.New("server and label_selector are mutually exclusive")
	}
	return nil
}

// GetType returns the Target type based on the given options
func (o LoadBalancerTargetOpts) getType() LoadBalancerTargetType {
	if o.LabelSelector != "" {
		return LoadBalancerTargetTypeLabelSelector
	}
	return LoadBalancerTargetTypeServer

}

// AddTarget adds a target to a Load Balancer.
func (c *LoadBalancerClient) AddTarget(ctx context.Context, loadBalancer *LoadBalancer, opts LoadBalancerTargetOpts) (*Action, *Response, error) {
	if err := opts.Validate(); err != nil {
		return nil, nil, err
	}
	targetType := opts.getType()
	reqBody := schema.LoadBalancerActionTargetRequest{
		Type: string(targetType),
	}
	switch targetType {
	case LoadBalancerTargetTypeServer:
		reqBody.Server = &schema.LoadBalancerTargetServer{
			ID: opts.Server.ID,
		}
	case LoadBalancerTargetTypeLabelSelector:
		reqBody.LabelSelector = &schema.LoadBalancerTargetLabelSelector{
			Selector: opts.LabelSelector,
		}
	}
	reqBodyData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, nil, err
	}

	path := fmt.Sprintf("/load_balancers/%d/actions/add_target", loadBalancer.ID)
	req, err := c.client.NewRequest(ctx, "POST", path, bytes.NewReader(reqBodyData))
	if err != nil {
		return nil, nil, err
	}

	var respBody schema.LoadBalancerActionTargetResponse
	resp, err := c.client.Do(req, &respBody)
	if err != nil {
		return nil, resp, err
	}
	return ActionFromSchema(respBody.Action), resp, nil
}

// RemoveTarget removes a target from a Load Balancer.
func (c *LoadBalancerClient) RemoveTarget(ctx context.Context, loadBalancer *LoadBalancer, opts LoadBalancerTargetOpts) (*Action, *Response, error) {
	if err := opts.Validate(); err != nil {
		return nil, nil, err
	}
	targetType := opts.getType()
	reqBody := schema.LoadBalancerActionTargetRequest{
		Type: string(targetType),
	}
	switch targetType {
	case LoadBalancerTargetTypeServer:
		reqBody.Server = &schema.LoadBalancerTargetServer{
			ID: opts.Server.ID,
		}
	case LoadBalancerTargetTypeLabelSelector:
		reqBody.LabelSelector = &schema.LoadBalancerTargetLabelSelector{
			Selector: opts.LabelSelector,
		}
	}
	reqBodyData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, nil, err
	}

	path := fmt.Sprintf("/load_balancers/%d/actions/remove_target", loadBalancer.ID)
	req, err := c.client.NewRequest(ctx, "POST", path, bytes.NewReader(reqBodyData))
	if err != nil {
		return nil, nil, err
	}

	var respBody schema.LoadBalancerActionTargetResponse
	resp, err := c.client.Do(req, &respBody)
	if err != nil {
		return nil, resp, err
	}
	return ActionFromSchema(respBody.Action), resp, nil
}

// LoadBalancerAddServiceOpts specifies options for adding service to a Load Balancer.
type LoadBalancerAddServiceOpts struct {
	Protocol        LoadBalancerServiceProtocol
	ListenPort      int
	DestinationPort int
	ProxyProtocol   *bool
	HTTP            *LoadBalancerServiceHTTP
	HealthCheck     *LoadBalancerServiceHealthCheck
}

// AddService adds a service to a Load Balancer.
func (c *LoadBalancerClient) AddService(ctx context.Context, loadBalancer *LoadBalancer, opts LoadBalancerAddServiceOpts) (*Action, *Response, error) {
	reqBody := schema.LoadBalancerActionAddServiceRequest{
		Protocol:        string(opts.Protocol),
		ListenPort:      opts.ListenPort,
		DestinationPort: opts.DestinationPort,
		ProxyProtocol:   opts.ProxyProtocol,
	}

	if opts.HTTP != nil {
		reqBody.HTTP = &schema.LoadBalancerServiceHTTP{
			CookieName:     opts.HTTP.CookieName,
			CookieLifetime: opts.HTTP.CookieLifeTime,
		}
	}

	if opts.HealthCheck != nil {
		reqBody.HealthCheck = &schema.LoadBalancerServiceHealthCheck{
			Protocol: opts.HealthCheck.Protocol,
			Port:     opts.HealthCheck.Port,
			Interval: opts.HealthCheck.Interval,
			Timeout:  opts.HealthCheck.Timeout,
			Retries:  opts.HealthCheck.Retries,
		}
		if opts.HealthCheck.HTTP != nil {
			reqBody.HealthCheck.HTTP = &schema.LoadBalancerServiceHealthCheckHTTP{
				Domain: opts.HealthCheck.HTTP.Domain,
				Path:   opts.HealthCheck.HTTP.Path,
			}
		}
	}

	reqBodyData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, nil, err
	}

	path := fmt.Sprintf("/load_balancers/%d/actions/add_service", loadBalancer.ID)
	req, err := c.client.NewRequest(ctx, "POST", path, bytes.NewReader(reqBodyData))
	if err != nil {
		return nil, nil, err
	}

	var respBody schema.LoadBalancerActionAddServiceResponse
	resp, err := c.client.Do(req, &respBody)
	if err != nil {
		return nil, resp, err
	}
	return ActionFromSchema(respBody.Action), resp, nil
}

// DeleteService removes a server from a Load Balancer.
func (c *LoadBalancerClient) DeleteService(ctx context.Context, loadBalancer *LoadBalancer, listenPort int) (*Action, *Response, error) {
	reqBody := schema.LoadBalancerDeleteServiceRequest{
		ListenPort: listenPort,
	}
	reqBodyData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, nil, err
	}

	path := fmt.Sprintf("/load_balancers/%d/actions/delete_service", loadBalancer.ID)
	req, err := c.client.NewRequest(ctx, "POST", path, bytes.NewReader(reqBodyData))
	if err != nil {
		return nil, nil, err
	}

	var respBody schema.LoadBalancerDeleteServiceResponse
	resp, err := c.client.Do(req, &respBody)
	if err != nil {
		return nil, resp, err
	}
	return ActionFromSchema(respBody.Action), resp, nil
}

// LoadBalancerChangeProtectionOpts specifies options for changing the resource protection level of a Load Balancer.
type LoadBalancerChangeProtectionOpts struct {
	Delete *bool
}

// ChangeProtection changes the resource protection level of a Load Balancer.
func (c *LoadBalancerClient) ChangeProtection(ctx context.Context, loadBalancer *LoadBalancer, opts LoadBalancerChangeProtectionOpts) (*Action, *Response, error) {
	reqBody := schema.LoadBalancerActionChangeProtectionRequest{
		Delete: opts.Delete,
	}
	reqBodyData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, nil, err
	}

	path := fmt.Sprintf("/load_balancers/%d/actions/change_protection", loadBalancer.ID)
	req, err := c.client.NewRequest(ctx, "POST", path, bytes.NewReader(reqBodyData))
	if err != nil {
		return nil, nil, err
	}

	respBody := schema.LoadBalancerActionChangeProtectionResponse{}
	resp, err := c.client.Do(req, &respBody)
	if err != nil {
		return nil, resp, err
	}
	return ActionFromSchema(respBody.Action), resp, err
}

// LoadBalancerChangeAlgorithmOpts specifies options for changing the algorithm of a Load Balancer.
type LoadBalancerChangeAlgorithmOpts struct {
	Type LoadBalancerAlgorithmType
}

// ChangeAlgorithm changes the algorithm of a Load Balancer.
func (c *LoadBalancerClient) ChangeAlgorithm(ctx context.Context, loadBalancer *LoadBalancer, opts LoadBalancerChangeAlgorithmOpts) (*Action, *Response, error) {
	reqBody := schema.LoadBalancerActionChangeAlgorithmRequest{
		Type: string(opts.Type),
	}
	reqBodyData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, nil, err
	}

	path := fmt.Sprintf("/load_balancers/%d/actions/change_algorithm", loadBalancer.ID)
	req, err := c.client.NewRequest(ctx, "POST", path, bytes.NewReader(reqBodyData))
	if err != nil {
		return nil, nil, err
	}

	respBody := schema.LoadBalancerActionChangeAlgorithmResponse{}
	resp, err := c.client.Do(req, &respBody)
	if err != nil {
		return nil, resp, err
	}
	return ActionFromSchema(respBody.Action), resp, err
}

// LoadBalancerUpdateHealthCheckOpts specifies options for updating a health check of a service from a Load Balancer.
type LoadBalancerUpdateHealthCheckOpts struct {
	ListenPort  int
	HealthCheck LoadBalancerServiceHealthCheck
}

// UpdateHealthCheck updates the health check from a service.
func (c *LoadBalancerClient) UpdateHealthCheck(ctx context.Context, loadBalancer *LoadBalancer, opts LoadBalancerUpdateHealthCheckOpts) (*Action, *Response, error) {
	reqBody := schema.LoadBalancerActionUpdateHealthCheckRequest{
		ListenPort: opts.ListenPort,
		HealthCheck: schema.LoadBalancerServiceHealthCheck{
			Protocol: opts.HealthCheck.Protocol,
			Port:     opts.HealthCheck.Port,
			Interval: opts.HealthCheck.Interval,
			Timeout:  opts.HealthCheck.Timeout,
			Retries:  opts.HealthCheck.Retries,
		},
	}

	if opts.HealthCheck.HTTP != nil {
		reqBody.HealthCheck.HTTP = &schema.LoadBalancerServiceHealthCheckHTTP{
			Domain: opts.HealthCheck.HTTP.Domain,
			Path:   opts.HealthCheck.HTTP.Path,
		}
	}
	reqBodyData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, nil, err
	}

	path := fmt.Sprintf("/load_balancers/%d/actions/update_health_check", loadBalancer.ID)
	req, err := c.client.NewRequest(ctx, "POST", path, bytes.NewReader(reqBodyData))
	if err != nil {
		return nil, nil, err
	}

	respBody := schema.LoadBalancerActionChangeAlgorithmResponse{}
	resp, err := c.client.Do(req, &respBody)
	if err != nil {
		return nil, resp, err
	}
	return ActionFromSchema(respBody.Action), resp, err
}
