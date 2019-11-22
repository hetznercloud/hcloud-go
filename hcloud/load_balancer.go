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

// LoadBalancerType represents a LoadBalancer in the Hetzner Cloud.
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
type LoadBalancerService struct {
	Protocol        string
	ListenPort      int
	DestinationPort int
	ProxyProtocol   bool
	HTTP            *LoadBalancerServiceHTTP
	HealthCheck     LoadBalancerServiceHealthCheck
}
type LoadBalancerServiceHTTP struct {
	CookieName     string
	CookieLifeTime int
}
type LoadBalancerServiceHealthCheck struct {
	Protocol string
	Port     int
	Interval int
	Timeout  int
	Retries  int
	HTTP     *LoadBalancerServiceHealthCheckHTTP
}

type LoadBalancerServiceHealthCheckHTTP struct {
	Domain string
	Path   string
}

type LoadBalancerAlgorithm struct {
	Type string
}

type LoadBalancerTarget struct {
	Type string
	*LoadBalancerTargetServer
	*LoadBalancerTargetLabelSelector
}

type LoadBalancerTargetServer struct {
	Server Server
}

type LoadBalancerTargetLabelSelector struct {
	LabelSelector struct {
		Selector string
	}
}

// LoadBalancerProtection represents the protection level of a Load Balancer.
type LoadBalancerProtection struct {
	Delete bool
}

// LoadBalancerClient is a client for the server types API.
type LoadBalancerClient struct {
	client *Client
}

// GetByID retrieves a server type by its ID. If the server type does not exist, nil is returned.
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

// GetByName retrieves a server type by its name. If the server type does not exist, nil is returned.
func (c *LoadBalancerClient) GetByName(ctx context.Context, name string) (*LoadBalancer, *Response, error) {
	LoadBalancer, response, err := c.List(ctx, LoadBalancerListOpts{Name: name})
	if len(LoadBalancer) == 0 {
		return nil, response, err
	}
	return LoadBalancer[0], response, err
}

// Get retrieves a server type by its ID if the input can be parsed as an integer, otherwise it
// retrieves a server type by its name. If the server type does not exist, nil is returned.
func (c *LoadBalancerClient) Get(ctx context.Context, idOrName string) (*LoadBalancer, *Response, error) {
	if id, err := strconv.Atoi(idOrName); err == nil {
		return c.GetByID(ctx, int(id))
	}
	return c.GetByName(ctx, idOrName)
}

// LoadBalancerListOpts specifies options for listing server types.
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

// List returns a list of server types for a specific page.
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

// All returns all server types.
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

// AllWithOpts returns all LoadBalancers for the given options.
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

// LoadBalancerCreateOpts specifies options for creating a new LoadBalancer.
type LoadBalancerCreateOpts struct {
	Name string
}

// Validate checks if options are valid.
func (o LoadBalancerCreateOpts) Validate() error {
	if o.Name == "" {
		return errors.New("missing name")
	}

	return nil
}

// Create creates a new LoadBalancer.
func (c *LoadBalancerClient) Create(ctx context.Context, opts LoadBalancerCreateOpts) (*LoadBalancer, *Response, error) {
	if err := opts.Validate(); err != nil {
		return nil, nil, err
	}
	reqBody := schema.LoadBalancerCreateRequest{
		Name: opts.Name,
	}
	reqBodyData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, nil, err
	}
	req, err := c.client.NewRequest(ctx, "POST", "/load_balancers", bytes.NewReader(reqBodyData))
	if err != nil {
		return nil, nil, err
	}

	respBody := schema.LoadBalancerCreateResponse{}
	resp, err := c.client.Do(req, &respBody)
	if err != nil {
		return nil, resp, err
	}
	return LoadBalancerFromSchema(respBody.LoadBalancer), resp, nil
}

// Delete deletes a load balancer.
func (c *LoadBalancerClient) Delete(ctx context.Context, loadBalancer *LoadBalancer) (*Response, error) {
	req, err := c.client.NewRequest(ctx, "DELETE", fmt.Sprintf("/load_balancers/%d", loadBalancer.ID), nil)
	if err != nil {
		return nil, err
	}
	return c.client.Do(req, nil)
}

// Detach detaches a volume from a server.
func (c *LoadBalancerClient) AddServerTarget(ctx context.Context, loadBalancer *LoadBalancer, server *Server) (*Action, *Response, error) {
	// TODO: Move to generic removeTarget after POC
	reqBody := schema.LoadBalancerTargetRequest{
		Type:   "server",
		Server: &schema.Server{ID: server.ID},
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

	var respBody schema.LoadBalancerTargetResponse
	resp, err := c.client.Do(req, &respBody)
	if err != nil {
		return nil, resp, err
	}
	return ActionFromSchema(respBody.Action), resp, nil
}

// Detach detaches a volume from a server.
func (c *LoadBalancerClient) RemoveServerTarget(ctx context.Context, loadBalancer *LoadBalancer, server *Server) (*Action, *Response, error) {
	// TODO: Move to generic removeTarget after POC
	reqBody := schema.LoadBalancerTargetRequest{
		Type:   "server",
		Server: &schema.Server{ID: server.ID},
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

	var respBody schema.LoadBalancerTargetResponse
	resp, err := c.client.Do(req, &respBody)
	if err != nil {
		return nil, resp, err
	}
	return ActionFromSchema(respBody.Action), resp, nil
}
