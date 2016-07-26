package user_data_get

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

func TestHandleRequest(t *testing.T) {
	counter := 0
	h := New(func(userName model.UserName) (map[string]string, error) {
		counter++
		if err := AssertThat(string(userName), Is("tester")); err != nil {
			t.Fatal(err)
		}
		return map[string]string{"keya": "valuea"}, nil
	})
	resp := mock.NewHttpResponseWriterMock()
	rb := requestbuilder.NewHttpRequestBuilder("http://example.com/user/tester/data")
	rb.SetMethod("GET")
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
	if err = AssertThat(string(resp.Bytes()), Is(fmt.Sprintln(`{"keya":"valuea"}`))); err != nil {
		t.Fatal(err)
	}
}
