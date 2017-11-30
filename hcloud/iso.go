package hcloud

import (
	"encoding/json"
)

// ISO represents an ISO image in the Hetzner Cloud.
type ISO struct {
	ID          int
	Name        string
	Description string
	Type        ISOType
}

// ISOType specifies the type of an ISO image.
type ISOType string

const (
	// ISOTypePublic is the type of a public ISO image.
	ISOTypePublic ISOType = "public"

	// ISOTypePrivate is the type of a private ISO image.
	ISOTypePrivate = "private"
)

// UnmarshalJSON implements json.Unmarshaler.
func (i *ISO) UnmarshalJSON(data []byte) error {
	var v struct {
		ID          int     `json:"id"`
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Type        ISOType `json:"type"`
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	i.ID = v.ID
	i.Name = v.Name
	i.Description = v.Description
	i.Type = v.Type

	return nil
}
