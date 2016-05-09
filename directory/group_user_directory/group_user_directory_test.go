package group_user_directory

import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestImplementsGroupDirectory(t *testing.T) {
	object := New(nil)
	var expected *GroupUserDirectory
	err := AssertThat(object, Implements(expected))
	if err != nil {
		t.Fatal(err)
	}
}
