package mock

import (
	hcloud "github.com/hetznercloud/hcloud-go/v2/hcloud"
	gomock "go.uber.org/mock/gomock"
)

type MockClient struct {
	Action           *MockActionClient
	Certificate      *MockCertificateClient
	Datacenter       *MockDatacenterClient
	Firewall         *MockFirewallClient
	FloatingIP       *MockFloatingIPClient
	Image            *MockImageClient
	ISO              *MockISOClient
	LoadBalancer     *MockLoadBalancerClient
	LoadBalancerType *MockLoadBalancerTypeClient
	Location         *MockLocationClient
	Network          *MockNetworkClient
	PlacementGroup   *MockPlacementGroupClient
	PrimaryIP        *MockPrimaryIPClient
	RDNS             *MockRDNSClient
	Server           *MockServerClient
	ServerType       *MockServerTypeClient
	SSHKey           *MockSSHKeyClient
	Volume           *MockVolumeClient
}

func NewMockClient(ctrl *gomock.Controller) *MockClient {
	return &MockClient{
		Action:           NewMockActionClient(ctrl),
		Certificate:      NewMockCertificateClient(ctrl),
		Datacenter:       NewMockDatacenterClient(ctrl),
		Firewall:         NewMockFirewallClient(ctrl),
		FloatingIP:       NewMockFloatingIPClient(ctrl),
		Image:            NewMockImageClient(ctrl),
		ISO:              NewMockISOClient(ctrl),
		LoadBalancer:     NewMockLoadBalancerClient(ctrl),
		LoadBalancerType: NewMockLoadBalancerTypeClient(ctrl),
		Location:         NewMockLocationClient(ctrl),
		Network:          NewMockNetworkClient(ctrl),
		PlacementGroup:   NewMockPlacementGroupClient(ctrl),
		PrimaryIP:        NewMockPrimaryIPClient(ctrl),
		RDNS:             NewMockRDNSClient(ctrl),
		Server:           NewMockServerClient(ctrl),
		ServerType:       NewMockServerTypeClient(ctrl),
		SSHKey:           NewMockSSHKeyClient(ctrl),
		Volume:           NewMockVolumeClient(ctrl),
	}
}

func (*MockClient) WithOpts(...hcloud.ClientOption) {}
