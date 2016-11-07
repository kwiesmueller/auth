package rest

import (
	"testing"

	"os"

	"net/http"

	. "github.com/bborbe/assert"
	"github.com/bborbe/http/header"
	"github.com/golang/glog"
)

func TestMain(m *testing.M) {
	exit := m.Run()
	glog.Flush()
	os.Exit(exit)
}

func TestNew(t *testing.T) {
	c := New(nil, "http://example.com", "user123", "password123")
	if err := AssertThat(c, NotNilValue()); err != nil {
		t.Fatal(err)
	}
}

func TestToken(t *testing.T) {
	callcounter := 0
	c := New(func(url string, method string, request interface{}, response interface{}, h http.Header) error {
		callcounter++
		value := h.Get("Authorization")
		if err := AssertThat(len(value) > 0, Is(true)); err != nil {
			t.Fatal(err)
		}
		value, err := header.ParseAuthorizationHeaderSimple("Bearer", value)
		if err := AssertThat(err, NilValue()); err != nil {
			t.Fatal(err)
		}
		if err := AssertThat(value, Is("user123:password123")); err != nil {
			t.Fatal(err)
		}
		return nil
	}, "http://example.com", "user123", "password123")
	err := c.Call("/path", http.MethodGet, nil, nil)
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(callcounter, Is(1)); err != nil {
		t.Fatal(err)
	}
}
