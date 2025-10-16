package internal

import (
	"net/url"
	"testing"

	"go.llib.dev/testcase/assert"
)

func TestCanonicalPathQuery(t *testing.T) {
	tests := []struct {
		name     string
		inputURL string
		want     string
	}{
		{
			name:     "order query params alphabetically",
			inputURL: "/test1?product=pixel&color=black",
			want:     "/test1?color=black&product=pixel",
		},
		{
			name:     "no query params just path",
			inputURL: "/simple",
			want:     "/simple",
		},
		{
			name:     "empty path becomes slash",
			inputURL: "",
			want:     "/",
		},
		{
			name:     "multiple values for key are sorted",
			inputURL: "/multi?foo=b&foo=a",
			want:     "/multi?foo=a&foo=b",
		},
		{
			name:     "sorted keys with multiple values",
			inputURL: "/mixed?z=3&a=2&z=1",
			want:     "/mixed?a=2&z=1&z=3",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u, err := url.Parse(tc.inputURL)
			assert.Must(t).Nil(err)

			got := canonicalPathQuery(u)
			assert.Must(t).Equal(tc.want, got)
		})
	}
}
