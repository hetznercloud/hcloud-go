package mock

//go:generate go run go.uber.org/mock/mockgen -package mock -destination zz_action_client_mock.go -mock_names IActionClient=MockActionClient github.com/hetznercloud/hcloud-go/v2/hcloud IActionClient
//go:generate go run go.uber.org/mock/mockgen -package mock -destination zz_certificate_client_mock.go -mock_names ICertificateClient=MockCertificateClient github.com/hetznercloud/hcloud-go/v2/hcloud ICertificateClient
//go:generate go run go.uber.org/mock/mockgen -package mock -destination zz_datacenter_client_mock.go -mock_names IDatacenterClient=MockDatacenterClient github.com/hetznercloud/hcloud-go/v2/hcloud IDatacenterClient
//go:generate go run go.uber.org/mock/mockgen -package mock -destination zz_firewall_client_mock.go -mock_names IFirewallClient=MockFirewallClient github.com/hetznercloud/hcloud-go/v2/hcloud IFirewallClient
//go:generate go run go.uber.org/mock/mockgen -package mock -destination zz_floating_ip_client_mock.go -mock_names IFloatingIPClient=MockFloatingIPClient github.com/hetznercloud/hcloud-go/v2/hcloud IFloatingIPClient
//go:generate go run go.uber.org/mock/mockgen -package mock -destination zz_image_client_mock.go -mock_names IImageClient=MockImageClient github.com/hetznercloud/hcloud-go/v2/hcloud IImageClient
//go:generate go run go.uber.org/mock/mockgen -package mock -destination zz_iso_client_mock.go -mock_names IISOClient=MockISOClient github.com/hetznercloud/hcloud-go/v2/hcloud IISOClient
//go:generate go run go.uber.org/mock/mockgen -package mock -destination zz_loadbalancer_client_mock.go -mock_names ILoadBalancerClient=MockLoadBalancerClient github.com/hetznercloud/hcloud-go/v2/hcloud ILoadBalancerClient
//go:generate go run go.uber.org/mock/mockgen -package mock -destination zz_loadbalancer_type_client_mock.go -mock_names ILoadBalancerTypeClient=MockLoadBalancerTypeClient github.com/hetznercloud/hcloud-go/v2/hcloud ILoadBalancerTypeClient
//go:generate go run go.uber.org/mock/mockgen -package mock -destination zz_location_client_mock.go -mock_names ILocationClient=MockLocationClient github.com/hetznercloud/hcloud-go/v2/hcloud ILocationClient
//go:generate go run go.uber.org/mock/mockgen -package mock -destination zz_network_client_mock.go -mock_names INetworkClient=MockNetworkClient github.com/hetznercloud/hcloud-go/v2/hcloud INetworkClient
//go:generate go run go.uber.org/mock/mockgen -package mock -destination zz_placement_group_client_mock.go -mock_names IPlacementGroupClient=MockPlacementGroupClient github.com/hetznercloud/hcloud-go/v2/hcloud IPlacementGroupClient
//go:generate go run go.uber.org/mock/mockgen -package mock -destination zz_primary_ip_client_mock.go -mock_names IPrimaryIPClient=MockPrimaryIPClient github.com/hetznercloud/hcloud-go/v2/hcloud IPrimaryIPClient
//go:generate go run go.uber.org/mock/mockgen -package mock -destination zz_rdns_client_mock.go -mock_names IRDNSClient=MockRDNSClient github.com/hetznercloud/hcloud-go/v2/hcloud IRDNSClient
//go:generate go run go.uber.org/mock/mockgen -package mock -destination zz_server_client_mock.go -mock_names IServerClient=MockServerClient github.com/hetznercloud/hcloud-go/v2/hcloud IServerClient
//go:generate go run go.uber.org/mock/mockgen -package mock -destination zz_server_type_client_mock.go -mock_names IServerTypeClient=MockServerTypeClient github.com/hetznercloud/hcloud-go/v2/hcloud IServerTypeClient
//go:generate go run go.uber.org/mock/mockgen -package mock -destination zz_ssh_key_client_mock.go -mock_names ISSHKeyClient=MockSSHKeyClient github.com/hetznercloud/hcloud-go/v2/hcloud ISSHKeyClient
//go:generate go run go.uber.org/mock/mockgen -package mock -destination zz_volume_client_mock.go -mock_names IVolumeClient=MockVolumeClient github.com/hetznercloud/hcloud-go/v2/hcloud IVolumeClient
