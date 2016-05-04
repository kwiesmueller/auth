package router

import (
	"testing"

	. "github.com/bborbe/assert"

	"net/http"

	"github.com/bborbe/http/requestbuilder"
	"github.com/bborbe/server/mock"
)

type counter struct {
	check             int
	login             int
	applicationCreate int
}

func Create(counter *int) func(http.ResponseWriter, *http.Request) {
	return func(http.ResponseWriter, *http.Request) {
		*counter++
	}
}

func TestHealthz(t *testing.T) {
	c := counter{}
	r := New(Create(&c.check), Create(&c.login), Create(&c.applicationCreate))
	resp := mock.NewHttpResponseWriterMock()

	rb := requestbuilder.NewHttpRequestBuilder("http://example.com/healthz")
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
	c := counter{}
	r := New(Create(&c.check), Create(&c.login), Create(&c.applicationCreate))
	resp := mock.NewHttpResponseWriterMock()

	rb := requestbuilder.NewHttpRequestBuilder("http://example.com/readiness")
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

func TestLogin(t *testing.T) {
	c := counter{}
	r := New(Create(&c.check), Create(&c.login), Create(&c.applicationCreate))
	resp := mock.NewHttpResponseWriterMock()

	rb := requestbuilder.NewHttpRequestBuilder("http://example.com/login")
	rb.SetMethod("POST")
	req, err := rb.Build()
	if err = AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err = AssertThat(c.login, Is(0)); err != nil {
		t.Fatal(err)
	}
	r.ServeHTTP(resp, req)
	if err = AssertThat(c.login, Is(1)); err != nil {
		t.Fatal(err)
	}
}

func TestApplicationCreate(t *testing.T) {
	c := counter{}
	r := New(Create(&c.check), Create(&c.login), Create(&c.applicationCreate))
	resp := mock.NewHttpResponseWriterMock()

	rb := requestbuilder.NewHttpRequestBuilder("http://example.com/application")
	rb.SetMethod("POST")
	req, err := rb.Build()
	if err = AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err = AssertThat(c.applicationCreate, Is(0)); err != nil {
		t.Fatal(err)
	}
	r.ServeHTTP(resp, req)
	if err = AssertThat(c.applicationCreate, Is(1)); err != nil {
		t.Fatal(err)
	}
}
