module github.com/hetznercloud/hcloud-go

// This is being kept at a lower value on purpose as raising this would require
// all dependends to update to the new version.
// As long as we do not depend on any newer language feature this can be kept at the current value.
// It should never be higher than the lowest currently supported version of Go.
go 1.21

require (
	github.com/google/go-cmp v0.5.9
	github.com/prometheus/client_golang v1.16.0
	github.com/stretchr/testify v1.8.4
	golang.org/x/net v0.12.0
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_model v0.3.0 // indirect
	github.com/prometheus/common v0.42.0 // indirect
	github.com/prometheus/procfs v0.10.1 // indirect
	github.com/rogpeppe/go-internal v1.12.0 // indirect
	golang.org/x/sys v0.10.0 // indirect
	golang.org/x/text v0.11.0 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
