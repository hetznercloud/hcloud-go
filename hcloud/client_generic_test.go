package hcloud

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/hetznercloud/hcloud-go/v2/hcloud/exp/ctxutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/exp/mockutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/schema"
)

func TestGenericRequest(t *testing.T) {
	t.Run("GET", func(t *testing.T) {
		ctx, server, client := makeTestUtils(t)

		server.Expect([]mockutil.Request{
			{
				Method: "GET", Path: "/resource",
				Status: 200,
				JSON:   schema.ActionGetResponse{Action: schema.Action{ID: 1234}},
			},
		})

		const opPath = "/resource"
		ctx = ctxutil.SetOpPath(ctx, opPath)

		respBody, resp, err := getRequest[schema.ActionGetResponse](ctx, client, opPath)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NotNil(t, respBody)

		require.Equal(t, int64(1234), respBody.Action.ID)
	})

	t.Run("POST", func(t *testing.T) {
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

		const opPath = "/resource"
		ctx = ctxutil.SetOpPath(ctx, opPath)

		respBody, resp, err := postRequest[schema.ActionGetResponse](ctx, client, opPath, map[string]string{"hello": "world"})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NotNil(t, respBody)

		require.Equal(t, int64(1234), respBody.Action.ID)
	})

	t.Run("PUT", func(t *testing.T) {
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

		const opPath = "/resource"
		ctx = ctxutil.SetOpPath(ctx, opPath)

		respBody, resp, err := putRequest[schema.ActionGetResponse](ctx, client, opPath, map[string]string{"hello": "world"})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NotNil(t, respBody)

		require.Equal(t, int64(1234), respBody.Action.ID)
	})

	t.Run("DELETE", func(t *testing.T) {
		ctx, server, client := makeTestUtils(t)

		server.Expect([]mockutil.Request{
			{
				Method: "DELETE", Path: "/resource",
				Status: 200,
				JSON:   schema.ActionGetResponse{Action: schema.Action{ID: 1234}},
			},
		})

		const opPath = "/resource"
		ctx = ctxutil.SetOpPath(ctx, opPath)

		respBody, resp, err := deleteRequest[schema.ActionGetResponse](ctx, client, opPath)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NotNil(t, respBody)

		require.Equal(t, int64(1234), respBody.Action.ID)
	})

	t.Run("DELETE no result", func(t *testing.T) {
		ctx, server, client := makeTestUtils(t)

		server.Expect([]mockutil.Request{
			{
				Method: "DELETE", Path: "/resource",
				Status: 204,
			},
		})

		const opPath = "/resource"
		ctx = ctxutil.SetOpPath(ctx, opPath)

		resp, err := deleteRequestNoResult(ctx, client, opPath)
		require.NoError(t, err)
		require.NotNil(t, resp)
	})
}
