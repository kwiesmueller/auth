package group_checker

import (
	"testing"

	auth_model "github.com/bborbe/auth/model"

	"fmt"
	"os"

	. "github.com/bborbe/assert"
	"github.com/golang/glog"
)

func TestMain(m *testing.M) {
	exit := m.Run()
	glog.Flush()
	os.Exit(exit)
}

func TestImplementsGroupChecker(t *testing.T) {
	object := New(nil)
	var expected *GroupChecker
	if err := AssertThat(object, Implements(expected)); err != nil {
		t.Fatal(err)
	}
}

func TestHasGroupsError(t *testing.T) {
	group := auth_model.GroupName("test")
	token := auth_model.AuthToken("token")
	g := New(func(authToken auth_model.AuthToken, requiredGroups []auth_model.GroupName) (*auth_model.UserName, error) {
		return nil, fmt.Errorf("foo")
	}, group)
	if err := AssertThat(g.HasRequiredGroups(token), Is(false)); err != nil {
		t.Fatal(err)
	}
}

func TestHasGroupsUserNil(t *testing.T) {
	group := auth_model.GroupName("test")
	token := auth_model.AuthToken("token")
	g := New(func(authToken auth_model.AuthToken, requiredGroups []auth_model.GroupName) (*auth_model.UserName, error) {
		return nil, nil
	}, group)
	if err := AssertThat(g.HasRequiredGroups(token), Is(false)); err != nil {
		t.Fatal(err)
	}
}

func TestHasGroupsUserEmpty(t *testing.T) {
	group := auth_model.GroupName("test")
	token := auth_model.AuthToken("token")
	g := New(func(authToken auth_model.AuthToken, requiredGroups []auth_model.GroupName) (*auth_model.UserName, error) {
		user := auth_model.UserName("")
		return &user, nil
	}, group)
	if err := AssertThat(g.HasRequiredGroups(token), Is(false)); err != nil {
		t.Fatal(err)
	}
}

func TestHasGroupsValid(t *testing.T) {
	group := auth_model.GroupName("test")
	token := auth_model.AuthToken("token")
	g := New(func(authToken auth_model.AuthToken, requiredGroups []auth_model.GroupName) (*auth_model.UserName, error) {
		user := auth_model.UserName("tester")
		return &user, nil
	}, group)
	if err := AssertThat(g.HasRequiredGroups(token), Is(true)); err != nil {
		t.Fatal(err)
	}
}
