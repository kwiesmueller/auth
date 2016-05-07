package token_user_directory

import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestImplementsTokenUserDirectory(t *testing.T) {
	object := New(nil)
	var expected *TokenUserDirectory
	err := AssertThat(object, Implements(expected))
	if err != nil {
		t.Fatal(err)
	}
}
