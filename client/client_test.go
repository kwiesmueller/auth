package client

import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestImplementsAuthClient(t *testing.T) {
	object := New(nil, nil, "", "", "")
	var expected *AuthClient
	err := AssertThat(object, Implements(expected))
	if err != nil {
		t.Fatal(err)
	}
}
