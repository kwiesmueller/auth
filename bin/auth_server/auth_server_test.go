package main

import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestServerFails(t *testing.T) {
	_, err := createServer(0)
	if err := AssertThat(err, NotNilValue()); err != nil {
		t.Fatal(err)
	}
}

func TestServerSuccess(t *testing.T) {
	srv, err := createServer(1337)
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(srv, NotNilValue()); err != nil {
		t.Fatal(err)
	}
}
