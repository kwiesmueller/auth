package application_directory

import (
	"testing"

	. "github.com/bborbe/assert"
	"github.com/bborbe/auth/api"
)

func TestImplementsApplicationDirectory(t *testing.T) {
	object := New(api.ApplicationPassword(""))
	var expected *ApplicationDirectory
	err := AssertThat(object, Implements(expected))
	if err != nil {
		t.Fatal(err)
	}
}
