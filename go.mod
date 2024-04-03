module github.com/hetznercloud/hcloud-go/v2

// This is being kept at a lower value on purpose as raising this would require
// all dependends to update to the new version.
// As long as we do not depend on any newer language feature this can be kept at the current value.
// It should never be higher than the lowest currently supported version of Go.
go 1.20

require (
	github.com/google/go-cmp v0.6.0
	github.com/jmattheis/goverter v1.4.0
	github.com/prometheus/client_golang v1.19.0
	github.com/stretchr/testify v1.9.0
	github.com/vburenin/ifacemaker v1.2.1
	golang.org/x/net v0.23.0
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dave/jennifer v1.6.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/jessevdk/go-flags v1.4.1-0.20181029123624-5de817a9aa20 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_model v0.5.0 // indirect
	github.com/prometheus/common v0.48.0 // indirect
	github.com/prometheus/procfs v0.12.0 // indirect
	github.com/rogpeppe/go-internal v1.11.0 // indirect
	golang.org/x/mod v0.14.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/tools v0.17.0 // indirect
	google.golang.org/protobuf v1.32.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
