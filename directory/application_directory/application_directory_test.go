package application_directory

import (
	"testing"

	. "github.com/bborbe/assert"
	"github.com/bborbe/auth/model"
)

func TestImplementsApplicationDirectory(t *testing.T) {
	object := New(nil)
	var expected *ApplicationDirectory
	err := AssertThat(object, Implements(expected))
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateKey(t *testing.T) {
	key := createKey(model.ApplicationName("test"))
	if err := AssertThat(key, Is("application:test")); err != nil {
		t.Fatal(err)
	}
}
