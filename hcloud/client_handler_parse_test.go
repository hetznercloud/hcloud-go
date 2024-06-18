package hcloud

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseHandler(t *testing.T) {
	type SomeStruct struct {
		Data string `json:"data"`
	}

	testCases := []struct {
		name    string
		wrapped func(req *http.Request, v any) (*Response, error)
		want    func(t *testing.T, v SomeStruct, resp *Response, err error)
	}{
		{
			name: "no error",
			wrapped: func(_ *http.Request, _ any) (*Response, error) {
				return fakeResponse(t, 200, `{"data": "Hello", "meta": {"pagination": {"page": 1}}}`, true), nil
			},
			want: func(t *testing.T, v SomeStruct, resp *Response, err error) {
				assert.NoError(t, err)
				assert.Equal(t, v.Data, "Hello")
				assert.Equal(t, resp.Meta.Pagination.Page, 1)
			},
		},
		{
			name: "any error",
			wrapped: func(_ *http.Request, _ any) (*Response, error) {
				return nil, fmt.Errorf("any error")
			},
			want: func(t *testing.T, v SomeStruct, resp *Response, err error) {
				assert.EqualError(t, err, "any error")
				assert.Equal(t, v.Data, "")
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			m := &mockHandler{testCase.wrapped}
			h := wrapParseHandler(m)

			s := SomeStruct{}

			resp, err := h.Do(nil, &s)

			testCase.want(t, s, resp, err)
		})
	}
}
