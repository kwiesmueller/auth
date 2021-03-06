package user

import (
	"testing"

	. "github.com/bborbe/assert"
	"github.com/bborbe/auth/directory/token_username_directory"
	"github.com/bborbe/auth/directory/username_data_directory"
	"github.com/bborbe/auth/directory/username_groupname_directory"
	"github.com/bborbe/auth/directory/username_token_directory"
	"github.com/bborbe/auth/model"
	"github.com/bborbe/auth/service"
	ledis "github.com/bborbe/redis_client/mock"
)

func TestImplementsService(t *testing.T) {
	object := New(nil, nil, nil, nil)
	var expected *service.UserService
	err := AssertThat(object, Implements(expected))
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeleteUserData(t *testing.T) {
	client := ledis.New()
	userTokenDirectory := username_token_directory.New(client)
	userGroupDirectory := username_groupname_directory.New(client)
	tokenUserDirectory := token_username_directory.New(client)
	userDataDirectory := username_data_directory.New(client)
	userService := New(userTokenDirectory, userGroupDirectory, tokenUserDirectory, userDataDirectory)

	user := model.UserName("test")
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
