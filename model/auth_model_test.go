package model

import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestAuthTokenByUsernameAndPassword(t *testing.T) {
	token := AuthTokenByUsernameAndPassword("user123", "pass123")
	err := AssertThat(token.String(), Is("dXNlcjEyMzpwYXNzMTIz"))
	if err != nil {
		t.Fatal(err)
	}
}
