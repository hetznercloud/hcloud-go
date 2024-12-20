package hcloud

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/hetznercloud/hcloud-go/v2/hcloud/exp/mockutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/schema"
)

func TestGenericRequest(t *testing.T) {
	t.Run("get", func(t *testing.T) {
		ctx, server, client := makeTestUtils(t)

		server.Expect([]mockutil.Request{
			{
				Method: "GET", Path: "/resource",
				Status: 200,
				JSON:   schema.ActionGetResponse{Action: schema.Action{ID: 1234}},
			},
		})

		respBody, resp, err := getRequest[schema.ActionGetResponse](ctx, client, "/resource")
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NotNil(t, respBody)

		require.Equal(t, int64(1234), respBody.Action.ID)
	})

	t.Run("post", func(t *testing.T) {
		ctx, server, client := makeTestUtils(t)

		server.Expect([]mockutil.Request{
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

		respBody, resp, err := postRequest[schema.ActionGetResponse](ctx, client, "/resource", map[string]string{"hello": "world"})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NotNil(t, respBody)

		require.Equal(t, int64(1234), respBody.Action.ID)
	})

	t.Run("put", func(t *testing.T) {
		ctx, server, client := makeTestUtils(t)

		server.Expect([]mockutil.Request{
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

		respBody, resp, err := putRequest[schema.ActionGetResponse](ctx, client, "/resource", map[string]string{"hello": "world"})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NotNil(t, respBody)

		require.Equal(t, int64(1234), respBody.Action.ID)
	})

	t.Run("delete", func(t *testing.T) {
		ctx, server, client := makeTestUtils(t)

		server.Expect([]mockutil.Request{
			{
				Method: "DELETE", Path: "/resource",
				Status: 200,
				JSON:   schema.ActionGetResponse{Action: schema.Action{ID: 1234}},
			},
		})

		respBody, resp, err := deleteRequest[schema.ActionGetResponse](ctx, client, "/resource")
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NotNil(t, respBody)

		require.Equal(t, int64(1234), respBody.Action.ID)
	})

	t.Run("delete no result", func(t *testing.T) {
		ctx, server, client := makeTestUtils(t)

		server.Expect([]mockutil.Request{
			{
				Method: "DELETE", Path: "/resource",
				Status: 204,
			},
		})

		resp, err := deleteRequestNoResult(ctx, client, "/resource")
		require.NoError(t, err)
		require.NotNil(t, resp)
	})
}
