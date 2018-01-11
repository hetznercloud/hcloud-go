package hcloud

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
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
	CreatedFrom *Server
	BoundTo     *Server
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

// GetByID retrieves an image by its ID.
func (c *ImageClient) GetByID(ctx context.Context, id int) (*Image, *Response, error) {
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

// GetByName retrieves an image by its name.
func (c *ImageClient) GetByName(ctx context.Context, name string) (*Image, *Response, error) {
	path := "/images?name=" + url.QueryEscape(name)
	req, err := c.client.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	var body schema.ImageListResponse
	resp, err := c.client.Do(req, &body)
	if err != nil {
		return nil, nil, err
	}

	if len(body.Images) == 0 {
		return nil, resp, nil
	}
	return ImageFromSchema(body.Images[0]), resp, nil
}

// Get retrieves an image by its ID if the input can be parsed as an integer, otherwise it retrieves an image by its name.
func (c *ImageClient) Get(ctx context.Context, idOrName string) (*Image, *Response, error) {
	if id, err := strconv.Atoi(idOrName); err == nil {
		return c.GetByID(ctx, int(id))
	}
	return c.GetByName(ctx, idOrName)
}

// ImagePage serves as accessor of the images API pagination.
type ImagePage struct {
	Page
	content []*Image
}

// Content contains the content of the current page.
func (p *ImagePage) Content() []*Image {
	return p.content
}

// ImageListOpts specifies options for listing images.
type ImageListOpts struct {
	ListOpts
	Types   []ImageType
	BoundTo []*Server
}

// URLValues returns the list opts as url.Values.
func (o ImageListOpts) URLValues() url.Values {
	vals := o.ListOpts.URLValues()
	for _, t := range o.Types {
		vals.Add("type", string(t))
	}
	for _, bt := range o.BoundTo {
		vals.Add("bound_to", strconv.Itoa(bt.ID))
	}
	return vals
}

// List returns an accessor to control the images API pagination.
func (c *ImageClient) List(ctx context.Context, opts ImageListOpts) *ImagePage {
	page := &ImagePage{}
	page.pageGetter = pageGetter(func(start, end int) (resp *Response, exhausted bool, err error) {
		allImages := []*Image{}
		if opts.PerPage == 0 {
			opts.PerPage = 50
		}

		resp, exhausted, err = c.client.all(func(page int) (*Response, error) {
			opts.Page = page
			images, resp, err := c.list(ctx, opts)
			if err != nil {
				return resp, err
			}
			allImages = append(allImages, images...)
			return resp, nil
		}, start, end)
		page.content = allImages
		return
	})
	return page
}

// list returns a list of images for a specific page.
func (c *ImageClient) list(ctx context.Context, opts ImageListOpts) ([]*Image, *Response, error) {
	path := "/images?" + opts.URLValues().Encode()
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
	opts := ImageListOpts{}
	opts.PerPage = 50
	page := c.List(ctx, opts)
	if page.All(); page.Err() != nil {
		return nil, page.Err()
	}
	return page.Content(), nil
}

// Delete deletes an image.
func (c *ImageClient) Delete(ctx context.Context, image *Image) (*Response, error) {
	req, err := c.client.NewRequest(ctx, "DELETE", fmt.Sprintf("/images/%d", image.ID), nil)
	if err != nil {
		return nil, err
	}
	return c.client.Do(req, nil)
}

// ImageUpdateOpts specifies options for updating an image.
type ImageUpdateOpts struct {
	Description *string
	Type        ImageType
}

// Update updates an image.
func (c *ImageClient) Update(ctx context.Context, image *Image, opts ImageUpdateOpts) (*Image, *Response, error) {
	reqBody := schema.ImageUpdateRequest{
		Description: opts.Description,
	}
	if opts.Type != "" {
		reqBody.Type = String(string(opts.Type))
	}
	reqBodyData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, nil, err
	}

	path := fmt.Sprintf("/images/%d", image.ID)
	req, err := c.client.NewRequest(ctx, "PUT", path, bytes.NewReader(reqBodyData))
	if err != nil {
		return nil, nil, err
	}

	respBody := schema.ImageUpdateResponse{}
	resp, err := c.client.Do(req, &respBody)
	if err != nil {
		return nil, resp, err
	}
	return ImageFromSchema(respBody.Image), resp, nil
}
