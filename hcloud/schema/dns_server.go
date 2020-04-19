package schema

import (
	"strings"
	"time"
)

type specialTime struct {
	time.Time
}

// UnmarshalJSON converts a specialTime.
func (st *specialTime) UnmarshalJSON(b []byte) error {
	// todo: Currently dns-api doest not support Iso date-time formats. Additionally: different time formats are delivered. And time.Parse in golang is not very smart.
	// Possible values: "2020-04-19 13:03:28.975 +0000 UTC" , "2020-04-19T13:03:30Z" => time.RFC3339Nano
	input := string(b)
	input = strings.Trim(input, `"`)
	input = strings.Replace(input, " +0000 UTC", "Z", -1)
	input = strings.Replace(input, " ", "T", -1)
	t, err := time.Parse(time.RFC3339Nano, input)
	if err != nil {
		return err
	}
	st.Time = t
	return nil
}

// Record defines the schema of a Dns server record.
type Record struct {
	ID       string      `json:"id"`
	Name     string      `json:"name"`
	Created  specialTime `json:"created"`
	Modified specialTime `json:"modified"`
	TTL      uint64      `json:"ttl"`
	Type     string      `json:"type"`
	Value    string      `json:"value"`
	ZoneID   string      `json:"zone_id"`
}

// RecordResponse defines the schema of a Dns server record response.
type RecordResponse struct {
	Record Record `json:"record"`
}

// RecordsResponse defines the schema of a Dns server records response.
type RecordsResponse struct {
	Records []Record `json:"records"`
}

// CreateRecordRequest defines the schema of a create Dns server record request.
type CreateRecordRequest struct {
	Name   string  `json:"name"`
	TTL    *uint64 `json:"ttl,omitempty"`
	Type   string  `json:"type"`
	Value  string  `json:"value"`
	ZoneID string  `json:"zone_id"`
}

// BulkCreateRecordRequest defines the schema of a create Dns server records bulk request.
type BulkCreateRecordRequest struct {
	Records []CreateRecordRequest `json:"records"`
}

// BulkCreateRecordResponse defines the schema of a create Dns server records bulk response.
type BulkCreateRecordResponse struct {
	Records        []Record `json:"records"`
	ValidRecords   []Record `json:"valid_records"`
	InvalidRecords []Record `json:"invalid_records"`
}

// UpdateRecordRequest defines the schema of a update Dns server record request.
type UpdateRecordRequest struct {
	Name   string  `json:"name"`
	TTL    *uint64 `json:"ttl,omitempty"`
	Type   string  `json:"type"`
	Value  string  `json:"value"`
	ZoneID string  `json:"zone_id"`
}

// BulkUpdateRecordRequest defines the schema of a update Dns server record bulk request.
type BulkUpdateRecordRequest struct {
	Records []UpdateRecordRequest `json:"records"`
}

// BulkUpdateRecordResponse defines the schema of a update Dns server record bulk response.
type BulkUpdateRecordResponse struct {
	Records       []Record `json:"records"`
	FailedRecords []Record `json:"failed_records"`
}

// ZoneResponse defines the schema of a update Dns server zone response.
type ZoneResponse struct {
	Zone Zone `json:"zone"`
}

// ZonesResponse defines the schema of a update Dns server zones response.
type ZonesResponse struct {
	Zones []Zone `json:"zones"`
}

// Zone defines the schema of a update Dns server zone.
type Zone struct {
	ID              string          `json:"id"`
	Name            string          `json:"name"`
	Ns              []string        `json:"ns"`
	LegacyDNSHost   string          `json:"legacy_dns_host"`
	LegacyNs        []string        `json:"legacy_ns"`
	Created         specialTime     `json:"created"`
	Modified        specialTime     `json:"modified"`
	Verified        specialTime     `json:"verified"`
	Owner           string          `json:"owner"`
	Paused          bool            `json:"paused"`
	Permission      string          `json:"permission"`
	Project         string          `json:"project"`
	Registrar       string          `json:"registrar"`
	Status          string          `json:"status"` // todo typed (which status types are available?)
	TTL             uint64          `json:"ttl"`
	IsSecondaryDNS  bool            `json:"is_secondary_dns"`
	TxtVerification TxtVerification `json:"txt_verification"`
	ZoneType        ZoneType        `json:"zone_type"`
	RecordsCount    int             `json:"records_count"`
}

// TxtVerification defines the schema of a update Dns server TXT verification.
type TxtVerification struct {
	Name  string `json:"name"`
	Token string `json:"token"`
}

// ZoneType defines the schema of a Dns server zone type.
type ZoneType struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Prices      []Price `json:"prices"`
}

// CreateZoneRequest defines the schema of a create Dns server zone request.
type CreateZoneRequest struct {
	Name string  `json:"name"`
	TTL  *uint64 `json:"ttl,omitempty"`
}

// UpdateZoneRequest defines the schema of a update Dns server zone request.
type UpdateZoneRequest struct {
	ID   string  `json:"id"`
	Name string  `json:"name"`
	TTL  *uint64 `json:"ttl,omitempty"`
}
