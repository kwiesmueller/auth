package router

import (
	"testing"

	. "github.com/bborbe/assert"

	"net/http"

	"github.com/bborbe/server/mock"
)

func TestLogin(t *testing.T) {
	counter := 0
	r := New(func(http.ResponseWriter, *http.Request) {
		counter++
	}, nil)
	resp := mock.NewHttpResponseWriterMock()
	req, err := mock.NewHttpRequestMock("http://example.com")
	if err = AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	r.ServeHTTP(resp, req)
}
