package hcloud

import (
	"context"
	"fmt"
	"time"

	"github.com/hetznercloud/hcloud-go/hcloud/schema"
)

// Image represents an Image in the Hetzner Cloud.
type Image struct {
	ID          int
	Name        string
	Type        ImageType
	Status      ImageStatus
	Description string
	ImageSize   float32
	DiskSize    float32
	Created     time.Time
	RapidDeploy bool

	OSFlavor  string
	OSVersion string
}

// ImageType specifies the type of an image.
type ImageType string

const (
	// ImageTypeSnapshot represents a snapshot image.
	ImageTypeSnapshot = "snapshot"
	// ImageTypeBackup represents a backup image.
	ImageTypeBackup = "backup"
	// ImageTypeSystem represents a system image.
	ImageTypeSystem = "system"
)

// ImageStatus specifies the status of an image.
type ImageStatus string

const (
	// ImageStatusCreating is the status when an image is being created.
	ImageStatusCreating = "creating"
	// ImageStatusAvailable is the stats when an image is available.
	ImageStatusAvailable = "available"
)

// ImageClient is a client for the image API.
type ImageClient struct {
	client *Client
}

// Get retrieves an image.
func (c *ImageClient) Get(ctx context.Context, id int) (*Image, *Response, error) {
	req, err := c.client.NewRequest(ctx, "GET", fmt.Sprintf("/images/%d", id), nil)
	if err != nil {
		return nil, nil, err
	}

	var body schema.ImageGetResponse
	resp, err := c.client.Do(req, &body)
	if err != nil {
		if IsError(err, ErrorCodeNotFound) {
			return nil, resp, nil
		}
		return nil, nil, err
	}
	return ImageFromSchema(body.Image), resp, nil
}

// ImageListOpts specifies options for listing images.
type ImageListOpts struct {
	ListOpts
}

// List returns a list of images for a specific page.
func (c *ImageClient) List(ctx context.Context, opts ImageListOpts) ([]*Image, *Response, error) {
	path := "/images?" + valuesForListOpts(opts.ListOpts).Encode()
	req, err := c.client.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	var body schema.ImageListResponse
	resp, err := c.client.Do(req, &body)
	if err != nil {
		return nil, nil, err
	}
	images := make([]*Image, 0, len(body.Images))
	for _, i := range body.Images {
		images = append(images, ImageFromSchema(i))
	}
	return images, resp, nil
}

// All returns all images.
func (c *ImageClient) All(ctx context.Context) ([]*Image, error) {
	allImages := []*Image{}

	opts := ImageListOpts{}
	opts.PerPage = 50

	_, err := c.client.all(func(page int) (*Response, error) {
		opts.Page = page
		images, resp, err := c.List(ctx, opts)
		if err != nil {
			return resp, err
		}
		allImages = append(allImages, images...)
		return resp, nil
	})
	if err != nil {
		return nil, err
	}

	return allImages, nil
}

// Delete deletes an image.
func (c *ImageClient) Delete(ctx context.Context, id int) (*Response, error) {
	req, err := c.client.NewRequest(ctx, "DELETE", fmt.Sprintf("/images/%d", id), nil)
	if err != nil {
		return nil, err
	}
	return c.client.Do(req, nil)
}
