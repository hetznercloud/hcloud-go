package hcloud

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError_Error(t *testing.T) {
	type fields struct {
		Code     ErrorCode
		Message  string
		Details  interface{}
		response *Response
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "simple error",
			fields: fields{
				Code:    ErrorCodeUnauthorized,
				Message: "unable to authenticate",
			},
			want: "unable to authenticate (unauthorized)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Error{
				Code:     tt.fields.Code,
				Message:  tt.fields.Message,
				Details:  tt.fields.Details,
				response: tt.fields.response,
			}
			assert.Equalf(t, tt.want, e.Error(), "Error()")
		})
	}
}

func TestIsError(t *testing.T) {
	type args struct {
		err  error
		code ErrorCode
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "hcloud error with code",
			args: args{
				err:  Error{Code: ErrorCodeUnauthorized},
				code: ErrorCodeUnauthorized,
			},
			want: true,
		},
		{
			name: "hcloud error with different code",
			args: args{
				err:  Error{Code: ErrorCodeConflict},
				code: ErrorCodeUnauthorized,
			},
			want: false,
		},
		{
			name: "wrapped hcloud error with code",
			args: args{
				err:  fmt.Errorf("wrapped: %w", Error{Code: ErrorCodeUnauthorized}),
				code: ErrorCodeUnauthorized,
			},
			want: true,
		},
		{
			name: "wrapped hcloud error with different code",
			args: args{
				err:  fmt.Errorf("wrapped: %w", Error{Code: ErrorCodeConflict}),
				code: ErrorCodeUnauthorized,
			},
			want: false,
		},
		{
			name: "non-hcloud error",
			args: args{
				err:  fmt.Errorf("something went wrong"),
				code: ErrorCodeUnauthorized,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, IsError(tt.args.err, tt.args.code), "IsError(%v, %v)", tt.args.err, tt.args.code)
		})
	}
}
