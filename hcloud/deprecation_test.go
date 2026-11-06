package hcloud

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/hcloud-go/v2/hcloud/schema"
)

type TestDeprecatableResource struct {
	DeprecatableResource
}

// Interface is implemented.
var _ Deprecatable = TestDeprecatableResource{}

func TestDeprecatableResource_IsDeprecated(t *testing.T) {
	tests := []struct {
		name     string
		resource TestDeprecatableResource
		want     bool
	}{
		{name: "nil returns false", resource: TestDeprecatableResource{DeprecatableResource: DeprecatableResource{Deprecation: nil}}, want: false},
		{name: "struct returns true", resource: TestDeprecatableResource{DeprecatableResource: DeprecatableResource{Deprecation: &DeprecationInfo{}}}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.resource.IsDeprecated(), "IsDeprecated()")
		})
	}
}

func TestDeprecatableResource_DeprecationAnnounced(t *testing.T) {
	tests := []struct {
		name     string
		resource TestDeprecatableResource
		want     time.Time
	}{
		{
			name:     "nil returns default time",
			resource: TestDeprecatableResource{DeprecatableResource: DeprecatableResource{Deprecation: nil}},
			want:     time.Unix(0, 0)},
		{
			name:     "actual value is returned",
			resource: TestDeprecatableResource{DeprecatableResource: DeprecatableResource{Deprecation: &DeprecationInfo{Announced: mustParseTime(t, "2023-06-01T00:00:00+00:00")}}},
			want:     mustParseTime(t, "2023-06-01T00:00:00+00:00")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.resource.DeprecationAnnounced(), "DeprecationAnnounced()")
		})
	}
}

func TestDeprecatableResource_UnavailableAfter(t *testing.T) {
	tests := []struct {
		name     string
		resource TestDeprecatableResource
		want     time.Time
	}{
		{
			name:     "nil returns default time",
			resource: TestDeprecatableResource{DeprecatableResource: DeprecatableResource{Deprecation: nil}},
			want:     time.Unix(0, 0)},
		{
			name:     "actual value is returned",
			resource: TestDeprecatableResource{DeprecatableResource: DeprecatableResource{Deprecation: &DeprecationInfo{UnavailableAfter: mustParseTime(t, "2023-06-01T00:00:00+00:00")}}},
			want:     mustParseTime(t, "2023-06-01T00:00:00+00:00")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.resource.UnavailableAfter(), "UnavailableAfter()")
		})
	}
}

func TestDeprecationSchema(t *testing.T) {
	testCases := map[string]string{
		"deprecated_resource": `{
			"deprecation": {
				"announced": "2023-06-01T00:00:00+00:00",
				"unavailable_after": "2023-09-01T00:00:00+00:00"
			}
		}`,
		"non_deprecated_resource": `{
			"deprecation": null
		}`,
	}

	for name, data := range testCases {
		t.Run(name, func(t *testing.T) {
			var s *schema.DeprecationInfo
			assert.NoError(t, json.Unmarshal([]byte(data), &s))

			assert.Equal(t, s, SchemaFromDeprecation(DeprecationFromSchema(s)))

			d := DeprecationFromSchema(s)
			assert.Equal(t, d, DeprecationFromSchema(SchemaFromDeprecation(d)))
		})
	}
}
