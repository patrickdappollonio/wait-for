package wait

import (
	"net/url"
	"reflect"
	"testing"
)

func TestParseHost(t *testing.T) {
	tests := []struct {
		name    string
		hostStr string
		want    *url.URL
		wantErr bool
	}{
		{
			name:    "No scheme",
			hostStr: "example.com",
			want:    &url.URL{Scheme: "tcp", Host: "example.com"},
			wantErr: false,
		},
		{
			name:    "With scheme",
			hostStr: "http://example.com",
			want:    &url.URL{Scheme: "http", Host: "example.com"},
			wantErr: false,
		},
		{
			name:    "Invalid URL",
			hostStr: "://example.com",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseHost(tt.hostStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseHost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseHost() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringifyHosts(t *testing.T) {
	tests := []struct {
		name string
		urls []matchedURLItem
		want string
	}{
		{
			name: "Single URL",
			urls: []matchedURLItem{
				{
					Raw: "http://example.com",
				},
			},
			want: `"http://example.com"`,
		},
		{
			name: "Multiple URLs",
			urls: []matchedURLItem{
				{
					Raw: "http://example.com",
				},
				{
					Raw: "https://example.org",
				},
			},
			want: `"http://example.com", "https://example.org"`,
		},
		{
			name: "No URLs",
			urls: []matchedURLItem{},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stringifyHosts(tt.urls); got != tt.want {
				t.Errorf("stringifyHosts() = %v, want %v", got, tt.want)
			}
		})
	}
}
