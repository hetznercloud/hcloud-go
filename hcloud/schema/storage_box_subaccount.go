package schema

import "time"

type StorageBoxSubaccount struct {
	ID             int64                               `json:"id"`
	Username       string                              `json:"username"`
	HomeDirectory  string                              `json:"home_directory"`
	Server         string                              `json:"server"`
	AccessSettings *StorageBoxSubaccountAccessSettings `json:"access_settings"`
	Description    string                              `json:"description"`
	Labels         map[string]string                   `json:"labels"`
	Created        time.Time                           `json:"created"`
	StorageBox     int64                               `json:"storage_box"`
}

type StorageBoxSubaccountAccessSettings struct {
	ReachableExternally *bool `json:"reachable_externally,omitempty"`
	Readonly            *bool `json:"readonly,omitempty"`
	SambaEnabled        *bool `json:"samba_enabled,omitempty"`
	SSHEnabled          *bool `json:"ssh_enabled,omitempty"`
	WebDAVEnabled       *bool `json:"webdav_enabled,omitempty"`
}

type StorageBoxSubaccountGetResponse struct {
	Subaccount StorageBoxSubaccount `json:"subaccount"`
}

type StorageBoxSubaccountListResponse struct {
	Subaccounts []StorageBoxSubaccount `json:"subaccounts"`
}

type StorageBoxSubaccountCreateRequest struct {
	Password       string                              `json:"password"`
	HomeDirectory  *string                             `json:"home_directory"`
	AccessSettings *StorageBoxSubaccountAccessSettings `json:"access_settings,omitempty"`
	Description    *string                             `json:"description,omitempty"`
	Labels         map[string]string                   `json:"labels,omitempty"`
}

type StorageBoxSubaccountCreateResponse struct {
	Subaccount struct {
		ID         int64 `json:"id"`
		StorageBox int64 `json:"storage_box"`
	} `json:"subaccount"`
	Action Action `json:"action"`
}

type StorageBoxSubaccountUpdateRequest struct {
	Labels      map[string]string `json:"labels"`
	Description string            `json:"description"`
}

type StorageBoxSubaccountUpdateResponse struct {
	Subaccount StorageBoxSubaccount `json:"subaccount"`
}
