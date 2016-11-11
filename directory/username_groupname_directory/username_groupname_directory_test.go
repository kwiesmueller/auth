package username_groupname_directory

import (
	"testing"

	. "github.com/bborbe/assert"
	"github.com/bborbe/auth/model"
)

func TestImplementsGroupDirectory(t *testing.T) {
	object := New(nil)
	var expected *UsernameGroupnameDirectory
	err := AssertThat(object, Implements(expected))
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateKey(t *testing.T) {
	key := createKey(model.UserName("test"))
	if err := AssertThat(key, Is("user_group:test")); err != nil {
		t.Fatal(err)
	}
}
