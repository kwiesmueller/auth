package user_directory

import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestImplementsUserDirectory(t *testing.T) {
	object := New()
	var expected *UserDirectory
	err := AssertThat(object, Implements(expected))
	if err != nil {
		t.Fatal(err)
	}
}
