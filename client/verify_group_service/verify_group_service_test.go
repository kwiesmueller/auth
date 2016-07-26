package verify_group_service

import (
	"net/http"
	"testing"

	. "github.com/bborbe/assert"
	"github.com/bborbe/auth/model"
	http_requestbuilder "github.com/bborbe/http/requestbuilder"
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
	httpRequestBuilderProvider := http_requestbuilder.NewHttpRequestBuilderProvider()
	c := New(func(req *http.Request) (resp *http.Response, err error) {
		counter++
		if err := AssertThat(req.URL.String(), Is("http://auth-api.auth.svc.cluster.local:8080/auth/api/1.0/login")); err != nil {
			t.Fatal(err)
		}
		if err := AssertThat(req.Method, Is("POST")); err != nil {
			t.Fatal(err)
		}
		return &http.Response{}, nil
	}, httpRequestBuilderProvider, "http://auth-api.auth.svc.cluster.local:8080/auth", "", "")
	if err := AssertThat(counter, Is(0)); err != nil {
		t.Fatal(err)
	}
	c.Auth(model.AuthToken("abc"), []model.GroupName{})
	if err := AssertThat(counter, Is(1)); err != nil {
		t.Fatal(err)
	}
}
