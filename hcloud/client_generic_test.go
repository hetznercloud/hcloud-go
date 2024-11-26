package hcloud

import (
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/hetznercloud/hcloud-go/v2/hcloud/exp/mockutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/schema"
)

func TestGenericRequest(t *testing.T) {
	ctx := context.Background()

	t.Run("get", func(t *testing.T) {
		server := mockutil.NewServer(t, []mockutil.Request{
			{
				Method: "GET", Path: "/resource",
				Status: 200,
				JSON:   schema.ActionGetResponse{Action: schema.Action{ID: 1234}},
			},
		})
		client := NewClient(WithEndpoint(server.URL))

		respBody, resp, err := getRequest[schema.ActionGetResponse](ctx, client, "/resource")
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NotNil(t, respBody)

		require.Equal(t, int64(1234), respBody.Action.ID)
	})

	t.Run("post", func(t *testing.T) {
		server := mockutil.NewServer(t, []mockutil.Request{
			{
				Method: "POST", Path: "/resource",
				Want: func(t *testing.T, r *http.Request) {
					bodyBytes, err := io.ReadAll(r.Body)
					require.NoError(t, err)
					require.JSONEq(t, `{"hello": "world"}`, string(bodyBytes))
				},
				Status: 200,
				JSON:   schema.ActionGetResponse{Action: schema.Action{ID: 1234}},
			},
		})
		client := NewClient(WithEndpoint(server.URL))

		respBody, resp, err := postRequest[schema.ActionGetResponse](ctx, client, "/resource", map[string]string{"hello": "world"})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NotNil(t, respBody)

		require.Equal(t, int64(1234), respBody.Action.ID)
	})

	t.Run("put", func(t *testing.T) {
		server := mockutil.NewServer(t, []mockutil.Request{
			{
				Method: "PUT", Path: "/resource",
				Want: func(t *testing.T, r *http.Request) {
					bodyBytes, err := io.ReadAll(r.Body)
					require.NoError(t, err)
					require.JSONEq(t, `{"hello": "world"}`, string(bodyBytes))
				},
				Status: 200,
				JSON:   schema.ActionGetResponse{Action: schema.Action{ID: 1234}},
			},
		})
		client := NewClient(WithEndpoint(server.URL))

		respBody, resp, err := putRequest[schema.ActionGetResponse](ctx, client, "/resource", map[string]string{"hello": "world"})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NotNil(t, respBody)

		require.Equal(t, int64(1234), respBody.Action.ID)
	})

	t.Run("delete", func(t *testing.T) {
		server := mockutil.NewServer(t, []mockutil.Request{
			{
				Method: "DELETE", Path: "/resource",
				Status: 204,
			},
		})
		client := NewClient(WithEndpoint(server.URL))

		resp, err := deleteRequest(ctx, client, "/resource")
		require.NoError(t, err)
		require.NotNil(t, resp)
	})
}
