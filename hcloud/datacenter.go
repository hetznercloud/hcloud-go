package hcloud

// Datacenter represents a datacenter in the Hetzner Cloud.
type Datacenter struct {
	ID          int
	Name        string
	Description string
	Location    *Location
}
