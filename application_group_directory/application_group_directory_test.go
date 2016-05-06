package application_group_directory

import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestImplementsApplicationGroupDirectory(t *testing.T) {
	object := New()
	var expected *ApplicationGroupDirectory
	err := AssertThat(object, Implements(expected))
	if err != nil {
		t.Fatal(err)
	}
}
