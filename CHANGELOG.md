# Changelog

## [2.9.0](https://github.com/hetznercloud/hcloud-go/compare/v2.8.0...v2.9.0) (2024-05-29)


### Features

* **exp:** add `AppendNextActions` function ([#440](https://github.com/hetznercloud/hcloud-go/issues/440)) ([b07d7ad](https://github.com/hetznercloud/hcloud-go/commit/b07d7adf0bb08dd372f1f6ff630d16af8cb1265c))
* **exp:** add ssh key functions ([#441](https://github.com/hetznercloud/hcloud-go/issues/441)) ([d766e96](https://github.com/hetznercloud/hcloud-go/commit/d766e96e2220b069a0a8067a3249afe48cf68b2e))


### Bug Fixes

* **exp:** rename to `sshutils` package name ([#450](https://github.com/hetznercloud/hcloud-go/issues/450)) ([6d4100d](https://github.com/hetznercloud/hcloud-go/commit/6d4100dbceb9a43a66bad3ce91cee039198130c3))

## [2.8.0](https://github.com/hetznercloud/hcloud-go/compare/v2.7.2...v2.8.0) (2024-05-06)


### Features

* **error:** show internal correlation id in error messages ([#411](https://github.com/hetznercloud/hcloud-go/issues/411)) ([6c96d19](https://github.com/hetznercloud/hcloud-go/commit/6c96d19dacde736b52abd7b8fc8879c8f721f23b))
* implement actions waiter ([#407](https://github.com/hetznercloud/hcloud-go/issues/407)) ([1e3fa70](https://github.com/hetznercloud/hcloud-go/commit/1e3fa7033d8a1cea1c0a9007a6255798979f0771))
* require Go &gt;= 1.21 ([#424](https://github.com/hetznercloud/hcloud-go/issues/424)) ([d4f4000](https://github.com/hetznercloud/hcloud-go/commit/d4f40009016c3ed5bb14ef9ae16bcf6aefa90fee))


### Bug Fixes

* improve error message format with correlation id ([#430](https://github.com/hetznercloud/hcloud-go/issues/430)) ([013477f](https://github.com/hetznercloud/hcloud-go/commit/013477f4227bdc69f376d8f13a875b09c32171f6))

## [2.7.2](https://github.com/hetznercloud/hcloud-go/compare/v2.7.1...v2.7.2) (2024-04-19)


### Bug Fixes

* creating firewall with rules not working correctly ([#412](https://github.com/hetznercloud/hcloud-go/issues/412)) ([43b2bab](https://github.com/hetznercloud/hcloud-go/commit/43b2bab8c82cab276c07e9deb9f4e422e10ccc82))

## [2.7.1](https://github.com/hetznercloud/hcloud-go/compare/v2.7.0...v2.7.1) (2024-04-18)


### Bug Fixes

* setting firewall rules not working correctly ([#405](https://github.com/hetznercloud/hcloud-go/issues/405)) ([16daea0](https://github.com/hetznercloud/hcloud-go/commit/16daea0dfb32a5e2b8a71fce88d1c897cd4134eb))

## [2.7.0](https://github.com/hetznercloud/hcloud-go/compare/v2.6.0...v2.7.0) (2024-03-27)


### Features

* add volume format property ([#397](https://github.com/hetznercloud/hcloud-go/issues/397)) ([c0940af](https://github.com/hetznercloud/hcloud-go/commit/c0940afce9eb01c0e6838502c91aa569ab411a03))
* **error:** handle wrapped errors in IsError() ([#374](https://github.com/hetznercloud/hcloud-go/issues/374)) ([83df108](https://github.com/hetznercloud/hcloud-go/commit/83df108228519617a919470e5ebbb4a1a2587b34))
* require Go &gt;= 1.20 ([#392](https://github.com/hetznercloud/hcloud-go/issues/392)) ([299f181](https://github.com/hetznercloud/hcloud-go/commit/299f181c469a48e977743f2587e232a293cd9db5))
* **server:** deprecate ServerRescueTypeLinux32 ([#378](https://github.com/hetznercloud/hcloud-go/issues/378)) ([2f334c3](https://github.com/hetznercloud/hcloud-go/commit/2f334c3de2aecfb7aaf179424bb9a1a7b0da53ea))
* test with Go 1.22 ([#391](https://github.com/hetznercloud/hcloud-go/issues/391)) ([49be506](https://github.com/hetznercloud/hcloud-go/commit/49be50664f56e4e315e687543ee0bf7fbb2af186))


### Bug Fixes

* falsely omitted fields in firewall schema ([#396](https://github.com/hetznercloud/hcloud-go/issues/396)) ([a3509b6](https://github.com/hetznercloud/hcloud-go/commit/a3509b6b0f2762b5d1e27f2374d3fb5fb53206e2))
* LoadBalancerUpdateServiceOpts not converted correctly ([#394](https://github.com/hetznercloud/hcloud-go/issues/394)) ([0f187ce](https://github.com/hetznercloud/hcloud-go/commit/0f187cef1f568d87a32b2ec063847dd9ff314740))
* primary ip assignee id not mapped to nil ([#395](https://github.com/hetznercloud/hcloud-go/issues/395)) ([b5fea38](https://github.com/hetznercloud/hcloud-go/commit/b5fea38f5d9d62c88af470ba69195a73ee075f09))

## [2.6.0](https://github.com/hetznercloud/hcloud-go/compare/v2.5.1...v2.6.0) (2024-01-09)


### Features

* alias deprecated field to deprecation info struct ([#371](https://github.com/hetznercloud/hcloud-go/issues/371)) ([e961be9](https://github.com/hetznercloud/hcloud-go/commit/e961be9615452fc63c4c71e6561d4e86f8e8e95a))
* **instrumentation:** allow passing in any prometheus.Registerer ([#369](https://github.com/hetznercloud/hcloud-go/issues/369)) ([0821c07](https://github.com/hetznercloud/hcloud-go/commit/0821c078900910fa9e3ca6c6c0af48a73f00c7c6))

## [2.5.1](https://github.com/hetznercloud/hcloud-go/compare/v2.5.0...v2.5.1) (2023-12-13)


### Bug Fixes

* schema conversion outputs debug messages to stdout ([#354](https://github.com/hetznercloud/hcloud-go/issues/354)) ([ade8fbd](https://github.com/hetznercloud/hcloud-go/commit/ade8fbd60a88a648c95391f00cfe3ccc09be8f37))

## [2.5.0](https://github.com/hetznercloud/hcloud-go/compare/v2.4.0...v2.5.0) (2023-12-12)


### Features

* add conversion methods from schema to hcloud objects ([#343](https://github.com/hetznercloud/hcloud-go/issues/343)) ([6feda4d](https://github.com/hetznercloud/hcloud-go/commit/6feda4d9b0e7cf3f5a17a4b38504abbe5213883d))
* add interfaces for client structs ([#342](https://github.com/hetznercloud/hcloud-go/issues/342)) ([4f9390f](https://github.com/hetznercloud/hcloud-go/commit/4f9390f8387d1c86330156adbb6801aacba7a8f0))
* add missing properties ([#349](https://github.com/hetznercloud/hcloud-go/issues/349)) ([c8a28d0](https://github.com/hetznercloud/hcloud-go/commit/c8a28d0dbf0c84364401282a60b44ccea1da6423))
* **error:** include http response in api errors ([#320](https://github.com/hetznercloud/hcloud-go/issues/320)) ([9558239](https://github.com/hetznercloud/hcloud-go/commit/95582395dfb1039f4ce4f10a1ac9c068db93a867))


### Bug Fixes

* make schemas consistent with API ([#348](https://github.com/hetznercloud/hcloud-go/issues/348)) ([b0d7055](https://github.com/hetznercloud/hcloud-go/commit/b0d7055543669fb96af1726daa7c1458fb1b65a2))

## [2.4.0](https://github.com/hetznercloud/hcloud-go/compare/v2.3.0...v2.4.0) (2023-10-12)


### Features

* **iso:** support deprecation info API ([#331](https://github.com/hetznercloud/hcloud-go/issues/331)) ([b3a3621](https://github.com/hetznercloud/hcloud-go/commit/b3a36214b21ab1c5add94c2b1896a995342757b6))

## [2.3.0](https://github.com/hetznercloud/hcloud-go/compare/v2.2.0...v2.3.0) (2023-10-04)


### Features

* add error code unauthorized  ([#313](https://github.com/hetznercloud/hcloud-go/issues/313)) ([b77d9e0](https://github.com/hetznercloud/hcloud-go/commit/b77d9e04ca903448cc1a22c242f440e67a81a028))
* test with Go 1.21 ([#319](https://github.com/hetznercloud/hcloud-go/issues/319)) ([7ddb2ec](https://github.com/hetznercloud/hcloud-go/commit/7ddb2ec057d0e165abb5b8cbd74e95d5bb4add49))


### Bug Fixes

* ensure the pollBackoffFunc is correctly set ([#322](https://github.com/hetznercloud/hcloud-go/issues/322)) ([2b2f869](https://github.com/hetznercloud/hcloud-go/commit/2b2f8697aa09f67dd508c33fe36f7a59b9d3f192))

## [2.2.0](https://github.com/hetznercloud/hcloud-go/compare/v2.1.1...v2.2.0) (2023-08-24)


### Features

* support resource-specific Action endpoints ([#295](https://github.com/hetznercloud/hcloud-go/issues/295)) ([ddc2ac4](https://github.com/hetznercloud/hcloud-go/commit/ddc2ac45489c48a7563806c425222236ab1f8aa0))

## [2.1.1](https://github.com/hetznercloud/hcloud-go/compare/v2.1.0...v2.1.1) (2023-08-15)


### Bug Fixes

* send more precise progress values ([#304](https://github.com/hetznercloud/hcloud-go/issues/304)) ([867aa63](https://github.com/hetznercloud/hcloud-go/commit/867aa632521ad3acfb04beb52b6307330740fc68))

## [2.1.0](https://github.com/hetznercloud/hcloud-go/compare/v2.0.0...v2.1.0) (2023-08-08)


### Features

* **metadata:** add timeout option ([#293](https://github.com/hetznercloud/hcloud-go/issues/293)) ([913bf74](https://github.com/hetznercloud/hcloud-go/commit/913bf74071e03a9c79fd0f8a5a37b1e11b350ae1))


### Bug Fixes

* **action:** unexpected behaviour when watching non-existing Actions ([#298](https://github.com/hetznercloud/hcloud-go/issues/298)) ([0727d42](https://github.com/hetznercloud/hcloud-go/commit/0727d42e26a8112923a84b84bb506c980e07262d))
* **instrumentation:** multiple instrumented clients cause panic ([#289](https://github.com/hetznercloud/hcloud-go/issues/289)) ([c0ef9b6](https://github.com/hetznercloud/hcloud-go/commit/c0ef9b6e6e3f36d8c2282c2b7aa9d8687141f291))

## [2.0.0](https://github.com/hetznercloud/hcloud-go/compare/v1.47.0...v2.0.0) (2023-07-12)


### ⚠ BREAKING CHANGES

* use int64 for ID fields ([#282](https://github.com/hetznercloud/hcloud-go/issues/282))

### Features

* make ListOpts.Values method public ([#285](https://github.com/hetznercloud/hcloud-go/issues/285)) ([c82ea59](https://github.com/hetznercloud/hcloud-go/commit/c82ea5971d94b812fd9302bfc4151f4ebfa43413))
* use int64 for ID fields ([#282](https://github.com/hetznercloud/hcloud-go/issues/282)) ([359c389](https://github.com/hetznercloud/hcloud-go/commit/359c3894641f2dcca4a049537e256a20853b5ad9))


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
