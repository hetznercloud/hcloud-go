package schema

import "time"

// Image defines the schema of an image.
type Image struct {
	ID          int       `json:"id"`
	Status      string    `json:"status"`
	Type        string    `json:"type"`
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description"`
	ImageSize   float32   `json:"image_size"`
	DiskSize    float32   `json:"disk_size"`
	Created     time.Time `json:"created"`
	Version     string    `json:"version,omitempty"`
	OSFlavor    string    `json:"os_flavor,omitempty"`
	OSVersion   string    `json:"os_version,omitempty"`
	RapidDeploy bool      `json:"rapid_deploy"`
}

// ImageGetResponse defines the schema of the response when
// retrieving a single image.
type ImageGetResponse struct {
	Image Image `json:"image"`
}

// ImageListResponse defines the schema of the response when
// listing images.
type ImageListResponse struct {
	Images []Image `json:"images"`
}
