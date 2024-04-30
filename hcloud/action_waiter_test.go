package hcloud

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWaitFor(t *testing.T) {
	RunMockedTestCases(t,
		[]MockedTestCase{
			{
				Name: "succeed",
				WantRequests: []MockedRequest{
					{"GET", "/actions?id=1509772237&page=1&sort=status&sort=id", nil, 200,
						`{
							"actions": [
								{ "id": 1509772237, "status": "running", "progress": 0 }
							],
							"meta": { "pagination": { "page": 1 }}
						}`},
					{"GET", "/actions?id=1509772237&page=1&sort=status&sort=id", nil, 200,
						`{
							"actions": [
								{ "id": 1509772237, "status": "success", "progress": 100 }
							],
							"meta": { "pagination": { "page": 1 }}
						}`},
				},
				Run: func(env testEnv) {
					actions := []*Action{{ID: 1509772237, Status: "running"}}

					err := env.Client.Action.WaitFor(context.Background(), actions...)
					assert.NoError(t, err)
				},
			},
			{
				Name: "succeed with already succeeded action",
				Run: func(env testEnv) {
					actions := []*Action{{ID: 1509772237, Status: "success"}}

					err := env.Client.Action.WaitFor(context.Background(), actions...)
					assert.NoError(t, err)
				},
			},
			{
				Name: "fail with unknown action",
				WantRequests: []MockedRequest{
					{"GET", "/actions?id=1509772237&page=1&sort=status&sort=id", nil, 200,
						`{
							"actions": [],
							"meta": { "pagination": { "page": 1 }}
						}`},
				},
				Run: func(env testEnv) {
					actions := []*Action{{ID: 1509772237, Status: "running"}}

					err := env.Client.Action.WaitFor(context.Background(), actions...)
					assert.Error(t, err)
					assert.Equal(t, "actions not found: [1509772237]", err.Error())
				},
			},
			{
				Name: "fail with canceled context",
				Run: func(env testEnv) {
					actions := []*Action{{ID: 1509772237, Status: "running"}}

					ctx, cancelFunc := context.WithCancel(context.Background())
					cancelFunc()
					err := env.Client.Action.WaitFor(ctx, actions...)
					assert.Error(t, err)
				},
			},
			{
				Name: "fail with api error",
				WantRequests: []MockedRequest{
					{"GET", "/actions?id=1509772237&page=1&sort=status&sort=id", nil, 503, ""},
				},
				Run: func(env testEnv) {
					actions := []*Action{{ID: 1509772237, Status: "running"}}

					err := env.Client.Action.WaitFor(context.Background(), actions...)
					assert.Error(t, err)
					assert.Equal(t, "hcloud: server responded with status code 503", err.Error())
				},
			},
		},
	)
}

func TestWaitForFunc(t *testing.T) {
	RunMockedTestCases(t,
		[]MockedTestCase{
			{
				Name: "succeed",
				WantRequests: []MockedRequest{
					{"GET", "/actions?id=1509772237&id=1509772238&page=1&sort=status&sort=id", nil, 200,
						`{
							"actions": [
								{ "id": 1509772237, "status": "running", "progress": 40 },
								{ "id": 1509772238, "status": "running", "progress": 0 }
							],
							"meta": { "pagination": { "page": 1 }}
						}`},
					{"GET", "/actions?id=1509772237&id=1509772238&page=1&sort=status&sort=id", nil, 200,
						`{
							"actions": [
								{ "id": 1509772237, "status": "running", "progress": 60 },
								{ "id": 1509772238, "status": "running", "progress": 50 }
							],
							"meta": { "pagination": { "page": 1 }}
						}`},
					{"GET", "/actions?id=1509772237&id=1509772238&page=1&sort=status&sort=id", nil, 200,
						`{
							"actions": [
								{ "id": 1509772237, "status": "success", "progress": 100 },
								{ "id": 1509772238, "status": "running", "progress": 75 }
							],
							"meta": { "pagination": { "page": 1 }}
						}`},
					{"GET", "/actions?id=1509772238&page=1&sort=status&sort=id", nil, 200,
						`{
							"actions": [
								{ "id": 1509772238, "status": "error", "progress": 75, 
									"error": {
										"code": "action_failed", 
										"message": "Something went wrong with the action"
									}
								}
							],
							"meta": { "pagination": { "page": 1 }}
						}`},
				},
				Run: func(env testEnv) {
					actions := []*Action{
						{ID: 1509772236, Status: "success"},
						{ID: 1509772237, Status: "running"},
						{ID: 1509772238, Status: "running"},
					}
					progress := make([]int, 0)

					progressByAction := make(map[int64]int, len(actions))
					err := env.Client.Action.WaitForFunc(context.Background(), func(update *Action) error {
						switch update.Status {
						case ActionStatusRunning:
							progressByAction[update.ID] = update.Progress
						case ActionStatusSuccess:
							progressByAction[update.ID] = 100
						case ActionStatusError:
							progressByAction[update.ID] = 100
						}

						sum := 0
						for _, value := range progressByAction {
							sum += value
						}
						progress = append(progress, sum/len(actions))

						return nil
					}, actions...)

					assert.Nil(t, err)
					assert.Equal(t, []int{33, 46, 46, 53, 70, 83, 91, 100}, progress)
				},
			},
		},
	)
}
