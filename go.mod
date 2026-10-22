module github.com/hetznercloud/hcloud-go/v2

// Since golang.org/x dependencies always requires version 1.(N-1), this is effectively
// the same version we will be using. (See http://go.dev/issue/69095)
go 1.26.10

toolchain go1.27.3

require (
	github.com/google/go-cmp v0.7.0
	github.com/prometheus/client_golang v1.25.0
	github.com/stretchr/testify v1.12.1
	golang.org/x/crypto v0.58.0
	golang.org/x/net v0.61.0
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/prometheus/client_model v0.6.3 // indirect
	github.com/prometheus/common v0.72.0 // indirect
	github.com/prometheus/procfs v0.22.0 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/sys v0.49.0 // indirect
	golang.org/x/text v0.43.0 // indirect
	google.golang.org/protobuf v1.36.12 // indirect
)
