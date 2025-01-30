# hcloud: A Go library for the Hetzner Cloud API

[![GitHub Actions status](https://github.com/hetznercloud/hcloud-go/workflows/Continuous%20Integration/badge.svg)](https://github.com/hetznercloud/hcloud-go/actions)
[![Codecov](https://codecov.io/github/hetznercloud/hcloud-go/graph/badge.svg?token=4IAbGIwNYp)](https://codecov.io/github/hetznercloud/hcloud-go/tree/main)
[![Go Reference](https://pkg.go.dev/badge/github.com/hetznercloud/hcloud-go/v2/hcloud.svg)](https://pkg.go.dev/github.com/hetznercloud/hcloud-go/v2/hcloud)

Package hcloud is a library for the Hetzner Cloud API.

The library’s documentation is available at [pkg.go.dev](https://godoc.org/github.com/hetznercloud/hcloud-go/v2/hcloud),
the public API documentation is available at [docs.hetzner.cloud](https://docs.hetzner.cloud/).

> [!IMPORTANT]
> Make sure to follow our API changelog available at
> [docs.hetzner.cloud/changelog](https://docs.hetzner.cloud/changelog) (or the RRS feed
> available at
> [docs.hetzner.cloud/changelog/feed.rss](https://docs.hetzner.cloud/changelog/feed.rss))
> to be notified about additions, deprecations and removals.

## Installation

```sh
go get github.com/hetznercloud/hcloud-go/v2/hcloud
```

## Example

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/hetznercloud/hcloud-go/v2/hcloud"
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

## Upgrading

### Support

- `v2` is actively maintained by Hetzner Cloud
- `v1` is supported until September 1st 2023 and will continue to receive new features until then. See [#263](https://github.com/hetznercloud/hcloud-go/issues/263).

### From v1 to v2

Version 2.0.0 was published because we changed the datatype of all `ID` fields from `int` to `int64`.

To migrate to the new version, replace all your imports to reference the new module path:

```diff
 import (
-  "github.com/hetznercloud/hcloud-go/hcloud"
+  "github.com/hetznercloud/hcloud-go/v2/hcloud"
 )
```

When you compile your code, it will show any invalid usages of `int` in your code that you need to fix. We commonly found these changes while updating our integrations:

- `strconv.Atoi(idString)` (parsing integers) needs to be replaced by `strconv.ParseInt(idString, 10, 64)`
- `strconv.Itoa(id)` (formatting integers) needs to be replaced by `strconv.FormatInt(id, 10)`

## Go Version Support

The library supports the latest two Go minor versions, e.g. at the time Go 1.19 is released, it supports Go 1.18 and 1.19.

This matches the official [Go Release Policy](https://go.dev/doc/devel/release#policy).

When the minimum required Go version is changed, it is announced in the release notes for that version.

## License

MIT license
