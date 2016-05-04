package router

import (
	"testing"

	. "github.com/bborbe/assert"

	"net/http"

	"github.com/bborbe/http/requestbuilder"
	"github.com/bborbe/server/mock"
)

func TestHealthz(t *testing.T) {
	counterCheck := 0
	counterLogin := 0
	r := New(Create(&counterCheck), Create(&counterLogin))
	resp := mock.NewHttpResponseWriterMock()

	rb := requestbuilder.NewHttpRequestBuilder("http://example.com/healthz")
	rb.SetMethod("GET")
	req, err := rb.Build()
	if err = AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err = AssertThat(counterCheck, Is(0)); err != nil {
		t.Fatal(err)
	}
	r.ServeHTTP(resp, req)
	if err = AssertThat(counterCheck, Is(1)); err != nil {
		t.Fatal(err)
	}
}

func TestReadiness(t *testing.T) {
	counterCheck := 0
	counterLogin := 0
	r := New(Create(&counterCheck), Create(&counterLogin))
	resp := mock.NewHttpResponseWriterMock()

	rb := requestbuilder.NewHttpRequestBuilder("http://example.com/readiness")
	rb.SetMethod("GET")
	req, err := rb.Build()
	if err = AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err = AssertThat(counterCheck, Is(0)); err != nil {
		t.Fatal(err)
	}
	r.ServeHTTP(resp, req)
	if err = AssertThat(counterCheck, Is(1)); err != nil {
		t.Fatal(err)
	}
}

func TestLogin(t *testing.T) {
	counterCheck := 0
	counterLogin := 0
	r := New(Create(&counterCheck), Create(&counterLogin))
	resp := mock.NewHttpResponseWriterMock()

	rb := requestbuilder.NewHttpRequestBuilder("http://example.com/login")
	rb.SetMethod("POST")
	req, err := rb.Build()
	if err = AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err = AssertThat(counterLogin, Is(0)); err != nil {
		t.Fatal(err)
	}
	r.ServeHTTP(resp, req)
	if err = AssertThat(counterLogin, Is(1)); err != nil {
		t.Fatal(err)
	}
}

func Create(counter *int) func(http.ResponseWriter, *http.Request) {
	return func(http.ResponseWriter, *http.Request) {
		*counter++
	}
}
