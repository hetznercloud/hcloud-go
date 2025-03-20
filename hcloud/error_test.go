package hcloud

import (
	"fmt"
	"net/http"
	"reflect"
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
		{
			name: "internal server error with correlation id",
			fields: fields{
				Code:    ErrorCodeUnknownError,
				Message: "Creating image failed because of an unknown error.",
				response: &Response{
					Response: &http.Response{
						StatusCode: http.StatusInternalServerError,
						Header: func() http.Header {
							headers := http.Header{}
							// [http.Header] requires normalized header names, easiest to do by using the Set method
							headers.Set("X-Correlation-ID", "d5064a1f0bb9de4b")
							return headers
						}(),
					},
				},
			},
			want: "Creating image failed because of an unknown error. (unknown_error, d5064a1f0bb9de4b)",
		},
		{
			name: "internal server error without correlation id",
			fields: fields{
				Code:    ErrorCodeUnknownError,
				Message: "Creating image failed because of an unknown error.",
				response: &Response{
					Response: &http.Response{
						StatusCode: http.StatusInternalServerError,
					},
				},
			},
			want: "Creating image failed because of an unknown error. (unknown_error)",
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
		code []ErrorCode
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
				code: []ErrorCode{ErrorCodeUnauthorized},
			},
			want: true,
		},
		{
			name: "hcloud error with many codes",
			args: args{
				err:  Error{Code: ErrorCodeUnauthorized},
				code: []ErrorCode{ErrorCodeUnauthorized, ErrorCodeConflict},
			},
			want: true,
		},
		{
			name: "hcloud error with different code",
			args: args{
				err:  Error{Code: ErrorCodeConflict},
				code: []ErrorCode{ErrorCodeUnauthorized},
			},
			want: false,
		},
		{
			name: "wrapped hcloud error with code",
			args: args{
				err:  fmt.Errorf("wrapped: %w", Error{Code: ErrorCodeUnauthorized}),
				code: []ErrorCode{ErrorCodeUnauthorized},
			},
			want: true,
		},
		{
			name: "wrapped hcloud error with different code",
			args: args{
				err:  fmt.Errorf("wrapped: %w", Error{Code: ErrorCodeConflict}),
				code: []ErrorCode{ErrorCodeUnauthorized},
			},
			want: false,
		},
		{
			name: "non-hcloud error",
			args: args{
				err:  fmt.Errorf("something went wrong"),
				code: []ErrorCode{ErrorCodeUnauthorized},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, IsError(tt.args.err, tt.args.code...), "IsError(%v, %v)", tt.args.err, tt.args.code)
		})
	}
}

func TestInvalidArgumentError(t *testing.T) {
	type Something struct{ Name string }
	something := Something{Name: "hello"}

	for _, tt := range []struct {
		err  error
		want string
	}{
		{
			err: &InvalidArgumentError{
				Message: "unexpected nil value",
				Struct:  reflect.TypeOf(something).String(),
			},
			want: "unexpected nil value for struct [hcloud.Something]",
		},
		{
			err: &InvalidArgumentError{
				Message: "unexpected empty value",
				Struct:  reflect.TypeOf(something).String(),
				Field:   "Name",
			},
			want: "unexpected empty value for field [Name] in struct [hcloud.Something]",
		},
		{
			err: &InvalidArgumentError{
				Message: "unexpected value",
				Value:   "invalid",
				Struct:  reflect.TypeOf(something).String(),
				Field:   "Name",
			},
			want: "unexpected value 'invalid' for field [Name] in struct [hcloud.Something]",
		},
		{
			err: &InvalidArgumentError{
				Message: "invalid argument",
			},
			want: "invalid argument",
		},
	} {
		t.Run("", func(t *testing.T) {
			assert.EqualError(t, tt.err, tt.want)
		})
	}
}
