package hcloud

import "github.com/hetznercloud/hcloud-go/hcloud/schema"

// Location represents a location in the Hetzner Cloud.
type Location struct {
	ID          int
	Name        string
	Description string
	Country     string
	City        string
	Latitude    float64
	Longitude   float64
}

// LocationFromSchema converts a schema.Location to a Location.
func LocationFromSchema(s schema.Location) *Location {
	return &Location{
		ID:          s.ID,
		Name:        s.Name,
		Description: s.Description,
		Country:     s.Country,
		City:        s.City,
		Latitude:    s.Latitude,
		Longitude:   s.Longitude,
	}
}
