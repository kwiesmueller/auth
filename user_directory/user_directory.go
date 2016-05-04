package user_directory

import (
	"fmt"

	"github.com/bborbe/auth/api"
	"github.com/bborbe/http/bearer"
	"github.com/bborbe/log"
)

type userDirectory struct {
}

type UserDirectory interface {
	FindUserByAuthToken(applicationName api.ApplicationName, authToken api.AuthToken) (*api.User, error)
	IsUserNotFound(err error) bool
}

var (
	NOT_FOUND = fmt.Errorf("user not found")
	logger    = log.DefaultLogger
)

func New() *userDirectory {
	return new(userDirectory)
}

func (u *userDirectory) FindUserByAuthToken(applicationId api.ApplicationName, authToken api.AuthToken) (*api.User, error) {
	return findUser(authToken)
}

func findUser(authToken api.AuthToken) (*api.User, error) {
	logger.Debugf("find user with auth token: %s", authToken)
	source, id, err := bearer.ParseBearerToken(string(authToken))
	if err != nil {
		return nil, err
	}
	if source == "hipchat" && id == "130647" {
		user := api.User("bborbe")
		return &user, nil
	}
	if source == "telegram" && id == "112230768" {
		user := api.User("bborbe")
		return &user, nil
	}
	return nil, NOT_FOUND
}

func (u *userDirectory) IsUserNotFound(err error) bool {
	return err == NOT_FOUND
}
