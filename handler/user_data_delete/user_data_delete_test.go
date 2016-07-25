package user_data_delete

import (
	"net/http"
	"testing"

	. "github.com/bborbe/assert"
	"github.com/bborbe/auth/api"
	"github.com/bborbe/http/mock"
	"github.com/bborbe/http/requestbuilder"
)

func TestImplementsHandler(t *testing.T) {
	h := New(nil)
	var expected *http.Handler
	err := AssertThat(h, Implements(expected))
	if err != nil {
		t.Fatal(err)
	}
}

func TestHandleRequest(t *testing.T) {
	counter := 0
	h := New(func(userName api.UserName) error {
		counter++
		if err := AssertThat(string(userName), Is("tester")); err != nil {
			t.Fatal(err)
		}
		return nil
	})
	resp := mock.NewHttpResponseWriterMock()
	rb := requestbuilder.NewHttpRequestBuilder("http://example.com/user/tester/data")
	rb.SetMethod("DELETE")
	req, err := rb.Build()
	if err = AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	err = h.serveHTTP(resp, req)
	if err = AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err = AssertThat(counter, Is(1)); err != nil {
		t.Fatal(err)
	}
}
