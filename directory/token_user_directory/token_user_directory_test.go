package token_user_directory

import (
	"testing"

	. "github.com/bborbe/assert"
	"github.com/bborbe/auth/model"
)

func TestImplementsTokenUserDirectory(t *testing.T) {
	object := New(nil)
	var expected *TokenUserDirectory
	err := AssertThat(object, Implements(expected))
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateKey(t *testing.T) {
	key := createKey(model.AuthToken("test"))
	if err := AssertThat(key, Is("token_user:test")); err != nil {
		t.Fatal(err)
	}
}
