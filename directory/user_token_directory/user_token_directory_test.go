package user_token_directory

import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestImplementsUserTokenDirectory(t *testing.T) {
	object := New(nil)
	var expected *UserTokenDirectory
	err := AssertThat(object, Implements(expected))
	if err != nil {
		t.Fatal(err)
	}
}
