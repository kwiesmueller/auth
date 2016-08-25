package user_token_directory

import (
	"testing"

	. "github.com/bborbe/assert"
	"github.com/bborbe/auth/model"
	ledis "github.com/bborbe/ledis/mock"
)

func TestImplementsUserTokenDirectory(t *testing.T) {
	object := New(nil)
	var expected *UserTokenDirectory
	if err := AssertThat(object, Implements(expected)); err != nil {
		t.Fatal(err)
	}
}

func TestList(t *testing.T) {
	userTokenDirectory := New(ledis.New())
	userName := model.UserName("user")
	authToken := model.AuthToken("token")
	var err error
	var userNames []model.UserName
	userNames, err = userTokenDirectory.List()
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(len(userNames), Is(0)); err != nil {
		t.Fatal(err)
	}
	err = userTokenDirectory.Add(userName, authToken)
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	userNames, err = userTokenDirectory.List()
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(userNames, NotNilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(len(userNames), Is(1)); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(userNames[0], Is(userName)); err != nil {
		t.Fatal(err)
	}
}

func TestCreateKey(t *testing.T) {
	key := createKey(model.UserName("test"))
	if err := AssertThat(key, Is("user_token:test")); err != nil {
		t.Fatal(err)
	}
}
