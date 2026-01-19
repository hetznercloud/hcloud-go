# hcloud: A Go library for the Hetzner Cloud API

[![GitHub Actions status](https://github.com/hetznercloud/hcloud-go/workflows/Continuous%20Integration/badge.svg)](https://github.com/hetznercloud/hcloud-go/actions)
[![Codecov](https://codecov.io/github/hetznercloud/hcloud-go/graph/badge.svg?token=4IAbGIwNYp)](https://codecov.io/github/hetznercloud/hcloud-go/tree/main)
[![Go Reference](https://pkg.go.dev/badge/github.com/hetznercloud/hcloud-go/v2/hcloud.svg)](https://pkg.go.dev/github.com/hetznercloud/hcloud-go/v2/hcloud)

Package hcloud is a library for the Hetzner Cloud API.

The library’s documentation is available at [pkg.go.dev](https://pkg.go.dev/github.com/hetznercloud/hcloud-go/v2/hcloud),
the public API documentation is available at [docs.hetzner.cloud](https://docs.hetzner.cloud/).

> [!IMPORTANT]
> Make sure to follow our API changelog available at
> [docs.hetzner.cloud/changelog](https://docs.hetzner.cloud/changelog) (or the RSS feed
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
	"github.com/hetznercloud/hcloud-go/v2/hcloud/exp/actionutil"
)

func main() {
	ctx := context.Background()

	client := hcloud.NewClient(
		hcloud.WithToken("token"),
		hcloud.WithApplication("my-tool", "v1.0.0"),
	)

	result, _, err := client.Server.Create(ctx, hcloud.ServerCreateOpts{
		Name:       "Foo",
		Image:      &hcloud.Image{Name: "ubuntu-24.0"},
		ServerType: &hcloud.ServerType{Name: "cpx22"},
		Location:   &hcloud.Location{Name: "hel1"},
	})
	if err != nil {
		log.Fatalf("error creating server: %s\n", err)
	}

	// Always await any returned actions, to make sure the async process is completed before you use the result:
	err = client.Action.WaitFor(ctx, actionutil.AppendNext(result.Action, result.NextActions)...)
	if err != nil {
		log.Fatalf("error creating server: %s\n", err)
	}

	server, _, err := client.Server.GetByID(ctx, result.Server.ID)
	if err != nil {
		log.Fatalf("error retrieving server: %s\n", err)
	}
	if server != nil {
		fmt.Printf("server is called %q\n", server.Name) // prints 'server is called "Foo"'
	} else {
		fmt.Println("server not found")
	}
}
```

## Experimental features

Experimental features are published as part of our regular releases (e.g. a product
public beta). During an experimental phase, breaking changes on those features may occur
within minor releases.

The stability of experimental features is not related to the stability of its upstream API.

Experimental features have different levels of maturity (e.g. experimental, alpha, beta)
based on the maturity of the upstream API.

While experimental features will be announced in the release notes, you can also find
whether a struct or function is experimental in its Go code comment:

```go
// Experimental: $PRODUCT is $MATURITY, breaking changes may occur within minor releases.
// See https://docs.hetzner.cloud/changelog#$SLUG for more details.
```

## Upgrading

### Support

- `v2` is actively maintained by Hetzner Cloud
- `v1` is unsupported since February 2025.

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

## Development

### Experimental Features

When adding an experimental feature:

1. Add the marker comment above the declaration:

   ```go
   // Experimental: $PRODUCT is $MATURITY, breaking changes may occur within minor releases.
   ```

2. Include a link to the changelog entry:

   ```go
   // See https://docs.hetzner.cloud/changelog#slug for more details.
   ```

3. Add an announcement to the release notes.

Example:

```go
// String returns a pointer to the passed string s.
//
// Experimental: Product is beta, breaking changes may occur within minor releases.
// See https://docs.hetzner.cloud/changelog#slug for more details.
func String(s string) *string { return Ptr(s) }
```

## License

MIT license
