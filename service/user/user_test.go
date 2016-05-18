package user

import (
	"testing"

	. "github.com/bborbe/assert"
	"github.com/bborbe/auth/api"
	"github.com/bborbe/auth/directory/token_user_directory"
	"github.com/bborbe/auth/directory/user_data_directory"
	"github.com/bborbe/auth/directory/user_group_directory"
	"github.com/bborbe/auth/directory/user_token_directory"
	"github.com/bborbe/ledis/mock"
)

func TestImplementsService(t *testing.T) {
	object := New(nil, nil, nil, nil)
	var expected *Service
	err := AssertThat(object, Implements(expected))
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeleteUserData(t *testing.T) {
	userTokenDirectory := user_token_directory.New(mock.NewSet())
	userGroupDirectory := user_group_directory.New(mock.NewSet())
	tokenUserDirectory := token_user_directory.New(mock.NewKv())
	userDataDirectory := user_data_directory.New(mock.NewHash())
	userService := New(userTokenDirectory, userGroupDirectory, tokenUserDirectory, userDataDirectory)

	user := api.UserName("test")
	key := "key"
	value := "value"
	if err := userDataDirectory.SetValue(user, key, value); err != nil {
		t.Fatal(err)
	}
	v, err := userDataDirectory.GetValue(user, key)
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(v, Is(value)); err != nil {
		t.Fatal(err)
	}

	if err := userService.DeleteUser(user); err != nil {
		t.Fatal(err)
	}
	_, err = userDataDirectory.GetValue(user, key)
	if err := AssertThat(err, NotNilValue()); err != nil {
		t.Fatal(err)
	}
}
