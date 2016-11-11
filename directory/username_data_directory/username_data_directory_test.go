package username_data_directory

import (
	"testing"

	. "github.com/bborbe/assert"
	"github.com/bborbe/auth/model"
)

func TestImplementsUserDataDirectory(t *testing.T) {
	object := New(nil)
	var expected *UsernameDataDirectory
	err := AssertThat(object, Implements(expected))
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateKey(t *testing.T) {
	key := createKey(model.UserName("test"))
	if err := AssertThat(key, Is("user_data:test")); err != nil {
		t.Fatal(err)
	}
}
