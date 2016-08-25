package verify_group_service

import (
	"net/http"
	"testing"

	"bytes"
	. "github.com/bborbe/assert"
	"github.com/bborbe/auth/model"
	http_requestbuilder "github.com/bborbe/http/requestbuilder"
	"github.com/bborbe/io/reader_nop_close"
)

func TestImplementsAuthClient(t *testing.T) {
	object := New(nil, nil, "", "", "")
	var expected *AuthClient
	err := AssertThat(object, Implements(expected))
	if err != nil {
		t.Fatal(err)
	}
}

func TestRequest(t *testing.T) {
	counter := 0
	httpRequestBuilderProvider := http_requestbuilder.NewHTTPRequestBuilderProvider()
	c := New(func(req *http.Request) (resp *http.Response, err error) {
		counter++
		if err := AssertThat(req.URL.String(), Is("http://auth-api.auth.svc.cluster.local:8080/auth/api/1.0/login")); err != nil {
			t.Fatal(err)
		}
		if err := AssertThat(req.Method, Is("POST")); err != nil {
			t.Fatal(err)
		}
		return &http.Response{
			StatusCode: 200,
			Body:       reader_nop_close.New(bytes.NewBufferString("{}")),
		}, nil
	}, httpRequestBuilderProvider, "http://auth-api.auth.svc.cluster.local:8080/auth", "", "")
	if err := AssertThat(counter, Is(0)); err != nil {
		t.Fatal(err)
	}
	_, err := c.Auth(model.AuthToken("abc"), []model.GroupName{})
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(counter, Is(1)); err != nil {
		t.Fatal(err)
	}
}
