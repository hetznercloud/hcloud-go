package hcloudmock

//go:generate go run go.uber.org/mock/mockgen -package hcloudmock -destination zz_action_client_mock.go -mock_names IActionClient=ActionClient github.com/hetznercloud/hcloud-go/v2/hcloud IActionClient
//go:generate go run go.uber.org/mock/mockgen -package hcloudmock -destination zz_certificate_client_mock.go -mock_names ICertificateClient=CertificateClient github.com/hetznercloud/hcloud-go/v2/hcloud ICertificateClient
//go:generate go run go.uber.org/mock/mockgen -package hcloudmock -destination zz_datacenter_client_mock.go -mock_names IDatacenterClient=DatacenterClient github.com/hetznercloud/hcloud-go/v2/hcloud IDatacenterClient
//go:generate go run go.uber.org/mock/mockgen -package hcloudmock -destination zz_firewall_client_mock.go -mock_names IFirewallClient=FirewallClient github.com/hetznercloud/hcloud-go/v2/hcloud IFirewallClient
//go:generate go run go.uber.org/mock/mockgen -package hcloudmock -destination zz_floating_ip_client_mock.go -mock_names IFloatingIPClient=FloatingIPClient github.com/hetznercloud/hcloud-go/v2/hcloud IFloatingIPClient
//go:generate go run go.uber.org/mock/mockgen -package hcloudmock -destination zz_image_client_mock.go -mock_names IImageClient=ImageClient github.com/hetznercloud/hcloud-go/v2/hcloud IImageClient
//go:generate go run go.uber.org/mock/mockgen -package hcloudmock -destination zz_iso_client_mock.go -mock_names IISOClient=ISOClient github.com/hetznercloud/hcloud-go/v2/hcloud IISOClient
//go:generate go run go.uber.org/mock/mockgen -package hcloudmock -destination zz_loadbalancer_client_mock.go -mock_names ILoadBalancerClient=LoadBalancerClient github.com/hetznercloud/hcloud-go/v2/hcloud ILoadBalancerClient
//go:generate go run go.uber.org/mock/mockgen -package hcloudmock -destination zz_loadbalancer_type_client_mock.go -mock_names ILoadBalancerTypeClient=LoadBalancerTypeClient github.com/hetznercloud/hcloud-go/v2/hcloud ILoadBalancerTypeClient
//go:generate go run go.uber.org/mock/mockgen -package hcloudmock -destination zz_location_client_mock.go -mock_names ILocationClient=LocationClient github.com/hetznercloud/hcloud-go/v2/hcloud ILocationClient
//go:generate go run go.uber.org/mock/mockgen -package hcloudmock -destination zz_network_client_mock.go -mock_names INetworkClient=NetworkClient github.com/hetznercloud/hcloud-go/v2/hcloud INetworkClient
//go:generate go run go.uber.org/mock/mockgen -package hcloudmock -destination zz_placement_group_client_mock.go -mock_names IPlacementGroupClient=PlacementGroupClient github.com/hetznercloud/hcloud-go/v2/hcloud IPlacementGroupClient
//go:generate go run go.uber.org/mock/mockgen -package hcloudmock -destination zz_primary_ip_client_mock.go -mock_names IPrimaryIPClient=PrimaryIPClient github.com/hetznercloud/hcloud-go/v2/hcloud IPrimaryIPClient
//go:generate go run go.uber.org/mock/mockgen -package hcloudmock -destination zz_rdns_client_mock.go -mock_names IRDNSClient=RDNSClient github.com/hetznercloud/hcloud-go/v2/hcloud IRDNSClient
//go:generate go run go.uber.org/mock/mockgen -package hcloudmock -destination zz_server_client_mock.go -mock_names IServerClient=ServerClient github.com/hetznercloud/hcloud-go/v2/hcloud IServerClient
//go:generate go run go.uber.org/mock/mockgen -package hcloudmock -destination zz_server_type_client_mock.go -mock_names IServerTypeClient=ServerTypeClient github.com/hetznercloud/hcloud-go/v2/hcloud IServerTypeClient
//go:generate go run go.uber.org/mock/mockgen -package hcloudmock -destination zz_ssh_key_client_mock.go -mock_names ISSHKeyClient=SSHKeyClient github.com/hetznercloud/hcloud-go/v2/hcloud ISSHKeyClient
//go:generate go run go.uber.org/mock/mockgen -package hcloudmock -destination zz_volume_client_mock.go -mock_names IVolumeClient=VolumeClient github.com/hetznercloud/hcloud-go/v2/hcloud IVolumeClient
