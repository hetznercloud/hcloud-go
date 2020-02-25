package schema

import "time"

// Certificate defines the schema of an certificate.
type Certificate struct {
	ID                 int                   `json:"id"`
	Name               string                `json:"name"`
	Type               string                `json:"type"`
	Protection         CertificateProtection `json:"protection"`
	Labels             map[string]string     `json:"labels"`
	Certificate        string                `json:"certificate"`
	Chain              string                `json:"chain"`
	Created            time.Time             `json:"created"`
	NotValidBefore     time.Time             `json:"not_valid_before"`
	NotValidAfter      time.Time             `json:"not_valid_after"`
	DomainNames        []string              `json:"domain_names"`
	Fingerprint        string                `json:"fingerprint"`
	Issuer             string                `json:"issuer"`
	PublicKeyInfo      string                `json:"public_key_info"`
	SignatureAlgorithm string                `json:"signature_algorithm"`
}

// CertificateProtection represents the protection level of a certificate.
type CertificateProtection struct {
	Delete bool `json:"delete"`
}

// CertificateListResponse defines the schema of the response when
// listing Certificates.
type CertificateListResponse struct {
	Certificates []Certificate `json:"certificates"`
}

// CertificateGetResponse defines the schema of the response when
// retrieving a single Certificate.
type CertificateGetResponse struct {
	Certificate Certificate `json:"certificate"`
}

// CertificateCreateRequest defines the schema of the request to create a certificate.
type CertificateCreateRequest struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Certificate string `json:"certificate"`
	Chain       string `json:"chain,omitempty"`
	PrivateKey  string `json:"private_key"`
}

// CertificateCreateResponse defines the schema of the response when creating a certificate.
type CertificateCreateResponse struct {
	Certificate Certificate `json:"certificate"`
}

// CertificateUpdateRequest defines the schema of the request to update a certificate.
type CertificateUpdateRequest struct {
	Name   string             `json:"name,omitempty"`
	Labels *map[string]string `json:"labels,omitempty"`
}

// CertificateUpdateResponse defines the schema of the response when updating a certificate.
type CertificateUpdateResponse struct {
	Certificate Certificate `json:"certificate"`
}

// CertificateActionChangeProtectionRequest defines the schema of the request to
// change the resource protection of a certificate.
type CertificateActionChangeProtectionRequest struct {
	Delete *bool `json:"delete,omitempty"`
}

// CertificateActionChangeProtectionResponse defines the schema of the response when
// changing the resource protection of a certificate.
type CertificateActionChangeProtectionResponse struct {
	Action Action `json:"action"`
}
