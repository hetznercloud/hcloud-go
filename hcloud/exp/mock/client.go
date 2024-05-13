package mock

import (
	gomock "go.uber.org/mock/gomock"

	hcloud "github.com/hetznercloud/hcloud-go/v2/hcloud"
)

type Client struct {
	Action           *ActionClient
	Certificate      *CertificateClient
	Datacenter       *DatacenterClient
	Firewall         *FirewallClient
	FloatingIP       *FloatingIPClient
	Image            *ImageClient
	ISO              *ISOClient
	LoadBalancer     *LoadBalancerClient
	LoadBalancerType *LoadBalancerTypeClient
	Location         *LocationClient
	Network          *NetworkClient
	PlacementGroup   *PlacementGroupClient
	PrimaryIP        *PrimaryIPClient
	RDNS             *RDNSClient
	Server           *ServerClient
	ServerType       *ServerTypeClient
	SSHKey           *SSHKeyClient
	Volume           *VolumeClient
}

func NewMockClient(ctrl *gomock.Controller) *Client {
	return &Client{
		Action:           NewActionClient(ctrl),
		Certificate:      NewCertificateClient(ctrl),
		Datacenter:       NewDatacenterClient(ctrl),
		Firewall:         NewFirewallClient(ctrl),
		FloatingIP:       NewFloatingIPClient(ctrl),
		Image:            NewImageClient(ctrl),
		ISO:              NewISOClient(ctrl),
		LoadBalancer:     NewLoadBalancerClient(ctrl),
		LoadBalancerType: NewLoadBalancerTypeClient(ctrl),
		Location:         NewLocationClient(ctrl),
		Network:          NewNetworkClient(ctrl),
		PlacementGroup:   NewPlacementGroupClient(ctrl),
		PrimaryIP:        NewPrimaryIPClient(ctrl),
		RDNS:             NewRDNSClient(ctrl),
		Server:           NewServerClient(ctrl),
		ServerType:       NewServerTypeClient(ctrl),
		SSHKey:           NewSSHKeyClient(ctrl),
		Volume:           NewVolumeClient(ctrl),
	}
}

func (*Client) WithOpts(...hcloud.ClientOption) {}
