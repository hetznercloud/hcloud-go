package errutil

import (
	"errors"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestLogValue(t *testing.T) {
	testCases := []struct {
		desc string
		err  error
		want string
	}{
		{
			desc: "random error",
			err:  errors.New(`random error`),
			want: `err="random error"`,
		},
		{
			desc: "api error without details",
			err: &hcloud.Error{
				Code:    `unauthorized`,
				Message: `unauthorized request`,
			},
			want: `err.msg="unauthorized request" err.code=unauthorized`,
		},
		{
			desc: "api error with invalid input single message",
			err: &hcloud.Error{
				Code:    `invalid_input`,
				Message: `invalid input for "labels"`,
				Details: hcloud.ErrorDetailsInvalidInput{
					Fields: []hcloud.ErrorDetailsInvalidInputField{
						{Name: "labels", Messages: []string{"value is too long"}},
					},
				},
			},
			want: `err.msg="invalid input for \"labels\"" err.code=invalid_input err.details="[{labels [value is too long]}]"`,
		},
		{
			desc: "api error with invalid input many message",
			err: &hcloud.Error{
				Code:    `invalid_input`,
				Message: `invalid input for "name", "labels"`,
				Details: hcloud.ErrorDetailsInvalidInput{
					Fields: []hcloud.ErrorDetailsInvalidInputField{
						{Name: "name", Messages: []string{"value is too long"}},
						{Name: "labels", Messages: []string{"value is too long"}},
					},
				},
			},
			want: `err.msg="invalid input for \"name\", \"labels\"" err.code=invalid_input err.details="[{name [value is too long]} {labels [value is too long]}]"`,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.desc, func(t *testing.T) {
			var out strings.Builder
			logger := slog.New(slog.NewTextHandler(&out, nil))
			r := slog.NewRecord(time.Time{}, slog.LevelInfo, "test", 0)
			r.Add("err", LogValue(testCase.err))
			logger.Handler().Handle(t.Context(), r)

			// Remove leading log level and msg
			got := strings.TrimSpace(out.String()[20:])
			assert.Equal(t, testCase.want, got)
		})
	}
}
