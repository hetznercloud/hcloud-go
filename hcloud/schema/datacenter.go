package schema

// Datacenter defines the schema of a datacenter.
type Datacenter struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Location    Location `json:"location"`
}
