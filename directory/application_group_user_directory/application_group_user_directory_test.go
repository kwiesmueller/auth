package application_group_user_directory

import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestImplementsApplicationGroupDirectory(t *testing.T) {
	object := New()
	var expected *ApplicationGroupUserDirectory
	err := AssertThat(object, Implements(expected))
	if err != nil {
		t.Fatal(err)
	}
}
