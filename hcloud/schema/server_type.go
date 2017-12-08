package schema

// ServerType defines the schema of a server type.
type ServerType struct {
	ID          int                  `json:"id"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Cores       int                  `json:"cores"`
	Memory      float32              `json:"memory"`
	Disk        int                  `json:"disk"`
	StorageType string               `json:"storage_type"`
	Prices      []ServerTypeLocation `json:"prices"`
}

type ServerTypeLocation struct {
	Location     string          `json:"location"`
	PriceHourly  ServerTypePrice `json:"price_hourly"`
	PriceMonthly ServerTypePrice `json:"price_monthly"`
}

type ServerTypePrice struct {
	Net   float32 `json:"net"`
	Gross float32 `json:"gross"`
}

// ServerTypeListResponse defines the schema of the response when
// listing servers.
type ServerTypeListResponse struct {
	ServerTypes []ServerType `json:"server_types"`
}

// ServerTypeGetResponse defines the schema of the response when
// listing servers.
type ServerTypeGetResponse struct {
	ServerType ServerType `json:"server_type"`
}
