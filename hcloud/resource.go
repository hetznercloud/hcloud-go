package hcloud

// Resource defines the schema of a resource.
type Resource struct {
	ID   int64
	Type string
}

type SupportsActions interface {
	pathID() string
}

var _ SupportsActions = &Certificate{}
var _ SupportsActions = &Firewall{}
var _ SupportsActions = &FloatingIP{}
var _ SupportsActions = &Image{}
var _ SupportsActions = &LoadBalancer{}
var _ SupportsActions = &Network{}
var _ SupportsActions = &PrimaryIP{}
var _ SupportsActions = &Server{}
var _ SupportsActions = &Volume{}
var _ SupportsActions = &Zone{}
var _ SupportsActions = &StorageBox{}
