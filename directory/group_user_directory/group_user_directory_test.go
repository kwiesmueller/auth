package group_user_directory

import (
	"testing"

	. "github.com/bborbe/assert"
	"github.com/bborbe/auth/model"
)

func TestImplementsGroupDirectory(t *testing.T) {
	object := New(nil)
	var expected *GroupUserDirectory
	err := AssertThat(object, Implements(expected))
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateKey(t *testing.T) {
	key := createKey(model.GroupName("test"))
	if err := AssertThat(key, Is("group_user:test")); err != nil {
		t.Fatal(err)
	}
}
