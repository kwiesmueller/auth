package user_list

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
	handler := New(func() ([]model.UserName, error) {
		return nil, fmt.Errorf("foo")
	})
	resp := mock.NewHttpResponseWriterMock()
	rb := requestbuilder.NewHttpRequestBuilder("http://example.com/prefix/healthz")
	req, err := rb.Build()
	if AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	err = handler.serveHTTP(resp, req)
	if err = AssertThat(err, NotNilValue()); err != nil {
		t.Fatal(err)
	}
}

func TestResponseSuccess(t *testing.T) {
	handler := New(func() ([]model.UserName, error) {
		return []model.UserName{model.UserName("foo")}, nil
	})
	resp := mock.NewHttpResponseWriterMock()
	rb := requestbuilder.NewHttpRequestBuilder("http://example.com/prefix/healthz")
	req, err := rb.Build()
	if AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	err = handler.serveHTTP(resp, req)
	if AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(string(resp.Bytes()), Is("[{\"username\":\"foo\"}]\n")); err != nil {
		t.Fatal(err)
	}
}
