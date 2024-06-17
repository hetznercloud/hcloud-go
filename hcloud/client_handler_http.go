package hcloud

import (
	"bytes"
	"io"
	"net/http"
)

func newHTTPHandler(httpClient *http.Client) handler {
	return &httpHandler{httpClient}
}

type httpHandler struct {
	httpClient *http.Client
}

func (h *httpHandler) Do(req *http.Request, _ interface{}) (*Response, error) {
	httpResponse, err := h.httpClient.Do(req) //nolint: bodyclose
	resp := &Response{Response: httpResponse}
	if err != nil {
		return resp, err
	}

	// Read full response body and save it for later use
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return resp, err
	}
	resp.body = body

	// Restore the body as if it was untouched, as it might be read by external users
	resp.Body = io.NopCloser(bytes.NewReader(body))

	return resp, err
}
