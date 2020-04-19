package hcloud

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/hetznercloud/hcloud-go/hcloud/schema"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// DNSEndpoint is the base URL of the DNS server API.
const DNSEndpoint = "https://dns.hetzner.com/api/v1"

// RecordType represents an record's type.
type RecordType string

const (
	// A represents a 'Host address
	A RecordType = "A"
	// AAAA represents a 'IPv6 host address'
	AAAA = "AAAA"
	// NS represents a 'Name Server'
	NS = "NS"
	// MX represents a 'Mail eXchange'
	MX = "MX"
	// CNAME represents a 'Canonical name for an alias'
	CNAME = "CNAME"
	// RP represents a 'Responsible person'
	RP = "RP"
	// TXT represents a 'Descriptive text'
	TXT = "TXT"
	// SOA represents a 'Start Of Authority'
	SOA = "SOA"
	// HINFO represents a 'Host information'
	HINFO = "HINFO"
	// SRV represents a 'Location of service'
	SRV = "SRV"
	// DANE represents a 'DNS-based Authentication of Named Entities'
	DANE = "DANE"
	// TLSA represents a 'Transport Layer Security Authentication'
	TLSA = "TLSA"
	// DS represents a 'Delegation Signer'
	DS = "DS"
	// CAA represents a 'Certification Authority Authorization'
	CAA = "CAA"
)

// Record defines the schema of a Dns server record.
type Record struct {
	ID       string
	Name     string
	Created  time.Time
	Modified time.Time
	TTL      uint64
	Type     RecordType
	Value    string
	ZoneID   string
}

// CreateOrUpdateRecord  defines the schema for creating or updating a Dns server record.
type CreateOrUpdateRecord struct {
	Name   string     // required
	Type   RecordType // required
	Value  string     // required
	ZoneID string     // required
	TTL    *uint64    // optional
}

// Zone defines the schema of a Dns server zone.
type Zone struct {
	ID              string
	Name            string
	Ns              []string
	LegacyDNSHost   string
	LegacyNs        []string
	Created         time.Time
	Modified        time.Time
	Verified        time.Time
	Owner           string
	Paused          bool
	Permission      string
	Project         string
	Registrar       string
	Status          string
	TTL             uint64
	IsSecondaryDNS  bool
	TxtVerification TxtVerification
	ZoneType        ZoneType
	RecordsCount    int
}

// TxtVerification defines the schema of a update Dns server TXT verification.
type TxtVerification struct {
	Name  string
	Token string
}

// ZoneType defines the schema of a Dns server zone type.
type ZoneType struct {
	ID          string
	Name        string
	Description string
	Prices      []Price
}

// CreateOrUpdateZone defines the schema for creating a Dns server zone.
type CreateOrUpdateZone struct {
	Name string  // required
	TTL  *uint64 // optional
}

// DNSServerClient is the Api client for Dns server management.
type DNSServerClient struct {
	client *Client
}

// NewRequest creates a new request for the DNS server client.
func (c *DNSServerClient) NewRequest(ctx context.Context, method, path string, body io.Reader) (*http.Request, error) {
	// todo: Currently it is not supported to override the endpoint. Maybe endpoint should not be http-client dependent, every api-client should have a default endpoint.
	url := DNSEndpoint + path
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", c.client.userAgent)
	req.Header.Set("Auth-API-Token", c.client.token)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req = req.WithContext(ctx)
	return req, nil
}

// RecordListOpts specifies options for listing record resources.
type RecordListOpts struct {
	Page    int // Page (starting at 1)
	PerPage int // Items per page (0 means default)
}

// ZoneListOpts specifies options for listing record resources.
type ZoneListOpts struct {
	Page       int     // Page (starting at 1)
	PerPage    int     // Items per page (0 means default)
	SearchName *string // Partial name of a zone. Will return all zones that contain the searched string (nul means default)
}

func (l RecordListOpts) values() url.Values {
	vals := url.Values{}
	if l.Page > 0 {
		vals.Add("page", strconv.Itoa(l.Page))
	}
	if l.PerPage > 0 {
		vals.Add("per_page", strconv.Itoa(l.PerPage))
	}
	return vals
}

func (l ZoneListOpts) values() url.Values {
	vals := url.Values{}
	if l.Page > 0 {
		vals.Add("page", strconv.Itoa(l.Page))
	}
	if l.PerPage > 0 {
		vals.Add("per_page", strconv.Itoa(l.PerPage))
	}
	if l.SearchName != nil {
		vals.Add("search_name", *l.SearchName)
	}
	return vals
}

// GetAllRecords returns all records associated with user.
func (c *DNSServerClient) GetAllRecords(ctx context.Context, zoneID string, opts RecordListOpts) ([]*Record, *Response, error) {
	path := "/records?zone_id=" + zoneID
	params := opts.values().Encode()
	if len(params) > 0 {
		path += "&" + params
	}

	req, err := c.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	var body schema.RecordsResponse
	resp, err := c.client.Do(req, &body)
	if err != nil {
		return nil, nil, err
	}
	records := make([]*Record, 0, len(body.Records))
	for _, s := range body.Records {
		records = append(records, RecordFromSchema(s))
	}
	return records, resp, nil
}

