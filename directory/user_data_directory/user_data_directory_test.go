package user_data_directory

import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestImplementsUserDataDirectory(t *testing.T) {
	object := New(nil)
	var expected *UserDataDirectory
	err := AssertThat(object, Implements(expected))
	if err != nil {
		t.Fatal(err)
	}
}
