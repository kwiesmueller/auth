package user_finder

import (
	"fmt"

	"github.com/bborbe/auth/api"
	"github.com/bborbe/log"
)

type userFinder struct {
}

type UserFinder interface {
	FindUserByAuthToken(applicationId api.ApplicationId, authToken api.AuthToken) (*api.User, error)
	IsUserNotFound(err error) bool
}

var (
	NOT_FOUND = fmt.Errorf("user not found")
	logger    = log.DefaultLogger
)

func New() *userFinder {
	return new(userFinder)
}

func (u *userFinder) FindUserByAuthToken(applicationId api.ApplicationId, authToken api.AuthToken) (*api.User, error) {
	return findUser(authToken)
}

func findUser(authToken api.AuthToken) (*api.User, error) {
	logger.Debugf("find user with auth token: %s", authToken)
	if authToken == api.AuthToken("hipchat:130647") {
		user := api.User("bborbe")
		return &user, nil
	}
	if authToken == api.AuthToken("telegram:112230768") {
		user := api.User("bborbe")
		return &user, nil
	}
	return nil, NOT_FOUND
}

func (u *userFinder) IsUserNotFound(err error) bool {
	return err == NOT_FOUND
}
