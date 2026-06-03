package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuildUserAgent(t *testing.T) {
	testCases := []struct {
		desc    string
		name    string
		version string
		want    string
	}{
		{"with application name and version", "test", "1.0", "test/1.0 hcloud-go/1.42.0"},
		{"with application name but no version", "test", "", "test hcloud-go/1.42.0"},
		{"without application name and version", "", "", "hcloud-go/1.42.0"},
	}

	for _, tt := range testCases {
		t.Run(tt.desc, func(t *testing.T) {
			got := BuildUserAgent(tt.name, tt.version, "hcloud-go/1.42.0")
			require.Equal(t, tt.want, got)
		})
	}
}
