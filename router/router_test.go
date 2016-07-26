package router

import (
	"testing"

	. "github.com/bborbe/assert"

	"net/http"

	"github.com/bborbe/http/mock"
	"github.com/bborbe/http/requestbuilder"
)

type counter struct {
	notFound   int
	check      int
	v1Router   int
	fileServer int
}

func Create(counter *int) func(http.ResponseWriter, *http.Request) {
	return func(http.ResponseWriter, *http.Request) {
		*counter++
	}
}

func newWithCounter(c *counter) *handler {
	return New(
		"/prefix",
		Create(&c.notFound),
		Create(&c.check),
		Create(&c.fileServer),
		Create(&c.v1Router),
	)
}

func TestHealthz(t *testing.T) {
	c := new(counter)
	r := newWithCounter(c)
	resp := mock.NewHttpResponseWriterMock()

	rb := requestbuilder.NewHttpRequestBuilder("http://example.com/prefix/healthz")
	rb.SetMethod("GET")
	req, err := rb.Build()
	if err = AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err = AssertThat(c.check, Is(0)); err != nil {
		t.Fatal(err)
	}
	r.ServeHTTP(resp, req)
	if err = AssertThat(c.check, Is(1)); err != nil {
		t.Fatal(err)
	}
}

func TestReadiness(t *testing.T) {
	c := new(counter)
	r := newWithCounter(c)
	resp := mock.NewHttpResponseWriterMock()

	rb := requestbuilder.NewHttpRequestBuilder("http://example.com/prefix/readiness")
	rb.SetMethod("GET")
	req, err := rb.Build()
	if err = AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err = AssertThat(c.check, Is(0)); err != nil {
		t.Fatal(err)
	}
	r.ServeHTTP(resp, req)
	if err = AssertThat(c.check, Is(1)); err != nil {
		t.Fatal(err)
	}
}

func TestFiles(t *testing.T) {
	c := new(counter)
	r := newWithCounter(c)
	resp := mock.NewHttpResponseWriterMock()
	rb := requestbuilder.NewHttpRequestBuilder("http://example.com/prefix/foo")
	rb.SetMethod("GET")
	req, err := rb.Build()
	if err = AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err = AssertThat(c.fileServer, Is(0)); err != nil {
		t.Fatal(err)
	}
	r.ServeHTTP(resp, req)
	if err = AssertThat(c.fileServer, Is(1)); err != nil {
		t.Fatal(err)
	}
}

func TestV1Router(t *testing.T) {
	c := new(counter)
	r := newWithCounter(c)
	resp := mock.NewHttpResponseWriterMock()

	rb := requestbuilder.NewHttpRequestBuilder("http://example.com/prefix/api/1.0/version")
	rb.SetMethod("GET")
	req, err := rb.Build()
	if err = AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err = AssertThat(c.v1Router, Is(0)); err != nil {
		t.Fatal(err)
	}
	r.ServeHTTP(resp, req)
	if err = AssertThat(c.v1Router, Is(1)); err != nil {
		t.Fatal(err)
	}
}
