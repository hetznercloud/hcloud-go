package retryutils

import (
	"time"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// TODO: generate the opts from the [hcloud.Client].
type Opts struct {
	Backoff    hcloud.BackoffFunc
	MaxRetries int
	Policy     func(err error) bool
}

func Retry[T any](opts Opts, request func() (T, *hcloud.Response, error)) (T, *hcloud.Response, error) {
	retries := 0
	for {
		result, resp, err := request()
		if err != nil {
			if opts.Policy(err) && retries < opts.MaxRetries {
				select {
				case <-resp.Request.Context().Done():
					break
				case <-time.After(opts.Backoff(retries)):
					retries++
					continue
				}
			}
		}
		return result, resp, err
	}
}

func RetryNoResult(opts Opts, request func() (*hcloud.Response, error)) (*hcloud.Response, error) {
	_, resp, err := Retry(opts, func() (any, *hcloud.Response, error) {
		resp, err := request()
		return nil, resp, err
	})

	return resp, err
}
