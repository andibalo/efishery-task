package voerrors

import (
	"auth-app/internal/response"
	"net/http"
	"testing"
)

func TestMapErrorsToCode(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want response.Code
	}{
		{
			err:  ErrUnauthorized,
			want: response.Unauthorized,
		},
		{
			err:  ErrTimeout,
			want: response.GatewayTimeout,
		},
		{
			err:  ErrBadRequest,
			want: response.BadRequest,
		},
		{
			err:  ErrNotFound,
			want: response.ServerError,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			if got := MapErrorsToCode(tt.err); got != tt.want {
				t.Errorf("MapErrorsToCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMapErrorsToStatusCode(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Bad Request",
			args: args{
				err: ErrBadRequest,
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Timeout",
			args: args{
				err: ErrTimeout,
			},
			want: http.StatusGatewayTimeout,
		},
		{
			name: "Bad Request",
			args: args{
				err: ErrBadRequest,
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Default",
			args: args{
				err: ErrNotFound,
			},
			want: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MapErrorsToStatusCode(tt.args.err); got != tt.want {
				t.Errorf("MapErrorsToHTTPCode() = %v, want %v", got, tt.want)
			}
		})
	}
}
