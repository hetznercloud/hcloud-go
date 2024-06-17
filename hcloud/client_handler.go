package hcloud

import (
	"net/http"
)

// handler is an interface representing a client request transaction. The handler are
// meant to be chained, similarly to the [http.RoundTripper] interface.
//
// The handler chain is placed between the [Client] API operations and the
// [http.Client].
type handler interface {
	Do(req *http.Request, v any) (resp *Response, err error)
}

// assembleHandlerChain assembles the chain of handlers used to make API requests.
//
// The order of the handlers is important.
func assembleHandlerChain(client *Client) handler {
	// Start down the chain: sending the http request
	h := newHTTPHandler(client.httpClient)

	// Read rate limit headers
	h = wrapRateLimitHandler(h)

	// Build error from response
	h = wrapErrorHandler(h)

	// Insert debug writer if enabled
	if client.debugWriter != nil {
		h = wrapDebugHandler(h, client.debugWriter)
	}

	// Retry request if condition are met
	h = wrapRetryHandler(h, client.backoffFunc)

	// Finally parse the response body into the provided schema
	h = wrapParseHandler(h)

	return h
}
