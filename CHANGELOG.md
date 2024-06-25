# Changelog

## [1.57.0](https://github.com/hetznercloud/hcloud-go/compare/v1.56.0...v1.57.0) (2024-06-25)


### Features

* filter out nil actions in action waiter ([#464](https://github.com/hetznercloud/hcloud-go/issues/464)) ([#471](https://github.com/hetznercloud/hcloud-go/issues/471)) ([14bf589](https://github.com/hetznercloud/hcloud-go/commit/14bf589531514c49de5ea374696a880dea09a355))

## [1.56.0](https://github.com/hetznercloud/hcloud-go/compare/v1.55.0...v1.56.0) (2024-05-22)


### Features

* add volume format property ([5ff15f8](https://github.com/hetznercloud/hcloud-go/commit/5ff15f8461b1ffb69418209b332f5b5905a024c7))

## [1.55.0](https://github.com/hetznercloud/hcloud-go/compare/v1.54.1...v1.55.0) (2024-05-06)


### Features

* **error:** show internal correlation id in error messages ([#417](https://github.com/hetznercloud/hcloud-go/issues/417)) ([f7f96a5](https://github.com/hetznercloud/hcloud-go/commit/f7f96a5d5ebc0d8239b946b928d5bd534cc2bf21))
* implement actions waiter ([c4d8be8](https://github.com/hetznercloud/hcloud-go/commit/c4d8be8f053469948d564e93b043af05c61afb41))
* require Go &gt;= 1.21 ([18adaf6](https://github.com/hetznercloud/hcloud-go/commit/18adaf6abe5caeee38af030af832af98baeaaf0e))


### Bug Fixes

* improve error message format with correlation id ([1939cfb](https://github.com/hetznercloud/hcloud-go/commit/1939cfb6c8bbcd495e79c059e11a30ffd7c16948))

## [1.54.1](https://github.com/hetznercloud/hcloud-go/compare/v1.54.0...v1.54.1) (2024-04-19)


### Bug Fixes

* creating firewall with rules not working correctly ([#414](https://github.com/hetznercloud/hcloud-go/issues/414)) ([2a19325](https://github.com/hetznercloud/hcloud-go/commit/2a193259d091f289707758aac87534e92bcbe272))

## [1.54.0](https://github.com/hetznercloud/hcloud-go/compare/v1.53.0...v1.54.0) (2024-04-18)


### Features

* **error:** handle wrapped errors in IsError() ([64902c5](https://github.com/hetznercloud/hcloud-go/commit/64902c582618a4042382d36a4005720528bfaf51))
* **server:** deprecate ServerRescueTypeLinux32 [#379](https://github.com/hetznercloud/hcloud-go/issues/379) ([7337010](https://github.com/hetznercloud/hcloud-go/commit/73370105e02bcfca0b97a2db8d30cf244866fdaf))


### Bug Fixes

* falsely omitted fields in firewall schema ([7f39bfa](https://github.com/hetznercloud/hcloud-go/commit/7f39bfa7da9eaf4ea7e5d7baae9354d1d4d1b3d6))

## [1.53.0](https://github.com/hetznercloud/hcloud-go/compare/v1.52.0...v1.53.0) (2023-10-20)


### Features

* **error:** include http response in api errors ([#335](https://github.com/hetznercloud/hcloud-go/issues/335)) ([4984025](https://github.com/hetznercloud/hcloud-go/commit/49840253cdf566ba47a86917e5f8255438b0fc05))

## [1.52.0](https://github.com/hetznercloud/hcloud-go/compare/v1.51.0...v1.52.0) (2023-10-12)


### Features

* **iso:** support deprecation info API ([#333](https://github.com/hetznercloud/hcloud-go/issues/333)) ([5f89c57](https://github.com/hetznercloud/hcloud-go/commit/5f89c57c2e88d1521aafedb5c5e9b9c06160c3c0))

## [1.51.0](https://github.com/hetznercloud/hcloud-go/compare/v1.50.0...v1.51.0) (2023-10-04)


### Features

* add error code unauthorized  ([#313](https://github.com/hetznercloud/hcloud-go/issues/313)) ([#315](https://github.com/hetznercloud/hcloud-go/issues/315)) ([5542b11](https://github.com/hetznercloud/hcloud-go/commit/5542b11a3fd9923e8a8d6359dddf0b2cbc15a957))


### Bug Fixes

* ensure the pollBackoffFunc is correctly set ([#322](https://github.com/hetznercloud/hcloud-go/issues/322)) ([#323](https://github.com/hetznercloud/hcloud-go/issues/323)) ([9e19dea](https://github.com/hetznercloud/hcloud-go/commit/9e19deac2e65a80d7d76c30c1ce85c44f5d7d628))

## [1.50.0](https://github.com/hetznercloud/hcloud-go/compare/v1.49.1...v1.50.0) (2023-08-24)


### Features

* support resource-specific Action endpoints ([#295](https://github.com/hetznercloud/hcloud-go/issues/295)) ([#309](https://github.com/hetznercloud/hcloud-go/issues/309)) ([b492d68](https://github.com/hetznercloud/hcloud-go/commit/b492d68b5345afd875ae11b31a29125df3775d26))

## [1.49.1](https://github.com/hetznercloud/hcloud-go/compare/v1.49.0...v1.49.1) (2023-08-17)


### Bug Fixes

* send more precise progress values ([#304](https://github.com/hetznercloud/hcloud-go/issues/304)) ([#306](https://github.com/hetznercloud/hcloud-go/issues/306)) ([298f1f9](https://github.com/hetznercloud/hcloud-go/commit/298f1f91a2775ba704c733b4a05b63cfaf330019))

## [1.49.0](https://github.com/hetznercloud/hcloud-go/compare/v1.48.0...v1.49.0) (2023-08-08)


### Features

* **metadata:** add timeout option ([#293](https://github.com/hetznercloud/hcloud-go/issues/293)) ([#294](https://github.com/hetznercloud/hcloud-go/issues/294)) ([6b1d465](https://github.com/hetznercloud/hcloud-go/commit/6b1d4650a743e6d72a8498d262e1c270e38ab590))


### Bug Fixes

* **action:** unexpected behaviour when watching non-existing Actions ([#298](https://github.com/hetznercloud/hcloud-go/issues/298)) ([#300](https://github.com/hetznercloud/hcloud-go/issues/300)) ([c437163](https://github.com/hetznercloud/hcloud-go/commit/c437163010bd5e898f182fc15691830fd5a4e4a7))
* **instrumentation:** multiple instrumented clients cause panic ([#289](https://github.com/hetznercloud/hcloud-go/issues/289)) ([#290](https://github.com/hetznercloud/hcloud-go/issues/290)) ([11e2297](https://github.com/hetznercloud/hcloud-go/commit/11e22979044177a9fd1f6dac395581cf1b1b590c))

## [1.48.0](https://github.com/hetznercloud/hcloud-go/compare/v1.47.0...v1.48.0) (2023-07-12)


### Features

* make ListOpts.Values method public ([#285](https://github.com/hetznercloud/hcloud-go/issues/285)) ([c82ea59](https://github.com/hetznercloud/hcloud-go/commit/c82ea5971d94b812fd9302bfc4151f4ebfa43413))


### Bug Fixes

* **action:** show accurate progress in WatchOverallProgress ([#281](https://github.com/hetznercloud/hcloud-go/issues/281)) ([cae9e57](https://github.com/hetznercloud/hcloud-go/commit/cae9e5789b20bdb7d9213ba88897435ee1abce86))
* **iso:** invalid field include_wildcard_architecture ([188b68a](https://github.com/hetznercloud/hcloud-go/commit/188b68ad674066d302b2614432cf4f7d5b47f41a))

## [1.47.0](https://github.com/hetznercloud/hcloud-go/compare/v1.46.1...v1.47.0) (2023-06-21)


### Features

* **network:** add field expose_routes_to_vswitch ([#277](https://github.com/hetznercloud/hcloud-go/issues/277)) ([e73c52d](https://github.com/hetznercloud/hcloud-go/commit/e73c52dce563e00c8ccba7528fd054936a52b64c))

## [1.46.1](https://github.com/hetznercloud/hcloud-go/compare/v1.46.0...v1.46.1) (2023-06-16)


### Bug Fixes

* adjust label validation for max length of 63 characters ([#273](https://github.com/hetznercloud/hcloud-go/issues/273)) ([6382808](https://github.com/hetznercloud/hcloud-go/commit/63828086e89413115f4f5ee326835f93752cbb51))
* **deps:** update module golang.org/x/net to v0.11.0 ([#258](https://github.com/hetznercloud/hcloud-go/issues/258)) ([7918f21](https://github.com/hetznercloud/hcloud-go/commit/7918f2172ac000f72a6ed52a2c69ae03c7b6513f))

## [1.46.0](https://github.com/hetznercloud/hcloud-go/compare/v1.45.1...v1.46.0) (2023-06-13)


### Features

* provide `.AllWithOpts` method for all clients ([#266](https://github.com/hetznercloud/hcloud-go/issues/266)) ([2a7249e](https://github.com/hetznercloud/hcloud-go/commit/2a7249ed646bf9b5d91890fc6698b04bfdaf7806))
* **servertype:** implement new Deprecation api field ([#268](https://github.com/hetznercloud/hcloud-go/issues/268)) ([ac5ae2e](https://github.com/hetznercloud/hcloud-go/commit/ac5ae2e80c361775bd14c776e23ccec4ce5849e7))

## [1.45.1](https://github.com/hetznercloud/hcloud-go/compare/v1.45.0...v1.45.1) (2023-05-11)


### Bug Fixes

* **servertype:** use int64 to fit TB sizes on 32-bit platforms ([#261](https://github.com/hetznercloud/hcloud-go/issues/261)) ([2b19245](https://github.com/hetznercloud/hcloud-go/commit/2b1924575148a8675de7ca65f5578aeb70ef750f))

## [1.45.0](https://github.com/hetznercloud/hcloud-go/compare/v1.44.0...v1.45.0) (2023-05-11)


### Features

* **servertype:** add field for included traffic ([#259](https://github.com/hetznercloud/hcloud-go/issues/259)) ([d3b012a](https://github.com/hetznercloud/hcloud-go/commit/d3b012a678ee7012a54bb6088c3ee3b3efb5978d))

## [1.44.0](https://github.com/hetznercloud/hcloud-go/compare/v1.43.0...v1.44.0) (2023-05-05)


### Features

* **iso:** extend ISOClient by AllWithOpts method ([#254](https://github.com/hetznercloud/hcloud-go/issues/254)) ([c42f69b](https://github.com/hetznercloud/hcloud-go/commit/c42f69b05bf732b92561285874e4ac2b6cc01d1c))


### Bug Fixes

* **deps:** update module github.com/prometheus/client_golang to v1.15.1 ([#253](https://github.com/hetznercloud/hcloud-go/issues/253)) ([275f0fd](https://github.com/hetznercloud/hcloud-go/commit/275f0fd4d4316fed0acddd13400d8a3143889f49))

## [1.43.0](https://github.com/hetznercloud/hcloud-go/compare/v1.42.0...v1.43.0) (2023-04-26)


### Features

* **primary-ip:** implement RDNSSupporter ([#252](https://github.com/hetznercloud/hcloud-go/issues/252)) ([41a4c5a](https://github.com/hetznercloud/hcloud-go/commit/41a4c5a1d7f70fa6a279c6a2834830e806d456de))


### Bug Fixes

* **deps:** update module github.com/prometheus/client_golang to v1.15.0 ([#250](https://github.com/hetznercloud/hcloud-go/issues/250)) ([f10e804](https://github.com/hetznercloud/hcloud-go/commit/f10e8042ac12e6195824b80a791700ce857111bc))

## [1.42.0](https://github.com/hetznercloud/hcloud-go/compare/v1.41.0...v1.42.0) (2023-04-12)


### Features

* add support for ARM APIs ([#249](https://github.com/hetznercloud/hcloud-go/issues/249)) ([ce9859f](https://github.com/hetznercloud/hcloud-go/commit/ce9859f178078c99f3d15de41c7f0266c2e885e1))


### Bug Fixes

* **deps:** update module golang.org/x/net to v0.9.0 ([#247](https://github.com/hetznercloud/hcloud-go/issues/247)) ([962afeb](https://github.com/hetznercloud/hcloud-go/commit/962afebeed76560687103777efae20ff70fbcf16))

## [1.41.0](https://github.com/hetznercloud/hcloud-go/compare/v1.40.0...v1.41.0) (2023-03-06)


### Features

* add ServerClient.RebuildWithResult to return root password ([#245](https://github.com/hetznercloud/hcloud-go/issues/245)) ([82f97cf](https://github.com/hetznercloud/hcloud-go/commit/82f97cf48695848e2569b38f8ff24bb050966ee4))


### Bug Fixes

* **deps:** update module github.com/google/go-cmp to v0.5.9 ([#237](https://github.com/hetznercloud/hcloud-go/issues/237)) ([2237ff7](https://github.com/hetznercloud/hcloud-go/commit/2237ff795cbaf1e75759cdd396b3dfe5491c0e24))
* **deps:** update module github.com/prometheus/client_golang to v1.14.0 ([#241](https://github.com/hetznercloud/hcloud-go/issues/241)) ([75a4a01](https://github.com/hetznercloud/hcloud-go/commit/75a4a0140216eb476990e50ab9b13b60881404be))
* **deps:** update module github.com/stretchr/testify to v1.8.2 ([#242](https://github.com/hetznercloud/hcloud-go/issues/242)) ([4b51f1e](https://github.com/hetznercloud/hcloud-go/commit/4b51f1e8a13f1f859211910f1dce2daebb583b04))
* **deps:** update module golang.org/x/net to v0.7.0 [security] ([#236](https://github.com/hetznercloud/hcloud-go/issues/236)) ([774a560](https://github.com/hetznercloud/hcloud-go/commit/774a560b3d167c5c55cd3cbc4f83872ecc878670))
* **deps:** update module golang.org/x/net to v0.8.0 ([#243](https://github.com/hetznercloud/hcloud-go/issues/243)) ([8ae14f3](https://github.com/hetznercloud/hcloud-go/commit/8ae14f36021a32f5bab21a74d2467aa2487b348d))

## [1.40.0](https://github.com/hetznercloud/hcloud-go/compare/v1.39.0...v1.40.0) (2023-02-08)


### Features

* **action:** use configurable backoff to wait for action progress ([#227](https://github.com/hetznercloud/hcloud-go/issues/227)) ([8da6417](https://github.com/hetznercloud/hcloud-go/commit/8da6417cf7d87bf44117ede9cd2839d7dc055f66))
* support go v1.20 and drop v1.18 ([#231](https://github.com/hetznercloud/hcloud-go/issues/231)) ([44af6e5](https://github.com/hetznercloud/hcloud-go/commit/44af6e5beade11432b5ca30575781875cbd08343))

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
