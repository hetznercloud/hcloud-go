package hcloud

import (
	"net/http/httptest"
	"testing"

	"github.com/hetznercloud/hcloud-go/v2/hcloud/exp/mockutils"
)

type MockedTestCase struct {
	Name         string
	WantRequests []mockutils.Request
	Run          func(env testEnv)
}

func RunMockedTestCases(t *testing.T, testCases []MockedTestCase) {
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			testServer := httptest.NewServer(mockutils.Handler(t, testCase.WantRequests))

			env := newTestEnvWithServer(testServer, nil)
			defer env.Teardown()

			testCase.Run(env)
		})
	}
}
