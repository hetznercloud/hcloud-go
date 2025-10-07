# Changelog

## [v2.27.0](https://github.com/hetznercloud/hcloud-go/releases/tag/v2.27.0)

### DNS API Beta

This release adds support for the new [DNS API](https://docs.hetzner.cloud/reference/cloud#dns).

The DNS API is currently in **beta**, which will likely end on 10 November 2025. After the beta ended, it will no longer be possible to create new zones in the old DNS system. See the [DNS Beta FAQ](https://docs.hetzner.com/networking/dns/faq/beta) for more details.

Future minor releases of this project may include breaking changes for features that are related to the DNS API.

See the [DNS API Beta changelog](https://docs.hetzner.cloud/changelog#2025-10-07-dns-beta) for more details.

**Examples**

```go
result, _, err := client.Zone.Create(ctx, hcloud.ZoneCreateOpts{
	Name:   "example.com",
	Mode:   hcloud.ZoneModePrimary,
	Labels: map[string]string{"key": "value"},
	RRSets: []hcloud.ZoneCreateOptsRRSet{
		{
			Name: "@",
			Type: hcloud.ZoneRRSetTypeA,
			Records: []hcloud.ZoneRRSetRecord{
				{Value: "201.180.75.2", Comment: "server1"},
			},
		},
	},
})

err = client.Action.WaitFor(ctx, result.Action)
zone = result.Zone
```

### Features

- support the new DNS API (#740)

## [v2.26.0](https://github.com/hetznercloud/hcloud-go/releases/tag/v2.26.0)

### Features

- support for go1.25 and drop go1.23 (#738)

### Bug Fixes

- **exp**: remove dots from deprecation messages (#736)

## [v2.25.1](https://github.com/hetznercloud/hcloud-go/releases/tag/v2.25.1)

### Bug Fixes

- **exp**: improve deprecation message helpers (#734)

## [v2.25.0](https://github.com/hetznercloud/hcloud-go/releases/tag/v2.25.0)

[Server Types](https://docs.hetzner.cloud/reference/cloud#server-types) now depend on [Locations](https://docs.hetzner.cloud/reference/cloud#locations).

- We added a new `locations` property to the [Server Types](https://docs.hetzner.cloud/reference/cloud#server-types) resource. The new property defines a list of supported [Locations](https://docs.hetzner.cloud/reference/cloud#locations) and additional per [Locations](https://docs.hetzner.cloud/reference/cloud#locations) details such as deprecations information.

- We deprecated the `deprecation` property from the [Server Types](https://docs.hetzner.cloud/reference/cloud#server-types) resource. The property will gradually be phased out as per [Locations](https://docs.hetzner.cloud/reference/cloud#locations) deprecations are being announced. Please use the new per [Locations](https://docs.hetzner.cloud/reference/cloud#locations) deprecation information instead.

See our [changelog](https://docs.hetzner.cloud/changelog#2025-09-24-per-location-server-types) for more details.

#### Upgrading

```go
// Before
func ValidateServerType(serverType *hcloud.ServerType) error {
	if serverType.IsDeprecated() {
		return fmt.Errorf("server type %s is deprecated", serverType.Name)
	}
	return nil
}
```

```go
// After
func ValidateServerType(serverType *hcloud.ServerType, location *hcloud.Location) error {
	serverTypeLocationIndex := slices.IndexFunc(serverType.Locations, func(e hcloud.ServerTypeLocation) bool {
		return e.Location.Name == location.Name
	})
	if serverTypeLocationIndex < 0 {
		return fmt.Errorf("server type %s is not supported in location %q", serverType.Name, location.Name)
	}

	if serverType.Locations[serverTypeLocationIndex].IsDeprecated() {
		return fmt.Errorf("server type %q is deprecated in location %q", serverType.Name, location.Name)
	}

	return nil
}
```

### Features

- **exp**: add `sliceutil.Transform` function (#731)
- per locations server types (#730)

## [v2.24.0](https://github.com/hetznercloud/hcloud-go/releases/tag/v2.24.0)

### Features

- add new ip_range param to ServerAttachToNetwork (#723)
- add new ip_range param to LoadBalancerAttachToNetwork (#724)

## [v2.23.0](https://github.com/hetznercloud/hcloud-go/releases/tag/v2.23.0)

### Features

- expose account currency information (#695)
- add category property to server type (#717)

### Bug Fixes

- **primary-ip**: labels not clearable (#699)

## [v2.22.0](https://github.com/hetznercloud/hcloud-go/releases/tag/v2.22.0)

### Features

- **firewall**: add `applied_to_resources` property (#667)
- add missing error codes (#671)

## [v2.21.1](https://github.com/hetznercloud/hcloud-go/releases/tag/v2.21.1)

### Bug Fixes

- add arg name in missing or invalid argument errors (#652)

## [v2.21.0](https://github.com/hetznercloud/hcloud-go/releases/tag/v2.21.0)

### Features

- consistent argument errors (#637)

### Bug Fixes

- http transport ignored when using instrumentation (#642)

## [v2.20.1](https://github.com/hetznercloud/hcloud-go/releases/tag/v2.20.1)

### Bug Fixes

- panic when a request did not set the OpPath for instrumentation (#635)

## [v2.20.0](https://github.com/hetznercloud/hcloud-go/releases/tag/v2.20.0)

In this release, the `api_endpoint` metric labels changed for all our API requests. Please make sure to update your setup if you were relying on them. The new labels are now easier to understand, see the example below:

- the path `/volumes/64314930` now has the label `/volumes/-` instead of `/volumes/`
- the path `/volumes/64314930/actions/resize` now has the label `/volumes/-/actions/resize` instead of `/volumes/actions/resize`

### Features

- support go v1.24 (#604)
- drop go v1.21 (#604)
- **exp**: remove sliceutil package (#610)
- drop go v1.22 (#602) (#621)
- redefine `api_endpoint` metric labels (#626)

### Bug Fixes

- request in batches of 25 when waiting for actions (#611)
- missing response from requests return values (#613)
- move primary ip client request/response structs to schema package (#633)

## [2.19.1](https://github.com/hetznercloud/hcloud-go/compare/v2.19.0...v2.19.1) (2025-02-03)

### Bug Fixes

- ignore nil backoff function when configuring poll or retry options ([#587](https://github.com/hetznercloud/hcloud-go/issues/587)) ([8d1b665](https://github.com/hetznercloud/hcloud-go/commit/8d1b665b88068a1956ec9f4f582b197cc7239ca7))

## [2.19.0](https://github.com/hetznercloud/hcloud-go/compare/v2.18.0...v2.19.0) (2025-01-30)

### Features

- add deprecated_api_endpoint error code ([#582](https://github.com/hetznercloud/hcloud-go/issues/582)) ([de07960](https://github.com/hetznercloud/hcloud-go/commit/de079601fdde0c240f8a8360c48facdc9f874043))
- update ActionClient for actions list removal ([#583](https://github.com/hetznercloud/hcloud-go/issues/583)) ([cf2d426](https://github.com/hetznercloud/hcloud-go/commit/cf2d4268758241eff30f5ca714372e2452c0a648))

## [2.18.0](https://github.com/hetznercloud/hcloud-go/compare/v2.17.1...v2.18.0) (2025-01-22)

### Features

- **exp:** add `randutil.GenerateID` helper ([#580](https://github.com/hetznercloud/hcloud-go/issues/580)) ([9515be7](https://github.com/hetznercloud/hcloud-go/commit/9515be7c2b44e5656bd756837e5c89ccd19665f1))

### Bug Fixes

- allow getting resources with number as name ([#571](https://github.com/hetznercloud/hcloud-go/issues/571)) ([b743a33](https://github.com/hetznercloud/hcloud-go/commit/b743a332578904748129240974142b3f6ea20845))

## [2.17.1](https://github.com/hetznercloud/hcloud-go/compare/v2.17.0...v2.17.1) (2024-12-17)

### Bug Fixes

- missing load-balancers property in network schema ([#568](https://github.com/hetznercloud/hcloud-go/issues/568)) ([103cf5e](https://github.com/hetznercloud/hcloud-go/commit/103cf5e3392b424c6ee16e9f590d8082be6977b7))

## [2.17.0](https://github.com/hetznercloud/hcloud-go/compare/v2.16.0...v2.17.0) (2024-11-20)

### Features

- **exp:** add mock requests after mock server creation ([#558](https://github.com/hetznercloud/hcloud-go/issues/558)) ([f9879fd](https://github.com/hetznercloud/hcloud-go/commit/f9879fdd5c610e76f9170ac876b65dfd38033867))

### Bug Fixes

- do not send primary IPs ID opts to the API ([#552](https://github.com/hetznercloud/hcloud-go/issues/552)) ([4e70340](https://github.com/hetznercloud/hcloud-go/commit/4e7034075363baa598a94c07083aa2a0d2779701))

## [2.16.0](https://github.com/hetznercloud/hcloud-go/compare/v2.15.0...v2.16.0) (2024-11-05)

### Features

- use custom IDOrName type for schemas ([#545](https://github.com/hetznercloud/hcloud-go/issues/545)) ([1d97017](https://github.com/hetznercloud/hcloud-go/commit/1d97017b303caa29379e7202a04814985604bea2))

### Bug Fixes

- **metadata:** trim responses before parsing ([#547](https://github.com/hetznercloud/hcloud-go/issues/547)) ([283363f](https://github.com/hetznercloud/hcloud-go/commit/283363f2c875f6cf7611d99a8f1590bfab157af6)), closes [#540](https://github.com/hetznercloud/hcloud-go/issues/540)

## [2.15.0](https://github.com/hetznercloud/hcloud-go/compare/v2.14.0...v2.15.0) (2024-10-31)

### Features

- **exp:** add mockutil.Server helper ([#543](https://github.com/hetznercloud/hcloud-go/issues/543)) ([fa1069b](https://github.com/hetznercloud/hcloud-go/commit/fa1069b1f677325726bf1f9dc14922a707e16440))
- **exp:** fail test when mock calls were expected but not received ([#542](https://github.com/hetznercloud/hcloud-go/issues/542)) ([01392cc](https://github.com/hetznercloud/hcloud-go/commit/01392cc77e05d1aa1add84e17aa89a8a26ea61d2))

## [2.14.0](https://github.com/hetznercloud/hcloud-go/compare/v2.13.1...v2.14.0) (2024-10-21)

### Features

- add support for Go 1.23 ([#532](https://github.com/hetznercloud/hcloud-go/issues/532)) ([838f878](https://github.com/hetznercloud/hcloud-go/commit/838f878189bb46fc071066e77bac421880d9b61e))
- allow retrieving the action from an action error ([#538](https://github.com/hetznercloud/hcloud-go/issues/538)) ([07727d3](https://github.com/hetznercloud/hcloud-go/commit/07727d3362ae3ebd0966ff09e5616afe9965b17f))
- include action ID in action error string ([#539](https://github.com/hetznercloud/hcloud-go/issues/539)) ([ad5417f](https://github.com/hetznercloud/hcloud-go/commit/ad5417f9086278521937ccbae00f31a2b3c8421e))

### Bug Fixes

- deprecate pricing floating ip field ([#523](https://github.com/hetznercloud/hcloud-go/issues/523)) ([1089d40](https://github.com/hetznercloud/hcloud-go/commit/1089d40577b60caeeffcfe30940f831ea7ce3999))
- deprecate unused window parameter in the `EnableBackup` call ([#531](https://github.com/hetznercloud/hcloud-go/issues/531)) ([584f6c2](https://github.com/hetznercloud/hcloud-go/commit/584f6c2a5091ebb2f4761f752c22d9562d46a6f6)), closes [#525](https://github.com/hetznercloud/hcloud-go/issues/525)

## [2.13.1](https://github.com/hetznercloud/hcloud-go/compare/v2.13.0...v2.13.1) (2024-08-09)

### Bug Fixes

- `invalid_input` API errors may not return details ([#507](https://github.com/hetznercloud/hcloud-go/issues/507)) ([ca78af2](https://github.com/hetznercloud/hcloud-go/commit/ca78af2af8acd375460c0f9705ea9d62d5ee1cc4))

## [2.13.0](https://github.com/hetznercloud/hcloud-go/compare/v2.12.0...v2.13.0) (2024-08-06)

### Features

- **network:** add new network zone ap-southeast ([#501](https://github.com/hetznercloud/hcloud-go/issues/501)) ([a79a06b](https://github.com/hetznercloud/hcloud-go/commit/a79a06b0dda182e968a7a6e5cef9a4181414f29e))

### Bug Fixes

- update `NetworkSubnetType` constants ([#499](https://github.com/hetznercloud/hcloud-go/issues/499)) ([ce497fe](https://github.com/hetznercloud/hcloud-go/commit/ce497fe89ccd9cccd8dc84ccd854427484fdd031))

## [2.12.0](https://github.com/hetznercloud/hcloud-go/compare/v2.11.0...v2.12.0) (2024-07-25)

### API Changes for Traffic Prices and Server Type Included Traffic

There will be a breaking change in the API regarding Traffic Prices and Server Type Included Traffic on 2024-08-05. This release marks the affected fields as `Deprecated`. Please check if this affects any of your code and switch to the replacement fields where necessary.

You can learn more about this change in [our changelog](https://docs.hetzner.cloud/changelog#2024-07-25-cloud-api-returns-traffic-information-in-different-format).

#### Upgrading

##### Server Type Included Traffic

If you were using the field `hcloud.ServerType.IncludedTraffic`, you can now get the information through `hcloud.ServerType.Pricings`:

```go
func main() {
// previous
includedTraffic := serverType.IncludedTraffic

    // now
    locationOfInterest := "fsn1"
    var includedTraffic uint64
    for _, price := range serverType.Pricings {
        if price.Location.Name == locationOfInterest {
            includedTraffic = price.IncludedTraffic
            break
        }
    }
}
```

##### Traffic Prices

If you were using the field `hcloud.Pricing.Traffic`, you can now get the information through `hcloud.Pricing.ServerTypes` or `hcloud.Pricing.LoadBalancerTypes`:

```go
func main() {
// previous
trafficPrice := pricing.Traffic

    // now
    serverTypeOfInterest := "cx22"
    locationOfInterest := "fsn1"

    var trafficPrice hcloud.Price
    for _, serverTypePricings := range pricing.ServerTypes {
        if serverTypePricings.ServerType.Name == serverTypeOfInterest {
            for _, price := range serverTypePricings {
               if price.Location.Name == locationOfInterest {
                   trafficPrice = price.PerTBTraffic
                   break
               }
            }
        }
    }
}
```

### Features

- add jitter in the client default retry exponential backoff ([#492](https://github.com/hetznercloud/hcloud-go/issues/492)) ([6205076](https://github.com/hetznercloud/hcloud-go/commit/6205076b89350bdbf08bc6b771a1d1267a3ac422))
- add new `WithPollOpts` client option ([#493](https://github.com/hetznercloud/hcloud-go/issues/493)) ([2c1a2d6](https://github.com/hetznercloud/hcloud-go/commit/2c1a2d65596bcbe282ff004c1a9da89950e754df))
- allow checking multiple errors codes in `IsError` ([#491](https://github.com/hetznercloud/hcloud-go/issues/491)) ([af59ab8](https://github.com/hetznercloud/hcloud-go/commit/af59ab846665abd735c9717eb2a47c0a8c79776d))
- **load-balancer-type:** new traffic price fields ([94e0f44](https://github.com/hetznercloud/hcloud-go/commit/94e0f44d269fdb5138485a6e69dae9105690e4b0))
- **pricing:** mark traffic field as deprecated ([94e0f44](https://github.com/hetznercloud/hcloud-go/commit/94e0f44d269fdb5138485a6e69dae9105690e4b0))
- **server-type:** mark included traffic field as deprecated ([94e0f44](https://github.com/hetznercloud/hcloud-go/commit/94e0f44d269fdb5138485a6e69dae9105690e4b0))
- **server-type:** new traffic price fields ([94e0f44](https://github.com/hetznercloud/hcloud-go/commit/94e0f44d269fdb5138485a6e69dae9105690e4b0))

## [2.11.0](https://github.com/hetznercloud/hcloud-go/compare/v2.10.2...v2.11.0) (2024-07-23)

### Features

- add truncated exponential backoff with full jitter ([#459](https://github.com/hetznercloud/hcloud-go/issues/459)) ([fd1f46c](https://github.com/hetznercloud/hcloud-go/commit/fd1f46cc35e61dde1e524399eef88c38a757636e))
- allow configuring retry options ([#488](https://github.com/hetznercloud/hcloud-go/issues/488)) ([2db9575](https://github.com/hetznercloud/hcloud-go/commit/2db95753e2c826aeafa3bd9b864e95efd89ace7f))
- **exp:** add sliceutil package ([#489](https://github.com/hetznercloud/hcloud-go/issues/489)) ([f4ad6bc](https://github.com/hetznercloud/hcloud-go/commit/f4ad6bc93ff5017dda1b71a1606b67a79b56eb57))
- **exp:** rename `*utils` package to `*util` ([#487](https://github.com/hetznercloud/hcloud-go/issues/487)) ([19da475](https://github.com/hetznercloud/hcloud-go/commit/19da4759f4cbee7ed94ed6996350b45650f8b0b9))
- respect cancelled contexts during retry sleep ([#470](https://github.com/hetznercloud/hcloud-go/issues/470)) ([756f605](https://github.com/hetznercloud/hcloud-go/commit/756f605c97ac570adec531a479fc61c1ed27ab72))
- retry requests when the api gateway errors ([#470](https://github.com/hetznercloud/hcloud-go/issues/470)) ([756f605](https://github.com/hetznercloud/hcloud-go/commit/756f605c97ac570adec531a479fc61c1ed27ab72))
- retry requests when the network timed out ([#470](https://github.com/hetznercloud/hcloud-go/issues/470)) ([756f605](https://github.com/hetznercloud/hcloud-go/commit/756f605c97ac570adec531a479fc61c1ed27ab72))
- retry requests when the rate limit was reached ([#470](https://github.com/hetznercloud/hcloud-go/issues/470)) ([756f605](https://github.com/hetznercloud/hcloud-go/commit/756f605c97ac570adec531a479fc61c1ed27ab72))

### Bug Fixes

- **exp:** set capacity for each batch ([#490](https://github.com/hetznercloud/hcloud-go/issues/490)) ([57f53c1](https://github.com/hetznercloud/hcloud-go/commit/57f53c1dca54ee79a33bab35740e7de7ece3b75f))

## [2.10.2](https://github.com/hetznercloud/hcloud-go/compare/v2.10.1...v2.10.2) (2024-06-26)

### Bug Fixes

- **exp:** allow request path matching in the want function ([#475](https://github.com/hetznercloud/hcloud-go/issues/475)) ([267879b](https://github.com/hetznercloud/hcloud-go/commit/267879b78989ae870d581e9c929105ff76c60fb0))

## [2.10.1](https://github.com/hetznercloud/hcloud-go/compare/v2.10.0...v2.10.1) (2024-06-25)

### Bug Fixes

- **exp:** configure response headers before sending them ([#473](https://github.com/hetznercloud/hcloud-go/issues/473)) ([07d4a35](https://github.com/hetznercloud/hcloud-go/commit/07d4a356dec0e2f44a6f3eec1ea3affec8932c22))

## [2.10.0](https://github.com/hetznercloud/hcloud-go/compare/v2.9.0...v2.10.0) (2024-06-25)

### Features

- **exp:** add envutils package ([#466](https://github.com/hetznercloud/hcloud-go/issues/466)) ([a7636bd](https://github.com/hetznercloud/hcloud-go/commit/a7636bdcf5e4d55860f40da684c64cb72f8ddc03))
- **exp:** add labelutils with selector ([#465](https://github.com/hetznercloud/hcloud-go/issues/465)) ([1a55a7e](https://github.com/hetznercloud/hcloud-go/commit/1a55a7ed65bffdbf73000ffcf1ef22c9e55650f8))
- **exp:** add mock utils package ([#460](https://github.com/hetznercloud/hcloud-go/issues/460)) ([92f7c62](https://github.com/hetznercloud/hcloud-go/commit/92f7c624edfc76e06abe8c9c60e9c78d4b28f12f))
- **exp:** rename `AppendNextActions` to `AppendNext` ([#452](https://github.com/hetznercloud/hcloud-go/issues/452)) ([9b6239a](https://github.com/hetznercloud/hcloud-go/commit/9b6239ad188e601f22bc28e0072603c07fea201c))
- filter out nil actions in action waiter ([#464](https://github.com/hetznercloud/hcloud-go/issues/464)) ([4fc9a40](https://github.com/hetznercloud/hcloud-go/commit/4fc9a4039d45071124a435121642ca396a8237c0))

### Bug Fixes

- nil check against the embedded `http.Response` ([#469](https://github.com/hetznercloud/hcloud-go/issues/469)) ([46e489a](https://github.com/hetznercloud/hcloud-go/commit/46e489a1782e8477d1c5a234dc203fa356c2a583))

## [2.9.0](https://github.com/hetznercloud/hcloud-go/compare/v2.8.0...v2.9.0) (2024-05-29)

### Features

- **exp:** add `AppendNextActions` function ([#440](https://github.com/hetznercloud/hcloud-go/issues/440)) ([b07d7ad](https://github.com/hetznercloud/hcloud-go/commit/b07d7adf0bb08dd372f1f6ff630d16af8cb1265c))
- **exp:** add ssh key functions ([#441](https://github.com/hetznercloud/hcloud-go/issues/441)) ([d766e96](https://github.com/hetznercloud/hcloud-go/commit/d766e96e2220b069a0a8067a3249afe48cf68b2e))

### Bug Fixes

- **exp:** rename to `sshutils` package name ([#450](https://github.com/hetznercloud/hcloud-go/issues/450)) ([6d4100d](https://github.com/hetznercloud/hcloud-go/commit/6d4100dbceb9a43a66bad3ce91cee039198130c3))

## [2.8.0](https://github.com/hetznercloud/hcloud-go/compare/v2.7.2...v2.8.0) (2024-05-06)

### Features

- **error:** show internal correlation id in error messages ([#411](https://github.com/hetznercloud/hcloud-go/issues/411)) ([6c96d19](https://github.com/hetznercloud/hcloud-go/commit/6c96d19dacde736b52abd7b8fc8879c8f721f23b))
- implement actions waiter ([#407](https://github.com/hetznercloud/hcloud-go/issues/407)) ([1e3fa70](https://github.com/hetznercloud/hcloud-go/commit/1e3fa7033d8a1cea1c0a9007a6255798979f0771))
- require Go &gt;= 1.21 ([#424](https://github.com/hetznercloud/hcloud-go/issues/424)) ([d4f4000](https://github.com/hetznercloud/hcloud-go/commit/d4f40009016c3ed5bb14ef9ae16bcf6aefa90fee))

### Bug Fixes

- improve error message format with correlation id ([#430](https://github.com/hetznercloud/hcloud-go/issues/430)) ([013477f](https://github.com/hetznercloud/hcloud-go/commit/013477f4227bdc69f376d8f13a875b09c32171f6))

## [2.7.2](https://github.com/hetznercloud/hcloud-go/compare/v2.7.1...v2.7.2) (2024-04-19)

### Bug Fixes

- creating firewall with rules not working correctly ([#412](https://github.com/hetznercloud/hcloud-go/issues/412)) ([43b2bab](https://github.com/hetznercloud/hcloud-go/commit/43b2bab8c82cab276c07e9deb9f4e422e10ccc82))

## [2.7.1](https://github.com/hetznercloud/hcloud-go/compare/v2.7.0...v2.7.1) (2024-04-18)

### Bug Fixes

- setting firewall rules not working correctly ([#405](https://github.com/hetznercloud/hcloud-go/issues/405)) ([16daea0](https://github.com/hetznercloud/hcloud-go/commit/16daea0dfb32a5e2b8a71fce88d1c897cd4134eb))

## [2.7.0](https://github.com/hetznercloud/hcloud-go/compare/v2.6.0...v2.7.0) (2024-03-27)

### Features

- add volume format property ([#397](https://github.com/hetznercloud/hcloud-go/issues/397)) ([c0940af](https://github.com/hetznercloud/hcloud-go/commit/c0940afce9eb01c0e6838502c91aa569ab411a03))
- **error:** handle wrapped errors in IsError() ([#374](https://github.com/hetznercloud/hcloud-go/issues/374)) ([83df108](https://github.com/hetznercloud/hcloud-go/commit/83df108228519617a919470e5ebbb4a1a2587b34))
- require Go &gt;= 1.20 ([#392](https://github.com/hetznercloud/hcloud-go/issues/392)) ([299f181](https://github.com/hetznercloud/hcloud-go/commit/299f181c469a48e977743f2587e232a293cd9db5))
- **server:** deprecate ServerRescueTypeLinux32 ([#378](https://github.com/hetznercloud/hcloud-go/issues/378)) ([2f334c3](https://github.com/hetznercloud/hcloud-go/commit/2f334c3de2aecfb7aaf179424bb9a1a7b0da53ea))
- test with Go 1.22 ([#391](https://github.com/hetznercloud/hcloud-go/issues/391)) ([49be506](https://github.com/hetznercloud/hcloud-go/commit/49be50664f56e4e315e687543ee0bf7fbb2af186))

### Bug Fixes

- falsely omitted fields in firewall schema ([#396](https://github.com/hetznercloud/hcloud-go/issues/396)) ([a3509b6](https://github.com/hetznercloud/hcloud-go/commit/a3509b6b0f2762b5d1e27f2374d3fb5fb53206e2))
- LoadBalancerUpdateServiceOpts not converted correctly ([#394](https://github.com/hetznercloud/hcloud-go/issues/394)) ([0f187ce](https://github.com/hetznercloud/hcloud-go/commit/0f187cef1f568d87a32b2ec063847dd9ff314740))
- primary ip assignee id not mapped to nil ([#395](https://github.com/hetznercloud/hcloud-go/issues/395)) ([b5fea38](https://github.com/hetznercloud/hcloud-go/commit/b5fea38f5d9d62c88af470ba69195a73ee075f09))

## [2.6.0](https://github.com/hetznercloud/hcloud-go/compare/v2.5.1...v2.6.0) (2024-01-09)

### Features

- alias deprecated field to deprecation info struct ([#371](https://github.com/hetznercloud/hcloud-go/issues/371)) ([e961be9](https://github.com/hetznercloud/hcloud-go/commit/e961be9615452fc63c4c71e6561d4e86f8e8e95a))
- **instrumentation:** allow passing in any prometheus.Registerer ([#369](https://github.com/hetznercloud/hcloud-go/issues/369)) ([0821c07](https://github.com/hetznercloud/hcloud-go/commit/0821c078900910fa9e3ca6c6c0af48a73f00c7c6))

## [2.5.1](https://github.com/hetznercloud/hcloud-go/compare/v2.5.0...v2.5.1) (2023-12-13)

### Bug Fixes

- schema conversion outputs debug messages to stdout ([#354](https://github.com/hetznercloud/hcloud-go/issues/354)) ([ade8fbd](https://github.com/hetznercloud/hcloud-go/commit/ade8fbd60a88a648c95391f00cfe3ccc09be8f37))

## [2.5.0](https://github.com/hetznercloud/hcloud-go/compare/v2.4.0...v2.5.0) (2023-12-12)

### Features

- add conversion methods from schema to hcloud objects ([#343](https://github.com/hetznercloud/hcloud-go/issues/343)) ([6feda4d](https://github.com/hetznercloud/hcloud-go/commit/6feda4d9b0e7cf3f5a17a4b38504abbe5213883d))
- add interfaces for client structs ([#342](https://github.com/hetznercloud/hcloud-go/issues/342)) ([4f9390f](https://github.com/hetznercloud/hcloud-go/commit/4f9390f8387d1c86330156adbb6801aacba7a8f0))
- add missing properties ([#349](https://github.com/hetznercloud/hcloud-go/issues/349)) ([c8a28d0](https://github.com/hetznercloud/hcloud-go/commit/c8a28d0dbf0c84364401282a60b44ccea1da6423))
- **error:** include http response in api errors ([#320](https://github.com/hetznercloud/hcloud-go/issues/320)) ([9558239](https://github.com/hetznercloud/hcloud-go/commit/95582395dfb1039f4ce4f10a1ac9c068db93a867))

### Bug Fixes

- make schemas consistent with API ([#348](https://github.com/hetznercloud/hcloud-go/issues/348)) ([b0d7055](https://github.com/hetznercloud/hcloud-go/commit/b0d7055543669fb96af1726daa7c1458fb1b65a2))

## [2.4.0](https://github.com/hetznercloud/hcloud-go/compare/v2.3.0...v2.4.0) (2023-10-12)

### Features

- **iso:** support deprecation info API ([#331](https://github.com/hetznercloud/hcloud-go/issues/331)) ([b3a3621](https://github.com/hetznercloud/hcloud-go/commit/b3a36214b21ab1c5add94c2b1896a995342757b6))

## [2.3.0](https://github.com/hetznercloud/hcloud-go/compare/v2.2.0...v2.3.0) (2023-10-04)

### Features

- add error code unauthorized ([#313](https://github.com/hetznercloud/hcloud-go/issues/313)) ([b77d9e0](https://github.com/hetznercloud/hcloud-go/commit/b77d9e04ca903448cc1a22c242f440e67a81a028))
- test with Go 1.21 ([#319](https://github.com/hetznercloud/hcloud-go/issues/319)) ([7ddb2ec](https://github.com/hetznercloud/hcloud-go/commit/7ddb2ec057d0e165abb5b8cbd74e95d5bb4add49))

### Bug Fixes

- ensure the pollBackoffFunc is correctly set ([#322](https://github.com/hetznercloud/hcloud-go/issues/322)) ([2b2f869](https://github.com/hetznercloud/hcloud-go/commit/2b2f8697aa09f67dd508c33fe36f7a59b9d3f192))

## [2.2.0](https://github.com/hetznercloud/hcloud-go/compare/v2.1.1...v2.2.0) (2023-08-24)

### Features

- support resource-specific Action endpoints ([#295](https://github.com/hetznercloud/hcloud-go/issues/295)) ([ddc2ac4](https://github.com/hetznercloud/hcloud-go/commit/ddc2ac45489c48a7563806c425222236ab1f8aa0))

## [2.1.1](https://github.com/hetznercloud/hcloud-go/compare/v2.1.0...v2.1.1) (2023-08-15)

### Bug Fixes

- send more precise progress values ([#304](https://github.com/hetznercloud/hcloud-go/issues/304)) ([867aa63](https://github.com/hetznercloud/hcloud-go/commit/867aa632521ad3acfb04beb52b6307330740fc68))

## [2.1.0](https://github.com/hetznercloud/hcloud-go/compare/v2.0.0...v2.1.0) (2023-08-08)

### Features

- **metadata:** add timeout option ([#293](https://github.com/hetznercloud/hcloud-go/issues/293)) ([913bf74](https://github.com/hetznercloud/hcloud-go/commit/913bf74071e03a9c79fd0f8a5a37b1e11b350ae1))

### Bug Fixes

- **action:** unexpected behaviour when watching non-existing Actions ([#298](https://github.com/hetznercloud/hcloud-go/issues/298)) ([0727d42](https://github.com/hetznercloud/hcloud-go/commit/0727d42e26a8112923a84b84bb506c980e07262d))
- **instrumentation:** multiple instrumented clients cause panic ([#289](https://github.com/hetznercloud/hcloud-go/issues/289)) ([c0ef9b6](https://github.com/hetznercloud/hcloud-go/commit/c0ef9b6e6e3f36d8c2282c2b7aa9d8687141f291))

## [2.0.0](https://github.com/hetznercloud/hcloud-go/compare/v1.47.0...v2.0.0) (2023-07-12)

### ⚠ BREAKING CHANGES

- use int64 for ID fields ([#282](https://github.com/hetznercloud/hcloud-go/issues/282))

### Features

- make ListOpts.Values method public ([#285](https://github.com/hetznercloud/hcloud-go/issues/285)) ([c82ea59](https://github.com/hetznercloud/hcloud-go/commit/c82ea5971d94b812fd9302bfc4151f4ebfa43413))
- use int64 for ID fields ([#282](https://github.com/hetznercloud/hcloud-go/issues/282)) ([359c389](https://github.com/hetznercloud/hcloud-go/commit/359c3894641f2dcca4a049537e256a20853b5ad9))

### Bug Fixes

- **action:** show accurate progress in WatchOverallProgress ([#281](https://github.com/hetznercloud/hcloud-go/issues/281)) ([cae9e57](https://github.com/hetznercloud/hcloud-go/commit/cae9e5789b20bdb7d9213ba88897435ee1abce86))
- **iso:** invalid field include_wildcard_architecture ([188b68a](https://github.com/hetznercloud/hcloud-go/commit/188b68ad674066d302b2614432cf4f7d5b47f41a))

## [1.47.0](https://github.com/hetznercloud/hcloud-go/compare/v1.46.1...v1.47.0) (2023-06-21)

### Features

- **network:** add field expose_routes_to_vswitch ([#277](https://github.com/hetznercloud/hcloud-go/issues/277)) ([e73c52d](https://github.com/hetznercloud/hcloud-go/commit/e73c52dce563e00c8ccba7528fd054936a52b64c))

## [1.46.1](https://github.com/hetznercloud/hcloud-go/compare/v1.46.0...v1.46.1) (2023-06-16)

### Bug Fixes

- adjust label validation for max length of 63 characters ([#273](https://github.com/hetznercloud/hcloud-go/issues/273)) ([6382808](https://github.com/hetznercloud/hcloud-go/commit/63828086e89413115f4f5ee326835f93752cbb51))
- **deps:** update module golang.org/x/net to v0.11.0 ([#258](https://github.com/hetznercloud/hcloud-go/issues/258)) ([7918f21](https://github.com/hetznercloud/hcloud-go/commit/7918f2172ac000f72a6ed52a2c69ae03c7b6513f))

## [1.46.0](https://github.com/hetznercloud/hcloud-go/compare/v1.45.1...v1.46.0) (2023-06-13)

### Features

- provide `.AllWithOpts` method for all clients ([#266](https://github.com/hetznercloud/hcloud-go/issues/266)) ([2a7249e](https://github.com/hetznercloud/hcloud-go/commit/2a7249ed646bf9b5d91890fc6698b04bfdaf7806))
- **servertype:** implement new Deprecation api field ([#268](https://github.com/hetznercloud/hcloud-go/issues/268)) ([ac5ae2e](https://github.com/hetznercloud/hcloud-go/commit/ac5ae2e80c361775bd14c776e23ccec4ce5849e7))

## [1.45.1](https://github.com/hetznercloud/hcloud-go/compare/v1.45.0...v1.45.1) (2023-05-11)

### Bug Fixes

- **servertype:** use int64 to fit TB sizes on 32-bit platforms ([#261](https://github.com/hetznercloud/hcloud-go/issues/261)) ([2b19245](https://github.com/hetznercloud/hcloud-go/commit/2b1924575148a8675de7ca65f5578aeb70ef750f))

## [1.45.0](https://github.com/hetznercloud/hcloud-go/compare/v1.44.0...v1.45.0) (2023-05-11)

### Features

- **servertype:** add field for included traffic ([#259](https://github.com/hetznercloud/hcloud-go/issues/259)) ([d3b012a](https://github.com/hetznercloud/hcloud-go/commit/d3b012a678ee7012a54bb6088c3ee3b3efb5978d))

## [1.44.0](https://github.com/hetznercloud/hcloud-go/compare/v1.43.0...v1.44.0) (2023-05-05)

### Features

- **iso:** extend ISOClient by AllWithOpts method ([#254](https://github.com/hetznercloud/hcloud-go/issues/254)) ([c42f69b](https://github.com/hetznercloud/hcloud-go/commit/c42f69b05bf732b92561285874e4ac2b6cc01d1c))

### Bug Fixes

- **deps:** update module github.com/prometheus/client_golang to v1.15.1 ([#253](https://github.com/hetznercloud/hcloud-go/issues/253)) ([275f0fd](https://github.com/hetznercloud/hcloud-go/commit/275f0fd4d4316fed0acddd13400d8a3143889f49))

## [1.43.0](https://github.com/hetznercloud/hcloud-go/compare/v1.42.0...v1.43.0) (2023-04-26)

### Features

- **primary-ip:** implement RDNSSupporter ([#252](https://github.com/hetznercloud/hcloud-go/issues/252)) ([41a4c5a](https://github.com/hetznercloud/hcloud-go/commit/41a4c5a1d7f70fa6a279c6a2834830e806d456de))

### Bug Fixes

- **deps:** update module github.com/prometheus/client_golang to v1.15.0 ([#250](https://github.com/hetznercloud/hcloud-go/issues/250)) ([f10e804](https://github.com/hetznercloud/hcloud-go/commit/f10e8042ac12e6195824b80a791700ce857111bc))

## [1.42.0](https://github.com/hetznercloud/hcloud-go/compare/v1.41.0...v1.42.0) (2023-04-12)

### Features

- add support for ARM APIs ([#249](https://github.com/hetznercloud/hcloud-go/issues/249)) ([ce9859f](https://github.com/hetznercloud/hcloud-go/commit/ce9859f178078c99f3d15de41c7f0266c2e885e1))

### Bug Fixes

- **deps:** update module golang.org/x/net to v0.9.0 ([#247](https://github.com/hetznercloud/hcloud-go/issues/247)) ([962afeb](https://github.com/hetznercloud/hcloud-go/commit/962afebeed76560687103777efae20ff70fbcf16))

## [1.41.0](https://github.com/hetznercloud/hcloud-go/compare/v1.40.0...v1.41.0) (2023-03-06)

### Features

- add ServerClient.RebuildWithResult to return root password ([#245](https://github.com/hetznercloud/hcloud-go/issues/245)) ([82f97cf](https://github.com/hetznercloud/hcloud-go/commit/82f97cf48695848e2569b38f8ff24bb050966ee4))

### Bug Fixes

- **deps:** update module github.com/google/go-cmp to v0.5.9 ([#237](https://github.com/hetznercloud/hcloud-go/issues/237)) ([2237ff7](https://github.com/hetznercloud/hcloud-go/commit/2237ff795cbaf1e75759cdd396b3dfe5491c0e24))
- **deps:** update module github.com/prometheus/client_golang to v1.14.0 ([#241](https://github.com/hetznercloud/hcloud-go/issues/241)) ([75a4a01](https://github.com/hetznercloud/hcloud-go/commit/75a4a0140216eb476990e50ab9b13b60881404be))
- **deps:** update module github.com/stretchr/testify to v1.8.2 ([#242](https://github.com/hetznercloud/hcloud-go/issues/242)) ([4b51f1e](https://github.com/hetznercloud/hcloud-go/commit/4b51f1e8a13f1f859211910f1dce2daebb583b04))
- **deps:** update module golang.org/x/net to v0.7.0 [security] ([#236](https://github.com/hetznercloud/hcloud-go/issues/236)) ([774a560](https://github.com/hetznercloud/hcloud-go/commit/774a560b3d167c5c55cd3cbc4f83872ecc878670))
- **deps:** update module golang.org/x/net to v0.8.0 ([#243](https://github.com/hetznercloud/hcloud-go/issues/243)) ([8ae14f3](https://github.com/hetznercloud/hcloud-go/commit/8ae14f36021a32f5bab21a74d2467aa2487b348d))

## [1.40.0](https://github.com/hetznercloud/hcloud-go/compare/v1.39.0...v1.40.0) (2023-02-08)

### Features

- **action:** use configurable backoff to wait for action progress ([#227](https://github.com/hetznercloud/hcloud-go/issues/227)) ([8da6417](https://github.com/hetznercloud/hcloud-go/commit/8da6417cf7d87bf44117ede9cd2839d7dc055f66))
- support go v1.20 and drop v1.18 ([#231](https://github.com/hetznercloud/hcloud-go/issues/231)) ([44af6e5](https://github.com/hetznercloud/hcloud-go/commit/44af6e5beade11432b5ca30575781875cbd08343))

## [1.39.0](https://github.com/hetznercloud/hcloud-go/compare/v1.38.0...v1.39.0) (2022-12-29)

### Features

- Use generics to get pointers to types ([#219](https://github.com/hetznercloud/hcloud-go/issues/219)) ([a5cd797](https://github.com/hetznercloud/hcloud-go/commit/a5cd79782dc849b3137e46ada2da6b319d4093c8))

### Bug Fixes

- deprecate PricingPrimaryIPTypePrice.Datacenter for Location ([#222](https://github.com/hetznercloud/hcloud-go/issues/222)) ([e0e5a1e](https://github.com/hetznercloud/hcloud-go/commit/e0e5a1e08fd7c0864fd94a787ee86714b5e9afc5))

## v1.38.0

### What's Changed

- feat(network): add new Network Zone us-west by @apricote in https://github.com/hetznercloud/hcloud-go/pull/217
- chore: prepare v1.38.0 by @apricote in https://github.com/hetznercloud/hcloud-go/pull/218

**Full Changelog**: https://github.com/hetznercloud/hcloud-go/compare/v1.37.0...v1.38.0

## v1.37.0

### What's Changed

- PrimaryIPClient Add AllWithOpts by @LKaemmerling in https://github.com/hetznercloud/hcloud-go/pull/214
- fix: error when updating IPv6 Primary IP by @apricote in https://github.com/hetznercloud/hcloud-go/pull/215

**Full Changelog**: https://github.com/hetznercloud/hcloud-go/compare/v1.36.0...v1.37.0

## v1.36.0

### What's Changed

- feat: add ServerClient.DeleteWithResult method by @apricote in https://github.com/hetznercloud/hcloud-go/pull/213

### New Contributors

- @apricote made their first contribution in https://github.com/hetznercloud/hcloud-go/pull/213

**Full Changelog**: https://github.com/hetznercloud/hcloud-go/compare/v1.35.3...v1.36.0

## v1.35.3

### What's Changed

- Drop support for Go < 1.17 and add official tests on go 1.19 by @LKaemmerling in https://github.com/hetznercloud/hcloud-go/pull/211
- Stop automatic retrying on RateLimitExceeded by @LKaemmerling in https://github.com/hetznercloud/hcloud-go/pull/210

**Full Changelog**: https://github.com/hetznercloud/hcloud-go/compare/v1.35.2...v1.35.3

## v1.35.2

### What's Changed

- Allow empty labels by @4ND3R50N in https://github.com/hetznercloud/hcloud-go/pull/207

**Full Changelog**: https://github.com/hetznercloud/hcloud-go/compare/v1.35.1...v1.35.2

## v1.35.1

### What's Changed

- Accept no primary IPs with server create with StartAfterCreate = false by @4ND3R50N in https://github.com/hetznercloud/hcloud-go/pull/205

**Full Changelog**: https://github.com/hetznercloud/hcloud-go/compare/v1.35.0...v1.35.1

## v1.35.0

### What's Changed

- Catch invalid token values and error out without value exposure by @NotTheEvilOne in https://github.com/hetznercloud/hcloud-go/pull/194
- Remove ServerRescueTypeFreeBSD64 by @LKaemmerling in https://github.com/hetznercloud/hcloud-go/pull/203
- Add Primary IP Support by @LKaemmerling in https://github.com/hetznercloud/hcloud-go/pull/204

### New Contributors

- @NotTheEvilOne made their first contribution in https://github.com/hetznercloud/hcloud-go/pull/194

**Full Changelog**: https://github.com/hetznercloud/hcloud-go/compare/v1.34.0...v1.35.0

## v1.34.0

### What's Changed

- Test on Go 1.18 by @LKaemmerling in https://github.com/hetznercloud/hcloud-go/pull/202
- Add support for sorting the response of all list calls by @LKaemmerling in https://github.com/hetznercloud/hcloud-go/pull/201
- Set UsePrivateIP for targets when creating a LoadBalancer by @hakman in https://github.com/hetznercloud/hcloud-go/pull/198

### New Contributors

- @hakman made their first contribution in https://github.com/hetznercloud/hcloud-go/pull/198

**Full Changelog**: https://github.com/hetznercloud/hcloud-go/compare/v1.33.2...v1.34.0

## v1.33.2

### What's Changed

- Add constant for resource locked error code by @patrickschaffrath in https://github.com/hetznercloud/hcloud-go/pull/189
- Fix metadata client error detection by @choffmeister in https://github.com/hetznercloud/hcloud-go/pull/193
- Add labels.go to validate resource labels by @4ND3R50N in https://github.com/hetznercloud/hcloud-go/pull/197

### New Contributors

- @patrickschaffrath made their first contribution in https://github.com/hetznercloud/hcloud-go/pull/189
- @choffmeister made their first contribution in https://github.com/hetznercloud/hcloud-go/pull/193
- @4ND3R50N made their first contribution in https://github.com/hetznercloud/hcloud-go/pull/197

**Full Changelog**: https://github.com/hetznercloud/hcloud-go/compare/v1.33.1...v1.33.2

## v1.33.1

### Changelog

41fef2f Add constants for new firewall error code

## v1.33.0

### What's Changed

- Add us-east network zone by @LKaemmerling in https://github.com/hetznercloud/hcloud-go/pull/187

**Full Changelog**: https://github.com/hetznercloud/hcloud-go/compare/v1.32.0...v1.33.0

## v1.32.0

### Changelog

785896c Add Metadata Client (#184)
7778143 Bump version to 1.32.0
eabe683 Make it possible to instrument the hcloud-go (#185)

## v1.31.1

### Changelog

1cde0d7 Add RDNS client (#183)

## v1.31.0

### Changelog

0656bf9 Add support for Load Balancer DNS PTRs (#182)
9218970 Bump hcloud.Version
ebf9f25 Test on/with go 1.17 (#181)

## v1.30.0

### Changelog

592b198 Add new Floating IP Pricing structure (#180)
c43897a Bump hcloud.Version
ab0ebb2 Placement groups (#179)

## v1.29.1

### Changelog

e951e99 Fix firewall rule description

## v1.29.0

### Changelog

e3eca12 Add description field to firewall rules (#177)

## v1.28.0

### Changelog

ae5e1b8 Add goreleaser (#175)
4cf75f4 Add support for app images (#176)

## v1.27.0

Add support for firewall protocols ESP and GRE (#173)

## v1.26.2

- Fix AppliedTo Field in FirewallCreateRequest
- Add `deleted` and `IsDeleted()` to `Image` (https://github.com/hetznercloud/hcloud-go/pull/172)

## v1.26.1

- Fix validation error caused by firewall label selectors (https://github.com/hetznercloud/hcloud-go/pull/171)

## v1.26.0

- Add Firewall Resource Label Selector (#169)

## v1.25.0

- Support Hetzner Cloud managed Certificates (#167)

## v1.23.1

- Add removed `ErrorCodeServerAlreadyAttached` again

## v1.23.0

- Add missing constants for all resource specific error codes
- Expose metrics for Servers and Load Balancers
- Add support for vSwitch Subnetworks

## v1.22.0

- Add `PrimaryDiskSize` Field to `Server`

## v1.21.1

- Don't send `Authorization` Header when `WithToken` was not called

## v1.21.0

- Add `IncludeDeprecated` Field to `ImageListOpts`

## v1.20.0

- Add support for Load Balancer Label Selector targets
- Add support for Load Balancer IP targets

## v1.19.0

- Fix nil pointer dereference when creating a Load Balancer with HTTP(S)
  service and not providing HTTP-specific options
- Add `IncludedTraffic`, `OutgoingTraffic` and `IngoingTraffic` fields to `LoadBalancer`
- Add `ChangeType()` method to the Load Balancer client
- Fix retrying of requests that contain a body

## v1.18.2

- Retry API requests on conflict error

## v1.18.1

- Make all `GetByName` methods return `nil` when an empty name is provided
- Clarify that filters specified in options for List() calls are not taken
  into account when their value corresponds to their zero value or when
  they are empty.

## v1.18.0

- Add `Status` field to `Volume`
- Add subnet type `cloud`
- Add `WithHTTPClient` option to specify a custom `http.Client`
- Add API for requesting a VNC console
- Add support for load balancers and certificates (beta)

## v1.17.0

- Add `Created` field to `SSHKey`

## v1.16.0

- Make IP range optional when adding a subnet to a network
- Add support for names to Floating IPs

## v1.15.1

- Rename `MacAddress` to `MACAddress` on `ServerPrivateNet`

## v1.15.0

- Add `MacAddress` field to `ServerPrivateNet`
- Add `WithDebugWriter()` client option to provide an `io.Writer` to write debug output to

## v1.14.0

- Add `Created` field to `FloatingIP`
- Add support for networks

## v1.13.0

- Add missing fields to `*ListOpts` structs
- Fix error handling in `WatchProgress()`
- Add support for filtering volumes, images, and servers by status

## v1.12.0

- Add missing constants for all [documented error codes](https://docs.hetzner.cloud/reference/cloud#errors)
- Add support for automounting volumes
- Add support for attaching volumes when creating a server

## v1.11.0

- Add `NextActions` to `ServerCreateResult` and `VolumeCreateResult`

## v1.10.0

- Add `WithApplication()` client option to provide an application name and version
  that will be included in the `User-Agent` HTTP header
- Add support for volumes

## v1.9.0

- Add `AllWithOpts()` to server, Floating IP, image, and SSH key client
- Expose labels of servers, Floating IPs, images, and SSH Keys

## v1.8.0

- Add `WithPollInterval()` option to `Client` which allows to specify the polling interval
  ([issue #92](https://github.com/hetznercloud/hcloud-go/issues/92))
- Add `CPUType` field to `ServerType` ([issue #91](https://github.com/hetznercloud/hcloud-go/pull/91))

## v1.7.0

- Add `Deprecated ` field to `Image` ([issue #88](https://github.com/hetznercloud/hcloud-go/issues/88))
- Add `StartAfterCreate` flag to `ServerCreateOpts` ([issue #87](https://github.com/hetznercloud/hcloud-go/issues/87))
- Fix enum types ([issue #89](https://github.com/hetznercloud/hcloud-go/issues/89))

## v1.6.0

- Add `ChangeProtection()` to server, Floating IP, and image client
- Expose protection of servers, Floating IPs, and images

## v1.5.0

- Add `GetByFingerprint()` to SSH key client

## v1.4.0

- Retry all calls that triggered the API ratelimit
- Slow down `WatchProgress()` in action client from 100ms polling interval to 500ms

## v1.3.1

- Make clients using the old error code for ratelimiting work as expected
  ([issue #73](https://github.com/hetznercloud/hcloud-go/issues/73))

## v1.3.0

- Support passing user data on server creation ([issue #70](https://github.com/hetznercloud/hcloud-go/issues/70))
- Fix leaking response body by not closing it ([issue #68](https://github.com/hetznercloud/hcloud-go/issues/68))

## v1.2.0

- Add `WatchProgress()` to action client
- Use correct error code for ratelimit error (deprecated
  `ErrorCodeLimitReached`, added `ErrorCodeRateLimitExceeded`)

## v1.1.0

- Add `Image` field to `Server`
