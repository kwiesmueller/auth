package user_group_directory

import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestImplementsGroupDirectory(t *testing.T) {
	object := New(nil)
	var expected *UserGroupDirectory
	err := AssertThat(object, Implements(expected))
	if err != nil {
		t.Fatal(err)
	}
}
