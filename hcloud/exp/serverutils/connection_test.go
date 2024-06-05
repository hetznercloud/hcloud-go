package serverutils

import (
	"testing"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/schema"
	"github.com/stretchr/testify/assert"
)

func TestFirstAvailableIP(t *testing.T) {
	testCases := []struct {
		name   string
		server *hcloud.Server
		want   string
	}{
		{
			name:   "empty",
			server: &hcloud.Server{},
			want:   "",
		},
		{
			name: "public_ipv4",
			server: &hcloud.Server{
				PublicNet: hcloud.ServerPublicNetFromSchema(schema.ServerPublicNet{
					IPv4: schema.ServerPublicNetIPv4{ID: 1, IP: "1.2.3.4"},
					IPv6: schema.ServerPublicNetIPv6{ID: 2, IP: "2a01:4f8:1c19:1403::/64"},
				}),
				PrivateNet: []hcloud.ServerPrivateNet{
					hcloud.ServerPrivateNetFromSchema(schema.ServerPrivateNet{Network: 3, IP: "10.0.0.1"}),
				},
			},
			want: "1.2.3.4",
		},
		{
			name: "public_ipv6",
			server: &hcloud.Server{
				PublicNet: hcloud.ServerPublicNetFromSchema(schema.ServerPublicNet{
					IPv6: schema.ServerPublicNetIPv6{ID: 2, IP: "2a01:4f8:1c19:1403::/64"},
				}),
				PrivateNet: []hcloud.ServerPrivateNet{
					hcloud.ServerPrivateNetFromSchema(schema.ServerPrivateNet{Network: 3, IP: "10.0.0.1"}),
				},
			},
			want: "2a01:4f8:1c19:1403::1",
		},
		{
			name: "private_ipv4",
			server: &hcloud.Server{
				PrivateNet: []hcloud.ServerPrivateNet{
					hcloud.ServerPrivateNetFromSchema(schema.ServerPrivateNet{Network: 3, IP: "10.0.0.1"}),
				},
			},
			want: "10.0.0.1",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result := FirstAvailableIP(testCase.server)
			assert.Equal(t, testCase.want, result)
		})
	}
}
