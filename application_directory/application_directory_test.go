package application_directory

import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestImplementsApplicationDirectory(t *testing.T) {
	object := New(nil)
	var expected *ApplicationDirectory
	err := AssertThat(object, Implements(expected))
	if err != nil {
		t.Fatal(err)
	}
}
