module github.com/hetznercloud/hcloud-go/v2

// Since golang.org/x dependencies always requires version 1.(N-1), this is effectively
// the same version we will be using. (See http://go.dev/issue/69095)
go 1.24.0

toolchain go1.25.4

require (
	github.com/google/go-cmp v0.7.0
	github.com/jmattheis/goverter v1.9.2
	github.com/prometheus/client_golang v1.23.2
	github.com/stretchr/testify v1.11.1
	github.com/vburenin/ifacemaker v1.3.0
	golang.org/x/crypto v0.43.0
	golang.org/x/net v0.46.0
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dave/jennifer v1.6.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/jessevdk/go-flags v1.6.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_model v0.6.2 // indirect
	github.com/prometheus/common v0.66.1 // indirect
	github.com/prometheus/procfs v0.16.1 // indirect
	github.com/rogpeppe/go-internal v1.11.0 // indirect
	go.yaml.in/yaml/v2 v2.4.2 // indirect
	golang.org/x/mod v0.28.0 // indirect
	golang.org/x/sync v0.17.0 // indirect
	golang.org/x/sys v0.37.0 // indirect
	golang.org/x/text v0.30.0 // indirect
	golang.org/x/tools v0.37.0 // indirect
	google.golang.org/protobuf v1.36.8 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