// GetRecord returns information about a single record.
func (c *DNSServerClient) GetRecord(ctx context.Context, recordID string) (*Record, *Response, error) {
	path := "/records/" + recordID
	req, err := c.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	var body schema.RecordResponse
	resp, err := c.client.Do(req, &body)
	if err != nil {
		return nil, nil, err
	}
	return RecordFromSchema(body.Record), resp, nil
}

// CreateRecord creates a new record.
func (c *DNSServerClient) CreateRecord(ctx context.Context, record CreateOrUpdateRecord) (*Record, *Response, error) {
	reqBody := schema.CreateRecordRequest{
		Name:   record.Name,
		Type:   string(record.Type),
		Value:  record.Value,
		TTL:    record.TTL,
		ZoneID: record.ZoneID,
	}
	reqBodyData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, nil, err
	}

	path := "/records"
	req, err := c.NewRequest(ctx, "POST", path, bytes.NewReader(reqBodyData))
	if err != nil {
		return nil, nil, err
	}

	respBody := schema.RecordResponse{}
	resp, err := c.client.Do(req, &respBody)
	if err != nil {
		return nil, resp, err
	}
	return RecordFromSchema(respBody.Record), resp, nil
}

// UpdateRecord updates a record.
func (c *DNSServerClient) UpdateRecord(ctx context.Context, record CreateOrUpdateRecord, recordID string) (*Record, *Response, error) {
	reqBody := schema.UpdateRecordRequest{
		Name:   record.Name,
		Type:   string(record.Type),
		Value:  record.Value,
		TTL:    record.TTL,
		ZoneID: record.ZoneID,
	}
	reqBodyData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, nil, err
	}

	path := "/records/" + recordID
	req, err := c.NewRequest(ctx, "PUT", path, bytes.NewReader(reqBodyData))
	if err != nil {
		return nil, nil, err
	}

	respBody := schema.RecordResponse{}
	resp, err := c.client.Do(req, &respBody)
	if err != nil {
		return nil, resp, err
	}
	return RecordFromSchema(respBody.Record), resp, nil
}

// DeleteRecord deletes a record.
func (c *DNSServerClient) DeleteRecord(ctx context.Context, recordID string) (*Response, error) {
	path := "/records/" + recordID
	req, err := c.NewRequest(ctx, "DELETE", path, nil)
	if err != nil {
		return nil, err
	}

	return c.client.Do(req, nil)
}

// GetAllZones returns all zones associated with the user.
func (c *DNSServerClient) GetAllZones(ctx context.Context, opts ZoneListOpts) ([]*Zone, *Response, error) {
	path := "/zones"
	params := opts.values().Encode()
	if len(params) > 0 {
		path += "?" + params
	}
	req, err := c.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	var body schema.ZonesResponse
	resp, err := c.client.Do(req, &body)
	if err != nil {
		return nil, nil, err
	}
	zones := make([]*Zone, 0, len(body.Zones))
	for _, s := range body.Zones {
		zones = append(zones, ZoneFromSchema(s))
	}
	return zones, resp, nil
}

// GetZone returns an object containing all information about a zone. Zone to get is identified by 'ZoneID'.
func (c *DNSServerClient) GetZone(ctx context.Context, zoneID string) (*Zone, *Response, error) {
	path := "/zones/" + zoneID
	req, err := c.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	respBody := schema.ZoneResponse{}
	resp, err := c.client.Do(req, &respBody)
	if err != nil {
		return nil, resp, err
	}
	return ZoneFromSchema(respBody.Zone), resp, nil
}

// UpdateZone updates a zone.
func (c *DNSServerClient) UpdateZone(ctx context.Context, zoneID string, updateZoneData CreateOrUpdateZone) (*Zone, *Response, error) {
	reqBody := schema.UpdateZoneRequest{
		Name: updateZoneData.Name,
		TTL:  updateZoneData.TTL,
	}
	reqBodyData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, nil, err
	}

	path := "/zones/" + zoneID
	req, err := c.NewRequest(ctx, "PUT", path, bytes.NewReader(reqBodyData))
	if err != nil {
		return nil, nil, err
	}

	respBody := schema.ZoneResponse{}
	resp, err := c.client.Do(req, &respBody)
	if err != nil {
		return nil, resp, err
	}
	return ZoneFromSchema(respBody.Zone), resp, nil
}

// DeleteZone deletes a zone.
func (c *DNSServerClient) DeleteZone(ctx context.Context, zoneID string) (*Response, error) {
	path := "/zones/" + zoneID
	req, err := c.NewRequest(ctx, "DELETE", path, nil)
	if err != nil {
		return nil, err
	}

	return c.client.Do(req, nil)
}

// CreateZone creates a new zone.
func (c *DNSServerClient) CreateZone(ctx context.Context, zone CreateOrUpdateZone) (*Zone, *Response, error) {
	reqBody := schema.CreateZoneRequest{
		Name: zone.Name,
		TTL:  zone.TTL,
	}
	reqBodyData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, nil, err
	}

	path := "/zones"
	req, err := c.NewRequest(ctx, "POST", path, bytes.NewReader(reqBodyData))
	if err != nil {
		return nil, nil, err
	}

	respBody := schema.ZoneResponse{}
	resp, err := c.client.Do(req, &respBody)
	if err != nil {
		return nil, resp, err
	}
	return ZoneFromSchema(respBody.Zone), resp, nil
}
