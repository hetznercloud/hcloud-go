# hcloud: A Go library for the Hetzner Cloud API

[![GitHub Actions status](https://github.com/hetznercloud/hcloud-go/workflows/Continuous%20Integration/badge.svg)](https://github.com/hetznercloud/hcloud-go/actions)
[![GoDoc](https://godoc.org/github.com/hetznercloud/hcloud-go/hcloud?status.svg)](https://godoc.org/github.com/hetznercloud/hcloud-go/hcloud)

Package hcloud is a library for the Hetzner Cloud API.

The library’s documentation is available at [GoDoc](https://godoc.org/github.com/hetznercloud/hcloud-go/hcloud),
the public API documentation is available at [docs.hetzner.cloud](https://docs.hetzner.cloud/).

> [!IMPORTANT]
> v1 is unsupported since February 2025.

## Example

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/hetznercloud/hcloud-go/hcloud"
)

func main() {
    client := hcloud.NewClient(hcloud.WithToken("token"))

    server, _, err := client.Server.GetByID(context.Background(), 1)
    if err != nil {
        log.Fatalf("error retrieving server: %s\n", err)
    }
    if server != nil {
        fmt.Printf("server 1 is called %q\n", server.Name)
    } else {
        fmt.Println("server 1 not found")
    }
}
```

## Go Version Support

The library supports the latest two Go minor versions, e.g. at the time Go 1.19 is released, it supports Go 1.18 and 1.19.

This matches the official [Go Release Policy](https://go.dev/doc/devel/release#policy).

When the minimum required Go version is changed, it is announced in the release notes for that version.

## License

MIT license
