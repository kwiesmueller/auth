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
	applicationDelete int
	applicationGet    int
	userRegister      int
	userUnregister    int
	tokenAdd          int
	tokenRemove       int
}

func Create(counter *int) func(http.ResponseWriter, *http.Request) {
	return func(http.ResponseWriter, *http.Request) {
		*counter++
	}
}

func newWithCounter(c *counter) *handler {
	return New(Create(&c.check), Create(&c.login), Create(&c.applicationCreate), Create(&c.applicationDelete), Create(&c.applicationGet), Create(&c.userRegister), Create(&c.userUnregister), Create(&c.tokenAdd), Create(&c.tokenRemove))
}

func TestHealthz(t *testing.T) {
	c := new(counter)
	r := newWithCounter(c)
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
	c := new(counter)
	r := newWithCounter(c)
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
	c := new(counter)
	r := newWithCounter(c)
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
	c := new(counter)
	r := newWithCounter(c)
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

func TestApplicationDelete(t *testing.T) {
	c := new(counter)
	r := newWithCounter(c)
	resp := mock.NewHttpResponseWriterMock()
	rb := requestbuilder.NewHttpRequestBuilder("http://example.com/application/test123")
	rb.SetMethod("DELETE")
	req, err := rb.Build()
	if err = AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err = AssertThat(c.applicationDelete, Is(0)); err != nil {
		t.Fatal(err)
	}
	r.ServeHTTP(resp, req)
	if err = AssertThat(c.applicationDelete, Is(1)); err != nil {
		t.Fatal(err)
	}
}

func TestApplicationGet(t *testing.T) {
	c := new(counter)
	r := newWithCounter(c)
	resp := mock.NewHttpResponseWriterMock()
	rb := requestbuilder.NewHttpRequestBuilder("http://example.com/application/test123")
	rb.SetMethod("GET")
	req, err := rb.Build()
	if err = AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err = AssertThat(c.applicationGet, Is(0)); err != nil {
		t.Fatal(err)
	}
	r.ServeHTTP(resp, req)
	if err = AssertThat(c.applicationGet, Is(1)); err != nil {
		t.Fatal(err)
	}
}

func TestUserRegister(t *testing.T) {
	c := new(counter)
	r := newWithCounter(c)
	resp := mock.NewHttpResponseWriterMock()
	rb := requestbuilder.NewHttpRequestBuilder("http://example.com/user")
	rb.SetMethod("POST")
	req, err := rb.Build()
	if err = AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err = AssertThat(c.userRegister, Is(0)); err != nil {
		t.Fatal(err)
	}
	r.ServeHTTP(resp, req)
	if err = AssertThat(c.userRegister, Is(1)); err != nil {
		t.Fatal(err)
	}
}

func TestUnregister(t *testing.T) {
	c := new(counter)
	r := newWithCounter(c)
	resp := mock.NewHttpResponseWriterMock()
	rb := requestbuilder.NewHttpRequestBuilder("http://example.com/user")
	rb.SetMethod("DELETE")
	req, err := rb.Build()
	if err = AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err = AssertThat(c.userUnregister, Is(0)); err != nil {
		t.Fatal(err)
	}
	r.ServeHTTP(resp, req)
	if err = AssertThat(c.userUnregister, Is(1)); err != nil {
		t.Fatal(err)
	}
}

func TestTokenAdd(t *testing.T) {
	c := new(counter)
	r := newWithCounter(c)
	resp := mock.NewHttpResponseWriterMock()
	rb := requestbuilder.NewHttpRequestBuilder("http://example.com/token")
	rb.SetMethod("POST")
	req, err := rb.Build()
	if err = AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err = AssertThat(c.tokenAdd, Is(0)); err != nil {
		t.Fatal(err)
	}
	r.ServeHTTP(resp, req)
	if err = AssertThat(c.tokenAdd, Is(1)); err != nil {
		t.Fatal(err)
	}
}

func TestTokenRemove(t *testing.T) {
	c := new(counter)
	r := newWithCounter(c)
	resp := mock.NewHttpResponseWriterMock()
	rb := requestbuilder.NewHttpRequestBuilder("http://example.com/token")
	rb.SetMethod("DELETE")
	req, err := rb.Build()
	if err = AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err = AssertThat(c.tokenRemove, Is(0)); err != nil {
		t.Fatal(err)
	}
	r.ServeHTTP(resp, req)
	if err = AssertThat(c.tokenRemove, Is(1)); err != nil {
		t.Fatal(err)
	}
}
