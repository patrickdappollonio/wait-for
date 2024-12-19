package probes

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOneOf(t *testing.T) {
	tests := []struct {
		name   string
		s      interface{}
		values []interface{}
		want   bool
	}{
		{
			name:   "String present in list",
			s:      "apple",
			values: []interface{}{"banana", "apple", "cherry"},
			want:   true,
		},
		{
			name:   "String not present in list",
			s:      "grape",
			values: []interface{}{"banana", "apple", "cherry"},
			want:   false,
		},
		{
			name:   "Integer present in list",
			s:      42,
			values: []interface{}{1, 2, 42, 100},
			want:   true,
		},
		{
			name:   "Integer not present in list",
			s:      99,
			values: []interface{}{1, 2, 42, 100},
			want:   false,
		},
		{
			name:   "Empty list",
			s:      "test",
			values: []interface{}{},
			want:   false,
		},
		{
			name:   "Single element list, present",
			s:      "single",
			values: []interface{}{"single"},
			want:   true,
		},
		{
			name:   "Single element list, not present",
			s:      "single",
			values: []interface{}{"not_single"},
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := oneOf(tt.s, tt.values...)
			if got != tt.want {
				t.Errorf("oneOf() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnwrapError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want error
	}{
		{
			name: "No wrapping",
			err:  errors.New("root error"),
			want: errors.New("root error"),
		},
		{
			name: "Single wrapping",
			err:  fmt.Errorf("wrapped: %w", errors.New("root error")),
			want: errors.New("root error"),
		},
		{
			name: "Double wrapping",
			err:  fmt.Errorf("wrapped again: %w", fmt.Errorf("wrapped: %w", errors.New("root error"))),
			want: errors.New("root error"),
		},
		{
			name: "Nil error",
			err:  nil,
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := unwrapError(tt.err)
			if got == nil && tt.want != nil {
				t.Errorf("unwrapError() = nil, want %v", tt.want)
			} else if got != nil && tt.want == nil {
				t.Errorf("unwrapError() = %v, want nil", got)
			} else if got != nil && got.Error() != tt.want.Error() {
				t.Errorf("unwrapError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDoGet(t *testing.T) {
	tests := []struct {
		name       string
		url        string
		statusCode int
		wantErr    bool
		errMsg     string
	}{
		{
			name:       "Status OK",
			url:        "/",
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "Status Created",
			url:        "/",
			statusCode: http.StatusCreated,
			wantErr:    false,
		},
		{
			name:       "Status Bad Request",
			url:        "/",
			statusCode: http.StatusBadRequest,
			wantErr:    true,
			errMsg:     "received non-2xx status code: 400 Bad Request",
		},
		{
			name:       "Status Internal Server Error",
			url:        "/",
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
			errMsg:     "received non-2xx status code: 500 Internal Server Error",
		},
		{
			name:    "Invalid URL",
			url:     "http://[::1]:namedport",
			wantErr: true,
			errMsg:  "error creating request: parse \"http://[::1]:namedport\": invalid port \":namedport\" after host",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test server that responds with the specified status code
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			client := server.Client()
			ctx := context.Background()

			url := server.URL + tt.url
			if tt.name == "Invalid URL" {
				url = tt.url // Use the invalid URL directly
			}

			err := doGet(ctx, client, url)
			if (err != nil) != tt.wantErr {
				t.Errorf("doGet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && err.Error() != tt.errMsg {
				t.Errorf("doGet() error = %v, errMsg %v", err, tt.errMsg)
			}
		})
	}
}
