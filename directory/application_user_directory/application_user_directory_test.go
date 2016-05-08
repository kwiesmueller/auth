package application_user_directory

import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestImplementsApplicationUserDirectory(t *testing.T) {
	object := New()
	var expected *ApplicationUserDirectory
	err := AssertThat(object, Implements(expected))
	if err != nil {
		t.Fatal(err)
	}
}
