package hcloud

import (
	"fmt"
	"net/http"
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

func TestStabilizeError(t *testing.T) {
	tests := []struct {
		name     string
		response *Response
	}{
		{
			name:     "error without correlation ID",
			response: nil,
		},
		{
			name: "error with correlation ID",
			response: &Response{
				Response: &http.Response{
					StatusCode: http.StatusNotFound,
					Header: func() http.Header {
						headers := http.Header{}
						// [http.Header] requires normalized header names, easiest to do by using the Set method
						headers.Set("X-Correlation-ID", "d5064a1f0bb9de4b")
						return headers
					}(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Error{
				Code:     ErrorCodeNotFound,
				Message:  "server not found",
				response: tt.response,
			}

			err := StabilizeError(e)

			assert.Equal(t, "server not found (not_found)", fmt.Sprintf("%v", err))
		})
	}
}

func TestArgumentError(t *testing.T) {
	type Something struct{ Name string }
	something := Something{Name: "hello"}
	var somethingNil *Something

	for _, tt := range []struct {
		err  error
		want string
	}{
		{
			err:  invalidArgument("something", something, fmt.Errorf("a custom error")),
			want: "invalid argument 'something' [hcloud.Something]: a custom error",
		},
		{
			err:  invalidArgument("something", somethingNil, emptyValue(somethingNil)),
			want: "invalid argument 'something' [*hcloud.Something]: empty value '<nil>'",
		},
		// Validation errors:
		{
			err:  invalidValue("hello"),
			want: "invalid value 'hello'",
		},
		{
			err:  invalidValue(somethingNil),
			want: "invalid value '<nil>'",
		},
		{
			err:  invalidValue(something),
			want: "invalid value '{hello}'",
		},
		{
			err:  emptyValue(""),
			want: "empty value ''",
		},
		{
			err:  emptyValue(somethingNil),
			want: "empty value '<nil>'",
		},
		{
			err:  missingField(something, "Name"),
			want: "missing field [Name] in [hcloud.Something]",
		},
		{
			err:  invalidFieldValue(something, "Name", something.Name),
			want: "invalid value 'hello' for field [Name] in [hcloud.Something]",
		},
		{
			err:  missingOneOfFields(something, "Name", "Type"),
			want: "missing one of fields [Name, Type] in [hcloud.Something]",
		},
		{
			err:  mutuallyExclusiveFields(something, "Name", "Type"),
			want: "found mutually exclusive fields [Name, Type] in [hcloud.Something]",
		},
		{
			err:  missingRequiredTogetherFields(something, "Name", "Type"),
			want: "missing required together fields [Name, Type] in [hcloud.Something]",
		},
	} {
		t.Run("", func(t *testing.T) {
			assert.EqualError(t, tt.err, tt.want)
		})
	}
}
