package internal_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"redirectware/internal"
	"redirectware/storage/inmem"
	"testing"

	"go.llib.dev/testcase/assert"
)

func TestHandler_ServeHTTP(t *testing.T) {
	storage := inmem.New()

	rule1, err := internal.NewFullMatchRule("/test1", "/shopware/test1")
	assert.Must(t).Nil(err)
	rule2, err := internal.NewFullMatchRule("/test2", "/shopware/test2")
	assert.Must(t).Nil(err)

	storage.AddFullMatchRule(context.Background(), rule1)
	storage.AddFullMatchRule(context.Background(), rule2)
	storage.SetWelcomePageURL(context.Background(), "https://www.google.com")

	handler := internal.NewHandler(storage)

	tests := []struct {
		name string
		path string
		want string
	}{
		{name: "test1", path: "/test1", want: "/shopware/test1"},
		{name: "test2", path: "/test2", want: "/shopware/test2"},
		{name: "test3", path: "/test3", want: "https://www.google.com"},
	}

	for _, test := range tests {
		req, err := http.NewRequest("GET", test.path, nil)
		assert.Must(t).NoError(err)

		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		assert.Must(t).Equal(w.Code, http.StatusMovedPermanently)
		assert.Must(t).Equal(w.Header().Get("Location"), test.want)
	}

}
