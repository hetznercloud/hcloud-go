# Changes

## master

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
