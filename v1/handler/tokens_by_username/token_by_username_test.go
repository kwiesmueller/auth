package tokens_by_username

import (
	"net/http"
	"testing"

	"fmt"

	. "github.com/bborbe/assert"
	"github.com/bborbe/auth/model"
	"github.com/bborbe/http/mock"
	"github.com/bborbe/http/requestbuilder"
)

func TestImplementsHandler(t *testing.T) {
	object := New(nil)
	var expected *http.Handler
	err := AssertThat(object, Implements(expected))
	if err != nil {
		t.Fatal(err)
	}
}

func TestResponseError(t *testing.T) {
	handler := New(func(model.UserName) ([]model.AuthToken, error) {
		return nil, fmt.Errorf("foo")
	})
	resp := mock.NewHttpResponseWriterMock()
	rb := requestbuilder.NewHTTPRequestBuilder("http://example.com/prefix/list?username=foo")
	req, err := rb.Build()
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	err = handler.serveHTTP(resp, req)
	if err := AssertThat(err, NotNilValue()); err != nil {
		t.Fatal(err)
	}
}

func TestResponseSuccess(t *testing.T) {
	handler := New(func(username model.UserName) ([]model.AuthToken, error) {
		if err := AssertThat(username.String(), Is("foo")); err != nil {
			t.Fatal(err)
		}
		return []model.AuthToken{"foo"}, nil
	})
	resp := mock.NewHttpResponseWriterMock()
	rb := requestbuilder.NewHTTPRequestBuilder("http://example.com/prefix/list?username=foo")
	req, err := rb.Build()
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	err = handler.serveHTTP(resp, req)
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(string(resp.Bytes()), Is("[\"foo\"]\n")); err != nil {
		t.Fatal(err)
	}
}
