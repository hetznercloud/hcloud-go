# Changelog

## [1.39.0](https://github.com/hetznercloud/hcloud-go/compare/v1.38.0...v1.39.0) (2022-12-29)


### Features

* Use generics to get pointers to types ([#219](https://github.com/hetznercloud/hcloud-go/issues/219)) ([a5cd797](https://github.com/hetznercloud/hcloud-go/commit/a5cd79782dc849b3137e46ada2da6b319d4093c8))


### Bug Fixes

* deprecate PricingPrimaryIPTypePrice.Datacenter for Location ([#222](https://github.com/hetznercloud/hcloud-go/issues/222)) ([e0e5a1e](https://github.com/hetznercloud/hcloud-go/commit/e0e5a1e08fd7c0864fd94a787ee86714b5e9afc5))

## v1.38.0

### What's Changed
* feat(network): add new Network Zone us-west by @apricote in https://github.com/hetznercloud/hcloud-go/pull/217
* chore: prepare v1.38.0 by @apricote in https://github.com/hetznercloud/hcloud-go/pull/218


**Full Changelog**: https://github.com/hetznercloud/hcloud-go/compare/v1.37.0...v1.38.0

## v1.37.0

### What's Changed
* PrimaryIPClient Add AllWithOpts by @LKaemmerling in https://github.com/hetznercloud/hcloud-go/pull/214
* fix: error when updating IPv6 Primary IP by @apricote in https://github.com/hetznercloud/hcloud-go/pull/215


**Full Changelog**: https://github.com/hetznercloud/hcloud-go/compare/v1.36.0...v1.37.0

## v1.36.0

### What's Changed
* feat: add ServerClient.DeleteWithResult method by @apricote in https://github.com/hetznercloud/hcloud-go/pull/213

### New Contributors
* @apricote made their first contribution in https://github.com/hetznercloud/hcloud-go/pull/213

**Full Changelog**: https://github.com/hetznercloud/hcloud-go/compare/v1.35.3...v1.36.0

## v1.35.3

### What's Changed
* Drop support for Go < 1.17 and add official tests on go 1.19 by @LKaemmerling in https://github.com/hetznercloud/hcloud-go/pull/211
* Stop automatic retrying on RateLimitExceeded by @LKaemmerling in https://github.com/hetznercloud/hcloud-go/pull/210


**Full Changelog**: https://github.com/hetznercloud/hcloud-go/compare/v1.35.2...v1.35.3

## v1.35.2

### What's Changed
* Allow empty labels by @4ND3R50N in https://github.com/hetznercloud/hcloud-go/pull/207


**Full Changelog**: https://github.com/hetznercloud/hcloud-go/compare/v1.35.1...v1.35.2

## v1.35.1

### What's Changed
* Accept no primary IPs with server create with StartAfterCreate = false by @4ND3R50N in https://github.com/hetznercloud/hcloud-go/pull/205


**Full Changelog**: https://github.com/hetznercloud/hcloud-go/compare/v1.35.0...v1.35.1

## v1.35.0

### What's Changed
* Catch invalid token values and error out without value exposure by @NotTheEvilOne in https://github.com/hetznercloud/hcloud-go/pull/194
* Remove ServerRescueTypeFreeBSD64 by @LKaemmerling in https://github.com/hetznercloud/hcloud-go/pull/203
* Add Primary IP Support by @LKaemmerling in https://github.com/hetznercloud/hcloud-go/pull/204

### New Contributors
* @NotTheEvilOne made their first contribution in https://github.com/hetznercloud/hcloud-go/pull/194

**Full Changelog**: https://github.com/hetznercloud/hcloud-go/compare/v1.34.0...v1.35.0

## v1.34.0

### What's Changed
* Test on Go 1.18 by @LKaemmerling in https://github.com/hetznercloud/hcloud-go/pull/202
* Add support for sorting the response of all list calls by @LKaemmerling in https://github.com/hetznercloud/hcloud-go/pull/201
* Set UsePrivateIP for targets when creating a LoadBalancer by @hakman in https://github.com/hetznercloud/hcloud-go/pull/198

### New Contributors
* @hakman made their first contribution in https://github.com/hetznercloud/hcloud-go/pull/198

**Full Changelog**: https://github.com/hetznercloud/hcloud-go/compare/v1.33.2...v1.34.0

## v1.33.2

### What's Changed
* Add constant for resource locked error code by @patrickschaffrath in https://github.com/hetznercloud/hcloud-go/pull/189
* Fix metadata client error detection by @choffmeister in https://github.com/hetznercloud/hcloud-go/pull/193
* Add labels.go to validate resource labels by @4ND3R50N in https://github.com/hetznercloud/hcloud-go/pull/197

### New Contributors
* @patrickschaffrath made their first contribution in https://github.com/hetznercloud/hcloud-go/pull/189
* @choffmeister made their first contribution in https://github.com/hetznercloud/hcloud-go/pull/193
* @4ND3R50N made their first contribution in https://github.com/hetznercloud/hcloud-go/pull/197

**Full Changelog**: https://github.com/hetznercloud/hcloud-go/compare/v1.33.1...v1.33.2

## v1.33.1

### Changelog

41fef2f Add constants for new firewall error code

## v1.33.0

### What's Changed
* Add us-east network zone by @LKaemmerling in https://github.com/hetznercloud/hcloud-go/pull/187


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

* Support Hetzner Cloud managed Certificates (#167)

## v1.23.1
* Add removed `ErrorCodeServerAlreadyAttached` again

## v1.23.0

* Add missing constants for all resource specific error codes
* Expose metrics for Servers and Load Balancers
* Add support for vSwitch Subnetworks

## v1.22.0

* Add `PrimaryDiskSize` Field to `Server`

## v1.21.1

* Don't send `Authorization` Header when `WithToken` was not called

## v1.21.0

* Add `IncludeDeprecated` Field to `ImageListOpts`

## v1.20.0

* Add support for Load Balancer Label Selector targets
* Add support for Load Balancer IP targets

## v1.19.0

* Fix nil pointer dereference when creating a Load Balancer with HTTP(S)
  service and not providing HTTP-specific options
* Add `IncludedTraffic`, `OutgoingTraffic` and `IngoingTraffic` fields to `LoadBalancer`
* Add `ChangeType()` method to the Load Balancer client
* Fix retrying of requests that contain a body

## v1.18.2

* Retry API requests on conflict error

## v1.18.1

* Make all `GetByName` methods return `nil` when an empty name is provided
* Clarify that filters specified in options for List() calls are not taken
  into account when their value corresponds to their zero value or when
  they are empty.

## v1.18.0

* Add `Status` field to `Volume`
* Add subnet type `cloud`
* Add `WithHTTPClient` option to specify a custom `http.Client`
* Add API for requesting a VNC console
* Add support for load balancers and certificates (beta)

## v1.17.0

* Add `Created` field to `SSHKey`

## v1.16.0

* Make IP range optional when adding a subnet to a network
* Add support for names to Floating IPs

## v1.15.1

* Rename `MacAddress` to `MACAddress` on `ServerPrivateNet`

## v1.15.0

* Add `MacAddress` field to `ServerPrivateNet`
* Add `WithDebugWriter()` client option to provide an `io.Writer` to write debug output to

## v1.14.0

* Add `Created` field to `FloatingIP`
* Add support for networks

## v1.13.0

* Add missing fields to `*ListOpts` structs
* Fix error handling in `WatchProgress()`
* Add support for filtering volumes, images, and servers by status

## v1.12.0

* Add missing constants for all [documented error codes](https://docs.hetzner.cloud/#overview-errors)
* Add support for automounting volumes
* Add support for attaching volumes when creating a server

## v1.11.0

* Add `NextActions` to `ServerCreateResult` and `VolumeCreateResult`

## v1.10.0

* Add `WithApplication()` client option to provide an application name and version
  that will be included in the `User-Agent` HTTP header
* Add support for volumes

## v1.9.0

* Add `AllWithOpts()` to server, Floating IP, image, and SSH key client
* Expose labels of servers, Floating IPs, images, and SSH Keys

## v1.8.0

* Add `WithPollInterval()` option to `Client` which allows to specify the polling interval
  ([issue #92](https://github.com/hetznercloud/hcloud-go/issues/92))
* Add `CPUType` field to `ServerType` ([issue #91](https://github.com/hetznercloud/hcloud-go/pull/91))

## v1.7.0

* Add `Deprecated ` field to `Image` ([issue #88](https://github.com/hetznercloud/hcloud-go/issues/88))
* Add `StartAfterCreate` flag to `ServerCreateOpts` ([issue #87](https://github.com/hetznercloud/hcloud-go/issues/87))
* Fix enum types ([issue #89](https://github.com/hetznercloud/hcloud-go/issues/89))

## v1.6.0

* Add `ChangeProtection()` to server, Floating IP, and image client
* Expose protection of servers, Floating IPs, and images

## v1.5.0

* Add `GetByFingerprint()` to SSH key client

## v1.4.0

* Retry all calls that triggered the API ratelimit
* Slow down `WatchProgress()` in action client from 100ms polling interval to 500ms

## v1.3.1

* Make clients using the old error code for ratelimiting work as expected
  ([issue #73](https://github.com/hetznercloud/hcloud-go/issues/73))

## v1.3.0

* Support passing user data on server creation ([issue #70](https://github.com/hetznercloud/hcloud-go/issues/70))
* Fix leaking response body by not closing it ([issue #68](https://github.com/hetznercloud/hcloud-go/issues/68))

## v1.2.0

* Add `WatchProgress()` to action client
* Use correct error code for ratelimit error (deprecated
  `ErrorCodeLimitReached`, added `ErrorCodeRateLimitExceeded`)

## v1.1.0

* Add `Image` field to `Server`
